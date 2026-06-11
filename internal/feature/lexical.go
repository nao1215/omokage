package feature

import "strings"

// englishFunctionWords is a closed-class list of high-frequency English words.
// Function-word frequencies are the canonical signal for authorship attribution
// (Burrows's Delta): they are used unconsciously, vary little with topic, and
// differ measurably between authors. They give omokage a language-independent
// fingerprint on top of the Japanese-centric script and register features.
var englishFunctionWords = []string{
	"the", "of", "and", "to", "a", "in", "that", "is", "it", "for", "as", "was",
	"with", "be", "by", "on", "not", "he", "i", "this", "are", "or", "his",
	"from", "at", "which", "but", "have", "an", "they", "you", "had", "we",
	"their", "one", "all", "there", "been", "if", "more", "when", "will", "would",
	"who", "so", "no", "what", "up", "out", "about", "into", "than", "them", "can",
	"only", "other", "some", "could", "these", "then", "do", "any", "my", "now",
	"such", "like", "our", "over", "me", "even", "most", "also", "did", "many",
	"before", "must", "through", "back", "where", "much", "your", "way", "well",
	"down", "should", "each", "just", "those", "how", "too", "very", "make",
	"still", "own", "see", "work", "here", "both", "between", "us", "its", "may",
}

// japaneseFunctionWords lists high-frequency Japanese particles and auxiliaries
// (plus a few common light verbs in dictionary form). When a Japanese analyzer is
// available they are matched as whole morphemes by surface; only the fallback path
// (no analyzer) counts them as substrings. Entries are dictionary forms, so the
// surface-match path counts a conjugating entry (する/ない/なる…) only when it
// appears in that bare form — see lexicalFrequencies for why folding inflections
// onto the lemma was rejected.
var japaneseFunctionWords = []string{
	"の", "は", "を", "に", "が", "と", "で", "も", "へ", "や", "か", "ね", "よ",
	"な", "し", "て", "た", "だ", "です", "ます", "から", "まで", "より", "ので",
	"のに", "けど", "ても", "でも", "という", "について", "における", "として",
	"による", "ように", "こと", "もの", "これ", "それ", "あれ", "この", "その",
	"ある", "いる", "する", "れる", "られる", "せる", "ない", "なる", "また",
	"ため", "など", "だけ", "しか", "ばかり", "ながら", "つつ",
}

// lexicalVocabulary is the ordered union of the English and Japanese function
// words. Every document's lexical vector is keyed by exactly this set so that
// aggregation and scoring can iterate a stable vocabulary.
var lexicalVocabulary = append(append([]string{}, englishFunctionWords...), japaneseFunctionWords...)

// LexicalVocabulary returns the function-word vocabulary backing the lexical
// fingerprint feature. Scoring iterates this list so every profile and target
// is compared over the same keys.
func LexicalVocabulary() []string {
	return lexicalVocabulary
}

// japaneseFunctionWordSet is the membership set of japaneseFunctionWords, for the
// morpheme-based counting path.
var japaneseFunctionWordSet = func() map[string]bool {
	s := make(map[string]bool, len(japaneseFunctionWords))
	for _, fw := range japaneseFunctionWords {
		s[fw] = true
	}
	return s
}()

// lexicalFrequencies builds the per-word relative-frequency vector for a
// document. English words are normalized by the word-token count. For Japanese,
// jpTokens carries the morphemes (P2): each function word is matched as a whole
// morpheme by surface and normalized by the script-character count, which removes
// the substring path's double counting (で in です) and false hits (は inside a
// content word). A nil jpTokens (non-Japanese or analyzer unavailable) falls back
// to the substring count so the feature is never empty.
func lexicalFrequencies(prose string, jpTokens []jpToken) map[string]float64 {
	freq := make(map[string]float64, len(lexicalVocabulary))

	words := splitWords(strings.ToLower(prose))
	wordCount := len(words)
	if wordCount > 0 {
		counts := make(map[string]int, len(words))
		for _, word := range words {
			counts[word]++
		}
		for _, fw := range englishFunctionWords {
			freq[fw] = float64(counts[fw]) / float64(wordCount)
		}
	} else {
		for _, fw := range englishFunctionWords {
			freq[fw] = 0
		}
	}

	// The Japanese denominator stays the script-character count (as before), so the
	// feature keeps its original scale and small-corpus stability; only the
	// numerator changes — exact whole-morpheme matches instead of substring counts,
	// which removes the double counting (で inside です) and false hits (は inside a
	// content word) that the substring path produced.
	kanji, hiragana, katakana := scriptCounts(prose)
	denominator := max(kanji+hiragana+katakana, 1)

	if len(jpTokens) > 0 {
		// Count whole-morpheme surfaces. Particles (は/が/を…) and the polite
		// auxiliaries (です/ます) appear as standalone token surfaces, so matching the
		// surface counts them exactly — without the substring path's double counting
		// (で inside です) or false hits (は inside a content word). The vocabulary's
		// dictionary-form verb/auxiliary entries (する/ない/なる…) therefore count only
		// their bare-form occurrences. Folding inflections onto the lemma was tried
		// and measured to LOWER author-attribution accuracy on the validation corpus
		// (89% → 87%): the base form of a common conjugating word is so frequent that
		// its per-document spread swamps the more discriminating particle counts. So
		// surface matching is the deliberate, validated choice, not an oversight.
		counts := make(map[string]int, len(japaneseFunctionWords))
		for _, t := range jpTokens {
			if isContentMorpheme(t) && japaneseFunctionWordSet[t.Surface] {
				counts[t.Surface]++
			}
		}
		for _, fw := range japaneseFunctionWords {
			freq[fw] = float64(counts[fw]) / float64(denominator)
		}
		return freq
	}

	for _, fw := range japaneseFunctionWords {
		freq[fw] = float64(strings.Count(prose, fw)) / float64(denominator)
	}
	return freq
}
