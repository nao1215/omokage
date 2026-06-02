// Package term extracts a learning corpus's notation preferences — which
// surface form an author actually uses for a recurring term (DB vs データベース,
// HTTP vs http) — without any LLM, network access, or external dictionary.
//
// The extraction is deterministic and reproducible: the same corpus always
// yields the same groups, preferred surfaces, and counts. It deliberately does
// not attempt open-ended synonym discovery. Only two things merge surfaces:
//
//   - normalization (normalized_key): case, full-width/half-width ASCII, and
//     surrounding punctuation differences are folded away deterministically, so
//     DB / db / ＤＢ share one normalized_key.
//   - alias bridges (group_key): two different normalized_keys are treated as the
//     same concept ONLY when the corpus itself spells out the bridge, e.g.
//     "データベース（DB）" or "データベース。以下、DB". A bare "A（B）" is never
//     enough; see aliasbridge.go for the conservative rules.
//
// normalized_key and group_key are kept as separate responsibilities so a
// consumer can always tell whether two surfaces were merged by plain
// normalization or by a corpus-declared alias bridge.
package term

// Variant is one distinct surface form of a term, scoped to a single profile
// (one profile == one SQLite database). Two surfaces that fold to the same
// normalized_key (DB and ＤＢ) are still two variants; they share NormalizedKey
// and GroupKey but keep their own Surface and counts.
type Variant struct {
	// Surface is the exact text as it appeared in the corpus.
	Surface string
	// NormalizedKey is the deterministic normalization of Surface (case,
	// full-width ASCII, and surrounding punctuation folded). It is the unit of
	// "these are the same spelling".
	NormalizedKey string
	// GroupKey is the identifier of the same-concept group this variant belongs
	// to. It equals "term:"+NormalizedKey for a surface that was never alias
	// bridged; it differs only when a corpus-declared bridge merged it with
	// another normalized_key.
	GroupKey string
	// Count is how many times Surface occurred across the corpus.
	Count int
	// DocCount is how many documents contained Surface at least once.
	DocCount int
}

// Group is a same-concept cluster of variants within one profile. Its members
// share a GroupKey. A group whose variants all share a single normalized_key was
// formed by normalization alone; a group spanning several normalized_keys was
// formed by at least one alias bridge.
type Group struct {
	GroupKey         string
	PreferredSurface string
	TotalCount       int
	DocCount         int
	Variants         []Variant
}

// Profile is a profile-local set of term preferences. It is stored alongside the
// style distribution in the same per-author SQLite database, never shared
// between authors, so profile A can prefer "DB" while profile B prefers
// "データベース".
type Profile struct {
	Groups []Group
}

// Occurrence is a position of a non-preferred surface inside a draft. The check
// path does not populate it yet; the field exists so a future version can add
// line/column information to a warning without changing the JSON shape.
type Occurrence struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Warning reports that a draft used a surface other than the group's preferred
// one. Occurrences is optional and omitted while only counts are tracked.
type Warning struct {
	GroupKey         string       `json:"group_key"`
	PreferredSurface string       `json:"preferred_surface"`
	UsedSurface      string       `json:"used_surface"`
	Count            int          `json:"count"`
	Occurrences      []Occurrence `json:"occurrences,omitempty"`
}
