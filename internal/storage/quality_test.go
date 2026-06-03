package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/nao1215/omokage/internal/profile"
)

// minimalRecord is the smallest valid profile to attach findings to.
func minimalRecord() profile.Record {
	return profile.Record{
		Author:    "me",
		TrainedAt: time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC),
		FileCount: 3,
	}
}

// TestSaveLoadQualityFindingsRoundTrip verifies stored findings read back byte-for-byte.
func TestSaveLoadQualityFindingsRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "profiles", "me.db")
	if err := SaveProfile(path, minimalRecord()); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	want := []byte(`[{"code":"few_documents","severity":"warning"}]`)
	if err := SaveQualityFindings(path, want); err != nil {
		t.Fatalf("save findings: %v", err)
	}
	got, err := LoadQualityFindings(path)
	if err != nil {
		t.Fatalf("load findings: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("round-trip mismatch:\n want %s\n  got %s", want, got)
	}
}

// TestLoadQualityFindingsDefaultsToEmpty verifies a profile without stored findings loads the empty default.
func TestLoadQualityFindingsDefaultsToEmpty(t *testing.T) {
	t.Parallel()

	// A profile trained before findings were stored (the column carries its default)
	// loads as an empty array, never an error, so show degrades to "no findings".
	path := filepath.Join(t.TempDir(), "profiles", "me.db")
	if err := SaveProfile(path, minimalRecord()); err != nil {
		t.Fatalf("save profile: %v", err)
	}
	got, err := LoadQualityFindings(path)
	if err != nil {
		t.Fatalf("load findings: %v", err)
	}
	if string(got) != "[]" {
		t.Fatalf("expected the empty-array default, got %s", got)
	}
}

// TestLoadQualityFindingsMissingProfile verifies a missing database yields the empty default, not an error.
func TestLoadQualityFindingsMissingProfile(t *testing.T) {
	t.Parallel()

	// No database file at all yields the empty default rather than an error.
	got, err := LoadQualityFindings(filepath.Join(t.TempDir(), "absent.db"))
	if err != nil {
		t.Fatalf("a missing profile should not error: %v", err)
	}
	if string(got) != "[]" {
		t.Fatalf("expected the empty-array default, got %s", got)
	}
}

// TestSaveQualityFindingsWithoutProfileFails verifies attaching findings before the profile row exists is reported.
func TestSaveQualityFindingsWithoutProfileFails(t *testing.T) {
	t.Parallel()

	// Attaching findings before the profile row exists is a programming error, and
	// must be reported rather than silently storing nothing.
	path := filepath.Join(t.TempDir(), "profiles", "me.db")
	if err := SaveQualityFindings(path, []byte("[]")); err == nil {
		t.Fatal("expected an error when no profile row exists yet")
	}
}
