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
# Current directory's GitHub repo
ghisu

# Explicit repo
ghisu --repo owner/repo
```

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
