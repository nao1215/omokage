#!/bin/sh
# shellcheck shell=sh
#
# run.sh is a thin wrapper kept for muscle memory and for anything that still
# calls it by path. The bootstrap itself moved to e2e/runner, a Go program: the
# end-to-end suite runs on Windows now, and a POSIX bootstrap would have made the
# Windows leg depend on Git for Windows being installed, which tests the runner
# image rather than omokage.
#
# Usage: e2e/run.sh [atago args...]        (e.g. e2e/run.sh --filter check)
set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
exec go run ./e2e/runner "$@"
