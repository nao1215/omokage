package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tapeDir holds the VHS scripts that render the GIFs README.md embeds.
const tapeDir = "../doc/img"

// TestDemoTapeCommandsStillWork runs every omokage command line the tapes type,
// in the order they type them, and requires each to exit 0.
//
// A GIF cannot be diffed to prove it is current -- VHS output varies with
// timing, fonts and the encoder -- so what CI can check is the thing that
// actually rots: the command lines inside the tape. Rename a subcommand, drop a
// flag, or move an example file and this fails here, instead of the next
// `make demo` quietly recording a terminal full of usage errors.
//
// Order matters and is preserved: `check` needs the profile `train` writes, and
// `train` needs the project `init` creates. Each tape gets its own working
// directory with a fresh copy of examples/, which is what `make demo` sets up
// under /tmp/omokage-demo before each run.
func TestDemoTapeCommandsStillWork(t *testing.T) {
	t.Parallel()

	tapes, err := filepath.Glob(filepath.Join(tapeDir, "*.tape"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tapes) == 0 {
		t.Fatalf("no tapes under %s", tapeDir)
	}

	examples, err := filepath.Abs("../examples")
	if err != nil {
		t.Fatal(err)
	}

	for _, tape := range tapes {
		t.Run(filepath.Base(tape), func(t *testing.T) {
			t.Parallel()

			commands := omokageCommandsInTape(t, tape)
			if len(commands) == 0 {
				t.Fatalf("%s types no omokage commands; the tape or this parser is wrong", tape)
			}

			workDir := t.TempDir()
			copyTree(t, examples, filepath.Join(workDir, "examples"))

			for _, command := range commands {
				args := strings.Fields(command)[1:]
				code, stdout, stderr := runApp(t, workDir, args...)
				if code != 0 {
					t.Fatalf("%q exited with %d\nstdout: %s\nstderr: %s", command, code, stdout, stderr)
				}
			}
		})
	}
}

// omokageCommandsInTape returns the omokage command lines a VHS tape types, in
// order. VHS keeps typed text in `Type "..."` directives; everything else in the
// file (Set, Sleep, Hide, Show) is not a command. A typed line may be a shell
// one-liner that chains several commands with `;` and throws output away with
// `>/dev/null`, so each fragment is looked at on its own and the redirect is
// dropped -- the test captures the output itself.
func omokageCommandsInTape(t *testing.T, path string) []string {
	t.Helper()

	data, err := os.ReadFile(path) //nolint:gosec // a path this test globbed inside the repository
	if err != nil {
		t.Fatal(err)
	}

	var commands []string
	for _, line := range strings.Split(string(data), "\n") {
		typed, ok := typedText(strings.TrimSpace(line))
		if !ok {
			continue
		}
		for _, fragment := range strings.Split(typed, ";") {
			fragment = strings.TrimSpace(fragment)
			if !strings.HasPrefix(fragment, "omokage ") {
				continue
			}
			commands = append(commands, strings.TrimSpace(strings.ReplaceAll(fragment, ">/dev/null", "")))
		}
	}
	return commands
}

// typedText extracts the quoted argument of a VHS `Type "..."` directive. A tape
// line may carry more directives after it, so the text ends at the closing quote
// rather than at the end of the line.
func typedText(line string) (string, bool) {
	const prefix = `Type "`
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := line[len(prefix):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return "", false
	}
	return rest[:end], true
}

// copyTree copies src into dst recursively. examples/ is a small tree of Markdown
// under en/ and ja/, so a plain walk is enough.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dst, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		from := filepath.Join(src, entry.Name())
		to := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			copyTree(t, from, to)
			continue
		}
		data, err := os.ReadFile(from) //nolint:gosec // fixture files inside the repository
		if err != nil {
			t.Fatal(err)
		}
		// Both ends are paths this test built: src is a fixed directory in the
		// repository and dst is under t.TempDir().
		if err := os.WriteFile(to, data, 0o600); err != nil { //nolint:gosec // see above
			t.Fatal(err)
		}
	}
}

// TestDemoTapesRenderTheGIFsTheReadmeEmbeds pins each tape's Output name. The
// README embeds doc/img/demo.gif, doctor.gif and explain.gif; `make demo` moves
// each tape's output there by that name, so a tape that renders something else
// leaves the README showing a stale image with nothing failing.
func TestDemoTapesRenderTheGIFsTheReadmeEmbeds(t *testing.T) {
	t.Parallel()

	for tape, want := range map[string]string{
		"demo.tape":    "Output demo.gif",
		"doctor.tape":  "Output doctor.gif",
		"explain.tape": "Output explain.gif",
	} {
		t.Run(tape, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(filepath.Join(tapeDir, tape)) //nolint:gosec // a fixed path inside the repository
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), want) {
				t.Errorf("%s does not contain %q, so `make demo` would not produce the GIF README.md embeds", tape, want)
			}
		})
	}
}
