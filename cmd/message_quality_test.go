package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArityRefusalsSayWhatWasWrong is the first half of the regression for #19.
// The zero-argument paths of check and diff say `missing FILE`, but the
// too-many-arguments paths fell through to a bare usage block, so the two
// mistakes that are easiest to make from a shell glob — `omokage check *.md` —
// answered with the one message that does not say what happened.
func TestArityRefusalsSayWhatWasWrong(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"check with two files", []string{"check", "one.md", "two.md"}, "check takes exactly one FILE"},
		{"check with three files", []string{"check", "a.md", "b.md", "c.md"}, "check takes exactly one FILE"},
		{"diff with three files", []string{"diff", "a.md", "b.md", "c.md"}, "diff takes exactly two files"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, _, stderr := runApp(t, t.TempDir(), tt.args...)
			if code != 1 {
				t.Fatalf("%v: exit %d, want 1", tt.args, code)
			}
			if !strings.Contains(stderr, tt.want) {
				t.Errorf("%v: stderr does not say %q:\n%s", tt.args, tt.want, stderr)
			}
		})
	}
}

// TestCheckOnADirectoryPointsAtTrain is the second half. `omokage check ./dir/`
// leaked the raw `read ./dir/: is a directory` from os.ReadFile, which names a
// syscall rather than the thing to do instead — while train, which is the
// command that does take a directory, has an IsDir branch of its own.
func TestCheckOnADirectoryPointsAtTrain(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	corpus := filepath.Join(workDir, "corpus")
	if err := os.MkdirAll(corpus, 0o750); err != nil {
		t.Fatal(err)
	}

	code, _, stderr := runApp(t, workDir, "check", "corpus")
	if code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
	if !strings.Contains(stderr, "is a directory") || !strings.Contains(stderr, "train") {
		t.Errorf("stderr does not name the directory and point at train:\n%s", stderr)
	}
}

// TestDiffOfEmptyDocumentsSaysThereIsNothingToMeasure is the third. Two empty
// files compared to 100% similar with the difference reported as "no enabled
// features configured" — a sentence about configuration, for a run where every
// feature was enabled and simply had nothing to measure. The two cases have
// different fixes, so they have to read differently.
func TestDiffOfEmptyDocumentsSaysThereIsNothingToMeasure(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	writeTestFile(t, filepath.Join(workDir, "empty.md"), "")

	code, stdout, stderr := runApp(t, workDir, "diff", "empty.md", "empty.md")
	if code != 0 {
		t.Fatalf("exit %d, want 0 (an empty document is not an error): %s", code, stderr)
	}
	if strings.Contains(stdout, "no enabled features configured") {
		t.Errorf("empty documents still blame the feature configuration:\n%s", stdout)
	}
	if !strings.Contains(stdout, "no measurable text") {
		t.Errorf("stdout does not say the documents have nothing to measure:\n%s", stdout)
	}
}
