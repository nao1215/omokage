package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/profile"
	"github.com/nao1215/omokage/internal/project"
	"github.com/nao1215/omokage/internal/storage"
)

// Output formats for `check`. text is the default human-readable report; json
// is the machine-readable explanation for an LLM revise-and-recheck loop.
const (
	formatText = "text"
	formatJSON = "json"
)

// devVersion is the sentinel reported for an untagged local build.
const devVersion = "dev"

// Version is the release version. goreleaser overrides it at build time via
// ldflags (-X github.com/nao1215/omokage/cmd.Version=...). When it is left at
// the default, resolveVersion falls back to the module version embedded by the
// Go toolchain, so `go install ...@v1.2.3` (or @latest) still reports the tag.
var Version = devVersion

// resolveVersion returns the version to print. A goreleaser build sets Version
// via ldflags and wins outright. Otherwise the binary was built with `go install`
// or `go build`, so the module version recorded in the build info is used: that
// is the git tag for `go install path@tag`/@latest, or "(devel)" for an untagged
// local build, which we report as devVersion.
func resolveVersion() string {
	if Version != devVersion {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return devVersion
}

type App struct {
	stdout  io.Writer
	stderr  io.Writer
	workDir string
}

func NewApp(stdout, stderr io.Writer, workDir string) *App {
	return &App{
		stdout:  stdout,
		stderr:  stderr,
		workDir: workDir,
	}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printRootHelp()
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printRootHelp()
		return 0
	case "init":
		return a.runInit(args[1:])
	case "train":
		return a.runTrain(args[1:])
	case "check":
		return a.runCheck(args[1:])
	case "diff":
		return a.runDiff(args[1:])
	case "list":
		return a.runList(args[1:])
	case "version", "-v", "--version":
		writef(a.stdout, "omokage %s\n", resolveVersion())
		return 0
	default:
		writef(a.stderr, "unknown command: %s\n\n", args[0])
		a.printRootHelp()
		return 1
	}
}

func (a *App) runInit(args []string) int {
	flagSet := newFlagSet("init", a.stderr)
	name := flagSet.String("name", filepath.Base(a.workDir), "project name")
	flagSet.Usage = func() {
		writef(a.stderr, "Create an omokage project in the current directory (omokage.toml, profiles/, cache/).\n")
		writef(a.stderr, "Usage: omokage init [--name NAME]\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	cfg, err := project.Init(a.workDir, *name)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	writeLine(a.stdout, "Initialized omokage project.")
	writef(a.stdout, "Config: %s\n", filepath.Join(a.workDir, project.ConfigFileName))
	writef(a.stdout, "Profiles: %s\n", filepath.Join(a.workDir, cfg.Storage.ProfileDir))
	writef(a.stdout, "Cache: %s\n", filepath.Join(a.workDir, cfg.Storage.CacheDir))
	return 0
}

func (a *App) runTrain(args []string) int {
	flagSet := newFlagSet("train", a.stderr)
	author := flagSet.String("author", "", "author profile name")
	flagSet.Usage = func() {
		writef(a.stderr, "Learn an author's style from every .md and .txt file in DIRECTORY.\n")
		writef(a.stderr, "Usage: omokage train --author AUTHOR DIRECTORY\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if strings.TrimSpace(*author) == "" || flagSet.NArg() != 1 {
		flagSet.Usage()
		return 1
	}

	root, cfg, err := project.Load(a.workDir)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	sourceDir, err := resolvePath(a.workDir, flagSet.Arg(0))
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	files, err := feature.CollectFiles(sourceDir)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	if len(files) == 0 {
		writef(a.stderr, "no supported files found in %s\n", sourceDir)
		return 1
	}

	distribution, err := feature.ExtractCorpus(files)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	// CollectFiles found supported files, but ExtractCorpus drops empty or
	// whitespace-only documents. If nothing usable is left, a saved profile would
	// be all zeros and every later check would score against noise, so refuse.
	if distribution.DocumentCount == 0 {
		writef(a.stderr, "no usable text found in %s (all files were empty)\n", sourceDir)
		return 1
	}

	profilePath, err := project.ProfilePath(root, cfg, *author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	record := profile.Record{
		Author:       *author,
		SourceDir:    sourceDir,
		TrainedAt:    time.Now().UTC(),
		FileCount:    len(files),
		Distribution: distribution,
	}
	if err := storage.SaveProfile(profilePath, record); err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	writef(a.stdout, "Trained author %q from %d files.\n", record.Author, record.FileCount)
	writef(a.stdout, "Profile: %s\n", profilePath)
	return 0
}

func (a *App) runCheck(args []string) int {
	flagSet := newFlagSet("check", a.stderr)
	author := flagSet.String("author", "", "author profile name")
	explain := flagSet.Bool("explain", false, "print a prioritized, numeric drift report instead of the top-3 summary")
	format := flagSet.String("format", formatText, "output format: text or json (json implies --explain)")
	flagSet.Usage = func() {
		writef(a.stderr, "Score how closely FILE matches AUTHOR's trained style, from 0 to 100.\n")
		writef(a.stderr, "Usage: omokage check --author AUTHOR [--explain] [--format text|json] FILE\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if strings.TrimSpace(*author) == "" || flagSet.NArg() != 1 {
		flagSet.Usage()
		return 1
	}
	if *format != formatText && *format != formatJSON {
		writef(a.stderr, "unknown --format %q: want text or json\n", *format)
		flagSet.Usage()
		return 1
	}

	root, cfg, err := project.Load(a.workDir)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	profilePath, err := project.ProfilePath(root, cfg, *author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	record, err := storage.LoadProfile(profilePath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	targetPath, err := resolvePath(a.workDir, flagSet.Arg(0))
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	// The plain path extracts whole-document metrics only. The explain/json path
	// additionally splits the document into paragraphs so it can localize drift;
	// that extra work runs only when the detailed output was requested.
	detailed := *explain || *format == formatJSON
	if !detailed {
		targetMetrics, err := feature.ExtractFile(targetPath)
		if err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		renderComparison(a.stdout, renderOptions{
			author:     record.Author,
			comparison: profile.Score(record.Distribution, targetMetrics, cfg.Features),
		})
		// A one-line pointer to the detailed report, but only when a person is
		// watching: it goes to stderr and only when stderr is a terminal. A pipe, a
		// redirect, a `$(...)` capture, a script, or an LLM harness gets clean output
		// on both streams; an interactive user gets the hint. The flags stay
		// discoverable for everyone through `check --help` and the root help.
		if isTerminal(a.stderr) {
			writeLine(a.stderr, "Tip: add --explain (or --format json) for per-feature drift, fix priority, and the paragraphs that drifted most.")
		}
		return 0
	}

	targetMetrics, segments, err := feature.ExtractFileWithSegments(targetPath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	explanation := profile.Explain(record.Distribution, targetMetrics, segments, cfg.Features)
	if *format == formatJSON {
		if err := renderExplanationJSON(a.stdout, record.Author, explanation); err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		return 0
	}
	renderExplanationText(a.stdout, record.Author, explanation)
	return 0
}

func (a *App) runDiff(args []string) int {
	flagSet := newFlagSet("diff", a.stderr)
	flagSet.Usage = func() {
		writef(a.stderr, "Compare two files directly and report their stylistic similarity, no profile needed.\n")
		writef(a.stderr, "Usage: omokage diff FILE_A FILE_B\n")
	}
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if flagSet.NArg() != 2 {
		flagSet.Usage()
		return 1
	}

	cfg, err := a.defaultOrProjectConfig()
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	leftPath, err := resolvePath(a.workDir, flagSet.Arg(0))
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	rightPath, err := resolvePath(a.workDir, flagSet.Arg(1))
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	leftMetrics, err := feature.ExtractFile(leftPath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	rightMetrics, err := feature.ExtractFile(rightPath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	renderComparison(a.stdout, renderOptions{
		leftPath:    flagSet.Arg(0),
		rightPath:   flagSet.Arg(1),
		comparison:  profile.Compare(leftMetrics, rightMetrics, cfg.Features),
		showSources: true,
	})
	return 0
}

func (a *App) runList(args []string) int {
	flagSet := newFlagSet("list", a.stderr)
	flagSet.Usage = func() {
		writef(a.stderr, "List the author profiles trained in this project.\n")
		writef(a.stderr, "Usage: omokage list\n")
	}
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	root, cfg, err := project.Load(a.workDir)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	authors, err := project.ListProfiles(root, cfg)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	for _, author := range authors {
		writeLine(a.stdout, author)
	}
	return 0
}

func (a *App) defaultOrProjectConfig() (config.Config, error) {
	cfg, found, err := project.LoadOptional(a.workDir)
	if err != nil {
		return config.Config{}, err
	}
	if found {
		return cfg, nil
	}
	return config.Default(filepath.Base(a.workDir)), nil
}

func (a *App) printRootHelp() {
	writeLine(a.stdout, "omokage analyzes writing style and compares it against learned author profiles.")
	writeLine(a.stdout, "It works on Japanese and English text and keeps each profile in a local SQLite database.")
	writeLine(a.stdout)
	writeLine(a.stdout, "Usage:")
	writeLine(a.stdout, "  omokage <command> [arguments]")
	writeLine(a.stdout)
	writeLine(a.stdout, "Commands:")
	writeLine(a.stdout, "  init     Create an omokage project here (omokage.toml, profiles/, cache/).")
	writeLine(a.stdout, "  train    Learn an author's style from a directory of .md and .txt files.")
	writeLine(a.stdout, "  check    Score how closely a file matches a trained author (--explain for details).")
	writeLine(a.stdout, "  diff     Compare two files directly, without a trained profile.")
	writeLine(a.stdout, "  list     List the author profiles trained in this project.")
	writeLine(a.stdout, "  version  Print the omokage version.")
	writeLine(a.stdout)
	writeLine(a.stdout, `Run "omokage <command> --help" to see a command's arguments.`)
}

func newFlagSet(name string, output io.Writer) *flag.FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(output)
	return flagSet
}

// printFlagDefaults lists a command's flags using the double-dash spelling shown
// in each Usage line. Go's flag package accepts both -flag and --flag but prints
// the single-dash form, which reads as inconsistent next to the "--author" usage
// strings; this keeps the help uniform.
func printFlagDefaults(w io.Writer, flagSet *flag.FlagSet) {
	flagSet.VisitAll(func(f *flag.Flag) {
		typeName, usage := flag.UnquoteUsage(f)
		if typeName != "" {
			writef(w, "  --%s %s\n", f.Name, typeName)
		} else {
			writef(w, "  --%s\n", f.Name)
		}
		writef(w, "        %s", usage)
		// Mirror flag.PrintDefaults: show a default only when it is meaningful
		// (a non-empty, non-false value), so boolean and empty-string flags stay
		// uncluttered.
		if f.DefValue != "" && f.DefValue != "false" {
			writef(w, " (default %q)", f.DefValue)
		}
		writeLine(w)
	})
}

func resolvePath(baseDir, target string) (string, error) {
	if filepath.IsAbs(target) {
		return filepath.Clean(target), nil
	}
	return filepath.Abs(filepath.Join(baseDir, target))
}

type renderOptions struct {
	author      string
	leftPath    string
	rightPath   string
	comparison  profile.Comparison
	showSources bool
}

func renderComparison(w io.Writer, opt renderOptions) {
	if opt.author != "" {
		writef(w, "Author: %s\n", opt.author)
	}
	if opt.showSources {
		writef(w, "Reference: %s\n", opt.leftPath)
		writef(w, "Target: %s\n", opt.rightPath)
	}
	writef(w, "Similarity: %d%%\n", opt.comparison.Similarity)
	writeLine(w)
	writeLine(w, "Differences:")
	for _, difference := range opt.comparison.Differences {
		writef(w, "- %s\n", difference)
	}
}

// isTerminal reports whether the writer is an interactive terminal. It is used to
// show the discoverability tip only to a human at a console, never into a pipe,
// redirect, or capture. A non-*os.File writer (e.g. the bytes.Buffer the tests
// inject) is never a terminal, so the tip stays out of programmatic output.
func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func writeLine(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(w, args...)
}

func writef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
