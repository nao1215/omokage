package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// jaCorpus seeds a directory with a small polite-register Japanese corpus, the
// same shape trainedProject uses, so the ergonomics tests can stand up extra
// authors without repeating the fixtures.
func jaCorpus(t *testing.T, dir string) {
	t.Helper()
	writeTestFile(t, filepath.Join(dir, "one.txt"), "今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。")
	writeTestFile(t, filepath.Join(dir, "two.txt"), "昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。")
	writeTestFile(t, filepath.Join(dir, "three.txt"), "週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。")
}

func TestCheckSingleProfileNeedsNoAuthor(t *testing.T) {
	t.Parallel()

	// One trained profile: a bare `check` must auto-select it, removing the
	// per-run --author friction in the common single-author case.
	workDir := trainedProject(t)
	writeTestFile(t, filepath.Join(workDir, "draft.txt"),
		"今日はとても良い天気です。散歩に出かけました。気持ちが良かったです。")

	code, stdout, stderr := runApp(t, workDir, "check", "draft.txt")
	if code != 0 {
		t.Fatalf("bare check failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Author: me") || !strings.Contains(stdout, "Similarity:") {
		t.Fatalf("expected the single profile to be auto-selected: %q", stdout)
	}
}

func TestCheckDefaultAuthorWins(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me"
	jaCorpus(t, filepath.Join(workDir, "other"))
	if code, _, stderr := runApp(t, workDir, "train", "--author", "other", "--default", "other"); code != 0 {
		t.Fatalf("train --default failed: %s", stderr)
	}
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。散歩に行きました。")

	// Two profiles exist, but default_author=other was set, so a bare check uses it
	// rather than erroring on ambiguity.
	code, stdout, stderr := runApp(t, workDir, "check", "draft.txt")
	if code != 0 {
		t.Fatalf("default-author check failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Author: other") {
		t.Fatalf("expected the default author to be used: %q", stdout)
	}

	// An explicit --author still overrides the default.
	code, stdout, _ = runApp(t, workDir, "check", "--author", "me", "draft.txt")
	if code != 0 || !strings.Contains(stdout, "Author: me") {
		t.Fatalf("explicit --author should override the default: %q", stdout)
	}
}

func TestCheckAmbiguousProfilesError(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me"
	jaCorpus(t, filepath.Join(workDir, "other"))
	if code, _, stderr := runApp(t, workDir, "train", "--author", "other", "other"); code != 0 {
		t.Fatalf("second train failed: %s", stderr)
	}
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。")

	code, _, stderr := runApp(t, workDir, "check", "draft.txt")
	if code == 0 {
		t.Fatal("expected ambiguous profiles to fail without --author")
	}
	if !strings.Contains(stderr, "multiple profiles") {
		t.Fatalf("expected an ambiguity error listing the authors, got %q", stderr)
	}
	if !strings.Contains(stderr, "me") || !strings.Contains(stderr, "other") {
		t.Fatalf("the ambiguity error should name the candidates, got %q", stderr)
	}
}

func TestCheckScoreOnly(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	writeTestFile(t, filepath.Join(workDir, "draft.txt"),
		"今日はとても良い天気です。散歩に出かけました。気持ちが良かったです。")

	code, stdout, stderr := runApp(t, workDir, "check", "--score-only", "draft.txt")
	if code != 0 {
		t.Fatalf("--score-only failed: stderr=%q", stderr)
	}
	trimmed := strings.TrimSpace(stdout)
	if trimmed == "" {
		t.Fatal("expected a numeric score on stdout")
	}
	// Output must be exactly the integer: no "Author:", no "Similarity:" label.
	if strings.ContainsAny(trimmed, "AuthorSimilarity:%") {
		t.Fatalf("--score-only must print only the integer, got %q", stdout)
	}
	for _, r := range trimmed {
		if r < '0' || r > '9' {
			t.Fatalf("--score-only output is not a bare integer: %q", stdout)
		}
	}
}

func TestCheckScoreOnlyRejectsStructuredFlags(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。")

	if code, _, stderr := runApp(t, workDir, "check", "--score-only", "--explain", "draft.txt"); code == 0 ||
		!strings.Contains(stderr, "--score-only cannot be combined") {
		t.Fatalf("expected --score-only + --explain to be rejected, code=%d stderr=%q", code, stderr)
	}
	if code, _, stderr := runApp(t, workDir, "check", "--score-only", "--format", "json", "draft.txt"); code == 0 ||
		!strings.Contains(stderr, "--score-only cannot be combined") {
		t.Fatalf("expected --score-only + --format json to be rejected, code=%d stderr=%q", code, stderr)
	}
}

func TestRemoveAuthor(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	if code, stdout, stderr := runApp(t, workDir, "remove", "--author", "me"); code != 0 ||
		!strings.Contains(stdout, `Removed author "me"`) {
		t.Fatalf("remove failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	// The profile is gone from the listing, and removing it again fails cleanly.
	if _, stdout, _ := runApp(t, workDir, "list"); strings.Contains(stdout, "me") {
		t.Fatalf("removed author still listed: %q", stdout)
	}
	if code, _, stderr := runApp(t, workDir, "remove", "--author", "me"); code == 0 ||
		!strings.Contains(stderr, "profile not found") {
		t.Fatalf("expected removing a missing author to fail, code=%d stderr=%q", code, stderr)
	}
}

func TestRemoveClearsDefault(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me"
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "--default", "posts"); code != 0 {
		t.Fatalf("setting default failed: %s", stderr)
	}
	code, stdout, stderr := runApp(t, workDir, "remove", "--author", "me")
	if code != 0 {
		t.Fatalf("remove failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Cleared the default author") {
		t.Fatalf("expected the default to be cleared when its author is removed: %q", stdout)
	}
}

func TestRenameAuthor(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me"
	code, stdout, stderr := runApp(t, workDir, "rename", "--author", "me", "--to", "watashi")
	if code != 0 {
		t.Fatalf("rename failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, `Renamed author "me" to "watashi"`) {
		t.Fatalf("unexpected rename output: %q", stdout)
	}

	// The old name is gone, the new name resolves, and the stored author column was
	// rewritten (so `show` reports the new name).
	_, listOut, _ := runApp(t, workDir, "list")
	if strings.Contains(listOut, "me") && !strings.Contains(listOut, "watashi") {
		t.Fatalf("rename did not move the profile: %q", listOut)
	}
	code, showOut, _ := runApp(t, workDir, "show", "--author", "watashi")
	if code != 0 || !strings.Contains(showOut, "Author: watashi") {
		t.Fatalf("renamed profile not readable as the new name: %q", showOut)
	}
}

func TestRenameRefusesToClobber(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me"
	jaCorpus(t, filepath.Join(workDir, "other"))
	if code, _, stderr := runApp(t, workDir, "train", "--author", "other", "other"); code != 0 {
		t.Fatalf("train other failed: %s", stderr)
	}
	code, _, stderr := runApp(t, workDir, "rename", "--author", "me", "--to", "other")
	if code == 0 || !strings.Contains(stderr, "already exists") {
		t.Fatalf("expected rename onto an existing author to be refused, code=%d stderr=%q", code, stderr)
	}
}

func TestRenameUpdatesDefault(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "--default", "posts"); code != 0 {
		t.Fatalf("setting default failed: %s", stderr)
	}
	code, stdout, stderr := runApp(t, workDir, "rename", "--author", "me", "--to", "watashi")
	if code != 0 {
		t.Fatalf("rename failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, `Default author updated to "watashi"`) {
		t.Fatalf("expected the default to follow the rename: %q", stdout)
	}
	// A bare check now resolves through the updated default.
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。")
	if code, out, _ := runApp(t, workDir, "check", "draft.txt"); code != 0 || !strings.Contains(out, "Author: watashi") {
		t.Fatalf("default did not follow rename: code=%d out=%q", code, out)
	}
}

func TestShowTextAndJSON(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t) // author "me", 3 files

	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me")
	if code != 0 {
		t.Fatalf("show failed: stderr=%q", stderr)
	}
	for _, want := range []string{"Author: me", "Trained:", "Files: 3", "Source:"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("show text missing %q:\n%s", want, stdout)
		}
	}

	code, stdout, stderr = runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show --format json failed: stderr=%q", stderr)
	}
	var payload profileSummaryJSON
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("show json invalid: %v\n%s", err, stdout)
	}
	if payload.Author != "me" || payload.FileCount != 3 {
		t.Fatalf("unexpected show json: %+v", payload)
	}
}

func TestListLong(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)

	// Short form stays a bare list of names.
	_, short, _ := runApp(t, workDir, "list")
	if strings.TrimSpace(short) != "me" {
		t.Fatalf("short list should be just names: %q", short)
	}

	code, long, stderr := runApp(t, workDir, "list", "--long")
	if code != 0 {
		t.Fatalf("list --long failed: stderr=%q", stderr)
	}
	for _, want := range []string{"AUTHOR", "TRAINED", "FILES", "SOURCE", "me"} {
		if !strings.Contains(long, want) {
			t.Fatalf("list --long missing %q:\n%s", want, long)
		}
	}
}

func TestGlobalInitAndCheck(t *testing.T) {
	t.Parallel()

	// No local project: work from a directory with no omokage.toml above it, and a
	// dedicated global home.
	workDir := t.TempDir()
	home := filepath.Join(t.TempDir(), "omokage-home")
	jaCorpus(t, filepath.Join(workDir, "posts"))
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。散歩に行きました。")

	if code, stdout, stderr := runAppHome(t, workDir, home, "init", "--global"); code != 0 ||
		!strings.Contains(stdout, "global store") {
		t.Fatalf("init --global failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if code, _, stderr := runAppHome(t, workDir, home, "train", "--global", "--author", "me", "posts"); code != 0 {
		t.Fatalf("global train failed: %s", stderr)
	}

	// With a global store present and no local project, a bare check falls back to
	// the global store and auto-selects the single profile.
	code, stdout, stderr := runAppHome(t, workDir, home, "check", "draft.txt")
	if code != 0 {
		t.Fatalf("global fallback check failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Author: me") {
		t.Fatalf("expected the global profile to be used: %q", stdout)
	}
}

func TestLocalWinsOverGlobal(t *testing.T) {
	t.Parallel()

	home := filepath.Join(t.TempDir(), "omokage-home")
	workDir := t.TempDir()
	jaCorpus(t, filepath.Join(workDir, "posts"))

	// Global store trained with "global_author".
	if code, _, stderr := runAppHome(t, workDir, home, "init", "--global"); code != 0 {
		t.Fatalf("init --global failed: %s", stderr)
	}
	if code, _, stderr := runAppHome(t, workDir, home, "train", "--global", "--author", "global_author", "posts"); code != 0 {
		t.Fatalf("global train failed: %s", stderr)
	}

	// Local project in workDir trained with "local_author".
	if code, _, stderr := runAppHome(t, workDir, home, "init"); code != 0 {
		t.Fatalf("local init failed: %s", stderr)
	}
	if code, _, stderr := runAppHome(t, workDir, home, "train", "--author", "local_author", "posts"); code != 0 {
		t.Fatalf("local train failed: %s", stderr)
	}

	// Inside the local project, list must show only the local author: local wins.
	code, stdout, stderr := runAppHome(t, workDir, home, "list")
	if code != 0 {
		t.Fatalf("list failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "local_author") || strings.Contains(stdout, "global_author") {
		t.Fatalf("local project should win over global: %q", stdout)
	}

	// Forcing --global reaches the global store from the same directory.
	code, stdout, _ = runAppHome(t, workDir, home, "list", "--global")
	if code != 0 || !strings.Contains(stdout, "global_author") {
		t.Fatalf("--global should reach the global store: %q", stdout)
	}
}

func TestProfileDirOverride(t *testing.T) {
	t.Parallel()

	// --profile-dir points at an arbitrary directory with no config or project.
	workDir := t.TempDir()
	profileDir := filepath.Join(t.TempDir(), "myprofiles")
	if err := os.MkdirAll(profileDir, 0o750); err != nil {
		t.Fatal(err)
	}
	jaCorpus(t, filepath.Join(workDir, "posts"))
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。散歩に行きました。")

	if code, _, stderr := runApp(t, workDir, "train", "--profile-dir", profileDir, "--author", "me", "posts"); code != 0 {
		t.Fatalf("train --profile-dir failed: %s", stderr)
	}
	if _, err := os.Stat(filepath.Join(profileDir, "me.db")); err != nil {
		t.Fatalf("profile not written to --profile-dir: %v", err)
	}
	code, stdout, stderr := runApp(t, workDir, "check", "--profile-dir", profileDir, "draft.txt")
	if code != 0 || !strings.Contains(stdout, "Author: me") {
		t.Fatalf("check --profile-dir failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}
