package feature

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractText(t *testing.T) {
	t.Parallel()

	metrics := ExtractText("# Heading\n\nそして文章です。だから続きます。\n- one\n- two\n")
	if metrics.SentenceCount == 0 {
		t.Fatal("expected sentences to be detected")
	}
	if metrics.BulletRatio <= 0 {
		t.Fatalf("expected bullet ratio > 0, got %f", metrics.BulletRatio)
	}
	if metrics.MarkdownStructureDensity <= 0 {
		t.Fatalf("expected markdown structure density > 0, got %f", metrics.MarkdownStructureDensity)
	}
	if metrics.KanjiRatio <= 0 && metrics.HiraganaRatio <= 0 {
		t.Fatalf("expected Japanese script ratios > 0, got kanji=%f hiragana=%f", metrics.KanjiRatio, metrics.HiraganaRatio)
	}
}

func TestExtractTextDetectsSentenceEndings(t *testing.T) {
	t.Parallel()

	polite := ExtractText("今日は晴れです。散歩に行きます。とても良い一日でした。")
	if polite.PoliteEndingRatio <= polite.PlainEndingRatio {
		t.Fatalf("expected polite register to dominate, got polite=%f plain=%f", polite.PoliteEndingRatio, polite.PlainEndingRatio)
	}

	plain := ExtractText("今日は晴れである。散歩に行く。とても良い一日だった。")
	if plain.PlainEndingRatio <= plain.PoliteEndingRatio {
		t.Fatalf("expected plain register to dominate, got polite=%f plain=%f", plain.PoliteEndingRatio, plain.PlainEndingRatio)
	}

	// English text has no Japanese sentence-ending forms.
	english := ExtractText("This is a plain English sentence. It has no Japanese endings.")
	if english.PoliteEndingRatio != 0 || english.PlainEndingRatio != 0 {
		t.Fatalf("expected zero ending ratios for English, got polite=%f plain=%f", english.PoliteEndingRatio, english.PlainEndingRatio)
	}
}

func TestExtractTextStripsCodeFromEveryFeature(t *testing.T) {
	t.Parallel()

	prose := "This is an ordinary English paragraph. It keeps two plain sentences.\n\n" +
		"And a second paragraph that holds the same calm voice all the way through."
	withCode := prose + "\n\n```go\nfunc main() {\n\tfor i := 0; i < 10; i++ {\n\t\t// a heading-looking # comment and a | pipe\n\t\tfmt.Println(i)\n\t}\n}\n```\n"

	bare := ExtractText(prose)
	coded := ExtractText(withCode)

	// Code is removed before *every* feature is measured, not only the lexical and
	// n-gram ones, so appending a fenced block must not move the structural
	// features. Otherwise a draft scores as drifting purely for containing code.
	if bare.SentenceCount != coded.SentenceCount {
		t.Fatalf("sentence count moved by code block: %d vs %d", bare.SentenceCount, coded.SentenceCount)
	}
	if bare.AverageSentenceLength != coded.AverageSentenceLength {
		t.Fatalf("average sentence length moved by code block: %f vs %f", bare.AverageSentenceLength, coded.AverageSentenceLength)
	}
	if bare.MarkdownStructureDensity != coded.MarkdownStructureDensity {
		t.Fatalf("markdown density moved by code block: %f vs %f", bare.MarkdownStructureDensity, coded.MarkdownStructureDensity)
	}
	if bare.CharacterCount != coded.CharacterCount {
		t.Fatalf("character count moved by code block: %d vs %d", bare.CharacterCount, coded.CharacterCount)
	}
	if bare.PunctuationFrequency != coded.PunctuationFrequency {
		t.Fatalf("punctuation frequency moved by code block: %f vs %f", bare.PunctuationFrequency, coded.PunctuationFrequency)
	}
}

func TestSplitSentencesKeepsInWordPeriodsTogether(t *testing.T) {
	t.Parallel()

	// Version numbers, domains, decimals, and abbreviations all carry an interior
	// period that must not be read as a sentence boundary, which would inflate the
	// sentence count and corrupt the sentence-length features for technical prose.
	metrics := ExtractText("The build moved from 1.2.3 to 1.10.0 today. See example.com or v2.1 for details.")
	if metrics.SentenceCount != 2 {
		t.Fatalf("expected 2 sentences across version numbers and domains, got %d", metrics.SentenceCount)
	}
}

func TestPlainEndingDetectsOpenClassPredicates(t *testing.T) {
	t.Parallel()

	// 常体 is the open class of plain-form predicates, not a short word list. Each
	// of these ends a plain sentence and must register as plain, never polite.
	plainCases := []string{
		"毎朝公園を走る。",     // godan/ichidan verb, dictionary form
		"彼は分厚い本を読んだ。",  // past plain (…んだ)
		"その山はとても高い。",   // i-adjective
		"もう時間がない。",     // negative ない
		"これはただの水だ。",    // copula だ
		"明日はきっと晴れである。", // copula である
	}
	for _, sentence := range plainCases {
		m := ExtractText(sentence)
		if m.PlainEndingRatio <= 0 {
			t.Fatalf("expected plain ending for %q, got plain=%f polite=%f", sentence, m.PlainEndingRatio, m.PoliteEndingRatio)
		}
		if m.PoliteEndingRatio != 0 {
			t.Fatalf("expected no polite ending for %q, got polite=%f", sentence, m.PoliteEndingRatio)
		}
	}

	// A polite sentence trailing a question particle is still polite: the tail
	// particle is trimmed before the form is read.
	question := ExtractText("もう行きますか。")
	if question.PoliteEndingRatio <= 0 || question.PlainEndingRatio != 0 {
		t.Fatalf("expected polite ending for trailing-か sentence, got polite=%f plain=%f", question.PoliteEndingRatio, question.PlainEndingRatio)
	}
}

func TestCollectFilesFiltersSupportedExtensions(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "a.md"), "a")
	mustWrite(t, filepath.Join(root, "b.txt"), "b")
	mustWrite(t, filepath.Join(root, "c.csv"), "c")

	files, err := CollectFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if filepath.Base(files[0]) != "a.md" || filepath.Base(files[1]) != "b.txt" {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func TestAggregateComputesMeanAndStdDev(t *testing.T) {
	t.Parallel()

	perDoc := []Metrics{
		{KanjiRatio: 0.2, AverageSentenceLength: 10, SentenceCount: 1, CharacterCount: 5},
		{KanjiRatio: 0.4, AverageSentenceLength: 30, SentenceCount: 2, CharacterCount: 7},
	}

	dist := Aggregate(perDoc)
	if dist.DocumentCount != 2 {
		t.Fatalf("expected document count 2, got %d", dist.DocumentCount)
	}
	if math.Abs(dist.Mean.KanjiRatio-0.3) > 1e-9 {
		t.Fatalf("expected mean kanji ratio 0.3, got %f", dist.Mean.KanjiRatio)
	}
	if dist.Mean.AverageSentenceLength != 20 {
		t.Fatalf("expected mean sentence length 20, got %f", dist.Mean.AverageSentenceLength)
	}
	// Population std dev of {0.2, 0.4} around 0.3 is 0.1.
	if math.Abs(dist.StdDev.KanjiRatio-0.1) > 1e-9 {
		t.Fatalf("expected std dev 0.1, got %f", dist.StdDev.KanjiRatio)
	}
	// Counts accumulate across the corpus.
	if dist.SentenceCount != 3 || dist.CharacterCount != 12 {
		t.Fatalf("unexpected totals: sentences=%d characters=%d", dist.SentenceCount, dist.CharacterCount)
	}
}

func TestAggregateEmptyCorpus(t *testing.T) {
	t.Parallel()

	dist := Aggregate(nil)
	if dist.DocumentCount != 0 {
		t.Fatalf("expected empty distribution, got document count %d", dist.DocumentCount)
	}
	if dist.Mean.KanjiRatio != 0 || dist.StdDev.KanjiRatio != 0 {
		t.Fatal("expected zero-valued metrics for an empty corpus")
	}
}

func TestExtractCorpus(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "a.md"), "# Title\n\nそして文章です。だから続きます。\n")
	mustWrite(t, filepath.Join(root, "b.md"), "- one\n- two\n- three\n")

	files, err := CollectFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	dist, err := ExtractCorpus(files)
	if err != nil {
		t.Fatal(err)
	}
	if dist.DocumentCount != 2 {
		t.Fatalf("expected 2 documents, got %d", dist.DocumentCount)
	}
}

func TestExtractCorpusMissingFile(t *testing.T) {
	t.Parallel()

	if _, err := ExtractCorpus([]string{filepath.Join(t.TempDir(), "missing.md")}); err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

func TestExtractSegments(t *testing.T) {
	t.Parallel()

	text := "最初の段落です。短い文を書きます。\n\n   \n\n二つ目の段落です。もう少し続きます。"
	segments := ExtractSegments(text)
	if len(segments) != 2 {
		t.Fatalf("expected 2 non-empty segments, got %d", len(segments))
	}
	if segments[0].Index != 1 || segments[1].Index != 2 {
		t.Fatalf("expected 1-based dense indexes, got %d and %d", segments[0].Index, segments[1].Index)
	}
	for _, segment := range segments {
		if segment.Kind != "paragraph" {
			t.Fatalf("unexpected segment kind: %q", segment.Kind)
		}
		if segment.Metrics.CharacterCount == 0 {
			t.Fatalf("segment %d should carry metrics", segment.Index)
		}
	}
}

func TestExtractFileWithSegments(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "doc.md")
	mustWrite(t, path, "段落一。これは本文です。\n\n段落二。これも本文です。")

	metrics, segments, err := ExtractFileWithSegments(path)
	if err != nil {
		t.Fatal(err)
	}
	if metrics.CharacterCount == 0 {
		t.Fatal("expected whole-document metrics")
	}
	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
