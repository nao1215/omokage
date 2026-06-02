package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nao1215/omokage/internal/config"
)

const ConfigFileName = "omokage.toml"

var ErrProjectNotFound = errors.New("omokage project not found")

// ErrStoreNotFound reports that a requested store does not exist yet: --global
// with no global store created, or a global fallback whose home is unknown.
// Commands that only need a feature set rather than a profile — diff — treat it,
// like ErrProjectNotFound, as a cue to fall back to the built-in defaults instead
// of failing.
var ErrStoreNotFound = errors.New("omokage store not found")

func Init(root string, name string) (config.Config, error) {
	configPath := filepath.Join(root, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return config.Config{}, fmt.Errorf("%s already exists", configPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return config.Config{}, err
	}

	cfg := config.Default(strings.TrimSpace(name))
	if cfg.Project.Name == "" {
		cfg.Project.Name = filepath.Base(root)
	}

	if err := config.Save(configPath, cfg); err != nil {
		return config.Config{}, err
	}
	if err := os.MkdirAll(filepath.Join(root, cfg.Storage.ProfileDir), 0o750); err != nil {
		return config.Config{}, err
	}
	if err := os.MkdirAll(filepath.Join(root, cfg.Storage.CacheDir), 0o750); err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

func Load(start string) (string, config.Config, error) {
	root, err := FindRoot(start)
	if err != nil {
		return "", config.Config{}, err
	}
	cfg, err := config.Load(filepath.Join(root, ConfigFileName))
	if err != nil {
		return "", config.Config{}, err
	}
	return root, cfg, nil
}

func LoadOptional(start string) (config.Config, bool, error) {
	root, err := FindRoot(start)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return config.Config{}, false, nil
		}
		return config.Config{}, false, err
	}
	cfg, err := config.Load(filepath.Join(root, ConfigFileName))
	if err != nil {
		return config.Config{}, false, err
	}
	return cfg, true, nil
}

func FindRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(current, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", ErrProjectNotFound
		}
		current = parent
	}
}

func ProfilePath(root string, cfg config.Config, author string) (string, error) {
	safeAuthor := strings.TrimSpace(author)
	if safeAuthor == "" {
		return "", errors.New("author must not be empty")
	}
	if safeAuthor != filepath.Base(safeAuthor) {
		return "", errors.New("author must not contain path separators")
	}
	return filepath.Join(root, cfg.Storage.ProfileDir, safeAuthor+".db"), nil
}

func ListProfiles(root string, cfg config.Config) ([]string, error) {
	return listProfilesInDir(filepath.Join(root, cfg.Storage.ProfileDir))
}
