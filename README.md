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

omokage (面影) is a Japanese word. It is written with 面 (face) and 影 (shadow, trace), and it means the remembered image of someone or something, the likeness that comes back to mind.

The name also comes from Omokage, a yokan (red-bean jelly) made by Toraya that I like (https://www.toraya-group.co.jp/products/collections/yokan-omokage). That was part of why I picked it.

## Why I built it

I often draft text with an LLM and then rework it so that it reads like something I wrote. Prompting and hand-editing only get me so far, so I wanted a tool that measures how close a draft is to my own style and points out where it drifts. omokage is that tool. You train it on your past writing, then check a draft against it.

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

## License

MIT. See [LICENSE](./LICENSE).
