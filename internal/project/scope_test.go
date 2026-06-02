package project

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePrefersLocalOverGlobal(t *testing.T) {
	t.Parallel()

	home := t.TempDir()
	if _, err := Init(home, "global"); err != nil {
		t.Fatal(err)
	}
	local := t.TempDir()
	if _, err := Init(local, "local"); err != nil {
		t.Fatal(err)
	}

	scope, err := Resolve(ResolveOptions{WorkDir: local, Home: home})
	if err != nil {
		t.Fatal(err)
	}
	if scope.Kind != ScopeLocal || scope.Root != local {
		t.Fatalf("expected local scope, got kind=%q root=%q", scope.Kind, scope.Root)
	}
}

func TestResolveFallsBackToGlobal(t *testing.T) {
	t.Parallel()

	home := t.TempDir()
	if _, err := Init(home, "global"); err != nil {
		t.Fatal(err)
	}
	// A working directory with no project above it falls back to the global store.
	work := t.TempDir()
	scope, err := Resolve(ResolveOptions{WorkDir: work, Home: home})
	if err != nil {
		t.Fatal(err)
	}
	if scope.Kind != ScopeGlobal || scope.Root != home {
		t.Fatalf("expected global scope, got kind=%q root=%q", scope.Kind, scope.Root)
	}
}

func TestResolveNoStoreIsNotFound(t *testing.T) {
	t.Parallel()

	// No local project, and a home that has no config: this must stay the original
	// "project not found" error rather than inventing a store.
	if _, err := Resolve(ResolveOptions{WorkDir: t.TempDir(), Home: t.TempDir()}); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
	// An empty home also yields not-found (the fallback is disabled).
	if _, err := Resolve(ResolveOptions{WorkDir: t.TempDir(), Home: ""}); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound with no home, got %v", err)
	}
}

func TestResolveGlobalFlagForces(t *testing.T) {
	t.Parallel()

	home := t.TempDir()
	if _, err := Init(home, "global"); err != nil {
		t.Fatal(err)
	}
	local := t.TempDir()
	if _, err := Init(local, "local"); err != nil {
		t.Fatal(err)
	}

	// Even inside a local project, --global forces the global store.
	scope, err := Resolve(ResolveOptions{WorkDir: local, Home: home, Global: true})
	if err != nil {
		t.Fatal(err)
	}
	if scope.Kind != ScopeGlobal {
		t.Fatalf("expected --global to force global, got %q", scope.Kind)
	}
}

func TestResolveExplicitProfileDir(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(t.TempDir(), "profiles")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	scope, err := Resolve(ResolveOptions{WorkDir: t.TempDir(), ProfileDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	if scope.Kind != ScopeExplicit || scope.ProfileDir != dir {
		t.Fatalf("unexpected explicit scope: kind=%q dir=%q", scope.Kind, scope.ProfileDir)
	}
	// A bare --profile-dir has no config file to write defaults into.
	if scope.ConfigPath != "" {
		t.Fatalf("explicit profile-dir scope should have no config path, got %q", scope.ConfigPath)
	}
}

func TestResolveExplicitConfig(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if _, err := Init(root, "explicit"); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(root, ConfigFileName)
	scope, err := Resolve(ResolveOptions{WorkDir: t.TempDir(), ConfigPath: cfgPath})
	if err != nil {
		t.Fatal(err)
	}
	if scope.Kind != ScopeExplicit {
		t.Fatalf("expected explicit scope, got %q", scope.Kind)
	}
	if scope.ProfileDir != filepath.Join(root, "profiles") {
		t.Fatalf("explicit config should resolve its profile dir, got %q", scope.ProfileDir)
	}
}
