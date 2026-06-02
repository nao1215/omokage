package storage

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"

	"github.com/nao1215/omokage/internal/term"
)

// SaveTerms persists a profile's term preferences into the same per-author
// database used by SaveProfile. The data is profile-local: this file holds one
// author, so the term tables describe only that author. A re-train replaces the
// previous term data wholesale, inside a transaction, so a reader never sees a
// half-written mix of old and new groups.
func SaveTerms(path string, profile term.Profile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	ctx := context.Background()
	db, err := openDB(ctx, path)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Rollback after a successful Commit is a no-op; the error only matters on the
	// failure paths, where Commit was never reached.
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM term_variant`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM term_group`); err != nil {
		return err
	}

	for _, group := range profile.Groups {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO term_group (group_key, preferred_surface, total_count, doc_count) VALUES (?, ?, ?, ?)`,
			group.GroupKey, group.PreferredSurface, group.TotalCount, group.DocCount,
		); err != nil {
			return err
		}
		for _, variant := range group.Variants {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO term_variant (group_key, normalized_key, surface, count, doc_count) VALUES (?, ?, ?, ?, ?)`,
				variant.GroupKey, variant.NormalizedKey, variant.Surface, variant.Count, variant.DocCount,
			); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// LoadTerms reads a profile's term preferences. A profile trained before term
// support (or one with no extracted terms) simply yields an empty profile, never
// an error, so show/check degrade gracefully. The returned groups and their
// variants are ordered deterministically (group_key ascending; variants by
// descending count then surface).
func LoadTerms(path string) (term.Profile, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return term.Profile{}, nil
		}
		return term.Profile{}, err
	}
	ctx := context.Background()
	db, err := openDB(ctx, path)
	if err != nil {
		return term.Profile{}, err
	}
	defer db.Close()

	variantsByGroup, err := loadVariants(ctx, db)
	if err != nil {
		return term.Profile{}, err
	}

	rows, err := db.QueryContext(ctx,
		`SELECT group_key, preferred_surface, total_count, doc_count FROM term_group ORDER BY group_key`)
	if err != nil {
		return term.Profile{}, err
	}
	defer rows.Close()

	var groups []term.Group
	for rows.Next() {
		var g term.Group
		if err := rows.Scan(&g.GroupKey, &g.PreferredSurface, &g.TotalCount, &g.DocCount); err != nil {
			return term.Profile{}, err
		}
		g.Variants = variantsByGroup[g.GroupKey]
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return term.Profile{}, err
	}
	return term.Profile{Groups: groups}, nil
}

func loadVariants(ctx context.Context, db *sql.DB) (map[string][]term.Variant, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT group_key, normalized_key, surface, count, doc_count FROM term_variant ORDER BY group_key, count DESC, surface`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string][]term.Variant)
	for rows.Next() {
		var v term.Variant
		if err := rows.Scan(&v.GroupKey, &v.NormalizedKey, &v.Surface, &v.Count, &v.DocCount); err != nil {
			return nil, err
		}
		out[v.GroupKey] = append(out[v.GroupKey], v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, variants := range out {
		sort.SliceStable(variants, func(i, j int) bool {
			if variants[i].Count != variants[j].Count {
				return variants[i].Count > variants[j].Count
			}
			return variants[i].Surface < variants[j].Surface
		})
	}
	return out, nil
}

// openDB opens the profile database and ensures the current schema (including the
// term tables) and migrations are applied, matching SaveProfile/LoadProfile.
func openDB(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
