package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/storage"
)

// TestDiffGlobalWithoutStoreFallsBack guards the rule that diff needs no store:
// passing --global when no global store exists must fall back to the built-in
// feature weights and still produce a comparison, exactly like a bare diff, so a
// wrapper or habit that always adds --global is not a trap.
func TestDiffGlobalWithoutStoreFallsBack(t *testing.T) {
	workDir := t.TempDir()
	writeTestFile(t, filepath.Join(workDir, "a.md"), "First file with plain prose. It has two sentences.")
	writeTestFile(t, filepath.Join(workDir, "b.md"), "An entirely different second file. Its voice diverges.")
	home := filepath.Join(t.TempDir(), "no-global-store")

	code, stdout, stderr := runAppHome(t, workDir, home, "diff", "--global", "a.md", "b.md")
	if code != 0 {
		t.Fatalf("diff --global with no global store should fall back, got code=%d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) == "" {
		t.Fatalf("expected a comparison report on stdout, got empty output (stderr=%q)", stderr)
	}

	// It must match the bare diff: same files, same defaults, same report.
	bareCode, bareStdout, _ := runAppHome(t, workDir, home, "diff", "a.md", "b.md")
	if bareCode != 0 || bareStdout != stdout {
		t.Fatalf("diff --global should equal bare diff: code=%d report match=%v", bareCode, bareStdout == stdout)
	}
}

// TestInitNestedWarnsButSucceeds guards the nested-store footgun: init inside an
// existing project's subtree still creates the store (nesting is occasionally
// intentional) but warns, naming the enclosing config, so a user does not
// silently end up with a shadowing store that hides the profiles they expect.
func TestInitNestedWarnsButSucceeds(t *testing.T) {
	root := t.TempDir()
	if code, _, stderr := runApp(t, root, "init"); code != 0 {
		t.Fatalf("init at root failed: %s", stderr)
	}

	sub := filepath.Join(root, "drafts")
	writeTestFile(t, filepath.Join(sub, ".keep"), "")

	code, _, stderr := runApp(t, sub, "init")
	if code != 0 {
		t.Fatalf("nested init should still succeed, got code=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "nested inside an existing omokage store") {
		t.Fatalf("expected a nested-store warning, got stderr=%q", stderr)
	}
	if !strings.Contains(stderr, filepath.Join(root, "omokage.toml")) {
		t.Fatalf("warning should name the enclosing config, got stderr=%q", stderr)
	}
}

// TestInitFreshDoesNotWarn pins the other side: a plain init with no enclosing
// store stays quiet, so the nested warning never fires spuriously.
func TestInitFreshDoesNotWarn(t *testing.T) {
	code, _, stderr := runApp(t, t.TempDir(), "init")
	if code != 0 {
		t.Fatalf("fresh init failed: %s", stderr)
	}
	if strings.Contains(stderr, "nested") {
		t.Fatalf("fresh init must not warn about nesting, got stderr=%q", stderr)
	}
}

// TestShowDisplaysLocalTime checks the trained time is rendered in the local
// zone. The expectation is derived from the same stored timestamp, so the test
// holds in any timezone: it asserts the displayed string equals the local
// rendering of the stored (UTC) time, not a hard-coded zone.
func TestShowDisplaysLocalTime(t *testing.T) {
	workDir := trainedProject(t)

	record, err := storage.LoadProfile(filepath.Join(workDir, "profiles", "me.db"))
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}

	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me")
	if code != 0 {
		t.Fatalf("show failed: %s", stderr)
	}

	want := record.TrainedAt.Local().Format("2006-01-02 15:04:05 MST")
	if !strings.Contains(stdout, want) {
		t.Fatalf("expected local trained time %q in show output, got:\n%s", want, stdout)
	}
}
