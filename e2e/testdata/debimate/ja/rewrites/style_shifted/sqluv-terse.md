---
title: "sqluv 機能追加メモ"
type: post
date: 2025-03-22
categories:
  - "linux"
---

### 概要

nao1215/sqluv を更新した。RDBMS と CSV／TSV／LTSV を扱う TUI クライアントである。https、S3、圧縮フォーマットに対応した。カラースキームも追加した。

### 対応フォーマット

v0.3.0 時点の対応は次のとおり。MySQL、PostgreSQL、SQLite3、SQL Server。CSV、TSV、LTSV。これらに SQL を実行できる。クエリ履歴も残る。

CSV／TSV／LTSV はローカルから読める。http、https、Amazon S3 からも読める。圧縮ファイルも読める。対象は .gz、.xz、.bz2、.zst の4種だ。

### UI

旧版は全カラムを1画面に出していた。読めなかった。新版は横スクロールする。サイドバーでテーブルを選ぶ。カラム情報が出る。これで足りる。

### 履歴

SQL の実行が成功する。sqluv は SQLite3 にクエリを保存する。History から参照する。テキストエリアにコピーする。それだけだ。

### リモート読み込み

sqluv は引数にパスを取る。スキームが無ければ file:// とみなす。https:// と s3:// はダウンロードしてから読む。拡張子でフォーマットを判定する。判定できなければエラーを返す。

DuckDB を参考にした。DuckDB は https と S3 に SQL を実行できる。同じ機能を1時間で書いた。難しくない。

### 制限

CSV／TSV／LTSV への UPDATE と DELETE はオリジナルに反映されない。これは仕様だ。直しづらい。当面は直さない。オラクルの DB サポートは CGO 依存のため入れない。純粋な Go だけを使う。

### まとめ

仕事で SQL を書く。sqluv を使う。少しずつ直す。スターが増える。開発が進む。以上。
