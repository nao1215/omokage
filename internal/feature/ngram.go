package feature

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// charNgramVocabularySize caps how many of an author's most frequent character
// n-grams are kept in the profile. Character n-grams are a strong, language
// independent authorship signal; keeping the most frequent ones captures the
// stable core of an author's habits while bounding the profile size.
const charNgramVocabularySize = 400

// charBigrams returns the relative frequency of every character bigram and
// trigram in the text. Newlines and runs of whitespace are collapsed to a single
// space so that layout does not dominate the counts. Bigrams capture broad
// habits while trigrams capture more distinctive sequences (especially in
// space-delimited English); each order is normalized by its own total so the two
// live comparably in one map. The same routine runs on training documents and on
// a check target; the profile later narrows the comparison to its own frequent
// n-grams.
func charBigrams(text string) map[string]float64 {
	runes := normalizeForNgram(text)
	freq := make(map[string]float64, len(runes)*2)
	addNgrams(freq, runes, 2)
	addNgrams(freq, runes, 3)
	return freq
}

// addNgrams accumulates the relative frequency of every n-gram of the given
// order into freq, normalized by the number of n-grams of that order.
func addNgrams(freq map[string]float64, runes []rune, order int) {
	if len(runes) < order {
		return
	}
	counts := make(map[string]int, len(runes))
	total := 0
	for i := 0; i+order <= len(runes); i++ {
		counts[string(runes[i:i+order])]++
		total++
	}
	// total is at least 1 here because len(runes) >= order was checked above.
	for ngram, count := range counts {
		freq[ngram] = float64(count) / float64(total)
	}
}

// inlineCodePattern matches an inline code span delimited by backticks.
var inlineCodePattern = regexp.MustCompile("`[^`]*`")

// Placeholders that topic-heavy technical tokens collapse to before the lexical
// and character n-gram fingerprints are measured. They are fixed strings, so the
// same identifier-shaped token always contributes the same n-grams regardless of
// which product, repository, or version a post happens to discuss. See
// MaskTechnicalTopicTokens.
const (
	placeholderRepo    = "<REPO>"
	placeholderIdent   = "<IDENT>"
	placeholderVersion = "<VERSION>"
	placeholderAcronym = "<ACRONYM>"
	placeholderFile    = "<FILE>"
	placeholderNumber  = "<NUMBER>"
)

// fileExtensions is the closed set of file extensions that mark a token as a
// filename/extension reference (data.csv, .parquet, main.go). It is deliberately
// small and concrete: a token whose final dotted segment is one of these is topic
// noise (which data format a post discusses), not voice.
var fileExtensions = map[string]bool{
	"csv": true, "tsv": true, "ltsv": true, "parquet": true, "xlsx": true,
	"xls": true, "json": true, "yaml": true, "yml": true, "toml": true,
	"md": true, "txt": true, "go": true, "rs": true, "py": true, "js": true,
	"ts": true, "jsx": true, "tsx": true, "java": true, "rb": true, "php": true,
	"sh": true, "sql": true, "html": true, "css": true, "png": true, "jpg": true,
	"jpeg": true, "gif": true, "svg": true, "pdf": true, "zip": true, "tar": true,
	"gz": true, "mod": true, "sum": true, "lock": true, "ini": true, "xml": true,
}

// technicalTokenPattern matches a maximal run of ASCII technical characters that
// begins and ends with an alphanumeric (so a trailing English sentence period or
// a leading bullet dash is never absorbed) plus an optional leading dot, which
// lets a bare extension reference (.csv) be caught. Japanese script characters
// are outside the class, so the kana/kanji that carry the Japanese stylistic
// signal are never touched — only the embedded Latin tokens are.
var technicalTokenPattern = regexp.MustCompile(`\.?[A-Za-z0-9](?:[A-Za-z0-9_./\-]*[A-Za-z0-9])?`)

// versionTokenPattern matches a version-like numeric string: an optional leading
// v then digits, optionally followed by dotted numeric segments (v1, v1.2.3,
// 1.25, 2025.12). A bare integer is intentionally excluded so it falls through to
// the <NUMBER> class instead.
var versionTokenPattern = regexp.MustCompile(`^v\d+(?:\.\d+)*$|^\d+(?:\.\d+)+$`)

// MaskTechnicalTopicTokens replaces topic-heavy technical tokens — repository
// names, dotted identifiers, snake_case/kebab-case/CamelCase identifiers, version
// strings, uppercase acronyms, filenames, and bare numbers — with fixed
// placeholders, so the lexical and character n-gram fingerprints measure how an
// author writes rather than which product, library, or version a particular post
// is about. It is applied ONLY to the inputs of the lexical/char-ngram features;
// every scalar feature, the POS n-grams, and the script ratios keep measuring the
// original prose, so the masking cannot move sentence length, punctuation, register,
// or kana/kanji balance. Replacement (not deletion) preserves token positions, so
// surrounding n-grams stay aligned.
func MaskTechnicalTopicTokens(prose string) string {
	return technicalTokenPattern.ReplaceAllStringFunc(prose, maskTechnicalToken)
}

// maskTechnicalToken classifies one technical run and returns its placeholder, or
// the token unchanged when it is an ordinary word (a lowercase word, a single
// capitalized word, a pronoun) that carries no topic noise. A leading dot is split
// off first: it only belongs to the token when the token is a bare extension
// (.csv); otherwise it is re-emitted so an English sentence boundary survives.
func maskTechnicalToken(token string) string {
	lead := ""
	if strings.HasPrefix(token, ".") {
		lead = "."
		token = token[1:]
		// A contiguous bare extension (.csv, .parquet) is a filename reference even
		// though the stem is a single lowercase word the classifier would otherwise
		// leave alone. The leading dot is what marks it, so it is absorbed here.
		if fileExtensions[strings.ToLower(token)] {
			return placeholderFile
		}
	}
	placeholder, ok := classifyTechnicalToken(token)
	if !ok {
		return lead + token
	}
	// A bare ".csv" absorbs its leading dot; any other classification keeps the dot
	// as separate punctuation so it is not silently swallowed.
	if lead == "." && placeholder == placeholderFile {
		return placeholderFile
	}
	return lead + placeholder
}

// classifyTechnicalToken decides which placeholder (if any) a separator-trimmed
// token maps to. The order is significant: repository and version shapes are
// checked before the generic dotted-identifier rule, and the all-uppercase
// acronym rule is checked before the mixed-case CamelCase rule.
func classifyTechnicalToken(token string) (string, bool) {
	if token == "" {
		return "", false
	}
	if strings.Contains(token, "/") {
		return placeholderRepo, true
	}
	if versionTokenPattern.MatchString(token) {
		return placeholderVersion, true
	}
	if strings.Contains(token, ".") {
		segments := strings.Split(token, ".")
		last := strings.ToLower(segments[len(segments)-1])
		if fileExtensions[last] {
			return placeholderFile, true
		}
		return placeholderIdent, true
	}
	if strings.Contains(token, "_") {
		return placeholderIdent, true // snake_case
	}
	if strings.Contains(token, "-") {
		return placeholderIdent, true // kebab-case
	}
	if isAllDigits(token) {
		return placeholderNumber, true
	}
	if isUpperAcronym(token) {
		return placeholderAcronym, true
	}
	if isCamelCase(token) {
		return placeholderIdent, true
	}
	return "", false
}

// isAllDigits reports whether the token is a bare integer (the <NUMBER> class).
func isAllDigits(token string) bool {
	for _, r := range token {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// isUpperAcronym reports whether the token is an all-uppercase acronym of at
// least two characters (CSV, LTSV, WASM). Trailing digits are allowed (UTF8), but
// at least one uppercase letter must be present so a bare number is not caught
// here. A single uppercase letter (the English pronoun "I", an initial) is left
// alone.
func isUpperAcronym(token string) bool {
	if len([]rune(token)) < 2 {
		return false
	}
	hasLetter := false
	for _, r := range token {
		switch {
		case r >= 'A' && r <= 'Z':
			hasLetter = true
		case r >= '0' && r <= '9':
			// allowed
		default:
			return false
		}
	}
	return hasLetter
}

// isCamelCase reports whether the token is a mixed-case identifier with an
// INTERNAL uppercase letter (GitHub, PascalCase, iOS, DataFrame). The internal
// requirement is deliberate: a single leading capital (The, Foo) is
// indistinguishable from an ordinary sentence-initial word or a proper noun, so
// masking it would corrupt the English fingerprint; only the unmistakable
// internal "hump" of an identifier qualifies.
func isCamelCase(token string) bool {
	runes := []rune(token)
	hasLower := false
	internalUpper := false
	for i, r := range runes {
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			if i > 0 {
				internalUpper = true
			}
		}
	}
	return hasLower && internalUpper
}

// Patterns for the non-prose constructs that survive code stripping but are
// layout/markup, not natural language, so they must not feed the features.
var (
	// frontMatterPattern matches a leading YAML front-matter block (Hugo/Jekyll),
	// whose keys and paths (title, image, categories) are metadata, not prose.
	frontMatterPattern = regexp.MustCompile(`(?s)\A\s*---\n.*?\n---\n`)
	// imagePattern matches a Markdown image; it is dropped whole (alt text and the
	// path are not the author's prose).
	imagePattern = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	// linkPattern matches a Markdown link; the visible text is kept and the URL
	// dropped, so a link target never contributes URL fragments (https, com, jp).
	linkPattern = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	urlPattern  = regexp.MustCompile(`https?://\S+`)
	// htmlTagPattern matches an HTML tag starting with a letter (so a stray "<" or
	// "a < b" in prose is left alone). Raw HTML embedded in Markdown is layout.
	htmlTagPattern    = regexp.MustCompile(`</?[A-Za-z][^>]*>`)
	htmlEntityPattern = regexp.MustCompile(`&[a-zA-Z]+;`)
)

// StripNonProse removes everything that is not natural-language prose — YAML
// front matter, fenced and inline code, Markdown images, Markdown link URLs (the
// visible link text is kept), raw URLs, HTML tags, and HTML entities — and
// returns the remaining prose with CRLF normalized. It is the single chokepoint
// the feature extractor and term extraction both run through, so they measure the
// same prose. Code shares vocabulary and character sequences across authors and
// would otherwise drown out the natural-language signal in technical writing;
// markup and metadata would manufacture drift from layout rather than voice.
func StripNonProse(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = frontMatterPattern.ReplaceAllString(text, "")
	text = stripCode(text)
	text = imagePattern.ReplaceAllString(text, " ")
	text = linkPattern.ReplaceAllString(text, "$1")
	text = urlPattern.ReplaceAllString(text, " ")
	text = htmlTagPattern.ReplaceAllString(text, " ")
	text = htmlEntityPattern.ReplaceAllString(text, " ")
	return text
}

// stripCode removes fenced code blocks and inline code spans from the text. Both
// CommonMark fence markers are recognized — backtick (``` … ```) and tilde
// (~~~ … ~~~) — and a block is closed only by its own marker, so a tilde line
// inside a backtick block (or vice versa) is treated as content, not a boundary.
func stripCode(text string) string {
	lines := strings.Split(text, "\n")
	kept := make([]string, 0, len(lines))
	fence := "" // marker that opened the current block ("```"/"~~~"), "" when outside
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if fence == "" {
			if marker := fenceMarker(trimmed); marker != "" {
				fence = marker
				continue
			}
		} else if strings.HasPrefix(trimmed, fence) {
			fence = ""
			continue
		}
		if fence != "" {
			continue
		}
		kept = append(kept, line)
	}
	return inlineCodePattern.ReplaceAllString(strings.Join(kept, "\n"), " ")
}

// fenceMarker reports the fenced-code marker a line opens with — "```" or "~~~",
// the two markers CommonMark allows — or "" when the line is not a code fence.
func fenceMarker(trimmed string) string {
	switch {
	case strings.HasPrefix(trimmed, "```"):
		return "```"
	case strings.HasPrefix(trimmed, "~~~"):
		return "~~~"
	default:
		return ""
	}
}

// normalizeForNgram lowercases the text and collapses any whitespace run into a
// single space, dropping a leading space so the first bigram is meaningful.
func normalizeForNgram(text string) []rune {
	lowered := strings.ToLower(text)
	out := make([]rune, 0, len(lowered))
	prevSpace := true
	for _, r := range lowered {
		if unicode.IsSpace(r) {
			if prevSpace {
				continue
			}
			out = append(out, ' ')
			prevSpace = true
			continue
		}
		out = append(out, r)
		prevSpace = false
	}
	return out
}

// aggregateCharNgrams selects the author's most frequent character bigrams and
// records their per-document mean and population standard deviation. A document
// missing a bigram contributes zero for it, so the spread reflects how
// consistently the author uses each bigram.
func aggregateCharNgrams(dist *Distribution, perDoc []Metrics, n float64) {
	totals := make(map[string]float64)
	for _, m := range perDoc {
		for ngram, freq := range m.CharNgrams {
			totals[ngram] += freq
		}
	}
	if len(totals) == 0 {
		dist.Mean.CharNgrams = map[string]float64{}
		dist.StdDev.CharNgrams = map[string]float64{}
		return
	}

	vocabulary := topNgrams(totals, charNgramVocabularySize)

	meanVec := make(map[string]float64, len(vocabulary))
	for _, m := range perDoc {
		for _, ngram := range vocabulary {
			meanVec[ngram] += m.CharNgrams[ngram]
		}
	}
	for ngram := range meanVec {
		meanVec[ngram] /= n
	}

	stdVec := make(map[string]float64, len(vocabulary))
	for _, m := range perDoc {
		for _, ngram := range vocabulary {
			stdVec[ngram] += square(m.CharNgrams[ngram] - meanVec[ngram])
		}
	}
	for ngram := range stdVec {
		stdVec[ngram] = math.Sqrt(stdVec[ngram] / n)
	}

	dist.Mean.CharNgrams = meanVec
	dist.StdDev.CharNgrams = stdVec
}

// topNgrams returns the keys with the highest accumulated frequency, breaking
// ties on the key itself so the selection is deterministic.
func topNgrams(totals map[string]float64, limit int) []string {
	keys := make([]string, 0, len(totals))
	for ngram := range totals {
		keys = append(keys, ngram)
	}
	sort.Slice(keys, func(i int, j int) bool {
		if totals[keys[i]] != totals[keys[j]] {
			return totals[keys[i]] > totals[keys[j]]
		}
		return keys[i] < keys[j]
	})
	if len(keys) > limit {
		keys = keys[:limit]
	}
	return keys
}
