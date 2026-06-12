package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/profile"
)

func TestSaveLoadProfile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "nao.db")
	expected := profile.Record{
		Author:    "nao",
		SourceDir: "/tmp/corpus",
		TrainedAt: time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC),
		FileCount: 3,
		SelfSimilarity: &profile.SelfSimilarityStats{
			MeanZ:       []float64{0.4, 0.6, 0.8},
			MeanZMedian: 0.6,
			MeanZSpread: 0.1632993161855452,
			MeanZMin:    0.4,
			MeanZMax:    0.8,
		},
		Distribution: feature.Distribution{
			Mean: feature.Metrics{
				AverageSentenceLength:    12,
				SentenceLengthVariance:   5,
				PunctuationFrequency:     0.12,
				NewlineFrequency:         0.03,
				BulletRatio:              0.1,
				ConjunctionFrequency:     0.04,
				KanjiRatio:               0.3,
				HiraganaRatio:            0.5,
				KatakanaRatio:            0.2,
				ParagraphLengthVariance:  7,
				MarkdownStructureDensity: 0.15,
				PoliteEndingRatio:        0.82,
				PlainEndingRatio:         0.04,
				LexicalFrequencies:       map[string]float64{"の": 0.08, "the": 0.0},
				CharNgrams:               map[string]float64{"です": 0.01, "th": 0.0},
			},
			StdDev: feature.Metrics{
				AverageSentenceLength:    2.5,
				SentenceLengthVariance:   1.5,
				PunctuationFrequency:     0.02,
				NewlineFrequency:         0.01,
				BulletRatio:              0.05,
				ConjunctionFrequency:     0.01,
				KanjiRatio:               0.04,
				HiraganaRatio:            0.06,
				KatakanaRatio:            0.03,
				ParagraphLengthVariance:  2.0,
				MarkdownStructureDensity: 0.05,
				PoliteEndingRatio:        0.12,
				PlainEndingRatio:         0.03,
				LexicalFrequencies:       map[string]float64{"の": 0.01, "the": 0.0},
				CharNgrams:               map[string]float64{"です": 0.002, "th": 0.0},
			},
			DocumentCount:  3,
			SentenceCount:  8,
			CharacterCount: 120,
		},
	}

	if err := SaveProfile(path, expected); err != nil {
		t.Fatal(err)
	}

	actual, err := LoadProfile(path)
	if err != nil {
		t.Fatal(err)
	}

	if actual.Author != expected.Author {
		t.Fatalf("author mismatch: got=%q want=%q", actual.Author, expected.Author)
	}
	if actual.FileCount != expected.FileCount {
		t.Fatalf("file count mismatch: got=%d want=%d", actual.FileCount, expected.FileCount)
	}
	if actual.Distribution.CharacterCount != expected.Distribution.CharacterCount {
		t.Fatalf("character count mismatch: got=%d want=%d", actual.Distribution.CharacterCount, expected.Distribution.CharacterCount)
	}
	if actual.Distribution.Mean.KanjiRatio != expected.Distribution.Mean.KanjiRatio {
		t.Fatalf("mean kanji ratio mismatch: got=%f want=%f", actual.Distribution.Mean.KanjiRatio, expected.Distribution.Mean.KanjiRatio)
	}
	if actual.Distribution.StdDev.KanjiRatio != expected.Distribution.StdDev.KanjiRatio {
		t.Fatalf("std kanji ratio mismatch: got=%f want=%f", actual.Distribution.StdDev.KanjiRatio, expected.Distribution.StdDev.KanjiRatio)
	}
	if actual.Distribution.Mean.PoliteEndingRatio != expected.Distribution.Mean.PoliteEndingRatio {
		t.Fatalf("mean polite ending ratio mismatch: got=%f want=%f", actual.Distribution.Mean.PoliteEndingRatio, expected.Distribution.Mean.PoliteEndingRatio)
	}
	if actual.Distribution.StdDev.PlainEndingRatio != expected.Distribution.StdDev.PlainEndingRatio {
		t.Fatalf("std plain ending ratio mismatch: got=%f want=%f", actual.Distribution.StdDev.PlainEndingRatio, expected.Distribution.StdDev.PlainEndingRatio)
	}
	if !actual.TrainedAt.Equal(expected.TrainedAt) {
		t.Fatalf("trained_at mismatch: got=%s want=%s", actual.TrainedAt, expected.TrainedAt)
	}
	if got := actual.Distribution.Mean.LexicalFrequencies["の"]; got != 0.08 {
		t.Fatalf("mean lexical frequency mismatch: got=%f want=%f", got, 0.08)
	}
	if got := actual.Distribution.StdDev.LexicalFrequencies["の"]; got != 0.01 {
		t.Fatalf("std lexical frequency mismatch: got=%f want=%f", got, 0.01)
	}
	if got := actual.Distribution.Mean.CharNgrams["です"]; got != 0.01 {
		t.Fatalf("mean char n-gram mismatch: got=%f want=%f", got, 0.01)
	}
	if got := actual.Distribution.StdDev.CharNgrams["です"]; got != 0.002 {
		t.Fatalf("std char n-gram mismatch: got=%f want=%f", got, 0.002)
	}
	if actual.SelfSimilarity == nil {
		t.Fatal("expected self-similarity stats to round-trip")
	}
	if actual.SelfSimilarity.MeanZMedian != expected.SelfSimilarity.MeanZMedian {
		t.Fatalf("self-similarity median mismatch: got=%f want=%f", actual.SelfSimilarity.MeanZMedian, expected.SelfSimilarity.MeanZMedian)
	}
	if len(actual.SelfSimilarity.MeanZ) != len(expected.SelfSimilarity.MeanZ) {
		t.Fatalf("self-similarity sample count mismatch: got=%d want=%d", len(actual.SelfSimilarity.MeanZ), len(expected.SelfSimilarity.MeanZ))
	}
}

func TestSaveLoadProfileRoundTripsSources(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "multi.db")
	record := profile.Record{
		Author:    "me",
		SourceDir: "/tmp/posts",
		Sources:   []string{"/tmp/posts", "/tmp/draft.md", "/tmp/note.txt"},
		TrainedAt: time.Date(2026, time.June, 2, 9, 0, 0, 0, time.UTC),
		FileCount: 5,
	}
	if err := SaveProfile(path, record); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Sources) != 3 {
		t.Fatalf("expected 3 sources, got %d: %v", len(loaded.Sources), loaded.Sources)
	}
	for i, want := range record.Sources {
		if loaded.Sources[i] != want {
			t.Fatalf("source %d mismatch: got=%q want=%q", i, loaded.Sources[i], want)
		}
	}
}

// A profile saved with no Sources (e.g. one trained before multi-input support,
// or any record that only set SourceDir) must load with Sources populated from
// SourceDir so every consumer can rely on it.
func TestLoadProfileFallsBackToSourceDir(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "legacy.db")
	record := profile.Record{
		Author:    "me",
		SourceDir: "/tmp/corpus",
		TrainedAt: time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC),
		FileCount: 3,
	}
	if err := SaveProfile(path, record); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Sources) != 1 || loaded.Sources[0] != "/tmp/corpus" {
		t.Fatalf("expected Sources to fall back to [SourceDir], got %v", loaded.Sources)
	}
}

func TestLoadProfileFromEmptyDatabase(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "absent.db")
	if _, err := LoadProfile(path); err == nil {
		t.Fatal("expected error when loading a profile that was never trained")
	}
}
