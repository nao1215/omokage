# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and per-release binaries and notes are published from git tags by GoReleaser.

## [Unreleased]

### Added

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
