package profile

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
)

// Record is a persisted author profile: the learned feature distribution plus
// the metadata describing how it was trained.
type Record struct {
	Author       string
	SourceDir    string
	TrainedAt    time.Time
	FileCount    int
	Distribution feature.Distribution
}

type Comparison struct {
	Similarity  int
	Differences []string
}

// zScoreScale maps the mean absolute z-score to a similarity percentage: a
// target sitting on average this many standard deviations away from the author
// scores 0%. Three standard deviations is treated as "clearly different".
const zScoreScale = 3.0

// driftThreshold is the minimum z-score for a feature to be reported as drift.
// Roughly one standard deviation away from the author's mean.
const driftThreshold = 1.0

// lexicalDistanceScale converts the small absolute gap between two documents'
// function-word frequencies into a [0,1] distance comparable to the ratio
// features, so the diff command weighs lexical drift alongside structure.
const lexicalDistanceScale = 8.0

// Feature categories group the stylistic features by what kind of signal they
// carry and, crucially, by how editable they are. The high-level categories
// (register, structure, script) map to concrete things a person or an LLM can
// change in a draft; the low-level categories (function word, character n-gram)
// are diffuse fingerprints that cannot be edited directly and so are reported as
// supporting detail. levelOf encodes that split.
const (
	categoryRegister     = "register"
	categoryStructure    = "structure"
	categoryScript       = "script"
	categoryFunctionWord = "function-word"
	categoryCharNgram    = "char-ngram"
)

const (
	levelHigh = "high"
	levelLow  = "low"
)

func levelOf(category string) string {
	switch category {
	case categoryFunctionWord, categoryCharNgram:
		return levelLow
	default:
		return levelHigh
	}
}

// featureSpec ties a stylistic feature to its config flag and accessor so that
// scoring, drift reporting, and direct comparison share a single definition.
//
// localizable marks the features that carry meaning at the granularity of a
// single paragraph and that a writer can act on there: register, script balance,
// average sentence length, and punctuation frequency. Each is normalized by a
// denominator large enough to stay stable in one paragraph (script-character
// count, sentence count, character count). The cleared features fall into two
// groups. Document-global features — the variance features (which need many
// sentences or paragraphs) and the layout features (newline, bullet, markdown
// structure) — collapse to a constant per paragraph (a prose paragraph has no
// headings, so its markdown structure frequency is always zero) and would
// otherwise rank every paragraph the same way for the same spurious reason.
// Conjunction frequency is cleared for a different reason: its denominator is the
// paragraph's word-token count, so a single connective spikes it to a large,
// misleading z in a short paragraph. Paragraph localization uses only the
// localizable set; whole-document scoring uses every feature.
type featureSpec struct {
	label       string
	category    string
	enabled     func(config.Features) bool
	value       func(feature.Metrics) float64
	isRatio     bool
	localizable bool
}

var featureSpecs = []featureSpec{
	{"average sentence length", categoryStructure, func(f config.Features) bool { return f.SentenceLength }, func(m feature.Metrics) float64 { return m.AverageSentenceLength }, false, true},
	{"sentence length variance", categoryStructure, func(f config.Features) bool { return f.SentenceLengthVariance }, func(m feature.Metrics) float64 { return m.SentenceLengthVariance }, false, false},
	{"punctuation frequency", categoryStructure, func(f config.Features) bool { return f.PunctuationFrequency }, func(m feature.Metrics) float64 { return m.PunctuationFrequency }, true, true},
	{"newline frequency", categoryStructure, func(f config.Features) bool { return f.NewlineFrequency }, func(m feature.Metrics) float64 { return m.NewlineFrequency }, true, false},
	{"bullet-list frequency", categoryStructure, func(f config.Features) bool { return f.BulletRatio }, func(m feature.Metrics) float64 { return m.BulletRatio }, true, false},
	{"conjunction frequency", categoryStructure, func(f config.Features) bool { return f.ConjunctionFrequency }, func(m feature.Metrics) float64 { return m.ConjunctionFrequency }, true, false},
	{"kanji ratio", categoryScript, func(f config.Features) bool { return f.KanjiRatio }, func(m feature.Metrics) float64 { return m.KanjiRatio }, true, true},
	{"hiragana ratio", categoryScript, func(f config.Features) bool { return f.HiraganaRatio }, func(m feature.Metrics) float64 { return m.HiraganaRatio }, true, true},
	{"katakana ratio", categoryScript, func(f config.Features) bool { return f.KatakanaRatio }, func(m feature.Metrics) float64 { return m.KatakanaRatio }, true, true},
	{"paragraph length variance", categoryStructure, func(f config.Features) bool { return f.ParagraphLengthVariance }, func(m feature.Metrics) float64 { return m.ParagraphLengthVariance }, false, false},
	{"markdown structure frequency", categoryStructure, func(f config.Features) bool { return f.MarkdownStructureDensity }, func(m feature.Metrics) float64 { return m.MarkdownStructureDensity }, true, false},
	{"polite sentence-ending ratio", categoryRegister, func(f config.Features) bool { return f.PoliteEndingRatio }, func(m feature.Metrics) float64 { return m.PoliteEndingRatio }, true, true},
	{"plain sentence-ending ratio", categoryRegister, func(f config.Features) bool { return f.PlainEndingRatio }, func(m feature.Metrics) float64 { return m.PlainEndingRatio }, true, true},
}

// FeatureDrift is the full, per-feature comparison that backs the explain output.
// Unlike the plain "X is higher than reference" string, it carries the numbers an
// editor needs: the target's value, the author's mean and spread, the z-score
// (how far out of the author's range it sits), and where it falls in the fix
// priority. Category/Level let a consumer separate the high-level, editable
// features from the low-level fingerprint.
type FeatureDrift struct {
	Feature   string
	Category  string
	Level     string
	Target    float64
	Mean      float64
	StdDev    float64
	Z         float64
	Direction string
	Priority  int
	// Actionable marks whether the drift exceeds driftThreshold (~1σ) — i.e. it is
	// far enough out of the author's range to be worth correcting, not noise.
	Actionable bool
}

// SegmentDrift localizes drift to a single paragraph and names the one editable,
// paragraph-local feature that strays most there, so a report can say not just
// where to look but what to change. Feature/Z correspond: Z is that feature's
// z-score, so the headline number and the named feature always agree. A paragraph
// is only reported when its strongest localizable drift is actually worth acting
// on, which keeps near-match documents from listing paragraphs with negligible or
// document-global drift.
type SegmentDrift struct {
	Index     int
	Kind      string
	Excerpt   string
	Feature   string
	Category  string
	Z         float64
	Direction string
}

// Explanation is the rich, opt-in result behind `check --explain`/`--format
// json`. Similarity is identical to Score's; Drifts adds the per-feature numbers
// (high-level first, then the capped low-level fingerprint) and Segments points
// at the paragraphs that drift most.
type Explanation struct {
	Similarity int
	Drifts     []FeatureDrift
	Segments   []SegmentDrift
}

// lowLevelExplainLimit caps how many low-level fingerprint drifts the explanation
// reports. The full set runs to hundreds of n-grams; surfacing the top few keeps
// the report readable while still flagging the strongest fingerprint movement.
const lowLevelExplainLimit = 10

// segmentExplainLimit caps how many drifting paragraphs the explanation reports,
// keeping attention on the few worst offenders.
const segmentExplainLimit = 5

// segmentDriftThreshold is the minimum z-score a paragraph's strongest
// localizable feature must reach to be reported. It matches driftThreshold (~1σ),
// the same bar the whole-document report uses to call a feature "actionable", so
// a paragraph is only surfaced when it holds drift genuinely worth editing.
const segmentDriftThreshold = driftThreshold

// minSegmentContentRunes is the minimum non-space character count for a paragraph
// to be localized. A heading or a one-line paragraph has no sentence ending, so
// its register and script ratios collapse to noise (a huge spurious z against an
// author who normally ends sentences politely); document-level features already
// capture whether such short lines are in character. Skipping them keeps the
// localization pointed at real prose.
const minSegmentContentRunes = 30

// Score measures how closely a target document matches a learned author
// distribution. Each feature is standardized against the author's own
// per-document spread (a Burrows's-Delta-style z-score), so a target is judged
// by how far it strays from the author's natural variation rather than from a
// single averaged value.
func Score(reference feature.Distribution, target feature.Metrics, flags config.Features) Comparison {
	drifts := featureDrifts(reference, target, flags)
	if len(drifts) == 0 {
		return Comparison{Similarity: 100, Differences: []string{"no enabled features configured"}}
	}
	return Comparison{
		Similarity:  similarityFromDrifts(drifts),
		Differences: topDifferences(drifts),
	}
}

// Explain produces the rich, editor-facing view of the same comparison Score
// makes. It reuses the identical per-feature z-scores (so the headline similarity
// matches Score exactly), then prioritizes them for editing and, when segments
// are supplied, localizes the drift to the worst paragraphs.
func Explain(reference feature.Distribution, target feature.Metrics, segments []feature.Segment, flags config.Features) Explanation {
	drifts := featureDrifts(reference, target, flags)
	similarity := 100
	if len(drifts) > 0 {
		similarity = similarityFromDrifts(drifts)
	}
	return Explanation{
		Similarity: similarity,
		Drifts:     prioritize(drifts),
		Segments:   locateSegmentDrift(reference, segments, flags),
	}
}

// featureDrifts computes the standardized drift of every active feature: the
// scalar style features, then the function-word fingerprint, then the character
// n-gram fingerprint. It is the shared core of Score and Explain — Score reduces
// it to a similarity and a top-3 list, Explain keeps the full detail. Features
// neither the author nor the target exhibits are dropped, exactly as before, so
// dead features (e.g. Japanese script in an English corpus) do not distort the
// result.
func featureDrifts(reference feature.Distribution, target feature.Metrics, flags config.Features) []FeatureDrift {
	drifts := scalarDrifts(reference, target, flags, func(featureSpec) bool { return true })

	if flags.LexicalFrequency {
		for _, word := range feature.LexicalVocabulary() {
			mean := reference.Mean.LexicalFrequencies[word]
			std := reference.StdDev.LexicalFrequencies[word]
			observed := target.LexicalFrequencies[word]
			if mean == 0 && std == 0 && observed == 0 {
				continue
			}
			drifts = append(drifts, FeatureDrift{
				Feature:   fmt.Sprintf("function word %q", word),
				Category:  categoryFunctionWord,
				Level:     levelLow,
				Target:    observed,
				Mean:      mean,
				StdDev:    std,
				Z:         math.Abs(observed-mean) / lexicalStdFloor(std, mean),
				Direction: direction(mean, observed),
			})
		}
	}

	if flags.CharNgramFrequency {
		for ngram, mean := range reference.Mean.CharNgrams {
			std := reference.StdDev.CharNgrams[ngram]
			observed := target.CharNgrams[ngram]
			if mean == 0 && std == 0 && observed == 0 {
				continue
			}
			drifts = append(drifts, FeatureDrift{
				Feature:   fmt.Sprintf("character n-gram %q", ngram),
				Category:  categoryCharNgram,
				Level:     levelLow,
				Target:    observed,
				Mean:      mean,
				StdDev:    std,
				Z:         math.Abs(observed-mean) / lexicalStdFloor(std, mean),
				Direction: direction(mean, observed),
			})
		}
	}

	return drifts
}

// similarityFromDrifts reduces per-feature drifts to the same similarity Score
// has always produced: each group's mean z is fed to combineDrift, which keeps
// the lexical fingerprint leading, charges register only for its excess, and lets
// the structural remainder nudge. Reconstructing the group means from the drift
// list keeps a single source of truth instead of duplicating the accumulation.
func similarityFromDrifts(drifts []FeatureDrift) int {
	var registerZ, otherZ, functionWordZ, ngramZ float64
	var registerCount, otherCount, functionWordCount, ngramCount int
	for _, drift := range drifts {
		switch drift.Category {
		case categoryRegister:
			registerZ += drift.Z
			registerCount++
		case categoryFunctionWord:
			functionWordZ += drift.Z
			functionWordCount++
		case categoryCharNgram:
			ngramZ += drift.Z
			ngramCount++
		default:
			otherZ += drift.Z
			otherCount++
		}
	}
	meanZ := combineDrift(groupDrift{
		register:     meanOf(registerZ, registerCount),
		other:        meanOf(otherZ, otherCount),
		functionWord: meanOf(functionWordZ, functionWordCount),
		ngram:        meanOf(ngramZ, ngramCount),
	})
	return clampPercent(int(math.Round((1 - meanZ/zScoreScale) * 100)))
}

// topDifferences renders the default `check` output: the three highest-z drifts
// above the threshold, phrased as before. Sorting a copy leaves the caller's
// slice order untouched.
func topDifferences(drifts []FeatureDrift) []string {
	sorted := append([]FeatureDrift(nil), drifts...)
	sort.SliceStable(sorted, func(i int, j int) bool {
		return sorted[i].Z > sorted[j].Z
	})

	differences := make([]string, 0, 3)
	for _, drift := range sorted {
		if drift.Z < driftThreshold {
			continue
		}
		differences = append(differences, fmt.Sprintf("%s is %s than reference", drift.Feature, drift.Direction))
		if len(differences) == 3 {
			break
		}
	}
	if len(differences) == 0 {
		differences = append(differences, "no significant stylistic drift detected")
	}
	return differences
}

// prioritize orders the drifts for an editor: the high-level, editable features
// first (sorted by how far out of range they sit), then the low-level fingerprint
// capped to the strongest few. Every drift gets a 1-based Priority and an
// Actionable flag, so a consumer can fix the highest-priority high-level item
// first and treat the fingerprint as supporting detail.
func prioritize(drifts []FeatureDrift) []FeatureDrift {
	high := make([]FeatureDrift, 0, len(drifts))
	low := make([]FeatureDrift, 0, len(drifts))
	for _, drift := range drifts {
		if drift.Level == levelHigh {
			high = append(high, drift)
		} else {
			low = append(low, drift)
		}
	}
	byZ := func(s []FeatureDrift) func(i, j int) bool {
		return func(i, j int) bool { return s[i].Z > s[j].Z }
	}
	sort.SliceStable(high, byZ(high))
	sort.SliceStable(low, byZ(low))
	if len(low) > lowLevelExplainLimit {
		low = low[:lowLevelExplainLimit]
	}

	ordered := make([]FeatureDrift, 0, len(high)+len(low))
	ordered = append(ordered, high...)
	ordered = append(ordered, low...)
	for i := range ordered {
		ordered[i].Priority = i + 1
		ordered[i].Actionable = ordered[i].Z >= driftThreshold
	}
	return ordered
}

// scalarDrifts computes the standardized drift of the scalar style features that
// pass the include predicate. featureDrifts uses it for the full set (and then
// appends the lexical and n-gram fingerprints); paragraph localization uses it
// for the localizable subset only, which also avoids recomputing the hundreds of
// fingerprint z-scores per paragraph. Features neither the author nor the target
// exhibits are dropped so dead features do not distort the result.
func scalarDrifts(reference feature.Distribution, target feature.Metrics, flags config.Features, include func(featureSpec) bool) []FeatureDrift {
	drifts := make([]FeatureDrift, 0, len(featureSpecs))
	for _, spec := range featureSpecs {
		if !spec.enabled(flags) || !include(spec) {
			continue
		}
		mean := spec.value(reference.Mean)
		std := spec.value(reference.StdDev)
		observed := spec.value(target)
		if mean == 0 && std == 0 && observed == 0 {
			continue
		}
		drifts = append(drifts, FeatureDrift{
			Feature:   spec.label,
			Category:  spec.category,
			Level:     levelOf(spec.category),
			Target:    observed,
			Mean:      mean,
			StdDev:    std,
			Z:         math.Abs(observed-mean) / stdFloor(std, mean, spec.isRatio),
			Direction: direction(mean, observed),
		})
	}
	return drifts
}

// locateSegmentDrift names, for each paragraph, the single editable feature that
// strays most there, and returns the worst few. It scores paragraphs on the
// localizable feature subset only: document-global features (layout, variance)
// are constant per paragraph and would otherwise rank every paragraph the same
// way for the same spurious reason. A paragraph is reported only when its
// strongest localizable drift clears segmentDriftThreshold, so a near-match
// document — where nothing local stands out — yields an empty list rather than
// misleading guidance. This is the one genuinely extra computation in the explain
// path (a feature extraction per paragraph), so callers pass segments in explain
// mode only.
func locateSegmentDrift(reference feature.Distribution, segments []feature.Segment, flags config.Features) []SegmentDrift {
	if len(segments) == 0 {
		return nil
	}
	out := make([]SegmentDrift, 0, len(segments))
	for _, segment := range segments {
		if segment.Metrics.CharacterCount < minSegmentContentRunes {
			continue
		}
		var top FeatureDrift
		found := false
		for _, drift := range scalarDrifts(reference, segment.Metrics, flags, localizableSpec) {
			if !found || drift.Z > top.Z {
				top = drift
				found = true
			}
		}
		if !found || top.Z < segmentDriftThreshold {
			continue
		}
		out = append(out, SegmentDrift{
			Index:     segment.Index,
			Kind:      segment.Kind,
			Excerpt:   excerpt(segment.Text),
			Feature:   top.Feature,
			Category:  top.Category,
			Z:         top.Z,
			Direction: top.Direction,
		})
	}
	sort.SliceStable(out, func(i int, j int) bool {
		return out[i].Z > out[j].Z
	})
	if len(out) > segmentExplainLimit {
		out = out[:segmentExplainLimit]
	}
	return out
}

// localizableSpec selects the features that are meaningful and editable at the
// granularity of a single paragraph. See featureSpec.localizable.
func localizableSpec(spec featureSpec) bool {
	return spec.localizable
}

// excerpt returns a short, single-line preview of a paragraph for a report,
// collapsing internal whitespace and truncating with an ellipsis.
func excerpt(text string) string {
	const maxRunes = 50
	collapsed := strings.Join(strings.Fields(text), " ")
	runes := []rune(collapsed)
	if len(runes) <= maxRunes {
		return collapsed
	}
	return string(runes[:maxRunes]) + "…"
}

// Compare measures the stylistic closeness of two individual documents. Unlike
// Score it has no learned distribution to standardize against, so it falls back
// to a relative distance per feature. It backs the `diff` command.
func Compare(reference feature.Metrics, target feature.Metrics, flags config.Features) Comparison {
	type scored struct {
		label     string
		distance  float64
		direction string
	}

	results := make([]scored, 0, len(featureSpecs))
	registerDist := 0.0
	registerCount := 0
	otherDist := 0.0
	otherCount := 0
	for _, spec := range featureSpecs {
		if !spec.enabled(flags) {
			continue
		}
		left := spec.value(reference)
		right := spec.value(target)
		// A feature absent from both documents carries no stylistic signal — e.g.
		// the Japanese script and register features for two English texts. Counting
		// it as a perfect match would inflate the similarity, so it is dropped, the
		// same way Score skips degenerate features.
		if left == 0 && right == 0 {
			continue
		}
		distance := relativeDistance(left, right)
		if spec.isRatio {
			distance = math.Min(1, math.Abs(left-right))
		}
		if registerLabels[spec.label] {
			registerDist += distance
			registerCount++
		} else {
			otherDist += distance
			otherCount++
		}
		results = append(results, scored{
			label:     spec.label,
			distance:  distance,
			direction: direction(left, right),
		})
	}

	functionWordDist := 0.0
	functionWordCount := 0
	if flags.LexicalFrequency {
		for _, word := range feature.LexicalVocabulary() {
			left := reference.LexicalFrequencies[word]
			right := target.LexicalFrequencies[word]
			if left == 0 && right == 0 {
				continue
			}
			// Function-word frequencies are tiny absolute numbers, so the raw gap is
			// scaled to land on the same [0,1] distance footing as the ratio features.
			distance := math.Min(1, math.Abs(left-right)*lexicalDistanceScale)
			functionWordDist += distance
			functionWordCount++
			results = append(results, scored{
				label:     fmt.Sprintf("function word %q", word),
				distance:  distance,
				direction: direction(left, right),
			})
		}
	}

	ngramDist := 0.0
	ngramCount := 0
	if flags.CharNgramFrequency {
		seen := make(map[string]struct{}, len(reference.CharNgrams)+len(target.CharNgrams))
		for ngram := range reference.CharNgrams {
			seen[ngram] = struct{}{}
		}
		for ngram := range target.CharNgrams {
			seen[ngram] = struct{}{}
		}
		for ngram := range seen {
			left := reference.CharNgrams[ngram]
			right := target.CharNgrams[ngram]
			if left == 0 && right == 0 {
				continue
			}
			distance := math.Min(1, math.Abs(left-right)*lexicalDistanceScale)
			ngramDist += distance
			ngramCount++
			results = append(results, scored{
				label:     fmt.Sprintf("character n-gram %q", ngram),
				distance:  distance,
				direction: direction(left, right),
			})
		}
	}

	if registerCount+otherCount+functionWordCount+ngramCount == 0 {
		return Comparison{Similarity: 100, Differences: []string{"no enabled features configured"}}
	}

	// Combine the groups the same way Score does so the diff stays consistent with
	// check: averaging within each group first keeps the many character n-grams
	// from drowning out a register difference between the two documents.
	drift := combineCompareDrift(groupDrift{
		register:     meanOf(registerDist, registerCount),
		other:        meanOf(otherDist, otherCount),
		functionWord: meanOf(functionWordDist, functionWordCount),
		ngram:        meanOf(ngramDist, ngramCount),
	})
	similarity := clampPercent(int(math.Round((1 - drift) * 100)))

	sort.SliceStable(results, func(i int, j int) bool {
		return results[i].distance > results[j].distance
	})

	differences := make([]string, 0, 3)
	for _, result := range results {
		if result.distance < 0.02 {
			continue
		}
		differences = append(differences, fmt.Sprintf("%s is %s than reference", result.label, result.direction))
		if len(differences) == 3 {
			break
		}
	}
	if len(differences) == 0 {
		differences = append(differences, "no significant stylistic drift detected")
	}

	return Comparison{Similarity: similarity, Differences: differences}
}

// stdFloor keeps the standardization stable when a feature barely varies across
// the training corpus (or the profile holds a single document). Without a floor
// a near-zero standard deviation would turn negligible differences into huge
// z-scores.
func stdFloor(std float64, mean float64, isRatio bool) float64 {
	if isRatio {
		return math.Max(std, 0.02)
	}
	return math.Max(std, math.Max(0.1*math.Abs(mean), 1.0))
}

// lexWeight is the share of the similarity drift attributed to the lexical
// fingerprint when both structural and lexical features are present. Averaging
// every feature together would let dozens of low-signal function words dilute a
// strong structural marker (such as a register shift), so the two groups are
// averaged independently and then blended. This keeps register/script detection
// intact while letting the lexical fingerprint separate same-register authors.
// registerLabels marks the sentence-ending features that form the register
// group. They are scored separately from the other structural features so a
// register shift is not diluted, and so they do not add noise to same-register
// comparisons. See combineDrift.
var registerLabels = map[string]bool{
	"polite sentence-ending ratio": true,
	"plain sentence-ending ratio":  true,
}

// registerWeight and otherStructWeight scale how much the register group and
// the remaining structural features add on top of the lexical fingerprint,
// which is the primary, language-independent authorship signal. The lexical
// group separates same-register and English authors; register is kept as a
// clean, undiluted term so a large register shift (an LLM imitation written in
// the opposite register, or cross-language text) still drives the score down,
// while an author's own mild register variation only nudges it. The structural
// remainder barely separates authors on its own, so it nudges least.
const (
	registerWeight    = 1.0
	otherStructWeight = 0.05
	// registerTolerance is the register z-score an author may reach through their
	// own variation (e.g. nao writes mostly 敬体 but slips into 常体 in some posts)
	// before it counts as a register shift. Only the excess above this hinge is
	// charged, so a genuine same-register document is untouched while a wholesale
	// register flip — an LLM imitation in the opposite register, or cross-language
	// text whose register features collapse to zero — is penalized sharply.
	registerTolerance = 2.5
)

// groupDrift holds the mean z-score of each feature group for a single
// comparison. A zero mean means the group had no active features. Function words
// and character n-grams are kept apart so the larger n-gram vocabulary cannot
// outweigh the function-word signal; combineDrift averages whichever of the two
// are present into a single lexical contribution.
type groupDrift struct {
	register     float64
	other        float64
	functionWord float64
	ngram        float64
}

// combineDrift fuses the feature groups into a single drift figure. The lexical
// fingerprint leads, since it separates same-register and English authors; it is
// the equal-weight mean of the function-word and character-n-gram sub-signals so
// that neither the ~150 function words nor the ~400 n-grams dominate by sheer
// count. The register group is added only for its excess above registerTolerance,
// so an author's own mild register variation is ignored while a wholesale
// register flip (an LLM imitation, cross-language text) is charged sharply. The
// noisy structural remainder only nudges the result.
func combineDrift(g groupDrift) float64 {
	lexical := meanOfPresent(g.functionWord, g.ngram)
	registerExcess := g.register - registerTolerance
	if registerExcess < 0 {
		registerExcess = 0
	}
	return lexical + registerWeight*registerExcess + otherStructWeight*g.other
}

// registerCompareWeight is how much a register difference between two documents
// contributes to the diff drift. Unlike Score there is no learned distribution,
// so there is no tolerance hinge: a register difference between two specific
// documents is taken at face value and weighted heavily, since it is one of the
// clearest stylistic divergences a reader notices.
const (
	registerCompareWeight = 0.6
	// otherCompareWeight is larger than otherStructWeight: a direct document
	// comparison has no authorship distribution to lean on, so structural
	// differences (layout, sentence length, punctuation) are themselves a
	// meaningful part of how two documents differ, not just noise.
	otherCompareWeight = 0.34
)

// combineCompareDrift mirrors combineDrift for the distribution-free diff path.
// The lexical fingerprint leads (the equal-weight mean of the function-word and
// n-gram distances), a register difference is added with a fixed weight, and the
// remaining structural features contribute a moderate share.
func combineCompareDrift(g groupDrift) float64 {
	lexical := meanOfPresent(g.functionWord, g.ngram)
	return lexical + registerCompareWeight*g.register + otherCompareWeight*g.other
}

// meanOfPresent averages the sub-signals that are actually present. A sub-signal
// of zero means its group had no active features (e.g. the function-word group
// for a target sharing no vocabulary, or either group when disabled), so it is
// excluded from the average rather than dragging it toward zero.
func meanOfPresent(values ...float64) float64 {
	sum := 0.0
	count := 0
	for _, value := range values {
		if value > 0 {
			sum += value
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// meanOf reduces a running z-score sum to its mean, returning zero for an empty
// group so an absent group contributes nothing to combineDrift.
func meanOf(sum float64, count int) float64 {
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// lexicalStdFloor stabilizes the standardization for function-word frequencies,
// which are small numbers (a common particle sits near 0.05). A fixed ratio
// floor like stdFloor's 0.02 would swamp them, so the floor scales with the
// word's own mean and keeps only a tiny absolute guard for near-constant words.
func lexicalStdFloor(std float64, mean float64) float64 {
	return math.Max(std, math.Max(0.12*mean, 0.0015))
}

func direction(reference float64, target float64) string {
	if target > reference {
		return "higher"
	}
	return "lower"
}

func relativeDistance(left float64, right float64) float64 {
	denominator := math.Max(math.Max(math.Abs(left), math.Abs(right)), 1)
	return math.Min(1, math.Abs(left-right)/denominator)
}

func clampPercent(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}
