window.BENCHMARK_DATA = {
  "lastUpdate": 1780377617476,
  "repoUrl": "https://github.com/nao1215/omokage",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "committer": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "distinct": true,
          "id": "52ba12f2a5bb354d86b7d0ea618513e280a5a78f",
          "message": "ci: add benchmarks with a regression gate; install shellspec via make tools and document e2e/bench",
          "timestamp": "2026-06-02T10:51:03+09:00",
          "tree_id": "197ff57dc2c323292a8adc62c6a9056686486a17",
          "url": "https://github.com/nao1215/omokage/commit/52ba12f2a5bb354d86b7d0ea618513e280a5a78f"
        },
        "date": 1780365291545,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1235781,
            "unit": "ns/op\t  326474 B/op\t    2820 allocs/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1235781,
            "unit": "ns/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326474,
            "unit": "B/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6214541,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "169 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6214541,
            "unit": "ns/op",
            "extra": "169 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "169 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "169 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 245458,
            "unit": "ns/op\t   57998 B/op\t     868 allocs/op",
            "extra": "4843 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 245458,
            "unit": "ns/op",
            "extra": "4843 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 57998,
            "unit": "B/op",
            "extra": "4843 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 868,
            "unit": "allocs/op",
            "extra": "4843 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 237325,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "5019 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 237325,
            "unit": "ns/op",
            "extra": "5019 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "5019 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5019 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "committer": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "distinct": true,
          "id": "059250d5188db98b183864c1d8036a4faefebe99",
          "message": "ci: add GoReleaser release workflow and security policy\n\nRelease on tag push (v*) via GoReleaser: cross-platform archives, checksums,\nand Linux packages (deb/rpm/apk). Version is injected through ldflags so the\nreleased binary reports the tag. Add SECURITY.md describing how to report\nvulnerabilities.",
          "timestamp": "2026-06-02T11:41:10+09:00",
          "tree_id": "4e8daecc2cbc7db06ec373ab42d14964b5ee2870",
          "url": "https://github.com/nao1215/omokage/commit/059250d5188db98b183864c1d8036a4faefebe99"
        },
        "date": 1780368138234,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1162863,
            "unit": "ns/op\t  326533 B/op\t    2820 allocs/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1162863,
            "unit": "ns/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326533,
            "unit": "B/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6165346,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "174 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6165346,
            "unit": "ns/op",
            "extra": "174 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "174 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "174 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 270448,
            "unit": "ns/op\t  164888 B/op\t     869 allocs/op",
            "extra": "4412 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 270448,
            "unit": "ns/op",
            "extra": "4412 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 164888,
            "unit": "B/op",
            "extra": "4412 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 869,
            "unit": "allocs/op",
            "extra": "4412 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 250624,
            "unit": "ns/op\t   83019 B/op\t     762 allocs/op",
            "extra": "4923 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 250624,
            "unit": "ns/op",
            "extra": "4923 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83019,
            "unit": "B/op",
            "extra": "4923 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4923 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 1424726,
            "unit": "ns/op\t 1154762 B/op\t    7830 allocs/op",
            "extra": "826 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 1424726,
            "unit": "ns/op",
            "extra": "826 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 1154762,
            "unit": "B/op",
            "extra": "826 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 7830,
            "unit": "allocs/op",
            "extra": "826 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "committer": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "distinct": true,
          "id": "0a42020d7b8094d9625a18ba97df1aadcb6f73ff",
          "message": "docs: add `make demo` and an --explain demo GIF\n\nAdd a `make demo` target that builds the CLI, seeds a throwaway project, runs\nthe VHS tapes, and writes doc/img/demo.gif and doc/img/explain.gif. Add a second\ntape and GIF showing the late-stage tuning view (`check --explain`) and embed it\nin the README next to the explain example. Refresh the overview GIF.",
          "timestamp": "2026-06-02T12:55:14+09:00",
          "tree_id": "aefbbcf51281f3c8230d5ff09cfb3349a0c854eb",
          "url": "https://github.com/nao1215/omokage/commit/0a42020d7b8094d9625a18ba97df1aadcb6f73ff"
        },
        "date": 1780372586626,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1211819,
            "unit": "ns/op\t  326224 B/op\t    2820 allocs/op",
            "extra": "934 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1211819,
            "unit": "ns/op",
            "extra": "934 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326224,
            "unit": "B/op",
            "extra": "934 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "934 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 7050211,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "165 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 7050211,
            "unit": "ns/op",
            "extra": "165 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "165 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "165 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 269878,
            "unit": "ns/op\t  164890 B/op\t     869 allocs/op",
            "extra": "4419 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 269878,
            "unit": "ns/op",
            "extra": "4419 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 164890,
            "unit": "B/op",
            "extra": "4419 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 869,
            "unit": "allocs/op",
            "extra": "4419 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 234417,
            "unit": "ns/op\t   83018 B/op\t     762 allocs/op",
            "extra": "4862 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 234417,
            "unit": "ns/op",
            "extra": "4862 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83018,
            "unit": "B/op",
            "extra": "4862 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4862 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 329181,
            "unit": "ns/op\t  242905 B/op\t     925 allocs/op",
            "extra": "3441 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 329181,
            "unit": "ns/op",
            "extra": "3441 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 242905,
            "unit": "B/op",
            "extra": "3441 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 925,
            "unit": "allocs/op",
            "extra": "3441 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "committer": {
            "email": "n.chika156@gmail.com",
            "name": "Naohiro CHIKAMATSU",
            "username": "nao1215"
          },
          "distinct": true,
          "id": "8a8fa973faeaca7060999c0d7827f2dbf93e6815",
          "message": "perf(profile): cut Score allocations back under the benchmark gate\n\nThe shared featureDrifts core regressed Score's memory ~2.8x: it grew the drift\nslice from a small capacity while appending hundreds of function-word and n-gram\nentries, and topDifferences then copied the whole slice before sorting.\n\n- Preallocate featureDrifts to the maximum drift count so the fingerprint\n  appends no longer reallocate and copy the backing array.\n- Sort in place in topDifferences instead of copying; Score is the only caller\n  and discards the slice afterwards.\n\nBenchmarkScore drops from ~165kB to ~86kB per op (back within the 2x regression\ngate) with identical output; Explain benefits from the same preallocation.",
          "timestamp": "2026-06-02T12:59:06+09:00",
          "tree_id": "244123ecb00d27b922ff6d08951c77497a7f2f86",
          "url": "https://github.com/nao1215/omokage/commit/8a8fa973faeaca7060999c0d7827f2dbf93e6815"
        },
        "date": 1780372798706,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1113833,
            "unit": "ns/op\t  326686 B/op\t    2820 allocs/op",
            "extra": "1344 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1113833,
            "unit": "ns/op",
            "extra": "1344 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326686,
            "unit": "B/op",
            "extra": "1344 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1344 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 5524912,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "202 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5524912,
            "unit": "ns/op",
            "extra": "202 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "202 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "202 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 253513,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4604 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 253513,
            "unit": "ns/op",
            "extra": "4604 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4604 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4604 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 226129,
            "unit": "ns/op\t   83018 B/op\t     762 allocs/op",
            "extra": "5121 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 226129,
            "unit": "ns/op",
            "extra": "5121 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83018,
            "unit": "B/op",
            "extra": "5121 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5121 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 316638,
            "unit": "ns/op\t  213200 B/op\t     921 allocs/op",
            "extra": "3506 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 316638,
            "unit": "ns/op",
            "extra": "3506 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213200,
            "unit": "B/op",
            "extra": "3506 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3506 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "n.chika156@gmail.com",
            "name": "CHIKAMATSU Naohiro",
            "username": "nao1215"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b3a1093f73f461cf52dbbaea40fe312e3741d076",
          "message": "Merge pull request #1 from nao1215/feat/cli-ergonomics\n\nfeat(cli): improve everyday ergonomics for single, multi, local, and global use",
          "timestamp": "2026-06-02T14:19:22+09:00",
          "tree_id": "d2940791b9040986865e6151295f1b2000b55d84",
          "url": "https://github.com/nao1215/omokage/commit/b3a1093f73f461cf52dbbaea40fe312e3741d076"
        },
        "date": 1780377616312,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1179850,
            "unit": "ns/op\t  326655 B/op\t    2820 allocs/op",
            "extra": "1088 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1179850,
            "unit": "ns/op",
            "extra": "1088 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326655,
            "unit": "B/op",
            "extra": "1088 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1088 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6920090,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "160 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6920090,
            "unit": "ns/op",
            "extra": "160 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "160 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "160 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 234778,
            "unit": "ns/op\t   86033 B/op\t     864 allocs/op",
            "extra": "4884 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 234778,
            "unit": "ns/op",
            "extra": "4884 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86033,
            "unit": "B/op",
            "extra": "4884 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4884 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 200690,
            "unit": "ns/op\t   83018 B/op\t     762 allocs/op",
            "extra": "5703 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 200690,
            "unit": "ns/op",
            "extra": "5703 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83018,
            "unit": "B/op",
            "extra": "5703 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5703 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 280826,
            "unit": "ns/op\t  213199 B/op\t     921 allocs/op",
            "extra": "3952 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 280826,
            "unit": "ns/op",
            "extra": "3952 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213199,
            "unit": "B/op",
            "extra": "3952 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3952 times\n2 procs"
          }
        ]
      }
    ]
  }
}