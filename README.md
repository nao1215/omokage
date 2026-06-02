<p align="center">
  <img src="doc/img/omokage-icon.jpg" alt="omokage" width="320">
</p>

[![Build](https://github.com/nao1215/omokage/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/build.yml)
[![MultiPlatformUnitTest](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/reviewdog.yml)
[![Coverage](https://github.com/nao1215/omokage/actions/workflows/coverage.yml/badge.svg)](https://github.com/nao1215/omokage/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/omokage.svg)](https://pkg.go.dev/github.com/nao1215/omokage)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/omokage)](https://goreportcard.com/report/github.com/nao1215/omokage)
![GitHub](https://img.shields.io/github/license/nao1215/omokage)

# omokage

omokage learns how you write from your past writing, then tells you how close a new draft is to that style. It runs locally and works on Japanese and English text.

![demo](./doc/img/demo.gif)

## About the name

omokage (面影) is a Japanese word. It is written with 面 (face) and 影 (shadow, trace), and it means the remembered image of someone or something, the likeness that comes back to mind. I took the name from [Omokage](https://www.toraya-group.co.jp/products/collections/yokan-omokage), a yokan (red-bean jelly) made by Toraya that I like.

## Why I built it

I often draft text with an LLM and then rework it so that it reads like something I wrote. Prompting and hand-editing only get me so far, so I wanted a tool that measures how close a draft is to my own style and points out where it drifts. omokage is that tool. You train it on your past writing, then check a draft against it.

It is meant to be used by an LLM as much as by a person. An agent can run `check` after each rewrite, read the similarity and the differences, and keep revising until the draft sits closer to the trained voice.

## Install

```shell
go install github.com/nao1215/omokage@latest
```

It runs on Windows, macOS, and Linux. Building from source needs Go 1.25 or later.

## Usage

The repository includes a small example corpus under [examples/](./examples) so you can follow along.

Create a project in the current directory. This writes `omokage.toml`, `profiles/`, and `cache/`.

```shell
$ omokage init
Initialized omokage project.
```

Learn a style from past writing.

```shell
$ omokage train --author me examples/posts
Trained author "me" from 8 files.
```

Check whether a draft still reads like that author.

```shell
$ omokage check --author me examples/draft-keeps-voice.md
Author: me
Similarity: 73%

Differences:
- character n-gram "持ち" is higher than reference
- character n-gram "気持" is higher than reference
- character n-gram "気持ち" is higher than reference
```

The same idea rewritten in a stiff, formal voice scores low, and omokage shows what changed.

```shell
$ omokage check --author me examples/draft-lost-voice.md
Author: me
Similarity: 0%

Differences:
- polite sentence-ending ratio is lower than reference
- kanji ratio is higher than reference
- hiragana ratio is lower than reference
```

You can also compare two documents directly, without training a profile.

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

Similarity runs from 0 to 100 and shows how close the text sits to the learned style. Differences lists the features that moved the most, such as the sentence-ending register (敬体 / 常体), the balance of kanji and kana, function-word and character n-gram usage, sentence length, and layout. omokage compares style rather than topic, so writing about something new in your usual voice still scores high.

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

## License

MIT. See [LICENSE](./LICENSE).
