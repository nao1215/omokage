package feature

import (
	"math"
	"strings"
	"sync"
	"unicode"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Part-of-speech tags used by the Japanese feature functions.
const (
	posVerb      = "動詞"
	posAdjective = "形容詞"
	posAux       = "助動詞"
	posParticle  = "助詞"
	posSymbol    = "記号"
	posConj      = "接続詞"
)

// jpToken is one morpheme from the Japanese analyzer. Surface is the text as it
// appears, Lemma the dictionary (base) form, and POS the part-of-speech path with
// the dictionary's "*" placeholders dropped.
type jpToken struct {
	Surface string
	Lemma   string
	POS     []string
}

// kagome's IPA tokenizer loads an embedded dictionary, so it is built once and
// reused; Tokenize is safe for concurrent use.
var (
	jpTokOnce sync.Once
	jpTok     *tokenizer.Tokenizer
	jpTokErr  error
)

func japaneseTokenizer() (*tokenizer.Tokenizer, error) {
	jpTokOnce.Do(func() {
		jpTok, jpTokErr = tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	})
	return jpTok, jpTokErr
}

// tokenizeJapanese splits prose into morphemes. On any analyzer error it returns
// nil, so callers fall back to their language-neutral path rather than failing.
func tokenizeJapanese(prose string) []jpToken {
	t, err := japaneseTokenizer()
	if err != nil {
		return nil
	}
	raw := t.Tokenize(prose)
	out := make([]jpToken, 0, len(raw))
	for _, tk := range raw {
		lemma := tk.Surface
		if base, ok := tk.BaseForm(); ok && base != "" && base != "*" {
			lemma = base
		}
		out = append(out, jpToken{Surface: tk.Surface, Lemma: lemma, POS: cleanPOS(tk.POS())})
	}
	return out
}

// cleanPOS drops the dictionary's "*" placeholder tags.
func cleanPOS(pos []string) []string {
	out := make([]string, 0, len(pos))
	for _, p := range pos {
		if p == "*" || p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

// isJapanese reports whether prose is Japanese enough to apply the morphological
// features. It compares Japanese script characters (kanji/kana) against Latin
// letters; a document that is mostly English stays on the language-neutral path.
func isJapanese(prose string) bool {
	var jp, latin int
	for _, r := range prose {
		switch {
		case unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana):
			jp++
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
			latin++
		}
	}
	if jp+latin == 0 {
		return false
	}
	// A modest threshold: any substantially Japanese document qualifies, while a
	// stray Japanese word in English prose does not flip it.
	return jp*100 >= (jp+latin)*30
}

// posNgramVocabularySize caps how many of an author's most frequent POS n-grams
// are kept in the profile, bounding the profile size while keeping the stable
// core of the author's syntactic habits.
const posNgramVocabularySize = 200

// posNgrams returns the relative frequency of part-of-speech bigrams and trigrams
// over the content morphemes (P4). It is empty for a non-Japanese document (nil
// tokens), so the feature simply does not contribute there. Each order is
// normalized by its own count so the two live comparably in one map, mirroring
// charBigrams.
func posNgrams(tokens []jpToken) map[string]float64 {
	seq := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if isContentMorpheme(t) {
			seq = append(seq, pos0(t))
		}
	}
	freq := make(map[string]float64, len(seq)*2)
	addPOSNgrams(freq, seq, 2)
	addPOSNgrams(freq, seq, 3)
	return freq
}

func addPOSNgrams(freq map[string]float64, seq []string, order int) {
	if len(seq) < order {
		return
	}
	counts := make(map[string]int, len(seq))
	total := 0
	for i := 0; i+order <= len(seq); i++ {
		counts[strings.Join(seq[i:i+order], "|")]++
		total++
	}
	for ngram, count := range counts {
		freq[ngram] = float64(count) / float64(total)
	}
}

// aggregatePOSNgrams selects an author's most frequent POS n-grams and records
// their per-document mean and population standard deviation, mirroring
// aggregateCharNgrams. A document missing an n-gram contributes zero for it.
func aggregatePOSNgrams(dist *Distribution, perDoc []Metrics, n float64) {
	totals := make(map[string]float64)
	for _, m := range perDoc {
		for ngram, freq := range m.POSNgrams {
			totals[ngram] += freq
		}
	}
	if len(totals) == 0 {
		dist.Mean.POSNgrams = map[string]float64{}
		dist.StdDev.POSNgrams = map[string]float64{}
		return
	}

	vocabulary := topNgrams(totals, posNgramVocabularySize)

	meanVec := make(map[string]float64, len(vocabulary))
	for _, m := range perDoc {
		for _, ngram := range vocabulary {
			meanVec[ngram] += m.POSNgrams[ngram]
		}
	}
	for ngram := range meanVec {
		meanVec[ngram] /= n
	}

	stdVec := make(map[string]float64, len(vocabulary))
	for _, m := range perDoc {
		for _, ngram := range vocabulary {
			stdVec[ngram] += square(m.POSNgrams[ngram] - meanVec[ngram])
		}
	}
	for ngram := range stdVec {
		stdVec[ngram] = math.Sqrt(stdVec[ngram] / n)
	}

	dist.Mean.POSNgrams = meanVec
	dist.StdDev.POSNgrams = stdVec
}

// typeTokenRatioJP is the lemma-based vocabulary richness (P5): the count of
// distinct content lemmas over the content morpheme count. Lemmatizing first
// (走る/走った/走って → 走る) measures how varied an author's word choice is
// independent of inflection. It is 0 for a non-Japanese document.
func typeTokenRatioJP(tokens []jpToken) float64 {
	uniq := make(map[string]struct{}, len(tokens))
	total := 0
	for _, t := range tokens {
		if !isContentMorpheme(t) {
			continue
		}
		total++
		uniq[t.Lemma] = struct{}{}
	}
	if total == 0 {
		return 0
	}
	return float64(len(uniq)) / float64(total)
}

// pos0 returns the primary part of speech of a token, or "" when absent.
func pos0(t jpToken) string {
	if len(t.POS) > 0 {
		return t.POS[0]
	}
	return ""
}

// isContentMorpheme reports whether a token counts as a real word for token
// statistics: not a symbol and not whitespace.
func isContentMorpheme(t jpToken) bool {
	if pos0(t) == posSymbol {
		return false
	}
	return strings.TrimSpace(t.Surface) != ""
}

// conjunctionStatsJP counts conjunction morphemes (POS 接続詞) over the content
// morpheme count. This is the P1 fix: Japanese has no whitespace, so the
// language-neutral splitWords counted a whole clause as one "token"; a real
// morpheme denominator makes the conjunction frequency meaningful.
func conjunctionStatsJP(tokens []jpToken) (conj int, total int) {
	for _, t := range tokens {
		if !isContentMorpheme(t) {
			continue
		}
		total++
		if pos0(t) == posConj {
			conj++
		}
	}
	return conj, total
}

// terminators marks sentence-final punctuation surfaces.
func isSentenceTerminator(surface string) bool {
	switch surface {
	case "。", "！", "？", "．", "!", "?", ".":
		return true
	default:
		return false
	}
}

// registerStatsJP classifies each sentence's closing predicate as polite (敬体:
// です/ます) or plain (常体: a plain verb/adjective/copula), using morphology
// rather than kana-suffix heuristics. This is the P3 accuracy fix for the
// register feature.
func registerStatsJP(tokens []jpToken) (polite int, plain int) {
	var sentence []jpToken
	flush := func() {
		switch classifySentenceJP(sentence) {
		case registerPolite:
			polite++
		case registerPlain:
			plain++
		}
		sentence = sentence[:0]
	}
	for _, t := range tokens {
		if isSentenceTerminator(t.Surface) {
			flush()
			continue
		}
		sentence = append(sentence, t)
	}
	// A trailing fragment with no terminator is not a completed sentence; drop it
	// (matches the 。！？ denominator ExtractText normalizes by).
	return polite, plain
}

const (
	registerNone = iota
	registerPolite
	registerPlain
)

// politeLemmas are the auxiliary verbs that mark 敬体.
var politeLemmas = map[string]bool{"です": true, "ます": true}

func classifySentenceJP(sentence []jpToken) int {
	// Drop trailing symbols and sentence-final particles (終助詞: ね/よ/か/な…) to
	// expose the predicate that carries the register.
	end := len(sentence)
	for end > 0 {
		t := sentence[end-1]
		if pos0(t) == posSymbol || pos0(t) == "助詞" {
			end--
			continue
		}
		break
	}
	if end == 0 {
		return registerNone
	}
	// Scan the trailing predicate run (動詞/形容詞/助動詞) for a polite auxiliary;
	// ました/ません put です・ます before the final た/ん, so look back a few tokens.
	for i := end - 1; i >= 0 && i >= end-3; i-- {
		t := sentence[i]
		p := pos0(t)
		if p != posVerb && p != posAdjective && p != posAux {
			break
		}
		if p == posAux && politeLemmas[t.Lemma] {
			return registerPolite
		}
	}
	last := sentence[end-1]
	switch pos0(last) {
	case posVerb, posAdjective, posAux:
		return registerPlain
	default:
		// 体言止め (noun stop) or a bare particle: neither register.
		return registerNone
	}
}
