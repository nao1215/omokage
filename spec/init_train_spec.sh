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

    It 'requires the author flag'
      When run omokage train ja/posts
      The status should be failure
      The stderr should include 'Usage: omokage train'
    End

    It 'fails on a missing directory'
      When run omokage train --author me ja/nope
      The status should be failure
      The stderr should be present
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
