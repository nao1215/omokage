#!/bin/sh
# shellcheck shell=sh
#
# init and train: project creation and learning a profile, plus the error and
# boundary cases (missing flags, missing/empty directories, empty corpora).

Describe 'omokage init and train'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  Describe 'init'
    BeforeEach 'make_workdir'
    AfterEach 'remove_project'

    It 'creates a project'
      When run omokage init
      The status should be success
      The output should include 'Initialized omokage project.'
      The path "$WORK/omokage.toml" should be file
      The path "$WORK/profiles" should be directory
    End

    It 'accepts a project name'
      When run omokage init --name my-style
      The status should be success
      The output should include 'Initialized omokage project.'
      The contents of file "$WORK/omokage.toml" should include 'my-style'
    End

    It 'refuses to re-initialize an existing project'
      init_project
      When run omokage init
      The status should be failure
      The stderr should include 'already exists'
    End
  End

  Describe 'train'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    It 'trains a Japanese profile'
      When run omokage train --author me ja/posts
      The status should be success
      The output should include 'Trained author "me"'
      The path "$WORK/profiles/me.db" should be file
    End

    It 'trains an English profile'
      When run omokage train --author me en/posts
      The status should be success
      The output should include 'Trained author "me"'
    End

    It 'trains from several files passed directly'
      When run omokage train --author me ja/keep.md ja/lost.md
      The status should be success
      The output should include 'Trained author "me" from 2 files'
      The path "$WORK/profiles/me.db" should be file
    End

    It 'trains from a directory and a file mixed together'
      # ja/posts holds 8 files; adding ja/keep.md makes 9.
      When run omokage train --author me ja/posts ja/keep.md
      The status should be success
      The output should include 'Trained author "me" from 9 files'
    End

    It 'does not double-count a file already inside a passed directory'
      # ja/posts holds 8 files; ja/posts/walk.md is one of them, so the count stays 8.
      When run omokage train --author me ja/posts ja/posts/walk.md
      The status should be success
      The output should include 'Trained author "me" from 8 files'
    End

    It 'requires the author flag'
      When run omokage train ja/posts
      The status should be failure
      The stderr should include 'Usage: omokage train'
    End

    It 'fails with a clear message when no INPUT is given'
      When run omokage train --author me
      The status should be failure
      The stderr should include 'missing INPUT'
    End

    It 'fails on a missing directory'
      When run omokage train --author me ja/nope
      The status should be failure
      The stderr should be present
    End

    It 'names the missing input path'
      When run omokage train --author me ja/nope.md
      The status should be failure
      The stderr should include 'input not found'
    End

    It 'rejects a directly passed unsupported file'
      printf 'not text\n' > "$WORK/notes.pdf"
      When run omokage train --author me notes.pdf
      The status should be failure
      The stderr should include 'unsupported file'
    End

    It 'rejects a URL input with a clear message'
      When run omokage train --author me https://example.com/post
      The status should be failure
      The stderr should include 'URL inputs are not supported'
    End

    It 'stops on a bad input mixed with valid ones and trains nothing'
      # A single bad input aborts the whole run: no profile is written, so the
      # user can drop the named argument and re-run.
      When run omokage train --author me ja/posts ja/nope.md
      The status should be failure
      The stderr should include 'input not found: ja/nope.md'
      The path "$WORK/profiles/me.db" should not be exist
    End

    It 'fails when the directory has no supported files'
      mkdir "$WORK/empty"
      When run omokage train --author me empty
      The status should be failure
      The stderr should include 'no supported files'
    End

    It 'rejects a corpus of only empty files'
      mkdir "$WORK/blank"
      : > "$WORK/blank/a.md"
      : > "$WORK/blank/b.md"
      When run omokage train --author me blank
      The status should be failure
      The stderr should include 'no usable text'
    End
  End
End
