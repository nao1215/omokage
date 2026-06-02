# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and per-release binaries and notes are published from git tags by GoReleaser.

## [Unreleased]

## [0.2.0] - 2026-06-02

### Added

- Term preferences: `train` now learns, per profile, which surface form you use
  for a recurring term (`DB` vs `データベース`, `HTTP` vs `http`) and stores it in
  the same per-author database — no LLM, no network, no dictionary, and the
  training text itself is never stored. `normalized_key` (case, full-width ASCII,
  and edge punctuation folded) and `group_key` are kept separate, so a
  normalization merge is distinguishable from a corpus-declared alias bridge such
  as `データベース（DB）`; bridges are conservative (a Japanese phrase paired with an
  uppercase acronym), so unrelated terms are never fused.
- `show --format json` adds a `term_preferences` array (group, preferred surface,
  counts, and variants). `check --format json` adds a `term_warnings` array for
  any draft surface that differs from its group's preferred form. Term warnings
  are a separate layer and never change the similarity score; both appear only in
  the JSON output, so plain `check` is unchanged.

### Fixed

- Feature extraction now measures natural-language prose only. YAML front matter,
  Markdown images and link URLs (visible text kept), raw URLs, HTML tags, and HTML
  entities are stripped alongside code, and per-paragraph extraction strips on the
  whole document before splitting, so a fenced block (a mermaid diagram or shell
  session) containing a blank line no longer leaks into the report. `check
  --explain`/`--format json` no longer points at diagrams, HTML, or front matter
  as the paragraphs that drift most, and inline markup no longer skews the script
  ratios.

## [0.1.0] - 2026-06-02

### Added

- Author resolution for `check` (and `show`): with a single trained profile
  `--author` is now optional, and a `default_author` setting (or `train
  --default`) picks the author when several exist. Two or more profiles with no
  default is a clear error that lists the candidates rather than guessing.
- Profile management commands so `profiles/*.db` never has to be edited by hand:
  `show` (how a profile was trained, with `--format json`), `rename` (keeps the
  trained data, refuses to overwrite), and `remove` (clears the default if it
  pointed there).
- A per-user global store alongside the existing project-local model:
  `omokage init --global`, the `OMOKAGE_HOME` environment variable, and the
  `--global`, `--config PATH`, and `--profile-dir PATH` flags. A local project
  always wins inside its tree; the global store is the fallback when none is
  found.
- `list --long` adds a table with each author's trained_at, file count, and
  source directory, marking the default; plain `list` still prints bare names.
- `check --score-only` prints just the integer similarity, for shell scripts.
- `check --explain`: a prioritized, numeric drift report that leads with the
  high-level, editable features (register, script balance, sentence and
  paragraph shape) and shows each one's target value, the trained mean and
  spread, the z-score, and a fix priority. The low-level function-word and
  character n-gram fingerprint follows as supporting detail.
- `check --format json`: the same analysis as a single JSON object, with the
  high-level and low-level drifts in separate arrays and the drifting
  paragraphs localized, for an LLM to read in a revise-and-recheck loop.
- Paragraph-level drift localization, surfaced by both `--explain` and
  `--format json`, pointing at the paragraphs that stray furthest from the
  trained style.
- Release tooling: GoReleaser configuration and a tag-triggered GitHub Actions
  release workflow that publishes cross-platform archives, checksums, and Linux
  packages (deb/rpm/apk).
- `SECURITY.md` describing how to report vulnerabilities.

### Changed

- `omokage version` now reports the release tag. GoReleaser injects it via
  ldflags, and `go install`/`go build` resolve it from the embedded module
  version, so installs from a tag report that tag instead of `dev`.

The opt-in detailed analysis runs only when requested, so the default `check`
output and its performance are unchanged.
