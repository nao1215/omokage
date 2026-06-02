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
    It 'describes check and lists its flags with double dashes'
      When run "$OMOKAGE_BIN" check --help
      The status should be failure
      The stderr should include 'Score how closely'
      The stderr should include 'Usage: omokage check'
      The stderr should include '--author'
      The stderr should include '--explain'
      The stderr should include '--format'
    End

    It 'surfaces the explain and json options in check help'
      When run "$OMOKAGE_BIN" check --help
      The status should be failure
      The stderr should include 'prioritized'
      The stderr should include 'text or json'
    End

    It 'does not show single-dash flag spellings in check help'
      When run "$OMOKAGE_BIN" check --help
      The status should be failure
      The stderr should not include '  -author'
    End

    It 'describes init with a double-dash flag'
      When run "$OMOKAGE_BIN" init --help
      The status should be failure
      The stderr should include 'Usage: omokage init'
      The stderr should include '--name'
    End

    It 'describes train'
      When run "$OMOKAGE_BIN" train --help
      The status should be failure
      The stderr should include 'Learn an author'
      The stderr should include '--author'
    End
  End
End
