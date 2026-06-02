package term

// minASCIILen is the shortest ASCII candidate kept. Two characters keeps the
// common acronyms (DB, CI) while dropping single letters like "a" or "I".
const minASCIILen = 2

// minJapaneseLen is the shortest Japanese candidate kept, in runes. Two runes
// keeps two-kanji compounds (文体) and short katakana words while dropping bare
// single characters.
const minJapaneseLen = 2

// asciiStopwords are common English function words that are never useful as term
// preferences. They are matched against the normalized_key (lower-cased), so the
// scanner never emits them as candidates and they never form a group.
var asciiStopwords = map[string]struct{}{
	"a": {}, "an": {}, "the": {}, "and": {}, "or": {}, "but": {}, "of": {},
	"to": {}, "in": {}, "on": {}, "at": {}, "by": {}, "for": {}, "with": {},
	"as": {}, "is": {}, "are": {}, "was": {}, "were": {}, "be": {}, "been": {},
	"it": {}, "its": {}, "this": {}, "that": {}, "these": {}, "those": {},
	"from": {}, "into": {}, "than": {}, "then": {}, "so": {}, "if": {}, "we": {},
	"you": {}, "they": {}, "he": {}, "she": {}, "i": {}, "my": {}, "our": {},
	"your": {}, "their": {}, "do": {}, "does": {}, "did": {}, "not": {},
	"no": {}, "can": {}, "will": {}, "would": {}, "should": {}, "could": {},
}

// scanCandidates walks prose once and returns the surface of every term
// candidate, in order, with duplicates preserved so the caller can count
// occurrences. Code must already be stripped by the caller.
//
// Two kinds of run are emitted:
//
//   - ASCII-like runs: ASCII or full-width ASCII letters/digits, length ≥
//     minASCIILen, containing at least one letter, and not a stopword. This
//     captures DB, API, HTTP, GoReleaser, database, priority.
//   - Japanese runs: kanji/katakana (with ー), length ≥ minJapaneseLen. This
//     captures データベース, 優先度, 文体, 東京タワー.
//
// A run ends at any rune of the other class or at any other character (spaces,
// punctuation, hiragana), so "データベースDB" yields the two candidates
// "データベース" and "DB".
func scanCandidates(prose string) []string {
	runes := []rune(prose)
	out := make([]string, 0, len(runes)/4)
	for i := 0; i < len(runes); {
		switch {
		case isASCIILike(runes[i]):
			j := i
			for j < len(runes) && isASCIILike(runes[j]) {
				j++
			}
			surface := string(runes[i:j])
			if keepASCII(surface) {
				out = append(out, surface)
			}
			i = j
		case isJapaneseTermRune(runes[i]):
			j := i
			for j < len(runes) && isJapaneseTermRune(runes[j]) {
				j++
			}
			if j-i >= minJapaneseLen {
				out = append(out, string(runes[i:j]))
			}
			i = j
		default:
			i++
		}
	}
	return out
}

// keepASCII reports whether an ASCII-like run is worth keeping as a candidate:
// long enough, holding at least one letter (so "2026" is dropped), and not a
// stopword.
func keepASCII(surface string) bool {
	key := normalizeKey(surface)
	if len([]rune(key)) < minASCIILen {
		return false
	}
	hasLetter := false
	for _, r := range key {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return false
	}
	_, stop := asciiStopwords[key]
	return !stop
}

// surfaceCounts tallies surface occurrences within a single document.
func surfaceCounts(prose string) map[string]int {
	counts := make(map[string]int)
	for _, surface := range scanCandidates(prose) {
		counts[surface]++
	}
	return counts
}
