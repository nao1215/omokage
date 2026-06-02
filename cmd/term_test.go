package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// The structs below mirror the term-related portions of the show/check JSON
// payloads so the end-to-end test reads the new fields back the way a consumer
// would.

type showTermVariant struct {
	Surface       string `json:"surface"`
	NormalizedKey string `json:"normalized_key"`
	Count         int    `json:"count"`
	DocCount      int    `json:"doc_count"`
}

type showTermGroup struct {
	GroupKey         string            `json:"group_key"`
	PreferredSurface string            `json:"preferred_surface"`
	DocCount         int               `json:"doc_count"`
	TotalCount       int               `json:"total_count"`
	Variants         []showTermVariant `json:"variants"`
}

type showTermsPayload struct {
	TermPreferences []showTermGroup `json:"term_preferences"`
}

type checkTermsPayload struct {
	Similarity   int `json:"similarity"`
	TermWarnings []struct {
		GroupKey         string `json:"group_key"`
		PreferredSurface string `json:"preferred_surface"`
		UsedSurface      string `json:"used_surface"`
		Count            int    `json:"count"`
	} `json:"term_warnings"`
}

func findShowGroup(p showTermsPayload, preferred string) *showTermGroup {
	for i := range p.TermPreferences {
		if p.TermPreferences[i].PreferredSurface == preferred {
			return &p.TermPreferences[i]
		}
	}
	return nil
}

func groupHasSurface(g *showTermGroup, surface string) bool {
	for _, v := range g.Variants {
		if v.Surface == surface {
			return true
		}
	}
	return false
}

// TestTermPreferencesEndToEnd exercises the full feature through the CLI:
// train extracts term preferences, show --format json reports them (including a
// corpus-declared alias bridge), and check --format json flags a draft that uses
// a non-preferred surface — without changing the similarity score.
func TestTermPreferencesEndToEnd(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	corpusDir := filepath.Join(workDir, "posts")
	// DB is used far more than データベース, so DB becomes the preferred surface;
	// the parenthetical "データベース（DB）" declares the alias bridge.
	writeTestFile(t, filepath.Join(corpusDir, "one.md"),
		"データベース（DB）を使う。DB は速い。DB を使う。DB が良い。")
	writeTestFile(t, filepath.Join(corpusDir, "two.md"),
		"DB を使う。DB が好き。DB を選ぶ。データベースも使う。")

	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "posts"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}

	// show --format json must surface the bridged DB/データベース group.
	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show --format json failed: %s", stderr)
	}
	var show showTermsPayload
	if err := json.Unmarshal([]byte(stdout), &show); err != nil {
		t.Fatalf("show JSON invalid: %v\n%s", err, stdout)
	}
	dbGroup := findShowGroup(show, "DB")
	if dbGroup == nil {
		t.Fatalf("expected a term preference group with preferred surface DB:\n%s", stdout)
	}
	if !groupHasSurface(dbGroup, "データベース") {
		t.Errorf("expected データベース to be bridged into the DB group, got %+v", dbGroup)
	}
	keys := map[string]struct{}{}
	for _, v := range dbGroup.Variants {
		keys[v.NormalizedKey] = struct{}{}
	}
	if len(keys) < 2 {
		t.Errorf("bridged group should span >1 normalized_key, got %v", keys)
	}

	// check --format json must flag a draft that uses the non-preferred surface.
	writeTestFile(t, filepath.Join(workDir, "draft.md"), "ＤＢ を整備する。データベースを設計する。")
	code, stdout, stderr = runApp(t, workDir, "check", "--author", "me", "--format", "json", "draft.md")
	if code != 0 {
		t.Fatalf("check --format json failed: %s", stderr)
	}
	var check checkTermsPayload
	if err := json.Unmarshal([]byte(stdout), &check); err != nil {
		t.Fatalf("check JSON invalid: %v\n%s", err, stdout)
	}
	if len(check.TermWarnings) == 0 {
		t.Fatalf("expected term warnings for ＤＢ/データベース, got none:\n%s", stdout)
	}
	for _, w := range check.TermWarnings {
		if w.PreferredSurface != "DB" {
			t.Errorf("warning %+v should point at preferred surface DB", w)
		}
		if w.UsedSurface == "DB" {
			t.Errorf("the preferred surface must never be flagged: %+v", w)
		}
	}

	// The plain (text) check must keep working and stay unaffected by the warnings
	// layer: it still reports output and exits 0.
	if code, plain, stderr := runApp(t, workDir, "check", "--author", "me", "draft.md"); code != 0 {
		t.Fatalf("plain check failed: %s", stderr)
	} else if !strings.Contains(plain, "Similarity") {
		t.Fatalf("plain check should still report a similarity line, got %q", plain)
	}
}

// TestShowJSONTermPreferencesAlwaysPresent checks the shape stays stable: even a
// corpus with no extractable bridged terms still emits a term_preferences array,
// not null.
func TestShowJSONTermPreferencesAlwaysPresent(t *testing.T) {
	t.Parallel()

	workDir := trainedProject(t)
	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show --format json failed: %s", stderr)
	}
	if !strings.Contains(stdout, `"term_preferences"`) {
		t.Fatalf("term_preferences must always be present:\n%s", stdout)
	}
}
