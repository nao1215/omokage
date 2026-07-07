#!/bin/sh
# shellcheck shell=sh
#
# Hermetic atago runner for omokage. It builds the binary and runs the E2E suite
# (e2e/atago/*.atago.yaml) inside a throwaway temp-backed HOME and profile-store
# sandbox, so the suite never reads or writes the developer's real config
# directory and local and CI runs are identical. The test DEFINITIONS are
# plain-YAML atago specs; this script is only the environment bootstrap. Any
# extra arguments are forwarded to `atago run` (for example `--filter check`).
set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if ! command -v atago >/dev/null 2>&1; then
	echo "e2e: atago is not installed. Install it from https://github.com/nao1215/atago" >&2
	echo "e2e: e.g. 'go install github.com/nao1215/atago@latest' (CI uses nao1215/setup-atago)" >&2
	exit 127
fi

# Build the binary the specs exercise; it is exposed on PATH below.
#
# When COVER is set (by scripts/coverage.sh) the binary is built with Go's
# coverage instrumentation instead of the plain `make build`. atago passes the
# environment through to the spec commands, so the omokage child processes
# inherit GOCOVERDIR and write their runtime covdata there. The default (unset
# COVER) path stays byte-for-byte identical, keeping `make test-e2e-atago` fast.
if [ -n "${COVER:-}" ]; then
	: "${GOCOVERDIR:?COVER set but GOCOVERDIR is empty; export GOCOVERDIR to collect e2e coverage}"
	# Mirror the Makefile's VERSION exactly (empty/"dev" when no tags are
	# reachable, e.g. on a shallow CI checkout) so `omokage --version` resolves
	# the same way the plain `make build` binary does.
	VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo dev)"
	env GO111MODULE=on CGO_ENABLED=0 \
		go build -cover -covermode=atomic -coverpkg=./... \
		-ldflags "-X github.com/nao1215/omokage/cmd.Version=${VERSION}" \
		-o omokage main.go
else
	make build
fi

# Create an isolated sandbox and remove it on exit, so no run leaves state behind.
SANDBOX="$(mktemp -d)"
trap 'rm -rf "$SANDBOX"' EXIT INT TERM

mkdir -p "$SANDBOX/home" "$SANDBOX/config" "$SANDBOX/data" "$SANDBOX/cache" "$SANDBOX/ohome" "$SANDBOX/gtest" "$SANDBOX/bin"
cp "$ROOT/omokage" "$SANDBOX/bin/omokage"

# Point HOME and every XDG base directory at the sandbox so any global-store or
# user-config lookup lands there instead of the developer's real home.
# USERPROFILE covers Windows-style home resolution if the suite ever runs there.
HOME="$SANDBOX/home"
export HOME
export USERPROFILE="$SANDBOX/home"
export XDG_CONFIG_HOME="$SANDBOX/config"
export XDG_DATA_HOME="$SANDBOX/data"
export XDG_CACHE_HOME="$SANDBOX/cache"

# Route the global profile store into the sandbox explicitly so a global-store
# lookup never touches a real per-user store. This directory is left empty for
# the whole run, so a no-local-project scenario deterministically reports
# "project not found" rather than depending on a store some earlier scenario
# happened to create.
export OMOKAGE_HOME="$SANDBOX/ohome"

# A separate, also-sandboxed global home used only by the `init --global`
# scenario. Pointing that one scenario here (via env in the spec) keeps its
# write out of both the tracked workdir and the shared OMOKAGE_HOME above, so it
# neither trips its own `changes: created: []` nor pollutes other scenarios.
export OMOKAGE_TEST_GLOBAL="$SANDBOX/gtest"

# The example corpus shipped in the repo is the fixture source for the specs,
# so they do not carry duplicate training data. Exposed as an absolute path
# because atago runs each scenario in its own isolated workdir.
export OMOKAGE_EXAMPLES="$ROOT/examples"

# The freshly built omokage goes first on PATH so the specs exercise that binary.
PATH="$SANDBOX/bin:$PATH"
export PATH

# No `exec`: it would replace the shell and skip the EXIT trap, leaking the
# sandbox. As the last command under `set -e`, atago's exit status is the
# script's exit status either way.
atago run --ci "$@" "$ROOT/e2e/atago"
