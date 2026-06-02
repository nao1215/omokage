#!/bin/sh
# shellcheck shell=sh
#
# CLI surface: help text, version, unknown commands, and per-subcommand help.
# These do not need a project, so they run the binary directly.

Describe 'omokage CLI surface'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  Describe 'root help'
    It 'prints help with no arguments'
      When run "$OMOKAGE_BIN"
      The status should be success
      The output should include 'omokage analyzes writing style'
      The output should include 'Commands:'
      The output should include 'check'
    End

    It 'prints help with --help'
      When run "$OMOKAGE_BIN" --help
      The status should be success
      The output should include 'Usage:'
      The output should include 'train'
    End
  End

  Describe 'version'
    It 'prints the version'
      When run "$OMOKAGE_BIN" version
      The status should be success
      The output should include 'omokage'
    End
  End

  Describe 'unknown command'
    It 'fails and shows the help'
      When run "$OMOKAGE_BIN" frobnicate
      The status should be failure
      The stderr should include 'unknown command'
      The stdout should include 'Commands:'
    End
  End

  Describe 'subcommand help describes the command'
    It 'describes check'
      When run "$OMOKAGE_BIN" check --help
      The status should be failure
      The stderr should include 'Score how closely'
      The stderr should include 'Usage: omokage check'
    End

    It 'describes train'
      When run "$OMOKAGE_BIN" train --help
      The status should be failure
      The stderr should include 'Learn an author'
    End
  End
End
