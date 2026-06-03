package storage

import (
	"context"
	"database/sql"
	"errors"
	"os"
)

// SaveQualityFindings stores the corpus-quality findings computed at train time
// in the profile's database, as the opaque JSON the caller hands over. Keeping it
// as raw bytes lets storage stay agnostic of the quality package (avoiding an
// import cycle, since quality depends on profile) while still letting `show`
// report what `doctor` found at training time, including the per-document
// findings the stored distribution alone cannot reproduce. The findings hold only
// counts, ratings, and file names — never the training text. It updates the
// single profile row, so SaveProfile must have run first.
func SaveQualityFindings(path string, data []byte) error {
	ctx := context.Background()
	db, err := openDB(ctx, path)
	if err != nil {
		return err
	}
	defer db.Close()

	result, err := db.ExecContext(ctx, `UPDATE profile SET quality_findings = ? WHERE id = 1`, string(data))
	if err != nil {
		return err
	}
	// A missing profile row means SaveProfile was not run first; report it rather
	// than silently storing nothing.
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return errors.New("no profile to attach quality findings to; save the profile first")
	}
	return nil
}

// LoadQualityFindings returns the stored quality findings as raw JSON. A profile
// trained before quality findings were stored (or one whose corpus looked clean)
// yields the empty-array default, never an error, so `show` degrades gracefully
// to "no findings". The caller unmarshals the bytes.
func LoadQualityFindings(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return []byte("[]"), nil
		}
		return nil, err
	}
	ctx := context.Background()
	db, err := openDB(ctx, path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var data string
	switch err := db.QueryRowContext(ctx, `SELECT quality_findings FROM profile WHERE id = 1`).Scan(&data); {
	case errors.Is(err, sql.ErrNoRows):
		return []byte("[]"), nil
	case err != nil:
		return nil, err
	}
	if data == "" {
		return []byte("[]"), nil
	}
	return []byte(data), nil
}
