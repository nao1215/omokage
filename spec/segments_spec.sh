#!/bin/sh
# shellcheck shell=sh
#
# check --format json segment localization: drift is localized to prose
# paragraphs only. Headings, bullet/table blocks, and other layout are not the
# paragraphs a writer edits for voice, so they must not show up as drifting
# segments.

Describe 'omokage check segment localization'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"
  BeforeEach 'fresh_project'
  AfterEach 'remove_project'

  It 'localizes drift to prose paragraphs, not headings, bullets, or tables'
    omokage train --author me ja/posts >/dev/null
    # A draft whose only prose paragraph drifts into the plain register, preceded by
    # a heading, a bullet list, and a table that must be ignored.
    {
      printf '# 見出しだけの段落\n\n'
      printf -- '- 箇条書きの項目一\n- 箇条書きの項目二\n- 箇条書きの項目三\n\n'
      printf '| 列A | 列B |\n| --- | --- |\n| 値1 | 値2 |\n\n'
      printf '本日は降雨である。外出を実施した。混雑は著しいものであった。状況は変化しない。\n'
    } > "$WORK/draft.md"

    When run omokage check --author me --format json draft.md
    The status should be success
    # The prose paragraph is localized...
    The output should include '"kind": "paragraph"'
    The output should include '本日は降雨である'
    # ...while the heading, bullet, and table text never appear as a drifting segment.
    The output should not include '見出しだけ'
    The output should not include '箇条書き'
    The output should not include '列A'
  End
End
