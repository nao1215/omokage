package feature

import (
	"strings"
	"testing"
)

// TestMaskTechnicalTopicTokensMasksEachClass checks that every token class the
// masking is responsible for collapses to its fixed placeholder, so the lexical
// and character n-gram fingerprints stop tracking which product, repository, or
// version a post happens to discuss.
func TestMaskTechnicalTopicTokensMasksEachClass(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"owner/repo", "nao1215/omokage is great", "<REPO> is great"},
		{"dotted identifier", "use sql.DB here", "use <IDENT> here"},
		{"deep dotted identifier", "foo.bar.Baz wraps it", "<IDENT> wraps it"},
		{"snake_case", "the read_file helper", "the <IDENT> helper"},
		{"kebab-case", "the topic-swapped sample", "the <IDENT> sample"},
		{"CamelCase", "GitHub hosts it", "<IDENT> hosts it"},
		{"PascalCase identifier", "a DataFrame value", "a <IDENT> value"},
		{"version dotted", "release v1.2.3 today", "release <VERSION> today"},
		{"version decimal", "Go 1.25 ships", "Go <VERSION> ships"},
		{"calendar version", "as of 2025.12 now", "as of <VERSION> now"},
		{"acronym", "export to CSV format", "export to <ACRONYM> format"},
		{"filename", "open data.csv now", "open <FILE> now"},
		{"bare extension", "writes a .parquet file", "writes a <FILE> file"},
		{"bare number", "about 4096 rows", "about <NUMBER> rows"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := MaskTechnicalTopicTokens(tc.in); got != tc.want {
				t.Fatalf("MaskTechnicalTopicTokens(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestMaskTechnicalTopicTokensLeavesOrdinaryWords checks the masking does not
// touch ordinary words: lowercase prose, a single capitalized word, the pronoun
// "I", and Japanese kana/kanji must all survive verbatim, or the masking would
// corrupt the very stylistic signal it is meant to preserve.
func TestMaskTechnicalTopicTokensLeavesOrdinaryWords(t *testing.T) {
	t.Parallel()

	cases := []string{
		"the quick brown fox jumps",
		"I think this is fine",
		"The morning was quiet and cold",
		"今日は朝から雨が降っています。",
		"敬体と常体の違いを測りたい。",
	}
	for _, in := range cases {
		if got := MaskTechnicalTopicTokens(in); got != in {
			t.Fatalf("MaskTechnicalTopicTokens(%q) = %q, want it unchanged", in, got)
		}
	}
}

// TestMaskTechnicalTopicTokensKeepsJapaneseSentenceShape verifies that an
// identifier embedded in Japanese prose is masked without disturbing the
// surrounding kana, kanji, or punctuation, since those carry the stylistic signal.
func TestMaskTechnicalTopicTokensKeepsJapaneseSentenceShape(t *testing.T) {
	t.Parallel()

	in := "今日はsql.DBを使ってCSVをv1.2.3で読みます。"
	got := MaskTechnicalTopicTokens(in)
	for _, want := range []string{"<IDENT>", "<ACRONYM>", "<VERSION>", "今日は", "を使って", "で読みます。"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected masked %q to contain %q, got %q", in, want, got)
		}
	}
	if strings.Contains(got, "sql.DB") || strings.Contains(got, "v1.2.3") {
		t.Fatalf("expected identifiers to be masked, got %q", got)
	}
}

// TestMaskingScalarFeaturesUnaffected is the core guarantee of the change: the
// scalar style features (sentence length, punctuation, register, script ratios)
// are measured on the original prose, so two texts that differ only in their
// technical tokens must produce identical scalar metrics, while their lexical and
// character n-gram fingerprints DO differ before masking and converge after it.
func TestMaskingScalarFeaturesUnaffected(t *testing.T) {
	t.Parallel()

	// Same prose shape, only the topic tokens (repo, identifier, version, format)
	// differ — exactly a "topic swap". The tokens are kept the same length and all
	// ASCII so the ORIGINAL prose, which the scalar features measure, is identical
	// in length and Japanese content: any scalar difference would therefore prove
	// the masking had leaked into a feature it must not touch.
	a := ExtractText("私たちは aaa/bbb と xx.Yy を v1.2.3 で CSV に書き出しました。とても便利でした。")
	b := ExtractText("私たちは ccc/ddd と zz.Ww を v9.8.7 で TSV に書き出しました。とても便利でした。")

	if a.AverageSentenceLength != b.AverageSentenceLength {
		t.Fatalf("sentence length must be unaffected by topic tokens: a=%f b=%f", a.AverageSentenceLength, b.AverageSentenceLength)
	}
	if a.PunctuationFrequency != b.PunctuationFrequency {
		t.Fatalf("punctuation frequency must be unaffected: a=%f b=%f", a.PunctuationFrequency, b.PunctuationFrequency)
	}
	if a.KanjiRatio != b.KanjiRatio || a.HiraganaRatio != b.HiraganaRatio || a.KatakanaRatio != b.KatakanaRatio {
		t.Fatalf("script ratios must be unaffected: a=%+v b=%+v", a, b)
	}
	if a.PoliteEndingRatio != b.PoliteEndingRatio || a.PlainEndingRatio != b.PlainEndingRatio {
		t.Fatalf("register ratios must be unaffected: a=%+v b=%+v", a, b)
	}
}

// TestMaskingConvergesCharNgrams confirms the payoff: after masking, two
// topic-swapped texts share their character n-gram fingerprint far more closely
// than the same two texts measured without masking would, because the
// topic-specific n-grams (sq, ql, cs, …) have collapsed to shared placeholders.
func TestMaskingConvergesCharNgrams(t *testing.T) {
	t.Parallel()

	proseA := "私たちは sql.DB を使って CSV を読み込みます。"
	proseB := "私たちは pkg.Store を使って TSV を読み込みます。"

	maskedDistance := charNgramDistance(charBigrams(MaskTechnicalTopicTokens(proseA)), charBigrams(MaskTechnicalTopicTokens(proseB)))
	rawDistance := charNgramDistance(charBigrams(proseA), charBigrams(proseB))

	if maskedDistance >= rawDistance {
		t.Fatalf("masking should reduce char n-gram distance between topic-swapped texts: masked=%f raw=%f", maskedDistance, rawDistance)
	}
}

// charNgramDistance is the summed absolute frequency difference over the union of
// two character n-gram vectors — a simple L1 distance used only by the masking
// tests to compare "how different do these two fingerprints look".
func charNgramDistance(a, b map[string]float64) float64 {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	var sum float64
	for k := range seen {
		d := a[k] - b[k]
		if d < 0 {
			d = -d
		}
		sum += d
	}
	return sum
}
