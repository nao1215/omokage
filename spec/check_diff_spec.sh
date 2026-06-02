#!/bin/sh
# shellcheck shell=sh
#
# check, diff, and list against trained profiles, in Japanese and English,
# including cross-language behavior and the error paths.

Describe 'omokage check, diff, and list'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  setup_all() {
    fresh_project
    omokage train --author ja_me ja/posts >/dev/null
    omokage train --author en_me en/posts >/dev/null
  }
  BeforeAll 'setup_all'
  AfterAll 'remove_project'

  Describe 'check (Japanese)'
    It 'scores a faithful draft higher than a drifted one'
      keep=$(similarity check --author ja_me ja/keep.md)
      lost=$(similarity check --author ja_me ja/lost.md)
      When call test "$keep" -gt "$lost"
      The status should be success
    End

    It 'names the register shift in a drifted draft'
      When run omokage check --author ja_me ja/lost.md
      The status should be success
      The output should include 'polite sentence-ending ratio is lower than reference'
    End
  End

  Describe 'check (English)'
    It 'scores a faithful draft higher than a different voice'
      keep=$(similarity check --author en_me en/keep.md)
      other=$(similarity check --author en_me en/other.md)
      When call test "$keep" -gt "$other"
      The status should be success
    End
  End

  Describe 'check (cross-language)'
    It 'scores English text low against a Japanese profile'
      value=$(similarity check --author ja_me en/keep.md)
      When call test "$value" -lt 20
      The status should be success
    End

    It 'scores Japanese text low against an English profile'
      value=$(similarity check --author en_me ja/keep.md)
      When call test "$value" -lt 20
      The status should be success
    End
  End

  Describe 'plain check output contract'
    It 'keeps captured output clean on both streams'
      # shellspec captures the streams (not a terminal), so the discoverability
      # tip must be suppressed entirely — scripts and pipes see only the result.
      When run omokage check --author ja_me ja/lost.md
      The status should be success
      The output should include 'Differences:'
      The output should not include 'High-level style'
      The output should not include 'Tip:'
      The stderr should equal ''
    End
  End

  Describe 'check --explain'
    It 'omits the tip in the detailed report itself'
      When run omokage check --author ja_me --explain ja/lost.md
      The status should be success
      The output should not include 'Tip:'
      The stderr should not include 'Tip:'
    End

    It 'leads with the editable high-level register drift'
      When run omokage check --author ja_me --explain ja/lost.md
      The status should be success
      The output should include 'High-level style differences'
      The output should include 'polite sentence-ending ratio is lower than reference'
    End

    It 'localizes register drift to the flipped paragraphs'
      When run omokage check --author ja_me --explain ja/lost.md
      The status should be success
      The output should include 'Paragraphs that drift most:'
      The output should include 'polite sentence-ending ratio lower'
    End

    It 'does not localize a near-match to a document-global feature'
      When run omokage check --author ja_me --explain ja/keep.md
      The status should be success
      The output should not include 'markdown structure frequency lower'
      The output should not include 'paragraph length variance lower'
    End
  End

  Describe 'check --format json'
    It 'emits a machine-readable report with reference statistics'
      When run omokage check --author ja_me --format json ja/lost.md
      The status should be success
      The output should include '"author": "ja_me"'
      The output should include '"similarity":'
      The output should include '"high_level_drift"'
      The output should include '"reference_mean"'
      The output should include '"priority": 1'
    End

    It 'localizes segments with the corresponding feature in JSON'
      When run omokage check --author ja_me --format json ja/lost.md
      The status should be success
      The output should include '"segments"'
      The output should include '"feature": "polite sentence-ending ratio"'
      The output should include '"category": "register"'
    End

    It 'rejects an unknown format'
      When run omokage check --author ja_me --format yaml ja/lost.md
      The status should be failure
      The stderr should include 'unknown --format'
    End
  End

  Describe 'check errors'
    It 'fails on an untrained author'
      When run omokage check --author ghost ja/keep.md
      The status should be failure
      The stderr should include 'profile not found'
    End

    It 'does not leave a stray profile behind for an untrained author'
      omokage check --author ghost ja/keep.md 2>/dev/null || true
      When run omokage list
      The status should be success
      The output should not include 'ghost'
    End

    It 'fails on a missing target file'
      When run omokage check --author ja_me ja/missing.md
      The status should be failure
      The stderr should be present
    End

    It 'errors when the author is ambiguous'
      # With two trained profiles and no default_author, a bare check must not
      # guess: it reports the candidates and asks for --author.
      When run omokage check ja/keep.md
      The status should be failure
      The stderr should include 'multiple profiles'
      The stderr should include 'ja_me'
    End
  End

  Describe 'diff'
    It 'reports 100% for a file against itself'
      When run omokage diff ja/keep.md ja/keep.md
      The status should be success
      The output should include 'Similarity: 100%'
    End

    It 'reports lower similarity across a register change'
      value=$(similarity diff ja/keep.md ja/lost.md)
      When call test "$value" -lt 100
      The status should be success
    End

    It 'fails with the wrong number of arguments'
      When run omokage diff ja/keep.md
      The status should be failure
      The stderr should include 'Usage: omokage diff'
    End

    It 'fails on a missing file'
      When run omokage diff ja/keep.md ja/missing.md
      The status should be failure
      The stderr should be present
    End
  End

  Describe 'list'
    It 'lists the trained authors'
      When run omokage list
      The status should be success
      The output should include 'ja_me'
      The output should include 'en_me'
    End
  End
End
