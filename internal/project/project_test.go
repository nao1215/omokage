package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nao1215/dyer/internal/config"
)

func TestFindRootAndListProfiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	cfg, err := Init(root, "sample")
	if err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "content", "drafts")
	if err := os.MkdirAll(nested, 0o750); err != nil {
		t.Fatal(err)
	}

	found, err := FindRoot(nested)
	if err != nil {
		t.Fatal(err)
	}
	if found != root {
		t.Fatalf("unexpected root: got=%q want=%q", found, root)
	}

	profilePath, err := ProfilePath(root, cfg, "nao")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(profilePath, []byte("placeholder"), 0o600); err != nil {
		t.Fatal(err)
	}

	authors, err := ListProfiles(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 || authors[0] != "nao" {
		t.Fatalf("unexpected authors: %#v", authors)
	}
}

func TestInitRejectsExistingProject(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if _, err := Init(root, "sample"); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(root, "sample"); err == nil {
		t.Fatal("expected re-initialization to fail")
	}
}

func TestInitDefaultsEmptyName(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	cfg, err := Init(root, "   ")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Project.Name != filepath.Base(root) {
		t.Fatalf("expected name to default to directory base, got %q", cfg.Project.Name)
	}
}

func TestLoadOptionalWithoutProject(t *testing.T) {
	t.Parallel()

	_, found, err := LoadOptional(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("expected no project to be found")
	}
}

func TestLoadOptionalWithProject(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if _, err := Init(root, "sample"); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadOptional(root)
	if err != nil {
		t.Fatal(err)
	}
	if !found || cfg.Project.Name != "sample" {
		t.Fatalf("unexpected load result: found=%v name=%q", found, cfg.Project.Name)
	}
}

func TestLoadWithoutProject(t *testing.T) {
	t.Parallel()

	if _, _, err := Load(t.TempDir()); err == nil {
		t.Fatal("expected Load to fail without a project")
	}
}

func TestProfilePathValidation(t *testing.T) {
	t.Parallel()

	cfg := config.Default("sample")
	root := t.TempDir()

	if _, err := ProfilePath(root, cfg, "  "); err == nil {
		t.Fatal("expected error for an empty author")
	}
	if _, err := ProfilePath(root, cfg, filepath.Join("..", "escape")); err == nil {
		t.Fatal("expected error for an author containing path separators")
	}
}

func TestListProfilesWithoutDirectory(t *testing.T) {
	t.Parallel()

	cfg := config.Default("sample")
	authors, err := ListProfiles(t.TempDir(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 0 {
		t.Fatalf("expected no authors, got %#v", authors)
	}
}
