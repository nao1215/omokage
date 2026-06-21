---
title: "【Golang】バリデーション付き CSV 読み込み機能と DataFrame 機能を持つパッケージを作った話"
type: post
date: 2025-11-15
categories:
  - "linux"
cover:
  image: "images/csv-logo.png"
  alt: "csv package logo"
  hidden: false
---

### 前書き：バリデーション付きCSV読み込み機能は2024年に開発

本記事で取り上げるのは、2024年に開発した [nao1215/csv](https://github.com/nao1215/csv) です。本来の予定では、新規に開発した機能の紹介だけ書く予定でした。しかし、本ブログで一度も nao1215/csv の説明をしていなかったようなので、まずは基本機能を説明した後に新機能（DataFrame）について紹介します。

---


### キッカケ：CSV 読み込みが辛い時期があった

君は、CSV を読み込んだことがあるだろうか。  
数万行、300列以上のCSVファイルから10個のエラー（書き間違い）を見つけるという作業をした経験は？ない？それは幸せな人生だ。私は何度か経験した。こめかみに青筋が浮いた。


さて、話と口調を戻しましょう。CSV 読み込みは、一筋縄で行かないときがあります。例えば、Excel から CSV エクスポートすると文字コードが UTF-8 ではなかったり、改行が含まれているカラムがあったり、仕様通りに記載されていないカラムが登場したりします。

私が特に辛かったのが、どの行で読み込みエラーが発生したか分からなかったことです。私は、Golang 標準ライブラリの csv パッケージや [shogo82148/go-header-csv](https://github.com/shogo82148/go-header-csv) を利用して CSV を読み込むケースが当時多かったです。しかし、これらのライブラリは、エラーが発生した行番号を教えてくれませんでした（注：[go-header-csv は、v0.1.0からエラー発生行番号を教えてくれます](https://shogo82148.github.io/blog/2025/03/27/convert-csv-to-structs-in-golang/)）。

また、CSV カラムに関するバリデーション要件が入ってくると、地獄の門が開門します。一度 CSV を読み込んだ後に、カラム単位でバリデーションしたり、「カラム A が◯◯の場合は、カラム B は☓☓でなければならない」などの条件を地道に実装する必要がでてきます。この要件を含むタスクと2度出会いましたが、実装者は苦戦していました。そう、私は実装していません。管理側の立場だったので、レビューしただけでした。このバリデーション処理はカラムが数百単位であったので、愚直に実装すると、かなり見通しの悪いコードが出来上がります。

---


### [go-playground/validator](https://github.com/go-playground/validator) を参考に [nao1215/csv](https://github.com/nao1215/csv) を開発

自分であれば CSV バリデーションをどのように実装するかを考えた時に、「[go-playground/validator](https://github.com/go-playground/validator) のように、 struct tag を使って CSV バリデーション を定義すれば良いのではないか」と発想しました。この発想を形にしたのが、[nao1215/csv](https://github.com/nao1215/csv) です。

以下に、サンプルコードを示します。  
読み込みたい CSV に対応する構造体（例：`person` 構造体）を用意し、その構造体フィールドに`validate:`タグおよびバリデーション条件を書くと、CSV 読み込み時にエラーを行番号付きで表示してくれます。実装時（2024年）のマイブームだったのか分かりませんが、i18n（国際化）対応されており、英語、日本語、ロシア語でエラーメッセージを表示できます。何故、i18n 対応したのか。理由が自分でも分かりません。

```go
package csv_test

import (
	"bytes"
	"fmt"

	"github.com/nao1215/csv"
)

func ExampleCSV() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf)
	if err != nil {
		panic(err)
	}

	type person struct {
		ID   int    `validate:"numeric"`
		Name string `validate:"alpha"`
		Age  int    `validate:"gt=24"`
	}
	people := make([]person, 0)

	errs := c.Decode(&people)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	// Output:
	// line:2 column age: target is not greater than the threshold value: threshold=24, value=23
	// line:3 column id: target is not a numeric character: value=a
	// line:4 column name: target is not an alphabetic character: value=Den1s
}

func ExampleWithJapaneseLanguage() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf, csv.WithJapaneseLanguage())
	if err != nil {
		panic(err)
	}

	type person struct {
		ID   int    `validate:"numeric"`
		Name string `validate:"alpha"`
		Age  int    `validate:"gt=24"`
	}
	people := make([]person, 0)

	errs := c.Decode(&people)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	// Output:
	// line:2 column age: 値がしきい値より大きくありません: threshold=24, value=23
	// line:3 column id: 値が数字ではありません: value=a
	// line:4 column name: 値がアルファベット文字ではありません: value=Den1s
}
```

---


### 対応しているバリデーションタグ一覧

[nao1215/csv](https://github.com/nao1215/csv)は、[go-playground/validator](https://github.com/go-playground/validator) が対応しているタグを部分的に実装しています。理論的には、[go-playground/validator](https://github.com/go-playground/validator)  の全タグを導入できます。しかし、[nao1215/csv](https://github.com/nao1215/csv) を自分でも使わないので、タグを増やしていません。2025年現在では、優秀な LLM が存在するので、ガッとタグを増やそうと思えば増やせます。

---


#### String rules

| Tag Name     | Description                              |
| ------------ | ---------------------------------------- |
| alpha        | Alphabetic characters only               |
| alphanumeric | Alphanumeric characters                  |
| ascii        | ASCII characters only                    |
| boolean      | Boolean values                           |
| contains     | Contains substring                       |
| containsany  | Contains any of the specified characters |
| lowercase    | Lowercase only                           |
| numeric      | Numeric only                             |
| uppercase    | Uppercase only                           |

---


#### Format rules

| Tag Name | Description         |
| -------- | ------------------- |
| email    | Valid email address |

---


#### Comparison rules

| Tag Name | Description                  |
| -------- | ---------------------------- |
| eq       | Equal to the specified value |
| gt       | Greater than                 |
| gte      | Greater or equal             |
| lt       | Less than                    |
| lte      | Less or equal                |
| ne       | Not equal                    |

---


#### Other rules

| Tag Name | Description                    |
| -------- | ------------------------------ |
| len      | Exact length                   |
| max      | Maximum value                  |
| min      | Minimum value                  |
| oneof    | Must match one of given values |
| required | Must not be empty              |

---


### 本題：v0.3.0で DataFrame 機能を追加

DataFrame 機能は、[pandas](https://pandas.pydata.org/) から着想を得ました。

最近、業務都合で機械学習を勉強しています。機械学習では、Python の [pandas](https://pandas.pydata.org/) を頻繁に使用します。この [pandas](https://pandas.pydata.org/)  は、表形式データを加工して、表示する機能を持ちます。ここでの表形式データは DataFrame と呼ばれ、以下のサンプルコードで示すようにデータ操作機能を提供します。

```csv
# 読み込み対象の data.csv。この行はないものとして考えてください。
name,age,score
Alice,25,70
Bob,32,88
Clara,41,92
Denis,29,60
```

```python
import pandas as pd

df = pd.read_csv("data.csv")                      # 読み込み
df = df[df["age"] >= 30]                          # フィルタ
df = df.sort_values("score", ascending=False)     # ソート
df = df.assign(age_group=df["age"] // 10 * 10)    # 新しい列
print(df[["name", "age", "score", "age_group"]])  # 必要な列だけ表示
```
```plaintext
# 結果
    name  age  score  age_group
2  Clara   41     92         40
1    Bob   32     88         30
```

上記のような処理を見た時に、「csv に対して、シンタックスシュガーで SQL を実行しているだけでは？」と思いました。私は、CSV に対して SQL を実行できる [nao1215/filesql](https://github.com/nao1215/filesql) を開発済みだったので、DataFrame を Golang で実現する土台があると気づきました。Golang で [pandas](https://pandas.pydata.org/) もどきを実現するメリットは正直ありませんが、[nao1215/filesql](https://github.com/nao1215/filesql) を使い倒したかったので（バグ出ししたかったので）、DataFrame を実装することにしました。

以下に、今回機能追加した DataFrame のサンプルコードを示します。
```go
package csv_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nao1215/csv"
)

func ExampleDataFrame_joinFilterSort() {
	users := csv.NewDataFrame(filepath.Join("testdata", "sample.csv")).
		Select("id", "name", "age").
		Mutate("age_bucket", "CASE WHEN age >= 30 THEN '30s' ELSE '20s' END")

	orders := csv.NewDataFrame(filepath.Join("testdata", "orders.csv")).
		Filter("total >= 100").
		Mutate("gross_total", "total + 5")

	depts := csv.NewDataFrame(filepath.Join("testdata", "departments.csv")).
		Select("id", "dept").
		Rename(map[string]string{"dept": "dept_name"})

	df := users.
		Join(orders, "id").
		Join(depts, "id").
		Filter("age >= 23").
		Sort("gross_total", false).
		Select("name", "dept_name", "gross_total", "age_bucket")

	var buf bytes.Buffer
	if err := df.Print(&buf); err != nil {
		panic(err)
	}
	fmt.Print(buf.String())

	// Output:
	// age_bucket  dept_name    gross_total  name
	// 30s         Engineering  155          Denis
	// 20s         Sales        105          Gina
}
```

---


### 誤算：DataFrame 同士を結合する仕様の実現で悩んだ

私は、[pandas](https://pandas.pydata.org/) で一つの DataFrame を使い回す例ばかりを見ていたので、DataFrame 同士が結合できることを知らずに機能追加を始めました。DataFrame 同士を結合できる仕様に気づいた時、私は実装を断念するつもりでした。理由は、素直なやり方で実装できないからです。

具体的に説明しましょう。[nao1215/filesql](https://github.com/nao1215/filesql) は内部的に SQLite3 を利用しており、初期化時に CSV ファイルパスを受け取ると sql.DB 構造体を返す仕様です。普通に実装すると、DataFrame 1個が sql.DB 構造体1個を持つ構成になります。つまり、DataFrame 同士を結合しようとすると、sql.DB 構造体 A から　sql.DB 構造体 B へデータをコピーする処理が発生します。巨大な CSV ファイルを扱う場合、このコピー処理は遅すぎて許容できません。  

VTuber の動画を見ながら導いた結論は、「必要なタイミングだけ、データ読み込みや SQL 実行をする」でした。例えば、DataFrame 初期化やフィルタリング設定をしている段階では、SQL を実行する必要がありません。ユーザーが表の中身を確認したくなった時に初めて、SQL を実行すれば十分です。Lazy Load や Deferred initcalls と発想は同じです。必要になるまで、実行を遅らせます。

この発想のおかげで、シンプルな DataFrame を実装できました。当然、[pandas](https://pandas.pydata.org/) と使い勝手は違いますし、機能が足りていません。しかし、私としては [nao1215/filesql](https://github.com/nao1215/filesql) の新しい使い方を提示できたので満足です。機能拡張は、Issue が飛んできたら行います。

---


### 最後に：CSV や SQL 系のツールばかり作っている

あまり CSV や SQL に詳しくないのですが、以下に示すように OSS がそっち系ばかりになってきました。上2つ（sqly, filesql）は、GitHub Star が100個超えたので、実装テーマが良かったのでしょう。なお、sqluv は開発停止する予定です。Text User Interface を調整するのが辛すぎるので......

- [sqly - eaisly execute SQL against CSV/TSV/LTSV/JSON and Microsoft Excel™ with shell.](https://github.com/nao1215/sqly)
- [filesql - sql driver for CSV, TSV, LTSV, Parquet, Excel with gzip, bzip2, xz, zstd support.](https://github.com/nao1215/filesql)
- [sqluv - simple terminal UI for RDBMS & CSV/TSV/LTSV at local/https/s3](https://github.com/nao1215/sqluv)
- [csv - read csv with validation and simple DataFrame feature in golang](https://github.com/nao1215/csv)
