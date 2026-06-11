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
  feature_version INTEGER NOT NULL DEFAULT 1,
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
  mean_type_token_ratio REAL NOT NULL DEFAULT 0,
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
  std_type_token_ratio REAL NOT NULL DEFAULT 0,
  document_count INTEGER NOT NULL,
  sentence_count INTEGER NOT NULL,
  character_count INTEGER NOT NULL,
  mean_lexical_frequencies TEXT NOT NULL DEFAULT '{}',
  std_lexical_frequencies TEXT NOT NULL DEFAULT '{}',
  mean_char_ngrams TEXT NOT NULL DEFAULT '{}',
  std_char_ngrams TEXT NOT NULL DEFAULT '{}',
  mean_pos_ngrams TEXT NOT NULL DEFAULT '{}',
  std_pos_ngrams TEXT NOT NULL DEFAULT '{}',
  sources TEXT NOT NULL DEFAULT '[]',
  quality_findings TEXT NOT NULL DEFAULT '[]'
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

// migrate brings an existing profile database up to the current schema. Columns
// added after the original schema shipped (`sources`, then `quality_findings`)
// are absent from a profile trained by an older omokage. CREATE TABLE IF NOT
// EXISTS never alters an existing table, so each new column is added explicitly,
// ignoring the "duplicate column name" error a freshly created (already-current)
// table returns.
func migrate(ctx context.Context, db *sql.DB) error {
	for _, column := range []string{
		`ALTER TABLE profile ADD COLUMN sources TEXT NOT NULL DEFAULT '[]'`,
		`ALTER TABLE profile ADD COLUMN quality_findings TEXT NOT NULL DEFAULT '[]'`,
		`ALTER TABLE profile ADD COLUMN mean_pos_ngrams TEXT NOT NULL DEFAULT '{}'`,
		`ALTER TABLE profile ADD COLUMN std_pos_ngrams TEXT NOT NULL DEFAULT '{}'`,
		`ALTER TABLE profile ADD COLUMN feature_version INTEGER NOT NULL DEFAULT 1`,
		`ALTER TABLE profile ADD COLUMN mean_type_token_ratio REAL NOT NULL DEFAULT 0`,
		`ALTER TABLE profile ADD COLUMN std_type_token_ratio REAL NOT NULL DEFAULT 0`,
	} {
		if _, err := db.ExecContext(ctx, column); err != nil && !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
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
  id, author, source_dir, trained_at, file_count, feature_version,
  mean_average_sentence_length, mean_sentence_length_variance, mean_punctuation_frequency,
  mean_newline_frequency, mean_bullet_ratio, mean_conjunction_frequency,
  mean_kanji_ratio, mean_hiragana_ratio, mean_katakana_ratio,
  mean_paragraph_length_variance, mean_markdown_structure_density,
  mean_polite_ending_ratio, mean_plain_ending_ratio, mean_type_token_ratio,
  std_average_sentence_length, std_sentence_length_variance, std_punctuation_frequency,
  std_newline_frequency, std_bullet_ratio, std_conjunction_frequency,
  std_kanji_ratio, std_hiragana_ratio, std_katakana_ratio,
  std_paragraph_length_variance, std_markdown_structure_density,
  std_polite_ending_ratio, std_plain_ending_ratio, std_type_token_ratio,
  document_count, sentence_count, character_count,
  mean_lexical_frequencies, std_lexical_frequencies,
  mean_char_ngrams, std_char_ngrams, mean_pos_ngrams, std_pos_ngrams, sources
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
  author = excluded.author,
  source_dir = excluded.source_dir,
  trained_at = excluded.trained_at,
  file_count = excluded.file_count,
  feature_version = excluded.feature_version,
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
  mean_type_token_ratio = excluded.mean_type_token_ratio,
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
  std_type_token_ratio = excluded.std_type_token_ratio,
  document_count = excluded.document_count,
  sentence_count = excluded.sentence_count,
  character_count = excluded.character_count,
  mean_lexical_frequencies = excluded.mean_lexical_frequencies,
  std_lexical_frequencies = excluded.std_lexical_frequencies,
  mean_char_ngrams = excluded.mean_char_ngrams,
  std_char_ngrams = excluded.std_char_ngrams,
  mean_pos_ngrams = excluded.mean_pos_ngrams,
  std_pos_ngrams = excluded.std_pos_ngrams,
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
		featureVersionOrDefault(record.FeatureVersion),
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
		mean.TypeTokenRatio,
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
		std.TypeTokenRatio,
		record.Distribution.DocumentCount,
		record.Distribution.SentenceCount,
		record.Distribution.CharacterCount,
		marshalLexical(mean.LexicalFrequencies),
		marshalLexical(std.LexicalFrequencies),
		marshalLexical(mean.CharNgrams),
		marshalLexical(std.CharNgrams),
		marshalLexical(mean.POSNgrams),
		marshalLexical(std.POSNgrams),
		marshalSources(record.Sources),
	)
	return err
}

// featureVersionOrDefault maps an unset feature version (0, e.g. a record built
// in a test or loaded from before the column existed) to 1, the original feature
// definitions, so the stored value is always a real version.
func featureVersionOrDefault(v int) int {
	if v < 1 {
		return 1
	}
	return v
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
  feature_version,
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
  mean_type_token_ratio,
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
  std_type_token_ratio,
  document_count,
  sentence_count,
  character_count,
  mean_lexical_frequencies,
  std_lexical_frequencies,
  mean_char_ngrams,
  std_char_ngrams,
  mean_pos_ngrams,
  std_pos_ngrams,
  sources
FROM profile
WHERE id = 1
`)

	var trainedAt string
	var meanLexicalJSON string
	var stdLexicalJSON string
	var meanNgramJSON string
	var stdNgramJSON string
	var meanPOSNgramJSON string
	var stdPOSNgramJSON string
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
		&record.FeatureVersion,
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
		&mean.TypeTokenRatio,
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
		&std.TypeTokenRatio,
		&dist.DocumentCount,
		&dist.SentenceCount,
		&dist.CharacterCount,
		&meanLexicalJSON,
		&stdLexicalJSON,
		&meanNgramJSON,
		&stdNgramJSON,
		&meanPOSNgramJSON,
		&stdPOSNgramJSON,
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
	mean.POSNgrams = unmarshalLexical(meanPOSNgramJSON)
	std.POSNgrams = unmarshalLexical(stdPOSNgramJSON)
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
