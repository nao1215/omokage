package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "omokage.toml")
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

func TestDefaultAuthorRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "omokage.toml")
	cfg := Default("writing-lab")
	cfg.Defaults.Author = "me"
	if err := Save(path, cfg); err != nil {
		t.Fatal(err)
	}
	actual, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if actual.Defaults.Author != "me" {
		t.Fatalf("default author not preserved: got %q", actual.Defaults.Author)
	}

	// A config without a [defaults] section parses to an empty default author,
	// preserving backward compatibility with files written before the field.
	legacy, err := Parse([]byte("[project]\nname = \"x\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	if legacy.Defaults.Author != "" {
		t.Fatalf("legacy config should have no default author, got %q", legacy.Defaults.Author)
	}
}
