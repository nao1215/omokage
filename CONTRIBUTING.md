# Contributing Guide

## Development Environment

- Go 1.25 or later
- `make`
- `git`

Clone the repository and install the optional helper tools:

```bash
make tools
```

## Common Commands

```bash
make build
make test
make lint
```

`make test` writes `coverage.out` and `coverage.html` to the repository root.

## Pull Request Expectations

- keep CLI behavior and error messages consistent
- add or update tests for new behavior
- run `make test` before opening a PR
- run `make lint` when changing Go code

## CI

GitHub Actions runs the following workflows:

- `build.yml`: verifies the project builds on Linux
- `unit_test.yml`: runs `go test ./...` on Linux, macOS, and Windows
- `coverage.yml`: uploads Octocov coverage reports
- `reviewdog.yml`: comments on lint, misspell, and workflow issues in pull requests

## Design Notes

- author profiles are stored as separate SQLite databases under `profiles/`
- `doc/reference/` is treated as a local reference area and is not part of commits
- the CLI currently focuses on `.md` and `.txt` corpora
