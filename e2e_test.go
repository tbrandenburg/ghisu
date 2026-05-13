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

func writeConfigFile(t *testing.T, cols []config.Column) string {
	t.Helper()
	data, err := json.Marshal(config.Config{Columns: cols})
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
	path := writeConfigFile(t, wantCols)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	cols := toModelColumns(cfg.Columns)
	m := model.New("tbrandenburg/ghisu", cols)

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

	m := model.New("", cols)
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

	wantCols := []config.Column{
		{Name: "Open", Query: "is:open"},
		{Name: "Closed", Query: "is:closed"},
	}
	path := writeConfigFile(t, wantCols)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	const repo = "tbrandenburg/ghisu"

	// Fetch each column and verify we get slices back (content may vary).
	for _, col := range cfg.Columns {
		issues, err := gh.FetchIssues(repo, col.Query)
		if err != nil {
			t.Errorf("FetchIssues(%q, %q): %v", repo, col.Query, err)
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
