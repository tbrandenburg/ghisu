// Package model contains the Bubble Tea model for ghisu.
package model

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tbrandenburg/ghisu/internal/gh"
)

const (
	refreshInterval = 10 * time.Second
	colWidth        = 28
	colHeight       = 20
)

// Column defines a kanban column with a display name and a GitHub search query.
type Column struct {
	Name  string
	Query string
}

// DefaultColumns returns the standard four-column board.
func DefaultColumns() []Column {
	return []Column{
		{Name: "Todo", Query: "is:open no:assignee"},
		{Name: "Bugs", Query: "is:open label:bug"},
		{Name: "Doing", Query: "is:open label:in-progress"},
		{Name: "Done", Query: "is:closed"},
	}
}

// tickMsg is sent on each refresh tick.
type tickMsg time.Time

// fetchedMsg carries new issues for a column index.
type fetchedMsg struct {
	col    int
	issues []gh.Issue
	err    error
}

// Model is the Bubble Tea application model.
type Model struct {
	columns     []Column
	issues      [][]gh.Issue
	errors      []error
	activeCol   int
	activeCur   []int
	repo        string
	hostname    string
	loading     []bool
	lastRefresh time.Time
}

// Columns returns the model's column definitions.
func (m Model) Columns() []Column {
	return m.columns
}

// New creates a new Model for the given repo and hostname.
// repo empty = current dir repo; hostname empty = github.com.
func New(repo, hostname string, columns []Column) Model {
	n := len(columns)
	return Model{
		columns:   columns,
		issues:    make([][]gh.Issue, n),
		errors:    make([]error, n),
		activeCur: make([]int, n),
		loading:   make([]bool, n),
		repo:      repo,
		hostname:  hostname,
	}
}

func tick() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchColumn(repo, hostname string, idx int, query string) tea.Cmd {
	return func() tea.Msg {
		issues, err := gh.FetchIssues(repo, hostname, query)
		return fetchedMsg{col: idx, issues: issues, err: err}
	}
}

// Init starts the first fetch and the refresh timer.
func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.columns)+1)
	for i, col := range m.columns {
		m.loading[i] = true
		cmds = append(cmds, fetchColumn(m.repo, m.hostname, i, col.Query))
	}
	cmds = append(cmds, tick())
	return tea.Batch(cmds...)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tickMsg:
		m.lastRefresh = time.Time(msg)
		cmds := make([]tea.Cmd, len(m.columns))
		for i, col := range m.columns {
			m.loading[i] = true
			cmds[i] = fetchColumn(m.repo, m.hostname, i, col.Query)
		}
		return m, tea.Batch(append(cmds, tick())...)

	case fetchedMsg:
		m.loading[msg.col] = false
		m.issues[msg.col] = msg.issues
		m.errors[msg.col] = msg.err
		// clamp cursor
		if m.activeCur[msg.col] >= len(msg.issues) && len(msg.issues) > 0 {
			m.activeCur[msg.col] = len(msg.issues) - 1
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	n := len(m.columns)
	col := m.activeCol
	cur := m.activeCur[col]
	issues := m.issues[col]

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "h", "left":
		if m.activeCol > 0 {
			m.activeCol--
		}
	case "l", "right":
		if m.activeCol < n-1 {
			m.activeCol++
		}
	case "j", "down":
		if cur < len(issues)-1 {
			m.activeCur[col] = cur + 1
		}
	case "k", "up":
		if cur > 0 {
			m.activeCur[col] = cur - 1
		}
	case "r":
		cmds := make([]tea.Cmd, len(m.columns))
		for i, c := range m.columns {
			m.loading[i] = true
			cmds[i] = fetchColumn(m.repo, m.hostname, i, c.Query)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View renders the board.
func (m Model) View() string {
	cols := make([]string, len(m.columns))
	for i, col := range m.columns {
		cols[i] = m.renderColumn(i, col)
	}

	board := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
	status := m.statusBar()
	return lipgloss.JoinVertical(lipgloss.Left, board, status)
}

var (
	activeColStyle = lipgloss.NewStyle().
			Width(colWidth).
			Height(colHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(0, 1)

	inactiveColStyle = lipgloss.NewStyle().
				Width(colWidth).
				Height(colHeight).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(0, 1)

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))

	inactiveTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("240"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true)

	itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)
)

func (m Model) renderColumn(idx int, col Column) string {
	active := idx == m.activeCol

	var sb strings.Builder

	// title
	ts := inactiveTitleStyle
	if active {
		ts = titleStyle
	}
	count := len(m.issues[idx])
	header := fmt.Sprintf("%s (%d)", col.Name, count)
	if m.loading[idx] {
		header += " …"
	}
	sb.WriteString(ts.Render(header))
	sb.WriteString("\n")

	if m.errors[idx] != nil {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("error: " + m.errors[idx].Error()))
	} else {
		for i, issue := range m.issues[idx] {
			line := fmt.Sprintf("#%d %s", issue.Number, truncate(issue.Title, colWidth-4))
			if active && i == m.activeCur[idx] {
				sb.WriteString(selectedItemStyle.Render("> " + line))
			} else {
				sb.WriteString(itemStyle.Render("  " + line))
			}
			sb.WriteString("\n")
		}
	}

	style := inactiveColStyle
	if active {
		style = activeColStyle
	}
	return style.Render(sb.String())
}

func (m Model) statusBar() string {
	refresh := "never"
	if !m.lastRefresh.IsZero() {
		refresh = m.lastRefresh.Format("15:04:05")
	}
	help := "h/l: col  j/k: item  r: refresh  q: quit"
	return statusStyle.Render(fmt.Sprintf("last refresh: %s  |  %s", refresh, help))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
