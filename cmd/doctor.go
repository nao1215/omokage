package cmd

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/quality"
)

// qualityDisclaimer is repeated wherever a quality report is shown so the report
// is never mistaken for a judgement of the writing itself. omokage measures
// style, and the quality report measures only whether there is enough consistent
// material to measure style reliably.
const qualityDisclaimer = "These checks look at sample size and consistency, not writing quality."

// qualityDocuments turns the extracted corpus documents into the quality view,
// naming each file relative to the working directory when it sits inside it so a
// finding can point at a readable path rather than a long absolute one.
func (a *App) qualityDocuments(docs []feature.Document) []quality.Document {
	out := make([]quality.Document, 0, len(docs))
	for _, doc := range docs {
		out = append(out, quality.Document{
			Name:    displayPath(a.workDir, doc.Path),
			Metrics: doc.Metrics,
		})
	}
	return out
}

// displayPath renders a file path for a report: relative to baseDir when it lives
// under it (the common case, an input inside the project), otherwise the cleaned
// absolute path so the user can still locate it. A relative path that climbs out
// of baseDir (starts with "..") is not an improvement over the absolute one, so
// the absolute path is kept.
func displayPath(baseDir, path string) string {
	rel, err := filepath.Rel(baseDir, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return filepath.Clean(path)
	}
	return rel
}

// runDoctor implements `omokage doctor INPUT...`: it extracts the corpus, assesses
// its quality, and prints the report as text or JSON. It trains and writes
// nothing, and needs no store (falling back to the default feature weights, like
// diff).
func (a *App) runDoctor(args []string) int {
	flagSet := newFlagSet("doctor", a.stderr)
	format := flagSet.String("format", formatText, "output format: text or json")
	scopeF := registerScopeFlags(flagSet)
	flagSet.Usage = func() {
		writef(a.stderr, "Check whether a corpus is solid enough to train a reliable profile.\n")
		writef(a.stderr, "doctor reads the files but trains nothing and writes nothing. It reports sample\n")
		writef(a.stderr, "size, document length, and how consistent the voice is, with a next step for each\n")
		writef(a.stderr, "issue, so you can curate the corpus before training.\n")
		writef(a.stderr, "Usage: omokage doctor [--format text|json] INPUT...\n")
		writef(a.stderr, "\nINPUT is one or more directories and/or .md/.txt files, exactly like train.\n")
		writef(a.stderr, "The flags below are optional and only select which feature weights to use.\n")
		printFlagDefaults(a.stderr, flagSet)
	}
	if code, ok := parseArgs(flagSet, args); !ok {
		return code
	}
	if flagSet.NArg() == 0 {
		return a.usageError(flagSet, "missing INPUT: pass one or more directories or .md/.txt files")
	}
	if *format != formatText && *format != formatJSON {
		writef(a.stderr, "unknown --format %q: want text or json\n", *format)
		flagSet.Usage()
		return 1
	}

	// doctor needs only the feature set, not a profile, so it works without any
	// store: an active scope supplies the feature weights, and the absence of a
	// store falls back to the defaults rather than erroring — the same lenient
	// resolution diff uses. A store that exists but is broken is still surfaced.
	features, ok := a.featuresOrDefault(scopeF)
	if !ok {
		return 1
	}

	sources, files, err := gatherTrainingInputs(a.workDir, flagSet.Args())
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	if len(files) == 0 {
		writef(a.stderr, "no supported files found in %s\n", strings.Join(sources, ", "))
		return 1
	}

	dist, docs, err := feature.ExtractCorpusDocuments(files)
	if err != nil {
		writeLine(a.stderr, err)
		return 1
	}
	if dist.DocumentCount == 0 {
		writef(a.stderr, "no usable text found in %s (all files were empty)\n", strings.Join(sources, ", "))
		return 1
	}

	report := quality.AssessCorpus(dist, a.qualityDocuments(docs), features)
	if *format == formatJSON {
		if err := renderDoctorJSON(a.stdout, report); err != nil {
			writeLine(a.stderr, err)
			return 1
		}
		return 0
	}
	renderDoctorText(a.stdout, report)
	return 0
}

// renderDoctorText prints the human-facing corpus report: a one-line size
// summary, the reliability headline, the findings (each with its detail and next
// action), and the disclaimer that this is about sample adequacy, not merit.
func renderDoctorText(w io.Writer, report quality.Report) {
	writef(w, "Corpus: %s, %d sentences, %d characters (avg %d per document)\n",
		quality.Documents(report.DocumentCount), report.SentenceCount, report.CharacterCount, report.AverageDocumentCharacters)
	writef(w, "Reliability: %s\n", report.Reliability())
	writeLine(w)

	if len(report.Findings) == 0 {
		writeLine(w, "No problems found: enough material, a consistent voice, and no obvious outliers.")
		writeLine(w)
		writeLine(w, qualityDisclaimer)
		return
	}

	writeLine(w, "Findings:")
	for _, finding := range report.Findings {
		writef(w, "- [%s] %s\n", finding.Severity, finding.Summary)
		if finding.Detail != "" {
			writef(w, "    %s\n", finding.Detail)
		}
		if finding.Action != "" {
			writef(w, "    → %s\n", finding.Action)
		}
	}
	writeLine(w)
	writeLine(w, qualityDisclaimer)
}

// renderQualityNotes prints the post-training hint to stderr: nothing when the
// corpus looks clean, otherwise the reliability headline, each finding's summary
// and next action, and a pointer to the full `doctor` report. It writes to stderr
// so the trained-profile confirmation on stdout stays clean and script-friendly.
func renderQualityNotes(w io.Writer, report quality.Report, inputs []string) {
	if len(report.Findings) == 0 {
		return
	}
	writeLine(w)
	writef(w, "Note: comparisons against this corpus may be noisy (reliability: %s).\n", report.Reliability())
	for _, finding := range report.Findings {
		writef(w, "- %s\n", finding.Summary)
		if finding.Action != "" {
			writef(w, "  → %s\n", finding.Action)
		}
	}
	writef(w, "Run 'omokage doctor %s' for the full report.\n", strings.Join(inputs, " "))
}

type doctorJSON struct {
	DocumentCount             int                  `json:"document_count"`
	SentenceCount             int                  `json:"sentence_count"`
	CharacterCount            int                  `json:"character_count"`
	AverageDocumentCharacters int                  `json:"average_document_characters"`
	Reliability               string               `json:"reliability"`
	Findings                  []qualityFindingJSON `json:"findings"`
}

// qualityFindingJSON is one corpus-quality finding in JSON form, shared by
// `doctor --format json` and the `show --format json` corpus summary.
type qualityFindingJSON struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail,omitempty"`
	Action   string `json:"action"`
}

// renderDoctorJSON emits the corpus report as one machine-readable JSON object
// for a tool or an LLM to consume.
func renderDoctorJSON(w io.Writer, report quality.Report) error {
	payload := doctorJSON{
		DocumentCount:             report.DocumentCount,
		SentenceCount:             report.SentenceCount,
		CharacterCount:            report.CharacterCount,
		AverageDocumentCharacters: report.AverageDocumentCharacters,
		Reliability:               report.Reliability(),
		Findings:                  toQualityFindingJSON(report.Findings),
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

// toQualityFindingJSON converts findings into their JSON shape, always returning
// a non-nil slice so the JSON shows an empty array rather than null.
func toQualityFindingJSON(findings []quality.Finding) []qualityFindingJSON {
	out := make([]qualityFindingJSON, 0, len(findings))
	for _, finding := range findings {
		out = append(out, qualityFindingJSON{
			Code:     finding.Code,
			Severity: finding.Severity.String(),
			Summary:  finding.Summary,
			Detail:   finding.Detail,
			Action:   finding.Action,
		})
	}
	return out
}
