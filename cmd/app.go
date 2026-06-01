package cmd

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/nao1215/dyer/internal/config"
	"github.com/nao1215/dyer/internal/feature"
	"github.com/nao1215/dyer/internal/profile"
	"github.com/nao1215/dyer/internal/project"
	"github.com/nao1215/dyer/internal/storage"
)

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
		writeLine(a.stdout, "dyer dev")
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
		writef(a.stderr, "Usage: dyer init [--name NAME]\n")
		flagSet.PrintDefaults()
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

	writeLine(a.stdout, "Initialized dyer project.")
	writef(a.stdout, "Config: %s\n", filepath.Join(a.workDir, project.ConfigFileName))
	writef(a.stdout, "Profiles: %s\n", filepath.Join(a.workDir, cfg.Storage.ProfileDir))
	writef(a.stdout, "Cache: %s\n", filepath.Join(a.workDir, cfg.Storage.CacheDir))
	return 0
}

func (a *App) runTrain(args []string) int {
	flagSet := newFlagSet("train", a.stderr)
	author := flagSet.String("author", "", "author profile name")
	flagSet.Usage = func() {
		writef(a.stderr, "Usage: dyer train --author AUTHOR DIRECTORY\n")
		flagSet.PrintDefaults()
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
	flagSet.Usage = func() {
		writef(a.stderr, "Usage: dyer check --author AUTHOR FILE\n")
		flagSet.PrintDefaults()
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

	targetMetrics, err := feature.ExtractFile(targetPath)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}

	comparison := profile.Score(record.Distribution, targetMetrics, cfg.Features)
	renderComparison(a.stdout, renderOptions{
		author:     record.Author,
		comparison: comparison,
	})
	return 0
}

func (a *App) runDiff(args []string) int {
	flagSet := newFlagSet("diff", a.stderr)
	flagSet.Usage = func() {
		writef(a.stderr, "Usage: dyer diff FILE_A FILE_B\n")
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
		leftPath:    leftPath,
		rightPath:   rightPath,
		comparison:  profile.Compare(leftMetrics, rightMetrics, cfg.Features),
		showSources: true,
	})
	return 0
}

func (a *App) runList(args []string) int {
	flagSet := newFlagSet("list", a.stderr)
	flagSet.Usage = func() {
		writef(a.stderr, "Usage: dyer list\n")
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
	writeLine(a.stdout, "dyer analyzes writing style and compares it against learned author profiles.")
	writeLine(a.stdout)
	writeLine(a.stdout, "Usage:")
	writeLine(a.stdout, "  dyer init [--name NAME]")
	writeLine(a.stdout, "  dyer train --author AUTHOR DIRECTORY")
	writeLine(a.stdout, "  dyer check --author AUTHOR FILE")
	writeLine(a.stdout, "  dyer diff FILE_A FILE_B")
	writeLine(a.stdout, "  dyer list")
}

func newFlagSet(name string, output io.Writer) *flag.FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(output)
	return flagSet
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

func writeLine(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(w, args...)
}

func writef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
