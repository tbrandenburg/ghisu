// Package gh wraps the GitHub CLI to fetch issues.
package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Issue represents a GitHub issue returned by the CLI.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

// FetchIssues runs `gh issue list` with the given search query and returns results.
func FetchIssues(repo, query string) ([]Issue, error) {
	args := []string{
		"issue", "list",
		"--search", query,
		"--limit", "50",
		"--json", "number,title,state,labels",
	}
	if repo != "" {
		args = append([]string{"--repo", repo}, args...)
	}

	out, err := exec.Command("gh", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("gh issue list: %w: %s", err, strings.TrimSpace(string(out)))
	}

	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("parse issues: %w", err)
	}
	return issues, nil
}
