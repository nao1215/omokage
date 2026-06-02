[![Build](https://github.com/nao1215/omokage/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/build.yml)
[![MultiPlatformUnitTest](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml)
[![Coverage](https://github.com/nao1215/omokage/actions/workflows/coverage.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/omokage.svg)](https://pkg.go.dev/github.com/nao1215/omokage)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/omokage)](https://goreportcard.com/report/github.com/nao1215/omokage)
![GitHub](https://img.shields.io/github/license/nao1215/omokage)

# omokage

**Keep the traces of your own writing voice.** `omokage` learns how you write from your past writing, then tells you how strongly a new draft still carries that voice.

![demo](./doc/img/demo.gif)

*omokage* (面影) is a Japanese word for the lingering likeness of someone — the traces of a face or presence that stay with you. `omokage` checks whether the *omokage* of your writing is still there: not what a text says, but whether it still sounds like you.

Everything runs on your machine. Your writing is never uploaded.

## Features

- Learn a writing style from a folder of Markdown or text files.
- Score how closely any draft matches that style, from 0 to 100%.
- See the concrete features that drifted — register, vocabulary, rhythm, and more.
- Compare two documents directly, with no training step.
- Works on both Japanese and English text.
- Local-first: a single static binary with no runtime dependencies.

## Install

```shell
go install github.com/nao1215/omokage@latest
```

Runs on Windows, macOS, and Linux. Building from source needs Go 1.25 or later.

## Quick start

The repository ships with a small example corpus under [`examples/`](./examples) so you can try the whole flow right away.

Create a project in the current directory (this writes `omokage.toml`, `profiles/`, and `cache/`):

```shell
$ omokage init
Initialized omokage project.
```

Learn an author's style from their past writing:

```shell
$ omokage train --author me examples/posts
Trained author "me" from 8 files.
```

Check whether a new draft still keeps that voice:

```shell
$ omokage check --author me examples/draft-keeps-voice.md
Author: me
Similarity: 73%

Differences:
- character n-gram "持ち" is higher than reference
- character n-gram "気持" is higher than reference
- character n-gram "気持ち" is higher than reference
```

The same idea rewritten in a stiff, formal voice drops sharply, and `omokage` shows what changed:

```shell
$ omokage check --author me examples/draft-lost-voice.md
Author: me
Similarity: 0%

Differences:
- polite sentence-ending ratio is lower than reference
- kanji ratio is higher than reference
- hiragana ratio is lower than reference
```

## Compare two drafts

No profile needed — compare any two pieces of writing directly:

```shell
$ omokage diff examples/draft-keeps-voice.md examples/draft-lost-voice.md
Reference: examples/draft-keeps-voice.md
Target: examples/draft-lost-voice.md
Similarity: 54%

Differences:
- polite sentence-ending ratio is lower than reference
- paragraph length variance is lower than reference
- sentence length variance is higher than reference
```

## Reading the result

- **Similarity** is how close the text sits to the learned voice, from 0 to 100%.
- **Differences** lists the few features that moved the most: sentence-ending register (敬体 / 常体), kanji and kana balance, function-word and character n-gram usage, sentence rhythm, and layout.

`omokage` compares *style*, not topic. You can write about something completely new and still score high, as long as the voice is yours.

## Use cases

- Keep one consistent voice across a long-running blog or documentation set.
- Confirm that an edited or co-written draft still sounds like you.
- See whether a rewrite — including an AI-assisted one — preserved your voice or flattened it.
- Discover which stylistic habits make your writing recognizable.

## How it is different

`omokage` is not an "AI text detector". It does not try to judge who or what wrote a text. It measures *likeness to a voice you trained it on* — which is what you want when the goal is to keep your own writing recognizable, however a draft was produced.

## License

MIT. See [LICENSE](./LICENSE).
