#!/bin/sh
# shellcheck shell=sh
#
# Everyday-use ergonomics: single-profile and default_author resolution for
# check, the profile-management commands (remove/rename/show), the list output
# modes, the script-friendly --score-only, and local/global precedence.

Describe 'omokage ergonomics'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  Describe 'check author resolution'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    It 'auto-selects the only profile when --author is omitted'
      omokage train --author me ja/posts >/dev/null
      When run omokage check ja/keep.md
      The status should be success
      The output should include 'Author: me'
      The output should include 'Similarity:'
    End

    It 'uses default_author when several profiles exist'
      omokage train --author me ja/posts >/dev/null
      omokage train --author other --default ja/posts >/dev/null
      When run omokage check ja/keep.md
      The status should be success
      The output should include 'Author: other'
    End

    It 'lets an explicit --author override the default'
      omokage train --author me ja/posts >/dev/null
      omokage train --author other --default ja/posts >/dev/null
      When run omokage check --author me ja/keep.md
      The status should be success
      The output should include 'Author: me'
    End

    It 'errors when no profile has been trained yet'
      When run omokage check ja/keep.md
      The status should be failure
      The stderr should include 'no author profiles'
    End
  End

  Describe 'check --score-only'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    It 'prints only the integer score'
      omokage train --author me ja/posts >/dev/null
      When run omokage check --score-only ja/keep.md
      The status should be success
      The output should not include 'Author'
      The output should not include 'Similarity'
      The output should not include '%'
    End

    It 'refuses to combine with --explain'
      omokage train --author me ja/posts >/dev/null
      When run omokage check --score-only --explain ja/keep.md
      The status should be failure
      The stderr should include 'cannot be combined'
    End
  End

  Describe 'profile management'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    It 'shows how a profile was trained'
      omokage train --author me ja/posts >/dev/null
      When run omokage show --author me
      The status should be success
      The output should include 'Author: me'
      The output should include 'Files:'
      The output should include 'Source:'
    End

    It 'removes a profile without touching the filesystem directly'
      omokage train --author me ja/posts >/dev/null
      When run omokage remove --author me
      The status should be success
      The output should include 'Removed author "me"'
      The path "$WORK/profiles/me.db" should not be exist
    End

    It 'fails to remove an author that does not exist'
      When run omokage remove --author ghost
      The status should be failure
      The stderr should include 'profile not found'
    End

    It 'renames a profile and keeps its data'
      omokage train --author me ja/posts >/dev/null
      omokage rename --author me --to watashi >/dev/null
      When run omokage show --author watashi
      The status should be success
      The output should include 'Author: watashi'
    End

    It 'refuses to rename onto an existing author'
      omokage train --author me ja/posts >/dev/null
      omokage train --author other ja/posts >/dev/null
      When run omokage rename --author me --to other
      The status should be failure
      The stderr should include 'already exists'
    End
  End

  Describe 'list output modes'
    BeforeEach 'fresh_project'
    AfterEach 'remove_project'

    It 'lists bare names by default'
      omokage train --author me ja/posts >/dev/null
      When run omokage list
      The status should be success
      The output should equal 'me'
    End

    It 'shows details with --long'
      omokage train --author me ja/posts >/dev/null
      When run omokage list --long
      The status should be success
      The output should include 'AUTHOR'
      The output should include 'TRAINED'
      The output should include 'SOURCE'
      The output should include 'me'
    End
  End

  Describe 'global store'
    setup_global() {
      make_workdir
      OMOKAGE_HOME="$WORK/home"
      export OMOKAGE_HOME
    }
    BeforeEach 'setup_global'
    AfterEach 'remove_project'

    It 'initializes and uses a global store with no local project'
      omokage init --global >/dev/null
      omokage train --global --author me ja/posts >/dev/null
      # No local omokage.toml in $WORK, so a bare check falls back to the global
      # store and auto-selects the single profile.
      When run omokage check ja/keep.md
      The status should be success
      The output should include 'Author: me'
    End

    It 'prefers a local project over the global store'
      omokage init --global >/dev/null
      omokage train --global --author global_author ja/posts >/dev/null
      omokage init >/dev/null
      omokage train --author local_author ja/posts >/dev/null
      When run omokage list
      The status should be success
      The output should include 'local_author'
      The output should not include 'global_author'
    End
  End
End
