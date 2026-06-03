package profile

import (
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
)

func TestDocumentDivergenceZeroForCorpusCenter(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()

	// A document sitting exactly on the corpus mean has nothing out of place, so
	// its divergence is zero — the outlier check must never flag the center.
	if z := DocumentDivergence(dist, dist.Mean, flags); z != 0 {
		t.Fatalf("a document at the mean should not diverge, got %.3f", z)
	}
}

func TestDocumentDivergenceRisesWithDistance(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()

	near := dist.Mean
	far := dist.Mean
	// Push the register and script features several standard deviations out: a
	// document unlike the corpus must diverge more than one near it.
	far.PoliteEndingRatio = clamp01ForTest(dist.Mean.PoliteEndingRatio + 0.5)
	far.KanjiRatio = clamp01ForTest(dist.Mean.KanjiRatio + 0.3)
	far.AverageSentenceLength = dist.Mean.AverageSentenceLength + 60

	if DocumentDivergence(dist, far, flags) <= DocumentDivergence(dist, near, flags) {
		t.Fatal("a document far from the corpus should diverge more than one at its center")
	}
}

func TestDocumentDivergenceIgnoresFingerprint(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()

	// Divergence is built on the interpretable localizable features only. A
	// document that differs only in its function-word fingerprint (which a person
	// cannot read as "out of place") must not register as an outlier.
	target := dist.Mean
	target.LexicalFrequencies = map[string]float64{"the": 0.9}
	if z := DocumentDivergence(dist, target, flags); z != 0 {
		t.Fatalf("fingerprint-only difference must not diverge, got %.3f", z)
	}
}

func TestLeaveOneOutFlagsTheOddDocument(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	steady := feature.Metrics{
		AverageSentenceLength: 30,
		PunctuationFrequency:  0.1,
		KanjiRatio:            0.3,
		HiraganaRatio:         0.6,
		KatakanaRatio:         0.1,
		PoliteEndingRatio:     0.9,
		PlainEndingRatio:      0.05,
		CharacterCount:        300,
	}
	odd := feature.Metrics{
		AverageSentenceLength: 30,
		PunctuationFrequency:  0.9,
		KanjiRatio:            0.9,
		HiraganaRatio:         0.05,
		KatakanaRatio:         0.05,
		PoliteEndingRatio:     0.05,
		PlainEndingRatio:      0.9,
		CharacterCount:        300,
	}
	samples := []feature.Metrics{steady, steady, steady, steady, steady, odd}

	divergences := LeaveOneOutDivergences(samples, flags)
	if len(divergences) != len(samples) {
		t.Fatalf("expected one divergence per document, got %d", len(divergences))
	}
	// Measured against the others, the odd document must stand far out while the
	// steady ones sit low — and far enough that the quality layer's 2.5σ bar catches
	// it even on this small corpus, which a whole-corpus z-score could not.
	if divergences[5] < 2.5 {
		t.Fatalf("the odd document should diverge sharply, got %.2fσ", divergences[5])
	}
	for i := range 5 {
		if divergences[i] >= divergences[5] {
			t.Fatalf("steady document %d should diverge less than the odd one (%.2f vs %.2f)", i, divergences[i], divergences[5])
		}
	}
}

func TestLeaveOneOutUniformCorpusHasNoOutliers(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	doc := feature.Metrics{
		AverageSentenceLength: 30,
		PoliteEndingRatio:     0.9,
		KanjiRatio:            0.3,
		HiraganaRatio:         0.6,
		CharacterCount:        300,
	}
	samples := []feature.Metrics{doc, doc, doc, doc}

	for i, z := range LeaveOneOutDivergences(samples, flags) {
		if z != 0 {
			t.Fatalf("a uniform corpus must produce no divergence, doc %d got %.3f", i, z)
		}
	}
}

func TestLeaveOneOutTooFewDocuments(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	single := []feature.Metrics{{AverageSentenceLength: 30, CharacterCount: 300}}
	got := LeaveOneOutDivergences(single, flags)
	if len(got) != 1 || got[0] != 0 {
		t.Fatalf("a single document has nothing to diverge from, got %+v", got)
	}
}

func TestHighLevelSpreadsRanksWidestFirstAndSkipsVariance(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := feature.Distribution{
		Mean: feature.Metrics{
			AverageSentenceLength:  30,
			PoliteEndingRatio:      0.5,
			KanjiRatio:             0.3,
			SentenceLengthVariance: 200, // a second-order feature, must be excluded
		},
		StdDev: feature.Metrics{
			AverageSentenceLength:  3,    // relative spread 0.10
			PoliteEndingRatio:      0.45, // relative spread 0.90 (widest)
			KanjiRatio:             0.06, // relative spread 0.20
			SentenceLengthVariance: 5000, // huge relative spread, but not localizable
		},
		DocumentCount: 8,
	}

	spreads := HighLevelSpreads(dist, flags)
	if len(spreads) == 0 {
		t.Fatal("expected spreads for the present localizable features")
	}
	if !strings.Contains(spreads[0].Feature, "polite") {
		t.Fatalf("the widest-spread feature (register) should lead, got %q", spreads[0].Feature)
	}
	for _, spread := range spreads {
		if strings.Contains(spread.Feature, "variance") {
			t.Fatalf("variance features are not interpretable and must be excluded, got %q", spread.Feature)
		}
	}
}

func TestHighLevelSpreadsSkipsAbsentFeatures(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	// Only the register feature is present; the script and length features sit at
	// (or below) the presence floor and should not be reported, so a near-empty
	// feature never manufactures a huge relative spread by dividing by ~0.
	dist := feature.Distribution{
		Mean:          feature.Metrics{PoliteEndingRatio: 0.6},
		StdDev:        feature.Metrics{PoliteEndingRatio: 0.1},
		DocumentCount: 5,
	}
	spreads := HighLevelSpreads(dist, flags)
	if len(spreads) != 1 || !strings.Contains(spreads[0].Feature, "polite") {
		t.Fatalf("only the present feature should be reported, got %+v", spreads)
	}
}

// clamp01ForTest keeps a pushed ratio within [0,1] so the synthetic far document
// stays a valid feature vector.
func clamp01ForTest(v float64) float64 {
	if v > 1 {
		return 1
	}
	if v < 0 {
		return 0
	}
	return v
}
