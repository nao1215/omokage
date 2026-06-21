package cmd

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// debimateFixtureDir is the Debimate validation corpus checked into the repo. The
// fixtures are real article bodies (tech, diary, and rewrites) used to assert the
// topic-masking and robust-aggregation behavior at the command level, against the
// same paths a user would pass to `train`/`check`.
const debimateFixtureDir = "../spec/testdata/debimate/ja"

// scoreFixture trains author "me" from the Debimate tech/train corpus once and
// returns the `check --score-only` score of one fixture file, so a test reads the
// same number a user would. The profile is trained per call into the test's own
// workDir, which keeps the tests independent.
func debimateScore(t *testing.T, workDir, file string) int {
	t.Helper()
	code, stdout, stderr := runApp(t, workDir, "check", "--author", "me", "--score-only", file)
	if code != 0 {
		t.Fatalf("check --score-only %s failed: code=%d stderr=%q", file, code, stderr)
	}
	n, err := strconv.Atoi(strings.TrimSpace(stdout))
	if err != nil {
		t.Fatalf("score for %s is not an integer: %q", file, stdout)
	}
	return n
}

func medianInts(values []int) float64 {
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[mid])
	}
	return float64(sorted[mid-1]+sorted[mid]) / 2
}

func globFixtures(t *testing.T, sub string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(debimateFixtureDir, sub, "*.md"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("no fixtures under %s (err=%v)", sub, err)
	}
	// check resolves a relative file argument against the run's workDir (a temp
	// dir), so the fixtures must be passed as absolute paths.
	for i, m := range matches {
		abs, err := filepath.Abs(m)
		if err != nil {
			t.Fatal(err)
		}
		matches[i] = abs
	}
	return matches
}

// fixtureAbs returns the absolute path of one fixture file, since check resolves
// a relative file argument against the run's temp workDir.
func fixtureAbs(t *testing.T, rel string) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join(debimateFixtureDir, rel))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// trainDebimate creates a project in a fresh workDir and trains author "me" from
// the tech/train fixtures, returning the workDir for subsequent checks.
func trainDebimate(t *testing.T) string {
	t.Helper()
	workDir := t.TempDir()
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	trainDir, err := filepath.Abs(filepath.Join(debimateFixtureDir, "tech", "train"))
	if err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := runApp(t, workDir, "train", "--author", "me", trainDir); code != 0 {
		t.Fatalf("train failed: %s", stderr)
	}
	return workDir
}

// TestDebimateGenreSeparation asserts the central ranking the change is meant to
// preserve: against a profile trained on nao's technical posts, his technical
// holdout reads as the same voice (high score), while his personal/reflection
// posts read as a different voice (lower score). The threshold is deliberately
// loose — the test asserts the ordering and a meaningful gap, not brittle exact
// scores.
func TestDebimateGenreSeparation(t *testing.T) {
	t.Parallel()
	workDir := trainDebimate(t)

	var tech, diary []int
	for _, f := range globFixtures(t, "tech/holdout") {
		tech = append(tech, debimateScore(t, workDir, f))
	}
	for _, f := range globFixtures(t, "diary/holdout") {
		diary = append(diary, debimateScore(t, workDir, f))
	}
	sameGenre := medianInts(tech)
	crossGenre := medianInts(diary)

	if sameGenre < 80 {
		t.Fatalf("technical holdout should still read as the author's voice (>=80), got median %.1f (%v)", sameGenre, tech)
	}
	if sameGenre-crossGenre < 15 {
		t.Fatalf("genre separation should be at least 15: same=%.1f cross=%.1f (tech=%v diary=%v)",
			sameGenre, crossGenre, tech, diary)
	}
}

// TestDebimateTopicRobustness asserts the masking goal: swapping only the topic
// tokens (product names, repository names, versions, formats) of a holdout article
// barely moves its score, because those tokens are masked out of the lexical and
// character n-gram fingerprints before scoring.
func TestDebimateTopicRobustness(t *testing.T) {
	t.Parallel()
	workDir := trainDebimate(t)

	pairs := []struct{ original, swapped string }{
		{"tech/holdout/sqluv-https.md", "rewrites/topic_swapped/dbpeek-https.md"},
		{"tech/holdout/omokage.md", "rewrites/topic_swapped/sukima.md"},
	}
	for _, p := range pairs {
		orig := debimateScore(t, workDir, fixtureAbs(t, p.original))
		swapped := debimateScore(t, workDir, fixtureAbs(t, p.swapped))
		diff := orig - swapped
		if diff < 0 {
			diff = -diff
		}
		if diff > 5 {
			t.Fatalf("topic swap should barely move the score (<=5): %s=%d %s=%d diff=%d",
				p.original, orig, p.swapped, swapped, diff)
		}
	}
}

// TestDebimateStyleSensitivity asserts the complementary guarantee: keeping the
// topic but deliberately rewriting the voice (register flipped to 常体, sentences
// chopped short) drops the score well below the same-voice technical holdout, so
// the masking has not blunted sensitivity to a genuine style change.
func TestDebimateStyleSensitivity(t *testing.T) {
	t.Parallel()
	workDir := trainDebimate(t)

	var tech, styleShifted []int
	for _, f := range globFixtures(t, "tech/holdout") {
		tech = append(tech, debimateScore(t, workDir, f))
	}
	for _, f := range globFixtures(t, "rewrites/style_shifted") {
		styleShifted = append(styleShifted, debimateScore(t, workDir, f))
	}
	sameGenre := medianInts(tech)
	shifted := medianInts(styleShifted)

	if sameGenre-shifted < 20 {
		t.Fatalf("a deliberate style shift should drop the score by at least 20: same=%.1f shifted=%.1f (tech=%v shifted=%v)",
			sameGenre, shifted, tech, styleShifted)
	}
}

// TestDebimateDoctorMixedCorpus asserts the corpus-quality contract: a pure
// technical corpus trains cleanly, while mixing in the personal/reflection posts
// (which shift register) makes doctor flag the corpus as less reliable. It asserts
// the relationship (mixed is noisier than pure), not the exact finding text.
func TestDebimateDoctorMixedCorpus(t *testing.T) {
	t.Parallel()
	workDir := t.TempDir()
	if code, _, stderr := runApp(t, workDir, "init"); code != 0 {
		t.Fatalf("init failed: %s", stderr)
	}
	techDir := fixtureAbs(t, filepath.Join("tech", "train"))
	diaryDir := fixtureAbs(t, filepath.Join("diary", "holdout"))

	code, pureOut, stderr := runApp(t, workDir, "doctor", techDir)
	if code != 0 {
		t.Fatalf("doctor on pure tech failed: %s", stderr)
	}
	code, mixedOut, stderr := runApp(t, workDir, "doctor", techDir, diaryDir)
	if code != 0 {
		t.Fatalf("doctor on mixed corpus failed: %s", stderr)
	}

	pureClean := strings.Contains(pureOut, "No problems found")
	mixedFlags := strings.Contains(mixedOut, "[warning]") || strings.Contains(mixedOut, "[notice]")
	if !pureClean {
		t.Fatalf("pure technical corpus should not gain noisy warnings, got:\n%s", pureOut)
	}
	if !mixedFlags {
		t.Fatalf("mixing tech and personal posts should make doctor flag the corpus, got:\n%s", mixedOut)
	}
}
