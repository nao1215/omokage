window.BENCHMARK_DATA = {
  "lastUpdate": 1781189604459,
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
          "id": "3dc8d81dfc959a442a69a86dcefca50db2e138ec",
          "message": "feat(cli): route help <command> to its --help and name missing arguments",
          "timestamp": "2026-06-02T14:58:04+09:00",
          "tree_id": "4112de1826845f1f1af6b9829fbace59cde5788c",
          "url": "https://github.com/nao1215/omokage/commit/3dc8d81dfc959a442a69a86dcefca50db2e138ec"
        },
        "date": 1780379942417,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1200092,
            "unit": "ns/op\t  326635 B/op\t    2820 allocs/op",
            "extra": "1110 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1200092,
            "unit": "ns/op",
            "extra": "1110 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326635,
            "unit": "B/op",
            "extra": "1110 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1110 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 7542962,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "151 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 7542962,
            "unit": "ns/op",
            "extra": "151 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "151 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "151 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 226883,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "5194 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 226883,
            "unit": "ns/op",
            "extra": "5194 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "5194 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "5194 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 203906,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "5768 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 203906,
            "unit": "ns/op",
            "extra": "5768 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "5768 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5768 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 289452,
            "unit": "ns/op\t  213200 B/op\t     921 allocs/op",
            "extra": "4154 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 289452,
            "unit": "ns/op",
            "extra": "4154 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213200,
            "unit": "B/op",
            "extra": "4154 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "4154 times\n2 procs"
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
          "id": "8fc7b50190bd7b7d6350da67bec568a96cd0bfd5",
          "message": "docs(cli): surface the help command in the root help",
          "timestamp": "2026-06-02T15:25:08+09:00",
          "tree_id": "f3512b1f89398b547a364dfd715ca05c26d920a8",
          "url": "https://github.com/nao1215/omokage/commit/8fc7b50190bd7b7d6350da67bec568a96cd0bfd5"
        },
        "date": 1780381561965,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1239906,
            "unit": "ns/op\t  326367 B/op\t    2820 allocs/op",
            "extra": "902 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1239906,
            "unit": "ns/op",
            "extra": "902 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326367,
            "unit": "B/op",
            "extra": "902 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "902 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6515326,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "187 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6515326,
            "unit": "ns/op",
            "extra": "187 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "187 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "187 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 255932,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4706 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 255932,
            "unit": "ns/op",
            "extra": "4706 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4706 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4706 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 233738,
            "unit": "ns/op\t   83016 B/op\t     762 allocs/op",
            "extra": "4863 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 233738,
            "unit": "ns/op",
            "extra": "4863 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83016,
            "unit": "B/op",
            "extra": "4863 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4863 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 320918,
            "unit": "ns/op\t  213207 B/op\t     921 allocs/op",
            "extra": "3656 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 320918,
            "unit": "ns/op",
            "extra": "3656 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213207,
            "unit": "B/op",
            "extra": "3656 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3656 times\n2 procs"
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
          "distinct": false,
          "id": "5c0361306a188ee252f0a5caacbcb1507e237c3c",
          "message": "feat(train): accept multiple file and directory inputs\n\nExtend `omokage train` from a single DIRECTORY to `INPUT...`: one or more\ndirectories and/or .md/.txt files, mixed freely. Inputs are de-duplicated by\nnormalized real path so a file reached through both a directory and a direct\nargument is learned once.\n\nValidate every input up front and fail fast by name: URLs are rejected outright\n(omokage reads local files only and never hits the network), missing paths and\ndirectly-passed unsupported extensions are reported by the exact argument the\nuser typed, and nothing is trained on a partial failure so the user can drop the\noffending input and re-run.\n\nRecord the full provenance: profiles gain a Sources list (stored as a JSON\ncolumn with an ALTER-based migration; older single-source profiles load with\nSources backfilled from SourceDir). `show` prints a numbered \"Sources (N):\"\nblock and a `sources` array in JSON for multiple inputs while keeping the single\n\"Source:\" line for one; `list --long` shows the first source with a \"(+N more)\"\nhint.\n\nUpdate train --help, the root help, and the README, which now documents the\nINPUT spec explicitly (local .md/.txt only, dedup, all-or-nothing, URL\nrejection). Add Go unit tests for input collection/dedup/URL detection and\nstorage round-trip, and ShellSpec E2E coverage for the new behaviors.",
          "timestamp": "2026-06-02T16:00:53+09:00",
          "tree_id": "2491b6c258721ad9e2e25cde456d3adf3e557a03",
          "url": "https://github.com/nao1215/omokage/commit/5c0361306a188ee252f0a5caacbcb1507e237c3c"
        },
        "date": 1780384495412,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1212591,
            "unit": "ns/op\t  326876 B/op\t    2820 allocs/op",
            "extra": "890 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1212591,
            "unit": "ns/op",
            "extra": "890 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326876,
            "unit": "B/op",
            "extra": "890 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "890 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 7843157,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "158 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 7843157,
            "unit": "ns/op",
            "extra": "158 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "158 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "158 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 226418,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "5229 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 226418,
            "unit": "ns/op",
            "extra": "5229 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "5229 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "5229 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 200889,
            "unit": "ns/op\t   83018 B/op\t     762 allocs/op",
            "extra": "5659 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 200889,
            "unit": "ns/op",
            "extra": "5659 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83018,
            "unit": "B/op",
            "extra": "5659 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5659 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 288072,
            "unit": "ns/op\t  213201 B/op\t     921 allocs/op",
            "extra": "4098 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 288072,
            "unit": "ns/op",
            "extra": "4098 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213201,
            "unit": "B/op",
            "extra": "4098 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "4098 times\n2 procs"
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
          "id": "2ace27731d2a3bb499f3edf36e7e168eacd3e87a",
          "message": "fix(train): dedup by real path, keep source_dir a directory, strict help\n\nAddress three review findings on multi-input train:\n\n- Symlink de-duplication: dedup keys now resolve the real path via\n  EvalSymlinks, so an alias and its target (a.md and alias.md -> a.md) are\n  recognized as one file instead of being learned twice and skewing the\n  distribution. Both the input and file levels use the resolved key.\n\n- source_dir semantics: the backward-compatible source_dir field is set only\n  when training from exactly one directory, and is empty for a single file or\n  several inputs, so it never holds a file path where a consumer expects a\n  directory. The full provenance lives in the sources list; `show` text drives\n  its Source line from that list.\n\n- help strictness: `omokage help <command> [args...]` forwards trailing tokens\n  to `<command> ... --help` instead of silently dropping them, so\n  `help check extra` fails exactly as `check extra --help` would and a typo or\n  wrapper bug is not hidden. `help`/`version` reject extra tokens too.\n\nUpdate the README input spec (real-path dedup incl. symlinks, source_dir\ncontract) and add unit + ShellSpec coverage for symlink dedup, the source_dir\ncontract, and help argument strictness.",
          "timestamp": "2026-06-02T16:22:05+09:00",
          "tree_id": "49d6db14a96efe5fc3469847e93f25b9cbaefaee",
          "url": "https://github.com/nao1215/omokage/commit/2ace27731d2a3bb499f3edf36e7e168eacd3e87a"
        },
        "date": 1780384986633,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1325013,
            "unit": "ns/op\t  326784 B/op\t    2820 allocs/op",
            "extra": "963 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1325013,
            "unit": "ns/op",
            "extra": "963 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326784,
            "unit": "B/op",
            "extra": "963 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "963 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 7635407,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "156 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 7635407,
            "unit": "ns/op",
            "extra": "156 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "156 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "156 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 257559,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4072 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 257559,
            "unit": "ns/op",
            "extra": "4072 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4072 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4072 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 240115,
            "unit": "ns/op\t   83016 B/op\t     762 allocs/op",
            "extra": "4860 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 240115,
            "unit": "ns/op",
            "extra": "4860 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83016,
            "unit": "B/op",
            "extra": "4860 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4860 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 319100,
            "unit": "ns/op\t  213200 B/op\t     921 allocs/op",
            "extra": "3626 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 319100,
            "unit": "ns/op",
            "extra": "3626 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213200,
            "unit": "B/op",
            "extra": "3626 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3626 times\n2 procs"
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
          "id": "a52413cd1c2744fcc14a7ff4ab9ceaa81d44514a",
          "message": "docs: tighten train/show README prose, drop heavy formatting\n\nCompress the multi-input train section and the show provenance notes into plain\nprose: remove the bold labels and the bulleted rules wall, keeping every example\nand fact (real-path/symlink dedup, all-or-nothing, URL rejection, source_dir\ncontract). The demo GIFs and the name-origin section are unchanged.",
          "timestamp": "2026-06-02T16:25:21+09:00",
          "tree_id": "611955ef47a036591eab1716b371c580117af2bf",
          "url": "https://github.com/nao1215/omokage/commit/a52413cd1c2744fcc14a7ff4ab9ceaa81d44514a"
        },
        "date": 1780385175590,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1185658,
            "unit": "ns/op\t  326849 B/op\t    2820 allocs/op",
            "extra": "1066 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1185658,
            "unit": "ns/op",
            "extra": "1066 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326849,
            "unit": "B/op",
            "extra": "1066 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1066 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 5969658,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "190 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5969658,
            "unit": "ns/op",
            "extra": "190 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "190 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "190 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 246819,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4738 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 246819,
            "unit": "ns/op",
            "extra": "4738 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4738 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4738 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 227081,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "4944 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 227081,
            "unit": "ns/op",
            "extra": "4944 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "4944 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4944 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 316243,
            "unit": "ns/op\t  213204 B/op\t     921 allocs/op",
            "extra": "3579 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 316243,
            "unit": "ns/op",
            "extra": "3579 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213204,
            "unit": "B/op",
            "extra": "3579 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3579 times\n2 procs"
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
          "id": "0d7c146d3ada386d90fef7d278932efe6e4fee7a",
          "message": "docs: move About the name section to the end",
          "timestamp": "2026-06-02T16:28:05+09:00",
          "tree_id": "1364643f92611f421e3fb164d8cee7cf4bc1c013",
          "url": "https://github.com/nao1215/omokage/commit/0d7c146d3ada386d90fef7d278932efe6e4fee7a"
        },
        "date": 1780385347100,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1281272,
            "unit": "ns/op\t  326689 B/op\t    2820 allocs/op",
            "extra": "1016 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1281272,
            "unit": "ns/op",
            "extra": "1016 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326689,
            "unit": "B/op",
            "extra": "1016 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1016 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 5857086,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "177 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5857086,
            "unit": "ns/op",
            "extra": "177 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "177 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "177 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 255241,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4575 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 255241,
            "unit": "ns/op",
            "extra": "4575 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4575 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4575 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 233867,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "4797 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 233867,
            "unit": "ns/op",
            "extra": "4797 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "4797 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4797 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 320844,
            "unit": "ns/op\t  213201 B/op\t     921 allocs/op",
            "extra": "3529 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 320844,
            "unit": "ns/op",
            "extra": "3529 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213201,
            "unit": "B/op",
            "extra": "3529 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3529 times\n2 procs"
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
          "id": "8d000ed6c84c6c79eba6d5807bb7b606db30a387",
          "message": "docs: refresh stale check --explain output in README\n\nVerified the README examples against the current binary and example corpus.\nThe --explain block was out of date: it showed the drifting paragraph at 11.0σ\n(now 50.0σ) and omitted the \"Low-level fingerprint drift\" section the tool now\nprints. Update the example and prose to match, and add the timezone that\nlist --long actually emits in the TRAINED column.",
          "timestamp": "2026-06-02T16:32:49+09:00",
          "tree_id": "5ce80a8c59f399f2ecf74b3d1e6ae5148e7dea9b",
          "url": "https://github.com/nao1215/omokage/commit/8d000ed6c84c6c79eba6d5807bb7b606db30a387"
        },
        "date": 1780385629722,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1257653,
            "unit": "ns/op\t  326544 B/op\t    2820 allocs/op",
            "extra": "1267 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1257653,
            "unit": "ns/op",
            "extra": "1267 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326544,
            "unit": "B/op",
            "extra": "1267 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1267 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 5284041,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "207 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5284041,
            "unit": "ns/op",
            "extra": "207 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "207 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "207 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 258288,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4561 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 258288,
            "unit": "ns/op",
            "extra": "4561 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4561 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4561 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 235238,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "4929 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 235238,
            "unit": "ns/op",
            "extra": "4929 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "4929 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4929 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 321860,
            "unit": "ns/op\t  213205 B/op\t     921 allocs/op",
            "extra": "3536 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 321860,
            "unit": "ns/op",
            "extra": "3536 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213205,
            "unit": "B/op",
            "extra": "3536 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3536 times\n2 procs"
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
          "id": "733b314850453a289e402590b18c6b85b5fdf85a",
          "message": "docs: make the primary demo English, keep Japanese as the showcase\n\nNon-Japanese visitors landed on a README whose demo GIF and usage examples\nwere entirely Japanese, which is hard to relate to even though the tool is\nbilingual. Lead with English instead, and keep Japanese where it shows the\ntool's distinctive depth.\n\n- Reorganize examples into examples/en and examples/ja (git mv preserves the\n  Japanese corpus history).\n- Add an original English corpus: eight short first-person posts in one casual\n  voice, a keeps-voice draft in that voice, and a stiff, formal lost-voice\n  rewrite of the same idea. Written from scratch, no copied text.\n- Regenerate demo.gif from the English corpus (train -> check keeps 70% ->\n  check stiff rewrite 26%); keep explain.gif on the Japanese corpus, where the\n  register (敬体/常体) and script-ratio features stand out.\n- README: train/check examples now English; diff and --explain stay Japanese\n  with a note explaining why, since English style separates less sharply. All\n  shown numbers verified against the current binary.",
          "timestamp": "2026-06-02T17:03:49+09:00",
          "tree_id": "86501b02da84d69d2ad6abb2e710c03357d66b9a",
          "url": "https://github.com/nao1215/omokage/commit/733b314850453a289e402590b18c6b85b5fdf85a"
        },
        "date": 1780387489637,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1171263,
            "unit": "ns/op\t  326763 B/op\t    2820 allocs/op",
            "extra": "1150 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1171263,
            "unit": "ns/op",
            "extra": "1150 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 326763,
            "unit": "B/op",
            "extra": "1150 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 2820,
            "unit": "allocs/op",
            "extra": "1150 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6114957,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "174 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6114957,
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
            "value": 270155,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4592 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 270155,
            "unit": "ns/op",
            "extra": "4592 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4592 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4592 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 229336,
            "unit": "ns/op\t   83018 B/op\t     762 allocs/op",
            "extra": "4934 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 229336,
            "unit": "ns/op",
            "extra": "4934 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83018,
            "unit": "B/op",
            "extra": "4934 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4934 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 316236,
            "unit": "ns/op\t  213205 B/op\t     921 allocs/op",
            "extra": "3366 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 316236,
            "unit": "ns/op",
            "extra": "3366 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213205,
            "unit": "B/op",
            "extra": "3366 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3366 times\n2 procs"
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
          "id": "59a833d20403abfb179216e0489138b24978f453",
          "message": "fix: measure every style feature on code-stripped prose and harden CLI edges\n\nAlign the implementation with the documented behavior and fix robustness gaps surfaced in review:\n\n- feature: strip code before measuring all features, not just lexical/n-gram, so a fenced block no longer manufactures false drift (matches the README promise)\n- feature: split sentences only on '.' at a boundary, keeping version numbers, domains, and decimals intact\n- feature: classify Japanese sentence endings by predicate so the plain register (常体) is detected for verbs, i-adjectives, ない, and だ, not only である/だった\n- cmd: let 'diff --global' fall back to default weights when no global store exists, matching bare diff\n- cmd: warn (but still proceed) when 'init' nests inside an existing store\n- cmd: render trained_at in local time for show and list --long\n\nAdd feature-level and end-to-end regression tests for each.",
          "timestamp": "2026-06-02T17:47:50+09:00",
          "tree_id": "d9e9d1de0d2bafc2847a3da55b097a19bbf1d2e6",
          "url": "https://github.com/nao1215/omokage/commit/59a833d20403abfb179216e0489138b24978f453"
        },
        "date": 1780390133968,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1009066,
            "unit": "ns/op\t  351039 B/op\t    3066 allocs/op",
            "extra": "1122 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1009066,
            "unit": "ns/op",
            "extra": "1122 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 351039,
            "unit": "B/op",
            "extra": "1122 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "1122 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3721518,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "303 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3721518,
            "unit": "ns/op",
            "extra": "303 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "303 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "303 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 191436,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "6082 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 191436,
            "unit": "ns/op",
            "extra": "6082 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "6082 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "6082 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 175792,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "6506 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 175792,
            "unit": "ns/op",
            "extra": "6506 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "6506 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "6506 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 245546,
            "unit": "ns/op\t  213207 B/op\t     921 allocs/op",
            "extra": "4698 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 245546,
            "unit": "ns/op",
            "extra": "4698 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213207,
            "unit": "B/op",
            "extra": "4698 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "4698 times\n2 procs"
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
          "id": "71b2403c6d1cc581f6adf34851bce87f18e6b705",
          "message": "fix(feature): strip tilde code fences too, not only backtick fences\n\nCommonMark allows ~~~ … ~~~ fences alongside , but stripCode only\nrecognized backtick fences, so tilde-fenced code leaked into the features (a\nself-similarity drop to 88% and a stray '~~' n-gram). Recognize both markers and\nclose a block only by its own marker. Add an end-to-end diff test and a\nfeature-level parity test.",
          "timestamp": "2026-06-02T17:53:47+09:00",
          "tree_id": "ebb6fce069cecb7794adf3bf9c8c167cc6266ce7",
          "url": "https://github.com/nao1215/omokage/commit/71b2403c6d1cc581f6adf34851bce87f18e6b705"
        },
        "date": 1780390478075,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1290743,
            "unit": "ns/op\t  350784 B/op\t    3066 allocs/op",
            "extra": "1009 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1290743,
            "unit": "ns/op",
            "extra": "1009 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 350784,
            "unit": "B/op",
            "extra": "1009 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "1009 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 6426799,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "172 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 6426799,
            "unit": "ns/op",
            "extra": "172 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "172 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "172 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 270456,
            "unit": "ns/op\t   86033 B/op\t     864 allocs/op",
            "extra": "4080 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 270456,
            "unit": "ns/op",
            "extra": "4080 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86033,
            "unit": "B/op",
            "extra": "4080 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4080 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 226850,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "4993 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 226850,
            "unit": "ns/op",
            "extra": "4993 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "4993 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4993 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 324832,
            "unit": "ns/op\t  213207 B/op\t     921 allocs/op",
            "extra": "3564 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 324832,
            "unit": "ns/op",
            "extra": "3564 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213207,
            "unit": "B/op",
            "extra": "3564 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3564 times\n2 procs"
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
          "id": "a016e01ee53c2e01b52c1f6ead4e5592a4d18114",
          "message": "fix(feature): strip tilde code fences too, not only backtick fences\n\nCommonMark allows tilde fences (~~~ ... ~~~) alongside backtick fences, but\nstripCode only recognized backtick fences, so tilde-fenced code leaked into the\nfeatures (a self-similarity drop to 88% and a stray \"~~\" n-gram). Recognize both\nmarkers and close a block only by its own marker. Add an end-to-end diff test\nand a feature-level parity test.",
          "timestamp": "2026-06-02T17:54:15+09:00",
          "tree_id": "ebb6fce069cecb7794adf3bf9c8c167cc6266ce7",
          "url": "https://github.com/nao1215/omokage/commit/a016e01ee53c2e01b52c1f6ead4e5592a4d18114"
        },
        "date": 1780390517108,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 1453730,
            "unit": "ns/op\t  351005 B/op\t    3066 allocs/op",
            "extra": "861 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 1453730,
            "unit": "ns/op",
            "extra": "861 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 351005,
            "unit": "B/op",
            "extra": "861 times\n2 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "861 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 5701508,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "210 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5701508,
            "unit": "ns/op",
            "extra": "210 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "210 times\n2 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "210 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 253735,
            "unit": "ns/op\t   86032 B/op\t     864 allocs/op",
            "extra": "4582 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 253735,
            "unit": "ns/op",
            "extra": "4582 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86032,
            "unit": "B/op",
            "extra": "4582 times\n2 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4582 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 232124,
            "unit": "ns/op\t   83017 B/op\t     762 allocs/op",
            "extra": "4916 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 232124,
            "unit": "ns/op",
            "extra": "4916 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83017,
            "unit": "B/op",
            "extra": "4916 times\n2 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4916 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 320341,
            "unit": "ns/op\t  213198 B/op\t     921 allocs/op",
            "extra": "3614 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 320341,
            "unit": "ns/op",
            "extra": "3614 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213198,
            "unit": "B/op",
            "extra": "3614 times\n2 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3614 times\n2 procs"
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
          "id": "345d016ae2bcb3dd7d63399e3bacbe101db71859",
          "message": "docs: cut the 0.1.0 release in the changelog (#2)",
          "timestamp": "2026-06-02T19:15:57+09:00",
          "tree_id": "8b422baccabcf82e4eb6f7b5bd3d5657d0624630",
          "url": "https://github.com/nao1215/omokage/commit/345d016ae2bcb3dd7d63399e3bacbe101db71859"
        },
        "date": 1780395403657,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 847229,
            "unit": "ns/op\t  352756 B/op\t    3066 allocs/op",
            "extra": "1810 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 847229,
            "unit": "ns/op",
            "extra": "1810 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 352756,
            "unit": "B/op",
            "extra": "1810 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "1810 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3721920,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3721920,
            "unit": "ns/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 258474,
            "unit": "ns/op\t   86041 B/op\t     864 allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 258474,
            "unit": "ns/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86041,
            "unit": "B/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 236697,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 236697,
            "unit": "ns/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 323057,
            "unit": "ns/op\t  213221 B/op\t     921 allocs/op",
            "extra": "3672 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 323057,
            "unit": "ns/op",
            "extra": "3672 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213221,
            "unit": "B/op",
            "extra": "3672 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3672 times\n4 procs"
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
          "id": "68f6ed0cf3a760bafcd233175ef96ea53fd5e6a2",
          "message": "fix(ci): grant contents:write so GoReleaser can publish the release (#3)",
          "timestamp": "2026-06-02T19:20:40+09:00",
          "tree_id": "2c21373dd232291b8f1cd6d84ef4cacbe71601fc",
          "url": "https://github.com/nao1215/omokage/commit/68f6ed0cf3a760bafcd233175ef96ea53fd5e6a2"
        },
        "date": 1780395686113,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 869263,
            "unit": "ns/op\t  352851 B/op\t    3066 allocs/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 869263,
            "unit": "ns/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 352851,
            "unit": "B/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 4223795,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 4223795,
            "unit": "ns/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 263978,
            "unit": "ns/op\t   86042 B/op\t     864 allocs/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 263978,
            "unit": "ns/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86042,
            "unit": "B/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 231279,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "5007 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 231279,
            "unit": "ns/op",
            "extra": "5007 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "5007 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5007 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 314042,
            "unit": "ns/op\t  213222 B/op\t     921 allocs/op",
            "extra": "3699 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 314042,
            "unit": "ns/op",
            "extra": "3699 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213222,
            "unit": "B/op",
            "extra": "3699 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3699 times\n4 procs"
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
          "id": "10ebb9e0997a5685670c36d37cdf9f6b0b9d47c5",
          "message": "feat(term): learn and report per-profile notation preferences without an LLM (#4)\n\n* feat(term): learn and report per-profile notation preferences without an LLM\n\nExtract a learning corpus's preferred surface forms (DB vs データベース) during\ntrain and store them profile-locally in SQLite, with no LLM, network, or\nexternal dictionary, deterministically.\n\n- internal/term: deterministic extraction. normalized_key folds case/full-width\n  ASCII/edge punctuation; group_key merges concepts only via corpus-declared\n  alias bridges (Japanese phrase ↔ uppercase acronym, or \"以下、X\"). A shared\n  acronym glossing several phrases is dropped rather than chaining them.\n  Preferred surface order: doc_count, then count, then ascending surface.\n- storage: term_group/term_variant tables added to the schema (auto-created for\n  older profiles), with SaveTerms/LoadTerms in a transaction.\n- show --format json gains term_preferences; check --format json gains\n  term_warnings (a separate layer that never affects the similarity score).\n  occurrences is reserved for future line/column data.\n- Plain check output and existing similarity scoring are unchanged.\n\n* fix(term): degrade show --format json when term load fails; condense README\n\n- runShow JSON path no longer exits 1 on a LoadTerms error; it warns on stderr\n  and emits the summary with an empty term list, matching the comment and the\n  check path (addresses CodeRabbit review).\n- Trim the README term-preferences section: plain prose, no bold or extra\n  decoration, drop the verbose JSON samples.\n\n* test(term): cover corpus IO, empty/noise inputs, union determinism, and storage errors\n\nClose behavioral gaps: ExtractCorpus file reading and read errors, empty/\nwhitespace/code-only documents yielding no groups, stripNoise excluding\nURL/front-matter/HTML tokens, group_key independence from union order, CheckText\non an empty profile, keepASCII filtering, and LoadTerms on a corrupt file.\ninternal/term coverage 92.9% -> 97.8%.",
          "timestamp": "2026-06-02T23:34:39+09:00",
          "tree_id": "aea07aa5845c26af1d6e46ff8d36484a89505043",
          "url": "https://github.com/nao1215/omokage/commit/10ebb9e0997a5685670c36d37cdf9f6b0b9d47c5"
        },
        "date": 1780410929667,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 864837,
            "unit": "ns/op\t  352087 B/op\t    3066 allocs/op",
            "extra": "1532 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 864837,
            "unit": "ns/op",
            "extra": "1532 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 352087,
            "unit": "B/op",
            "extra": "1532 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3066,
            "unit": "allocs/op",
            "extra": "1532 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3878208,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "289 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3878208,
            "unit": "ns/op",
            "extra": "289 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "289 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "289 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 260272,
            "unit": "ns/op\t   86042 B/op\t     864 allocs/op",
            "extra": "4608 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 260272,
            "unit": "ns/op",
            "extra": "4608 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86042,
            "unit": "B/op",
            "extra": "4608 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4608 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 236966,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 236966,
            "unit": "ns/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 320767,
            "unit": "ns/op\t  213220 B/op\t     921 allocs/op",
            "extra": "3414 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 320767,
            "unit": "ns/op",
            "extra": "3414 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213220,
            "unit": "B/op",
            "extra": "3414 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3414 times\n4 procs"
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
          "id": "14fbe86f58dafe6011bedfa22fe0dd6c6173e4bf",
          "message": "fix(feature): measure prose only — strip HTML/front matter/links and fix segment fence context (#6)\n\n* fix(feature): measure prose only — strip HTML/front matter/links and fix segment fence context (#5)\n\nTwo non-prose leaks made check --explain/--format json report drift on content\nthe author cannot edit (mermaid diagrams, HTML blocks, CLI transcripts, image and\nlink URLs, YAML front matter):\n\n- ExtractSegments split the raw document into paragraphs before stripping code, so\n  a fenced block containing a blank line lost its fence context and its interior\n  was measured as prose. It now strips on the whole document first, then splits.\n- Non-prose was only partly removed. Introduce feature.StripNonProse — the single\n  cleaner for front matter, fenced/inline code, Markdown images, link URLs (visible\n  text kept), raw URLs, HTML tags, and entities — and run both whole-document and\n  per-paragraph extraction through it. The term package now reuses it instead of\n  its own stripNoise (removed), so the two stay in sync.\n\nAdds regression tests (HTML stripping, fenced block with a blank line, HTML-only\nparagraph, StripNonProse). BenchmarkExtractText ~+10%, well under the gate.\n\n* test(e2e): add shellspec coverage for term preferences and prose-only extraction\n\nDrive the built binary end to end:\n- show --format json reports the bridged DB/データベース group with separate\n  normalized_key and group_key; show text stays a provenance summary.\n- check --format json flags a non-preferred surface (ＤＢ/データベース) and emits\n  an empty term_warnings array when the preferred surface is used; plain check is\n  unchanged and silent on stderr.\n- check --explain/--format json never reports a fenced mermaid diagram, an HTML\n  block, YAML front matter, or a link URL as a drifting paragraph (Issue #5).",
          "timestamp": "2026-06-02T23:48:56+09:00",
          "tree_id": "f57c07b387ec4f95ad55fddeea26e39f2e8e1c34",
          "url": "https://github.com/nao1215/omokage/commit/14fbe86f58dafe6011bedfa22fe0dd6c6173e4bf"
        },
        "date": 1780411790813,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 948800,
            "unit": "ns/op\t  389645 B/op\t    3084 allocs/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 948800,
            "unit": "ns/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 389645,
            "unit": "B/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3084,
            "unit": "allocs/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3721909,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3721909,
            "unit": "ns/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 256734,
            "unit": "ns/op\t   86043 B/op\t     864 allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 256734,
            "unit": "ns/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86043,
            "unit": "B/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 240331,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 240331,
            "unit": "ns/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 315635,
            "unit": "ns/op\t  213222 B/op\t     921 allocs/op",
            "extra": "3654 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 315635,
            "unit": "ns/op",
            "extra": "3654 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213222,
            "unit": "B/op",
            "extra": "3654 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3654 times\n4 procs"
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
          "id": "d457cdc0daf24d884f040ae7d7100e8da9de86cb",
          "message": "docs: cut the 0.2.0 release in the changelog (#7)",
          "timestamp": "2026-06-02T23:52:08+09:00",
          "tree_id": "93cec7e5726093d21779452dd103eab49c5c122d",
          "url": "https://github.com/nao1215/omokage/commit/d457cdc0daf24d884f040ae7d7100e8da9de86cb"
        },
        "date": 1780411975070,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 900235,
            "unit": "ns/op\t  389670 B/op\t    3084 allocs/op",
            "extra": "1669 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 900235,
            "unit": "ns/op",
            "extra": "1669 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 389670,
            "unit": "B/op",
            "extra": "1669 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3084,
            "unit": "allocs/op",
            "extra": "1669 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3972646,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "286 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3972646,
            "unit": "ns/op",
            "extra": "286 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "286 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "286 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 259555,
            "unit": "ns/op\t   86042 B/op\t     864 allocs/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 259555,
            "unit": "ns/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86042,
            "unit": "B/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 243516,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 243516,
            "unit": "ns/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 325426,
            "unit": "ns/op\t  213222 B/op\t     921 allocs/op",
            "extra": "3543 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 325426,
            "unit": "ns/op",
            "extra": "3543 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213222,
            "unit": "B/op",
            "extra": "3543 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3543 times\n4 procs"
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
          "id": "be58fe843d200389763f706b82539db5d163329f",
          "message": "feat(quality): add doctor command and corpus-quality guidance (#8)\n\n* feat(quality): add doctor command and corpus-quality guidance\n\nAdd a corpus-quality layer so users can tell whether their training\nmaterial is solid before they trust the scores. A new doctor command reads\na corpus (training and writing nothing) and rates the reliability of\ncomparisons against it, naming a next action for too few documents, short\ndocuments, a mixed voice, or an out-of-place document; --format json emits\nthe same report. The same assessment surfaces as a terminal-only note after\ntrain and as reliability/quality_findings in show --format json.\n\nStylometric judgements stay in internal/profile (DocumentDivergence,\nLeaveOneOutDivergences, HighLevelSpreads, reusing the existing localizable\nfeature set and std floor); internal/quality holds only thresholds and\nwording. Outliers are measured leave-one-out so a single odd file is caught\neven on a small corpus. Help and README were rewritten for scannability and\nhonesty: output modes, --author as a purpose-named profile, corpus caveats,\nand an LLM revise-and-recheck loop. Adds unit and shellspec coverage.\n\n* test(quality): avoid a floating-point boundary in the mixed-voice test\n\nA 50/50 register split puts the relative spread exactly on the 1.0 warning\nthreshold, which rounded to a notice on macOS and a warning on Linux. Use an\nuneven 5/3 split so the spread lands at ~1.29, well clear of the boundary, so\nthe test is deterministic across platforms.\n\n* docs(quality): add doc comments to quality helpers and tests\n\nDocument the per-check helpers, the doctor command entry points, and every\nnew test function so their intent is clear from the declaration, addressing\nCodeRabbit's docstring-coverage check.\n\n* feat(quality): keep findings through show, print reliability on train, denoise segments\n\nAddress UX gaps found in real use:\n\n- train now prints the corpus reliability on stdout for everyone (a person, a\n  script, an LLM), not only at a terminal, matching train --help. A thin or mixed\n  corpus lists the fixes and points at doctor; a clean corpus prints one line.\n- The quality findings are stored in the profile at train time, so show --format\n  json reports what doctor found — including the per-document outlier and\n  short-file findings the stored distribution alone could not reproduce. show text\n  gains a one-line Reliability field, and show --format json --summary omits the\n  large term_preferences list for a lighter LLM-facing payload (default unchanged).\n- check --explain/--format json localizes drift to running-prose paragraphs only;\n  headings, bullet and table blocks, and blockquotes are no longer reported as\n  drifting paragraphs (the whole-document score still measures them).\n- The README is more than half shorter and leads with what omokage does and does\n  not do; the root help points at doctor.\n\nFindings are persisted as opaque JSON bytes so storage stays decoupled from the\nquality package (which imports profile). Adds unit and shellspec coverage for each\nbehavior. Regenerates demo.gif for the new train output.\n\n* docs(storage): document the quality-findings storage tests",
          "timestamp": "2026-06-03T22:26:55+09:00",
          "tree_id": "6f6bcfe71dc1b5293897446103748c2dde04e53a",
          "url": "https://github.com/nao1215/omokage/commit/be58fe843d200389763f706b82539db5d163329f"
        },
        "date": 1780493265363,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 858083,
            "unit": "ns/op\t  389241 B/op\t    3084 allocs/op",
            "extra": "1923 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 858083,
            "unit": "ns/op",
            "extra": "1923 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 389241,
            "unit": "B/op",
            "extra": "1923 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3084,
            "unit": "allocs/op",
            "extra": "1923 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3867261,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "307 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3867261,
            "unit": "ns/op",
            "extra": "307 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "307 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "307 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 255202,
            "unit": "ns/op\t   86042 B/op\t     864 allocs/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 255202,
            "unit": "ns/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86042,
            "unit": "B/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 231815,
            "unit": "ns/op\t   83026 B/op\t     762 allocs/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 231815,
            "unit": "ns/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83026,
            "unit": "B/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 313707,
            "unit": "ns/op\t  213221 B/op\t     921 allocs/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 313707,
            "unit": "ns/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213221,
            "unit": "B/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "3658 times\n4 procs"
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
          "id": "9fe4dedaac5fb728366373a5e6bc3a5b5c6d42ae",
          "message": "docs: cut the 0.3.0 release in the changelog (#9)",
          "timestamp": "2026-06-03T22:35:15+09:00",
          "tree_id": "92e68ac515885d540478ed3a5dfd2ab99b758218",
          "url": "https://github.com/nao1215/omokage/commit/9fe4dedaac5fb728366373a5e6bc3a5b5c6d42ae"
        },
        "date": 1780493758875,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 918230,
            "unit": "ns/op\t  390191 B/op\t    3084 allocs/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 918230,
            "unit": "ns/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 390191,
            "unit": "B/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 3084,
            "unit": "allocs/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 3294296,
            "unit": "ns/op\t  102944 B/op\t      34 allocs/op",
            "extra": "348 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 3294296,
            "unit": "ns/op",
            "extra": "348 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 102944,
            "unit": "B/op",
            "extra": "348 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "348 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 224288,
            "unit": "ns/op\t   86043 B/op\t     864 allocs/op",
            "extra": "5120 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 224288,
            "unit": "ns/op",
            "extra": "5120 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 86043,
            "unit": "B/op",
            "extra": "5120 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 864,
            "unit": "allocs/op",
            "extra": "5120 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 200060,
            "unit": "ns/op\t   83025 B/op\t     762 allocs/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 200060,
            "unit": "ns/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 83025,
            "unit": "B/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 762,
            "unit": "allocs/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 277365,
            "unit": "ns/op\t  213221 B/op\t     921 allocs/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 277365,
            "unit": "ns/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 213221,
            "unit": "B/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 921,
            "unit": "allocs/op",
            "extra": "4077 times\n4 procs"
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
          "id": "0498985675265c0f1f3a77d08f95919e3318c5ad",
          "message": "feat(feature): use kagome morphological analysis for Japanese scoring (#10)\n\nJapanese has no whitespace, so the language-neutral tokenizer mismeasured several features. Use kagome (IPA dict) to count function words as whole morphemes (no で-in-です double counting), classify polite/plain register from each sentence's closing predicate, and give conjunction frequency a real morpheme denominator; a held-out author-attribution test improved over the heuristics. Add two opt-in Japanese features (POS n-gram fingerprint, lemma-based vocabulary richness) off by default since the test showed no gain, with additive storage migration. Record a feature-definition version per profile and warn on check when it differs from the build. English prose is unchanged.",
          "timestamp": "2026-06-11T23:49:32+09:00",
          "tree_id": "d020983d30954d1214f305d2760086ea831ba82c",
          "url": "https://github.com/nao1215/omokage/commit/0498985675265c0f1f3a77d08f95919e3318c5ad"
        },
        "date": 1781189430142,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 7839531,
            "unit": "ns/op\t 3202538 B/op\t   52560 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 7839531,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 3202538,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 52560,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 4369379,
            "unit": "ns/op\t  110864 B/op\t      52 allocs/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 4369379,
            "unit": "ns/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 110864,
            "unit": "B/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 52,
            "unit": "allocs/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 248570,
            "unit": "ns/op\t   85960 B/op\t     848 allocs/op",
            "extra": "5346 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 248570,
            "unit": "ns/op",
            "extra": "5346 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 85960,
            "unit": "B/op",
            "extra": "5346 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 848,
            "unit": "allocs/op",
            "extra": "5346 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 217554,
            "unit": "ns/op\t   82688 B/op\t     746 allocs/op",
            "extra": "6234 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 217554,
            "unit": "ns/op",
            "extra": "6234 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 82688,
            "unit": "B/op",
            "extra": "6234 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 746,
            "unit": "allocs/op",
            "extra": "6234 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 327640,
            "unit": "ns/op\t  215162 B/op\t     905 allocs/op",
            "extra": "4250 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 327640,
            "unit": "ns/op",
            "extra": "4250 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 215162,
            "unit": "B/op",
            "extra": "4250 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 905,
            "unit": "allocs/op",
            "extra": "4250 times\n4 procs"
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
          "id": "f6693d52d66b436f0b130fb460738797b1778f1f",
          "message": "docs: cut the 0.4.0 release in the changelog (#11)",
          "timestamp": "2026-06-11T23:52:38+09:00",
          "tree_id": "edd0cd7685928433af90d469244bd36743f856e3",
          "url": "https://github.com/nao1215/omokage/commit/f6693d52d66b436f0b130fb460738797b1778f1f"
        },
        "date": 1781189603948,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature)",
            "value": 5916068,
            "unit": "ns/op\t 2557590 B/op\t   39977 allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 5916068,
            "unit": "ns/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 2557590,
            "unit": "B/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkExtractText (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 39977,
            "unit": "allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature)",
            "value": 4133264,
            "unit": "ns/op\t  110864 B/op\t      52 allocs/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - ns/op",
            "value": 4133264,
            "unit": "ns/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - B/op",
            "value": 110864,
            "unit": "B/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkAggregate (github.com/nao1215/omokage/internal/feature) - allocs/op",
            "value": 52,
            "unit": "allocs/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile)",
            "value": 272342,
            "unit": "ns/op\t   85961 B/op\t     848 allocs/op",
            "extra": "4806 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 272342,
            "unit": "ns/op",
            "extra": "4806 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 85961,
            "unit": "B/op",
            "extra": "4806 times\n4 procs"
          },
          {
            "name": "BenchmarkScore (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 848,
            "unit": "allocs/op",
            "extra": "4806 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile)",
            "value": 234866,
            "unit": "ns/op\t   82688 B/op\t     746 allocs/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 234866,
            "unit": "ns/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 82688,
            "unit": "B/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkCompare (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 746,
            "unit": "allocs/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile)",
            "value": 332673,
            "unit": "ns/op\t  215162 B/op\t     905 allocs/op",
            "extra": "3886 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - ns/op",
            "value": 332673,
            "unit": "ns/op",
            "extra": "3886 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - B/op",
            "value": 215162,
            "unit": "B/op",
            "extra": "3886 times\n4 procs"
          },
          {
            "name": "BenchmarkExplain (github.com/nao1215/omokage/internal/profile) - allocs/op",
            "value": 905,
            "unit": "allocs/op",
            "extra": "3886 times\n4 procs"
          }
        ]
      }
    ]
  }
}