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

// stripCode removes fenced code blocks and inline code spans from the text so
// that source code does not contribute to the authorship features. Code shares
// vocabulary and character sequences across authors and would otherwise drown
// out the natural-language signal in technical writing. Both CommonMark fence
// markers are recognized — backtick (``` … ```) and tilde (~~~ … ~~~) — and a
// block is closed only by its own marker, so a tilde line inside a backtick
// block (or vice versa) is treated as content, not a boundary.
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
