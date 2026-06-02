package feature

import (
	"fmt"
	"testing"
)

// benchDoc is a mixed Japanese/English paragraph, repeated to a realistic blog
// post length. It exercises sentence splitting, script counting, function-word
// counting, and bigram/trigram extraction together.
func benchDoc(paragraphs int) string {
	const unit = "私は朝のしずかな時間が好きです。白湯を一杯ゆっくり飲みながら、窓の外をながめます。" +
		"On quiet mornings I drink a cup of warm water and watch the light change outside. " +
		"特別なことは何もしていません。それでも、こうした小さな習慣が一日をおだやかにする気がしています。\n\n"
	out := ""
	for range paragraphs {
		out += unit
	}
	return out
}

func benchCorpus(docs, paragraphs int) []Metrics {
	corpus := make([]Metrics, 0, docs)
	for i := range docs {
		// Vary each document slightly so the aggregated spread is non-degenerate.
		text := benchDoc(paragraphs) + fmt.Sprintf("追記その%d。\n", i)
		corpus = append(corpus, ExtractText(text))
	}
	return corpus
}

func BenchmarkExtractText(b *testing.B) {
	doc := benchDoc(8)
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractText(doc)
	}
}

func BenchmarkAggregate(b *testing.B) {
	corpus := benchCorpus(50, 8)
	b.ReportAllocs()
	for b.Loop() {
		_ = Aggregate(corpus)
	}
}
