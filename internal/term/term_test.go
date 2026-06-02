package term

import (
	"encoding/json"
	"strings"
	"testing"
)

// findGroup returns the group a surface belongs to, by matching any variant's
// surface, or nil when the surface is absent.
func findGroup(p Profile, surface string) *Group {
	for i := range p.Groups {
		for _, v := range p.Groups[i].Variants {
			if v.Surface == surface {
				return &p.Groups[i]
			}
		}
	}
	return nil
}

// TestNormalizeKeyFoldsCaseAndWidth pins the normalized_key contract: case,
// full-width ASCII, and surrounding punctuation fold away, so DB / db / ＤＢ share
// one key, while Japanese text is left intact.
func TestNormalizeKeyFoldsCaseAndWidth(t *testing.T) {
	t.Parallel()

	cases := []struct {
		surface string
		want    string
	}{
		{"DB", "db"},
		{"db", "db"},
		{"ＤＢ", "db"},   // full-width letters fold to ASCII
		{"(DB)", "db"}, // surrounding punctuation is trimmed
		{"API,", "api"},
		{"HTTP", "http"},
		{"データベース", "データベース"}, // Japanese is not folded
	}
	for _, c := range cases {
		if got := normalizeKey(c.surface); got != c.want {
			t.Errorf("normalizeKey(%q) = %q, want %q", c.surface, got, c.want)
		}
	}
}

// TestNormalizationGroupsSameKey checks that DB, db, and ＤＢ collapse to a single
// normalized_key and a single group, with three distinct surfaces, when no alias
// bridge is involved.
func TestNormalizationGroupsSameKey(t *testing.T) {
	t.Parallel()

	p := ExtractDocuments([]string{"DB db ＤＢ を使う。"})
	g := findGroup(p, "DB")
	if g == nil {
		t.Fatal("expected a group containing DB")
	}
	if g.GroupKey != "term:db" {
		t.Fatalf("group_key = %q, want term:db", g.GroupKey)
	}
	for _, v := range g.Variants {
		if v.NormalizedKey != "db" {
			t.Errorf("surface %q has normalized_key %q, want db", v.Surface, v.NormalizedKey)
		}
	}
	if len(g.Variants) != 3 {
		t.Fatalf("expected 3 surface variants (DB, db, ＤＢ), got %d", len(g.Variants))
	}
}

// TestNormalizedKeyAndGroupKeySeparation pins the responsibility split: a group
// formed by normalization alone has every variant sharing one normalized_key
// equal to the group_key suffix; a group formed by an alias bridge spans more
// than one normalized_key, which is how a reader tells the two apart.
func TestNormalizedKeyAndGroupKeySeparation(t *testing.T) {
	t.Parallel()

	// Normalization-only group: DB and db.
	norm := ExtractDocuments([]string{"DB db を使う。"})
	g := findGroup(norm, "DB")
	keys := distinctNormalizedKeys(g)
	if len(keys) != 1 {
		t.Fatalf("normalization-only group should have one normalized_key, got %v", keys)
	}
	if g.GroupKey != "term:"+keys[0] {
		t.Fatalf("group_key %q should equal term:%s for a normalization-only group", g.GroupKey, keys[0])
	}

	// Alias-bridged group: データベース（DB）.
	bridged := ExtractDocuments([]string{"データベース（DB）を使う。データベースは便利。DBも便利。"})
	bg := findGroup(bridged, "DB")
	if bg == nil {
		t.Fatal("expected a bridged group containing DB")
	}
	if len(distinctNormalizedKeys(bg)) < 2 {
		t.Fatalf("alias-bridged group should span >1 normalized_key, got %v", distinctNormalizedKeys(bg))
	}
}

func distinctNormalizedKeys(g *Group) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, v := range g.Variants {
		if _, ok := seen[v.NormalizedKey]; !ok {
			seen[v.NormalizedKey] = struct{}{}
			out = append(out, v.NormalizedKey)
		}
	}
	return out
}

// TestAliasBridgeMergesGroup checks the core bridge: when the corpus declares
// "データベース（DB）", データベース and DB end up in the same group.
func TestAliasBridgeMergesGroup(t *testing.T) {
	t.Parallel()

	patterns := []string{
		"データベース（DB）を採用する。データベースは速い。DBも速い。",
		"DB（データベース）を採用する。データベースは速い。DBも速い。",
		"データベース（以下 DB）を採用する。データベースは速い。DBも速い。",
		"データベース。以下、DB。データベースは速い。DBも速い。",
	}
	for _, doc := range patterns {
		p := ExtractDocuments([]string{doc})
		gJa := findGroup(p, "データベース")
		gEn := findGroup(p, "DB")
		if gJa == nil || gEn == nil {
			t.Fatalf("%q: expected both データベース and DB as candidates", doc)
		}
		if gJa.GroupKey != gEn.GroupKey {
			t.Errorf("%q: データベース (%s) and DB (%s) should share a group_key", doc, gJa.GroupKey, gEn.GroupKey)
		}
	}
}

// TestConservativeAliasDetectionRejectsArbitraryParens checks that a bare "A（B）"
// is never treated as an alias when neither side is a Japanese-phrase ↔ acronym
// pair.
func TestConservativeAliasDetectionRejectsArbitraryParens(t *testing.T) {
	t.Parallel()

	// 関数 (2 kanji, too short to be the long side) and function (lowercase, not an
	// acronym): must not merge.
	p := ExtractDocuments([]string{"関数（function）を定義する。関数は便利。functionも便利。"})
	gJa := findGroup(p, "関数")
	gEn := findGroup(p, "function")
	if gJa == nil || gEn == nil {
		t.Fatal("expected both 関数 and function as candidates")
	}
	if gJa.GroupKey == gEn.GroupKey {
		t.Errorf("関数 and function must not be alias-merged (no acronym bridge), both = %s", gJa.GroupKey)
	}
}

// TestIfukaFormBridges checks the "以下、X" abbreviation declaration links the
// long phrase and the acronym.
func TestIfukaFormBridges(t *testing.T) {
	t.Parallel()

	p := ExtractDocuments([]string{"継続的インテグレーション。以下、CI。継続的インテグレーションは重要。CIは重要。"})
	gJa := findGroup(p, "継続的インテグレーション")
	gEn := findGroup(p, "CI")
	if gJa == nil || gEn == nil {
		t.Fatal("expected both 継続的インテグレーション and CI as candidates")
	}
	if gJa.GroupKey != gEn.GroupKey {
		t.Errorf("「以下、CI」should bridge 継続的インテグレーション and CI: %s vs %s", gJa.GroupKey, gEn.GroupKey)
	}
}

// TestNoFalseMergeWithoutBridge checks that two genuine synonyms stay in separate
// groups when the corpus never declares a bridge between them.
func TestNoFalseMergeWithoutBridge(t *testing.T) {
	t.Parallel()

	p := ExtractDocuments([]string{"優先度を上げる。プライオリティを上げる。優先度は重要。プライオリティは重要。"})
	gA := findGroup(p, "優先度")
	gB := findGroup(p, "プライオリティ")
	if gA == nil || gB == nil {
		t.Fatal("expected both 優先度 and プライオリティ as candidates")
	}
	if gA.GroupKey == gB.GroupKey {
		t.Errorf("優先度 and プライオリティ must not merge without a bridge, both = %s", gA.GroupKey)
	}
}

// TestAmbiguousAcronymDoesNotChain checks that a short acronym glossing several
// different Japanese phrases (each "X（DB）") does not chain those phrases into
// one group: an ambiguous acronym is dropped rather than fusing unrelated terms.
func TestAmbiguousAcronymDoesNotChain(t *testing.T) {
	t.Parallel()

	p := ExtractDocuments([]string{
		"データベース（DB）を使う。データベースは速い。DBも速い。" +
			"ダッシュボード（DB）を見る。ダッシュボードは便利。DBは便利。",
	})
	gA := findGroup(p, "データベース")
	gB := findGroup(p, "ダッシュボード")
	if gA == nil || gB == nil {
		t.Fatal("expected both データベース and ダッシュボード as candidates")
	}
	if gA.GroupKey == gB.GroupKey {
		t.Errorf("an ambiguous acronym (DB→two phrases) must not chain データベース and ダッシュボード, both = %s", gA.GroupKey)
	}
}

// TestUnambiguousAcronymStillBridges guards against over-correction: the same
// phrase↔acronym pair repeated across documents must still bridge.
func TestUnambiguousAcronymStillBridges(t *testing.T) {
	t.Parallel()

	p := ExtractDocuments([]string{
		"データベース（DB）を使う。DBは速い。",
		"データベースを設計する。DBを設計する。",
	})
	gJa := findGroup(p, "データベース")
	gEn := findGroup(p, "DB")
	if gJa == nil || gEn == nil || gJa.GroupKey != gEn.GroupKey {
		t.Fatalf("a repeated unambiguous pair should still bridge: %v / %v", gJa, gEn)
	}
}

// TestPreferredSurfaceSelectionOrder pins the fixed selection order:
// doc_count, then count, then ascending surface as a stable tie-break.
func TestPreferredSurfaceSelectionOrder(t *testing.T) {
	t.Parallel()

	t.Run("doc_count wins over count", func(t *testing.T) {
		t.Parallel()
		// DB: 2 docs, count 2. db: 1 doc, count 5. Higher doc_count wins.
		p := ExtractDocuments([]string{"DB", "DB", "db db db db db"})
		if got := findGroup(p, "DB").PreferredSurface; got != "DB" {
			t.Fatalf("preferred = %q, want DB (higher doc_count beats higher count)", got)
		}
	})

	t.Run("count wins when doc_count ties", func(t *testing.T) {
		t.Parallel()
		// Both in 2 docs; db occurs more often.
		p := ExtractDocuments([]string{"DB db db", "DB db db"})
		if got := findGroup(p, "DB").PreferredSurface; got != "db" {
			t.Fatalf("preferred = %q, want db (higher count on equal doc_count)", got)
		}
	})

	t.Run("surface ascending breaks a full tie", func(t *testing.T) {
		t.Parallel()
		// Both in 2 docs with equal counts; ASCII 'D' < 'd' so DB wins.
		p := ExtractDocuments([]string{"DB db", "DB db"})
		if got := findGroup(p, "DB").PreferredSurface; got != "DB" {
			t.Fatalf("preferred = %q, want DB (ascending surface tie-break)", got)
		}
	})
}

// TestProfileLocality checks that two different corpora yield different preferred
// surfaces for the same concept: term data is profile-local, so profile A can
// prefer DB while profile B prefers db.
func TestProfileLocality(t *testing.T) {
	t.Parallel()

	profileA := ExtractDocuments([]string{"DB を使う。", "DB を使う。", "db"})
	profileB := ExtractDocuments([]string{"db を使う。", "db を使う。", "DB"})

	if got := findGroup(profileA, "DB").PreferredSurface; got != "DB" {
		t.Fatalf("profile A preferred = %q, want DB", got)
	}
	if got := findGroup(profileB, "db").PreferredSurface; got != "db" {
		t.Fatalf("profile B preferred = %q, want db", got)
	}
}

// TestCheckTextFlagsNonPreferredSurface checks that a draft using a non-preferred
// surface (including a full-width or bridged form) produces a warning, and a
// draft using the preferred surface does not.
func TestCheckTextFlagsNonPreferredSurface(t *testing.T) {
	t.Parallel()

	// Profile prefers DB; AI is bridged to 人工知能 (人工知能 wins on counts).
	p := ExtractDocuments([]string{
		"DB を使う。DB は速い。",
		"人工知能（AI）を使う。人工知能は賢い。人工知能はすごい。",
	})

	warnings := p.CheckText("ＤＢ に保存し、AI で処理する。")
	if len(warnings) == 0 {
		t.Fatal("expected warnings for non-preferred surfaces ＤＢ and AI")
	}
	gotBySurface := map[string]Warning{}
	for _, w := range warnings {
		gotBySurface[w.UsedSurface] = w
	}
	if w, ok := gotBySurface["ＤＢ"]; !ok || w.PreferredSurface != "DB" {
		t.Errorf("expected ＤＢ flagged with preferred DB, got %+v", w)
	}
	if w, ok := gotBySurface["AI"]; !ok || w.PreferredSurface != "人工知能" {
		t.Errorf("expected AI flagged with preferred 人工知能, got %+v", w)
	}

	// Using the preferred surface produces no warning for that group.
	for _, w := range p.CheckText("DB に保存する。") {
		if w.UsedSurface == "DB" {
			t.Errorf("DB is the preferred surface and must not be flagged: %+v", w)
		}
	}
}

// TestWarningOccurrencesShape locks the forward-compatible JSON shape: with no
// occurrences the field is omitted entirely; adding occurrences later surfaces a
// line/column array without changing any existing field.
func TestWarningOccurrencesShape(t *testing.T) {
	t.Parallel()

	without, err := json.Marshal(Warning{GroupKey: "term:db", PreferredSurface: "DB", UsedSurface: "db", Count: 2})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(without), "occurrences") {
		t.Fatalf("occurrences must be omitted when empty, got %s", without)
	}

	with, err := json.Marshal(Warning{
		GroupKey: "term:db", PreferredSurface: "DB", UsedSurface: "db", Count: 1,
		Occurrences: []Occurrence{{Line: 3, Column: 12}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"occurrences"`, `"line":3`, `"column":12`} {
		if !strings.Contains(string(with), want) {
			t.Fatalf("expected %s in %s", want, with)
		}
	}
}

// TestDeterministic checks that extraction is reproducible: the same documents
// always yield the same groups, keys, preferred surfaces, and counts.
func TestDeterministic(t *testing.T) {
	t.Parallel()

	docs := []string{
		"データベース（DB）を使う。DBは速い。継続的インテグレーション。以下、CI。CIは重要。",
		"DB と CI を組み合わせる。データベースは便利。",
	}
	first, err := json.Marshal(ExtractDocuments(docs))
	if err != nil {
		t.Fatal(err)
	}
	for range 5 {
		again, err := json.Marshal(ExtractDocuments(docs))
		if err != nil {
			t.Fatal(err)
		}
		if string(first) != string(again) {
			t.Fatalf("extraction is not deterministic:\n%s\nvs\n%s", first, again)
		}
	}
}
