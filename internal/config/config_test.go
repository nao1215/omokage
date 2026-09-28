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

// FuzzParseRoundTrip feeds arbitrary bytes to Parse, the reader of the
// user-edited omokage.toml. Parse may reject the input, but it must not panic,
// and whatever it accepts must survive Save's rendering: Parse(cfg.String())
// has to succeed and yield the same Config, otherwise `omokage init` style
// rewrites would silently change a project's name, default author, feature
// switches, or storage paths.
func FuzzParseRoundTrip(f *testing.F) {
	f.Add([]byte(Default("writing-lab").String()))
	f.Add([]byte("[project]\nname = \"x\"\n"))
	f.Add([]byte("[project]\nname = unquoted value # not a comment\n"))
	f.Add([]byte("[defaults]\ndefault_author = \"me\"\n"))
	f.Add([]byte("[features]\nkatakana_ratio = false\npos_ngram_frequency = 1\n"))
	f.Add([]byte("[storage]\nprofile_dir = \"./p\\n\\\"q\"\ncache_dir = ./c\r\n"))
	f.Add([]byte("[features]\nsentence_length = maybe\n"))
	f.Add([]byte("[project]\nname = \"unterminated\n"))
	f.Add([]byte("no equals sign\n"))
	f.Add([]byte("[project]\nname = \"\\xff\\u3042\"\n"))

	f.Fuzz(func(t *testing.T, data []byte) {
		cfg, err := Parse(data)
		if err != nil {
			return
		}
		rendered := cfg.String()
		again, err := Parse([]byte(rendered))
		if err != nil {
			t.Fatalf("Parse rejected its own rendering: %v\ninput: %q\nrendered: %q", err, data, rendered)
		}
		if again != cfg {
			t.Fatalf("round trip changed the config\ninput: %q\nfirst:  %+v\nsecond: %+v", data, cfg, again)
		}
		if again.String() != rendered {
			t.Fatalf("rendering is not stable\nfirst:  %q\nsecond: %q", rendered, again.String())
		}
	})
}
