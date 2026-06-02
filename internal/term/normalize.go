package term

import (
	"strings"
	"unicode"
)

// normalizeKey folds a surface into its deterministic normalized_key. The fold
// is intentionally limited to changes that never alter meaning, so two surfaces
// share a normalized_key only when they are genuinely the same spelling:
//
//   - surrounding punctuation, brackets, and whitespace are trimmed,
//   - full-width ASCII letters/digits (Ａ-Ｚ ａ-ｚ ０-９) collapse to ASCII,
//   - ASCII letters are lower-cased.
//
// So DB, db, and ＤＢ all normalize to "db". Japanese text is left as-is (data
// stays データベース): folding kana would risk merging distinct words, which is
// out of scope for this conservative initial implementation. Half-width katakana
// is not folded for the same reason (see the package limitations in the README).
func normalizeKey(surface string) string {
	folded := make([]rune, 0, len(surface))
	for _, r := range foldRune(surface) {
		folded = append(folded, r)
	}
	trimmed := strings.TrimFunc(string(folded), isTrimmable)
	return strings.ToLower(trimmed)
}

// foldRune returns surface with full-width ASCII folded to ASCII. It is split
// out so candidate scanning can recognize ＤＢ as an ASCII-like run while keeping
// the original surface intact.
func foldRune(surface string) string {
	var b strings.Builder
	b.Grow(len(surface))
	for _, r := range surface {
		b.WriteRune(foldWidth(r))
	}
	return b.String()
}

// foldWidth maps a single full-width ASCII rune to its half-width ASCII form and
// leaves every other rune unchanged.
func foldWidth(r rune) rune {
	if r >= '！' && r <= '～' { // U+FF01..U+FF5E -> U+0021..U+007E
		return r - 0xFEE0
	}
	return r
}

// isTrimmable reports whether a rune should be stripped from the edges of a
// surface when forming its normalized_key: spaces, and leading/trailing
// punctuation or symbols such as brackets and quotation marks.
func isTrimmable(r rune) bool {
	return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// isASCIILike reports whether a rune is an ASCII letter/digit or its full-width
// equivalent, so DB and ＤＢ scan as a single candidate run.
func isASCIILike(r rune) bool {
	f := foldWidth(r)
	return (f >= 'a' && f <= 'z') || (f >= 'A' && f <= 'Z') || (f >= '0' && f <= '9')
}

// isJapaneseTermRune reports whether a rune may be part of a Japanese term
// candidate. Kanji and katakana (with the prolonged-sound mark) form noun and
// proper-noun cores; hiragana is excluded because runs of hiragana are usually
// grammatical (particles, verb/adjective inflection) rather than terms.
func isJapaneseTermRune(r rune) bool {
	if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Katakana, r) {
		return true
	}
	return r == 'ー' // U+30FC prolonged sound mark, common inside katakana words
}
