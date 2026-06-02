package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/nao1215/omokage/internal/feature"
	"github.com/nao1215/omokage/internal/profile"
)

const schema = `
CREATE TABLE IF NOT EXISTS profile (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  author TEXT NOT NULL,
  source_dir TEXT NOT NULL,
  trained_at TEXT NOT NULL,
  file_count INTEGER NOT NULL,
  mean_average_sentence_length REAL NOT NULL,
  mean_sentence_length_variance REAL NOT NULL,
  mean_punctuation_frequency REAL NOT NULL,
  mean_newline_frequency REAL NOT NULL,
  mean_bullet_ratio REAL NOT NULL,
  mean_conjunction_frequency REAL NOT NULL,
  mean_kanji_ratio REAL NOT NULL,
  mean_hiragana_ratio REAL NOT NULL,
  mean_katakana_ratio REAL NOT NULL,
  mean_paragraph_length_variance REAL NOT NULL,
  mean_markdown_structure_density REAL NOT NULL,
  mean_polite_ending_ratio REAL NOT NULL,
  mean_plain_ending_ratio REAL NOT NULL,
  std_average_sentence_length REAL NOT NULL,
  std_sentence_length_variance REAL NOT NULL,
  std_punctuation_frequency REAL NOT NULL,
  std_newline_frequency REAL NOT NULL,
  std_bullet_ratio REAL NOT NULL,
  std_conjunction_frequency REAL NOT NULL,
  std_kanji_ratio REAL NOT NULL,
  std_hiragana_ratio REAL NOT NULL,
  std_katakana_ratio REAL NOT NULL,
  std_paragraph_length_variance REAL NOT NULL,
  std_markdown_structure_density REAL NOT NULL,
  std_polite_ending_ratio REAL NOT NULL,
  std_plain_ending_ratio REAL NOT NULL,
  document_count INTEGER NOT NULL,
  sentence_count INTEGER NOT NULL,
  character_count INTEGER NOT NULL,
  mean_lexical_frequencies TEXT NOT NULL DEFAULT '{}',
  std_lexical_frequencies TEXT NOT NULL DEFAULT '{}',
  mean_char_ngrams TEXT NOT NULL DEFAULT '{}',
  std_char_ngrams TEXT NOT NULL DEFAULT '{}',
  sources TEXT NOT NULL DEFAULT '[]'
);

-- Term preferences are profile-local: this database holds exactly one author's
-- profile, so these tables describe only that author's notation choices. No
-- profile_id column is needed. normalized_key and group_key are kept as separate
-- columns so a reader can tell whether two surfaces were merged by plain
-- normalization (same normalized_key) or by a corpus-declared alias bridge (same
-- group_key spanning different normalized_keys). Original training text is never
-- stored — only surfaces and their counts.
CREATE TABLE IF NOT EXISTS term_group (
  group_key TEXT PRIMARY KEY,
  preferred_surface TEXT NOT NULL,
  total_count INTEGER NOT NULL,
  doc_count INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS term_variant (
  group_key TEXT NOT NULL,
  normalized_key TEXT NOT NULL,
  surface TEXT NOT NULL,
  count INTEGER NOT NULL,
  doc_count INTEGER NOT NULL,
  PRIMARY KEY (group_key, surface)
);`

// migrate brings an existing profile database up to the current schema. The
// `sources` column was added after the original schema shipped, so a profile
// trained by an older omokage lacks it. CREATE TABLE IF NOT EXISTS never alters
// an existing table, so we add the column explicitly and ignore the
// "duplicate column name" error a freshly created (already-current) table returns.
func migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE profile ADD COLUMN sources TEXT NOT NULL DEFAULT '[]'`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}
	return nil
}

func SaveProfile(path string, record profile.Record) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	ctx := context.Background()

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err
	}
	if err := migrate(ctx, db); err != nil {
		return err
	}

	const query = `
INSERT INTO profile (
  id, author, source_dir, trained_at, file_count,
  mean_average_sentence_length, mean_sentence_length_variance, mean_punctuation_frequency,
  mean_newline_frequency, mean_bullet_ratio, mean_conjunction_frequency,
  mean_kanji_ratio, mean_hiragana_ratio, mean_katakana_ratio,
  mean_paragraph_length_variance, mean_markdown_structure_density,
  mean_polite_ending_ratio, mean_plain_ending_ratio,
  std_average_sentence_length, std_sentence_length_variance, std_punctuation_frequency,
  std_newline_frequency, std_bullet_ratio, std_conjunction_frequency,
  std_kanji_ratio, std_hiragana_ratio, std_katakana_ratio,
  std_paragraph_length_variance, std_markdown_structure_density,
  std_polite_ending_ratio, std_plain_ending_ratio,
  document_count, sentence_count, character_count,
  mean_lexical_frequencies, std_lexical_frequencies,
  mean_char_ngrams, std_char_ngrams, sources
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
  author = excluded.author,
  source_dir = excluded.source_dir,
  trained_at = excluded.trained_at,
  file_count = excluded.file_count,
  mean_average_sentence_length = excluded.mean_average_sentence_length,
  mean_sentence_length_variance = excluded.mean_sentence_length_variance,
  mean_punctuation_frequency = excluded.mean_punctuation_frequency,
  mean_newline_frequency = excluded.mean_newline_frequency,
  mean_bullet_ratio = excluded.mean_bullet_ratio,
  mean_conjunction_frequency = excluded.mean_conjunction_frequency,
  mean_kanji_ratio = excluded.mean_kanji_ratio,
  mean_hiragana_ratio = excluded.mean_hiragana_ratio,
  mean_katakana_ratio = excluded.mean_katakana_ratio,
  mean_paragraph_length_variance = excluded.mean_paragraph_length_variance,
  mean_markdown_structure_density = excluded.mean_markdown_structure_density,
  mean_polite_ending_ratio = excluded.mean_polite_ending_ratio,
  mean_plain_ending_ratio = excluded.mean_plain_ending_ratio,
  std_average_sentence_length = excluded.std_average_sentence_length,
  std_sentence_length_variance = excluded.std_sentence_length_variance,
  std_punctuation_frequency = excluded.std_punctuation_frequency,
  std_newline_frequency = excluded.std_newline_frequency,
  std_bullet_ratio = excluded.std_bullet_ratio,
  std_conjunction_frequency = excluded.std_conjunction_frequency,
  std_kanji_ratio = excluded.std_kanji_ratio,
  std_hiragana_ratio = excluded.std_hiragana_ratio,
  std_katakana_ratio = excluded.std_katakana_ratio,
  std_paragraph_length_variance = excluded.std_paragraph_length_variance,
  std_markdown_structure_density = excluded.std_markdown_structure_density,
  std_polite_ending_ratio = excluded.std_polite_ending_ratio,
  std_plain_ending_ratio = excluded.std_plain_ending_ratio,
  document_count = excluded.document_count,
  sentence_count = excluded.sentence_count,
  character_count = excluded.character_count,
  mean_lexical_frequencies = excluded.mean_lexical_frequencies,
  std_lexical_frequencies = excluded.std_lexical_frequencies,
  mean_char_ngrams = excluded.mean_char_ngrams,
  std_char_ngrams = excluded.std_char_ngrams,
  sources = excluded.sources;
`

	mean := record.Distribution.Mean
	std := record.Distribution.StdDev
	_, err = db.ExecContext(
		ctx,
		query,
		1,
		record.Author,
		record.SourceDir,
		record.TrainedAt.Format(time.RFC3339),
		record.FileCount,
		mean.AverageSentenceLength,
		mean.SentenceLengthVariance,
		mean.PunctuationFrequency,
		mean.NewlineFrequency,
		mean.BulletRatio,
		mean.ConjunctionFrequency,
		mean.KanjiRatio,
		mean.HiraganaRatio,
		mean.KatakanaRatio,
		mean.ParagraphLengthVariance,
		mean.MarkdownStructureDensity,
		mean.PoliteEndingRatio,
		mean.PlainEndingRatio,
		std.AverageSentenceLength,
		std.SentenceLengthVariance,
		std.PunctuationFrequency,
		std.NewlineFrequency,
		std.BulletRatio,
		std.ConjunctionFrequency,
		std.KanjiRatio,
		std.HiraganaRatio,
		std.KatakanaRatio,
		std.ParagraphLengthVariance,
		std.MarkdownStructureDensity,
		std.PoliteEndingRatio,
		std.PlainEndingRatio,
		record.Distribution.DocumentCount,
		record.Distribution.SentenceCount,
		record.Distribution.CharacterCount,
		marshalLexical(mean.LexicalFrequencies),
		marshalLexical(std.LexicalFrequencies),
		marshalLexical(mean.CharNgrams),
		marshalLexical(std.CharNgrams),
		marshalSources(record.Sources),
	)
	return err
}

// marshalSources serializes the list of learning-source paths to a JSON array
// for storage, defaulting to an empty array so the NOT NULL column always has a
// value.
func marshalSources(sources []string) string {
	if len(sources) == 0 {
		return "[]"
	}
	encoded, err := json.Marshal(sources)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

// marshalLexical serializes a lexical frequency vector to JSON for storage,
// defaulting to an empty object so the NOT NULL column always has a value.
func marshalLexical(vector map[string]float64) string {
	if len(vector) == 0 {
		return "{}"
	}
	encoded, err := json.Marshal(vector)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func LoadProfile(path string) (profile.Record, error) {
	ctx := context.Background()
	// Opening the database would create an empty file for an untrained author,
	// which then shows up in `list`. Reject a missing profile up front instead.
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return profile.Record{}, fmt.Errorf("profile not found: %s", path)
		}
		return profile.Record{}, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return profile.Record{}, err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return profile.Record{}, err
	}
	if err := migrate(ctx, db); err != nil {
		return profile.Record{}, err
	}

	row := db.QueryRowContext(ctx, `
SELECT
  author,
  source_dir,
  trained_at,
  file_count,
  mean_average_sentence_length,
  mean_sentence_length_variance,
  mean_punctuation_frequency,
  mean_newline_frequency,
  mean_bullet_ratio,
  mean_conjunction_frequency,
  mean_kanji_ratio,
  mean_hiragana_ratio,
  mean_katakana_ratio,
  mean_paragraph_length_variance,
  mean_markdown_structure_density,
  mean_polite_ending_ratio,
  mean_plain_ending_ratio,
  std_average_sentence_length,
  std_sentence_length_variance,
  std_punctuation_frequency,
  std_newline_frequency,
  std_bullet_ratio,
  std_conjunction_frequency,
  std_kanji_ratio,
  std_hiragana_ratio,
  std_katakana_ratio,
  std_paragraph_length_variance,
  std_markdown_structure_density,
  std_polite_ending_ratio,
  std_plain_ending_ratio,
  document_count,
  sentence_count,
  character_count,
  mean_lexical_frequencies,
  std_lexical_frequencies,
  mean_char_ngrams,
  std_char_ngrams,
  sources
FROM profile
WHERE id = 1
`)

	var trainedAt string
	var meanLexicalJSON string
	var stdLexicalJSON string
	var meanNgramJSON string
	var stdNgramJSON string
	var sourcesJSON string
	var record profile.Record
	var mean feature.Metrics
	var std feature.Metrics
	var dist feature.Distribution
	if err := row.Scan(
		&record.Author,
		&record.SourceDir,
		&trainedAt,
		&record.FileCount,
		&mean.AverageSentenceLength,
		&mean.SentenceLengthVariance,
		&mean.PunctuationFrequency,
		&mean.NewlineFrequency,
		&mean.BulletRatio,
		&mean.ConjunctionFrequency,
		&mean.KanjiRatio,
		&mean.HiraganaRatio,
		&mean.KatakanaRatio,
		&mean.ParagraphLengthVariance,
		&mean.MarkdownStructureDensity,
		&mean.PoliteEndingRatio,
		&mean.PlainEndingRatio,
		&std.AverageSentenceLength,
		&std.SentenceLengthVariance,
		&std.PunctuationFrequency,
		&std.NewlineFrequency,
		&std.BulletRatio,
		&std.ConjunctionFrequency,
		&std.KanjiRatio,
		&std.HiraganaRatio,
		&std.KatakanaRatio,
		&std.ParagraphLengthVariance,
		&std.MarkdownStructureDensity,
		&std.PoliteEndingRatio,
		&std.PlainEndingRatio,
		&dist.DocumentCount,
		&dist.SentenceCount,
		&dist.CharacterCount,
		&meanLexicalJSON,
		&stdLexicalJSON,
		&meanNgramJSON,
		&stdNgramJSON,
		&sourcesJSON,
	); err != nil {
		if err == sql.ErrNoRows {
			return profile.Record{}, fmt.Errorf("profile not found: %s", path)
		}
		return profile.Record{}, err
	}

	record.TrainedAt, err = time.Parse(time.RFC3339, trainedAt)
	if err != nil {
		return profile.Record{}, err
	}
	mean.LexicalFrequencies = unmarshalLexical(meanLexicalJSON)
	std.LexicalFrequencies = unmarshalLexical(stdLexicalJSON)
	mean.CharNgrams = unmarshalLexical(meanNgramJSON)
	std.CharNgrams = unmarshalLexical(stdNgramJSON)
	dist.Mean = mean
	dist.StdDev = std
	record.Distribution = dist
	record.Sources = unmarshalSources(sourcesJSON)
	// A profile trained before multi-input support has no sources list. Populate it
	// from the single SourceDir so every consumer can rely on Sources being set.
	if len(record.Sources) == 0 && record.SourceDir != "" {
		record.Sources = []string{record.SourceDir}
	}
	return record, nil
}

// unmarshalSources deserializes a stored JSON array of learning-source paths,
// returning an empty slice for missing or malformed data.
func unmarshalSources(encoded string) []string {
	if encoded == "" || encoded == "[]" {
		return nil
	}
	var sources []string
	if err := json.Unmarshal([]byte(encoded), &sources); err != nil {
		return nil
	}
	return sources
}

// unmarshalLexical deserializes a stored lexical frequency vector, returning an
// empty map for missing or malformed data so scoring can index it safely.
func unmarshalLexical(encoded string) map[string]float64 {
	vector := make(map[string]float64)
	if encoded == "" || encoded == "{}" {
		return vector
	}
	if err := json.Unmarshal([]byte(encoded), &vector); err != nil {
		return make(map[string]float64)
	}
	return vector
}
