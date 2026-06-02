package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/tabwriter"
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
	// home is the base directory of the global profile store. It is read from
	// OMOKAGE_HOME, falling back to the per-user config directory. An empty home
	// disables the global fallback, which keeps the tool from ever touching a
	// user-wide store when one was never configured.
	home string
}

func NewApp(stdout, stderr io.Writer, workDir string) *App {
	return &App{
		stdout:  stdout,
		stderr:  stderr,
		workDir: workDir,
		home:    resolveHome(),
	}
}

// resolveHome locates the global store directory: OMOKAGE_HOME if set, otherwise
// "<user config dir>/omokage" (e.g. ~/.config/omokage on Linux). It returns ""
// only when neither is available, in which case --global and the global fallback
// are simply unavailable rather than guessing a path.
func resolveHome() string {
	if h := strings.TrimSpace(os.Getenv("OMOKAGE_HOME")); h != "" {
		return h
	}
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "omokage")
	}
	return ""
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printRootHelp()
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		// `omokage help` is the root help; `omokage help <command>` is the same as
		// `omokage <command> --help`, so users can reach a command's usage either way.
		if args[0] == "help" && len(args) > 1 {
			return a.runHelp(args[1])
		}
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
	case "show":
		return a.runShow(args[1:])
	case "remove":
		return a.runRemove(args[1:])
	case "rename":
		return a.runRename(args[1:])
	case "version", "-v", "--version":
		writef(a.stdout, "omokage %s\n", resolveVersion())
		return 0
	default:
		writef(a.stderr, "unknown command: %s\n\n", args[0])
		a.printRootHelp()
		return 1
	}
}

// runHelp implements `omokage help <command>`. For a command that has a usage
// screen it dispatches to `<command> --help`, so the two spellings stay identical
// in content and exit code. An unknown name fails the same way the root dispatcher
// does, rather than silently falling back to the root help.
func (a *App) runHelp(name string) int {
	switch name {
	case "init", "train", "check", "diff", "list", "show", "remove", "rename":
		return a.Run([]string{name, "--help"})
	case "help", "version":
		// These have no flags and no per-command usage screen; the root help already
		// documents them, so point there with a success exit.
		a.printRootHelp()
		return 0
	default:
		writef(a.stderr, "unknown command: %s\n\n", name)
		a.printRootHelp()
		return 1
	}
}

// scopeFlags are the store-selection flags shared by every command that reads or
// writes profiles. They are registered per command (rather than parsed globally)
// so each subcommand's --help still lists them, while the definitions stay in one
// place.
type scopeFlags struct {
	global     *bool
	configPath *string
	profileDir *string
}

func registerScopeFlags(flagSet *flag.FlagSet) scopeFlags {
	return scopeFlags{
		global:     flagSet.Bool("global", false, "use the global profile store ($OMOKAGE_HOME, else your user config dir like ~/.config/omokage) instead of searching for a local project"),
		configPath: flagSet.String("config", "", "path to an omokage.toml to use, overriding project discovery"),
		profileDir: flagSet.String("profile-dir", "", "directory of author profiles to use, overriding the config"),
	}
}

func (a *App) resolveScope(sf scopeFlags) (project.Scope, error) {
	return project.Resolve(project.ResolveOptions{
		WorkDir:    a.workDir,
		Home:       a.home,
		Global:     *sf.global,
		ConfigPath: strings.TrimSpace(*sf.configPath),
		ProfileDir: strings.TrimSpace(*sf.profileDir),
	})
}

// writeScopeError prints a resolve error, expanding the bare "project not found"
// into an actionable hint that points at both the local and global entry points.
func (a *App) writeScopeError(err error) {
	if errors.Is(err, project.ErrProjectNotFound) {
		writeLine(a.stderr, "omokage project not found; run 'omokage init' here, or 'omokage init --global' for a per-user store.")
		return
	}
	writeLine(a.stderr, err)
}

// resolveAuthor decides which profile `check`/`show` act on when --author is
// omitted. The rules are intentionally unambiguous:
//
//  1. an explicit --author always wins;
//  2. else the scope's configured default_author;
//  3. else, if exactly one profile exists, that one (the single-author case);
//  4. else it is an error — zero profiles, or two-plus with no default, never
//     silently picks one.
func (a *App) resolveAuthor(scope project.Scope, explicit string) (string, error) {
	if name := strings.TrimSpace(explicit); name != "" {
		return name, nil
	}
	if name := strings.TrimSpace(scope.Config.Defaults.Author); name != "" {
		return name, nil
	}

	authors, err := scope.ListProfiles()
	if err != nil {
		return "", err
	}
	switch len(authors) {
	case 0:
		return "", errors.New("no author profiles found; train one with 'omokage train --author NAME DIRECTORY'")
	case 1:
		return authors[0], nil
	default:
		// A bare --profile-dir scope has no config file, so "set a default in the
		// config" is not an option there — only --author can disambiguate.
		if scope.ConfigPath == "" {
			return "", fmt.Errorf("multiple profiles found (%s); pass --author NAME (this --profile-dir store has no config file to record a default in)",
				strings.Join(authors, ", "))
		}
		return "", fmt.Errorf("multiple profiles found (%s); pass --author NAME, or set a default with 'omokage train --default' (saved as default_author in %s)",
			strings.Join(authors, ", "), scope.ConfigPath)
	}
}

func (a *App) runInit(args []string) int {
	flagSet := newFlagSet("init", a.stderr)
	global := flagSet.Bool("global", false, "create the per-user store ($OMOKAGE_HOME, else your user config dir like ~/.config/omokage) instead of the current directory")
	name := flagSet.String("name", "", "project name (defaults to the directory name)")
	flagSet.Usage = func() {
		writef(a.stderr, "Create an omokage store (omokage.toml, profiles/, cache/).\n")
		writef(a.stderr, "Usage: omokage init [--global] [--name NAME]\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	dir := a.workDir
	label := "project"
	if *global {
		if a.home == "" {
			writeLine(a.stderr, "cannot determine the global store location; set OMOKAGE_HOME")
			return 1
		}
		dir = a.home
		label = "global store"
		if err := os.MkdirAll(dir, 0o750); err != nil {
			writeLine(a.stderr, err)
			return 1
		}
	}

	cfg, err := project.Init(dir, *name)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	writef(a.stdout, "Initialized omokage %s.\n", label)
	writef(a.stdout, "Config: %s\n", filepath.Join(dir, project.ConfigFileName))
	writef(a.stdout, "Profiles: %s\n", filepath.Join(dir, cfg.Storage.ProfileDir))
	writef(a.stdout, "Cache: %s\n", filepath.Join(dir, cfg.Storage.CacheDir))
	return 0
}

func (a *App) runTrain(args []string) int {
	flagSet := newFlagSet("train", a.stderr)
	author := flagSet.String("author", "", "author profile name")
	makeDefault := flagSet.Bool("default", false, "set this author as the store's default (used by check/show when --author is omitted)")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Learn an author's style from every .md and .txt file in DIRECTORY.\n")
		writef(a.stderr, "Usage: omokage train --author AUTHOR [--default] DIRECTORY\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if strings.TrimSpace(*author) == "" {
		return a.usageError(flagSet, "missing --author")
	}
	switch flagSet.NArg() {
	case 1:
		// exactly one DIRECTORY, as required
	case 0:
		return a.usageError(flagSet, "missing DIRECTORY")
	default:
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}
	if *makeDefault && scope.ConfigPath == "" {
		writeLine(a.stderr, "cannot set --default: this store has no config file to record it in")
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

	profilePath, err := scope.ProfilePath(*author)
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

	if *makeDefault {
		// The profile is saved and valid at this point; setting it as the default is
		// a separate write that can fail on its own (e.g. a read-only config). A
		// freshly trained profile with no default is a consistent state, so on
		// failure we report the partial result honestly instead of pretending the
		// default was recorded.
		scope.Config.Defaults.Author = *author
		if err := scope.SaveConfig(); err != nil {
			writef(a.stderr, "warning: the profile was trained, but setting it as the default failed: %v\n", err)
			return 1
		}
		writef(a.stdout, "Default author set to %q.\n", *author)
	}
	return 0
}

func (a *App) runCheck(args []string) int {
	flagSet := newFlagSet("check", a.stderr)
	author := flagSet.String("author", "", "author profile name (optional: defaults to default_author or the only trained profile)")
	explain := flagSet.Bool("explain", false, "print a prioritized, numeric drift report instead of the top-3 summary")
	format := flagSet.String("format", formatText, "output format: text or json (json implies --explain)")
	scoreOnly := flagSet.Bool("score-only", false, "print only the integer similarity (0-100), for scripts")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Score how closely FILE matches AUTHOR's trained style, from 0 to 100.\n")
		writef(a.stderr, "Usage: omokage check [--author AUTHOR] [--explain] [--format text|json] [--score-only] FILE\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	switch flagSet.NArg() {
	case 1:
		// exactly one FILE, as required
	case 0:
		return a.usageError(flagSet, "missing FILE")
	default:
		flagSet.Usage()
		return 1
	}
	if *format != formatText && *format != formatJSON {
		writef(a.stderr, "unknown --format %q: want text or json\n", *format)
		flagSet.Usage()
		return 1
	}
	// --score-only is the scalar, scripting output; --explain/--format json are the
	// structured outputs. They answer different needs, so combining them is a
	// mistake worth catching rather than silently picking one.
	if *scoreOnly && (*explain || *format == formatJSON) {
		writef(a.stderr, "--score-only cannot be combined with --explain or --format json\n")
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}

	resolvedAuthor, err := a.resolveAuthor(scope, *author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	profilePath, err := scope.ProfilePath(resolvedAuthor)
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

	if *scoreOnly {
		targetMetrics, err := feature.ExtractFile(targetPath)
		if err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		comparison := profile.Score(record.Distribution, targetMetrics, scope.Config.Features)
		writef(a.stdout, "%d\n", comparison.Similarity)
		return 0
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
			comparison: profile.Score(record.Distribution, targetMetrics, scope.Config.Features),
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
	explanation := profile.Explain(record.Distribution, targetMetrics, segments, scope.Config.Features)
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
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Compare two files directly and report their stylistic similarity.\n")
		writef(a.stderr, "No init, training, or profile is needed: diff reads the two files and compares them.\n")
		writef(a.stderr, "Usage: omokage diff FILE_A FILE_B\n")
		writef(a.stderr, "\nThe flags below are optional and only select which feature weights to use:\n")
		writef(a.stderr, "diff uses the active store's config when one is found, and the built-in defaults otherwise.\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	switch flagSet.NArg() {
	case 2:
		// exactly two files, as required
	case 0:
		return a.usageError(flagSet, "missing FILE_A and FILE_B")
	case 1:
		return a.usageError(flagSet, "missing FILE_B")
	default:
		flagSet.Usage()
		return 1
	}

	// diff only needs the feature set, not a profile, so it works without any
	// store: an active scope supplies the features, and its absence falls back to
	// the defaults rather than erroring.
	cfg := config.Default(filepath.Base(a.workDir))
	if scope, err := a.resolveScope(scopeF); err == nil {
		cfg = scope.Config
	} else if !errors.Is(err, project.ErrProjectNotFound) {
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
	long := flagSet.Bool("long", false, "show trained_at, file count, and source directory for each author")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "List the author profiles trained in the active store.\n")
		writef(a.stderr, "Usage: omokage list [--long]\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}

	authors, err := scope.ListProfiles()
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	// The short form stays a bare list of names: one per line, no header, so it
	// pipes cleanly into other tools. --long is the human-facing, annotated view.
	if !*long {
		for _, author := range authors {
			writeLine(a.stdout, author)
		}
		return 0
	}

	if len(authors) == 0 {
		writeLine(a.stdout, "No author profiles trained yet.")
		return 0
	}

	defaultAuthor := strings.TrimSpace(scope.Config.Defaults.Author)
	tw := tabwriter.NewWriter(a.stdout, 0, 2, 2, ' ', 0)
	writeLine(tw, "AUTHOR\tTRAINED\tFILES\tSOURCE")
	for _, author := range authors {
		profilePath, err := scope.ProfilePath(author)
		if err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		record, err := storage.LoadProfile(profilePath)
		if err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		name := author
		if author == defaultAuthor {
			// Mark the default so --long doubles as "which profile does a bare check
			// use?" without a second command.
			name += " (default)"
		}
		writef(tw, "%s\t%s\t%d\t%s\n",
			name, record.TrainedAt.Format("2006-01-02 15:04 MST"), record.FileCount, record.SourceDir)
	}
	return flushTab(a.stderr, tw)
}

func (a *App) runShow(args []string) int {
	flagSet := newFlagSet("show", a.stderr)
	author := flagSet.String("author", "", "author profile to show (optional: defaults to default_author or the only trained profile)")
	format := flagSet.String("format", formatText, "output format: text or json")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Show how an author profile was trained: when, from how many files, and from where.\n")
		writef(a.stderr, "Usage: omokage show [--author AUTHOR] [--format text|json]\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}
	if *format != formatText && *format != formatJSON {
		writef(a.stderr, "unknown --format %q: want text or json\n", *format)
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}
	resolvedAuthor, err := a.resolveAuthor(scope, *author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	profilePath, err := scope.ProfilePath(resolvedAuthor)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	record, err := storage.LoadProfile(profilePath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	if *format == formatJSON {
		payload := profileSummaryJSON{
			Author:         record.Author,
			TrainedAt:      record.TrainedAt.Format(time.RFC3339),
			FileCount:      record.FileCount,
			SourceDir:      record.SourceDir,
			DocumentCount:  record.Distribution.DocumentCount,
			SentenceCount:  record.Distribution.SentenceCount,
			CharacterCount: record.Distribution.CharacterCount,
			Default:        record.Author == strings.TrimSpace(scope.Config.Defaults.Author),
		}
		encoder := json.NewEncoder(a.stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(payload); err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		return 0
	}

	writef(a.stdout, "Author: %s\n", record.Author)
	if record.Author == strings.TrimSpace(scope.Config.Defaults.Author) {
		writeLine(a.stdout, "Default: yes")
	}
	writef(a.stdout, "Trained: %s\n", record.TrainedAt.Format("2006-01-02 15:04:05 MST"))
	writef(a.stdout, "Files: %d\n", record.FileCount)
	writef(a.stdout, "Source: %s\n", record.SourceDir)
	writef(a.stdout, "Documents: %d\n", record.Distribution.DocumentCount)
	writef(a.stdout, "Sentences: %d\n", record.Distribution.SentenceCount)
	writef(a.stdout, "Characters: %d\n", record.Distribution.CharacterCount)
	return 0
}

func (a *App) runRemove(args []string) int {
	flagSet := newFlagSet("remove", a.stderr)
	author := flagSet.String("author", "", "author profile to remove")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Remove an author profile from the active store.\n")
		writef(a.stderr, "Usage: omokage remove --author AUTHOR\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if strings.TrimSpace(*author) == "" {
		return a.usageError(flagSet, "missing --author")
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}
	profilePath, err := scope.ProfilePath(*author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	if _, err := os.Stat(profilePath); err != nil {
		if os.IsNotExist(err) {
			writef(a.stderr, "profile not found for author %q\n", *author)
		} else {
			writeLine(a.stderr, err)
		}
		return 1
	}
	// Clearing a dangling default and deleting the profile must not half-apply.
	// Clear the default first (a read-only omokage.toml fails here, before anything
	// is destroyed). If the profile delete then fails, restore the default so the
	// store returns to its prior state rather than ending up with the profile still
	// present but no default — which could make a later bare check ambiguous.
	clearedDefault := false
	if scope.ConfigPath != "" && strings.TrimSpace(scope.Config.Defaults.Author) == *author {
		scope.Config.Defaults.Author = ""
		if err := scope.SaveConfig(); err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		clearedDefault = true
	}
	if err := os.Remove(profilePath); err != nil {
		if clearedDefault {
			scope.Config.Defaults.Author = *author
			_ = scope.SaveConfig()
		}
		writeLine(a.stderr, err)
		return 1
	}

	writef(a.stdout, "Removed author %q.\n", *author)
	if clearedDefault {
		writeLine(a.stdout, "Cleared the default author.")
	}
	return 0
}

func (a *App) runRename(args []string) int {
	flagSet := newFlagSet("rename", a.stderr)
	author := flagSet.String("author", "", "current author name")
	to := flagSet.String("to", "", "new author name")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Rename an author profile, keeping its trained data.\n")
		writef(a.stderr, "Usage: omokage rename --author OLD --to NEW\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if strings.TrimSpace(*author) == "" {
		return a.usageError(flagSet, "missing --author")
	}
	if strings.TrimSpace(*to) == "" {
		return a.usageError(flagSet, "missing --to")
	}
	if flagSet.NArg() != 0 {
		flagSet.Usage()
		return 1
	}

	scope, err := a.resolveScope(scopeF)
	if err != nil {
		a.writeScopeError(err)
		return 1
	}
	oldPath, err := scope.ProfilePath(*author)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	newPath, err := scope.ProfilePath(*to)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	if _, err := os.Stat(oldPath); err != nil {
		if os.IsNotExist(err) {
			writef(a.stderr, "profile not found for author %q\n", *author)
		} else {
			writeLine(a.stderr, err)
		}
		return 1
	}
	// Never overwrite an existing profile: a silent clobber would destroy trained
	// data the user did not ask to lose.
	if _, err := os.Stat(newPath); err == nil {
		writef(a.stderr, "an author named %q already exists\n", *to)
		return 1
	} else if !os.IsNotExist(err) {
		writeLine(a.stderr, err)
		return 1
	}

	// Re-save under the new name (rewriting the stored author column) and then drop
	// the old file, so the profile's data and its name never disagree.
	record, err := storage.LoadProfile(oldPath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	record.Author = *to
	if err := storage.SaveProfile(newPath, record); err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	// Order the writes so a failure never leaves a dangling default or a half-done
	// rename: the new profile exists now, but the old one is still in place. Update
	// the config next (the fragile write); if it fails, roll the new profile back
	// out so the store returns to exactly its prior state. Only once the default is
	// safely recorded do we drop the old profile and report success.
	updatedDefault := false
	if scope.ConfigPath != "" && strings.TrimSpace(scope.Config.Defaults.Author) == *author {
		scope.Config.Defaults.Author = *to
		if err := scope.SaveConfig(); err != nil {
			_ = os.Remove(newPath)
			scope.Config.Defaults.Author = *author
			writeLine(a.stderr, err)
			return 1
		}
		updatedDefault = true
	}
	if err := os.Remove(oldPath); err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	writef(a.stdout, "Renamed author %q to %q.\n", *author, *to)
	if updatedDefault {
		writef(a.stdout, "Default author updated to %q.\n", *to)
	}
	return 0
}

func (a *App) printRootHelp() {
	writeLine(a.stdout, "omokage analyzes writing style and compares it against learned author profiles.")
	writeLine(a.stdout, "It works on Japanese and English text and keeps each profile in a local SQLite database.")
	writeLine(a.stdout)
	writeLine(a.stdout, "Usage:")
	writeLine(a.stdout, "  omokage <command> [arguments]")
	writeLine(a.stdout)
	writeLine(a.stdout, "Commands:")
	writeLine(a.stdout, "  init     Create an omokage store here, or --global for a per-user one.")
	writeLine(a.stdout, "  train    Learn an author's style from a directory of .md and .txt files.")
	writeLine(a.stdout, "  check    Score how closely a file matches a trained author (--explain for details).")
	writeLine(a.stdout, "  diff     Compare two files directly, without a trained profile.")
	writeLine(a.stdout, "  list     List the author profiles in the store (--long for details).")
	writeLine(a.stdout, "  show     Show how an author profile was trained.")
	writeLine(a.stdout, "  rename   Rename an author profile.")
	writeLine(a.stdout, "  remove   Remove an author profile.")
	writeLine(a.stdout, "  version  Print the omokage version.")
	writeLine(a.stdout, "  help     Show this help, or 'omokage help <command>' for one command.")
	writeLine(a.stdout)
	writeLine(a.stdout, `omokage uses a local project (omokage.toml found by walking up from the current`)
	writeLine(a.stdout, `directory) when one exists, otherwise the global store at $OMOKAGE_HOME, or your`)
	writeLine(a.stdout, `user config directory (e.g. ~/.config/omokage). 'omokage init --global' prints`)
	writeLine(a.stdout, `the exact path it created.`)
	writeLine(a.stdout)
	writeLine(a.stdout, `check picks the author from --author, then default_author, then the only profile.`)
	writeLine(a.stdout)
	writeLine(a.stdout, `Run "omokage <command> --help" (or "omokage help <command>") to see a command's arguments.`)
}

// profileSummaryJSON is the machine-readable form of `show`.
type profileSummaryJSON struct {
	Author         string `json:"author"`
	TrainedAt      string `json:"trained_at"`
	FileCount      int    `json:"file_count"`
	SourceDir      string `json:"source_dir"`
	DocumentCount  int    `json:"document_count"`
	SentenceCount  int    `json:"sentence_count"`
	CharacterCount int    `json:"character_count"`
	Default        bool   `json:"default"`
}

func newFlagSet(name string, output io.Writer) *flag.FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(output)
	return flagSet
}

// parseArgs parses a subcommand's flags and reports whether the caller should
// proceed. An explicit -h/--help is a successful request for the usage text (the
// flag package has already printed it), so it maps to exit 0; any other parse
// error maps to exit 1. Callers do `if code, ok := parseArgs(...); !ok { return
// code }`.
func parseArgs(flagSet *flag.FlagSet, args []string) (code int, ok bool) {
	if err := flagSet.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0, false
		}
		return 1, false
	}
	return 0, true
}

// usageError reports a missing or invalid argument directly on stderr, then prints
// the command's usage and returns exit code 1. Centralizing it keeps the "what is
// missing" wording and the message-then-usage layout consistent across commands.
func (a *App) usageError(flagSet *flag.FlagSet, msg string) int {
	writef(a.stderr, "%s\n\n", msg)
	flagSet.Usage()
	return 1
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

// flushTab flushes a tabwriter, reporting any error on stderr and turning it into
// an exit code so a write failure is not swallowed.
func flushTab(stderr io.Writer, tw *tabwriter.Writer) int {
	if err := tw.Flush(); err != nil {
		writeLine(stderr, err)
		return 1
	}
	return 0
}

func writeLine(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(w, args...)
}

func writef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
