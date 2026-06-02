package term

// unionFind is a tiny disjoint-set over normalized_keys, used to merge keys that
// a corpus-declared alias bridge linked. It is deterministic: find always
// resolves to a stable root for a given set of unions, and group_key assignment
// then picks the smallest key in each component rather than the root itself, so
// the externally visible id never depends on union order.
type unionFind struct {
	parent map[string]string
}

func newUnionFind() *unionFind {
	return &unionFind{parent: make(map[string]string)}
}

func (u *unionFind) add(key string) {
	if _, ok := u.parent[key]; !ok {
		u.parent[key] = key
	}
}

func (u *unionFind) find(key string) string {
	root := key
	for u.parent[root] != root {
		root = u.parent[root]
	}
	// Path compression keeps repeated lookups cheap without affecting the result.
	for u.parent[key] != root {
		u.parent[key], key = root, u.parent[key]
	}
	return root
}

// union merges the two keys' components. The smaller key (lexicographically) is
// made the root so the structure is deterministic regardless of call order.
func (u *unionFind) union(a, b string) {
	ra, rb := u.find(a), u.find(b)
	if ra == rb {
		return
	}
	if ra < rb {
		u.parent[rb] = ra
	} else {
		u.parent[ra] = rb
	}
}
