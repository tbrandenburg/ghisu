# ghisu

A terminal Kanban board for GitHub issues, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Features

- Four-column board: **Todo**, **Bugs**, **Doing**, **Done**
- Automatic refresh every 10 seconds
- Full keyboard navigation
- Works with any GitHub repository via the `gh` CLI

## Requirements

- Go 1.24+
- [GitHub CLI (`gh`)](https://cli.github.com/) — authenticated (`gh auth login`)

## Installation

```bash
go install github.com/tbrandenburg/ghisu@latest
```

Or build from source:

```bash
make build
```

## Usage

```bash
# Current directory's GitHub repo (github.com)
ghisu

# Explicit repo
ghisu --repo owner/repo

# GitHub Enterprise
ghisu --repo owner/repo --hostname github.example.com

# Custom column config
ghisu --config ~/.config/ghisu/config.json
```

CLI flags (`--repo`, `--hostname`) always override values in the config file.

## Configuration

Columns, repo, and hostname can be defined in a JSON file. The default location
is `$XDG_CONFIG_HOME/ghisu/config.json` (or `~/.config/ghisu/config.json`).

```json
{
  "repo": "owner/repo",
  "hostname": "github.example.com",
  "columns": [
    { "name": "Todo",  "query": "is:open no:assignee" },
    { "name": "Bugs",  "query": "is:open label:bug" },
    { "name": "Doing", "query": "is:open label:in-progress" },
    { "name": "Done",  "query": "is:closed" }
  ]
}
```

| Field      | Required | Description |
|------------|----------|-------------|
| `repo`     | No       | `owner/repo`. Omit to use the current directory's repo. |
| `hostname` | No       | Bare hostname for GitHub Enterprise (e.g. `github.example.com`). Omit for github.com. Must not include a scheme (`https://`) or path. |
| `columns`  | No       | Custom columns. Omit to use the four built-in defaults. |

Each column entry requires a `name` (display label) and a `query` (any valid
[GitHub issue search query](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests)).

## Keyboard shortcuts

| Key       | Action                  |
|-----------|-------------------------|
| `h` / `l` | Switch column left/right |
| `j` / `k` | Move cursor up/down      |
| `r`       | Manual refresh           |
| `q`       | Quit                     |

## Columns

| Column | GitHub search query           |
|--------|-------------------------------|
| Todo   | `is:open no:assignee`         |
| Bugs   | `is:open label:bug`           |
| Doing  | `is:open label:in-progress`   |
| Done   | `is:closed`                   |

## Development

```bash
make install   # install dependencies
make build     # build binary
make test      # run tests
make lint      # format and vet
make run       # build and run
make clean     # remove build artifacts
```
