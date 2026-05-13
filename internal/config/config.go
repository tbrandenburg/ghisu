// Package config loads ghisu column configuration from a JSON file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Column defines a kanban column with a display name and a GitHub search query.
type Column struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Config is the top-level configuration structure.
type Config struct {
	// Repo is the GitHub repository in "owner/repo" format.
	// When empty the gh CLI resolves the repo from the current directory.
	Repo string `json:"repo"`

	// Hostname is the GitHub hostname (e.g. "github.example.com" for GHE).
	// When empty github.com is used (no --hostname flag is passed to gh).
	Hostname string `json:"hostname"`

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
	if c.Hostname != "" && strings.ContainsAny(c.Hostname, "/:") {
		return fmt.Errorf("hostname must be a bare hostname without scheme or path (got %q)", c.Hostname)
	}
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
