# 📊 Vex - Terminal Spreadsheet Editor

A beautiful, fast, and feature-rich terminal-based Excel and CSV editor with vim-style keybindings and formula support.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.0.0-brightgreen.svg)](https://github.com/CodeOne45/vex-tui/releases)

![Vex Demo](assets/vex-demo.gif)

## ✨ Features

### 🎨 Ten Beautiful Themes

**Dark Themes:**
- **Catppuccin Mocha** - Soft pastels, perfect for all-day use
- **Nord** - Cool Arctic blues, minimal and focused
- **Rosé Pine** - Elegant rose tones, sophisticated
- **Tokyo Night** - Vibrant cyberpunk aesthetic
- **Gruvbox** - Warm retro colors, comfortable
- **Dracula** - Classic high contrast theme

**Light Themes:**
- **Catppuccin Latte** - Gentle pastel light theme
- **Solarized Light** - Balanced contrast
- **GitHub Light** - Clean and minimal
- **One Light** - Soft Atom-inspired colors

### ✏️ Full Editing Capabilities

- **Edit cells** with formulas and values
- **Insert/delete** rows and columns
- **Copy/paste** single cells or ranges
- **Fill down/right** for quick data entry
- **Apply formulas** to entire ranges with automatic reference adjustment
- **Auto-save** with modification tracking
- **Undo-friendly** workflow with clear status messages

### 📐 Powerful Formula Engine

**Arithmetic:** `=A1+B1`, `=C1*2`, `=D1/E1-F1`

**15+ Built-in Functions:**
- `SUM(A1:A10)` - Sum range
- `AVERAGE(B1:B20)` / `AVG(...)` - Average values
- `COUNT(C1:C50)` - Count numbers
- `MAX(D1:D100)` / `MIN(...)` - Find max/min
- `IF(A1>100, "High", "Low")` - Conditional logic
- `CONCATENATE(A1, " ", B1)` / `CONCAT(...)` - Join text
- `UPPER(A1)` / `LOWER(A1)` - Change case
- `LEN(A1)` - Text length
- `ROUND(A1, 2)` - Round numbers
- `ABS(A1)` - Absolute value
- `SQRT(A1)` - Square root
- `POWER(2, 8)` / `POW(...)` - Exponentiation

**Auto-recalculation** when cells change

### 🔍 Powerful Navigation

- **Vim-style keybindings** (hjkl) and arrow keys
- **Jump to cell** (Ctrl+G) - supports `A100`, `500`, or `10,5` formats
- **Search** (/) across cells and formulas
- **Navigate results** (n/N)
- **Page Up/Down**, Home/End
- **Multi-sheet** support with Tab navigation

### 📋 Data Operations

- **Copy** cell (c) or entire row (C)
- **Paste** (p) with multi-cell support
- **Export** to CSV or JSON
- **Save** (Ctrl+S) with format preservation
- **Save As** (Ctrl+Shift+S) to new file
- **Toggle formula display** (f)
- **View cell details** (Enter)

### 📊 Live Data Visualization

- **Bar charts** - Compare values visually
- **Line charts** - Show trends over time
- **Sparklines** - Compact inline charts
- **Pie charts** - Display proportions

**How to use:**
1. Press `V` to start range selection
2. Move cursor and press `V` again
3. Press `v` to open visualization
4. Press 1-4 to switch chart types

### 📑 File Support

- **Excel files** (.xlsx, .xlsm, .xls) with formula preservation
- **CSV files** with formula support (saved as text)
- **Multiple sheets** with easy navigation
- **Large file optimization** with lazy loading
- **Safe saving** with backup on errors

## 🚀 Installation

### Using Homebrew (macOS/Linux)

```bash
brew install CodeOne45/tap/vex
```

### Using go install

```bash
go install github.com/CodeOne45/vex-tui@latest
```

### Download Binary

Download pre-built binaries from the [releases page](https://github.com/CodeOne45/vex-tui/releases).

**Available for:**
- macOS (Intel & Apple Silicon)
- Linux (x64 & ARM64)
- Windows (x64)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/CodeOne45/vex-tui.git
cd vex-tui

# Install dependencies
go mod download

# Build
go build -o vex .

# Optional: Install globally
sudo mv vex /usr/local/bin/
```

## 📖 Usage

```bash
# View a file (read-only until you press 'i')
vex data.xlsx

# Start with a specific theme
vex report.csv --theme nord

# Create new file (will be created on first save)
vex newfile.xlsx
```

## ⌨️ Keyboard Shortcuts

### Navigation

- `↑↓←→` or `hjkl` - Navigate cells
- `Page Up/Down` or `Ctrl+U/D` - Scroll by page
- `Home/End` or `0/$` - First/last column
- `g/G` - First/last column
- `Tab/Shift+Tab` - Next/previous sheet
- `Ctrl+G` - Jump to specific cell

### Editing

- `i` - Enter edit mode
- `Enter` - Commit changes (in edit) / View details (in normal)
- `Tab` - Save and move right (in edit mode)
- `Shift+Tab` - Save and move left (in edit mode)
- `Esc` - Cancel editing
- `x` - Delete cell content
- `dd` - Delete current row
- `dc` - Delete current column

### Cell Operations

- `c` - Copy cell
- `C` - Copy entire row
- `p` - Paste
- `o` - Insert row below
- `O` - Insert column right
- `Ctrl+J` - Fill down (requires selection)
- `Ctrl+L` - Fill right (requires selection)
- `Ctrl+A` - Apply formula to range (requires selection)

### File Operations

- `Ctrl+S` - Save file
- `Ctrl+Shift+S` - Save as
- `e` - Export to CSV/JSON
- `q` - Quit (press twice if unsaved changes)

### Search & Navigation

- `/` - Search
- `n/N` - Next/previous search result
- `Esc` - Clear search

### Visualization & Display

- `V` - Start/finish range selection
- `v` - Open visualization (after selection)
- `1-4` - Switch chart types (in viz mode)
- `f` - Toggle formula display
- `t` - Change theme
- `?` - Toggle help

## 💡 Quick Start Guide

### Viewing a File

```bash
vex mydata.xlsx
# Navigate with arrow keys or hjkl
# Press 'f' to toggle formula view
# Press 't' to change theme
```

### Editing Data

```bash
# 1. Open file
vex mydata.xlsx

# 2. Navigate to a cell
# 3. Press 'i' to edit
# 4. Type value or formula: =SUM(A1:A10)
# 5. Press Enter to save
# 6. Press Ctrl+S to save file
```

### Working with Formulas

```bash
# Create a formula
Press 'i' on cell C1
Type: =A1+B1
Press Enter

# Apply to entire column
Press 'V' to start selection at C1
Move to C10 and press 'V'
Move back to C1
Press Ctrl+A to apply formula
# Result: C1=A1+B1, C2=A2+B2, C3=A3+B3, etc.
```

### Bulk Data Entry

```bash
Press 'i' to start editing
Type value and press Tab
Continue typing in next cell
Press Tab to keep moving right
Press Ctrl+S when done
```

## 🏗️ Project Structure

```
vex-tui/
├── main.go                    # Application entry point
├── internal/
│   ├── app/
│   │   ├── model.go          # State management
│   │   ├── update.go         # Event handling
│   │   ├── view.go           # Rendering logic
│   │   ├── keys.go           # Keybindings
│   │   ├── edit.go           # Edit operations
│   │   └── formulas.go       # Formula evaluation engine
│   ├── loader/
│   │   ├── loader.go         # File loading
│   │   └── save.go           # File saving
│   ├── theme/
│   │   └── theme.go          # Theme definitions
│   └── ui/
│       └── ui.go             # UI utilities
└── pkg/
    └── models/
        └── models.go         # Data models
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
git clone https://github.com/CodeOne45/vex-tui.git
cd vex-tui

# Install dependencies
make deps

# Run tests
make test

# Build
make build

# Run with sample data
make run
```

### Code Style

- Run `make fmt` before committing
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Add tests for new features
- Update documentation

## 📊 Examples

### Sales Report with Formulas

```
     A        B        C          D
1  Item     Price    Qty      Total
2  Apples   1.50     10      =B2*C2
3  Oranges  2.00     5       =B3*C3
4  Total    -        -       =SUM(D2:D3)
```

Select D2:D4, cursor on D2, press Ctrl+A to apply formula pattern.

### Conditional Formatting (Formula-based)

```
     A           B
1  Score    Grade
2  85       =IF(A2>=90,"A",IF(A2>=80,"B","C"))
```

Apply to B2:B20 with Ctrl+A for instant grading.

### Data Summary

```
     A              B
1  Sales         1000
2  Costs          600
3  Profit      =A1-A2
4  Margin      =ROUND(A3/A1*100,2)
```

## 🔧 Advanced Features

### Formula Auto-adjustment

When you apply a formula to a range, cell references automatically adjust:

```
Source: A1 contains =B1+C1

Apply to A1:A5:
  A1: =B1+C1
  A2: =B2+C2
  A3: =B3+C3
  A4: =B4+C4
  A5: =B5+C5
```

### Multi-cell Paste

Copy ranges from other apps and paste into Vex:
- Tab-separated values paste as multiple columns
- Newline-separated values paste as multiple rows
- Formulas (starting with =) are preserved

### Keyboard-driven Workflow

```
i          Enter edit mode
Type data  
Tab        Save and move to next cell
Tab        Continue entering data
Tab        Keep going...
Ctrl+S     Save entire file
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Excel parsing by [Excelize](https://github.com/xuri/excelize)
- Clipboard by [atotto/clipboard](https://github.com/atotto/clipboard)
- Themes inspired by [Catppuccin](https://github.com/catppuccin/catppuccin), [Nord](https://www.nordtheme.com/), [Rosé Pine](https://rosepinetheme.com/), and others

## 🔒 Security

If you discover a security vulnerability, please create a private security advisory on GitHub or email the maintainers.

## 📮 Contact

- **GitHub**: [@CodeOne45](https://github.com/CodeOne45)
- **Issues**: [GitHub Issues](https://github.com/CodeOne45/vex-tui/issues)
- **Discussions**: [GitHub Discussions](https://github.com/CodeOne45/vex-tui/discussions)

## ⭐ Star History

If you find Vex useful, please consider giving it a star on GitHub!

---

Made with ❤️ for terminal enthusiasts and spreadsheet power users everywhere.
