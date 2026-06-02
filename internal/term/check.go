package term

import (
	"sort"

	"github.com/nao1215/omokage/internal/feature"
)

// CheckText reports where a draft uses a surface other than the preferred one for
// a term the profile knows. It is a separate layer from the style similarity
// score: it never changes the score, it only surfaces notation deviations.
//
// Matching is by normalized_key: a draft surface is mapped to its group through
// the variants' normalized_keys, so a draft that types "db" or the bridged
// "データベース" is flagged against a profile that prefers "DB". A surface that
// equals the group's preferred surface is not flagged. Occurrences is left empty
// for now; only counts are reported.
func (p Profile) CheckText(text string) []Warning {
	if len(p.Groups) == 0 {
		return nil
	}
	groupByKey := make(map[string]Group, len(p.Groups))
	groupKeyForNorm := make(map[string]string)
	for _, g := range p.Groups {
		groupByKey[g.GroupKey] = g
		for _, v := range g.Variants {
			// Extraction guarantees one normalized_key maps to one group, but a profile
			// loaded from an external database might not. Keep the first mapping (groups
			// arrive ordered by group_key) so a malformed store cannot make warnings
			// non-deterministic.
			if _, seen := groupKeyForNorm[v.NormalizedKey]; !seen {
				groupKeyForNorm[v.NormalizedKey] = g.GroupKey
			}
		}
	}

	type usage struct {
		groupKey string
		surface  string
	}
	counts := make(map[usage]int)
	for _, surface := range scanCandidates(feature.StripNonProse(text)) {
		groupKey, ok := groupKeyForNorm[normalizeKey(surface)]
		if !ok {
			continue
		}
		if surface == groupByKey[groupKey].PreferredSurface {
			continue
		}
		counts[usage{groupKey: groupKey, surface: surface}]++
	}

	warnings := make([]Warning, 0, len(counts))
	for u, count := range counts {
		warnings = append(warnings, Warning{
			GroupKey:         u.groupKey,
			PreferredSurface: groupByKey[u.groupKey].PreferredSurface,
			UsedSurface:      u.surface,
			Count:            count,
		})
	}
	sort.Slice(warnings, func(i, j int) bool {
		if warnings[i].GroupKey != warnings[j].GroupKey {
			return warnings[i].GroupKey < warnings[j].GroupKey
		}
		return warnings[i].UsedSurface < warnings[j].UsedSurface
	})
	return warnings
}
