package term

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExtractCorpusReadsFiles checks the file-reading entry point: counts are
// aggregated across files and doc_count reflects how many files held a surface.
func TestExtractCorpusReadsFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	write := func(name, body string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	a := write("a.md", "DB を使う。DB は速い。")
	b := write("b.md", "DB を使う。データベースも使う。")

	p, err := ExtractCorpus([]string{a, b})
	if err != nil {
		t.Fatalf("ExtractCorpus: %v", err)
	}
	g := findGroup(p, "DB")
	if g == nil {
		t.Fatal("expected a DB group")
	}
	var db *Variant
	for i := range g.Variants {
		if g.Variants[i].Surface == "DB" {
			db = &g.Variants[i]
		}
	}
	if db == nil {
		t.Fatal("expected a DB variant")
	}
	if db.Count != 3 || db.DocCount != 2 {
		t.Fatalf("DB variant = count %d / doc_count %d, want 3 / 2", db.Count, db.DocCount)
	}
}

// TestExtractCorpusReadError surfaces a read failure rather than silently
// producing an empty profile.
func TestExtractCorpusReadError(t *testing.T) {
	t.Parallel()

	_, err := ExtractCorpus([]string{filepath.Join(t.TempDir(), "does-not-exist.md")})
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

// TestExtractDocumentsEmptyInputs checks that empty, whitespace-only, and
// code-only documents yield no term groups (code is stripped before extraction).
func TestExtractDocumentsEmptyInputs(t *testing.T) {
	t.Parallel()

	cases := map[string][]string{
		"empty":      {""},
		"whitespace": {"   \n\n  \t "},
		"code only":  {"```go\npackage main\nfunc main() { DB() }\n```"},
	}
	for name, docs := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if p := ExtractDocuments(docs); len(p.Groups) != 0 {
				t.Fatalf("expected no groups for %s input, got %+v", name, p.Groups)
			}
		})
	}
}

// TestStripNoiseExcludesURLsFrontmatterAndHTML checks that link URLs, YAML front
// matter, and HTML tags never become term candidates, while visible prose does.
func TestExtractDocumentsExcludesNoiseTokens(t *testing.T) {
	t.Parallel()

	doc := "---\ntitle: foo\nimage: images/cover.jpg\n---\n\n" +
		"詳しくは [データベース入門](https://example.com/db/intro.html) を参照。\n" +
		`<img src="images/x.png" alt="DB">` + "\n\nデータベースは便利。"

	p := ExtractDocuments([]string{doc})
	for _, bad := range []string{"https", "com", "example", "jpg", "png", "images", "img", "src", "html", "title"} {
		if g := findGroup(p, bad); g != nil {
			t.Errorf("noise token %q must not be a candidate (group %s)", bad, g.GroupKey)
		}
	}
	if findGroup(p, "データベース") == nil {
		t.Error("visible prose term データベース should still be a candidate")
	}
}

// TestGroupKeyIndependentOfUnionOrder checks that the deterministic group_key
// (smallest normalized_key in the component) does not depend on the order links
// are applied — both union directions must give the same id.
func TestGroupKeyIndependentOfUnionOrder(t *testing.T) {
	t.Parallel()

	forward := newUnionFind()
	for _, k := range []string{"db", "data", "store"} {
		forward.add(k)
	}
	forward.union("db", "data")
	forward.union("data", "store")

	reverse := newUnionFind()
	for _, k := range []string{"db", "data", "store"} {
		reverse.add(k)
	}
	reverse.union("store", "data")
	reverse.union("data", "db")

	gf := assignGroupKeys(forward)
	gr := assignGroupKeys(reverse)
	for _, k := range []string{"db", "data", "store"} {
		if gf[forward.find(k)] != "term:data" || gr[reverse.find(k)] != "term:data" {
			t.Fatalf("group_key for %q: forward=%s reverse=%s, want term:data both", k, gf[forward.find(k)], gr[reverse.find(k)])
		}
	}
}

// TestCheckTextEmptyProfile checks that checking against a profile with no terms
// returns no warnings rather than panicking.
func TestCheckTextEmptyProfile(t *testing.T) {
	t.Parallel()

	if w := (Profile{}).CheckText("DB を使う。"); len(w) != 0 {
		t.Fatalf("empty profile should produce no warnings, got %+v", w)
	}
}

// TestKeepASCIIRejectsNoise pins the ASCII candidate filter: stopwords,
// single-letter runs, and pure-digit runs are dropped; real tokens are kept.
func TestKeepASCIIRejectsNoise(t *testing.T) {
	t.Parallel()

	kept := map[string]bool{"DB": true, "API": true, "GoReleaser": true, "the": false, "a": false, "2026": false}
	for surface, want := range kept {
		if got := keepASCII(surface); got != want {
			t.Errorf("keepASCII(%q) = %v, want %v", surface, got, want)
		}
	}
}
