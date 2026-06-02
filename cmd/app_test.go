package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppLifecycle(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	corpusDir := filepath.Join(workDir, "posts")
	if err := os.MkdirAll(corpusDir, 0o750); err != nil {
		t.Fatal(err)
	}

	writeTestFile(t, filepath.Join(corpusDir, "one.md"), "# Title\n\nI write short notes. However, I still use markdown.\n- bullet\n- bullet\n")
	writeTestFile(t, filepath.Join(corpusDir, "two.txt"), "そして今日は静かです。だから文章は短めです。")
	writeTestFile(t, filepath.Join(workDir, "target.md"), "# Draft\n\nI write short notes. But this draft uses different pacing.\n")

	code, stdout, stderr := runApp(t, workDir, "init", "--name", "sample-style")
	if code != 0 {
		t.Fatalf("init failed: stdout=%q stderr=%q", stdout, stderr)
	}
	if !strings.Contains(stdout, "Initialized omokage project.") {
		t.Fatalf("unexpected init stdout: %q", stdout)
	}

	code, stdout, stderr = runApp(t, workDir, "train", "--author", "nao", "posts")
	if code != 0 {
		t.Fatalf("train failed: stdout=%q stderr=%q", stdout, stderr)
	}
	if !strings.Contains(stdout, `Trained author "nao"`) {
		t.Fatalf("unexpected train stdout: %q", stdout)
	}

	code, stdout, stderr = runApp(t, workDir, "list")
	if code != 0 {
		t.Fatalf("list failed: stdout=%q stderr=%q", stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "nao" {
		t.Fatalf("unexpected list stdout: %q", stdout)
	}

	code, stdout, stderr = runApp(t, workDir, "check", "--author", "nao", "target.md")
	if code != 0 {
		t.Fatalf("check failed: stdout=%q stderr=%q", stdout, stderr)
	}
	if !strings.Contains(stdout, "Author: nao") || !strings.Contains(stdout, "Similarity:") {
		t.Fatalf("unexpected check stdout: %q", stdout)
	}

	code, stdout, stderr = runApp(t, workDir, "diff", "posts/one.md", "target.md")
	if code != 0 {
		t.Fatalf("diff failed: stdout=%q stderr=%q", stdout, stderr)
	}
	if !strings.Contains(stdout, "Reference:") || !strings.Contains(stdout, "Target:") {
		t.Fatalf("unexpected diff stdout: %q", stdout)
	}
}

func TestInitRejectsExistingProject(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("first init failed: %s", stderr)
	}

	code, _, stderr := runApp(t, workDir, "init")
	if code == 0 {
		t.Fatal("expected second init to fail")
	}
	if !strings.Contains(stderr, "already exists") {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

func runApp(t *testing.T, workDir string, args ...string) (int, string, string) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app := NewApp(&stdout, &stderr, workDir)
	code := app.Run(args)
	return code, stdout.String(), stderr.String()
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
