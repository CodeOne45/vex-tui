# Changelog

All notable changes to Vex will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2026-04-17

### Added

#### Row Filter (`f`)
- Press `f` to filter rows on the current column or all columns (Tab toggles scope)
- Supports plain text (contains), `>`, `<`, `>=`, `<=`, `=` expressions for numeric and exact-match comparisons
- Active filter shown in the title bar as `▾ <expr>`; **Esc** clears it without touching the underlying data

#### Data Profile (`W`)
- New modal showing a per-column breakdown: detected type (number / text / mixed), non-null count, null count
- Numeric columns show Σ / avg / min / max inline; text columns show unique count and top sample values
- Highlights the column currently under the cursor

#### Sort Columns (`s` / `S` / `u`)
- `s` sorts the current column ascending, `S` descending
- `u` restores the original row order (pre-sort state is saved in memory)
- Sort indicator shown in the column header and title bar; freeze-header aware

#### Freeze Header Row (`Ctrl+R`)
- Toggles the first row pinned at all times while scrolling
- A divider line separates the frozen row from scrollable data
- `⌶` indicator in title bar; sort respects the frozen header

#### Live Column Statistics
- When the cursor is on any numeric column, the status bar automatically shows Σ / avg / min / max / n
- No selection needed; updates as you navigate

#### Chart PNG Export (`p` in chart view)
- Press `p` inside the chart modal to export the current chart as a beautiful PNG via [`freeze`](https://github.com/charmbracelet/freeze)
- Output file is `<filename>-chart.png` saved next to the source file
- If `freeze` is not on PATH the binary is located across all common Homebrew and Linux install paths automatically
- If still not found, shows the exact install command as a status bar hint

#### CSV Delimiter Auto-Detection + CLI Flag
- Vex reads the first line of every CSV and picks the best delimiter from `,` `;` `\t` `|`
- New `-d` / `--delimiter` flag for explicit override: `vex data.csv -d ';'`, `vex data.tsv -d '\t'`

### Changed

#### UI / UX
- **Cell type coloring**: numbers render in accent color, formulas in secondary color, text in default — no configuration needed
- **Status bar redesign**: left cluster (cell ref + dimensions + mode flags), middle cluster (column stats or selection size), right cluster (search + status message)
- **Title bar**: `●` modified indicator appears immediately on first edit; freeze / sort / filter state shown as compact glyphs
- **Selection highlight**: uses primary color background, much more visible than before
- **Responsive modals**: theme picker and chart modal now scale to terminal width/height instead of hardcoded sizes
- **Chart bar chart**: supports negative values (shown in red), `░` unfilled track, cleaner value labels
- **Chart line chart**: Bresenham algorithm draws actual connecting lines (`·`) between data points (`●`)
- **Chart data extraction**: auto-detects label vs value columns — works with or without a leading text column
- **Chart type tabs**: rendered as proper tab buttons instead of a plain list

#### Keybinding Changes
- `f` — now opens row filter (was: toggle formula display)
- `ctrl+f` — toggle formula display (moved from `f`)
- `W` — data profile
- `s` / `S` / `u` — sort asc / desc / unsort
- `ctrl+r` — freeze header

#### Release Workflow
- GoReleaser config: grouped changelog, `test` stanza in Homebrew formula, Scoop bucket support
- `Makefile`: `make tag VER=v2.1.0` tags and pushes, triggering the full CI/release pipeline; `make snapshot` for local multi-platform builds

### Fixed

- **Esc cancels selection**: pressing Esc in normal mode while a range is selected now cancels it (previously Esc only cleared search)
- **Freeze binary lookup**: `freeze` is located across PATH, all Homebrew symlink paths, and the Cellar directly — fixing "not found" errors when the Homebrew symlink is broken

## [2.0.2] - 2025-02-13

### Fixed

- **Excel file crash**: Updated excelize library from v2.8.0 to v2.10.0, fixing `panic: index out of range` when opening certain Excel files
- **Delete row/column shortcuts**: `dd` (delete row) and `dc` (delete column) now work correctly as two-key sequences instead of `d` immediately deleting a row
- **Help shortcut**: Pressing `?` now properly toggles the full help view showing all keybindings
- **CLI `-h` flag**: Added `-h` as shorthand for `--help` so both `-h` and `--help` work
- **Multiline cell display**: Cells containing newlines now show `↵` indicators in grid view instead of collapsing to spaces
- **Multiline cell editing**: Edit mode now uses a multi-line text area, preserving newlines in cell values (Enter adds newline, Ctrl+S confirms edit)
- **CI pipeline**: Fixed Go version compatibility and coverage tool errors

### Added

- **Copy selected range**: Press `V` to select a range, then `c` to copy it (tab-separated columns, newline-separated rows)
- **Vim-style file navigation**: `gg` jumps to top of file, `G` jumps to bottom of file (previously these only moved between columns)

### Changed

- Minimum Go version bumped to 1.24 (required by excelize v2.10.0)
- Release workflow now uses GoReleaser for automated GitHub releases and Homebrew formula updates
- Edit mode keybindings: Enter inserts newline, Ctrl+S saves edit, Esc cancels

## [2.0.1] - 2024-12-14

### Fixed

- Configure Homebrew tap for macOS installation
- Version description update in `printVersion` function

## [2.0.0] - 2024-12-14

### 🎉 Major Release - Full Editing Capabilities

This is a major release that transforms Vex from a viewer into a full-featured terminal spreadsheet editor with formula support.

### ✨ Added

#### Edit Mode & Cell Operations
- **Full cell editing** with `i` key to enter edit mode
- **Tab/Shift+Tab navigation** in edit mode to move between cells while editing
- **Formula support** with 15+ built-in functions
- **Delete operations**: Delete cell (`x`), row (`dd`), column (`dc`)
- **Insert operations**: Insert row (`o`), column (`O`)
- **Copy/Paste**: Enhanced clipboard operations with multi-cell paste support
- **Fill operations**: Fill down (`Ctrl+J`), fill right (`Ctrl+L`)
- **Apply formula to range** (`Ctrl+A`) - Apply current cell's formula to selected range with automatic reference adjustment

#### Formula Engine
- Arithmetic operations: `+`, `-`, `*`, `/`
- **SUM(range)** - Sum of values in range
- **AVERAGE(range) / AVG(range)** - Average of values
- **COUNT(range)** - Count numeric values
- **MAX(range)** - Maximum value in range
- **MIN(range)** - Minimum value in range
- **IF(condition, true_val, false_val)** - Conditional logic
- **CONCATENATE(...) / CONCAT(...)** - Combine text
- **UPPER(text) / LOWER(text)** - Change case
- **LEN(text)** - Text length
- **ROUND(number, digits)** - Round number
- **ABS(number)** - Absolute value
- **SQRT(number)** - Square root
- **POWER(base, exp) / POW(base, exp)** - Exponentiation
- **Automatic formula recalculation** when cells change
- **Relative cell references** that adjust when formulas are copied/applied

#### File Operations
- **Save functionality** (`Ctrl+S`) - Save changes to Excel or CSV
- **Save As** (`Ctrl+Shift+S`) - Save to new filename
- **Format preservation** - Excel files maintain formulas, CSV saves formulas as text
- **Modified file tracking** - Visual indicator when file has unsaved changes
- **Quit confirmation** - Double-tap `q` to quit without saving when modified

#### Enhanced UX
- **Edit mode indicator** - Clear visual feedback when in edit mode
- **Formula bar** - Shows current cell formula or value
- **Status messages** - Informative feedback for all operations
- **Modified indicator** - `[Modified]` tag when file has unsaved changes
- **Smart formula application** - Formulas automatically adjust cell references when applied to ranges

### 🔧 Changed

- **Enter key behavior** - Now opens cell detail view in normal mode, commits edit in edit mode
- **Enhanced selection** - Selection now shows formula application option
- **Improved status bar** - Shows edit mode, modification status, and current operation
- **Better error handling** - Clear error messages for invalid operations

### 🐛 Fixed

- **Quit confirmation** - Properly handles unsaved changes with double-quit pattern
- **Cell boundary checks** - No more crashes when editing at sheet boundaries
- **Formula evaluation** - Improved error handling for invalid formulas
- **Memory management** - Better handling of large sheets during edit operations

### 📚 Documentation

- Added comprehensive editing guide
- Formula reference documentation
- Keyboard shortcuts updated with all new bindings
- Usage examples for common editing scenarios

### ⌨️ New Keyboard Shortcuts

**Edit Mode:**
- `i` - Enter edit mode
- `Enter` - Commit changes (in edit mode)
- `Tab` - Save and move to next cell
- `Shift+Tab` - Save and move to previous cell
- `Esc` - Cancel editing

**Cell Operations:**
- `x` - Delete cell content
- `dd` - Delete row
- `dc` - Delete column
- `o` - Insert row below
- `O` - Insert column right
- `p` - Paste from clipboard

**File Operations:**
- `Ctrl+S` - Save file
- `Ctrl+Shift+S` - Save as
- `q` (twice) - Quit without saving (when modified)

**Formula Operations:**
- `Ctrl+A` - Apply formula to selected range
- `Ctrl+J` - Fill down
- `Ctrl+L` - Fill right
- `f` - Toggle formula display

### 💡 Usage Examples

**Basic Editing:**
```
1. Press 'i' on any cell to start editing
2. Type your value or formula (e.g., =A1+B1)
3. Press Enter to save
4. Press Ctrl+S to save the file
```

**Formula Application:**
```
1. Create a formula in cell A1: =B1+C1
2. Press 'V' to start selection at A1
3. Move to A10 and press 'V' again
4. Move cursor to A1
5. Press Ctrl+A to apply formula to entire range
   (A1=B1+C1, A2=B2+C2, A3=B3+C3, etc.)
```

**Quick Data Entry:**
```
1. Press 'i' to edit cell
2. Type value and press Tab
3. Continue typing in next cell
4. Repeat for fast data entry
```

## [1.1.1] - 2024-12-04

### Fixed

- Fixed GitHub CI "no such tool covdata" error by removing race detector from coverage tests
- Fixed CI workflow to properly run across different Go versions and platforms

### Added

- Added `--version` / `-v` flag to display version information
- Added `--help` / `-h` flag to display usage information
- Improved CLI argument parsing with proper flag handling
- Added comprehensive help text with examples and keyboard shortcuts

### Changed

- Improved error messages and user feedback
- Enhanced CI workflow with separate lint job
- Updated release workflow for better reliability

## [1.1.0] - 2024-11-27

### Added

- Live ASCII charts (Bar, Line, Sparkline, Pie)
- 'v' visualization window
- Auto-scaling and grid-based rendering

## [1.0.0] - 2024-11-26

### Added

#### Themes & Visuals

- Six professional themes (Catppuccin, Nord, Rosé Pine, Tokyo Night, Gruvbox, Dracula)
- Theme switcher accessible with `t` key
- CLI flag `--theme` for setting theme on launch
- Visual highlighting for rows, columns, and search matches
- Color-coded status messages (info, success, warning, error)
- Dynamic style system that updates on theme change

#### Search & Navigation

- Vim-style search bar at bottom of screen
- Persistent search display with active query
- Search highlighting with yellow background
- Jump to cell feature (Ctrl+G) supporting multiple formats
- Viewport auto-centering when jumping
- Smart viewport scrolling

#### Cell Operations

- Cell detail modal (Enter key)
- Copy entire row feature (Shift+C)
- Enhanced copy cell with preview
- Formula display toggle

#### UI Improvements

- Formula bar showing current cell info
- Enhanced status bar with position, mode, and search results
- Compact help display
- Beautiful centered modals for dialogs
- Real-time status messages for all operations

### Changed

- Complete code restructure following Go best practices
- Modular architecture with clean separation of concerns
- Improved error handling throughout
- Better state management
- Enhanced performance for large files

### Fixed

- Panic on startup with uninitialized terminal size
- Negative viewport calculations
- Empty cell handling
- Memory leaks with file operations

### Security

- Added input validation and sanitization
- Safe file handling with proper cleanup
- No code execution from formulas
- Read-only file access by default

## [0.0.1] - 2024-01-26

### Added

- Initial release
- Multi-format support (.xlsx, .xlsm, .xls, .csv)
- Basic TUI with Bubble Tea framework
- Search functionality
- Formula display
- Clipboard support
- Export to CSV/JSON
- Vim-style navigation
- Multiple sheet support

[2.1.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v2.1.0
[2.0.2]: https://github.com/CodeOne45/vex-tui/releases/tag/v2.0.2
[2.0.1]: https://github.com/CodeOne45/vex-tui/releases/tag/v2.0.1
[2.0.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v2.0.0
[1.1.1]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.1.1
[1.1.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.1.0
[1.0.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.0.0
[0.0.1]: https://github.com/CodeOne45/vex-tui/releases/tag/v0.0.1
