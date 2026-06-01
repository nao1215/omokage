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

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
