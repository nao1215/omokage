package profile

import (
	"fmt"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
)

func benchDoc(paragraphs, salt int) string {
	const unit = "私は朝のしずかな時間が好きです。白湯を一杯ゆっくり飲みながら、窓の外をながめます。" +
		"On quiet mornings I drink a cup of warm water and watch the light change outside. " +
		"特別なことは何もしていません。それでも、こうした小さな習慣が一日をおだやかにする気がしています。\n\n"
	out := ""
	for range paragraphs {
		out += unit
	}
	return out + fmt.Sprintf("追記その%d。\n", salt)
}

func benchSetup() (feature.Distribution, feature.Metrics, feature.Metrics, config.Features) {
	corpus := make([]feature.Metrics, 0, 50)
	for i := range 50 {
		corpus = append(corpus, feature.ExtractText(benchDoc(8, i)))
	}
	dist := feature.Aggregate(corpus)
	target := feature.ExtractText(benchDoc(8, 999))
	other := feature.ExtractText(benchDoc(6, 1234))
	return dist, target, other, config.Default("bench").Features
}

func BenchmarkScore(b *testing.B) {
	dist, target, _, flags := benchSetup()
	b.ReportAllocs()
	for b.Loop() {
		_ = Score(dist, target, flags)
	}
}

func BenchmarkCompare(b *testing.B) {
	_, target, other, flags := benchSetup()
	b.ReportAllocs()
	for b.Loop() {
		_ = Compare(target, other, flags)
	}
}

// BenchmarkExplain measures the opt-in detailed path, including per-paragraph
// localization, against BenchmarkScore to bound how much --explain/--format json
// adds over the plain check.
func BenchmarkExplain(b *testing.B) {
	dist, target, _, flags := benchSetup()
	segments := feature.ExtractSegments(benchDoc(8, 999))
	b.ReportAllocs()
	for b.Loop() {
		_ = Explain(dist, target, segments, flags)
	}
}
