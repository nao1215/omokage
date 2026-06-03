package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nao1215/omokage/internal/config"
	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/quality"
)

// thinCorpus seeds a directory with a few very short documents — too few and too
// short to train a reliable profile — so the quality checks have something to
// flag.
func thinCorpus(t *testing.T, dir string) {
	t.Helper()
	writeTestFile(t, filepath.Join(dir, "a.md"), "今日は晴れ。散歩した。")
	writeTestFile(t, filepath.Join(dir, "b.md"), "本を読んだ。良かった。")
	writeTestFile(t, filepath.Join(dir, "c.md"), "映画を見た。面白い。")
}

// richDocument is a single document long enough (well over the short-document
// threshold) and written in a steady polite voice, so a corpus of several copies
// reads as solid: enough material, one consistent voice, no outliers.
const richDocument = "今日は朝から良い天気でした。近所の公園までゆっくり散歩に出かけました。" +
	"道の途中で猫に出会い、しばらくその様子を眺めていました。帰り道にパン屋へ寄り、" +
	"焼きたてのパンをいくつか買って帰りました。午後は部屋で本を読みながら、" +
	"温かいお茶を何杯も飲みました。夕方になると、空が少しずつ赤く染まっていきました。" +
	"窓の外を眺めながら、今日は本当に穏やかな一日だったと感じました。"

// richCorpus seeds dir with n copies of richDocument, a clean corpus that should
// raise no quality findings.
func richCorpus(t *testing.T, dir string, n int) {
	t.Helper()
	for i := range n {
		writeTestFile(t, filepath.Join(dir, fmt.Sprintf("post%d.md", i)), richDocument)
	}
}

// TestDoctorReportsThinCorpus verifies doctor rates a thin corpus weak, lists findings, and writes no store.
func TestDoctorReportsThinCorpus(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))

	// doctor needs no init: it reads the corpus and reports, training nothing.
	code, stdout, stderr := runApp(t, workDir, "doctor", "posts")
	if code != 0 {
		t.Fatalf("doctor failed: stderr=%q", stderr)
	}
	for _, want := range []string{"Reliability: weak", "Findings:", "documents", "→", "not writing quality"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("doctor output missing %q:\n%s", want, stdout)
		}
	}
	// doctor must not create a store or profile: it is read-only.
	if _, err := os.Stat(filepath.Join(workDir, "omokage.toml")); err == nil {
		t.Fatal("doctor must not initialize a project")
	}
}

// TestDoctorCleanCorpusReportsGood verifies doctor rates a solid corpus good.
func TestDoctorCleanCorpusReportsGood(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	// Eight longer, consistent documents: enough material, one voice, no outliers.
	richCorpus(t, filepath.Join(workDir, "posts"), 8)

	code, stdout, stderr := runApp(t, workDir, "doctor", "posts")
	if code != 0 {
		t.Fatalf("doctor failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, "Reliability: good") {
		t.Fatalf("a clean corpus should report good reliability:\n%s", stdout)
	}
	if !strings.Contains(stdout, "No problems found") {
		t.Fatalf("a clean corpus should say so:\n%s", stdout)
	}
}

// TestDoctorJSON verifies doctor --format json emits valid, complete findings.
func TestDoctorJSON(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))

	code, stdout, stderr := runApp(t, workDir, "doctor", "--format", "json", "posts")
	if code != 0 {
		t.Fatalf("doctor --format json failed: stderr=%q", stderr)
	}
	var payload struct {
		DocumentCount int    `json:"document_count"`
		Reliability   string `json:"reliability"`
		Findings      []struct {
			Code     string `json:"code"`
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
			Action   string `json:"action"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("doctor --format json did not emit valid JSON: %v\n%s", err, stdout)
	}
	if payload.DocumentCount != 3 {
		t.Fatalf("expected 3 documents, got %d", payload.DocumentCount)
	}
	if payload.Reliability == "" {
		t.Fatal("expected a reliability rating")
	}
	if len(payload.Findings) == 0 {
		t.Fatal("a thin corpus should produce findings")
	}
	for _, finding := range payload.Findings {
		if finding.Code == "" || finding.Severity == "" || finding.Action == "" {
			t.Fatalf("every finding needs a code, severity, and action: %+v", finding)
		}
	}
}

// TestDoctorRejectsUnknownFormat verifies an unknown --format fails.
func TestDoctorRejectsUnknownFormat(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))

	code, _, stderr := runApp(t, workDir, "doctor", "--format", "yaml", "posts")
	if code == 0 {
		t.Fatal("an unknown format should fail")
	}
	if !strings.Contains(stderr, "unknown --format") {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

// TestDoctorRequiresInput verifies doctor with no input fails with a clear message.
func TestDoctorRequiresInput(t *testing.T) {
	t.Parallel()

	code, _, stderr := runApp(t, t.TempDir(), "doctor")
	if code == 0 {
		t.Fatal("doctor with no input should fail")
	}
	if !strings.Contains(stderr, "missing INPUT") {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

// TestDoctorRejectsURL verifies a URL input is rejected like train.
func TestDoctorRejectsURL(t *testing.T) {
	t.Parallel()

	code, _, stderr := runApp(t, t.TempDir(), "doctor", "https://example.com/post")
	if code == 0 {
		t.Fatal("a URL input should fail")
	}
	if !strings.Contains(stderr, "URL inputs are not supported") {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

// TestTrainPrintsReliabilityOnStdout verifies the post-training summary reports a
// thin corpus weak on stdout (not gated behind a terminal) with a pointer to
// doctor, while the trained-profile confirmation stays present.
func TestTrainPrintsReliabilityOnStdout(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}

	code, stdout, stderr := runApp(t, workDir, "train", "--author", "me", "posts")
	if code != 0 {
		t.Fatalf("train failed: stderr=%q", stderr)
	}
	if !strings.Contains(stdout, `Trained author "me"`) {
		t.Fatalf("train confirmation missing from stdout: %q", stdout)
	}
	// The reliability summary is on stdout so a person, a script, and an LLM all see
	// it — matching what train --help promises.
	if !strings.Contains(stdout, "Corpus reliability: weak") {
		t.Fatalf("train should print the corpus reliability on stdout: %q", stdout)
	}
	if !strings.Contains(stdout, "omokage doctor") {
		t.Fatalf("a weak corpus should point at doctor: %q", stdout)
	}
}

// TestTrainCleanCorpusReportsGood verifies a solid corpus trains with a single
// good reliability line and no findings or doctor pointer.
func TestTrainCleanCorpusReportsGood(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	richCorpus(t, filepath.Join(workDir, "posts"), 8)
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}

	_, stdout, _ := runApp(t, workDir, "train", "--author", "me", "posts")
	if !strings.Contains(stdout, "Corpus reliability: good.") {
		t.Fatalf("a clean corpus should report good reliability: %q", stdout)
	}
	if strings.Contains(stdout, "omokage doctor") || strings.Contains(stdout, "→") {
		t.Fatalf("a clean corpus should not list findings or point at doctor: %q", stdout)
	}
}

// TestRenderTrainQualitySummaryContent verifies the summary renderer: a thin
// corpus yields a weak headline, an actionable line, and a doctor pointer, while a
// clean report yields a single good line.
func TestRenderTrainQualitySummaryContent(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	postsDir := filepath.Join(workDir, "posts")
	thinCorpus(t, postsDir)
	files, err := feature.CollectFiles(postsDir)
	if err != nil {
		t.Fatal(err)
	}
	dist, docs, err := feature.ExtractCorpusDocuments(files)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{workDir: workDir}
	report := quality.AssessCorpus(dist, app.qualityDocuments(docs), config.Default("test").Features)

	var buf bytes.Buffer
	renderTrainQualitySummary(&buf, report, []string{"posts"})
	out := buf.String()
	for _, want := range []string{"Corpus reliability: weak", "→", "omokage doctor posts"} {
		if !strings.Contains(out, want) {
			t.Fatalf("quality summary missing %q:\n%s", want, out)
		}
	}

	// A clean report renders just the one-line good rating, no findings.
	var clean bytes.Buffer
	renderTrainQualitySummary(&clean, quality.Report{}, []string{"posts"})
	if !strings.Contains(clean.String(), "Corpus reliability: good.") {
		t.Fatalf("a clean report should render the good line, got %q", clean.String())
	}
	if strings.Contains(clean.String(), "doctor") {
		t.Fatalf("a clean report should not point at doctor, got %q", clean.String())
	}
}

// TestShowJSONReportsReliability verifies show --format json carries a reliability rating and quality findings.
func TestShowJSONReportsReliability(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "posts"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}

	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show --format json failed: stderr=%q", stderr)
	}
	var payload struct {
		Reliability     string `json:"reliability"`
		QualityFindings []struct {
			Code   string `json:"code"`
			Action string `json:"action"`
		} `json:"quality_findings"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("show --format json invalid: %v\n%s", err, stdout)
	}
	if payload.Reliability == "" {
		t.Fatal("show --format json should report a reliability rating")
	}
	if len(payload.QualityFindings) == 0 {
		t.Fatal("a thin corpus profile should carry quality findings")
	}
}

// TestShowJSONKeepsOutlierFindingFromTrain verifies that a per-document finding
// (an outlier), which the stored distribution alone cannot reproduce, is recorded
// at train time and surfaces through show — the gap the user reported.
func TestShowJSONKeepsOutlierFindingFromTrain(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	dir := filepath.Join(workDir, "posts")
	// Many consistent polite documents and one in a wholly different voice, so the
	// odd one is recorded as an outlier finding.
	for i := range 7 {
		writeTestFile(t, filepath.Join(dir, fmt.Sprintf("p%d.md", i)),
			"今日は良い天気でした。近所の公園まで散歩に出かけました。猫に出会い、しばらく眺めていました。とても穏やかな一日でした。")
	}
	writeTestFile(t, filepath.Join(dir, "odd.md"),
		"OUTPUT FORMAT: JSON. RUN: omokage check --format json. SEE: docs/spec.md, table 3.4, row 12.")

	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "posts"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}

	code, stdout, stderr := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if code != 0 {
		t.Fatalf("show --format json failed: stderr=%q", stderr)
	}
	var payload struct {
		QualityFindings []struct {
			Code string `json:"code"`
		} `json:"quality_findings"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("show --format json invalid: %v\n%s", err, stdout)
	}
	found := false
	for _, f := range payload.QualityFindings {
		if f.Code == "outlier_documents" {
			found = true
		}
	}
	if !found {
		t.Fatalf("the train-time outlier finding should survive into show: %+v", payload.QualityFindings)
	}
}

// TestShowSummaryOmitsTermPreferences verifies --summary drops the term list for a
// lighter JSON while keeping reliability and quality findings, and that the
// default output keeps the term_preferences key.
func TestShowSummaryOmitsTermPreferences(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	thinCorpus(t, filepath.Join(workDir, "posts"))
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", "posts"); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}

	// Default JSON keeps the term_preferences key (stable shape).
	_, full, _ := runApp(t, workDir, "show", "--author", "me", "--format", "json")
	if !strings.Contains(full, `"term_preferences"`) {
		t.Fatalf("default show --format json should keep term_preferences: %s", full)
	}

	// --summary drops it but keeps reliability and quality findings.
	_, summary, _ := runApp(t, workDir, "show", "--author", "me", "--format", "json", "--summary")
	if strings.Contains(summary, `"term_preferences"`) {
		t.Fatalf("--summary should omit term_preferences: %s", summary)
	}
	if !strings.Contains(summary, `"reliability"`) || !strings.Contains(summary, `"quality_findings"`) {
		t.Fatalf("--summary should keep reliability and quality findings: %s", summary)
	}
	// The lighter payload must still be valid JSON.
	var probe map[string]any
	if err := json.Unmarshal([]byte(summary), &probe); err != nil {
		t.Fatalf("--summary JSON invalid: %v\n%s", err, summary)
	}
}
