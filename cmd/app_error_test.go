package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveVersionPrefersLdflags(t *testing.T) { //nolint:paralleltest // mutates package-level Version
	// Not parallel: it mutates the package-level Version that ldflags sets.
	original := Version
	t.Cleanup(func() { Version = original })

	Version = "v1.2.3"
	if got := resolveVersion(); got != "v1.2.3" {
		t.Fatalf("expected the ldflags version to win, got %q", got)
	}
}

func TestResolveVersionFallsBackWhenUnset(t *testing.T) { //nolint:paralleltest // mutates package-level Version
	// Not parallel: it mutates the package-level Version that ldflags sets.
	original := Version
	t.Cleanup(func() { Version = original })

	// With the default sentinel and no module version recorded in the test
	// binary's build info, resolveVersion must still return a non-empty label.
	Version = "dev"
	if got := resolveVersion(); got == "" {
		t.Fatal("expected a non-empty version label")
	}
}

func TestRunNoArgsShowsHelp(t *testing.T) {
	t.Parallel()

	code, stdout, _ := runApp(t, t.TempDir())
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "omokage compares writing style") {
		t.Fatalf("expected root help, got %q", stdout)
	}
}

func TestRunHelpAndVersion(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"help", "-h", "--help"} {
		if code, stdout, _ := runApp(t, t.TempDir(), arg); code != 0 || !strings.Contains(stdout, "Usage:") {
			t.Fatalf("help via %q failed: code=%d stdout=%q", arg, code, stdout)
		}
	}
	for _, arg := range []string{"version", "-v", "--version"} {
		if code, stdout, _ := runApp(t, t.TempDir(), arg); code != 0 || !strings.Contains(stdout, "omokage") {
			t.Fatalf("version via %q failed: code=%d stdout=%q", arg, code, stdout)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	t.Parallel()

	code, _, stderr := runApp(t, t.TempDir(), "frobnicate")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Fatalf("expected unknown command message, got %q", stderr)
	}
}

func TestRunInitRejectsExtraArgs(t *testing.T) {
	t.Parallel()

	if code, _, _ := runApp(t, t.TempDir(), "init", "extra"); code != 1 {
		t.Fatalf("expected exit 1 for extra init args, got %d", code)
	}
}

func TestRunTrainValidation(t *testing.T) {
	t.Parallel()

	// Missing --author.
	if code, _, _ := runApp(t, t.TempDir(), "train", "posts"); code != 1 {
		t.Fatal("expected failure when --author is missing")
	}
	// Missing directory argument.
	if code, _, _ := runApp(t, t.TempDir(), "train", "--author", "nao"); code != 1 {
		t.Fatal("expected failure when directory argument is missing")
	}
}

func TestRunTrainWithoutProject(t *testing.T) {
	t.Parallel()

	code, _, stderr := runApp(t, t.TempDir(), "train", "--author", "nao", "posts")
	if code != 1 {
		t.Fatalf("expected exit 1 without a project, got %d", code)
	}
	if !strings.Contains(stderr, "project not found") {
		t.Fatalf("expected project-not-found error, got %q", stderr)
	}
}

func TestRunTrainEmptyDirectory(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	if code, _, _ := runApp(t, workDir, "init"); code != 0 {
		t.Fatal("init failed")
	}
	writeTestFile(t, filepath.Join(workDir, "empty", "notes.csv"), "unsupported")

	code, _, stderr := runApp(t, workDir, "train", "--author", "nao", "empty")
	if code != 1 {
		t.Fatalf("expected exit 1 for a directory with no supported files, got %d", code)
	}
	if !strings.Contains(stderr, "no supported files") {
		t.Fatalf("expected no-supported-files error, got %q", stderr)
	}
}

func TestRunCheckValidation(t *testing.T) {
	t.Parallel()

	// Missing --author.
	if code, _, _ := runApp(t, t.TempDir(), "check", "target.md"); code != 1 {
		t.Fatal("expected failure when --author is missing")
	}
	// No project.
	if code, _, _ := runApp(t, t.TempDir(), "check", "--author", "nao", "target.md"); code != 1 {
		t.Fatal("expected failure without a project")
	}
}

func TestRunCheckMissingProfile(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	if code, _, _ := runApp(t, workDir, "init"); code != 0 {
		t.Fatal("init failed")
	}
	writeTestFile(t, filepath.Join(workDir, "target.md"), "some text\n")

	code, _, stderr := runApp(t, workDir, "check", "--author", "ghost", "target.md")
	if code != 1 {
		t.Fatalf("expected exit 1 for a missing profile, got %d", code)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Fatal("expected an error message for a missing profile")
	}
}

func TestRunDiffValidation(t *testing.T) {
	t.Parallel()

	// Wrong number of arguments.
	if code, _, _ := runApp(t, t.TempDir(), "diff", "only-one.md"); code != 1 {
		t.Fatal("expected failure with a single diff argument")
	}
	// Missing files.
	workDir := t.TempDir()
	if code, _, _ := runApp(t, workDir, "diff", "a.md", "b.md"); code != 1 {
		t.Fatal("expected failure when diff files do not exist")
	}
}

func TestRunListRejectsExtraArgs(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	if code, _, _ := runApp(t, workDir, "init"); code != 0 {
		t.Fatal("init failed")
	}
	if code, _, _ := runApp(t, workDir, "list", "extra"); code != 1 {
		t.Fatalf("expected exit 1 for extra list args, got %d", code)
	}
}

func TestRunListWithoutProject(t *testing.T) {
	t.Parallel()

	if code, _, _ := runApp(t, t.TempDir(), "list"); code != 1 {
		t.Fatal("expected failure listing without a project")
	}
}

func TestRunHelpSubcommandMatchesDashHelp(t *testing.T) {
	t.Parallel()

	// `omokage help <command>` must be identical to `omokage <command> --help` in
	// both content and exit code.
	for _, name := range []string{"check", "init", "train", "diff", "list", "show", "remove", "rename"} {
		helpCode, helpOut, helpErr := runApp(t, t.TempDir(), "help", name)
		dashCode, dashOut, dashErr := runApp(t, t.TempDir(), name, "--help")
		if helpCode != 0 || dashCode != 0 {
			t.Fatalf("%s: expected exit 0, got help=%d dash=%d", name, helpCode, dashCode)
		}
		if helpOut != dashOut || helpErr != dashErr {
			t.Fatalf("%s: `help %s` differs from `%s --help`\nhelp stdout=%q stderr=%q\ndash stdout=%q stderr=%q",
				name, name, name, helpOut, helpErr, dashOut, dashErr)
		}
	}
}

func TestRunHelpUnknownCommandFails(t *testing.T) {
	t.Parallel()

	code, _, stderr := runApp(t, t.TempDir(), "help", "frobnicate")
	if code != 1 {
		t.Fatalf("expected exit 1 for `help frobnicate`, got %d", code)
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Fatalf("expected unknown command message, got %q", stderr)
	}
}

func TestRunMissingArgMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"check missing FILE", []string{"check"}, "missing FILE"},
		{"diff missing FILE_B", []string{"diff", "a.md"}, "missing FILE_B"},
		{"train missing INPUT", []string{"train", "--author", "me"}, "missing INPUT"},
		{"train missing --author", []string{"train", "examples/posts"}, "missing --author"},
		{"remove missing --author", []string{"remove"}, "missing --author"},
		{"rename missing --to", []string{"rename", "--author", "me"}, "missing --to"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			code, _, stderr := runApp(t, t.TempDir(), tc.args...)
			if code != 1 {
				t.Fatalf("expected exit 1, got %d", code)
			}
			if !strings.Contains(stderr, tc.want) {
				t.Fatalf("expected stderr to contain %q, got %q", tc.want, stderr)
			}
			// The usage text must still follow the direct error.
			if !strings.Contains(stderr, "Usage: omokage") {
				t.Fatalf("expected usage after the error, got %q", stderr)
			}
		})
	}
}

func TestRunDiffWithoutProjectUsesDefaults(t *testing.T) {
	t.Parallel()

	// diff must work without an initialized project, falling back to defaults.
	workDir := t.TempDir()
	writeTestFile(t, filepath.Join(workDir, "a.md"), "# A\n\nそして文章です。だから続きます。\n")
	writeTestFile(t, filepath.Join(workDir, "b.md"), "# B\n\nしかし今日は違う書き方をしています。\n")

	code, stdout, stderr := runApp(t, workDir, "diff", "a.md", "b.md")
	if code != 0 {
		t.Fatalf("diff without project failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Similarity:") {
		t.Fatalf("expected a similarity line, got %q", stdout)
	}
}
