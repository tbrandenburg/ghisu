# ghisu

A terminal Kanban board for GitHub issues, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Columns, their names, and their GitHub search queries are **fully configurable**.

## Preview

```
╭─────────────────────────────────────────────────────────────────────────╮
│  ghisu — tbrandenburg/ghisu                    last refresh: 14:03:22   │
├──────────────────┬──────────────────┬──────────────────┬────────────────┤
│ Todo (3)         │ Bugs (1)         │ Doing (1)        │ Done (1)       │
├──────────────────┼──────────────────┼──────────────────┼────────────────┤
│                  │                  │                  │                │
│ > #4 Another     │   #2 Sample bug  │   #3 Sample      │   #1 Sample    │
│     todo task    │      issue       │      in-progress │      todo      │
│                  │                  │      issue       │      issue     │
│   #3 Sample      │                  │                  │                │
│      in-progress │                  │                  │                │
│      issue       │                  │                  │                │
│                  │                  │                  │                │
│   #2 Sample bug  │                  │                  │                │
│      issue       │                  │                  │                │
│                  │                  │                  │                │
├──────────────────┴──────────────────┴──────────────────┴────────────────┤
│  h/l: col   j/k: item   r: refresh   q: quit                            │
╰─────────────────────────────────────────────────────────────────────────╯
```

> Active column and selected item are highlighted in the terminal.
> Column names and queries shown above are the built-in defaults — all fully configurable.

## Features

- **Fully configurable columns** — any name, any [GitHub issue search query](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests)
- Four built-in default columns: **Todo**, **Bugs**, **Doing**, **Done**
- Automatic refresh every 10 seconds
- Full keyboard navigation
- Works with any GitHub repository via the `gh` CLI
- GitHub Enterprise support via `--hostname`

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
# Current directory's GitHub repo — columns from .ghisu/config.json or defaults
ghisu

# Explicit repo
ghisu --repo owner/repo

# GitHub Enterprise
ghisu --repo owner/repo --hostname github.example.com

# Explicit config file with custom columns
ghisu --config /path/to/config.json
```

CLI flags (`--repo`, `--hostname`) always override values in the config file.

## Configuration

Columns are fully configurable via a JSON file. ghisu looks for a config file
in the following order, using the first one found:

1. `.ghisu/config.json` in the **current directory** — project-local config
2. `$XDG_CONFIG_HOME/ghisu/config.json` (or `~/.config/ghisu/config.json`) — user-global config

The `--config` flag overrides both.

### Config file format

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
| `columns`  | No       | **Custom columns.** Any number, any name, any query. Omit to use the four built-in defaults. |

Each column requires:

| Field   | Required | Description |
|---------|----------|-------------|
| `name`  | Yes      | Display label shown in the column header. |
| `query` | Yes      | Any valid [GitHub issue search query](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests). |

### Example: custom columns for a project with priority labels

```json
{
  "repo": "owner/repo",
  "columns": [
    { "name": "Critical", "query": "is:open label:priority/high" },
    { "name": "Backlog",  "query": "is:open label:priority/medium" },
    { "name": "Review",   "query": "is:open label:in-review" },
    { "name": "Done",     "query": "is:closed" }
  ]
}
```

### Project-local config

Drop a `.ghisu/config.json` into any repository and ghisu will use it
automatically when run from that directory — no flags needed:

```
my-repo/
├── .ghisu/
│   └── config.json   ← picked up automatically
└── ...
```

## Keyboard shortcuts

| Key       | Action                   |
|-----------|--------------------------|
| `h` / `l` | Switch column left/right |
| `j` / `k` | Move cursor up/down      |
| `r`       | Manual refresh           |
| `q`       | Quit                     |

## Default columns

When no config file is present the following columns are used. All are
**overridable** via the config file.

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
