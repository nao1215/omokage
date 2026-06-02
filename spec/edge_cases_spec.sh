#!/bin/sh
# shellcheck shell=sh
#
# Adversarial cases: empty inputs, a directory passed where a file is expected,
# path-separator injection in an author name, and mixed file extensions. These
# are meant to break the binary rather than confirm the happy path.

Describe 'omokage edge cases'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  BeforeEach 'fresh_project'
  AfterEach 'remove_project'

  It 'checks an empty target file without crashing'
    : > "$WORK/empty.md"
    omokage train --author me ja/posts >/dev/null
    When run omokage check --author me empty.md
    The status should be success
    The output should include 'Similarity:'
  End

  It 'fails when the target is a directory, not a file'
    omokage train --author me ja/posts >/dev/null
    When run omokage check --author me ja
    The status should be failure
    The stderr should be present
  End

  It 'rejects an author name containing a path separator'
    When run omokage train --author ../evil ja/posts
    The status should be failure
    The stderr should include 'path separators'
    The path "$WORK/../evil.db" should not be exist
  End

  It 'trains only from supported extensions in a mixed directory'
    mkdir "$WORK/mixed"
    cp "$WORK/ja/posts/coffee.md" "$WORK/mixed/a.md"
    printf '{"not":"prose"}\n' > "$WORK/mixed/b.json"
    printf 'Just a short line of plain text here.\n' > "$WORK/mixed/c.txt"
    When run omokage train --author mixed mixed
    The status should be success
    The output should include 'Trained author "mixed"'
  End

  It 'diffs two empty files without crashing'
    : > "$WORK/a.md"
    : > "$WORK/b.md"
    When run omokage diff a.md b.md
    The status should be success
    The output should include 'Similarity:'
  End
End
