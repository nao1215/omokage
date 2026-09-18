# Contributing Guide

## Development Environment

- Go 1.25 or later
- `make`
- `git`

Clone the repository and install the helper tools. This installs the linter,
the coverage tool, and `shellspec` (used by the end-to-end tests):

```bash
make tools
```

`make tools` installs shellspec under `~/.local`, so make sure `~/.local/bin` is
on your `PATH`.

## Common Commands

```bash
make build      # build the omokage binary
make test       # unit tests with coverage (writes coverage.out / coverage.html)
make test-e2e   # shellspec end-to-end tests against the built binary
make bench      # measure the CLI end to end with himorime (bench/README.md)
make lint       # golangci-lint
```

The end-to-end tests live under `spec/` and exercise the built binary the way a
user does, using the fixtures in `spec/testdata/` (a Japanese and an English
corpus). Run them with `make test-e2e`, or directly with `shellspec --shell sh`
after `make build`.

`make bench` measures `train`, `check` and `diff` end to end with [himorime](https://github.com/nao1215/himorime), and `make bench-compare` compares main with your working tree on the same machine. The benchmark CI does the same with the base of each pull request and fails on a confirmed slowdown, CPU time or memory increase, so run `make bench-compare` locally when you touch feature extraction or scoring. `bench/README.md` lists what is measured.

## Pull Request Expectations

- keep CLI behavior and error messages consistent
- add or update tests for new behavior, including a `spec/` test for CLI changes
- run `make test` and `make test-e2e` before opening a PR
- run `make lint` when changing Go code
- run `make bench-compare` when changing feature extraction or scoring

## CI

GitHub Actions runs the following workflows, and every gate is reproducible
locally with the `make` targets above:

- `build.yml`: verifies the project builds on Linux
- `unit_test.yml`: runs `go test ./...` on Linux, macOS, and Windows (`make test`)
- `e2e_test.yml`: runs the shellspec end-to-end tests on Linux and macOS (`make test-e2e`)
- `bench.yml`: runs `himorime ci bench`, which measures the base and the head of a pull request on one runner and fails on a confirmed regression (`make bench-compare`)
- `coverage.yml`: uploads Octocov coverage reports
- `reviewdog.yml`: comments on lint, misspell, and workflow issues in pull requests

## Design Notes

- author profiles are stored as separate SQLite databases under `profiles/`
- the CLI currently focuses on `.md` and `.txt` corpora
