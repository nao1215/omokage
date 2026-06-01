package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "dyer.toml")
	expected := Default("writing-lab")
	expected.Features.KatakanaRatio = false

	if err := Save(path, expected); err != nil {
		t.Fatal(err)
	}

	actual, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if actual.Project.Name != expected.Project.Name {
		t.Fatalf("project name mismatch: got=%q want=%q", actual.Project.Name, expected.Project.Name)
	}
	if actual.Features.KatakanaRatio != expected.Features.KatakanaRatio {
		t.Fatalf("katakana ratio mismatch: got=%v want=%v", actual.Features.KatakanaRatio, expected.Features.KatakanaRatio)
	}
	if actual.Storage.ProfileDir != "./profiles" {
		t.Fatalf("unexpected profile dir: %q", actual.Storage.ProfileDir)
	}
}
