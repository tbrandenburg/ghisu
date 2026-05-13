// Package config loads ghisu column configuration from a JSON file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Column defines a kanban column with a display name and a GitHub search query.
type Column struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Config is the top-level configuration structure.
type Config struct {
	Columns []Column `json:"columns"`
}

// DefaultPath returns the default config file location:
// $XDG_CONFIG_HOME/ghisu/config.json or ~/.config/ghisu/config.json.
func DefaultPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ghisu", "config.json")
}

// Load reads the config file at path. If path is empty or the file does not
// exist, an empty Config is returned with no error (caller uses defaults).
func Load(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config %s: %w", path, err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	for i, col := range c.Columns {
		if col.Name == "" {
			return fmt.Errorf("column[%d]: name is required", i)
		}
		if col.Query == "" {
			return fmt.Errorf("column[%d] %q: query is required", i, col.Name)
		}
	}
	return nil
}
