# Vex - Excel Viewer

A TUI Excel and CSV viewer in Go 

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.1.0-brightgreen.svg)](https://github.com/odesaur/vex-tui/releases)

![Vex Demo](assets/vex-demo.gif)


## Installation

### Using go install (quickest)

```bash
go install github.com/odesaur/vex-tui@latest
```

This places the `vex-tui` binary in `$(go env GOPATH)/bin` (or `$(go env GOBIN)` if set). Ensure that directory is on your `PATH`.

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

`make install` will produce an optimized build and install it to your Go bin directory so you can run `vex-tui` from anywhere.

## Usage

```bash
# Basic usage
vex-tui data.xlsx

# With a specific theme
vex-tui report.csv --theme nord

# Short flag
vex-tui sales.xlsx -t tokyo-night
```

## Keyboard Shortcuts

### Navigation

- `↑↓←→` or `hjkl` - Navigate cells
- `Page Up/Down` - Scroll by page
- `Ctrl+U/D` - Alternative page scroll
- `Home/End` or `0/$` - First/last column
- `g/G` - First/last column
- `Tab/Shift+Tab` - Next/previous sheet

### Search & Actions

- `/` - Search (vim-style)
- `n/N` - Next/previous result
- `Ctrl+G` - Jump to cell
- `Enter` - View cell details
- `c` - Copy cell
- `C` - Copy entire row
- `f` - Toggle formula display
- `e` - Export sheet
- `t` - Theme selector
- `?` - Toggle help
- `q` or `Ctrl+C` - Quit

### Data Visualization

Step 1: Select Data Range

1. Navigate to your data
2. Press 'V' (shift+v) to start selection
3. Move cursor to select range (arrows/hjkl)
4. Press 'V' again to finish selection

Step 2: Visualize

1. Press 'v' (lowercase) to open visualization
2. Press 1-4 to switch between chart types:
   - 1: Bar Chart
   - 2: Line Chart
   - 3: Sparkline
   - 4: Pie Chart
3. Press Esc to close

## Project Structure

```
vex-tui/
├── main.go                 # Application entry point
├── internal/
│   ├── app/               # Application logic
│   │   ├── model.go       # State management
│   │   ├── update.go      # Event handling
│   │   ├── view.go        # Rendering logic
│   │   └── keys.go        # Keybindings
│   ├── loader/            # File I/O operations
│   │   └── loader.go
│   ├── theme/             # Theme management
│   │   └── theme.go
│   └── ui/                # UI utilities
│       └── ui.go
└── pkg/
    └── models/            # Data models
        └── models.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/odesaur/vex-tui.git
cd vex-tui

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o vex-tui .
```


## Acknowledgments

- Built with the amazing [Charm](https://charm.sh/) ecosystem
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
  - [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
  - [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- Excel parsing by [Excelize](https://github.com/xuri/excelize)
- Clipboard support by [clipboard](https://github.com/atotto/clipboard)
- Themes inspired by [Catppuccin](https://github.com/catppuccin/catppuccin), [Nord](https://www.nordtheme.com/), [Rosé Pine](https://rosepinetheme.com/), [Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme), [Gruvbox](https://github.com/morhetz/gruvbox), and [Dracula](https://draculatheme.com/)
