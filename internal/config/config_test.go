package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "ghisu-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadValidConfig(t *testing.T) {
	path := writeTemp(t, `{
		"columns": [
			{"name": "Todo",  "query": "is:open no:assignee"},
			{"name": "Bugs",  "query": "is:open label:bug"}
		]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cfg.Columns))
	}
	if cfg.Columns[0].Name != "Todo" {
		t.Errorf("col[0].Name: want Todo, got %q", cfg.Columns[0].Name)
	}
	if cfg.Columns[1].Query != "is:open label:bug" {
		t.Errorf("col[1].Query: want 'is:open label:bug', got %q", cfg.Columns[1].Query)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Columns) != 0 {
		t.Errorf("expected 0 columns for empty path, got %d", len(cfg.Columns))
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatalf("missing file should not error, got: %v", err)
	}
	if len(cfg.Columns) != 0 {
		t.Errorf("expected 0 columns for missing file, got %d", len(cfg.Columns))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := writeTemp(t, `{ not valid json }`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoadMissingName(t *testing.T) {
	path := writeTemp(t, `{"columns": [{"name": "", "query": "is:open"}]}`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing column name, got nil")
	}
}

func TestLoadMissingQuery(t *testing.T) {
	path := writeTemp(t, `{"columns": [{"name": "Todo", "query": ""}]}`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing column query, got nil")
	}
}

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()
	if path == "" {
		t.Error("DefaultPath() returned empty string")
	}
	// must end with the expected suffix
	want := filepath.Join("ghisu", "config.json")
	if filepath.Base(filepath.Dir(path)) != "ghisu" || filepath.Base(path) != "config.json" {
		t.Errorf("DefaultPath() = %q, expected to end with %q", path, want)
	}
}
