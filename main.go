package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tbrandenburg/ghisu/internal/model"
)

func main() {
	repo := flag.String("repo", "", "GitHub repo (owner/name). Defaults to current directory repo.")
	flag.Parse()

	m := model.New(*repo, model.DefaultColumns())

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
