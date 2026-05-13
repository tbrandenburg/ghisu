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

func TestLoadRepoAndHostname(t *testing.T) {
	path := writeTemp(t, `{
		"repo": "myorg/myrepo",
		"hostname": "github.example.com",
		"columns": [
			{"name": "Open", "query": "is:open"}
		]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Repo != "myorg/myrepo" {
		t.Errorf("Repo: want %q, got %q", "myorg/myrepo", cfg.Repo)
	}
	if cfg.Hostname != "github.example.com" {
		t.Errorf("Hostname: want %q, got %q", "github.example.com", cfg.Hostname)
	}
}

func TestLoadHostnameWithSchemeRejected(t *testing.T) {
	path := writeTemp(t, `{
		"hostname": "https://github.example.com",
		"columns": [{"name": "Open", "query": "is:open"}]
	}`)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for hostname with scheme, got nil")
	}
}

func TestLoadHostnameWithPathRejected(t *testing.T) {
	path := writeTemp(t, `{
		"hostname": "github.example.com/extra",
		"columns": [{"name": "Open", "query": "is:open"}]
	}`)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for hostname with path, got nil")
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
	path := ResolvePath()
	if path == "" {
		t.Error("ResolvePath() returned empty string")
	}
	// must end with ghisu/config.json
	if filepath.Base(filepath.Dir(path)) != "ghisu" || filepath.Base(path) != "config.json" {
		t.Errorf("ResolvePath() = %q, expected to end with ghisu/config.json", path)
	}
}

func TestResolvePathPrefersLocal(t *testing.T) {
	// Create a temporary directory to act as the working directory.
	dir := t.TempDir()
	localDir := filepath.Join(dir, ".ghisu")
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	localCfg := filepath.Join(localDir, "config.json")
	if err := os.WriteFile(localCfg, []byte(`{"columns":[]}`), 0o644); err != nil {
		t.Fatalf("write local config: %v", err)
	}

	// Change working directory to the temp dir for this test.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	got := ResolvePath()
	if got != localCfg {
		t.Errorf("ResolvePath() = %q, want local path %q", got, localCfg)
	}
}

func TestResolvePathFallsBackToGlobal(t *testing.T) {
	// Use a temp dir with no .ghisu/config.json.
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	got := ResolvePath()
	global := globalPath()
	if got != global {
		t.Errorf("ResolvePath() = %q, want global path %q", got, global)
	}
}
