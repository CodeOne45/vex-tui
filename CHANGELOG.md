# Changelog

All notable changes to Excel TUI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-02-01

### Added

- Live ASCII charts (Bar, Line, Sparkline, Pie)
- 'v' visualization window
- Auto-scaling and grid-based rendering

## [1.0.0] - 2025-01-26

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
- Jump to cell feature (Ctrl+G) supporting multiple formats:
  - Column letter + row number (e.g., A100)
  - Row number only (e.g., 500)
  - Row, column coordinates (e.g., 10,5)
- Viewport auto-centering when jumping
- Smart viewport scrolling

#### Cell Operations

- Cell detail modal (Enter key) showing:
  - Cell reference
  - Full value with text wrapping
  - Formula if present
  - Cell type detection
- Copy entire row feature (Shift+C)
- Enhanced copy cell with preview
- Formula display toggle

#### UI Improvements

- Formula bar showing current cell info
- Enhanced status bar with position, mode, and search results
- Compact help display
- Beautiful centered modals for dialogs
- Real-time status messages for all operations

#### Navigation Enhancements

- `g` key - Jump to first column
- `G` key - Jump to last column
- `Ctrl+U` - Alternative Page Up
- `Ctrl+D` - Alternative Page Down
- `0` - Alternative Home
- `$` - Alternative End

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

## [0.0.1] - 2025-01-26

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

[1.0.0]: https://github.com/odesaur/vex-tui/releases/tag/v1.0.0
