package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeInputFile creates a file with content, making parent directories as needed.
func writeInputFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestGatherTrainingInputsMixedDirAndFiles(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	writeInputFile(t, filepath.Join(work, "posts", "a.md"), "alpha")
	writeInputFile(t, filepath.Join(work, "posts", "b.txt"), "beta")
	writeInputFile(t, filepath.Join(work, "draft.md"), "gamma")

	sources, files, err := gatherTrainingInputs(work, []string{"posts", "draft.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("expected 2 sources, got %d: %v", len(sources), sources)
	}
	if len(files) != 3 {
		t.Fatalf("expected 3 collected files, got %d: %v", len(files), files)
	}
}

func TestGatherTrainingInputsDeduplicatesDirAndFileOverlap(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	writeInputFile(t, filepath.Join(work, "posts", "a.md"), "alpha")
	writeInputFile(t, filepath.Join(work, "posts", "b.md"), "beta")

	// posts/ already contains a.md; passing both must not learn a.md twice.
	_, files, err := gatherTrainingInputs(work, []string{"posts", "posts/a.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 unique files after dedup, got %d: %v", len(files), files)
	}
}

func TestGatherTrainingInputsDeduplicatesRepeatedInput(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	writeInputFile(t, filepath.Join(work, "a.md"), "alpha")

	sources, files, err := gatherTrainingInputs(work, []string{"a.md", "a.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("expected the repeated input to collapse to 1 source, got %d: %v", len(sources), sources)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 unique file, got %d: %v", len(files), files)
	}
}

func TestGatherTrainingInputsRejectsUnsupportedFile(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	writeInputFile(t, filepath.Join(work, "notes.pdf"), "binary-ish")

	_, _, err := gatherTrainingInputs(work, []string{"notes.pdf"})
	if err == nil {
		t.Fatal("expected an error for an unsupported file extension")
	}
	if !strings.Contains(err.Error(), "unsupported file") {
		t.Fatalf("expected an unsupported-file error, got %v", err)
	}
}

func TestGatherTrainingInputsRejectsMissingPath(t *testing.T) {
	t.Parallel()

	work := t.TempDir()

	_, _, err := gatherTrainingInputs(work, []string{"nope.md"})
	if err == nil {
		t.Fatal("expected an error for a missing path")
	}
	if !strings.Contains(err.Error(), "input not found") {
		t.Fatalf("expected an input-not-found error, got %v", err)
	}
}

func TestGatherTrainingInputsDeduplicatesSymlink(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	target := filepath.Join(work, "a.md")
	writeInputFile(t, target, "alpha")
	alias := filepath.Join(work, "alias.md")
	if err := os.Symlink(target, alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	// a.md and alias.md -> a.md are the same real file; passing both must learn it
	// once, keyed by the resolved real path, so the distribution is not skewed.
	sources, files, err := gatherTrainingInputs(work, []string{"a.md", "alias.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected the symlink alias to collapse to 1 file, got %d: %v", len(files), files)
	}
	if len(sources) != 1 {
		t.Fatalf("expected 1 source after real-path dedup, got %d: %v", len(sources), sources)
	}
}

func TestLegacySourceDir(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	dir := filepath.Join(work, "posts")
	writeInputFile(t, filepath.Join(dir, "a.md"), "alpha")
	file := filepath.Join(work, "b.md")
	writeInputFile(t, file, "beta")

	// A single directory input keeps source_dir meaningful (a directory).
	if got := legacySourceDir([]string{dir}); got != dir {
		t.Fatalf("single directory: got %q, want %q", got, dir)
	}
	// A single file input must not put a file path into source_dir.
	if got := legacySourceDir([]string{file}); got != "" {
		t.Fatalf("single file: expected empty source_dir, got %q", got)
	}
	// Several inputs leave source_dir empty; provenance lives in Sources.
	if got := legacySourceDir([]string{dir, file}); got != "" {
		t.Fatalf("multiple inputs: expected empty source_dir, got %q", got)
	}
}

func TestGatherTrainingInputsRejectsURL(t *testing.T) {
	t.Parallel()

	work := t.TempDir()
	writeInputFile(t, filepath.Join(work, "a.md"), "alpha")

	for _, url := range []string{
		"https://example.com/post",
		"http://example.com/post.md",
		"https://user:pass@example.com/private",
		"ftp://example.com/file.txt",
	} {
		// A URL mixed in with a valid local file must still stop the whole run.
		_, _, err := gatherTrainingInputs(work, []string{"a.md", url})
		if err == nil {
			t.Fatalf("expected an error for URL input %q", url)
		}
		if !strings.Contains(err.Error(), "URL inputs are not supported") {
			t.Fatalf("expected a URL-not-supported error for %q, got %v", url, err)
		}
	}
}

func TestLooksLikeURL(t *testing.T) {
	t.Parallel()

	urls := []string{"https://x.com", "http://x", "ftp://x", "file://x", "s3://bucket/key"}
	for _, u := range urls {
		if !looksLikeURL(u) {
			t.Errorf("expected %q to be detected as a URL", u)
		}
	}
	// Local paths and Windows-style paths must not be mistaken for URLs.
	paths := []string{"posts", "./draft.md", "/abs/posts", "../up/a.txt", "a:b.md", "C:\\Users\\me\\a.md"}
	for _, p := range paths {
		if looksLikeURL(p) {
			t.Errorf("expected %q not to be detected as a URL", p)
		}
	}
}

// Training from a directory plus an extra file end to end: the file count covers
// both inputs and show surfaces the full provenance in text and JSON.
func TestTrainMultipleInputsShowProvenance(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	writeTestFile(t, filepath.Join(workDir, "posts", "one.txt"), "今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。")
	writeTestFile(t, filepath.Join(workDir, "posts", "two.txt"), "昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。")
	writeTestFile(t, filepath.Join(workDir, "draft.md"), "週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。")

	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	code, stdout, stderr := runApp(t, workDir, "train", "--author", "me", "posts", "draft.md")
	if code != 0 {
		t.Fatalf("train failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "from 3 files") {
		t.Fatalf("expected 3 files trained (2 in posts + draft.md), got: %q", stdout)
	}

	// show text: multiple inputs switch to the numbered "Sources (N):" block.
	code, stdout, stderr = runApp(t, workDir, "show", "--author", "me")
	if code != 0 {
		t.Fatalf("show failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Sources (2):") {
		t.Fatalf("expected a Sources block listing both inputs, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "draft.md") {
		t.Fatalf("expected the file input in the provenance, got:\n%s", stdout)
	}

	// show JSON: the sources array carries both inputs.
	code, stdout, stderr = runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show json failed: stderr=%q", stderr)
	}
	var payload profileSummaryJSON
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("show json invalid: %v\n%s", err, stdout)
	}
	if len(payload.Sources) != 2 {
		t.Fatalf("expected 2 sources in json, got %d: %v", len(payload.Sources), payload.Sources)
	}
}

// Training from a single file must not leak a file path into the source_dir field
// (reserved for an actual directory); the file shows up only in sources and on the
// human-readable "Source:" line.
func TestTrainSingleFileSourceDirStaysEmpty(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	writeTestFile(t, filepath.Join(workDir, "draft.md"),
		"今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。")

	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "draft.md"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}

	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show json failed: stderr=%q", stderr)
	}
	var payload profileSummaryJSON
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("show json invalid: %v\n%s", err, stdout)
	}
	if payload.SourceDir != "" {
		t.Fatalf("expected empty source_dir for a single file, got %q", payload.SourceDir)
	}
	if len(payload.Sources) != 1 || !strings.HasSuffix(payload.Sources[0], "draft.md") {
		t.Fatalf("expected the file in sources, got %v", payload.Sources)
	}

	// The human-readable view still names the file under "Source:".
	_, text, _ := runApp(t, workDir, "show", "--author", "me")
	if !strings.Contains(text, "Source: ") || !strings.Contains(text, "draft.md") {
		t.Fatalf("expected the file on the Source line, got:\n%s", text)
	}
}

func TestHelpRejectsExtraArguments(t *testing.T) {
	t.Parallel()

	// `help check extra` must fail exactly as `check extra --help` would, not drop
	// the trailing token and succeed.
	code, _, stderr := runApp(t, t.TempDir(), "help", "check", "extra")
	if code == 0 {
		t.Fatalf("expected a non-zero exit for trailing tokens after 'help check', got 0")
	}
	if stderr == "" {
		t.Fatal("expected an error on stderr")
	}

	// A clean `help check` still mirrors `check --help` and succeeds.
	if code, _, _ := runApp(t, t.TempDir(), "help", "check"); code != 0 {
		t.Fatalf("expected 'help check' to succeed, got exit %d", code)
	}
}
