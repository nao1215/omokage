// Command runner bootstraps omokage's end-to-end suite and hands the specs to
// atago.
//
// It builds omokage from this checkout into a throwaway directory, puts that
// directory first on PATH so the specs exercise that exact binary, points HOME
// and every store location at a sandbox under the same temp tree, and runs the
// atago specs under e2e/atago. Nothing the suite does can reach the developer's
// real configuration directory or profile store, and a run leaves nothing
// behind.
//
// The test DEFINITIONS are the atago YAML; this program is only the environment
// bootstrap. It is Go rather than the shell script it replaces because the suite
// runs on Windows too, and a POSIX bootstrap would make the Windows leg depend
// on Git Bash being installed -- which tests the runner image rather than
// omokage. Process launching, temp trees and PATH handling are things the
// standard library already does portably.
//
// Usage:
//
//	go run ./e2e/runner                # every spec under e2e/atago
//	go run ./e2e/runner --filter check # extra flags are passed through to atago
//	COVER=1 GOCOVERDIR=... go run ./e2e/runner
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "e2e: %v\n", err)
		var exit *exitError
		if errors.As(err, &exit) {
			os.Exit(exit.code)
		}
		os.Exit(1)
	}
}

// exitError carries an exit status out of run, so a failing atago run exits with
// atago's own status instead of a generic 1. CI reads that status to tell "the
// suite failed" from "the bootstrap broke".
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string { return e.err.Error() }
func (e *exitError) Unwrap() error { return e.err }

func run(ctx context.Context, args []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	if _, err := exec.LookPath("atago"); err != nil {
		return &exitError{code: 127, err: fmt.Errorf(
			"atago is not installed. Install it from https://github.com/nao1215/atago\n"+
				"e2e: e.g. 'go install github.com/nao1215/atago@latest' (CI uses nao1215/setup-atago): %w", err)}
	}

	sandbox, err := os.MkdirTemp("", "omokage-e2e-")
	if err != nil {
		return fmt.Errorf("can not create the e2e temp tree: %w", err)
	}
	defer func() { _ = os.RemoveAll(sandbox) }()

	binDir := filepath.Join(sandbox, "bin")
	binary := filepath.Join(binDir, "omokage")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if err := buildOmokage(ctx, repoRoot, binary); err != nil {
		return err
	}

	if err := sandboxEnv(repoRoot, sandbox, binDir); err != nil {
		return err
	}

	// Extra args (e.g. --filter check) come before the path so atago's flag
	// parser sees them as flags rather than targets.
	atagoArgs := append([]string{"run", "--ci"}, args...)
	if !hasTarget(args) {
		atagoArgs = append(atagoArgs, filepath.Join(repoRoot, "e2e", "atago"))
	}

	atago := exec.CommandContext(ctx, "atago", atagoArgs...) //nolint:gosec // a fixed command with author-supplied spec targets
	atago.Dir = repoRoot
	atago.Stdout = os.Stdout
	atago.Stderr = os.Stderr
	atago.Stdin = os.Stdin
	if err := atago.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return &exitError{code: exit.ExitCode(), err: fmt.Errorf("atago run failed: %w", err)}
		}
		return fmt.Errorf("can not run atago: %w", err)
	}
	return nil
}

// sandboxEnv points every location omokage may read or write at the throwaway
// tree, and puts the freshly built binary first on PATH.
//
// The two store variables are separate on purpose. OMOKAGE_HOME stays empty for
// the whole run, so a scenario with no local project deterministically reports
// "project not found" rather than depending on a store some earlier scenario
// happened to create; OMOKAGE_TEST_GLOBAL is where the one `init --global`
// scenario is pointed, keeping its write out of both the tracked workdir and the
// shared store.
func sandboxEnv(repoRoot, sandbox, binDir string) error {
	home := filepath.Join(sandbox, "home")
	for _, dir := range []string{
		home,
		filepath.Join(sandbox, "config"),
		filepath.Join(sandbox, "data"),
		filepath.Join(sandbox, "cache"),
		filepath.Join(sandbox, "ohome"),
		filepath.Join(sandbox, "gtest"),
	} {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("can not create the e2e sandbox directory %s: %w", dir, err)
		}
	}

	env := map[string]string{
		// HOME and USERPROFILE cover POSIX and Windows home resolution.
		"HOME":            home,
		"USERPROFILE":     home,
		"XDG_CONFIG_HOME": filepath.Join(sandbox, "config"),
		"XDG_DATA_HOME":   filepath.Join(sandbox, "data"),
		"XDG_CACHE_HOME":  filepath.Join(sandbox, "cache"),
		"OMOKAGE_HOME":    filepath.Join(sandbox, "ohome"),

		"OMOKAGE_TEST_GLOBAL": filepath.Join(sandbox, "gtest"),
		// The example corpus in the repository is the fixture source for the
		// specs, so they carry no duplicate training data. Absolute, because
		// atago runs each scenario in its own isolated workdir.
		"OMOKAGE_EXAMPLES": filepath.Join(repoRoot, "examples"),

		"PATH": binDir + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	for name, value := range env {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("can not set %s: %w", name, err)
		}
	}
	return nil
}

// buildOmokage compiles the binary under test. COVER=1 (used by
// scripts/coverage.sh) builds a coverage-instrumented omokage so the E2E run's
// covdata can be merged with unit coverage; the omokage processes atago spawns
// inherit GOCOVERDIR and each writes raw covdata on exit. With COVER unset the
// build is a plain one, so `make test-e2e-atago` stays a test of the shipped
// binary.
func buildOmokage(ctx context.Context, repoRoot, output string) error {
	fmt.Println("e2e: building omokage...")

	buildArgs := []string{"build"}
	if os.Getenv("COVER") != "" {
		if os.Getenv("GOCOVERDIR") == "" {
			return errors.New("COVER=1 requires GOCOVERDIR to be set (see scripts/coverage.sh)")
		}
		buildArgs = append(buildArgs, "-cover", "-covermode=atomic", "-coverpkg=./...")
	}
	// Mirror the Makefile's VERSION so `omokage version` answers the same way
	// the `make build` binary does; "dev" when no tag is reachable, which is what
	// a shallow CI checkout gives.
	buildArgs = append(buildArgs,
		"-ldflags", "-X github.com/nao1215/omokage/cmd.Version="+describeVersion(ctx, repoRoot),
		"-o", output, "main.go")

	cmd := exec.CommandContext(ctx, "go", buildArgs...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GO111MODULE=on")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can not build omokage: %w", err)
	}
	return nil
}

// describeVersion returns the newest reachable tag, or "dev" when there is none
// -- the same answer the Makefile's `git describe --tags --abbrev=0` gives.
func describeVersion(ctx context.Context, repoRoot string) string {
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return "dev"
	}
	if version := strings.TrimSpace(string(out)); version != "" {
		return version
	}
	return "dev"
}

// findRepoRoot walks up from the working directory to the checkout root. go run
// compiles into a temp directory, so os.Executable is unreliable here; the
// module root is found by looking for go.mod above the current directory.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can not determine the working directory: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("can not find the repository root (no go.mod above the working directory)")
		}
		dir = parent
	}
}

// valueFlags are the `atago run` options that take their value as the NEXT
// argument. Without this list the runner cannot tell `--filter check` (a flag
// and its value) from `--filter` followed by a spec path, and would then skip
// adding the default target -- silently running nothing.
var valueFlags = map[string]bool{ //nolint:gochecknoglobals // a fixed table read by hasTarget
	"artifacts-dir": true,
	"filter":        true,
	"parallel":      true,
	"profile":       true,
	"repeat":        true,
	"report":        true,
	"retry-failed":  true,
	"skip-tag":      true,
	"tag":           true,
}

// hasTarget reports whether the caller named spec files or directories of their
// own, as opposed to passing only atago flags.
func hasTarget(args []string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			return true
		}
		name := strings.TrimLeft(arg, "-")
		if strings.Contains(name, "=") {
			continue
		}
		if valueFlags[name] {
			i++
		}
	}
	return false
}
