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
	SentenceCount            int
	CharacterCount           int
	// LexicalFrequencies holds the per-function-word relative frequency vector
	// keyed by LexicalVocabulary(). On a Distribution's Mean/StdDev it carries
	// the per-word mean and standard deviation across the corpus.
	LexicalFrequencies map[string]float64
	// CharNgrams holds the relative frequency of character bigrams. For a single
	// document it covers every bigram that occurs; on a Distribution's Mean/StdDev
	// it is restricted to the author's most frequent bigrams (see aggregateCharNgrams),
	// which form a language-independent stylometric fingerprint.
	CharNgrams map[string]float64
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

func ExtractFile(path string) (Metrics, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Metrics{}, err
	}
	return ExtractText(string(data)), nil
}

// ExtractCorpus extracts features from each file independently and aggregates
// them into a Distribution describing the mean and standard deviation of every
// feature across the corpus.
func ExtractCorpus(paths []string) (Distribution, error) {
	perDoc := make([]Metrics, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return Distribution{}, fmt.Errorf("read %s: %w", path, err)
		}
		metrics := ExtractText(string(data))
		// Empty or whitespace-only files contribute an all-zero feature vector
		// that would drag the mean down and inflate the standard deviation, so
		// they are excluded from the learned distribution.
		if metrics.CharacterCount == 0 {
			continue
		}
		perDoc = append(perDoc, metrics)
	}
	return Aggregate(perDoc), nil
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

	for _, m := range perDoc {
		accumulateSquaredError(&dist.StdDev, m, dist.Mean)
	}
	finalizeStdDev(&dist.StdDev, n)

	aggregateLexical(&dist, perDoc, n)
	aggregateCharNgrams(&dist, perDoc, n)
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
}

func square(value float64) float64 {
	return value * value
}

func ExtractText(text string) Metrics {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	// prose drops code so the authorship features (function words, character
	// n-grams) measure the author's natural-language habits rather than the
	// shared vocabulary of code samples — which otherwise dominates technical
	// posts and masks the difference between two technical authors. The
	// structural features below still see the original text.
	prose := stripCode(normalized)
	sentences := splitSentences(normalized)
	sentenceLengths := make([]float64, 0, len(sentences))
	for _, sentence := range sentences {
		sentenceLengths = append(sentenceLengths, float64(contentRuneCount(sentence)))
	}

	paragraphs := splitParagraphs(normalized)
	paragraphLengths := make([]float64, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		paragraphLengths = append(paragraphLengths, float64(contentRuneCount(paragraph)))
	}

	lines := strings.Split(normalized, "\n")
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

	characterCount := contentRuneCount(normalized)
	if characterCount == 0 {
		return Metrics{}
	}

	newlineCount := strings.Count(normalized, "\n")
	punctuationCount := punctuationRuneCount(normalized)
	conjunctionCount, tokenCount := conjunctionStats(normalized)
	kanjiCount, hiraganaCount, katakanaCount := scriptCounts(normalized)
	scriptTotal := max(kanjiCount+hiraganaCount+katakanaCount, 1)
	politeCount, plainCount := sentenceEndingStats(normalized)
	// Normalize sentence-ending forms by Japanese terminators (。！？) rather than
	// len(sentences): the latter also splits on Latin "." and would over-count
	// sentences in technical prose (version numbers, decimals), diluting the ratio.
	japaneseSentences := max(strings.Count(normalized, "。")+strings.Count(normalized, "！")+strings.Count(normalized, "？"), 1)

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
		SentenceCount:            len(sentences),
		CharacterCount:           characterCount,
		LexicalFrequencies:       lexicalFrequencies(prose),
		CharNgrams:               charBigrams(prose),
	}
}

func supportedExtension(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	return extension == ".md" || extension == ".txt"
}

func splitSentences(text string) []string {
	replacer := strings.NewReplacer("!", ". ", "?", ". ", "。", ". ", "！", ". ", "？", ". ")
	normalized := replacer.Replace(text)
	parts := strings.Split(normalized, ".")
	sentences := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			sentences = append(sentences, trimmed)
		}
	}
	if len(sentences) == 0 {
		trimmed := strings.TrimSpace(text)
		if trimmed != "" {
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
// (常体: だ・である) is one of the strongest stylistic markers in Japanese
// writing, and a sharp shift between the two is exactly the kind of drift omokage
// aims to surface.
func sentenceEndingStats(text string) (polite int, plain int) {
	politeForms := []string{"ます", "です", "ました", "でした", "ません"}
	plainForms := []string{"である", "だった", "だが"}
	for _, form := range politeForms {
		polite += strings.Count(text, form)
	}
	for _, form := range plainForms {
		plain += strings.Count(text, form)
	}
	return polite, plain
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
