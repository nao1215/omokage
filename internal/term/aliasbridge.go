package term

import (
	"regexp"
	"strings"
)

// aliasLink is a corpus-declared bridge between a spelled-out Japanese phrase and
// its acronym, e.g. phrase="データベース", acronym="db" from "データベース（DB）".
// The two sides are kept distinct (not just an unordered pair) so the extractor
// can detect an ambiguous acronym that glosses several different phrases and
// refuse to chain them together.
type aliasLink struct {
	phrase  string
	acronym string
}

// Bridge sizing. Detection is deliberately narrow: one side must be a genuine
// Japanese phrase (kanji/katakana), the other a short UPPERCASE ASCII acronym
// (DB, API, HTTP, OSS). This is the "長い日本語語句 + 英字略語" pattern the spec
// centers on. Two ordinary words, two long phrases, or two lowercase words never
// link, which also stops a shared short word from transitively chaining many
// unrelated terms into one giant group.
const (
	minLongJapaneseLen = 3
	minAcronymLen      = 2
	maxAcronymLen      = 8
)

// parenBridge matches a term immediately followed by a parenthesized gloss, in
// either ASCII or full-width parentheses. The left side is captured greedily as a
// trailing ASCII/Japanese run; the inside is captured raw so the "以下" forms can
// be handled separately. It is only a coarse candidate matcher — bridge()
// applies the real Japanese-phrase ↔ acronym test. Examples it feeds to bridge():
//
//	データベース（DB）   DB（データベース）   データベース（以下 DB）
//
// An English-only pair like "database (DB)" is matched here but rejected by
// bridge(): English phrase ↔ acronym is intentionally out of scope (it would
// chain too eagerly), so only Japanese-phrase ↔ acronym pairs bridge.
var parenBridge = regexp.MustCompile(`([\p{Han}\p{Katakana}ーA-Za-z0-9Ａ-Ｚａ-ｚ０-９]+)\s*[（(]([^）)]*)[）)]`)

// inlineBridge matches the no-parenthesis "X。以下、Y" abbreviation declaration,
// e.g. "データベース。以下、DB" or "データベース。以下DB".
var inlineBridge = regexp.MustCompile(`([\p{Han}\p{Katakana}ーA-Za-z0-9Ａ-Ｚａ-ｚ０-９]{2,})。\s*以下[、,]?\s*([A-Za-zＡ-Ｚａ-ｚ][\p{Han}\p{Katakana}ーA-Za-z0-9Ａ-Ｚａ-ｚ０-９]*)`)

// detectAliasLinks scans prose for conservative, corpus-declared alias bridges
// and returns the normalized_key links they imply. It never infers a bridge from
// a bare "A（B）"; one side must be a long phrase and the other a short ASCII
// abbreviation (the orientation may be either way around), or the text must use
// an explicit "以下 X" declaration.
func detectAliasLinks(prose string) []aliasLink {
	var links []aliasLink

	for _, m := range parenBridge.FindAllStringSubmatch(prose, -1) {
		left := m[1]
		inner := strings.TrimSpace(m[2])
		if rest, ok := stripIfuka(inner); ok {
			// "以下 DB" explicitly declares DB as the abbreviation of left.
			if link, ok := bridge(left, rest); ok {
				links = append(links, link)
			}
			continue
		}
		if link, ok := bridge(left, inner); ok {
			links = append(links, link)
		}
	}

	for _, m := range inlineBridge.FindAllStringSubmatch(prose, -1) {
		if link, ok := bridge(m[1], m[2]); ok {
			links = append(links, link)
		}
	}

	return links
}

// stripIfuka strips a leading "以下" (optionally followed by "、" or ",") from a
// parenthetical gloss, reporting whether it was present. "以下 DB" -> "DB".
func stripIfuka(s string) (string, bool) {
	if !strings.HasPrefix(s, "以下") {
		return s, false
	}
	rest := strings.TrimPrefix(s, "以下")
	rest = strings.TrimLeft(rest, "、, \t")
	return strings.TrimSpace(rest), true
}

// bridge validates that one side is a long Japanese phrase and the other a short
// uppercase ASCII acronym (in either order) and, if so, returns the
// normalized_key link between them. A pair of ordinary words, two phrases, or a
// lowercase word is rejected, so a bare "A（B）" never links and shared short
// words cannot chain unrelated terms together.
func bridge(left, right string) (aliasLink, bool) {
	lk := normalizeKey(left)
	rk := normalizeKey(right)
	if lk == "" || rk == "" || lk == rk {
		return aliasLink{}, false
	}
	if isJapanesePhrase(left) && isAcronym(right) {
		return aliasLink{phrase: lk, acronym: rk}, true
	}
	if isAcronym(left) && isJapanesePhrase(right) {
		return aliasLink{phrase: rk, acronym: lk}, true
	}
	return aliasLink{}, false
}

// isJapanesePhrase reports whether a surface is a long-enough Japanese term: at
// least minLongJapaneseLen kanji/katakana runes and no ASCII letters (so it is a
// spelled-out word, not an acronym).
func isJapanesePhrase(surface string) bool {
	japanese := 0
	for _, r := range surface {
		switch {
		case isJapaneseTermRune(r):
			japanese++
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			return false
		}
	}
	return japanese >= minLongJapaneseLen
}

// isAcronym reports whether a surface looks like an English acronym: only
// UPPERCASE ASCII letters and digits (after folding full-width forms), at least
// one letter, within the length bound. DB, API, HTTP, OSS qualify; "Debian",
// "linux", and "testing" do not, which is what keeps bridges from chaining.
func isAcronym(surface string) bool {
	runes := []rune(foldRune(strings.TrimFunc(surface, isTrimmable)))
	if len(runes) < minAcronymLen || len(runes) > maxAcronymLen {
		return false
	}
	hasLetter := false
	for _, r := range runes {
		switch {
		case r >= 'A' && r <= 'Z':
			hasLetter = true
		case r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return hasLetter
}
