package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/nao1215/omokage/internal/profile"
	"github.com/nao1215/omokage/internal/term"
)

// renderExplanationText prints the human-facing detailed report. It leads with
// the high-level, editable drifts (the ones a person or LLM can act on), then
// lists the low-level fingerprint as supporting detail, and finally points at the
// paragraphs that drift most. The header lines (Author, Similarity) match the
// plain `check` output so the two formats read consistently.
func renderExplanationText(w io.Writer, author string, explanation profile.Explanation) {
	if author != "" {
		writef(w, "Author: %s\n", author)
	}
	writeSimilarityLine(w, explanation.Similarity, explanation.SelfSimilarity)
	if explanation.ScoreDriver != "" {
		writef(w, "Score driver: %s\n", explanation.ScoreDriver)
	}
	if explanation.ScoreNote != "" {
		writef(w, "Scoring note: %s\n", explanation.ScoreNote)
	}

	high, low := splitDrifts(explanation.Drifts)

	writeLine(w)
	actionable := actionableDrifts(high)
	if len(actionable) == 0 {
		writeLine(w, "High-level style: within your usual range; no editable feature drifts.")
	} else {
		writeLine(w, "High-level style differences (fix these first):")
		for _, drift := range actionable {
			writef(w, "  %d. %s is %s than reference [%s]\n",
				drift.Priority, drift.Feature, drift.Direction, drift.Category)
			writef(w, "       target %s  reference %s ± %s  (%.1fσ)\n",
				formatValue(drift.Target), formatValue(drift.Mean), formatValue(drift.StdDev), drift.Z)
		}
	}

	lowActionable := actionableDrifts(low)
	if len(lowActionable) > 0 {
		writeLine(w)
		writeLine(w, "Low-level fingerprint drift (supporting detail):")
		for _, drift := range lowActionable {
			writef(w, "  - %s is %s than reference  (%.1fσ)\n",
				drift.Feature, drift.Direction, drift.Z)
		}
	}

	writeLine(w)
	if len(explanation.Segments) == 0 {
		// Reported only in the detailed view, so silence here is informative, not a
		// gap: no single paragraph holds an editable, paragraph-local drift worth
		// acting on. Saying so keeps late-stage tuning from chasing a paragraph that
		// is not actually the problem.
		writeLine(w, "No single paragraph stands out; the remaining drift is spread across the document or in the low-level fingerprint.")
		return
	}
	writeLine(w, "Paragraphs that drift most:")
	for _, segment := range explanation.Segments {
		writef(w, "  #%d (%.1fσ; %s %s): %s\n",
			segment.Index, segment.Z, segment.Feature, segment.Direction, segment.Excerpt)
	}
}

// renderExplanationJSON emits the explanation as a single JSON object designed to
// be read back by an LLM in a revise-and-recheck loop: the high-level (editable)
// drifts and the low-level fingerprint are kept in separate arrays, each carrying
// the target value, reference mean and spread, z-score, and fix priority.
func renderExplanationJSON(w io.Writer, author string, explanation profile.Explanation, warnings []term.Warning) error {
	high, low := splitDrifts(explanation.Drifts)
	payload := explanationJSON{
		Author:         author,
		Similarity:     explanation.Similarity,
		SelfSimilarity: toAnchorJSON(explanation.SelfSimilarity),
		ScoreDriver:    explanation.ScoreDriver,
		ScoreNote:      explanation.ScoreNote,
		HighLevelDrift: toDriftJSON(high),
		LowLevelDrift:  toDriftJSON(low),
		Segments:       toSegmentJSON(explanation.Segments),
		TermWarnings:   toTermWarningJSON(warnings),
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

type explanationJSON struct {
	Author         string             `json:"author"`
	Similarity     int                `json:"similarity"`
	SelfSimilarity *similarityAnchorJSON `json:"self_similarity_anchor,omitempty"`
	ScoreDriver    string             `json:"score_driver,omitempty"`
	ScoreNote      string             `json:"score_note,omitempty"`
	HighLevelDrift []featureDriftJSON `json:"high_level_drift"`
	LowLevelDrift  []featureDriftJSON `json:"low_level_drift"`
	Segments       []segmentJSON      `json:"segments"`
	// TermWarnings reports notation deviations (a non-preferred surface in the
	// draft). It is a separate layer from the similarity score and never affects
	// it. Always present (empty array, not null) so the shape is stable.
	TermWarnings []termWarningJSON `json:"term_warnings"`
}

// termWarningJSON is one notation deviation in `check --format json`. Occurrences
// is reserved for future line/column reporting and is omitted while only counts
// are tracked, so adding it later will not change existing fields.
type termWarningJSON struct {
	GroupKey         string               `json:"group_key"`
	PreferredSurface string               `json:"preferred_surface"`
	UsedSurface      string               `json:"used_surface"`
	Count            int                  `json:"count"`
	Occurrences      []termOccurrenceJSON `json:"occurrences,omitempty"`
}

type termOccurrenceJSON struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// toTermWarningJSON converts term warnings into the check payload, always
// returning a non-nil slice so the JSON shows an empty array rather than null.
func toTermWarningJSON(warnings []term.Warning) []termWarningJSON {
	out := make([]termWarningJSON, 0, len(warnings))
	for _, w := range warnings {
		warning := termWarningJSON{
			GroupKey:         w.GroupKey,
			PreferredSurface: w.PreferredSurface,
			UsedSurface:      w.UsedSurface,
			Count:            w.Count,
		}
		for _, occ := range w.Occurrences {
			warning.Occurrences = append(warning.Occurrences, termOccurrenceJSON{Line: occ.Line, Column: occ.Column})
		}
		out = append(out, warning)
	}
	return out
}

type featureDriftJSON struct {
	Feature         string  `json:"feature"`
	Category        string  `json:"category"`
	Target          float64 `json:"target"`
	ReferenceMean   float64 `json:"reference_mean"`
	ReferenceStddev float64 `json:"reference_stddev"`
	Z               float64 `json:"z"`
	Direction       string  `json:"direction"`
	Priority        int     `json:"priority"`
	Actionable      bool    `json:"actionable"`
}

type segmentJSON struct {
	Index     int     `json:"index"`
	Kind      string  `json:"kind"`
	Excerpt   string  `json:"excerpt"`
	Feature   string  `json:"feature"`
	Category  string  `json:"category"`
	Z         float64 `json:"z"`
	Direction string  `json:"direction"`
}

type similarityAnchorJSON struct {
	Median  int `json:"median"`
	Low     int `json:"low"`
	High    int `json:"high"`
	Samples int `json:"samples"`
}

// splitDrifts separates the prioritized drift list into its high-level and
// low-level halves while preserving the priority order within each.
func splitDrifts(drifts []profile.FeatureDrift) (high, low []profile.FeatureDrift) {
	for _, drift := range drifts {
		if drift.Level == "low" {
			low = append(low, drift)
		} else {
			high = append(high, drift)
		}
	}
	return high, low
}

// actionableDrifts keeps only the drifts worth correcting (those past the drift
// threshold), so the text report does not list features that already match.
func actionableDrifts(drifts []profile.FeatureDrift) []profile.FeatureDrift {
	out := make([]profile.FeatureDrift, 0, len(drifts))
	for _, drift := range drifts {
		if drift.Actionable {
			out = append(out, drift)
		}
	}
	return out
}

func toDriftJSON(drifts []profile.FeatureDrift) []featureDriftJSON {
	out := make([]featureDriftJSON, 0, len(drifts))
	for _, drift := range drifts {
		out = append(out, featureDriftJSON{
			Feature:         drift.Feature,
			Category:        drift.Category,
			Target:          round4(drift.Target),
			ReferenceMean:   round4(drift.Mean),
			ReferenceStddev: round4(drift.StdDev),
			Z:               round4(drift.Z),
			Direction:       drift.Direction,
			Priority:        drift.Priority,
			Actionable:      drift.Actionable,
		})
	}
	return out
}

func toSegmentJSON(segments []profile.SegmentDrift) []segmentJSON {
	out := make([]segmentJSON, 0, len(segments))
	for _, segment := range segments {
		out = append(out, segmentJSON{
			Index:     segment.Index,
			Kind:      segment.Kind,
			Excerpt:   segment.Excerpt,
			Feature:   segment.Feature,
			Category:  segment.Category,
			Z:         round4(segment.Z),
			Direction: segment.Direction,
		})
	}
	return out
}

func toAnchorJSON(anchor *profile.SimilarityAnchor) *similarityAnchorJSON {
	if anchor == nil {
		return nil
	}
	return &similarityAnchorJSON{
		Median:  anchor.Median,
		Low:     anchor.Low,
		High:    anchor.High,
		Samples: anchor.Samples,
	}
}

// formatValue trims feature values to a readable width: ratios sit in [0,1] and
// read best with three decimals, while larger scalars (sentence/paragraph
// lengths) read better rounded to one.
func formatValue(value float64) string {
	if math.Abs(value) < 10 {
		return fmt.Sprintf("%.3f", value)
	}
	return fmt.Sprintf("%.1f", value)
}

// round4 keeps the JSON numbers compact and stable for an LLM to parse without
// drowning in float noise.
func round4(value float64) float64 {
	return math.Round(value*10000) / 10000
}
