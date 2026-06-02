package cmd

import (
	"bytes"
	"encoding/json"
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

// trainedProject sets up a project with one author trained on a polite-register
// Japanese corpus and returns the working directory. It backs the explain/json
// output tests.
func trainedProject(t *testing.T) string {
	t.Helper()

	workDir := t.TempDir()
	corpusDir := filepath.Join(workDir, "posts")
	writeTestFile(t, filepath.Join(corpusDir, "one.txt"), "今日は朝から雨が降っています。傘を持って出かけました。電車はとても混んでいました。")
	writeTestFile(t, filepath.Join(corpusDir, "two.txt"), "昨日は友人と食事に行きました。料理はどれも美味しかったです。また行きたいと思います。")
	writeTestFile(t, filepath.Join(corpusDir, "three.txt"), "週末は近くの公園を散歩しました。空気が澄んでいて気持ちが良かったです。")

	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "posts"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}
	return workDir
}

func TestCheckExplainTextOutput(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	// A register-flipped draft: the high-level register feature should lead.
	writeTestFile(t, filepath.Join(workDir, "draft.txt"),
		"本日は降雨である。外出を実施した。混雑は著しいものであった。対応を継続するものとする。")

	code, stdout, stderr := runApp(t, workDir, "check", "--author", "me", "--explain", "draft.txt")
	if code != 0 {
		t.Fatalf("check --explain failed: stderr=%q", stderr)
	}
	for _, want := range []string{"Author: me", "Similarity:", "High-level style", "σ)"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("explain output missing %q:\n%s", want, stdout)
		}
	}
}

func TestCheckJSONOutput(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	writeTestFile(t, filepath.Join(workDir, "draft.txt"),
		"本日は降雨である。外出を実施した。混雑は著しいものであった。対応を継続するものとする。")

	code, stdout, stderr := runApp(t, workDir, "check", "--author", "me", "--format", "json", "draft.txt")
	if code != 0 {
		t.Fatalf("check --format json failed: stderr=%q", stderr)
	}

	var payload struct {
		Author     string `json:"author"`
		Similarity int    `json:"similarity"`
		HighLevel  []struct {
			Feature       string  `json:"feature"`
			Category      string  `json:"category"`
			Target        float64 `json:"target"`
			ReferenceMean float64 `json:"reference_mean"`
			Priority      int     `json:"priority"`
			Actionable    bool    `json:"actionable"`
		} `json:"high_level_drift"`
		LowLevel []struct {
			Feature string `json:"feature"`
		} `json:"low_level_drift"`
		Segments []struct {
			Index      int    `json:"index"`
			TopFeature string `json:"top_feature"`
		} `json:"segments"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("check --format json did not emit valid JSON: %v\n%s", err, stdout)
	}
	if payload.Author != "me" {
		t.Fatalf("unexpected author: %q", payload.Author)
	}
	if payload.Similarity < 0 || payload.Similarity > 100 {
		t.Fatalf("similarity out of range: %d", payload.Similarity)
	}
	if len(payload.HighLevel) == 0 {
		t.Fatal("expected high-level drift entries in JSON output")
	}
	if payload.HighLevel[0].Priority != 1 {
		t.Fatalf("expected the first high-level drift to have priority 1, got %d", payload.HighLevel[0].Priority)
	}
}

func TestCheckRejectsUnknownFormat(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	writeTestFile(t, filepath.Join(workDir, "draft.txt"), "本日は晴天なり。")

	code, _, stderr := runApp(t, workDir, "check", "--author", "me", "--format", "yaml", "draft.txt")
	if code == 0 {
		t.Fatal("expected an unknown format to fail")
	}
	if !strings.Contains(stderr, "unknown --format") {
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
