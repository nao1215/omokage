package feature

import "strings"

// englishFunctionWords is a closed-class list of high-frequency English words.
// Function-word frequencies are the canonical signal for authorship attribution
// (Burrows's Delta): they are used unconsciously, vary little with topic, and
// differ measurably between authors. They give dyer a language-independent
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

// japaneseFunctionWords lists high-frequency Japanese particles and auxiliaries.
// Japanese has no word delimiters, so these are counted as substrings; the same
// counting is applied to every document, so the per-author z-score stays
// meaningful even though overlapping forms (で / です / でも) are double counted.
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

// lexicalFrequencies builds the per-word relative-frequency vector for a
// document. English words are normalized by the word-token count; Japanese
// particles are normalized by the Japanese script-character count. Mixing
// denominators is harmless because each word is later standardized against the
// author's own spread for that word independently.
func lexicalFrequencies(prose string) map[string]float64 {
	freq := make(map[string]float64, len(lexicalVocabulary))

	kanji, hiragana, katakana := scriptCounts(prose)
	scriptTotal := kanji + hiragana + katakana

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

	denominator := max(scriptTotal, 1)
	for _, fw := range japaneseFunctionWords {
		freq[fw] = float64(strings.Count(prose, fw)) / float64(denominator)
	}

	return freq
}
