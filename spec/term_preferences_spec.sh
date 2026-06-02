#!/bin/sh
# shellcheck shell=sh
#
# Term preferences end-to-end: train learns per-profile notation preferences,
# `show --format json` exposes them (including a corpus-declared alias bridge),
# and `check --format json` flags a draft that uses a non-preferred surface,
# without changing the similarity score. Also covers the prose-only contract:
# code, diagrams, HTML, and front matter never surface as drifting paragraphs.

Describe 'omokage term preferences'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  # seed_terms builds a corpus where DB is used far more than データベース, and the
  # parenthetical "データベース（DB）" declares the alias bridge between them.
  seed_terms() {
    fresh_project
    mkdir "$WORK/terms"
    printf 'データベース（DB）を使う。DB は速い。DB を使う。DB が良い。\n' > "$WORK/terms/a.md"
    printf 'DB を使う。DB が好き。DB を選ぶ。データベースも使う。\n' > "$WORK/terms/b.md"
    omokage train --author tm terms >/dev/null
  }
  BeforeAll 'seed_terms'
  AfterAll 'remove_project'

  Describe 'show --format json'
    It 'reports the bridged group with separate normalized_key and group_key'
      When run omokage show --author tm --format json
      The status should be success
      The output should include '"term_preferences"'
      The output should include '"group_key": "term:db"'
      The output should include '"preferred_surface": "DB"'
      The output should include '"normalized_key"'
      # データベース is bridged into the same group as DB.
      The output should include 'データベース'
    End
  End

  Describe 'show text output'
    It 'stays a short provenance summary without a term dump'
      When run omokage show --author tm
      The status should be success
      The output should include 'Author: tm'
      The output should not include 'term_preferences'
      The output should not include 'group_key'
    End
  End

  Describe 'check --format json'
    It 'flags a non-preferred surface as a term warning'
      printf 'ＤＢ を整備する。データベースを設計する。\n' > "$WORK/draft.md"
      When run omokage check --author tm --format json draft.md
      The status should be success
      The output should include '"term_warnings"'
      The output should include '"preferred_surface": "DB"'
      The output should include '"used_surface"'
      # The score layer is still present and unaffected.
      The output should include '"similarity":'
    End

    It 'does not warn when the draft uses the preferred surface'
      printf 'DB を整備する。DB を設計する。\n' > "$WORK/clean.md"
      When run omokage check --author tm --format json clean.md
      The status should be success
      The output should include '"term_warnings": []'
    End
  End

  Describe 'plain check is unchanged by the term layer'
    It 'keeps the plain check output and exit code intact'
      printf 'ＤＢ を整備する。\n' > "$WORK/draft.md"
      When run omokage check --author tm draft.md
      The status should be success
      The output should include 'Similarity:'
      The output should not include 'term_warnings'
      The stderr should equal ''
    End
  End
End

Describe 'omokage prose-only extraction'
  Include "$SHELLSPEC_SPECDIR/spec_helper.sh"

  setup_prose() {
    fresh_project
    omokage train --author pm ja/posts >/dev/null
  }
  BeforeAll 'setup_prose'
  AfterAll 'remove_project'

  Describe 'check --explain ignores non-prose'
    write_mixed_draft() {
      # Prose interleaved with a fenced mermaid diagram that contains a blank
      # line, an HTML block, and YAML front matter — none of which is prose.
      cat > "$WORK/mixed.md" <<'EOF'
---
title: テスト記事
image: images/cover.jpg
---

最初の段落です。これは普通の本文であり、敬体で書かれています。

```mermaid
flowchart TD
    subgraph Linter["omokage"]

    A --> B
    end
```

<p align="center">
  <img src="images/x.jpg" alt="omokage" width="280">
</p>

詳しくは [データベース入門](https://example.com/db.html) を参照してください。

最後の段落です。これも本文として扱われるべき部分です。
EOF
    }
    BeforeEach 'write_mixed_draft'

    It 'never reports a diagram, HTML, or front matter as a drifting paragraph'
      When run omokage check --author pm --explain mixed.md
      The status should be success
      The output should not include 'subgraph'
      The output should not include 'flowchart'
      The output should not include '<img'
      The output should not include 'align='
      The output should not include 'title:'
      The output should not include 'images/cover.jpg'
    End

    It 'does not leak a link URL into the JSON report'
      When run omokage check --author pm --format json mixed.md
      The status should be success
      The output should not include 'example.com'
      The output should not include 'db.html'
    End
  End
End
