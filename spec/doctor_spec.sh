#!/bin/sh
# shellcheck shell=sh
#
# doctor: the corpus-quality check, and the post-training quality note. These
# drive the built binary the way a user curating a corpus does, so they catch
# regressions in the user-facing wording and exit codes, not just the logic.

Describe 'omokage doctor'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  Describe 'reporting'
    BeforeEach 'make_workdir'
    AfterEach 'remove_project'

    It 'reports a solid corpus as good without needing init'
      # doctor reads and reports only; it never trains or writes a store, so it
      # works in a bare directory with no omokage project.
      When run omokage doctor ja/posts
      The status should be success
      The output should include 'Reliability: good'
      The output should include 'No problems found'
      The output should include 'not writing quality'
      The path "$WORK/omokage.toml" should not be exist
    End

    It 'flags a thin corpus and suggests what to do'
      mkdir "$WORK/thin"
      printf '今日は晴れ。散歩した。\n' > "$WORK/thin/a.md"
      printf '本を読んだ。良かった。\n' > "$WORK/thin/b.md"
      When run omokage doctor thin
      The status should be success
      The output should include 'Reliability: weak'
      The output should include 'Findings:'
      The output should include '→'
    End

    It 'emits machine-readable JSON'
      When run omokage doctor --format json en/posts
      The status should be success
      The output should include '"reliability"'
      The output should include '"findings"'
      The output should include '"action"'
    End
  End

  Describe 'errors'
    BeforeEach 'make_workdir'
    AfterEach 'remove_project'

    It 'fails with a clear message when no INPUT is given'
      When run omokage doctor
      The status should be failure
      The stderr should include 'missing INPUT'
    End

    It 'rejects an unknown format'
      When run omokage doctor --format yaml ja/posts
      The status should be failure
      The stderr should include 'unknown --format'
    End

    It 'rejects a URL input like train does'
      When run omokage doctor https://example.com/post
      The status should be failure
      The stderr should include 'URL inputs are not supported'
    End
  End

  Describe 'post-training summary'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    # The corpus-reliability summary prints on stdout (not gated behind a terminal),
    # so a person, a script, and an LLM all see whether to curate the corpus. A thin
    # corpus reports weak and points at doctor; stderr stays clean.
    It 'reports a thin corpus as weak and points at doctor'
      mkdir "$WORK/thin"
      printf '今日は晴れ。散歩した。\n' > "$WORK/thin/a.md"
      printf '本を読んだ。良かった。\n' > "$WORK/thin/b.md"
      When run omokage train --author thin thin
      The status should be success
      The output should include 'Trained author "thin"'
      The output should include 'Corpus reliability: weak'
      The output should include 'omokage doctor'
      The stderr should equal ''
    End

    It 'reports a solid corpus as good without pointing at doctor'
      When run omokage train --author me ja/posts
      The status should be success
      The output should include 'Corpus reliability: good.'
      The output should not include 'omokage doctor'
    End
  End
End
