package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nao1215/omokage/internal/config"
)

// Scope is the resolved place a command reads and writes profiles. It unifies
// the three ways omokage can be pointed at a profile store — a local project
// found by walking up from the working directory, a global store under
// OMOKAGE_HOME, or an explicit --config/--profile-dir override — behind a single
// value so the command layer never has to special-case them.
//
// ProfileDir and CacheDir are absolute so they are stable regardless of the
// working directory. ConfigPath is the omokage.toml backing this scope, or "" for
// a bare --profile-dir scope that has no config file to write defaults into.
type Scope struct {
	// Kind is "local", "global", or "explicit"; used only for human-facing hints.
	Kind       string
	Root       string
	ConfigPath string
	ProfileDir string
	CacheDir   string
	Config     config.Config
}

const (
	// ScopeLocal is a project found by walking up from the working directory.
	ScopeLocal = "local"
	// ScopeGlobal is the per-user store under OMOKAGE_HOME.
	ScopeGlobal = "global"
	// ScopeExplicit is a store named directly with --config or --profile-dir.
	ScopeExplicit = "explicit"
)

// ResolveOptions describes how to locate the active scope. The precedence is
// fixed in Resolve; these are the inputs that drive it.
type ResolveOptions struct {
	// WorkDir is the directory the command was invoked from.
	WorkDir string
	// Home is the global store base directory (OMOKAGE_HOME or its default). An
	// empty Home disables the global fallback entirely, which keeps tests and
	// project-only setups from ever touching a user-wide store by accident.
	Home string
	// Global forces the global store, skipping the upward project search.
	Global bool
	// ConfigPath names an omokage.toml to use directly.
	ConfigPath string
	// ProfileDir names a profile directory to use directly.
	ProfileDir string
}

// Resolve picks the active scope using a fixed precedence:
//
//  1. --config / --profile-dir (explicit) — the caller named the store.
//  2. --global — the per-user store under Home.
//  3. a local project found by walking up from WorkDir — inside a project, local
//     always wins, so the existing project-local model is untouched.
//  4. the global store, but only when a global config already exists — this is
//     the "use omokage anywhere" convenience and never silently invents a store.
//  5. otherwise ErrProjectNotFound, exactly as before.
func Resolve(o ResolveOptions) (Scope, error) {
	if o.ConfigPath != "" || o.ProfileDir != "" {
		return resolveExplicit(o)
	}
	if o.Global {
		return resolveGlobal(o.Home)
	}

	root, err := FindRoot(o.WorkDir)
	if err == nil {
		cfg, loadErr := config.Load(filepath.Join(root, ConfigFileName))
		if loadErr != nil {
			return Scope{}, loadErr
		}
		return scopeFromRoot(ScopeLocal, root, cfg), nil
	}
	if !errors.Is(err, ErrProjectNotFound) {
		return Scope{}, err
	}

	if o.Home != "" {
		if _, statErr := os.Stat(filepath.Join(o.Home, ConfigFileName)); statErr == nil {
			return resolveGlobal(o.Home)
		}
	}
	return Scope{}, ErrProjectNotFound
}

func resolveExplicit(o ResolveOptions) (Scope, error) {
	scope := Scope{Kind: ScopeExplicit, Config: config.Default("omokage")}

	if o.ConfigPath != "" {
		cfg, err := config.Load(o.ConfigPath)
		if err != nil {
			return Scope{}, err
		}
		abs, err := filepath.Abs(o.ConfigPath)
		if err != nil {
			return Scope{}, err
		}
		base := filepath.Dir(abs)
		scope.Config = cfg
		scope.ConfigPath = abs
		scope.Root = base
		scope.ProfileDir = absDir(base, cfg.Storage.ProfileDir)
		scope.CacheDir = absDir(base, cfg.Storage.CacheDir)
	}

	if o.ProfileDir != "" {
		abs, err := filepath.Abs(o.ProfileDir)
		if err != nil {
			return Scope{}, err
		}
		// An explicit --profile-dir wins over whatever a --config pointed at, so the
		// two flags compose: --config supplies the feature set and default author,
		// --profile-dir relocates only the profiles.
		scope.ProfileDir = abs
	}

	if scope.ProfileDir == "" {
		return Scope{}, errors.New("--config or --profile-dir did not yield a profile directory")
	}
	return scope, nil
}

func resolveGlobal(home string) (Scope, error) {
	if home == "" {
		return Scope{}, errors.New("global store location is unknown; set OMOKAGE_HOME")
	}
	cfgPath := filepath.Join(home, ConfigFileName)
	cfg, err := config.Load(cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Scope{}, fmt.Errorf("no global omokage store at %s; run 'omokage init --global'", home)
		}
		return Scope{}, err
	}
	return scopeFromRoot(ScopeGlobal, home, cfg), nil
}

func scopeFromRoot(kind, root string, cfg config.Config) Scope {
	return Scope{
		Kind:       kind,
		Root:       root,
		ConfigPath: filepath.Join(root, ConfigFileName),
		ProfileDir: absDir(root, cfg.Storage.ProfileDir),
		CacheDir:   absDir(root, cfg.Storage.CacheDir),
		Config:     cfg,
	}
}

// absDir resolves a possibly-relative storage directory against base, so the
// "./profiles" stored in a config becomes an absolute path under its project.
func absDir(base, dir string) string {
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir)
	}
	return filepath.Clean(filepath.Join(base, dir))
}

// ProfilePath returns the SQLite path for an author inside this scope. The author
// is validated the same way ProfilePath(root, cfg, author) is, so a name with a
// path separator is rejected before it can escape the profile directory.
func (s Scope) ProfilePath(author string) (string, error) {
	safe := strings.TrimSpace(author)
	if safe == "" {
		return "", errors.New("author must not be empty")
	}
	if safe != filepath.Base(safe) {
		return "", errors.New("author must not contain path separators")
	}
	return filepath.Join(s.ProfileDir, safe+".db"), nil
}

// ListProfiles returns the author names trained in this scope, sorted.
func (s Scope) ListProfiles() ([]string, error) {
	return listProfilesInDir(s.ProfileDir)
}

func listProfilesInDir(profileDir string) ([]string, error) {
	entries, err := os.ReadDir(profileDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	authors := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".db") {
			authors = append(authors, strings.TrimSuffix(name, filepath.Ext(name)))
		}
	}
	sort.Strings(authors)
	return authors, nil
}

// SaveConfig persists the scope's config back to its omokage.toml. It is used by
// the commands that mutate scope-wide state (train --default, remove, rename) so
// the default author stays in sync without the user editing the file by hand. A
// scope with no config file (a bare --profile-dir) has nowhere to write to.
func (s Scope) SaveConfig() error {
	if s.ConfigPath == "" {
		return errors.New("this scope has no config file to update (use --config or a project/global store)")
	}
	return config.Save(s.ConfigPath, s.Config)
}
