package model

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tbrandenburg/ghisu/internal/gh"
)

func TestNewModel(t *testing.T) {
	cols := DefaultColumns()
	m := New("", "", cols)

	if len(m.columns) != len(cols) {
		t.Fatalf("expected %d columns, got %d", len(cols), len(m.columns))
	}
	if len(m.issues) != len(cols) {
		t.Fatalf("expected %d issue slices, got %d", len(cols), len(m.issues))
	}
	if m.activeCol != 0 {
		t.Errorf("expected activeCol=0, got %d", m.activeCol)
	}
}

func TestDefaultColumns(t *testing.T) {
	cols := DefaultColumns()
	if len(cols) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(cols))
	}
	names := []string{"Todo", "Bugs", "Doing", "Done"}
	for i, want := range names {
		if cols[i].Name != want {
			t.Errorf("col[%d]: want %q, got %q", i, want, cols[i].Name)
		}
	}
}

func TestKeyNavigation(t *testing.T) {
	cols := DefaultColumns()
	m := New("", "", cols)
	// pre-populate issues so cursor movement works
	m.issues[0] = []gh.Issue{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}

	// move right
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	nm := next.(Model)
	if nm.activeCol != 1 {
		t.Errorf("after l: want col 1, got %d", nm.activeCol)
	}

	// move left back
	next, _ = nm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	nm = next.(Model)
	if nm.activeCol != 0 {
		t.Errorf("after h: want col 0, got %d", nm.activeCol)
	}

	// move cursor down
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	nm = next.(Model)
	if nm.activeCur[0] != 1 {
		t.Errorf("after j: want cursor 1, got %d", nm.activeCur[0])
	}

	// move cursor up
	next, _ = nm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	nm = next.(Model)
	if nm.activeCur[0] != 0 {
		t.Errorf("after k: want cursor 0, got %d", nm.activeCur[0])
	}
}

func TestKeyNavBoundaries(t *testing.T) {
	cols := DefaultColumns()
	m := New("", "", cols)

	// can't go left from first column
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	nm := next.(Model)
	if nm.activeCol != 0 {
		t.Errorf("want col 0, got %d", nm.activeCol)
	}

	// navigate to last column
	m.activeCol = len(cols) - 1
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	nm = next.(Model)
	if nm.activeCol != len(cols)-1 {
		t.Errorf("want col %d, got %d", len(cols)-1, nm.activeCol)
	}
}

func TestFetchedMsg(t *testing.T) {
	m := New("", "", DefaultColumns())
	issues := []gh.Issue{{Number: 5, Title: "hello"}}

	next, _ := m.Update(fetchedMsg{col: 0, issues: issues})
	nm := next.(Model)

	if len(nm.issues[0]) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(nm.issues[0]))
	}
	if nm.issues[0][0].Number != 5 {
		t.Errorf("expected issue #5, got #%d", nm.issues[0][0].Number)
	}
	if nm.loading[0] {
		t.Error("loading should be false after fetchedMsg")
	}
}

func TestTickUpdatesRefresh(t *testing.T) {
	m := New("", "", DefaultColumns())
	tick := tickMsg(time.Now())

	next, _ := m.Update(tick)
	nm := next.(Model)

	if nm.lastRefresh.IsZero() {
		t.Error("lastRefresh should be set after tick")
	}
}

func TestView(t *testing.T) {
	m := New("", "", DefaultColumns())
	m.issues[0] = []gh.Issue{{Number: 1, Title: "test issue"}}

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
	if !contains(view, "Todo") {
		t.Error("View() should contain column name 'Todo'")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"short", 10, "short"},
		{"exactly ten!", 12, "exactly ten!"},
		{"this is a very long title", 10, "this is a…"},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.max)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
