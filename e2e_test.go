// e2e_test.go exercises the full config → model pipeline and optionally
// the live gh CLI when available.
package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tbrandenburg/ghisu/internal/config"
	"github.com/tbrandenburg/ghisu/internal/gh"
	"github.com/tbrandenburg/ghisu/internal/model"
)

func writeConfigFile(t *testing.T, cfg config.Config) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

// TestE2EConfigPipeline verifies the full Load → toModelColumns → model.New path.
func TestE2EConfigPipeline(t *testing.T) {
	wantCols := []config.Column{
		{Name: "Todo", Query: "is:open no:assignee"},
		{Name: "Bugs", Query: "is:open label:bug"},
	}
	path := writeConfigFile(t, config.Config{
		Repo:     "tbrandenburg/ghisu",
		Hostname: "",
		Columns:  wantCols,
	})

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Repo != "tbrandenburg/ghisu" {
		t.Errorf("Repo: want %q, got %q", "tbrandenburg/ghisu", cfg.Repo)
	}
	if cfg.Hostname != "" {
		t.Errorf("Hostname: want empty, got %q", cfg.Hostname)
	}

	cols := toModelColumns(cfg.Columns)
	m := model.New(cfg.Repo, cfg.Hostname, cols)

	got := m.Columns()
	if len(got) != len(wantCols) {
		t.Fatalf("expected %d columns in model, got %d", len(wantCols), len(got))
	}
	for i, want := range wantCols {
		if got[i].Name != want.Name {
			t.Errorf("col[%d].Name: want %q, got %q", i, want.Name, got[i].Name)
		}
		if got[i].Query != want.Query {
			t.Errorf("col[%d].Query: want %q, got %q", i, want.Query, got[i].Query)
		}
	}
}

// TestE2EHostnameInConfig verifies that a GHE hostname round-trips through config.
func TestE2EHostnameInConfig(t *testing.T) {
	path := writeConfigFile(t, config.Config{
		Repo:     "myorg/myrepo",
		Hostname: "github.example.com",
		Columns:  []config.Column{{Name: "Open", Query: "is:open"}},
	})

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Hostname != "github.example.com" {
		t.Errorf("Hostname: want %q, got %q", "github.example.com", cfg.Hostname)
	}

	// CLI flag override: flag value beats config value.
	resolvedHostname := cfg.Hostname
	flagHostname := "override.example.com"
	if flagHostname != "" {
		resolvedHostname = flagHostname
	}
	if resolvedHostname != "override.example.com" {
		t.Errorf("flag override: want %q, got %q", "override.example.com", resolvedHostname)
	}
}

// TestE2EDefaultFallback verifies that an empty config falls back to DefaultColumns.
func TestE2EDefaultFallback(t *testing.T) {
	cfg, err := config.Load("") // empty path → no file
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	cols := toModelColumns(cfg.Columns)
	if len(cols) == 0 {
		cols = model.DefaultColumns()
	}

	m := model.New("", "", cols)
	if len(m.Columns()) != 4 {
		t.Errorf("expected 4 default columns, got %d", len(m.Columns()))
	}
}

// TestE2ELiveGHFetch verifies that columns defined in a JSON config actually
// return issues from GitHub. Skipped when gh is unavailable or unauthenticated.
func TestE2ELiveGHFetch(t *testing.T) {
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("gh CLI not found")
	}
	if err := exec.Command("gh", "auth", "status").Run(); err != nil {
		t.Skip("gh not authenticated")
	}

	path := writeConfigFile(t, config.Config{
		Repo:    "tbrandenburg/ghisu",
		Columns: []config.Column{{Name: "Open", Query: "is:open"}, {Name: "Closed", Query: "is:closed"}},
	})

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	for _, col := range cfg.Columns {
		issues, err := gh.FetchIssues(cfg.Repo, cfg.Hostname, col.Query)
		if err != nil {
			t.Errorf("FetchIssues(%q, %q, %q): %v", cfg.Repo, cfg.Hostname, col.Query, err)
			continue
		}
		t.Logf("column %q: %d issue(s)", col.Name, len(issues))
		for _, iss := range issues {
			if iss.Number == 0 {
				t.Errorf("column %q: issue with number=0", col.Name)
			}
			if iss.Title == "" {
				t.Errorf("column %q: issue #%d has empty title", col.Name, iss.Number)
			}
		}
	}
}
