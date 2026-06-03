#!/bin/sh
# shellcheck shell=sh
#
# show quality: the corpus-quality information a trained profile carries, and the
# lighter --summary JSON for handing to an LLM. These pin the user-facing contract
# that what `doctor` found at train time stays visible through `show`.

Describe 'omokage show quality'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"
  BeforeEach 'fresh_project'
  AfterEach 'remove_project'

  It 'reports reliability in the text view'
    omokage train --author me ja/keep.md ja/lost.md >/dev/null
    When run omokage show --author me
    The status should be success
    The output should include 'Reliability:'
  End

  It 'keeps the train-time findings in --format json, including outliers'
    # A larger, register-mixed corpus so doctor-style findings (mixed voice and an
    # outlier the stored distribution alone could not reproduce) are recorded and
    # then surface through show.
    omokage train --author me ja/posts >/dev/null
    When run omokage show --author me --format json
    The status should be success
    The output should include '"reliability"'
    The output should include '"quality_findings"'
    The output should include '"term_preferences"'
  End

  It 'omits term_preferences with --summary for a lighter JSON'
    omokage train --author me ja/posts >/dev/null
    When run omokage show --author me --format json --summary
    The status should be success
    The output should include '"reliability"'
    The output should include '"quality_findings"'
    The output should not include '"term_preferences"'
  End
End
