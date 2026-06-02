#!/bin/sh
# shellcheck shell=sh
#
# shellspec helper for omokage end-to-end tests. These drive the built binary
# the way a user does (subcommands, flags, exit codes, files on disk) so they
# catch regressions the Go unit tests cannot. Each test runs inside a throwaway
# project directory created with mktemp, so nothing touches the repository.

set -eu

PROJECT_ROOT="$(cd "$SHELLSPEC_SPECDIR/.." && pwd)"
export PROJECT_ROOT

# OMOKAGE_BIN points at the binary built by `make build`. Override to test
# another build.
OMOKAGE_BIN="${OMOKAGE_BIN:-$PROJECT_ROOT/omokage}"
export OMOKAGE_BIN

# omokage runs the built binary inside the current project directory ($WORK) so
# that omokage.toml, profiles/, and relative fixture paths resolve there.
omokage() {
  ( cd "$WORK" && "$OMOKAGE_BIN" "$@" )
}

# similarity prints just the integer percentage from a check/diff invocation.
similarity() {
  omokage "$@" | awk '/Similarity/ { value = $2; gsub(/[^0-9]/, "", value); print value }'
}

# make_workdir creates a fresh project directory seeded with the test fixtures
# (the ja/ and en/ corpora under spec/testdata) but without running init.
make_workdir() {
  WORK="$(mktemp -d)"
  export WORK
  cp -R "$PROJECT_ROOT/spec/testdata/." "$WORK/"
}

init_project() {
  omokage init >/dev/null
}

# fresh_project is the common setup: a new working directory with an initialized
# omokage project.
fresh_project() {
  make_workdir
  init_project
}

remove_project() {
  if [ -n "${WORK:-}" ]; then
    rm -rf "$WORK"
  fi
}
