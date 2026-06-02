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
}

func TestLoadProfileFromEmptyDatabase(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "absent.db")
	if _, err := LoadProfile(path); err == nil {
		t.Fatal("expected error when loading a profile that was never trained")
	}
}
