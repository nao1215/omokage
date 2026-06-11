package feature

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// Version identifies the feature-definition generation a profile was trained
// with. It is bumped whenever an existing feature's measurement changes (not when
// a new feature is merely added), so a check can warn that a profile trained by
// an older omokage is being compared with newly-defined target metrics. Version 1
// was the original whitespace/heuristic definitions; version 2 introduced the
// kagome morphological measurement of conjunction frequency, the polite/plain
// register, and the Japanese function-word fingerprint.
const Version = 2

type Metrics struct {
	AverageSentenceLength    float64
	SentenceLengthVariance   float64
	PunctuationFrequency     float64
	NewlineFrequency         float64
	BulletRatio              float64
	ConjunctionFrequency     float64
	KanjiRatio               float64
	HiraganaRatio            float64
	KatakanaRatio            float64
	ParagraphLengthVariance  float64
	MarkdownStructureDensity float64
	PoliteEndingRatio        float64
	PlainEndingRatio         float64
	// TypeTokenRatio is the lemma-based vocabulary richness (P5): distinct content
	// lemmas over content morphemes, for Japanese prose (0 otherwise).
	TypeTokenRatio float64
	SentenceCount  int
	CharacterCount int
	// LexicalFrequencies holds the per-function-word relative frequency vector
	// keyed by LexicalVocabulary(). On a Distribution's Mean/StdDev it carries
	// the per-word mean and standard deviation across the corpus.
	LexicalFrequencies map[string]float64
	// CharNgrams holds the relative frequency of character bigrams. For a single
	// document it covers every bigram that occurs; on a Distribution's Mean/StdDev
	// it is restricted to the author's most frequent bigrams (see aggregateCharNgrams),
	// which form a language-independent stylometric fingerprint.
	CharNgrams map[string]float64
	// POSNgrams holds the relative frequency of part-of-speech bigrams and trigrams
	// (P4), capturing how a Japanese author builds sentences independently of
	// vocabulary. It is populated only for Japanese prose; English documents leave
	// it empty. On a Distribution it is restricted to the author's most frequent
	// POS n-grams (see aggregatePOSNgrams).
	POSNgrams map[string]float64
}

// Distribution captures the per-document spread of a corpus. An author profile
// is modelled as the mean and standard deviation of each feature across the
// individual documents, which lets a target be scored by how many standard
// deviations it sits from the author's own writing rather than from a single
// concatenated point estimate.
type Distribution struct {
	Mean           Metrics
	StdDev         Metrics
	DocumentCount  int
	SentenceCount  int
	CharacterCount int
}

func CollectFiles(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if !supportedExtension(root) {
			return nil, nil
		}
		return []string{filepath.Clean(root)}, nil
	}

	files := make([]string, 0, 16)
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if supportedExtension(path) {
			files = append(files, filepath.Clean(path))
		}
		return nil
	}); err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// IsSupportedFile reports whether path names a file omokage can learn from: a
// .md or .txt file. It lets callers reject an explicitly passed file with the
// wrong extension before any collection happens, using the same rule CollectFiles
// applies when walking a directory.
func IsSupportedFile(path string) bool {
	return supportedExtension(path)
}

func ExtractFile(path string) (Metrics, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Metrics{}, err
	}
	return ExtractText(string(data)), nil
}

// Segment is one localizable unit of a document — a paragraph — paired with the
// features extracted from it. It lets the explain path point at the specific
// paragraph that drifts most, rather than only reporting whole-document drift.
type Segment struct {
	Index   int
	Kind    string
	Text    string
	Metrics Metrics
}

// ExtractFileWithSegments reads a file once and returns both its whole-document
// metrics and its per-paragraph segments. It backs `check --explain`/`--format
// json`; the plain `check` path stays on ExtractFile so it does the lighter
// whole-document work only.
func ExtractFileWithSegments(path string) (Metrics, []Segment, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Metrics{}, nil, err
	}
	text := string(data)
	return ExtractText(text), ExtractSegments(text), nil
}

// ExtractSegments splits a document into prose paragraphs and extracts the
// features of each, dropping whitespace-only paragraphs and ones that are not
// running prose — a lone heading, a bullet or table block, a blockquote. Those
// are layout, not sentences a writer edits for voice, so localizing drift to them
// is noise; the whole-document score still measures them. The 1-based Index is
// dense over the prose paragraphs that survive, so a report can name "paragraph
// #N".
func ExtractSegments(text string) []Segment {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	// Strip code and HTML on the whole document before splitting into paragraphs.
	// Splitting first would break a fenced block that contains a blank line (mermaid
	// diagrams and shell sessions routinely do): each fragment would lose its
	// opening fence and be measured as prose, so the report would point at a diagram
	// or CLI transcript as a drifting "paragraph". Cleaning first keeps segment
	// metrics and excerpts consistent with the whole-document measurement.
	prose := StripNonProse(normalized)
	paragraphs := splitParagraphs(prose)
	segments := make([]Segment, 0, len(paragraphs))
	index := 0
	for _, paragraph := range paragraphs {
		if !looksLikeProse(paragraph) {
			continue
		}
		metrics := ExtractText(paragraph)
		if metrics.CharacterCount == 0 {
			continue
		}
		index++
		segments = append(segments, Segment{
			Index:   index,
			Kind:    "paragraph",
			Text:    paragraph,
			Metrics: metrics,
		})
	}
	return segments
}

// looksLikeProse reports whether a paragraph is running prose worth localizing
// drift to, as opposed to layout. It rejects a paragraph whose lines are mostly
// Markdown structure (a heading, a bullet or ordered list, a table, a blockquote)
// and one with no sentence terminator at all (a bare heading or label). A normal
// prose paragraph — sentences, few or no structure lines — passes. This narrows
// only the paragraph-localization in the explain/JSON output; the whole-document
// score still measures every paragraph.
func looksLikeProse(paragraph string) bool {
	var nonEmpty, structure int
	for line := range strings.SplitSeq(paragraph, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmpty++
		if isMarkdownStructureLine(trimmed) {
			structure++
		}
	}
	if nonEmpty == 0 {
		return false
	}
	// Mostly structure lines (a heading, a bullet/table/quote block) is layout, not
	// prose to edit for voice.
	if structure*2 > nonEmpty {
		return false
	}
	// Prose carries at least one sentence; a heading or a label line carries none.
	return hasSentenceTerminator(paragraph)
}

// hasSentenceTerminator reports whether the text contains any sentence-ending
// punctuation, in either script. It is the cheap "is there a sentence here at
// all" test looksLikeProse uses to drop heading- and label-only paragraphs.
func hasSentenceTerminator(text string) bool {
	return strings.ContainsAny(text, ".!?。！？")
}

// Document pairs a corpus file with the features extracted from it. It lets a
// caller inspect the corpus per document (for a quality check that flags short
// or out-of-place files) while still aggregating the same metrics into the
// learned Distribution, so the two views never disagree.
type Document struct {
	Path    string
	Metrics Metrics
}

// ExtractCorpus extracts features from each file independently and aggregates
// them into a Distribution describing the mean and standard deviation of every
// feature across the corpus.
func ExtractCorpus(paths []string) (Distribution, error) {
	dist, _, err := ExtractCorpusDocuments(paths)
	return dist, err
}

// ExtractCorpusDocuments is ExtractCorpus plus the per-document metrics that fed
// the aggregate. It reads each file once and returns both the Distribution and
// the surviving documents (empty or whitespace-only files are dropped from both,
// exactly as ExtractCorpus drops them from the aggregate), so a caller can assess
// the corpus document by document without a second pass over the files.
func ExtractCorpusDocuments(paths []string) (Distribution, []Document, error) {
	docs := make([]Document, 0, len(paths))
	perDoc := make([]Metrics, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return Distribution{}, nil, fmt.Errorf("read %s: %w", path, err)
		}
		metrics := ExtractText(string(data))
		// Empty or whitespace-only files contribute an all-zero feature vector
		// that would drag the mean down and inflate the standard deviation, so
		// they are excluded from the learned distribution and from the per-document
		// view.
		if metrics.CharacterCount == 0 {
			continue
		}
		perDoc = append(perDoc, metrics)
		docs = append(docs, Document{Path: filepath.Clean(path), Metrics: metrics})
	}
	return Aggregate(perDoc), docs, nil
}

// Aggregate reduces per-document metrics into a Distribution. Feature means and
// population standard deviations are computed per field; sentence and character
// counts are summed across the corpus.
func Aggregate(perDoc []Metrics) Distribution {
	dist := Distribution{DocumentCount: len(perDoc)}
	if len(perDoc) == 0 {
		return dist
	}

	for _, m := range perDoc {
		dist.Mean.AverageSentenceLength += m.AverageSentenceLength
		dist.Mean.SentenceLengthVariance += m.SentenceLengthVariance
		dist.Mean.PunctuationFrequency += m.PunctuationFrequency
		dist.Mean.NewlineFrequency += m.NewlineFrequency
		dist.Mean.BulletRatio += m.BulletRatio
		dist.Mean.ConjunctionFrequency += m.ConjunctionFrequency
		dist.Mean.KanjiRatio += m.KanjiRatio
		dist.Mean.HiraganaRatio += m.HiraganaRatio
		dist.Mean.KatakanaRatio += m.KatakanaRatio
		dist.Mean.ParagraphLengthVariance += m.ParagraphLengthVariance
		dist.Mean.MarkdownStructureDensity += m.MarkdownStructureDensity
		dist.Mean.PoliteEndingRatio += m.PoliteEndingRatio
		dist.Mean.PlainEndingRatio += m.PlainEndingRatio
		dist.Mean.TypeTokenRatio += m.TypeTokenRatio
		dist.SentenceCount += m.SentenceCount
		dist.CharacterCount += m.CharacterCount
	}

	n := float64(len(perDoc))
	dist.Mean.AverageSentenceLength /= n
	dist.Mean.SentenceLengthVariance /= n
	dist.Mean.PunctuationFrequency /= n
	dist.Mean.NewlineFrequency /= n
	dist.Mean.BulletRatio /= n
	dist.Mean.ConjunctionFrequency /= n
	dist.Mean.KanjiRatio /= n
	dist.Mean.HiraganaRatio /= n
	dist.Mean.KatakanaRatio /= n
	dist.Mean.ParagraphLengthVariance /= n
	dist.Mean.MarkdownStructureDensity /= n
	dist.Mean.PoliteEndingRatio /= n
	dist.Mean.PlainEndingRatio /= n
	dist.Mean.TypeTokenRatio /= n

	for _, m := range perDoc {
		accumulateSquaredError(&dist.StdDev, m, dist.Mean)
	}
	finalizeStdDev(&dist.StdDev, n)

	aggregateLexical(&dist, perDoc, n)
	aggregateCharNgrams(&dist, perDoc, n)
	aggregatePOSNgrams(&dist, perDoc, n)
	return dist
}

// aggregateLexical computes the per-word mean and population standard deviation
// of the lexical frequency vector across the corpus. The vocabulary is fixed,
// so a word absent from a document simply contributes a zero for that document.
func aggregateLexical(dist *Distribution, perDoc []Metrics, n float64) {
	meanVec := make(map[string]float64, len(lexicalVocabulary))
	for _, m := range perDoc {
		for _, word := range lexicalVocabulary {
			meanVec[word] += m.LexicalFrequencies[word]
		}
	}
	for word := range meanVec {
		meanVec[word] /= n
	}

	stdVec := make(map[string]float64, len(lexicalVocabulary))
	for _, m := range perDoc {
		for _, word := range lexicalVocabulary {
			stdVec[word] += square(m.LexicalFrequencies[word] - meanVec[word])
		}
	}
	for word := range stdVec {
		stdVec[word] = math.Sqrt(stdVec[word] / n)
	}

	dist.Mean.LexicalFrequencies = meanVec
	dist.StdDev.LexicalFrequencies = stdVec
}

func accumulateSquaredError(acc *Metrics, sample Metrics, mean Metrics) {
	acc.AverageSentenceLength += square(sample.AverageSentenceLength - mean.AverageSentenceLength)
	acc.SentenceLengthVariance += square(sample.SentenceLengthVariance - mean.SentenceLengthVariance)
	acc.PunctuationFrequency += square(sample.PunctuationFrequency - mean.PunctuationFrequency)
	acc.NewlineFrequency += square(sample.NewlineFrequency - mean.NewlineFrequency)
	acc.BulletRatio += square(sample.BulletRatio - mean.BulletRatio)
	acc.ConjunctionFrequency += square(sample.ConjunctionFrequency - mean.ConjunctionFrequency)
	acc.KanjiRatio += square(sample.KanjiRatio - mean.KanjiRatio)
	acc.HiraganaRatio += square(sample.HiraganaRatio - mean.HiraganaRatio)
	acc.KatakanaRatio += square(sample.KatakanaRatio - mean.KatakanaRatio)
	acc.ParagraphLengthVariance += square(sample.ParagraphLengthVariance - mean.ParagraphLengthVariance)
	acc.MarkdownStructureDensity += square(sample.MarkdownStructureDensity - mean.MarkdownStructureDensity)
	acc.PoliteEndingRatio += square(sample.PoliteEndingRatio - mean.PoliteEndingRatio)
	acc.PlainEndingRatio += square(sample.PlainEndingRatio - mean.PlainEndingRatio)
	acc.TypeTokenRatio += square(sample.TypeTokenRatio - mean.TypeTokenRatio)
}

func finalizeStdDev(acc *Metrics, n float64) {
	acc.AverageSentenceLength = math.Sqrt(acc.AverageSentenceLength / n)
	acc.SentenceLengthVariance = math.Sqrt(acc.SentenceLengthVariance / n)
	acc.PunctuationFrequency = math.Sqrt(acc.PunctuationFrequency / n)
	acc.NewlineFrequency = math.Sqrt(acc.NewlineFrequency / n)
	acc.BulletRatio = math.Sqrt(acc.BulletRatio / n)
	acc.ConjunctionFrequency = math.Sqrt(acc.ConjunctionFrequency / n)
	acc.KanjiRatio = math.Sqrt(acc.KanjiRatio / n)
	acc.HiraganaRatio = math.Sqrt(acc.HiraganaRatio / n)
	acc.KatakanaRatio = math.Sqrt(acc.KatakanaRatio / n)
	acc.ParagraphLengthVariance = math.Sqrt(acc.ParagraphLengthVariance / n)
	acc.MarkdownStructureDensity = math.Sqrt(acc.MarkdownStructureDensity / n)
	acc.PoliteEndingRatio = math.Sqrt(acc.PoliteEndingRatio / n)
	acc.PlainEndingRatio = math.Sqrt(acc.PlainEndingRatio / n)
	acc.TypeTokenRatio = math.Sqrt(acc.TypeTokenRatio / n)
}

func square(value float64) float64 {
	return value * value
}

func ExtractText(text string) Metrics {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	// Every feature is measured on the code-stripped prose, not the raw text, so
	// the score reflects the author's natural-language habits rather than the
	// shared vocabulary and layout of code samples. Authorship features (function
	// words, character n-grams) need this most — code otherwise dominates technical
	// posts and masks the difference between two technical authors — but the
	// structural features (sentence length, punctuation, Markdown density) are
	// measured on the same prose too, so adding a fenced code block to a draft no
	// longer manufactures false drift. This matches the documented promise that
	// code blocks are removed before the features are measured. HTML tags are
	// stripped too: raw HTML embedded in Markdown is layout, not prose.
	prose := StripNonProse(normalized)
	sentences := splitSentences(prose)
	sentenceLengths := make([]float64, 0, len(sentences))
	for _, sentence := range sentences {
		sentenceLengths = append(sentenceLengths, float64(contentRuneCount(sentence)))
	}

	paragraphs := splitParagraphs(prose)
	paragraphLengths := make([]float64, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		paragraphLengths = append(paragraphLengths, float64(contentRuneCount(paragraph)))
	}

	lines := strings.Split(prose, "\n")
	var bulletLines int
	var structureLines int
	var nonEmptyLines int
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmptyLines++
		if isBulletLine(trimmed) {
			bulletLines++
		}
		if isMarkdownStructureLine(trimmed) {
			structureLines++
		}
	}

	characterCount := contentRuneCount(prose)
	if characterCount == 0 {
		return Metrics{}
	}

	newlineCount := strings.Count(prose, "\n")
	punctuationCount := punctuationRuneCount(prose)
	conjunctionCount, tokenCount := conjunctionStats(prose)
	kanjiCount, hiraganaCount, katakanaCount := scriptCounts(prose)
	scriptTotal := max(kanjiCount+hiraganaCount+katakanaCount, 1)
	politeCount, plainCount := sentenceEndingStats(prose)
	// Tokenize Japanese prose once and reuse the morphemes across every feature
	// that whitespace tokenization cannot measure correctly: conjunction frequency
	// (P1), the polite/plain register (P3), the function-word fingerprint (P2), and
	// the POS n-gram fingerprint (P4). English prose keeps the language-neutral path
	// and jpTokens stays nil.
	var jpTokens []jpToken
	var typeTokenRatio float64
	if isJapanese(prose) {
		jpTokens = tokenizeJapanese(prose)
	}
	if len(jpTokens) > 0 {
		conjunctionCount, tokenCount = conjunctionStatsJP(jpTokens)
		politeCount, plainCount = registerStatsJP(jpTokens)
		typeTokenRatio = typeTokenRatioJP(jpTokens)
	}
	// Normalize sentence-ending forms by Japanese terminators (。！？) rather than
	// len(sentences): the latter also splits on Latin "." and would over-count
	// sentences in technical prose (version numbers, decimals), diluting the ratio.
	japaneseSentences := max(strings.Count(prose, "。")+strings.Count(prose, "！")+strings.Count(prose, "？"), 1)

	return Metrics{
		AverageSentenceLength:    mean(sentenceLengths),
		SentenceLengthVariance:   variance(sentenceLengths),
		PunctuationFrequency:     clamp01(float64(punctuationCount) / float64(characterCount)),
		NewlineFrequency:         clamp01(float64(newlineCount) / float64(characterCount)),
		BulletRatio:              ratio(bulletLines, max(nonEmptyLines, 1)),
		ConjunctionFrequency:     ratio(conjunctionCount, max(tokenCount, 1)),
		KanjiRatio:               ratio(kanjiCount, scriptTotal),
		HiraganaRatio:            ratio(hiraganaCount, scriptTotal),
		KatakanaRatio:            ratio(katakanaCount, scriptTotal),
		ParagraphLengthVariance:  variance(paragraphLengths),
		MarkdownStructureDensity: ratio(structureLines, max(nonEmptyLines, 1)),
		PoliteEndingRatio:        ratio(politeCount, japaneseSentences),
		PlainEndingRatio:         ratio(plainCount, japaneseSentences),
		TypeTokenRatio:           typeTokenRatio,
		SentenceCount:            len(sentences),
		CharacterCount:           characterCount,
		LexicalFrequencies:       lexicalFrequencies(prose, jpTokens),
		CharNgrams:               charBigrams(prose),
		POSNgrams:                posNgrams(jpTokens),
	}
}

func supportedExtension(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	return extension == ".md" || extension == ".txt"
}

func splitSentences(text string) []string {
	replacer := strings.NewReplacer("!", ". ", "?", ". ", "。", ". ", "！", ". ", "？", ". ")
	runes := []rune(replacer.Replace(text))
	sentences := make([]string, 0)
	start := 0
	for i, r := range runes {
		if r != '.' {
			continue
		}
		// A period ends a sentence only when it is the last character or is followed
		// by whitespace. A period wedged between non-space characters — version
		// numbers (1.2.3), domains (example.com), decimals (3.14), abbreviations —
		// is part of a token, not a boundary, so it must not split the sentence.
		if i+1 < len(runes) && !unicode.IsSpace(runes[i+1]) {
			continue
		}
		if trimmed := strings.TrimSpace(string(runes[start : i+1])); trimmed != "" {
			sentences = append(sentences, trimmed)
		}
		start = i + 1
	}
	if start < len(runes) {
		if trimmed := strings.TrimSpace(string(runes[start:])); trimmed != "" {
			sentences = append(sentences, trimmed)
		}
	}
	if len(sentences) == 0 {
		if trimmed := strings.TrimSpace(text); trimmed != "" {
			return []string{trimmed}
		}
	}
	return sentences
}

func splitParagraphs(text string) []string {
	raw := strings.Split(text, "\n\n")
	paragraphs := make([]string, 0, len(raw))
	for _, paragraph := range raw {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed != "" {
			paragraphs = append(paragraphs, trimmed)
		}
	}
	if len(paragraphs) == 0 {
		trimmed := strings.TrimSpace(text)
		if trimmed != "" {
			return []string{trimmed}
		}
	}
	return paragraphs
}

func contentRuneCount(text string) int {
	var total int
	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		total++
	}
	return total
}

func punctuationRuneCount(text string) int {
	var total int
	for _, r := range text {
		if unicode.IsPunct(r) || strings.ContainsRune("、。！？「」（）【】", r) {
			total++
		}
	}
	return total
}

func conjunctionStats(text string) (int, int) {
	englishConjunctions := map[string]struct{}{
		"and": {}, "but": {}, "or": {}, "so": {}, "yet": {}, "for": {}, "nor": {},
		"because": {}, "although": {}, "however": {}, "therefore": {}, "thus": {},
		"meanwhile": {}, "moreover": {}, "nevertheless": {},
	}
	japaneseConjunctions := []string{
		"そして", "しかし", "だから", "なので", "また", "さらに", "ただし", "一方", "つまり", "ところが",
	}

	words := splitWords(strings.ToLower(text))
	count := 0
	for _, word := range words {
		if _, ok := englishConjunctions[word]; ok {
			count++
		}
	}
	for _, token := range japaneseConjunctions {
		count += strings.Count(text, token)
	}
	return count, len(words)
}

func splitWords(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// sentenceEndingStats counts Japanese sentence-ending forms. The contrast
// between the polite register (敬体: です・ます) and the plain register
// (常体: だ・である・verb/adjective stop) is one of the strongest stylistic
// markers in Japanese writing, and a sharp shift between the two is exactly the
// kind of drift omokage aims to surface.
//
// Each clause terminated by a Japanese full stop (。！？) is classified by the
// predicate that closes it rather than by counting fixed substrings anywhere in
// the text. That is what lets the plain register be detected at all: 常体 is not
// a short list of words (である・だった) but the open class of plain-form
// predicates — verbs (する・行く・読んだ), i-adjectives (高い・ない), and the
// copula (だ) — which a substring scan cannot enumerate.
func sentenceEndingStats(text string) (polite int, plain int) {
	for _, clause := range japaneseClauses(text) {
		switch classifyEnding(clause) {
		case endingPolite:
			polite++
		case endingPlain:
			plain++
		}
	}
	return polite, plain
}

const (
	endingNone = iota
	endingPolite
	endingPlain
)

// japaneseClauses splits text into the substrings that each end at a Japanese
// sentence terminator (。！？). A trailing fragment with no terminator is not a
// completed sentence and is dropped, so the count matches the 。！？ denominator
// ExtractText normalizes by.
func japaneseClauses(text string) []string {
	clauses := make([]string, 0)
	var current strings.Builder
	for _, r := range text {
		if r == '。' || r == '！' || r == '？' {
			if clause := strings.TrimSpace(current.String()); clause != "" {
				clauses = append(clauses, clause)
			}
			current.Reset()
			continue
		}
		current.WriteRune(r)
	}
	return clauses
}

// classifyEnding decides whether a clause ends in the polite or the plain
// register, or in neither (English, a noun stop/体言止め, a bare particle). The
// polite check runs first because polite auxiliaries (です・ます) end in kana
// that the plain heuristic would otherwise claim.
func classifyEnding(clause string) int {
	core := trimSentenceTail(clause)
	if core == "" {
		return endingNone
	}
	if hasPoliteEnding(core) {
		return endingPolite
	}
	if hasPlainEnding(core) {
		return endingPlain
	}
	return endingNone
}

// sentenceTailRunes are closing quotes/brackets and interjective sentence-final
// particles that sit after the predicate. Trimming them exposes the verb,
// adjective, or copula that carries the register (行きますか → 行きます,
// するなよ → するな).
const sentenceTailRunes = "かねよわなぞぜさのっ」』）)】"

func trimSentenceTail(clause string) string {
	runes := []rune(clause)
	for len(runes) > 0 && strings.ContainsRune(sentenceTailRunes, runes[len(runes)-1]) {
		runes = runes[:len(runes)-1]
	}
	return string(runes)
}

func hasPoliteEnding(core string) bool {
	politeForms := []string{"です", "ます", "でした", "ました", "ません", "でしょう", "ましょう"}
	for _, form := range politeForms {
		if strings.HasSuffix(core, form) {
			return true
		}
	}
	return false
}

// plainEndingKana are the kana that close a plain-form predicate: the godan and
// ichidan verb endings (る・く・す…) and the i-adjective い. A clause that is not
// polite and ends in one of these is treated as 常体.
const plainEndingKana = "うくぐすつづぬふぶむゆるい"

func hasPlainEnding(core string) bool {
	plainForms := []string{"である", "であった", "だった", "だろう", "なかった", "ない", "た", "だ"}
	for _, form := range plainForms {
		if strings.HasSuffix(core, form) {
			return true
		}
	}
	runes := []rune(core)
	return strings.ContainsRune(plainEndingKana, runes[len(runes)-1])
}

func scriptCounts(text string) (kanji int, hiragana int, katakana int) {
	for _, r := range text {
		switch {
		case unicode.In(r, unicode.Han):
			kanji++
		case unicode.In(r, unicode.Hiragana):
			hiragana++
		case unicode.In(r, unicode.Katakana):
			katakana++
		}
	}
	return kanji, hiragana, katakana
}

func isBulletLine(line string) bool {
	return strings.HasPrefix(line, "- ") ||
		strings.HasPrefix(line, "* ") ||
		strings.HasPrefix(line, "+ ") ||
		isOrderedList(line)
}

func isOrderedList(line string) bool {
	index := 0
	for index < len(line) && line[index] >= '0' && line[index] <= '9' {
		index++
	}
	if index == 0 || index+1 >= len(line) {
		return false
	}
	return line[index] == '.' && line[index+1] == ' '
}

func isMarkdownStructureLine(line string) bool {
	return strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, ">") ||
		strings.HasPrefix(line, "```") ||
		strings.HasPrefix(line, "|") ||
		isBulletLine(line)
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func variance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	avg := mean(values)
	sum := 0.0
	for _, value := range values {
		diff := value - avg
		sum += diff * diff
	}
	return sum / float64(len(values))
}

func ratio(numerator int, denominator int) float64 {
	if denominator <= 0 {
		return 0
	}
	return clamp01(float64(numerator) / float64(denominator))
}

func clamp01(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
