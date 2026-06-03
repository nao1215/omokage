[![Build](https://github.com/nao1215/omokage/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/build.yml)
[![MultiPlatformUnitTest](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml)
[![Coverage](https://github.com/nao1215/omokage/actions/workflows/coverage.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/omokage.svg)](https://pkg.go.dev/github.com/nao1215/omokage)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/omokage)](https://goreportcard.com/report/github.com/nao1215/omokage)
![GitHub](https://img.shields.io/github/license/nao1215/omokage)

<p align="center">
  <img src="doc/img/omokage-icon.jpg" alt="omokage" width="320">
</p>

omokage learns how you write from your past writing, then scores how close a new draft is to that style. It runs locally, works on Japanese and English, and never uses the network.

![demo](./doc/img/demo.gif)

## What it does (and doesn't)

- **Does**: compare *style* — sentence shape, register (敬体 / 常体), kanji/kana balance, word and character patterns — between a draft and a trained author, and point out where they differ.
- **Doesn't**: judge meaning, correctness, originality, or quality. It is not an AI-text detector. A high score means "this reads like the voice you trained," nothing more.

It is built for an LLM as much as for a person: an agent can run `check` after each rewrite, read the differences, and revise until the draft sits closer to the voice.

## Install

```shell
go install github.com/nao1215/omokage@latest
```

Runs on Windows, macOS, and Linux. Building from source needs Go 1.25 or later.

## Quick start

The repo ships a small example corpus under [examples/](./examples) to follow along.

```shell
$ omokage init                                   # writes omokage.toml, profiles/, cache/
$ omokage train --author me examples/en/posts    # learn a voice from .md/.txt files
Trained author "me" from 8 files.
Profile: /home/me/blog/profiles/me.db

Corpus reliability: good.
$ omokage check examples/en/draft-keeps-voice.md  # score a draft (one profile needs no --author)
Author: me
Similarity: 70%

Differences:
- character n-gram "gh" is higher than reference
- function word "at" is higher than reference
- character n-gram "ht" is higher than reference
```

`train` takes any mix of directories (scanned for `.md`/`.txt`) and individual files; a file reached twice is learned once. It reads local files only — a URL, a missing path, or an unsupported extension stops the run by name and trains nothing.

The same idea rewritten in a stiff, formal voice scores low:

```shell
$ omokage check --author me examples/en/draft-lost-voice.md
Author: me
Similarity: 26%

Differences:
- average sentence length is higher than reference
- paragraph length variance is higher than reference
- sentence length variance is higher than reference
```

`omokage diff A B` compares two files directly, without training a profile.

## Checking a corpus

Scores are only as steady as the corpus behind them. A good corpus is several documents (aim for eight or more), each a few paragraphs long, all in one consistent voice. `doctor` rates a corpus — training and writing nothing — and names what to fix:

```shell
$ omokage doctor ~/writing/posts
Corpus: 8 documents, 142 sentences, 5210 characters (avg 651 per document)
Reliability: good

No problems found: enough material, a consistent voice, and no obvious outliers.

These checks look at sample size and consistency, not writing quality.
```

```shell
$ omokage doctor ~/drafts
Corpus: 3 documents, 9 sentences, 140 characters (avg 46 per document)
Reliability: weak

Findings:
- [warning] Only 3 documents. The measured spread is barely an estimate, so scores will be noisy.
    → Add more samples of this voice; 8 or more documents give steadier scores.
- [warning] 3 of 3 documents are short (under 150 characters).
    a.md, b.md, c.md
    → Short samples make per-document features jumpy; prefer samples of a few paragraphs.

These checks look at sample size and consistency, not writing quality.
```

![doctor demo](./doc/img/doctor.gif)

`doctor --format json` gives the same report as data. `train` prints the reliability too (with the findings and a pointer to `doctor` when a corpus is thin or mixed), and `show --format json` carries the rating and findings so you can read a profile's standing later. A mixed corpus is flagged by the feature it disagrees on (often the register or kanji/kana balance) — the fix is to split it into one profile per voice.

## Choosing the author

`--author` is just a profile name; it need not be a person. Name a profile for a purpose — `--author blog`, `--author docs` — and train each on the writing that belongs to it. `check` and `show` resolve the author as: `--author` if given, else `default_author`, else the only profile, else an error (they never silently pick one). Set a default with `train --author me --default ...`.

## Output modes

`check` reads one file; pick how you want the result:

| Mode | Output | For |
| --- | --- | --- |
| (default) | similarity score + top differences | quick, human-facing checks |
| `--score-only` | the integer 0-100 | shell pipelines, pass/fail gates |
| `--explain` | per-feature drift (value, mean ± spread, z-score) + the paragraphs that drift most | final by-hand tuning |
| `--format json` | the `--explain` detail as JSON, plus `term_warnings` | an LLM or tool reading between rewrites |

`--explain` and `--format json` split the draft into paragraphs, so they are opt-in and plain `check` stays fast. `--score-only` can't be combined with them.

```shell
$ score=$(omokage check --score-only draft.md)
$ [ "$score" -ge 70 ] && echo "close enough"

$ omokage check --author me --explain examples/ja/draft-lost-voice.md
Author: me
Similarity: 0%

High-level style differences (fix these first):
  1. polite sentence-ending ratio is lower than reference [register]
       target 0.000  reference 1.000 ± 0.000  (50.0σ)
  ...

Paragraphs that drift most:
  #2 (50.0σ; polite sentence-ending ratio lower): 雨天時は在宅で過ごすケースが多い。特段の活動は行わない…
```

![explain demo](./doc/img/explain.gif)

## Using omokage with an LLM

Train once, then on each rewrite have the agent run `check --format json` and read it back. The JSON leads with `high_level_drift` — the editable features, each with a `priority` and `actionable` flag — so the agent knows what to change first; `segments` points at the paragraphs that drift most, and `term_warnings` flags notation that differs from your learned preference (never part of the score). For a lighter payload to hand an agent, `show --author me --format json --summary` returns provenance and the quality rating without the (often large) term list. omokage tells the agent how close the draft sits to your voice and where it strays — not whether it is correct or good, so keep a human in the loop.

## Term preferences

`train` also learns which surface form you use for a recurring term (`DB` vs `データベース`, `HTTP` vs `http`), stored in the same per-author database — no LLM, no network, no dictionary, and only surfaces and counts are kept, never the text. A `normalized_key` folds case and full/half-width ASCII so `DB`, `db`, and `ＤＢ` share a key; a `group_key` merges a Japanese phrase with its acronym only when the corpus declares the bridge (`データベース（DB）`). `show --format json` lists them under `term_preferences`, and `check --format json` adds `term_warnings`; both appear only in JSON, so plain `check` is unchanged.

## Managing profiles and stores

```shell
$ omokage list [--long]                # names, or trained_at / file count / source(s)
$ omokage show --author me             # how a profile was trained (--format json for more)
$ omokage rename --author me --to watashi
$ omokage remove --author watashi
```

By default omokage finds an `omokage.toml` by walking up from the current directory (a project-local store). `omokage init --global` makes a per-user store under `$OMOKAGE_HOME` (or your user config dir) that any directory falls back to; a local project always wins inside its tree. `--config PATH` / `--profile-dir PATH` point at a specific store.

## How it scores

Training measures a set of stylistic features per document and stores their mean and spread in a SQLite database under `profiles/` (one per author) — the numbers only, never the text. The features are register (敬体 / 常体), script balance (kanji/hiragana/katakana), function words, character n-grams, and shape (sentence and paragraph length, punctuation, layout). A check measures the same features on the draft and scores each by how far it strays from your usual range, as a z-score in the spirit of Burrows's Delta: the function-word and n-gram fingerprint carries most of the signal, a clear register shift is penalized on top, and shape only nudges. Code blocks are stripped first, so the score reflects prose.

## Limits

omokage looks at style, not meaning: it cannot tell whether a draft is correct, original, or well written, only whether it resembles the voice it was trained on. It needs a reasonable amount of text — with a few short documents the spread is wide and scores are noisy, which `doctor` and the `reliability` rating flag (they measure sample adequacy, not writing quality). It separates Japanese authors more sharply than English ones, and two people who write in the same register look more alike than they are. It is not an AI-text detector.

## About the name

omokage (面影) is written with 面 (face) and 影 (shadow, trace): the remembered image of someone, the likeness that comes back to mind. The name is borrowed from [Omokage](https://www.toraya-group.co.jp/products/collections/yokan-omokage), a yokan by Toraya that I like.

## License

MIT. See [LICENSE](./LICENSE).
