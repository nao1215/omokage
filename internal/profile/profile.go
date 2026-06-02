package profile

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/nao1215/dyer/internal/config"
	"github.com/nao1215/dyer/internal/feature"
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

// featureSpec ties a stylistic feature to its config flag and accessor so that
// scoring, drift reporting, and direct comparison share a single definition.
type featureSpec struct {
	label   string
	enabled func(config.Features) bool
	value   func(feature.Metrics) float64
	isRatio bool
}

var featureSpecs = []featureSpec{
	{"average sentence length", func(f config.Features) bool { return f.SentenceLength }, func(m feature.Metrics) float64 { return m.AverageSentenceLength }, false},
	{"sentence length variance", func(f config.Features) bool { return f.SentenceLengthVariance }, func(m feature.Metrics) float64 { return m.SentenceLengthVariance }, false},
	{"punctuation frequency", func(f config.Features) bool { return f.PunctuationFrequency }, func(m feature.Metrics) float64 { return m.PunctuationFrequency }, true},
	{"newline frequency", func(f config.Features) bool { return f.NewlineFrequency }, func(m feature.Metrics) float64 { return m.NewlineFrequency }, true},
	{"bullet-list frequency", func(f config.Features) bool { return f.BulletRatio }, func(m feature.Metrics) float64 { return m.BulletRatio }, true},
	{"conjunction frequency", func(f config.Features) bool { return f.ConjunctionFrequency }, func(m feature.Metrics) float64 { return m.ConjunctionFrequency }, true},
	{"kanji ratio", func(f config.Features) bool { return f.KanjiRatio }, func(m feature.Metrics) float64 { return m.KanjiRatio }, true},
	{"hiragana ratio", func(f config.Features) bool { return f.HiraganaRatio }, func(m feature.Metrics) float64 { return m.HiraganaRatio }, true},
	{"katakana ratio", func(f config.Features) bool { return f.KatakanaRatio }, func(m feature.Metrics) float64 { return m.KatakanaRatio }, true},
	{"paragraph length variance", func(f config.Features) bool { return f.ParagraphLengthVariance }, func(m feature.Metrics) float64 { return m.ParagraphLengthVariance }, false},
	{"markdown structure frequency", func(f config.Features) bool { return f.MarkdownStructureDensity }, func(m feature.Metrics) float64 { return m.MarkdownStructureDensity }, true},
	{"polite sentence-ending ratio", func(f config.Features) bool { return f.PoliteEndingRatio }, func(m feature.Metrics) float64 { return m.PoliteEndingRatio }, true},
	{"plain sentence-ending ratio", func(f config.Features) bool { return f.PlainEndingRatio }, func(m feature.Metrics) float64 { return m.PlainEndingRatio }, true},
}

// Score measures how closely a target document matches a learned author
// distribution. Each feature is standardized against the author's own
// per-document spread (a Burrows's-Delta-style z-score), so a target is judged
// by how far it strays from the author's natural variation rather than from a
// single averaged value.
func Score(reference feature.Distribution, target feature.Metrics, flags config.Features) Comparison {
	type scored struct {
		label     string
		z         float64
		direction string
	}

	results := make([]scored, 0, len(featureSpecs))
	registerZ := 0.0
	registerCount := 0
	otherZ := 0.0
	otherCount := 0
	for _, spec := range featureSpecs {
		if !spec.enabled(flags) {
			continue
		}
		mean := spec.value(reference.Mean)
		std := spec.value(reference.StdDev)
		observed := spec.value(target)
		// A feature the author never exhibits (mean and spread both zero) carries
		// no stylistic signal when the target also lacks it — e.g. the Japanese
		// script and sentence-ending features are identically zero for an English
		// corpus. Counting them as a perfect match would inflate every score, so
		// they are dropped. A target that *does* exhibit the feature still counts
		// as drift via the floor below.
		if mean == 0 && std == 0 && observed == 0 {
			continue
		}
		z := math.Abs(observed-mean) / stdFloor(std, mean, spec.isRatio)
		// The register markers (polite/plain endings) are kept in their own group:
		// a register shift is a decisive impostor signal (an LLM imitation), but
		// averaged in with the other structural features it gets diluted, and those
		// features add noise to same-author comparisons. See combineDrift.
		if registerLabels[spec.label] {
			registerZ += z
			registerCount++
		} else {
			otherZ += z
			otherCount++
		}
		results = append(results, scored{
			label:     spec.label,
			z:         z,
			direction: direction(mean, observed),
		})
	}

	functionWordZ := 0.0
	functionWordCount := 0
	if flags.LexicalFrequency {
		for _, word := range feature.LexicalVocabulary() {
			mean := reference.Mean.LexicalFrequencies[word]
			std := reference.StdDev.LexicalFrequencies[word]
			observed := target.LexicalFrequencies[word]
			// A function word neither the author nor the target ever uses (e.g.
			// English words in a Japanese-only corpus) carries no signal, so it is
			// dropped exactly like the degenerate scalar features above.
			if mean == 0 && std == 0 && observed == 0 {
				continue
			}
			z := math.Abs(observed-mean) / lexicalStdFloor(std, mean)
			functionWordZ += z
			functionWordCount++
			results = append(results, scored{
				label:     fmt.Sprintf("function word %q", word),
				z:         z,
				direction: direction(mean, observed),
			})
		}
	}

	ngramZ := 0.0
	ngramCount := 0
	if flags.CharNgramFrequency {
		for ngram, mean := range reference.Mean.CharNgrams {
			std := reference.StdDev.CharNgrams[ngram]
			observed := target.CharNgrams[ngram]
			if mean == 0 && std == 0 && observed == 0 {
				continue
			}
			z := math.Abs(observed-mean) / lexicalStdFloor(std, mean)
			ngramZ += z
			ngramCount++
			results = append(results, scored{
				label:     fmt.Sprintf("character n-gram %q", ngram),
				z:         z,
				direction: direction(mean, observed),
			})
		}
	}

	if registerCount+otherCount+functionWordCount+ngramCount == 0 {
		return Comparison{Similarity: 100, Differences: []string{"no enabled features configured"}}
	}

	meanZ := combineDrift(groupDrift{
		register:     meanOf(registerZ, registerCount),
		other:        meanOf(otherZ, otherCount),
		functionWord: meanOf(functionWordZ, functionWordCount),
		ngram:        meanOf(ngramZ, ngramCount),
	})
	similarity := clampPercent(int(math.Round((1 - meanZ/zScoreScale) * 100)))

	sort.SliceStable(results, func(i int, j int) bool {
		return results[i].z > results[j].z
	})

	differences := make([]string, 0, 3)
	for _, result := range results {
		if result.z < driftThreshold {
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
	total := 0.0
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
		total += distance
		results = append(results, scored{
			label:     spec.label,
			distance:  distance,
			direction: direction(left, right),
		})
	}

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
			total += distance
			results = append(results, scored{
				label:     fmt.Sprintf("function word %q", word),
				distance:  distance,
				direction: direction(left, right),
			})
		}
	}

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
			total += distance
			results = append(results, scored{
				label:     fmt.Sprintf("character n-gram %q", ngram),
				distance:  distance,
				direction: direction(left, right),
			})
		}
	}

	if len(results) == 0 {
		return Comparison{Similarity: 100, Differences: []string{"no enabled features configured"}}
	}

	similarity := clampPercent(int(math.Round((1 - total/float64(len(results))) * 100)))

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
