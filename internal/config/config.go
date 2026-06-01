package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Project  Project
	Features Features
	Storage  Storage
}

type Project struct {
	Name string
}

type Features struct {
	SentenceLength           bool
	SentenceLengthVariance   bool
	PunctuationFrequency     bool
	NewlineFrequency         bool
	BulletRatio              bool
	KanjiRatio               bool
	HiraganaRatio            bool
	KatakanaRatio            bool
	ConjunctionFrequency     bool
	ParagraphLengthVariance  bool
	MarkdownStructureDensity bool
}

type Storage struct {
	ProfileDir string
	CacheDir   string
}

func Default(projectName string) Config {
	return Config{
		Project: Project{
			Name: projectName,
		},
		Features: Features{
			SentenceLength:           true,
			SentenceLengthVariance:   true,
			PunctuationFrequency:     true,
			NewlineFrequency:         true,
			BulletRatio:              true,
			KanjiRatio:               true,
			HiraganaRatio:            true,
			KatakanaRatio:            true,
			ConjunctionFrequency:     true,
			ParagraphLengthVariance:  true,
			MarkdownStructureDensity: true,
		},
		Storage: Storage{
			ProfileDir: "./profiles",
			CacheDir:   "./cache",
		},
	}
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Config{}, err
	}
	return Parse(data)
}

func Parse(data []byte) (Config, error) {
	cfg := Default("dyer-project")
	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return Config{}, fmt.Errorf("invalid config line: %s", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch section {
		case "project":
			if key == "name" {
				parsed, err := parseString(value)
				if err != nil {
					return Config{}, err
				}
				cfg.Project.Name = parsed
			}
		case "features":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("invalid boolean for %s: %w", key, err)
			}
			switch key {
			case "sentence_length":
				cfg.Features.SentenceLength = parsed
			case "sentence_length_variance":
				cfg.Features.SentenceLengthVariance = parsed
			case "punctuation_frequency":
				cfg.Features.PunctuationFrequency = parsed
			case "newline_frequency":
				cfg.Features.NewlineFrequency = parsed
			case "bullet_ratio":
				cfg.Features.BulletRatio = parsed
			case "kanji_ratio":
				cfg.Features.KanjiRatio = parsed
			case "hiragana_ratio":
				cfg.Features.HiraganaRatio = parsed
			case "katakana_ratio":
				cfg.Features.KatakanaRatio = parsed
			case "conjunction_frequency":
				cfg.Features.ConjunctionFrequency = parsed
			case "paragraph_length_variance":
				cfg.Features.ParagraphLengthVariance = parsed
			case "markdown_structure_frequency":
				cfg.Features.MarkdownStructureDensity = parsed
			}
		case "storage":
			parsed, err := parseString(value)
			if err != nil {
				return Config{}, err
			}
			switch key {
			case "profile_dir":
				cfg.Storage.ProfileDir = parsed
			case "cache_dir":
				cfg.Storage.CacheDir = parsed
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Save(path string, cfg Config) error {
	return os.WriteFile(filepath.Clean(path), []byte(cfg.String()), 0o600)
}

func (c Config) String() string {
	return fmt.Sprintf(`[project]
name = %q

[features]
sentence_length = %t
sentence_length_variance = %t
punctuation_frequency = %t
newline_frequency = %t
bullet_ratio = %t
kanji_ratio = %t
hiragana_ratio = %t
katakana_ratio = %t
conjunction_frequency = %t
paragraph_length_variance = %t
markdown_structure_frequency = %t

[storage]
profile_dir = %q
cache_dir = %q
`,
		c.Project.Name,
		c.Features.SentenceLength,
		c.Features.SentenceLengthVariance,
		c.Features.PunctuationFrequency,
		c.Features.NewlineFrequency,
		c.Features.BulletRatio,
		c.Features.KanjiRatio,
		c.Features.HiraganaRatio,
		c.Features.KatakanaRatio,
		c.Features.ConjunctionFrequency,
		c.Features.ParagraphLengthVariance,
		c.Features.MarkdownStructureDensity,
		c.Storage.ProfileDir,
		c.Storage.CacheDir,
	)
}

func parseString(raw string) (string, error) {
	if strings.HasPrefix(raw, "\"") {
		return strconv.Unquote(raw)
	}
	return raw, nil
}
