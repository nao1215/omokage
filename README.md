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

omokage learns how you write from your past writing, then tells you how close a new draft is to that style. It runs locally and works on Japanese and English text.

![demo](./doc/img/demo.gif)

## Why I built it

I often draft text with an LLM and then rework it so that it reads like something I wrote. Prompting and hand-editing only get me so far, so I wanted a tool that measures how close a draft is to my own style and points out where it drifts. omokage is that tool. You train it on your past writing, then check a draft against it.

It is meant to be used by an LLM as much as by a person. An agent can run `check` after each rewrite, read the similarity and the differences, and keep revising until the draft sits closer to the trained voice.

## Install

```shell
go install github.com/nao1215/omokage@latest
```

It runs on Windows, macOS, and Linux. Building from source needs Go 1.25 or later.

## Usage

The repository includes a small example corpus under [examples/](./examples) so you can follow along — English under [examples/en/](./examples/en) and Japanese under [examples/ja/](./examples/ja).

Create a project in the current directory. This writes `omokage.toml`, `profiles/`, and `cache/`.

```shell
$ omokage init
Initialized omokage project.
Config: /home/me/blog/omokage.toml
Profiles: /home/me/blog/profiles
Cache: /home/me/blog/cache
```

Learn a style from past writing. `train` takes one or more inputs: directories
(scanned recursively for `.md` and `.txt`) and individual `.md`/`.txt` files, in
any mix, so you need not gather everything into one folder first.

```shell
$ omokage train --author me examples/en/posts
Trained author "me" from 8 files.
Profile: /home/me/blog/profiles/me.db
```

Paths may be relative or absolute. `examples/en/posts` holds 8 files, so adding
one more makes 9; a file reached twice — through a containing directory or a
symlink — is learned once, matched by its real path.

```shell
$ omokage train --author me examples/en/posts examples/en/draft-keeps-voice.md
Trained author "me" from 9 files.
Profile: /home/me/blog/profiles/me.db
```

omokage reads local files only and never touches the network. An unsupported
extension, a missing path, or a URL stops the run by name and trains nothing, so
you can drop the bad argument and retry.

```shell
$ omokage train --author me https://example.com/post
URL inputs are not supported: https://example.com/post (omokage trains from local files only; save the page as a .md or .txt file and pass that path instead)

$ omokage train --author me posts notes.pdf
unsupported file notes.pdf: omokage learns only .md and .txt files

$ omokage train --author me posts missing.md
input not found: missing.md
```

Check whether a draft still reads like that author. With a single trained
profile you can drop `--author` entirely — omokage selects the only one.

```shell
$ omokage check examples/en/draft-keeps-voice.md
Author: me
Similarity: 70%

Differences:
- character n-gram "gh" is higher than reference
- function word "at" is higher than reference
- character n-gram "ht" is higher than reference
```

The same idea rewritten in a stiff, formal voice scores low, and omokage shows what changed.

```shell
$ omokage check --author me examples/en/draft-lost-voice.md
Author: me
Similarity: 26%

Differences:
- average sentence length is higher than reference
- paragraph length variance is higher than reference
- sentence length variance is higher than reference
```

omokage works the same on Japanese, where the sentence-ending register (敬体 / 常体) and the kanji/kana balance make the difference even clearer. The `diff` and `--explain` examples below use the Japanese corpus to show that depth.

You can also compare two documents directly, without training a profile.

```shell
$ omokage diff examples/ja/draft-keeps-voice.md examples/ja/draft-lost-voice.md
Reference: examples/ja/draft-keeps-voice.md
Target: examples/ja/draft-lost-voice.md
Similarity: 54%

Differences:
- polite sentence-ending ratio is lower than reference
- paragraph length variance is lower than reference
- sentence length variance is higher than reference
```

Similarity runs from 0 to 100 and shows how close the text sits to the learned style. Differences lists the features that moved the most, such as the sentence-ending register (敬体 / 常体), the balance of kanji and kana, function-word and character n-gram usage, sentence length, and layout. omokage compares style rather than topic, so writing about something new in your usual voice still scores high.

For final tuning, add `--explain` to lead with the high-level, editable features (register, script balance, structure), each with the draft's value, your trained mean ± spread, a z-score, and a fix priority, followed by the low-level function-word and n-gram drift as supporting detail and the paragraphs that drift most. `--format json` prints the same data for an LLM to read between rewrites. Both are opt-in, so plain `check` stays fast.

```shell
$ omokage check --author me --explain examples/ja/draft-lost-voice.md
Author: me
Similarity: 0%

High-level style differences (fix these first):
  1. polite sentence-ending ratio is lower than reference [register]
       target 0.000  reference 1.000 ± 0.000  (50.0σ)
  2. kanji ratio is higher than reference [script]
       target 0.489  reference 0.213 ± 0.025  (10.9σ)
  ...

Low-level fingerprint drift (supporting detail):
  - character n-gram "する" is higher than reference  (9.7σ)
  - function word "する" is higher than reference  (9.6σ)
  ...

Paragraphs that drift most:
  #2 (50.0σ; polite sentence-ending ratio lower): 雨天時は在宅で過ごすケースが多い。特段の活動は行わない。ただし窓に当たる降雨音を聴取することで、精神…
```

![explain demo](./doc/img/explain.gif)

## Choosing the author

`check` and `show` resolve the author in this order, so single-author use needs
no flags and multi-author use stays unambiguous:

1. `--author NAME`, if given;
2. otherwise `default_author` from the config;
3. otherwise the only trained profile;
4. otherwise it is an error — zero profiles, or two or more with no default,
   never silently picks one.

Set a default without editing the config by hand:

```shell
$ omokage train --author me --default examples/en/posts
```

## Managing profiles

You never have to touch `profiles/*.db` directly.

```shell
$ omokage list                 # bare names, one per line (pipe-friendly)
me
$ omokage list --long          # trained_at, file count, and source(s)
AUTHOR  TRAINED               FILES  SOURCE
me      2026-06-01 09:14 JST  9      /home/me/writing/posts (+1 more)
$ omokage show --author me      # how a profile was trained (--format json too)
$ omokage rename --author me --to watashi
$ omokage remove --author watashi
```

Trained from several inputs, `list --long` shows the first source with a
`(+N more)` hint and `show` lists them all:

```shell
$ omokage show --author me
Author: me
Trained: 2026-06-01 09:14:32 JST
Files: 9
Sources (2):
  - /home/me/writing/posts
  - /home/me/writing/draft-keeps-voice.md
Documents: 9
Sentences: 142
Characters: 5210
```

`show --format json` exposes the same provenance: read the `sources` array for
the full list. `source_dir` is kept for backward compatibility — it holds the
training directory only when a single directory was used, and is empty otherwise.

`rename` keeps the trained data and refuses to overwrite an existing author;
`remove` clears `default_author` if it pointed at the removed profile.

## Local and global stores

By default omokage looks for an `omokage.toml` by walking up from the current
directory — a project-local store, good for keeping separate writing contexts
apart. For a single voice you can use anywhere, create a per-user store instead:

```shell
$ omokage init --global                       # under $OMOKAGE_HOME or ~/.config/omokage
$ omokage train --global --author me ~/writing
$ cd ~/anywhere && omokage check draft.md      # falls back to the global store
```

When both exist, a local project always wins inside its directory tree; the
global store is the fallback used only when no local project is found. `--global`
forces the global store from anywhere, and `--config PATH` / `--profile-dir PATH`
point omokage at a specific store.

## Scripting

`--score-only` prints just the integer similarity, for shell pipelines:

```shell
$ score=$(omokage check --score-only draft.md)
$ [ "$score" -ge 70 ] && echo "close enough"
```

Use `--format json` (with `check` or `show`) when a tool or LLM needs the full
structured report instead of a single number.

## How it scores

When you train an author, omokage reads every file, measures a set of stylistic features for each document, and stores their mean and spread in a SQLite database under `profiles/` (one database per author). The text itself is not kept, only the numbers.

The features fall into a few groups:

- Register: how often sentences end in the polite form (敬体) or the plain form (常体).
- Script balance: the ratio of kanji, hiragana, and katakana.
- Function words: the frequency of common particles and English function words (の, は, に, the, of, and, …).
- Character n-grams: the most frequent two- and three-character sequences.
- Shape: sentence length and its variation, punctuation and newline frequency, bullet and Markdown usage, paragraph length variation.

A check measures the same features on the draft and compares each one to how much it normally varies across your own writing, as a z-score in the spirit of Burrows's Delta. A feature that stays within your usual range costs nothing; one that strays far lowers the score. The function-word and n-gram fingerprint carries most of the signal, a clear register shift is penalized on top of that, and the shape features only nudge the result. Code blocks are removed before the features are measured, so the score reflects prose rather than code. `diff` uses the same features to compare two documents directly, without a stored profile.

## Limits

omokage looks at style, not meaning. It cannot tell whether a draft is correct, original, or well written, only whether it resembles the voice it was trained on. It needs a reasonable amount of training text: with a few short documents the measured spread is wide and the scores are noisy. It separates Japanese authors more sharply than English ones, and two people who write in the same register will look more alike than they really are. It is not an AI-text detector; the score is similarity to a voice you trained, nothing more.

## About the name

omokage (面影) is a Japanese word. It is written with 面 (face) and 影 (shadow, trace), and it means the remembered image of someone or something, the likeness that comes back to mind. I took the name from [Omokage](https://www.toraya-group.co.jp/products/collections/yokan-omokage), a yokan (red-bean jelly) made by Toraya that I like.

## License

MIT. See [LICENSE](./LICENSE).
