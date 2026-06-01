[![Build](https://github.com/nao1215/dyer/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/dyer/actions/workflows/build.yml)
[![MultiPlatformUnitTest](https://github.com/nao1215/dyer/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/dyer/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/dyer/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/dyer/actions/workflows/reviewdog.yml)
[![Coverage](https://github.com/nao1215/dyer/actions/workflows/coverage.yml/badge.svg)](https://github.com/nao1215/dyer/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/dyer.svg)](https://pkg.go.dev/github.com/nao1215/dyer)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/dyer)](https://goreportcard.com/report/github.com/nao1215/dyer)

# dyer

`dyer` is a local-first CLI for learning an author's writing style and measuring how closely a target text matches that profile.

The current implementation supports:

- project initialization with `dyer.toml`, `profiles/`, and `cache/`
- training author profiles from `.md` and `.txt` files
- profile storage in SQLite
- profile-to-document comparison with a similarity score
- direct document-to-document comparison

## Requirements

- Go 1.25 or later
- Linux, macOS, or Windows

## Install

```bash
go install github.com/nao1215/dyer@latest
```

## Quick Start

Initialize a project:

```bash
dyer init
```

Train a profile from a directory of Markdown and text files:

```bash
dyer train --author nao ./posts
```

Check a target document against the stored profile:

```bash
dyer check --author nao article.md
```

Compare two documents directly:

```bash
dyer diff before.md after.md
```

List registered profiles:

```bash
dyer list
```

## Project Layout

```text
dyer-project/
├── dyer.toml
├── profiles/
│   └── nao.db
└── cache/
```

`dyer` stores each author profile in a separate SQLite database under `profiles/`.

## Default Configuration

```toml
[project]
name = "my-writing-profiles"

[features]
sentence_length = true
sentence_length_variance = true
punctuation_frequency = true
newline_frequency = true
bullet_ratio = true
kanji_ratio = true
hiragana_ratio = true
katakana_ratio = true
conjunction_frequency = true
paragraph_length_variance = true
markdown_structure_frequency = true

[storage]
profile_dir = "./profiles"
cache_dir = "./cache"
```

## Implemented Style Features

- average sentence length
- sentence length variance
- punctuation frequency
- newline frequency
- bullet-list ratio
- conjunction frequency
- kanji / hiragana / katakana ratio
- paragraph length variance
- markdown structure frequency

## Development

```bash
make build
make test
make lint
```

Install local development tools:

```bash
make tools
```

More developer guidance is available in [CONTRIBUTING.md](./CONTRIBUTING.md) and [doc/development.md](./doc/development.md).
