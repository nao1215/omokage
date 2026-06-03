// Package quality assesses how reliable comparisons against a corpus will be,
// so a user can tell whether their training material is good enough before they
// lean on the scores it produces.
//
// It judges three things only: is there enough material, are the individual
// documents long enough to measure, and is the voice consistent (or does the
// corpus mix different kinds of writing). It does NOT judge whether the writing
// is good, correct, or original — omokage compares style, not merit, and the
// quality report keeps that promise. Every finding names a concrete next action
// so the report helps a user curate a corpus rather than only describing it.
//
// The thresholds and messages live here so the new `doctor` command, the
// post-training notes, and the `show --format json` summary all speak with one
// voice. The stylometric judgements (how far a document strays, how widely a
// feature varies) are delegated to internal/profile, which owns the scoring, so
// quality stays policy and wording on top of a single source of truth.
package quality

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/profile"
)

// Severity ranks a finding by how much it should worry the user. Notice is "this
// is worth knowing"; Warning is "the scores from this corpus are likely to
// mislead until you act". There is no error level: a corpus that trains at all is
// never rejected here — train enforces the one hard rule (non-empty) itself.
type Severity int

const (
	// SeverityNotice is advisory: the corpus works, but could be better.
	SeverityNotice Severity = iota + 1
	// SeverityWarning means the corpus is thin or inconsistent enough that scores
	// will be noisy or misleading.
	SeverityWarning
)

// String renders a severity for human and JSON output.
func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityNotice:
		return "notice"
	default:
		return "notice"
	}
}

// Finding is one observation about a corpus, written to be acted on: Summary says
// what is true, Detail gives the specifics (counts, file names), and Action says
// what to do about it.
type Finding struct {
	Code     string
	Severity Severity
	Summary  string
	Detail   string
	Action   string
}

// Document pairs a corpus file's display name with its extracted features, so a
// per-document finding can name the files involved.
type Document struct {
	Name    string
	Metrics feature.Metrics
}

// Report is the assessment of a corpus: its size, and the findings (empty when
// nothing stands out). Reliability summarizes the findings into one word.
type Report struct {
	DocumentCount             int
	SentenceCount             int
	CharacterCount            int
	AverageDocumentCharacters int
	Findings                  []Finding
}

// Reliability ratings, the single-word headline summarizing a Report. They
// describe how dependable comparisons against the corpus are, not the quality of
// the writing.
const (
	ReliabilityGood = "good"
	ReliabilityFair = "fair"
	ReliabilityWeak = "weak"
)

// Reliability collapses the findings into a single rating: weak if any warning
// fired, fair if only notices did, good if the corpus looks clean.
func (r Report) Reliability() string {
	level := ReliabilityGood
	for _, f := range r.Findings {
		switch f.Severity {
		case SeverityWarning:
			return ReliabilityWeak
		case SeverityNotice:
			level = ReliabilityFair
		}
	}
	return level
}

const (
	// minReliableDocuments is the count below which scores start to wobble: with
	// fewer samples the measured spread is a poor estimate of the author's real
	// range. It is a notice, not a hard floor.
	minReliableDocuments = 8
	// minUsableDocuments is the count below which the spread is barely an estimate
	// at all (one or two documents). It is a warning.
	minUsableDocuments = 3
	// shortDocumentChars is the non-space character count below which a single
	// document is too short to measure stable per-document features.
	shortDocumentChars = 150
	// shortAverageChars is the same idea at the corpus level, used when only the
	// aggregate is available (a stored profile, where per-document lengths are not
	// kept).
	shortAverageChars = 200
	// manyShortFraction is the share of short documents past which the shortness is
	// a property of the corpus, not a couple of stragglers, and so warrants a
	// warning rather than a notice.
	manyShortFraction = 0.5
	// mixedSpreadRatio is the relative spread (std / |mean|) past which a
	// high-level feature swings widely enough to suggest the corpus mixes different
	// kinds of writing.
	mixedSpreadRatio = 0.6
	// strongSpreadRatio promotes a mixed-voice notice to a warning: a feature that
	// varies as much as its own mean is no longer mild author variation.
	strongSpreadRatio = 1.0
	// outlierZ is the mean high-level z-score past which a single document reads
	// far enough from the rest of the corpus to flag for review.
	outlierZ = 2.5
	// minDocsForOutliers is the smallest corpus where calling one document an
	// outlier is meaningful; below it "outlier" is just "small sample". Because the
	// outlier check measures each document against the others (leave-one-out), this
	// stays meaningful even at the floor rather than needing a large corpus.
	minDocsForOutliers = 4
	// maxNamedDocuments caps how many file names a finding lists before summarizing
	// the rest, so the report stays scannable on a large corpus.
	maxNamedDocuments = 5
)

// AssessCorpus assesses a corpus from its aggregate distribution and the
// per-document metrics that produced it. It runs every check, including the ones
// that need individual documents (short files, outliers). It backs `doctor` and
// the post-training notes.
func AssessCorpus(dist feature.Distribution, docs []Document, flags config.Features) Report {
	report := baseReport(dist)
	appendFinding(&report, checkDocumentCount(dist.DocumentCount))
	appendFinding(&report, checkShortDocuments(docs))
	appendFinding(&report, checkMixedVoice(dist, flags))
	appendFinding(&report, checkOutliers(docs, flags))
	return report
}

// AssessProfile assesses a corpus from a stored profile's distribution alone. The
// training text and per-document metrics are not kept, so it runs only the checks
// the aggregate can answer (size, average length, consistency); it cannot name a
// short file or an outlier. It backs `show --format json`.
func AssessProfile(dist feature.Distribution, flags config.Features) Report {
	report := baseReport(dist)
	appendFinding(&report, checkDocumentCount(dist.DocumentCount))
	appendFinding(&report, checkAverageLength(dist))
	appendFinding(&report, checkMixedVoice(dist, flags))
	return report
}

// baseReport fills in the size fields shared by every assessment (counts and the
// average document length), leaving Findings for the checks to append.
func baseReport(dist feature.Distribution) Report {
	average := 0
	if dist.DocumentCount > 0 {
		average = dist.CharacterCount / dist.DocumentCount
	}
	return Report{
		DocumentCount:             dist.DocumentCount,
		SentenceCount:             dist.SentenceCount,
		CharacterCount:            dist.CharacterCount,
		AverageDocumentCharacters: average,
	}
}

// appendFinding adds a check's result to the report, ignoring the nil a check
// returns when it has nothing to report.
func appendFinding(report *Report, finding *Finding) {
	if finding != nil {
		report.Findings = append(report.Findings, *finding)
	}
}

// checkDocumentCount flags a corpus that is too small for the spread to be a
// trustworthy estimate: a warning below minUsableDocuments, a notice below
// minReliableDocuments, nothing above.
func checkDocumentCount(n int) *Finding {
	switch {
	case n <= 0:
		return nil
	case n < minUsableDocuments:
		return &Finding{
			Code:     "few_documents",
			Severity: SeverityWarning,
			Summary:  fmt.Sprintf("Only %s. The measured spread is barely an estimate, so scores will be noisy.", Documents(n)),
			Action:   fmt.Sprintf("Add more samples of this voice; %d or more documents give steadier scores.", minReliableDocuments),
		}
	case n < minReliableDocuments:
		return &Finding{
			Code:     "few_documents",
			Severity: SeverityNotice,
			Summary:  fmt.Sprintf("Only %s. A few more would tighten the spread and steady the scores.", Documents(n)),
			Action:   fmt.Sprintf("Aim for %d or more documents of this voice.", minReliableDocuments),
		}
	default:
		return nil
	}
}

// checkShortDocuments flags documents too short to measure stable per-document
// features, naming them. It warns when short documents are at least half the
// corpus (the shortness is a property of the corpus, not a few stragglers) and
// notices otherwise.
func checkShortDocuments(docs []Document) *Finding {
	short := make([]string, 0, len(docs))
	for _, doc := range docs {
		if doc.Metrics.CharacterCount < shortDocumentChars {
			short = append(short, doc.Name)
		}
	}
	if len(short) == 0 || len(docs) == 0 {
		return nil
	}
	severity := SeverityNotice
	if float64(len(short))/float64(len(docs)) >= manyShortFraction {
		severity = SeverityWarning
	}
	return &Finding{
		Code:     "short_documents",
		Severity: severity,
		Summary:  fmt.Sprintf("%d of %d documents are short (under %d characters).", len(short), len(docs), shortDocumentChars),
		Detail:   nameList(short),
		Action:   "Short samples make per-document features (sentence length, register) jumpy; prefer samples of a few paragraphs, or merge very short notes.",
	}
}

// checkAverageLength is the short-document check for a stored profile, which
// keeps only the aggregate. It stands in for checkShortDocuments when the
// per-document lengths are gone, judging by the corpus average instead.
func checkAverageLength(dist feature.Distribution) *Finding {
	if dist.DocumentCount == 0 {
		return nil
	}
	average := dist.CharacterCount / dist.DocumentCount
	if average >= shortAverageChars {
		return nil
	}
	// The profile keeps only the aggregate, so this stands in for the per-document
	// short-file check doctor runs. An average below the single-document floor means
	// most documents are short, which is the warning doctor would raise; a merely
	// below-average corpus is an advisory notice. Tiering the severity here keeps a
	// profile's reliability rating aligned with what doctor reported on the corpus.
	severity := SeverityNotice
	if average < shortDocumentChars {
		severity = SeverityWarning
	}
	return &Finding{
		Code:     "short_documents",
		Severity: severity,
		Summary:  fmt.Sprintf("Documents average %d characters, which is short for steady per-document features.", average),
		Action:   "Prefer samples of a few paragraphs each, or merge very short notes, then retrain.",
	}
}

// checkMixedVoice flags a corpus whose interpretable features swing widely enough
// to suggest it mixes different kinds of writing, naming the feature it disagrees
// on most. A feature that varies as much as its own mean (strongSpreadRatio)
// promotes the notice to a warning.
func checkMixedVoice(dist feature.Distribution, flags config.Features) *Finding {
	// Spread is only meaningful once a few documents back it; with one or two the
	// document-count finding already carries the message.
	if dist.DocumentCount < minUsableDocuments {
		return nil
	}
	wide := make([]profile.FeatureSpread, 0)
	strong := false
	for _, spread := range profile.HighLevelSpreads(dist, flags) {
		if spread.RelativeSpread < mixedSpreadRatio {
			continue
		}
		if spread.RelativeSpread >= strongSpreadRatio {
			strong = true
		}
		wide = append(wide, spread)
	}
	if len(wide) == 0 {
		return nil
	}
	severity := SeverityNotice
	if strong {
		severity = SeverityWarning
	}
	if len(wide) > 2 {
		wide = wide[:2]
	}
	named := wide
	features := make([]string, 0, len(named))
	details := make([]string, 0, len(named))
	for _, spread := range named {
		features = append(features, spread.Feature)
		details = append(details, fmt.Sprintf("%s: mean %.3f, spread %.3f (relative %.1f)",
			spread.Feature, spread.Mean, spread.StdDev, spread.RelativeSpread))
	}
	return &Finding{
		Code:     "mixed_voice",
		Severity: severity,
		Summary:  fmt.Sprintf("Your samples vary widely in %s.", strings.Join(features, " and ")),
		Detail:   strings.Join(details, "; "),
		Action:   "If these samples mix different kinds of writing (formal posts and casual notes, say), train a separate author for each so every voice stays consistent.",
	}
}

// checkOutliers flags documents that read differently from the rest of the
// corpus, naming them with how far out they sit. It measures each document
// leave-one-out (against the others) so a lone outlier does not mask itself by
// widening the spread it is judged against.
func checkOutliers(docs []Document, flags config.Features) *Finding {
	if len(docs) < minDocsForOutliers {
		return nil
	}
	metrics := make([]feature.Metrics, len(docs))
	for i, doc := range docs {
		metrics[i] = doc.Metrics
	}
	// Leave-one-out: each document is measured against the distribution of the
	// others, so a lone outlier does not blunt its own signal by widening the
	// spread it is judged against.
	divergences := profile.LeaveOneOutDivergences(metrics, flags)

	type outlier struct {
		name string
		z    float64
	}
	outliers := make([]outlier, 0)
	for i, doc := range docs {
		if divergences[i] >= outlierZ {
			outliers = append(outliers, outlier{name: doc.Name, z: divergences[i]})
		}
	}
	if len(outliers) == 0 {
		return nil
	}
	sort.SliceStable(outliers, func(i int, j int) bool {
		return outliers[i].z > outliers[j].z
	})
	details := make([]string, 0, len(outliers))
	for _, o := range outliers {
		details = append(details, fmt.Sprintf("%s (%.1fσ)", o.name, o.z))
	}
	return &Finding{
		Code:     "outlier_documents",
		Severity: SeverityNotice,
		Summary:  fmt.Sprintf("%d of %d documents read differently from the rest.", len(outliers), len(docs)),
		Detail:   nameList(details),
		Action:   "Check the named files belong to this voice; drop the ones that don't and retrain.",
	}
}

// Documents renders a document count with the correct singular/plural noun. It is
// shared by the finding summaries here and the size line the doctor command
// prints, so the two never disagree on wording.
func Documents(n int) string {
	if n == 1 {
		return "1 document"
	}
	return fmt.Sprintf("%d documents", n)
}

// nameList joins names for a finding's Detail, capping the visible names so a
// large corpus does not print a wall of file paths.
func nameList(names []string) string {
	if len(names) <= maxNamedDocuments {
		return strings.Join(names, ", ")
	}
	shown := names[:maxNamedDocuments]
	return fmt.Sprintf("%s, and %d more", strings.Join(shown, ", "), len(names)-maxNamedDocuments)
}
