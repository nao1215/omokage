package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/nao1215/dyer/internal/feature"
	"github.com/nao1215/dyer/internal/profile"
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
  document_count INTEGER NOT NULL,
  sentence_count INTEGER NOT NULL,
  character_count INTEGER NOT NULL
);`

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

	const query = `
INSERT INTO profile (
  id, author, source_dir, trained_at, file_count,
  mean_average_sentence_length, mean_sentence_length_variance, mean_punctuation_frequency,
  mean_newline_frequency, mean_bullet_ratio, mean_conjunction_frequency,
  mean_kanji_ratio, mean_hiragana_ratio, mean_katakana_ratio,
  mean_paragraph_length_variance, mean_markdown_structure_density,
  std_average_sentence_length, std_sentence_length_variance, std_punctuation_frequency,
  std_newline_frequency, std_bullet_ratio, std_conjunction_frequency,
  std_kanji_ratio, std_hiragana_ratio, std_katakana_ratio,
  std_paragraph_length_variance, std_markdown_structure_density,
  document_count, sentence_count, character_count
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
  document_count = excluded.document_count,
  sentence_count = excluded.sentence_count,
  character_count = excluded.character_count;
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
		record.Distribution.DocumentCount,
		record.Distribution.SentenceCount,
		record.Distribution.CharacterCount,
	)
	return err
}

func LoadProfile(path string) (profile.Record, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return profile.Record{}, err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, schema); err != nil {
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
  document_count,
  sentence_count,
  character_count
FROM profile
WHERE id = 1
`)

	var trainedAt string
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
		&dist.DocumentCount,
		&dist.SentenceCount,
		&dist.CharacterCount,
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
	dist.Mean = mean
	dist.StdDev = std
	record.Distribution = dist
	return record, nil
}
