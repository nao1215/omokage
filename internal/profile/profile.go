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
	totalZ := 0.0
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
		totalZ += z
		results = append(results, scored{
			label:     spec.label,
			z:         z,
			direction: direction(mean, observed),
		})
	}

	if len(results) == 0 {
		return Comparison{Similarity: 100, Differences: []string{"no enabled features configured"}}
	}

	meanZ := totalZ / float64(len(results))
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
