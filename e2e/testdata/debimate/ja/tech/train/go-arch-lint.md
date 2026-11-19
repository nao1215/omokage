---
title: "【Golang】fe3dback/go-arch-lintでアーキテクチャの破壊を防ぐ"
type: post
date: 2025-02-13
categories:
  - "linux"
cover:
  image: "images/Screenshot-from-2025-02-13-23-26-45.png"
  alt: "【Golang】fe3dback/go-arch-lintでアーキテクチャの破壊を防ぐ"
  hidden: false
---

### 前書き：アーキテクチャは容易に壊される

アーキテクチャリンターである[fe3dback/go-arch-lint](https://github.com/fe3dback/go-arch-lint)を[nao1215/sqly](https://github.com/nao1215/sqly)に導入したので、使用方法のメモを記事として残します。結論としては、初期設定が面倒ですが、期待通りの効果が得られました。なお、既存コードがカオスなアーキテクチャの場合、go-arch-lintを採用できないと思われます。

まず、アーキテクチャをリンターでチェックする発想に至った理由から、説明します。以前、ペアプロ中にドライバ側（実装する人）がアーキテクチャルールに反しているのを偶然目撃しました。違反内容は、「外部サービス操作用パッケージ内でのみ使用できる構造体をユースケースレイヤーから呼び出した」というものです。構造体の定義場所が悪いと思いつつも、ルール違反してしまう理由はアーキテクチャを理解していないからだと考え、対策が必要と考えました（ちなみに、ドライバの方に「その使い方、ダメですよ」と声をかけたら、「そうなんですか？」と返答がありました）

前提条件ですが、当時は以下のような状況で開発していました。

- アーキテクチャに関するドキュメントが存在
- Pull Request（以降PR）単位のレビューで、アーキテクチャルール違反をレビューアが検知
- 必ずしもアーキテクチャルールが遵守されていたわけではない（グレーゾーンや暗黙の了解があった）

実装者のスキルレベルに合わせて、PRレビューの確認観点を意識的に変えるのは、それなりの難しさがあります。レビュー時間もかかります。レビューが長引くと、疲労によって余計な一言をコメントしてしまうリスクも高まります。となると、機械（リンター）ができることは機械にやらせよう、という発想に辿り着きます。機械から指摘された方が、イラッとしませんしね。

---


### リンター候補

[fe3dback/go-arch-lint](https://github.com/fe3dback/go-arch-lint)と[arch-go/arch-go](https://github.com/arch-go/arch-go)が候補でした。どちらもインポート対象パッケージ（依存パッケージ）をチェックする機能があります。

これらのリンターの差分は何でしょうか。go-arch-lintは、依存関係をグラフ化する機能があります。しかし、それ以外はarch-goの方が多機能です。例えば、依存パッケージ内に含められる定義（例：インターフェースのみ）を制御できたり、パラメータや返り値の数、ファイル単位のパブリック関数の数や関数の行数、命名規則のチェックなどができます。また、アーキテクチャを `go test` でき、コンプライアンス（リンター設定）遵守レベルのしきい値チェック機能があります。

しかし、go-arch-lintを採用しました。その理由は、「機能が少ない分、相対的に設定が楽そう」「既存プロジェクトは、arch-goの厳しい設定をパスできない」と考えたからです。プロジェクト特性に合わせて、好きなリンターを選択すれば良いかなといったレベル感です。

---


### fe3dback/go-arch-lintのインストール方法

```
go install github.com/fe3dback/go-arch-lint@latest
```

---


### go-arch-lintの設定

\`.go-arch-lint.yml\` ファイルに設定を書き、プロジェクトのルートディレクトリに配置します。設定読み込みは、\`go-arch-lint check\`を実行すれば、自動的に設定値が反映された状態でリンターが動作します。

\`.go-arch-lint.yml\` に記載する設定項目は、[GitHub](https://github.com/fe3dback/go-arch-lint/blob/master/docs/syntax/README.md)に詳細説明が書かれています。下表に、2025年2月13日時点の各設定項目を示します（GitHubに書かれている説明を訳したもの）

| パス | 必須 | 型 | 説明 |
| --- | --- | --- | --- |
| version | 必須 | int | スキーマバージョン（最新: 3） |
| workdir | 任意 | str | 解析対象の相対ディレクトリ |
| allow | 任意 | map | グローバルルール |
| . depOnAnyVendor | 任意 | bool | プロジェクト内のすべてのファイルが任意のベンダーコードをインポートできるか |
| . deepScan | 任意 | bool | 高度なASTコード解析を使用（v3以降デフォルト \`true\`） |
| exclude | 任意 | \[\]str | 解析対象から除外するディレクトリのリスト（相対パス） |
| excludeFiles | 任意 | \[\]str | ファイル名の正規表現ルール。該当ファイルとそのパッケージを解析対象から除外 |
| components | 必須 | map | Goパッケージの抽象化。1つのコンポーネント = 1つ以上のGoパッケージ |
| . %name% | 必須 | str | コンポーネントの名前 |
| . . in | 必須 | str, \[\]str | 1つ以上の相対ディレクトリ名。グロブパターン対応（例: src/\*/engine/\*\*） |
| vendors | 任意 | map | ベンダーライブラリ（go.mod） |
| . %name% | 必須 | str | ベンダーコンポーネントの名前 |
| . . in | 必須 | str, \[\]str | 1つ以上のベンダーライブラリのインポートパス（例: github.com/abc/\*/engine/\*\*） |
| deps | 必須 | map | 依存関係のルール |
| . %name% | 必須 | str | コンポーネントの名前（"components" セクションで定義したものと同一） |
| . . mayDependOn | 任意 | \[\]str | このコンポーネントがインポート可能なコンポーネントのリスト |
| . . canUse | 任意 | \[\]str | このコンポーネントがインポート可能なベンダーのリスト |

基本的な設定は、以下のような流れで行います。

1. excludeFilesに、除外対象ファイルを設定
2. vendorsに、サードパーティライブラリのエイリアス名を設定
3. componentsに、開発対象パッケージのエイリアス名を設定
4. commonVendorsに、どのパッケージからも呼び出せるサードパーティライブラリ名を設定
5. commonComponentsに、どのパッケージからも呼び出せるパッケージ名（componentsで定義したパッケージ）を設定
6. depsに、各パッケージ（componentsで定義したパッケージ）の依存関係および利用するサードパーティライブラリを設定

---


### 設定例：公式の例、nao1215/sqlyの例

[**公式の設定例**](https://github.com/fe3dback/go-arch-lint/blob/master/.go-arch-lint.yml)

```
version: 3
workdir: internal
allow:
  depOnAnyVendor: false

excludeFiles:
  - "^.*_test\\.go$"
  - "^.*\/test\/.*$"

vendors:
  go-common:           { in: golang.org/x/sync/errgroup }
  go-ast:              { in: [ golang.org/x/mod/modfile, golang.org/x/tools/go/packages ] }
  3rd-cobra:           { in: github.com/spf13/cobra }
  3rd-color-fmt:       { in: github.com/logrusorgru/aurora/v3 }
  3rd-code-highlight:  { in: github.com/alecthomas/chroma/* }
  3rd-json-scheme:     { in: github.com/xeipuuv/gojsonschema }
  3rd-graph:           { in: oss.terrastruct.com/d2/** }
  3rd-yaml:
    in:
      - github.com/goccy/go-yaml
      - github.com/goccy/go-yaml/**
      - github.com/fe3dback/go-yaml    # custom fork (need propose back PR)
      - github.com/fe3dback/go-yaml/** # custom fork (need propose back PR)

components:
  main:                { in: app }
  container:           { in: app/internal/container/** }
  operations:          { in: operations/* }
  services:            { in: services/** }
  view:                { in: view }
  models:              { in: models/** }

commonVendors:
  - go-common

commonComponents:
  - models

deps:
  main:
    mayDependOn:
      - container

  container:
    anyVendorDeps: true
    mayDependOn:
      - operations
      - services
      - view

  operations:
    mayDependOn:
      - services
    canUse:
      - 3rd-graph

  services:
    mayDependOn:
      - services
    canUse:
      - go-ast
      - 3rd-yaml
      - 3rd-color-fmt
      - 3rd-code-highlight
      - 3rd-json-scheme

```

[**nao1215/sqlyでの設定例**](https://github.com/nao1215/sqly/blob/main/.go-arch-lint.yml)

```
version: 3
workdir: .

excludeFiles:
  - "^.*_test\\.go$"
  - "^.*\/test\/.*$"

vendors:
  color: { in: github.com/fatih/color } 
  pflag: { in: github.com/spf13/pflag }
  go-colorable: { in: github.com/mattn/go-colorable }
  xdg: { in: github.com/adrg/xdg }
  env: { in: github.com/caarlos0/env/v6 }
  sqlite: { in: modernc.org/sqlite }
  wire: { in: [github.com/google/wire, github.com/google/wire/cmd/wire] }
  tablewriter: { in: github.com/olekukonko/tablewriter }
  diffmatchpatch: { in: github.com/sergi/go-diff/diffmatchpatch }
  difflib: { in: github.com/pmezard/go-difflib/difflib }
  excelize: { in: github.com/xuri/excelize/v2 }
  gomock: { in: go.uber.org/mock/gomock }
  go-prompt: { in: [github.com/c-bata/go-prompt, github.com/c-bata/go-prompt/completer] }

components:
  cmd:   { in: . }
  shell: { in: shell }
  domain: { in: domain }
  model: { in: domain/model }
  repository: { in: domain/repository }
  infrastructure: { in: [infrastructure, infrastructure/mock/**] }
  memory-infra: { in: infrastructure/memory }
  persistence-infra: { in: infrastructure/persistence }
  usecase: { in: usecase }
  interactor: { in: interactor/** }
  config: { in: config }
  golden: { in: golden }
  di: { in: di }
  mock: { in: [] }

commonVendors:
  - wire
  - gomock
  - color

commonComponents:
  - model
  - golden
  - config

deps:
  di:
    mayDependOn:
      - model
      - shell
      - usecase
      - interactor
      - repository
      - config
      - infrastructure
      - memory-infra
      - persistence-infra
  golden:
    canUse:
      - diffmatchpatch
      - difflib
  config:
    canUse:
      - color
      - pflag
      - go-colorable
      - xdg
      - env
      - sqlite
      - wire
  model:
    canUse:
      - tablewriter
    mayDependOn:
      - domain
  cmd:
    mayDependOn:
      - shell
      - di
  shell:
    canUse:
      - go-prompt
      - tablewriter
      - go-colorable
    mayDependOn:
      - model
      - usecase
  usecase:
    mayDependOn:
      - model
      - repository
  interactor:
    mayDependOn:
      - model
      - usecase
      - repository
  repository:
    mayDependOn:
      - model
  infrastructure:
    mayDependOn:
      - model
  memory-infra:
    mayDependOn:
      - model
      - repository
      - infrastructure
  persistence-infra:
    canUse:
      - excelize
    mayDependOn:
      - model
      - repository
      - infrastructure

```

---


### GitHub Actionsによるアーキテクチャルール違反の検知

`.github/workflows/arch-lint.yml`に以下の設定を書くと、PR作成時にアーキテクチャが期待通りに実装されているかをチェックできます。

```
name: LintArchitecture

on:
  workflow_dispatch:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  check_generate_file:
    name: Lint architecture
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version-file: "go.mod"
          cache-dependency-path: "go.sum"

      - name: Install tools
        run: |
          go install github.com/fe3dback/go-arch-lint@latest
      
      - name: Lint architecture
        run: |
          go-arch-lint check

```

---


### pre-commitによるコミット前チェック

[公式のREADME](https://github.com/fe3dback/go-arch-lint?tab=readme-ov-file#pre-commit)では、[pre-commit](https://pre-commit.com/)を利用して、コミット前にアーキテクチャルール違反チェックおよびグラフ更新チェックを行う方法が提示されています。私は、python環境を構築するのがそれなりに手間だと考えているので、この方法を採用していません。

以下、公式のREADMEに書かれた手順の転載です。

1. pre-commitを[https://pre-commit.com/#install](https://pre-commit.com/#install)からインストール
2. `.pre-commit-config.yaml`ファイルをプロジェクトルートディレクトリに配置し、後述のコードブロックの内容を記載
3. 必要であれば、`args`にフラグを設定
4. 設定を最新バージョンに更新するには、`pre-commit autoupdate`を実行
5. pre-commitを有効にするには、`pre-commit install`を実行

```
repos:
  - repo: https://github.com/fe3dback/go-arch-lint
    rev: master
    hooks:
      - id: go-arch-lint-check
      - id: go-arch-lint-graph
        args: ['--include-vendors=true', '--out=go-arch-lint-graph.svg']

```

---


### 最後に

OSS開発ではクリーンアーキテクチャを採用しないので、リンターでチェックする必要性がないかなと感じています。その一方で、業務の場合は、メンバのスキルやアーキテクチャに対する理解度がマチマチなので、リンターの出番があるかなと。
