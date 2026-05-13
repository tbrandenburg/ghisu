package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tbrandenburg/ghisu/internal/config"
	"github.com/tbrandenburg/ghisu/internal/model"
)

func main() {
	repo := flag.String("repo", "", "GitHub repo (owner/name). Defaults to current directory repo.")
	cfgPath := flag.String("config", config.DefaultPath(), "Path to JSON config file.")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cols := toModelColumns(cfg.Columns)
	if len(cols) == 0 {
		cols = model.DefaultColumns()
	}

	m := model.New(*repo, cols)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// toModelColumns converts config columns to model columns.
func toModelColumns(cols []config.Column) []model.Column {
	out := make([]model.Column, len(cols))
	for i, c := range cols {
		out[i] = model.Column{Name: c.Name, Query: c.Query}
	}
	return out
}
