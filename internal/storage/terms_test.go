package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nao1215/omokage/internal/profile"
	"github.com/nao1215/omokage/internal/term"
)

func sampleTermProfile() term.Profile {
	return term.Profile{
		Groups: []term.Group{
			{
				GroupKey:         "term:ci",
				PreferredSurface: "継続的インテグレーション",
				TotalCount:       9,
				DocCount:         3,
				Variants: []term.Variant{
					{Surface: "継続的インテグレーション", NormalizedKey: "継続的インテグレーション", GroupKey: "term:ci", Count: 6, DocCount: 3},
					{Surface: "CI", NormalizedKey: "ci", GroupKey: "term:ci", Count: 3, DocCount: 2},
				},
			},
			{
				GroupKey:         "term:db",
				PreferredSurface: "DB",
				TotalCount:       5,
				DocCount:         2,
				Variants: []term.Variant{
					{Surface: "DB", NormalizedKey: "db", GroupKey: "term:db", Count: 4, DocCount: 2},
					{Surface: "ＤＢ", NormalizedKey: "db", GroupKey: "term:db", Count: 1, DocCount: 1},
				},
			},
		},
	}
}

// TestSaveLoadTerms round-trips a term profile through the database and checks
// every field survives, including the separation of normalized_key (db for both
// DB and ＤＢ) from group_key.
func TestSaveLoadTerms(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "nao.db")
	want := sampleTermProfile()
	if err := SaveTerms(path, want); err != nil {
		t.Fatalf("SaveTerms: %v", err)
	}

	got, err := LoadTerms(path)
	if err != nil {
		t.Fatalf("LoadTerms: %v", err)
	}
	if len(got.Groups) != len(want.Groups) {
		t.Fatalf("group count = %d, want %d", len(got.Groups), len(want.Groups))
	}
	for i, wg := range want.Groups {
		gg := got.Groups[i]
		if gg.GroupKey != wg.GroupKey || gg.PreferredSurface != wg.PreferredSurface ||
			gg.TotalCount != wg.TotalCount || gg.DocCount != wg.DocCount {
			t.Errorf("group %d = %+v, want %+v", i, gg, wg)
		}
		if len(gg.Variants) != len(wg.Variants) {
			t.Fatalf("group %d variant count = %d, want %d", i, len(gg.Variants), len(wg.Variants))
		}
		for j, wv := range wg.Variants {
			if gg.Variants[j] != wv {
				t.Errorf("group %d variant %d = %+v, want %+v", i, j, gg.Variants[j], wv)
			}
		}
	}
}

// TestSaveTermsReplacesPrevious checks that re-saving replaces the prior term
// data wholesale rather than accumulating stale groups.
func TestSaveTermsReplacesPrevious(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "nao.db")
	if err := SaveTerms(path, sampleTermProfile()); err != nil {
		t.Fatal(err)
	}
	replacement := term.Profile{Groups: []term.Group{{
		GroupKey: "term:api", PreferredSurface: "API", TotalCount: 2, DocCount: 1,
		Variants: []term.Variant{{Surface: "API", NormalizedKey: "api", GroupKey: "term:api", Count: 2, DocCount: 1}},
	}}}
	if err := SaveTerms(path, replacement); err != nil {
		t.Fatal(err)
	}
	got, err := LoadTerms(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Groups) != 1 || got.Groups[0].GroupKey != "term:api" {
		t.Fatalf("expected only the replacement group, got %+v", got.Groups)
	}
}

// TestLoadTermsMissingFile returns an empty profile (not an error) so show/check
// can run against an author trained before term support existed.
func TestLoadTermsMissingFile(t *testing.T) {
	t.Parallel()

	got, err := LoadTerms(filepath.Join(t.TempDir(), "absent.db"))
	if err != nil {
		t.Fatalf("LoadTerms on a missing file should not error: %v", err)
	}
	if len(got.Groups) != 0 {
		t.Fatalf("expected an empty profile, got %+v", got.Groups)
	}
}

// TestLoadTermsCorruptFile checks that a non-database file is reported as an
// error (the caller degrades to an empty list), not loaded as empty.
func TestLoadTermsCorruptFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "broken.db")
	if err := os.WriteFile(path, []byte("this is not a sqlite database"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadTerms(path); err == nil {
		t.Fatal("expected an error loading a corrupt database file")
	}
}

// TestLoadTermsFromProfileWithoutTerms checks that a database created by
// SaveProfile alone (no terms) loads as an empty term profile through the term
// tables added by the schema migration, never an error.
func TestLoadTermsFromProfileWithoutTerms(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "nao.db")
	record := profile.Record{
		Author:    "nao",
		SourceDir: "/tmp/corpus",
		TrainedAt: time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC),
		FileCount: 1,
	}
	if err := SaveProfile(path, record); err != nil {
		t.Fatal(err)
	}
	got, err := LoadTerms(path)
	if err != nil {
		t.Fatalf("LoadTerms on a terms-less profile should not error: %v", err)
	}
	if len(got.Groups) != 0 {
		t.Fatalf("expected an empty term profile, got %+v", got.Groups)
	}
}
