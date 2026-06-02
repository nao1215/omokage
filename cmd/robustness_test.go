package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
)

func TestSubcommandHelpExitsZero(t *testing.T) {
	t.Parallel()

	// An explicit --help is a successful request for usage on every subcommand.
	for _, name := range []string{"init", "train", "check", "diff", "list", "show", "remove", "rename"} {
		if code, _, _ := runApp(t, t.TempDir(), name, "--help"); code != 0 {
			t.Fatalf("%q --help should exit 0, got %d", name, code)
		}
	}
	// A genuine parse error (unknown flag) still exits non-zero, so automation can
	// tell "show me the help" apart from "you used me wrong".
	if code, _, _ := runApp(t, t.TempDir(), "check", "--bogus"); code == 0 {
		t.Fatal("an unknown flag should exit non-zero")
	}
}

// readonlyConfig sets a project up with author "me" as the default, then makes
// omokage.toml read-only so any later config write fails, and returns the config
// path. It backs the atomicity tests for the default-author mutations.
func readonlyConfig(t *testing.T, workDir string) string {
	t.Helper()
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "--default", "posts"); code != 0 {
		t.Fatalf("seeding the default failed: %s", stderr)
	}
	configPath := filepath.Join(workDir, "omokage.toml")
	if err := os.Chmod(configPath, 0o400); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(configPath, 0o600) })
	return configPath
}

func defaultAuthor(t *testing.T, configPath string) string {
	t.Helper()
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	return cfg.Defaults.Author
}

func TestTrainDefaultReportsPartialOnConfigFailure(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	configPath := filepath.Join(workDir, "omokage.toml")
	if err := os.Chmod(configPath, 0o400); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(configPath, 0o600) })

	// Re-train "me" with --default against a read-only config. The profile is the
	// command's purpose and must still be written; only the default can't be set,
	// and the output must say so rather than claim the default was recorded.
	code, stdout, stderr := runApp(t, workDir, "train", "--author", "me", "--default", "posts")
	if code != 1 {
		t.Fatalf("expected exit 1 when the default cannot be saved, got %d", code)
	}
	if !strings.Contains(stdout, `Trained author "me"`) {
		t.Fatalf("the profile was trained, so that should be reported: %q", stdout)
	}
	if strings.Contains(stdout, "Default author set") {
		t.Fatalf("must not claim the default was set when it failed: %q", stdout)
	}
	if !strings.Contains(stderr, "warning") {
		t.Fatalf("expected a warning about the failed default: %q", stderr)
	}
	if got := defaultAuthor(t, configPath); got != "" {
		t.Fatalf("the default must not be recorded on failure, got %q", got)
	}
}

func TestRemoveLeavesStoreIntactOnConfigFailure(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	configPath := readonlyConfig(t, workDir)

	// Removing the default author needs a config write to clear it. That write is
	// done first, so when it fails nothing is deleted: the profile survives and the
	// default still points at it. No half-applied state, no dangling default.
	code, _, stderr := runApp(t, workDir, "remove", "--author", "me")
	if code != 1 {
		t.Fatalf("expected exit 1 when the config cannot be updated, got %d", code)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Fatal("expected an error explaining the failure")
	}
	if _, err := os.Stat(filepath.Join(workDir, "profiles", "me.db")); err != nil {
		t.Fatalf("the profile must not be deleted when clearing the default failed: %v", err)
	}
	if got := defaultAuthor(t, configPath); got != "me" {
		t.Fatalf("the default should be unchanged, got %q", got)
	}
}

func TestRemoveRollsBackDefaultOnDeleteFailure(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "--default", "posts"); code != 0 {
		t.Fatalf("seeding the default failed: %s", stderr)
	}
	configPath := filepath.Join(workDir, "omokage.toml")

	// Make the profiles directory non-writable so os.Remove fails while the config
	// stays writable. The default is cleared first, the delete then fails, and the
	// default must be restored — never left cleared with the profile still present.
	profilesDir := filepath.Join(workDir, "profiles")
	// A directory needs its execute bit to be traversable, so 0500/0700 (not the
	// 0600 gosec prefers for files) are the right modes for a read-only then
	// restored directory fixture.
	if err := os.Chmod(profilesDir, 0o500); err != nil { //nolint:gosec // directory perms need the execute bit
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(profilesDir, 0o700) }) //nolint:gosec // directory perms need the execute bit

	code, _, stderr := runApp(t, workDir, "remove", "--author", "me")
	if code != 1 {
		t.Fatalf("expected exit 1 when the profile cannot be deleted, got %d", code)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Fatal("expected an error explaining the failure")
	}
	if _, err := os.Stat(filepath.Join(profilesDir, "me.db")); err != nil {
		t.Fatalf("the profile should still be present after a failed delete: %v", err)
	}
	if got := defaultAuthor(t, configPath); got != "me" {
		t.Fatalf("the default must be rolled back, not left cleared, got %q", got)
	}
}

func TestProfileDirAmbiguityMentionsAuthorOnly(t *testing.T) {
	t.Parallel()

	// A bare --profile-dir scope has no config file, so the ambiguity error must
	// not suggest setting a default there — only --author can disambiguate.
	workDir := t.TempDir()
	profileDir := filepath.Join(t.TempDir(), "profiles")
	if err := os.MkdirAll(profileDir, 0o750); err != nil {
		t.Fatal(err)
	}
	jaCorpus(t, filepath.Join(workDir, "posts"))
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。")
	for _, name := range []string{"alpha", "beta"} {
		if code, _, stderr := runApp(t, workDir, "train", "--profile-dir", profileDir, "--author", name, "posts"); code != 0 {
			t.Fatalf("train %s failed: %s", name, stderr)
		}
	}

	code, _, stderr := runApp(t, workDir, "check", "--profile-dir", profileDir, "draft.txt")
	if code == 0 {
		t.Fatal("expected an ambiguity error with no --author")
	}
	if !strings.Contains(stderr, "--author") || !strings.Contains(stderr, "no config file") {
		t.Fatalf("error should steer to --author and note the missing config: %q", stderr)
	}
	if strings.Contains(stderr, "set a default") || strings.Contains(stderr, "default_author") {
		t.Fatalf("must not suggest an impossible default in a --profile-dir scope: %q", stderr)
	}
}

func TestRenameRollsBackOnConfigFailure(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	configPath := readonlyConfig(t, workDir)

	// Renaming the default author updates the config. When that write fails, the
	// just-created profile is rolled back so the store returns to exactly its prior
	// state: the old profile remains, the new one does not, and the default is
	// unchanged.
	code, _, stderr := runApp(t, workDir, "rename", "--author", "me", "--to", "watashi")
	if code != 1 {
		t.Fatalf("expected exit 1 when the config cannot be updated, got %d", code)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Fatal("expected an error explaining the failure")
	}
	if _, err := os.Stat(filepath.Join(workDir, "profiles", "me.db")); err != nil {
		t.Fatalf("the original profile must survive a failed rename: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workDir, "profiles", "watashi.db")); !os.IsNotExist(err) {
		t.Fatalf("the new profile must be rolled back on failure, stat err=%v", err)
	}
	if got := defaultAuthor(t, configPath); got != "me" {
		t.Fatalf("the default should still point at the original, got %q", got)
	}
}
