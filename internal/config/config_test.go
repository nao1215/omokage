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
	expected.Features.LexicalFrequency = false
	expected.Features.CharNgramFrequency = false

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
	if actual.Features.LexicalFrequency != expected.Features.LexicalFrequency {
		t.Fatalf("lexical frequency mismatch: got=%v want=%v", actual.Features.LexicalFrequency, expected.Features.LexicalFrequency)
	}
	if actual.Features.CharNgramFrequency != expected.Features.CharNgramFrequency {
		t.Fatalf("char n-gram frequency mismatch: got=%v want=%v", actual.Features.CharNgramFrequency, expected.Features.CharNgramFrequency)
	}
	if !Default("x").Features.LexicalFrequency || !Default("x").Features.CharNgramFrequency {
		t.Fatal("expected the new authorship features to default to enabled")
	}
	if actual.Storage.ProfileDir != "./profiles" {
		t.Fatalf("unexpected profile dir: %q", actual.Storage.ProfileDir)
	}
}
