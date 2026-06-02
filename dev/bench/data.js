window.BENCHMARK_DATA = {
  "lastUpdate": 1780390135472,
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
      }
    ]
  }
}