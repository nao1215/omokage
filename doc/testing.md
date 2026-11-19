# Testing

How each kind of test in this repository is run, what it covers and where it runs in CI. The README lists the bug classes each one catches.

## Unit tests

```shell
make test
```

The Go tests of every package. CI runs them on Linux, macOS and Windows (`.github/workflows/unit_test.yml`).

`cmd/debimate_regression_test.go` trains on real blog posts and checks, through the same `train` and `check` a user runs, that technical posts score above diary posts, that swapping only the topic words of a post barely moves its score, that rewriting its voice drops the score, and that `doctor` flags a corpus that mixes the two genres. It runs with the unit tests.

## End-to-end tests

```shell
make test-e2e
```

[atago](https://github.com/nao1215/atago) runs the real `omokage` binary from the scenarios in `e2e/atago/`. `e2e/runner` builds the binary into a sandbox and points `OMOKAGE_HOME` and the working directory of each scenario into it, so your own profiles are never read. CI runs the suite on Linux, macOS and Windows (`.github/workflows/e2e_test.yml`); `make coverage` runs it with a coverage-instrumented binary and merges the result with the unit tests (`.github/workflows/coverage.yml`).

## Property-based and metamorphic tests

`internal/profile/property_test.go`, written with [rapid](https://github.com/flyingmutant/rapid). The generator builds documents from Japanese and English sentences, punctuation, list items, headings and blank lines.

| Property | Kind |
|---|---|
| A document compared with itself scores 100%, and so does a draft against a corpus made only of copies of it | property |
| Every comparison and score is within 0..100 | property |
| The profile of a corpus, and the score of a draft against it, do not depend on the order of the documents | metamorphic |
| The same text with CRLF line endings, or with a leading UTF-8 byte order mark, measures the same | metamorphic |

Every `go test` draws 100 cases per property; the nightly `rapid properties` job of `.github/workflows/fuzz.yml` draws 20 000. A failure prints the seed and writes the case under `internal/profile/testdata/rapid/`.

```shell
go test ./internal/profile -run '^TestProperty' -rapid.checks=20000
```

## Fuzzing

`go test` runs every `Fuzz*` target over its seed corpus and `testdata/fuzz/`. `.github/workflows/fuzz.yml` discovers the targets and fuzzes each for five minutes every night.

| Target | Input |
|---|---|
| `FuzzExtractText` (`internal/feature/extractor_test.go`) | Markdown and plain text of a draft |
| `FuzzParseRoundTrip` (`internal/config/config_test.go`) | `omokage.toml` |
| `FuzzNormalizeKey`, `FuzzScanCandidates` (`internal/term/term_test.go`) | Term preference keys and the text scanned for them |

```shell
go test -run '^$' -fuzz '^FuzzExtractText$' -fuzztime 1m ./internal/feature
```

## Benchmarks

```shell
make bench
make bench-compare
```

[himorime](https://github.com/nao1215/himorime) measures the binary end to end from `bench/himorime.yaml`. On a pull request `.github/workflows/bench.yml` builds the base and the head and measures them in interleaved rounds on the same runner; see `bench/README.md`.
