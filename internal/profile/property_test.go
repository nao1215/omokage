package profile

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
	"pgregory.net/rapid"
)

// proseGen draws a short document from Japanese and English sentence pieces,
// punctuation, list items and blank lines, so the generated texts differ in
// every feature omokage measures.
func proseGen() *rapid.Generator[string] {
	pieces := []string{
		"今日は晴れです。", "私はGoを書いた。", "しかし、問題がある。", "そして結果を確認します。",
		"テストを追加しました！", "カタカナのサンプル。", "This is a test.", "However, it works.",
		"And then we ship it.", "- a bullet item", "1. first step", "## Heading", "\n", "\n\n", "、", "，",
	}
	return rapid.Custom(func(t *rapid.T) string {
		parts := rapid.SliceOfN(rapid.SampledFrom(pieces), 1, 12).Draw(t, "pieces")
		return strings.Join(parts, "")
	})
}

func metricsGen() *rapid.Generator[feature.Metrics] {
	return rapid.Custom(func(t *rapid.T) feature.Metrics {
		return feature.ExtractText(proseGen().Draw(t, "text"))
	})
}

func defaultFlags() config.Features {
	return config.Default("p").Features
}

// TestPropertyCompareWithItselfIsFullySimilar checks that a document compared
// with itself scores 100%.
func TestPropertyCompareWithItselfIsFullySimilar(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		m := metricsGen().Draw(t, "metrics")
		got := Compare(m, m, defaultFlags())
		if got.Similarity != 100 {
			t.Fatalf("Compare(m, m).Similarity = %d, differences %v", got.Similarity, got.Differences)
		}
	})
}

// TestPropertySimilarityIsAPercentage checks that every comparison and score
// stays within 0..100.
func TestPropertySimilarityIsAPercentage(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		a := metricsGen().Draw(t, "a")
		b := metricsGen().Draw(t, "b")
		corpus := rapid.SliceOfN(metricsGen(), 1, 5).Draw(t, "corpus")
		flags := defaultFlags()
		for name, s := range map[string]int{
			"Compare": Compare(a, b, flags).Similarity,
			"Score":   Score(feature.Aggregate(corpus), b, flags).Similarity,
		} {
			if s < 0 || s > 100 {
				t.Fatalf("%s similarity %d outside 0..100", name, s)
			}
		}
	})
}

// TestPropertyAggregateIgnoresDocumentOrder is a metamorphic relation: the
// profile of a corpus does not depend on the order its documents are read in,
// so neither does the score of a draft against it.
func TestPropertyAggregateIgnoresDocumentOrder(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		corpus := rapid.SliceOfN(metricsGen(), 2, 6).Draw(t, "corpus")
		perm := rapid.Permutation(corpus).Draw(t, "permutation")
		target := metricsGen().Draw(t, "target")

		a := feature.Aggregate(corpus)
		b := feature.Aggregate(perm)
		if !distributionsClose(a, b) {
			t.Fatalf("aggregate depends on the order:\n%+v\n%+v", a, b)
		}
		flags := defaultFlags()
		if sa, sb := Score(a, target, flags).Similarity, Score(b, target, flags).Similarity; sa != sb {
			t.Fatalf("score %d against the corpus, %d against the same corpus reordered", sa, sb)
		}
	})
}

// distributionsClose compares two distributions field by field, allowing for
// the rounding of floating-point sums taken in another order.
func distributionsClose(a, b feature.Distribution) bool {
	return valuesClose(reflect.ValueOf(a), reflect.ValueOf(b))
}

func valuesClose(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Float32, reflect.Float64:
		x, y := a.Float(), b.Float()
		return math.Abs(x-y) <= 1e-9*math.Max(1, math.Max(math.Abs(x), math.Abs(y)))
	case reflect.Struct:
		for i := range a.NumField() {
			if !valuesClose(a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}
		for _, k := range a.MapKeys() {
			bv := b.MapIndex(k)
			if !bv.IsValid() || !valuesClose(a.MapIndex(k), bv) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}
		for i := range a.Len() {
			if !valuesClose(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
}

// TestPropertyADraftEqualToEveryDocumentScoresFull checks the score against a
// corpus whose documents are all the draft itself: there is no spread and no
// distance, so the draft must score 100%.
func TestPropertyADraftEqualToEveryDocumentScoresFull(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		m := metricsGen().Draw(t, "metrics")
		n := rapid.IntRange(1, 4).Draw(t, "copies")
		corpus := make([]feature.Metrics, n)
		for i := range corpus {
			corpus[i] = m
		}
		got := Score(feature.Aggregate(corpus), m, defaultFlags())
		if got.Similarity != 100 {
			t.Fatalf("Score = %d with %d copies, differences %v", got.Similarity, n, got.Differences)
		}
	})
}

// TestPropertyLineEndingsDoNotChangeTheMetrics is a metamorphic relation: the
// same text saved with Windows line endings measures the same.
func TestPropertyLineEndingsDoNotChangeTheMetrics(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		text := proseGen().Draw(t, "text")
		lf := feature.ExtractText(text)
		crlf := feature.ExtractText(strings.ReplaceAll(text, "\n", "\r\n"))
		if !valuesClose(reflect.ValueOf(lf), reflect.ValueOf(crlf)) {
			t.Fatalf("metrics differ for CRLF:\n%+v\n%+v", lf, crlf)
		}
	})
}

// TestPropertyAByteOrderMarkDoesNotChangeTheMetrics is a metamorphic relation:
// a file saved with a UTF-8 byte order mark, as some Windows editors do,
// measures the same as without it.
func TestPropertyAByteOrderMarkDoesNotChangeTheMetrics(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		text := proseGen().Draw(t, "text")
		plain := feature.ExtractText(text)
		bom := feature.ExtractText(string([]byte{0xEF, 0xBB, 0xBF}) + text)
		if !valuesClose(reflect.ValueOf(plain), reflect.ValueOf(bom)) {
			t.Fatalf("metrics differ with a BOM:\n%+v\n%+v", plain, bom)
		}
	})
}
