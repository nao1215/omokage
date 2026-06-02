package term

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/nao1215/omokage/internal/feature"
)

// ExtractCorpus reads each file, strips code, and builds the profile-local term
// preferences for the corpus. It mirrors feature.ExtractCorpus so train can hand
// it the same file list. It uses no LLM, network, or external dictionary.
func ExtractCorpus(paths []string) (Profile, error) {
	docs := make([]string, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return Profile{}, fmt.Errorf("read %s: %w", path, err)
		}
		docs = append(docs, string(data))
	}
	return ExtractDocuments(docs), nil
}

// ExtractDocuments builds term preferences from in-memory documents. It is the
// deterministic core shared by ExtractCorpus and the tests:
//
//  1. count surface occurrences per document (code stripped first),
//  2. collect corpus-declared alias links between normalized_keys,
//  3. union normalized_keys into same-concept groups,
//  4. assign each group a deterministic group_key and a preferred surface.
func ExtractDocuments(docs []string) Profile {
	totalCount := make(map[string]int)       // surface -> total occurrences
	docsBySurface := make(map[string]intSet) // surface -> set of doc indices
	keyOfSurface := make(map[string]string)  // surface -> normalized_key
	uf := newUnionFind()
	var links []aliasLink

	for docIndex, doc := range docs {
		prose := feature.StripNonProse(doc)
		for surface, n := range surfaceCounts(prose) {
			totalCount[surface] += n
			if docsBySurface[surface] == nil {
				docsBySurface[surface] = intSet{}
			}
			docsBySurface[surface][docIndex] = struct{}{}
			key := normalizeKey(surface)
			keyOfSurface[surface] = key
			uf.add(key)
		}
		links = append(links, detectAliasLinks(prose)...)
	}

	mergeAliasLinks(uf, links)

	groupKeys := assignGroupKeys(uf)

	// Gather variants per group.
	type groupAcc struct {
		variants []Variant
		docs     intSet
	}
	accs := make(map[string]*groupAcc)
	for surface, key := range keyOfSurface {
		groupKey := groupKeys[uf.find(key)]
		acc := accs[groupKey]
		if acc == nil {
			acc = &groupAcc{docs: intSet{}}
			accs[groupKey] = acc
		}
		acc.variants = append(acc.variants, Variant{
			Surface:       surface,
			NormalizedKey: key,
			GroupKey:      groupKey,
			Count:         totalCount[surface],
			DocCount:      len(docsBySurface[surface]),
		})
		for docIndex := range docsBySurface[surface] {
			acc.docs[docIndex] = struct{}{}
		}
	}

	groups := make([]Group, 0, len(accs))
	for groupKey, acc := range accs {
		sortVariants(acc.variants)
		total := 0
		for _, v := range acc.variants {
			total += v.Count
		}
		groups = append(groups, Group{
			GroupKey:         groupKey,
			PreferredSurface: preferredSurface(acc.variants),
			TotalCount:       total,
			DocCount:         len(acc.docs),
			Variants:         acc.variants,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].GroupKey < groups[j].GroupKey })
	return Profile{Groups: groups}
}

// mergeAliasLinks unions the normalized_keys joined by corpus-declared bridges,
// but only when both sides are unambiguous. An acronym that glosses several
// different Japanese phrases (e.g. "データベース（DB）" and "ダッシュボード（DB）"
// both reducing to "db"), or a phrase glossed by several different acronyms, is
// ambiguous: merging through it would chain unrelated terms into one group, so
// those links are dropped and the terms stay separate. This keeps the bridge
// conservative — a shared short acronym never silently fuses distinct concepts.
func mergeAliasLinks(uf *unionFind, links []aliasLink) {
	phrasesByAcronym := make(map[string]map[string]struct{})
	acronymsByPhrase := make(map[string]map[string]struct{})
	for _, link := range links {
		addToSet(phrasesByAcronym, link.acronym, link.phrase)
		addToSet(acronymsByPhrase, link.phrase, link.acronym)
	}
	for _, link := range links {
		uf.add(link.acronym)
		uf.add(link.phrase)
		if len(phrasesByAcronym[link.acronym]) == 1 && len(acronymsByPhrase[link.phrase]) == 1 {
			uf.union(link.acronym, link.phrase)
		}
	}
}

func addToSet(m map[string]map[string]struct{}, key, value string) {
	if m[key] == nil {
		m[key] = make(map[string]struct{})
	}
	m[key][value] = struct{}{}
}

// assignGroupKeys gives every union-find component a deterministic group_key of
// the form "term:<representative>", where the representative is the
// lexicographically smallest normalized_key in the component. Choosing the
// smallest key makes the id stable regardless of the order surfaces were seen.
func assignGroupKeys(uf *unionFind) map[string]string {
	rep := make(map[string]string) // root -> smallest key in component
	for key := range uf.parent {
		root := uf.find(key)
		if cur, ok := rep[root]; !ok || key < cur {
			rep[root] = key
		}
	}
	out := make(map[string]string, len(rep))
	for root, smallest := range rep {
		out[root] = "term:" + smallest
	}
	return out
}

// preferredSurface picks a group's canonical spelling. The order is fixed and
// must stay stable, because it is part of the documented contract and is pinned
// by tests:
//
//  1. highest DocCount (used in the most documents),
//  2. then highest Count (used most often),
//  3. then the lexicographically smallest Surface, as a deterministic tie-break.
func preferredSurface(variants []Variant) string {
	best := variants[0]
	for _, v := range variants[1:] {
		if betterPreferred(v, best) {
			best = v
		}
	}
	return best.Surface
}

func betterPreferred(candidate, best Variant) bool {
	switch {
	case candidate.DocCount != best.DocCount:
		return candidate.DocCount > best.DocCount
	case candidate.Count != best.Count:
		return candidate.Count > best.Count
	default:
		return candidate.Surface < best.Surface
	}
}

// sortVariants orders a group's variants deterministically for storage and JSON
// output: by descending count, then ascending surface.
func sortVariants(variants []Variant) {
	sort.Slice(variants, func(i, j int) bool {
		if variants[i].Count != variants[j].Count {
			return variants[i].Count > variants[j].Count
		}
		return variants[i].Surface < variants[j].Surface
	})
}

type intSet map[int]struct{}
