#!/bin/sh
# shellcheck shell=sh
#
# Missing required-argument handling: each affected command must print a direct
# "missing X" line on stderr (not just the usage block) so the user sees what is
# absent. The validation runs before any store lookup, so these need no project
# and can drive the binary directly.

Describe 'omokage missing-argument errors'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  It 'check without a FILE reports the missing argument'
    When run "$OMOKAGE_BIN" check
    The status should be failure
    The stderr should include 'missing FILE'
    The stderr should include 'Usage: omokage check'
  End

  It 'diff with only FILE_A reports the missing second file'
    When run "$OMOKAGE_BIN" diff a.md
    The status should be failure
    The stderr should include 'missing FILE_B'
    The stderr should include 'Usage: omokage diff'
  End

  It 'train with --author but no directory reports the missing directory'
    When run "$OMOKAGE_BIN" train --author me
    The status should be failure
    The stderr should include 'missing DIRECTORY'
    The stderr should include 'Usage: omokage train'
  End

  It 'train with a directory but no --author reports the missing flag'
    When run "$OMOKAGE_BIN" train examples/posts
    The status should be failure
    The stderr should include 'missing --author'
    The stderr should include 'Usage: omokage train'
  End

  It 'remove without --author reports the missing flag'
    When run "$OMOKAGE_BIN" remove
    The status should be failure
    The stderr should include 'missing --author'
    The stderr should include 'Usage: omokage remove'
  End

  It 'rename with --author but no --to reports the missing flag'
    When run "$OMOKAGE_BIN" rename --author me
    The status should be failure
    The stderr should include 'missing --to'
    The stderr should include 'Usage: omokage rename'
  End
End
