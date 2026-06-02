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

func TestScoreDetectsRegisterShift(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	polite := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
		"新しい本を買いました。内容がとても面白くて一気に読み終えました。",
		"先週は仕事が忙しかったです。それでも毎日きちんと休めました。",
	}
	dist := distributionFromTexts(polite)

	politeTarget := feature.ExtractText("今日は良い天気です。散歩に出かけます。とても気持ちが良いです。")
	plainTarget := feature.ExtractText("今日は良い天気である。散歩に出かける。とても気持ちが良いのだった。")

	politeScore := Score(dist, politeTarget, flags)
	plainScore := Score(dist, plainTarget, flags)
	if plainScore.Similarity >= politeScore.Similarity {
		t.Fatalf("expected a register shift to lower similarity: polite=%d plain=%d",
			politeScore.Similarity, plainScore.Similarity)
	}
}

func TestScoreRejectsCrossLanguageText(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	japanese := []string{
		"私は毎朝コーヒーを飲みます。新聞を読みながらゆっくり過ごします。",
		"昼休みには散歩をします。季節の移り変わりを感じられて楽しいです。",
		"夜は本を読んでから眠ります。静かな時間がとても好きです。",
		"休日は料理を作ります。家族と一緒に食べる時間が幸せです。",
	}
	dist := distributionFromTexts(japanese)

	english := feature.ExtractText("This is an English paragraph about the weather and the city. " +
		"It contains the kind of function words that an English writer would use every day.")
	comparison := Score(dist, english, flags)
	if comparison.Similarity > 40 {
		t.Fatalf("expected cross-language text to score low, got %d", comparison.Similarity)
	}
}

func TestScoreSeparatesLexicalFingerprint(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	authorA := []string{
		"私はコーヒーが好きです。毎朝必ず一杯飲みます。香りがとても良いからです。",
		"昨日もカフェに行きました。新しい豆を試しました。深い味わいでした。",
		"週末は自分で豆を挽きます。手間をかけるほど美味しくなります。",
		"友人にもコーヒーを勧めました。彼もすっかり気に入ったようです。",
	}
	distA := distributionFromTexts(authorA)

	ownText := feature.ExtractText("今日もコーヒーを淹れました。豆の香りに癒やされます。やはり毎日飲みたいです。")
	otherText := feature.ExtractText("きょうもサッカーをしたよ。仲間と思いきり走ったんだ。最高に楽しい一日だったなあ。")

	ownScore := Score(distA, ownText, flags)
	otherScore := Score(distA, otherText, flags)
	if otherScore.Similarity >= ownScore.Similarity {
		t.Fatalf("expected a same-register impostor to score lower than the author's own text: own=%d other=%d",
			ownScore.Similarity, otherScore.Similarity)
	}
}

func distributionFromTexts(texts []string) feature.Distribution {
	perDoc := make([]feature.Metrics, 0, len(texts))
	for _, text := range texts {
		perDoc = append(perDoc, feature.ExtractText(text))
	}
	return feature.Aggregate(perDoc)
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

func TestScoreSkipsDegenerateFeatures(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features

	// Simulate an English profile: the Japanese-only features (kanji/hiragana/
	// katakana ratios and sentence-ending ratios) are identically zero, while the
	// language-agnostic features carry real values.
	dist := feature.Distribution{
		Mean: feature.Metrics{
			AverageSentenceLength:    18,
			PunctuationFrequency:     0.12,
			BulletRatio:              0.2,
			MarkdownStructureDensity: 0.3,
		},
		StdDev: feature.Metrics{
			AverageSentenceLength:    3,
			PunctuationFrequency:     0.02,
			BulletRatio:              0.05,
			MarkdownStructureDensity: 0.06,
		},
		DocumentCount: 50,
	}

	// A target that matches the live features but also has zero Japanese features
	// must score 100: the dead features must not be counted as perfect matches
	// that pull the score around. (Without the skip, they would still read as
	// matches, but they would dominate the average and mask real drift.)
	onProfile := Score(dist, dist.Mean, flags)
	if onProfile.Similarity != 100 {
		t.Fatalf("expected 100%% for a matching English target, got %d", onProfile.Similarity)
	}

	// Drift in a live feature must still be detected despite the dead features.
	target := dist.Mean
	target.BulletRatio = 0.9
	drifted := Score(dist, target, flags)
	if drifted.Similarity >= onProfile.Similarity {
		t.Fatalf("expected bullet drift to lower similarity, got %d", drifted.Similarity)
	}
	if !strings.Contains(drifted.Differences[0], "bullet") {
		t.Fatalf("expected bullet drift to be reported, got %q", drifted.Differences[0])
	}
}

func TestScoreFlagsDegenerateFeatureWhenTargetHasIt(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features

	// An author who never uses bullet lists (mean and std both zero).
	dist := feature.Distribution{
		Mean:          feature.Metrics{AverageSentenceLength: 20},
		StdDev:        feature.Metrics{AverageSentenceLength: 4},
		DocumentCount: 30,
	}

	// A target that suddenly uses many bullets is genuine drift and must be
	// penalized, not skipped.
	target := dist.Mean
	target.BulletRatio = 0.6
	drifted := Score(dist, target, flags)
	if drifted.Similarity >= 100 {
		t.Fatalf("expected drift from a never-used feature to lower similarity, got %d", drifted.Similarity)
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
