# Vex - Excel Viewer

A TUI Excel and CSV viewer in Go 

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.1.0-brightgreen.svg)](https://github.com/odesaur/vex-tui/releases)

![Vex Demo](assets/vex-demo.gif)


## Installation

### Using go install

```bash
GOPROXY=direct go install github.com/odesaur/vex-tui@latest
```

### From source

```bash
# Clone the repository
git clone https://github.com/odesaur/vex-tui.git
cd vex-tui

# Install dependencies
go mod download

# Build locally
go build -o vex-tui .
```

### One-command install via Make

```bash
git clone https://github.com/odesaur/vex-tui.git
cd vex-tui
make install
```

## Usage

```bash
# Basic usage (opens file directly)
vex-tui data.xlsx

# Start with a picker in the current directory
vex-tui

# With a specific theme
vex-tui report.csv --theme nord

# Short flag
vex-tui sales.xlsx -t tokyo-night
```

## Keyboard Shortcuts

Press `?` in the app anytime to open the built-in help with the full list. Quick reference:

| Keys                   | Action                        |
| ---------------------- | ----------------------------- |
| `↑/↓/←/→` or `hjkl`    | Move cursor                   |
| `PageUp/PageDown`      | Page scroll                   |
| `Tab` / `Shift+Tab`    | Next / previous sheet         |
| `/`, `n`, `N`          | Search, next result, previous |
| `Ctrl+G`               | Jump to cell                  |
| `Enter`                | Cell details                  |
| `c` / `C`              | Copy cell / row               |
| `f`                    | Toggle formulas               |
| `e`                    | Export sheet                  |
| `t`                    | Theme selector                |
| `q` or `Ctrl+C`        | Quit                          |

## Acknowledgments

- Built with the amazing [Charm](https://charm.sh/) ecosystem
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
  - [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
  - [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- Excel parsing by [Excelize](https://github.com/xuri/excelize)
- Clipboard support by [clipboard](https://github.com/atotto/clipboard)
- Themes inspired by [Catppuccin](https://github.com/catppuccin/catppuccin), [Nord](https://www.nordtheme.com/), [Rosé Pine](https://rosepinetheme.com/), [Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme), [Gruvbox](https://github.com/morhetz/gruvbox), and [Dracula](https://draculatheme.com/)
