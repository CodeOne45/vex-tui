# Changelog

All notable changes to Vex will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[2.0.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v2.0.0
[1.1.1]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.1.1
[1.1.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.1.0
[1.0.0]: https://github.com/CodeOne45/vex-tui/releases/tag/v1.0.0
[0.0.1]: https://github.com/CodeOne45/vex-tui/releases/tag/v0.0.1
