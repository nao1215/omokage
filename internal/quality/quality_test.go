package quality

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
)

// styleMetrics builds a metrics vector with the localizable, interpretable
// features set explicitly, so a test can control exactly how consistent or how
// out-of-place a document is. The remaining fields are left at zero; the quality
// checks read only size and the localizable features.
func styleMetrics(charCount int, sentenceLength, punctuation, kanji, hiragana, katakana, polite, plain float64) feature.Metrics {
	return feature.Metrics{
		AverageSentenceLength: sentenceLength,
		PunctuationFrequency:  punctuation,
		KanjiRatio:            kanji,
		HiraganaRatio:         hiragana,
		KatakanaRatio:         katakana,
		PoliteEndingRatio:     polite,
		PlainEndingRatio:      plain,
		SentenceCount:         10,
		CharacterCount:        charCount,
	}
}

// consistentDoc is the shape of a long, polite, consistent document. Repeating it
// yields a clean corpus that should raise no findings.
func consistentDoc() feature.Metrics {
	return styleMetrics(300, 30, 0.1, 0.3, 0.6, 0.1, 0.9, 0.05)
}

// corpus builds the distribution and the per-document quality view from a set of
// metrics, the same pairing AssessCorpus consumes, so a test sets up a corpus the
// way training produces one.
func corpus(t *testing.T, metrics ...feature.Metrics) (feature.Distribution, []Document) {
	t.Helper()
	docs := make([]Document, 0, len(metrics))
	for i, m := range metrics {
		docs = append(docs, Document{Name: namef(i), Metrics: m})
	}
	return feature.Aggregate(metrics), docs
}

func namef(i int) string {
	return fmt.Sprintf("doc%d.md", i)
}

func findingByCode(report Report, code string) (Finding, bool) {
	for _, f := range report.Findings {
		if f.Code == code {
			return f, true
		}
	}
	return Finding{}, false
}

func defaultFeatures() config.Features {
	return config.Default("test").Features
}

func TestCleanCorpusHasNoFindings(t *testing.T) {
	t.Parallel()

	metrics := make([]feature.Metrics, 8)
	for i := range metrics {
		metrics[i] = consistentDoc()
	}
	dist, docs := corpus(t, metrics...)

	report := AssessCorpus(dist, docs, defaultFeatures())
	if len(report.Findings) != 0 {
		t.Fatalf("a clean corpus should raise no findings, got %d: %+v", len(report.Findings), report.Findings)
	}
	if report.Reliability() != "good" {
		t.Fatalf("a clean corpus should be good, got %q", report.Reliability())
	}
	if report.DocumentCount != 8 {
		t.Fatalf("document count not reported: %d", report.DocumentCount)
	}
}

func TestFewDocumentsWarnsBelowUsable(t *testing.T) {
	t.Parallel()

	dist, docs := corpus(t, consistentDoc(), consistentDoc())
	report := AssessCorpus(dist, docs, defaultFeatures())

	finding, ok := findingByCode(report, "few_documents")
	if !ok {
		t.Fatalf("two documents should warn about corpus size: %+v", report.Findings)
	}
	if finding.Severity != SeverityWarning {
		t.Fatalf("a 2-document corpus should be a warning, got %v", finding.Severity)
	}
	if report.Reliability() != "weak" {
		t.Fatalf("a 2-document corpus should be weak, got %q", report.Reliability())
	}
	if finding.Action == "" {
		t.Fatal("a finding must carry a next action")
	}
}

func TestFewDocumentsNoticesBelowReliable(t *testing.T) {
	t.Parallel()

	metrics := make([]feature.Metrics, 5)
	for i := range metrics {
		metrics[i] = consistentDoc()
	}
	dist, docs := corpus(t, metrics...)
	report := AssessCorpus(dist, docs, defaultFeatures())

	finding, ok := findingByCode(report, "few_documents")
	if !ok {
		t.Fatalf("five documents should notice (but not warn) about size: %+v", report.Findings)
	}
	if finding.Severity != SeverityNotice {
		t.Fatalf("a 5-document corpus should be a notice, got %v", finding.Severity)
	}
	if report.Reliability() != "fair" {
		t.Fatalf("a notice-only corpus should be fair, got %q", report.Reliability())
	}
}

func TestShortDocumentsWarnWhenMostAreShort(t *testing.T) {
	t.Parallel()

	// Eight documents so size is fine, but all are short, so the short-document
	// finding (not the size finding) is what fires, and at warning level because
	// every document is short.
	short := styleMetrics(40, 12, 0.1, 0.3, 0.6, 0.1, 0.9, 0.05)
	metrics := make([]feature.Metrics, 8)
	for i := range metrics {
		metrics[i] = short
	}
	dist, docs := corpus(t, metrics...)
	report := AssessCorpus(dist, docs, defaultFeatures())

	if _, ok := findingByCode(report, "few_documents"); ok {
		t.Fatalf("eight documents should not warn about size: %+v", report.Findings)
	}
	finding, ok := findingByCode(report, "short_documents")
	if !ok {
		t.Fatalf("a corpus of short documents should be flagged: %+v", report.Findings)
	}
	if finding.Severity != SeverityWarning {
		t.Fatalf("all-short corpus should warn, got %v", finding.Severity)
	}
	// The finding must name the offending files so a user can act on them.
	if !strings.Contains(finding.Detail, namef(0)) {
		t.Fatalf("short-document finding should name the files, got %q", finding.Detail)
	}
}

func TestMixedVoiceFlagsAnInterpretableFeature(t *testing.T) {
	t.Parallel()

	// A register-split corpus: five mostly-polite documents and three plain ones.
	// The uneven split lands the plain-ending ratio's relative spread well above the
	// warning threshold (≈1.3, not on the 1.0 boundary), so the assertion does not
	// hinge on floating-point rounding. The mixed-voice finding should fire and name
	// a register feature, not a noisy second-order feature like a variance.
	polite := styleMetrics(300, 30, 0.1, 0.3, 0.6, 0.1, 0.9, 0.0)
	plain := styleMetrics(300, 30, 0.1, 0.3, 0.6, 0.1, 0.0, 0.9)
	dist, docs := corpus(t, polite, polite, polite, polite, polite, plain, plain, plain)
	report := AssessCorpus(dist, docs, defaultFeatures())

	finding, ok := findingByCode(report, "mixed_voice")
	if !ok {
		t.Fatalf("a register-split corpus should be flagged as mixed: %+v", report.Findings)
	}
	if !strings.Contains(finding.Summary, "ending ratio") {
		t.Fatalf("mixed-voice should name a register feature, got %q", finding.Summary)
	}
	// A wholesale register split swings the feature wider than its own mean, which
	// is a warning, not a mild notice.
	if finding.Severity != SeverityWarning {
		t.Fatalf("a register-split corpus should warn, got %v", finding.Severity)
	}
}

func TestOutlierDocumentIsNamed(t *testing.T) {
	t.Parallel()

	// Many consistent documents and one that differs across every interpretable
	// feature: the odd one out should be named as an outlier. A larger corpus is
	// used on purpose — with only a handful of documents a single outlier dominates
	// the spread it is measured against, exactly as it does on a real corpus, so the
	// outlier signal only becomes clear once the rest of the corpus is the majority.
	base := consistentDoc()
	odd := styleMetrics(300, 30, 0.9, 0.9, 0.05, 0.05, 0.05, 0.9)
	metrics := make([]feature.Metrics, 0, 16)
	for range 15 {
		metrics = append(metrics, base)
	}
	metrics = append(metrics, odd)
	dist, docs := corpus(t, metrics...)
	report := AssessCorpus(dist, docs, defaultFeatures())

	finding, ok := findingByCode(report, "outlier_documents")
	if !ok {
		t.Fatalf("the odd document should be flagged as an outlier: %+v", report.Findings)
	}
	if !strings.Contains(finding.Detail, namef(15)) {
		t.Fatalf("the outlier finding should name the odd document %q, got %q", namef(15), finding.Detail)
	}
}

func TestOutlierDetectedOnSmallCorpus(t *testing.T) {
	t.Parallel()

	// Leave-one-out measurement means the outlier check works even on a small
	// corpus, where a single odd document would otherwise inflate the very spread it
	// is judged against and hide itself. Five steady documents and one that differs
	// across every interpretable feature: the odd one must still be named.
	base := consistentDoc()
	odd := styleMetrics(300, 30, 0.9, 0.9, 0.05, 0.05, 0.05, 0.9)
	dist, docs := corpus(t, base, base, base, base, base, odd)
	report := AssessCorpus(dist, docs, defaultFeatures())

	finding, ok := findingByCode(report, "outlier_documents")
	if !ok {
		t.Fatalf("the odd document should be flagged even on a small corpus: %+v", report.Findings)
	}
	if !strings.Contains(finding.Detail, namef(5)) {
		t.Fatalf("the outlier finding should name the odd document %q, got %q", namef(5), finding.Detail)
	}
}

func TestNoOutlierOnConsistentCorpus(t *testing.T) {
	t.Parallel()

	metrics := make([]feature.Metrics, 6)
	for i := range metrics {
		metrics[i] = consistentDoc()
	}
	dist, docs := corpus(t, metrics...)
	report := AssessCorpus(dist, docs, defaultFeatures())

	if _, ok := findingByCode(report, "outlier_documents"); ok {
		t.Fatalf("a consistent corpus must not invent outliers: %+v", report.Findings)
	}
}

func TestAssessProfileShortAverageAlignsWithDoctor(t *testing.T) {
	t.Parallel()

	// AssessProfile sees only the aggregate, but its rating must not contradict what
	// doctor reported on the same corpus. A corpus whose documents average well
	// under the single-document floor is one doctor calls weak (most documents
	// short), so the profile view warns too — they agree the corpus is weak.
	short := styleMetrics(40, 12, 0.1, 0.3, 0.6, 0.1, 0.9, 0.05)
	metrics := make([]feature.Metrics, 8)
	for i := range metrics {
		metrics[i] = short
	}
	dist := feature.Aggregate(metrics)

	report := AssessProfile(dist, defaultFeatures())
	finding, ok := findingByCode(report, "short_documents")
	if !ok {
		t.Fatalf("a short-on-average corpus should be flagged: %+v", report.Findings)
	}
	if finding.Severity != SeverityWarning {
		t.Fatalf("a very-short-average corpus should warn (matching doctor), got %v", finding.Severity)
	}
	if report.Reliability() != ReliabilityWeak {
		t.Fatalf("the profile rating should be weak to match doctor, got %q", report.Reliability())
	}
	if !strings.Contains(finding.Summary, "average") {
		t.Fatalf("the profile-level finding should speak in averages, got %q", finding.Summary)
	}
}

func TestAssessProfileModerateAverageIsAdvisory(t *testing.T) {
	t.Parallel()

	// A corpus that is merely below the comfortable average, but whose documents are
	// individually long enough, is an advisory notice rather than a warning.
	moderate := styleMetrics(180, 25, 0.1, 0.3, 0.6, 0.1, 0.9, 0.05)
	metrics := make([]feature.Metrics, 8)
	for i := range metrics {
		metrics[i] = moderate
	}
	dist := feature.Aggregate(metrics)

	report := AssessProfile(dist, defaultFeatures())
	finding, ok := findingByCode(report, "short_documents")
	if !ok {
		t.Fatalf("a below-average corpus should still be flagged: %+v", report.Findings)
	}
	if finding.Severity != SeverityNotice {
		t.Fatalf("a moderate-average corpus is a notice, got %v", finding.Severity)
	}
}

func TestReliabilityWeakBeatsFair(t *testing.T) {
	t.Parallel()

	report := Report{Findings: []Finding{
		{Code: "a", Severity: SeverityNotice},
		{Code: "b", Severity: SeverityWarning},
	}}
	if report.Reliability() != "weak" {
		t.Fatalf("any warning makes the corpus weak, got %q", report.Reliability())
	}
}
