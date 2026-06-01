package profile

import (
	"strings"
	"testing"

	"github.com/nao1215/dyer/internal/config"
	"github.com/nao1215/dyer/internal/feature"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	reference := feature.Metrics{
		AverageSentenceLength:    10,
		SentenceLengthVariance:   4,
		PunctuationFrequency:     0.1,
		NewlineFrequency:         0.05,
		BulletRatio:              0.1,
		ConjunctionFrequency:     0.05,
		KanjiRatio:               0.4,
		HiraganaRatio:            0.5,
		KatakanaRatio:            0.1,
		ParagraphLengthVariance:  8,
		MarkdownStructureDensity: 0.2,
	}
	target := reference
	target.BulletRatio = 0.6
	target.MarkdownStructureDensity = 0.7

	comparison := Compare(reference, target, flags)
	if comparison.Similarity >= 100 {
		t.Fatalf("expected drift to lower similarity, got %d", comparison.Similarity)
	}
	if len(comparison.Differences) == 0 {
		t.Fatal("expected differences")
	}
	if !strings.Contains(comparison.Differences[0], "reference") {
		t.Fatalf("unexpected difference message: %q", comparison.Differences[0])
	}
}

func sampleDistribution() feature.Distribution {
	return feature.Distribution{
		Mean: feature.Metrics{
			AverageSentenceLength:    40,
			SentenceLengthVariance:   200,
			PunctuationFrequency:     0.15,
			NewlineFrequency:         0.04,
			BulletRatio:              0.12,
			ConjunctionFrequency:     0.01,
			KanjiRatio:               0.33,
			HiraganaRatio:            0.45,
			KatakanaRatio:            0.22,
			ParagraphLengthVariance:  300,
			MarkdownStructureDensity: 0.30,
		},
		StdDev: feature.Metrics{
			AverageSentenceLength:    8,
			SentenceLengthVariance:   60,
			PunctuationFrequency:     0.02,
			NewlineFrequency:         0.01,
			BulletRatio:              0.05,
			ConjunctionFrequency:     0.005,
			KanjiRatio:               0.04,
			HiraganaRatio:            0.05,
			KatakanaRatio:            0.06,
			ParagraphLengthVariance:  90,
			MarkdownStructureDensity: 0.08,
		},
		DocumentCount: 100,
	}
}

func TestScoreRewardsOnProfileText(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()

	// A target sitting exactly at the author's mean is maximally similar.
	onProfile := Score(dist, dist.Mean, flags)
	if onProfile.Similarity != 100 {
		t.Fatalf("expected 100%% similarity at the mean, got %d", onProfile.Similarity)
	}
	if len(onProfile.Differences) == 0 || !strings.Contains(onProfile.Differences[0], "no significant") {
		t.Fatalf("expected no drift message, got %#v", onProfile.Differences)
	}
}

func TestScorePenalizesDeviationOutsideSpread(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()

	target := dist.Mean
	// Push kanji ratio well outside the author's per-document spread.
	target.KanjiRatio = dist.Mean.KanjiRatio + 4*dist.StdDev.KanjiRatio

	deviated := Score(dist, target, flags)
	if deviated.Similarity >= 100 {
		t.Fatalf("expected a deviating target to score below 100, got %d", deviated.Similarity)
	}
	if !strings.Contains(deviated.Differences[0], "kanji ratio") {
		t.Fatalf("expected kanji ratio to be the top drift, got %q", deviated.Differences[0])
	}
	if !strings.Contains(deviated.Differences[0], "higher") {
		t.Fatalf("expected drift direction to be higher, got %q", deviated.Differences[0])
	}
}

func TestScoreWithNoEnabledFeatures(t *testing.T) {
	t.Parallel()

	comparison := Score(sampleDistribution(), feature.Metrics{}, config.Features{})
	if comparison.Similarity != 100 {
		t.Fatalf("expected 100%% when no features are enabled, got %d", comparison.Similarity)
	}
}

func TestScoreSurvivesZeroStdDev(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()
	dist.StdDev = feature.Metrics{} // a single-document profile has no spread

	// The std floor must keep z-scores finite rather than producing NaN/Inf.
	comparison := Score(dist, dist.Mean, flags)
	if comparison.Similarity < 0 || comparison.Similarity > 100 {
		t.Fatalf("similarity out of range with zero std dev: %d", comparison.Similarity)
	}
}
