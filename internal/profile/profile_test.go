package profile

import (
	"math"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
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

func TestRegisterPenaltySaturates(t *testing.T) {
	t.Parallel()

	values := []float64{
		registerPenalty(3),
		registerPenalty(5),
		registerPenalty(10),
		registerPenalty(50),
	}
	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			t.Fatalf("register penalty should be monotonic: %#v", values)
		}
	}
	if values[len(values)-1] > registerSaturation {
		t.Fatalf("register penalty should stay within the saturation cap %.2f, got %.6f",
			registerSaturation, values[len(values)-1])
	}
	if math.Abs(values[len(values)-1]-values[len(values)-2]) > 0.1 {
		t.Fatalf("register penalty should flatten out by large z, got %#v", values)
	}
}

func TestRegisterPenaltyLeavesToleranceRegionUntouched(t *testing.T) {
	t.Parallel()

	got := combineDrift(groupDrift{
		register:     registerTolerance,
		functionWord: 0.6,
		ngram:        0.4,
		other:        1.0,
	})
	want := 0.5 + otherStructWeight*1.0
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("register within tolerance should add nothing: got=%f want=%f", got, want)
	}
}

func TestRegisterFlipRetainsLexicalResolution(t *testing.T) {
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

	// Same lexical fingerprint and structure as the profile mean, but with a full
	// register flip. This used to collapse to 0% together with a genuinely
	// different author because the register excess was unbounded.
	registerFlip := dist.Mean
	registerFlip.PoliteEndingRatio = 0
	registerFlip.PlainEndingRatio = 1

	otherAuthor := feature.ExtractText("本日は降雨である。会議資料を作成した。結論を整理し、対応を決定した。")

	near := Score(dist, registerFlip, flags)
	far := Score(dist, otherAuthor, flags)
	if near.Similarity <= far.Similarity {
		t.Fatalf("register-flipped near match should outscore a different author: near=%d far=%d",
			near.Similarity, far.Similarity)
	}
}

func TestExactLexicalMatchBeatsTinyLexicalDrift(t *testing.T) {
	t.Parallel()

	exact := []FeatureDrift{
		{Category: categoryFunctionWord, Z: 0},
		{Category: categoryCharNgram, Z: 2},
	}
	epsilon := []FeatureDrift{
		{Category: categoryFunctionWord, Z: 0.1},
		{Category: categoryCharNgram, Z: 2},
	}

	exactBreakdown := summarizeDrifts(exact)
	epsilonBreakdown := summarizeDrifts(epsilon)
	if math.Abs(exactBreakdown.lexicalTerm-1.0) > 1e-9 {
		t.Fatalf("expected exact lexical match to average with the active n-gram group, got %f", exactBreakdown.lexicalTerm)
	}
	if exactBreakdown.meanZ >= epsilonBreakdown.meanZ {
		t.Fatalf("expected exact lexical match to produce a lower mean z: exact=%f epsilon=%f",
			exactBreakdown.meanZ, epsilonBreakdown.meanZ)
	}
	if similarityFromDrifts(exact) <= similarityFromDrifts(epsilon) {
		t.Fatalf("expected an exact lexical match to score better than a tiny drift: exact=%d epsilon=%d",
			similarityFromDrifts(exact), similarityFromDrifts(epsilon))
	}
}

func TestCompareCountsExactLexicalMatchesAsPresent(t *testing.T) {
	t.Parallel()

	got := combineCompareDriftWithCounts(groupDrift{
		functionWord: 0,
		ngram:        0.4,
	}, groupCounts{
		functionWord: 1,
		ngram:        1,
	})
	want := 0.2
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("expected exact function-word match to stay in the lexical mean: got=%f want=%f", got, want)
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

func TestScoreRecordFallsBackWithoutCalibration(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	polite := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
	}
	dist := distributionFromTexts(polite)
	target := feature.ExtractText("今日は良い天気です。散歩に出かけました。とても気持ちが良かったです。")

	legacy := Score(dist, target, flags)
	record := Record{Distribution: dist}
	calibrated := ScoreRecord(record, target, flags)
	if calibrated.Similarity != legacy.Similarity {
		t.Fatalf("expected legacy fallback without calibration: legacy=%d calibrated=%d",
			legacy.Similarity, calibrated.Similarity)
	}
}

func TestScoreRecordCalibratesSelfSimilarity(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	texts := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
		"新しい本を買いました。内容がとても面白くて一気に読み終えました。",
		"先週は仕事が忙しかったです。それでも毎日きちんと休めました。",
	}
	perDoc := make([]feature.Metrics, 0, len(texts))
	for _, text := range texts {
		perDoc = append(perDoc, feature.ExtractText(text))
	}
	record := Record{
		Distribution:   feature.Aggregate(perDoc),
		SelfSimilarity: ComputeSelfSimilarityStats(perDoc, flags),
	}

	own := feature.ExtractText("今日は朝から少し雨が降っています。傘を持って出かけました。空気が冷たくて静かでした。")
	other := feature.ExtractText("きょうもサッカーをしたよ。仲間と思いきり走ったんだ。最高に楽しい一日だったなあ。")

	legacyOwn := Score(record.Distribution, own, flags)
	calibratedOwn := ScoreRecord(record, own, flags)
	calibratedOther := ScoreRecord(record, other, flags)

	if calibratedOwn.SelfSimilarity == nil || calibratedOwn.SelfSimilarity.Median != calibratedMedianScore {
		t.Fatalf("expected calibrated self-similarity anchor at %d, got %+v",
			calibratedMedianScore, calibratedOwn.SelfSimilarity)
	}
	if calibratedOwn.Similarity <= legacyOwn.Similarity {
		t.Fatalf("expected calibration to lift own-text scores: legacy=%d calibrated=%d",
			legacyOwn.Similarity, calibratedOwn.Similarity)
	}
	if calibratedOwn.Similarity < 85 {
		t.Fatalf("expected own-text score to move near the calibrated anchor, got %d", calibratedOwn.Similarity)
	}
	if calibratedOther.Similarity >= calibratedOwn.Similarity {
		t.Fatalf("expected other author to remain below own text: own=%d other=%d",
			calibratedOwn.Similarity, calibratedOther.Similarity)
	}
	if calibratedOther.Similarity >= 50 {
		t.Fatalf("expected other author to stay clearly low after calibration, got %d", calibratedOther.Similarity)
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

func TestExplainMatchesScoreSimilarity(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := sampleDistribution()
	target := dist.Mean
	target.KanjiRatio = dist.Mean.KanjiRatio + 3*dist.StdDev.KanjiRatio

	score := Score(dist, target, flags)
	explanation := Explain(dist, target, nil, flags)
	if explanation.Similarity != score.Similarity {
		t.Fatalf("Explain similarity %d should match Score similarity %d",
			explanation.Similarity, score.Similarity)
	}
}

func TestExplainRecordMatchesScoreRecordSimilarity(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	texts := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
	}
	perDoc := make([]feature.Metrics, 0, len(texts))
	for _, text := range texts {
		perDoc = append(perDoc, feature.ExtractText(text))
	}
	record := Record{
		Distribution:   feature.Aggregate(perDoc),
		SelfSimilarity: ComputeSelfSimilarityStats(perDoc, flags),
	}
	target := feature.ExtractText("今日は良い天気です。散歩に出かけました。とても気持ちが良かったです。")

	score := ScoreRecord(record, target, flags)
	explanation := ExplainRecord(record, target, nil, flags)
	if explanation.Similarity != score.Similarity {
		t.Fatalf("ExplainRecord similarity %d should match ScoreRecord similarity %d",
			explanation.Similarity, score.Similarity)
	}
}

func TestExplainSurfacesHighLevelDriftFirst(t *testing.T) {
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

	plainText := "今日は良い天気である。散歩に出かける。とても気持ちが良いのだった。"
	target := feature.ExtractText(plainText)
	explanation := Explain(dist, target, feature.ExtractSegments(plainText), flags)

	if len(explanation.Drifts) == 0 {
		t.Fatal("expected drifts for a register-shifted target")
	}

	// The first actionable drift must be a high-level, editable feature, not a
	// function word or character n-gram.
	var firstActionable *FeatureDrift
	for i := range explanation.Drifts {
		if explanation.Drifts[i].Actionable {
			firstActionable = &explanation.Drifts[i]
			break
		}
	}
	if firstActionable == nil {
		t.Fatal("expected at least one actionable drift")
	}
	if firstActionable.Level != levelHigh {
		t.Fatalf("expected the first actionable drift to be high-level, got %q (%s)",
			firstActionable.Feature, firstActionable.Level)
	}

	// Priority must be a dense 1..N ranking with high-level features before
	// low-level ones.
	lastHighIndex := -1
	firstLowIndex := -1
	for i, drift := range explanation.Drifts {
		if drift.Priority != i+1 {
			t.Fatalf("priority should be 1-based dense order, got %d at index %d", drift.Priority, i)
		}
		if drift.Level == levelHigh {
			lastHighIndex = i
		} else if firstLowIndex == -1 {
			firstLowIndex = i
		}
	}
	if firstLowIndex != -1 && lastHighIndex > firstLowIndex {
		t.Fatalf("high-level drift at index %d should precede low-level at %d", lastHighIndex, firstLowIndex)
	}

	// The numeric detail an editor needs must be populated.
	if firstActionable.Mean == 0 && firstActionable.StdDev == 0 {
		t.Fatalf("expected reference statistics on drift %q", firstActionable.Feature)
	}
}

func TestExplainLocalizesDriftToParagraph(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	polite := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
		"新しい本を買いました。内容がとても面白くて一気に読み終えました。",
	}
	dist := distributionFromTexts(polite)

	// A faithful first paragraph followed by a register-flipped second one. The
	// drifting paragraph must be the one localized.
	keep := "今日は朝から良い天気です。少し散歩に出かけました。空気が澄んでいてとても気持ちが良かったです。"
	drift := "結論を以下に記述する。本件は重要である。対応を実施するものとする。詳細は別途共有することとする。"
	doc := keep + "\n\n" + drift

	explanation := Explain(dist, feature.ExtractText(doc), feature.ExtractSegments(doc), flags)
	if len(explanation.Segments) == 0 {
		t.Fatal("expected segment localization")
	}
	if explanation.Segments[0].Index != 2 {
		t.Fatalf("expected the second paragraph to drift most, got paragraph %d", explanation.Segments[0].Index)
	}
	if explanation.Segments[0].Feature == "" {
		t.Fatal("expected a top drifting feature for the worst paragraph")
	}
	// The headline z and the named feature must correspond, and the feature must
	// clear the reporting bar.
	if explanation.Segments[0].Z < 1.0 {
		t.Fatalf("a reported paragraph must hold an actionable drift, got %.2fσ", explanation.Segments[0].Z)
	}
}

func TestLocalizationExcludesDocumentGlobalFeatures(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	// A profile whose documents carry markdown structure and varied paragraph
	// lengths, so the document-global features have a non-zero mean and spread.
	dist := distributionFromTexts([]string{
		"# 見出し\n\n今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。\n\n- 箇条書き\n- もう一つ",
		"## 別の見出し\n\n昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。とても満足な一日でした。\n\n> 引用文",
		"### 三つ目\n\n週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
	})

	// A plain prose draft with no headings or bullets. Its per-paragraph markdown
	// structure / bullet / newline / paragraph-length-variance all collapse to a
	// constant, which used to dominate every paragraph as the "top" drift. None of
	// those document-global features may now be named as a paragraph's drift.
	doc := "今日は良い天気である。散歩に出かける。気持ちが良いのだった。\n\n" +
		"午後も外を歩いた。風がとても心地よかった。空はどこまでも青かった。"
	explanation := Explain(dist, feature.ExtractText(doc), feature.ExtractSegments(doc), flags)

	documentGlobal := map[string]bool{
		"markdown structure frequency": true,
		"bullet-list frequency":        true,
		"newline frequency":            true,
		"paragraph length variance":    true,
		"sentence length variance":     true,
	}
	for _, segment := range explanation.Segments {
		if documentGlobal[segment.Feature] {
			t.Fatalf("paragraph #%d localized to a document-global feature %q", segment.Index, segment.Feature)
		}
	}
}

func TestNearMatchProducesNoMisleadingSegments(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	corpus := []string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
		"週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。",
		"新しい本を買いました。内容がとても面白くて一気に読み終えました。",
	}
	dist := distributionFromTexts(corpus)

	// A draft built from the author's own sentences: a near-match. No paragraph
	// holds a genuinely actionable local drift, so none must be reported rather
	// than a list of paragraphs with negligible or document-global drift.
	doc := corpus[0] + "\n\n" + corpus[1]
	explanation := Explain(dist, feature.ExtractText(doc), feature.ExtractSegments(doc), flags)
	for _, segment := range explanation.Segments {
		if segment.Z < 1.0 {
			t.Fatalf("near-match paragraph #%d reported with sub-threshold drift %.2fσ", segment.Index, segment.Z)
		}
	}
}

func TestExplainSkipsShortSegments(t *testing.T) {
	t.Parallel()

	flags := config.Default("sample").Features
	dist := distributionFromTexts([]string{
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。",
		"昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。",
	})

	// A bare heading is too short to localize and must not appear as drift.
	doc := "# 見出し\n\n" + "今日はとても良い天気です。散歩に出かけました。気持ちが良かったです。"
	explanation := Explain(dist, feature.ExtractText(doc), feature.ExtractSegments(doc), flags)
	for _, segment := range explanation.Segments {
		if strings.Contains(segment.Excerpt, "見出し") {
			t.Fatalf("a short heading should be skipped, got %q", segment.Excerpt)
		}
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
