package feature

import (
	"math"
	"strings"
	"testing"
)

func TestExtractTextPopulatesLexicalAndNgrams(t *testing.T) {
	t.Parallel()

	metrics := ExtractText("私は今日も本を読みます。そして散歩に出かけます。")
	if len(metrics.LexicalFrequencies) == 0 {
		t.Fatal("expected lexical frequencies to be populated")
	}
	if metrics.LexicalFrequencies["は"] <= 0 {
		t.Fatalf("expected the particle は to register a frequency, got %f", metrics.LexicalFrequencies["は"])
	}
	if len(metrics.CharNgrams) == 0 {
		t.Fatal("expected character n-grams to be populated")
	}
}

func TestLexicalFrequenciesSeparatesLanguages(t *testing.T) {
	t.Parallel()

	japanese := ExtractText("これはテストの文章です。とても面白い本でした。")
	english := ExtractText("the quick brown fox jumps over the lazy dog and the cat")

	if japanese.LexicalFrequencies["の"] <= 0 {
		t.Fatalf("expected Japanese particle frequency > 0, got %f", japanese.LexicalFrequencies["の"])
	}
	if japanese.LexicalFrequencies["the"] != 0 {
		t.Fatalf("expected English function word to be absent from Japanese text, got %f", japanese.LexicalFrequencies["the"])
	}
	if english.LexicalFrequencies["the"] <= 0 {
		t.Fatalf("expected English function word frequency > 0, got %f", english.LexicalFrequencies["the"])
	}
	if english.LexicalFrequencies["の"] != 0 {
		t.Fatalf("expected Japanese particle to be absent from English text, got %f", english.LexicalFrequencies["の"])
	}
}

func TestStripCodeRemovesFencedAndInline(t *testing.T) {
	t.Parallel()

	text := "before\n```go\nfunc main() {}\n```\nafter `inline()` end"
	stripped := stripCode(text)

	for _, fragment := range []string{"func main", "inline()"} {
		if strings.Contains(stripped, fragment) {
			t.Fatalf("expected %q to be removed, got %q", fragment, stripped)
		}
	}
	for _, fragment := range []string{"before", "after", "end"} {
		if !strings.Contains(stripped, fragment) {
			t.Fatalf("expected %q to be preserved, got %q", fragment, stripped)
		}
	}
}

func TestCodeIsExcludedFromLexicalFeatures(t *testing.T) {
	t.Parallel()

	prose := ExtractText("これは説明の文章です。とても分かりやすいと思います。")
	withCode := ExtractText("これは説明の文章です。\n```go\nfunc main() { the the the and and or }\n```\nとても分かりやすいと思います。")

	// The English keywords live only inside the fenced block, so stripping code
	// must keep them out of the lexical vector.
	if withCode.LexicalFrequencies["the"] != 0 {
		t.Fatalf("expected code keywords to be excluded, got the=%f", withCode.LexicalFrequencies["the"])
	}
	if math.Abs(prose.LexicalFrequencies["の"]-withCode.LexicalFrequencies["の"]) > 0.05 {
		t.Fatalf("expected prose particle frequency to be stable across code, got %f vs %f",
			prose.LexicalFrequencies["の"], withCode.LexicalFrequencies["の"])
	}
}

func TestStripCodeDropsUnclosedFence(t *testing.T) {
	t.Parallel()

	// An unterminated fence is treated as code through end of document, dropping
	// everything after it. This is the safe interpretation and must not panic.
	stripped := stripCode("intro text\n```go\nfunc x() {}\nmore code without closing")
	if !strings.Contains(stripped, "intro text") {
		t.Fatalf("expected text before the fence to survive, got %q", stripped)
	}
	if strings.Contains(stripped, "func x") {
		t.Fatalf("expected code after an unclosed fence to be dropped, got %q", stripped)
	}
}

func TestAggregatePopulatesLexicalAndNgramDistribution(t *testing.T) {
	t.Parallel()

	docs := []Metrics{
		ExtractText("私は本を読みます。とても面白いです。"),
		ExtractText("私は今日も本を読みました。やはり面白いです。"),
		ExtractText("彼も本が好きです。一緒に読みます。"),
	}
	dist := Aggregate(docs)

	if len(dist.Mean.LexicalFrequencies) == 0 || len(dist.StdDev.LexicalFrequencies) == 0 {
		t.Fatal("expected aggregated lexical mean and std to be populated")
	}
	if len(dist.Mean.CharNgrams) == 0 || len(dist.StdDev.CharNgrams) == 0 {
		t.Fatal("expected aggregated char n-gram mean and std to be populated")
	}
	if dist.Mean.LexicalFrequencies["は"] <= 0 {
		t.Fatalf("expected a positive mean for a heavily used particle, got %f", dist.Mean.LexicalFrequencies["は"])
	}
}

func TestTopNgramsIsDeterministicOnTies(t *testing.T) {
	t.Parallel()

	totals := map[string]float64{"aa": 1, "bb": 1, "cc": 1}
	first := topNgrams(totals, 2)
	second := topNgrams(totals, 2)
	if len(first) != 2 {
		t.Fatalf("expected limit to cap the result at 2, got %d", len(first))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("expected deterministic tie-break, got %v then %v", first, second)
		}
	}
}

func TestCharBigramsIncludesBigramsAndTrigrams(t *testing.T) {
	t.Parallel()

	ngrams := charBigrams("hello")
	if _, ok := ngrams["he"]; !ok {
		t.Fatalf("expected bigram 'he' to be present, got %v keys", len(ngrams))
	}
	if _, ok := ngrams["hel"]; !ok {
		t.Fatalf("expected trigram 'hel' to be present, got %v keys", len(ngrams))
	}
}
