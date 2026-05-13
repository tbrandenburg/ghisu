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
	repo := flag.String("repo", "", "GitHub repo (owner/name). Overrides config. Defaults to current directory repo.")
	hostname := flag.String("hostname", "", "GitHub hostname for GHE (e.g. github.example.com). Overrides config.")
	cfgPath := flag.String("config", config.ResolvePath(), "Path to JSON config file.")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// CLI flags take precedence over config values.
	resolvedRepo := cfg.Repo
	if *repo != "" {
		resolvedRepo = *repo
	}

	resolvedHostname := cfg.Hostname
	if *hostname != "" {
		resolvedHostname = *hostname
	}

	cols := toModelColumns(cfg.Columns)
	if len(cols) == 0 {
		cols = model.DefaultColumns()
	}

	m := model.New(resolvedRepo, resolvedHostname, cols)

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
