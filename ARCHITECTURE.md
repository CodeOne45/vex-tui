# Architecture Documentation

This document describes the architecture and design decisions of Excel TUI.

## Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [Architecture Layers](#architecture-layers)
- [Data Flow](#data-flow)
- [Component Details](#component-details)
- [Design Patterns](#design-patterns)
- [Performance Considerations](#performance-considerations)
- [Extension Points](#extension-points)

## Overview

Excel TUI is a terminal-based spreadsheet viewer built using the Elm Architecture (via Bubble Tea). It follows a Model-View-Update pattern with clear separation between data, logic, and presentation.

### Tech Stack

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI framework)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) (terminal styling)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (UI components)
- **Excel Parsing**: [Excelize](https://github.com/xuri/excelize) (Excel file handling)
- **Clipboard**: [clipboard](https://github.com/atotto/clipboard) (cross-platform clipboard)

## Design Principles

### 1. Separation of Concerns

Each package has a single, well-defined responsibility:

- `app` - Application logic and state management
- `loader` - File I/O operations
- `theme` - Visual theme management
- `ui` - Reusable UI utilities
- `models` - Data structures

### 2. Dependency Direction

```
main.go
  ↓
internal/app (depends on ↓)
  ↓
internal/loader, internal/theme, internal/ui
  ↓
pkg/models (no dependencies)
```

### 3. Interface Segregation

Small, focused interfaces rather than large monolithic ones.

### 4. Immutability Where Possible

State changes are explicit and controlled through the Update function.

## Architecture Layers

### Layer 1: Entry Point (`main.go`)

**Responsibilities:**

- Parse command-line arguments
- Validate input
- Initialize application
- Start Bubble Tea program

**Dependencies:** `app`, `loader`

### Layer 2: Application (`internal/app`)

**Responsibilities:**

- State management (Model)
- Event handling (Update)
- View rendering (View)
- Keyboard bindings

**Key Files:**

- `model.go` - Application state
- `update.go` - Event handlers
- `view.go` - Rendering logic
- `keys.go` - Keybindings

### Layer 3: Services

#### Loader (`internal/loader`)

**Responsibilities:**

- File parsing (Excel, CSV)
- Data export (CSV, JSON)
- Search operations

**Key Functions:**

- `LoadFile(filename) ([]Sheet, error)`
- `ExportToCSV(sheet, filename) error`
- `ExportToJSON(sheet, filename) error`
- `SearchSheet(sheet, term) []Cell`

#### Theme (`internal/theme`)

**Responsibilities:**

- Theme definitions
- Theme switching
- Color management

**Key Functions:**

- `GetThemeNames() []string`
- `SetTheme(name) bool`
- `GetCurrentTheme() Theme`

#### UI (`internal/ui`)

**Responsibilities:**

- Style initialization
- Helper functions
- Rendering utilities

**Key Functions:**

- `InitStyles() *Styles`
- `ColIndexToLetter(index) string`
- `TruncateToWidth(s, width) string`
- `RenderModal(width, height, modal) string`

### Layer 4: Models (`pkg/models`)

**Responsibilities:**

- Data structures
- Constants
- Types

**Key Types:**

- `Cell` - Single spreadsheet cell
- `Sheet` - Worksheet with data
- `Mode` - Application mode enum
- `StatusMsg` - Status message

## Data Flow

### Initialization Flow

```
1. main.go
   ↓
2. Parse arguments & validate file
   ↓
3. loader.LoadFile() → []Sheet
   ↓
4. app.NewModel() → Model
   ↓
5. tea.NewProgram() → Start
```

### Event Processing Flow

```
1. User Input (keyboard)
   ↓
2. tea.Msg delivered to Update()
   ↓
3. Model state changes
   ↓
4. View() called with new state
   ↓
5. Rendered output to terminal
```

### Mode-Based Flow

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.mode {
        case ModeNormal:   return m.updateNormal(msg)
        case ModeSearch:   return m.updateSearch(msg)
        case ModeDetail:   return m.updateDetail(msg)
        case ModeJump:     return m.updateJump(msg)
        case ModeExport:   return m.updateExport(msg)
        case ModeTheme:    return m.updateTheme(msg)
        }
    }
    return m, nil
}
```

## Component Details

### Application State (Model)

```go
type Model struct {
    // Data
    sheets        []models.Sheet
    currentSheet  int

    // Cursor and viewport
    cursorRow     int
    cursorCol     int
    offsetRow     int
    offsetCol     int

    // UI state
    width         int
    height        int
    mode          models.Mode

    // Search state
    searchQuery   string
    searchResults []models.Cell
    searchIndex   int

    // UI components
    searchInput   textinput.Model
    jumpInput     textinput.Model
    exportInput   textinput.Model

    // Settings
    showFormulas  bool
    themeName     string

    // Metadata
    status        models.StatusMsg
    help          help.Model
    keys          KeyMap
    filename      string
    styles        *ui.Styles
}
```

### Viewport Management

The viewport system ensures the cursor is always visible:

```go
func (m *Model) adjustViewport() {
    visibleRows := max(1, m.height-9)
    visibleCols := max(1, (m.width-8)/(MinCellWidth+2))

    // Vertical adjustment
    if m.cursorRow < m.offsetRow {
        m.offsetRow = m.cursorRow
    } else if m.cursorRow >= m.offsetRow+visibleRows {
        m.offsetRow = m.cursorRow - visibleRows + 1
    }

    // Horizontal adjustment
    // ...
}
```

### Theme System

Themes are defined as color collections:

```go
type Theme struct {
    Name          string
    Primary       lipgloss.Color
    Secondary     lipgloss.Color
    Accent        lipgloss.Color
    // ... more colors
}
```

Styles are generated from themes:

```go
func InitStyles() *Styles {
    t := theme.GetCurrentTheme()
    return &Styles{
        Title: lipgloss.NewStyle().
            Foreground(t.Primary).
            Bold(true),
        // ... more styles
    }
}
```

## Design Patterns

### 1. Model-View-Update (Elm Architecture)

**Model**: Application state
**Update**: Pure function: `(Model, Msg) → (Model, Cmd)`
**View**: Pure function: `Model → String`

### 2. Strategy Pattern (Themes)

Different themes are strategies for coloring the UI. Themes can be swapped at runtime without changing application logic.

### 3. Factory Pattern (Model Creation)

```go
func NewModel(filename, sheets, themeName) Model {
    // Complex initialization
    // Returns fully configured model
}
```

### 4. Command Pattern (Bubble Tea Commands)

Actions that produce side effects are represented as commands:

```go
return m, textinput.Blink  // Command to blink cursor
```

### 5. Observer Pattern (Event System)

Bubble Tea's message passing is an implementation of the observer pattern.

## Performance Considerations

### 1. Lazy Rendering

Only visible cells are rendered:

```go
visibleRows := max(1, m.height-9)
endRow := min(m.offsetRow+visibleRows, sheet.MaxRows)
for row := m.offsetRow; row < endRow {
    // Render only visible rows
}
```

### 2. Efficient String Building

Use `strings.Builder` for concatenation:

```go
var b strings.Builder
b.WriteString(header)
b.WriteString(content)
return b.String()
```

### 3. Pre-allocated Slices

```go
sheets := make([]models.Sheet, 0, len(sheetList))
cellRow := make([]models.Cell, 0, len(row))
```

### 4. CSV Reader Optimization

```go
reader := csv.NewReader(file)
reader.ReuseRecord = true  // Reuse memory
```

### 5. Viewport Calculations

Cache viewport dimensions to avoid recalculation:

```go
visibleRows := max(1, m.height-9)
visibleCols := max(1, (m.width-8)/(MinCellWidth+2))
```

## Extension Points

### 1. Adding New Themes

1. Define theme in `internal/theme/theme.go`:

```go
themes["mytheme"] = Theme{
    Name: "My Theme",
    // ... colors
}
```

2. Add to `GetThemeNames()`

### 2. Adding New Modes

1. Add mode constant to `pkg/models/models.go`:

```go
const (
    // ...
    ModeCustom Mode = iota + 6
)
```

2. Add update handler in `internal/app/update.go`:

```go
case models.ModeCustom:
    return m.updateCustom(msg)
```

3. Add view renderer in `internal/app/view.go`:

```go
case models.ModeCustom:
    return m.renderCustom()
```

### 3. Adding New File Formats

1. Add loader in `internal/loader/loader.go`:

```go
func loadXML(filename string) ([]models.Sheet, error) {
    // Implementation
}
```

2. Update `LoadFile()` switch statement:

```go
case ".xml":
    return loadXML(filename)
```

### 4. Adding New Export Formats

Similar to file formats, add to `loader.go`:

```go
func ExportToXML(sheet models.Sheet, filename string) error {
    // Implementation
}
```

## Testing Strategy

### Unit Tests

Test individual functions in isolation:

```go
func TestColIndexToLetter(t *testing.T) {
    // Test helper functions
}

func TestLoadCSV(t *testing.T) {
    // Test file loading
}
```

### Integration Tests

Test component interaction:

```go
func TestSearchAndJump(t *testing.T) {
    // Test search → jump workflow
}
```

### End-to-End Tests

Test complete user workflows:

```go
func TestOpenFileAndExport(t *testing.T) {
    // Test file → view → export
}
```

## Error Handling Strategy

### Levels of Error Handling

1. **Fatal Errors**: Exit program with error message

   - File not found
   - Invalid file format
   - Permission denied

2. **Recoverable Errors**: Show status message

   - Search no results
   - Invalid jump reference
   - Export failed

3. **Warnings**: Log but continue
   - Sheet read error (skip sheet)
   - File close error

### Error Propagation

Use `fmt.Errorf` with `%w` for error wrapping:

```go
if err != nil {
    return nil, fmt.Errorf("failed to load Excel file: %w", err)
}
```

## Security Considerations

### Input Validation

- File path validation
- Cell reference validation
- Search term sanitization

### File Operations

- Read-only access by default
- Proper file closing
- No arbitrary code execution

### Formula Handling

- Formulas displayed but never executed
- No eval() or similar operations

## Scalability

### Memory Usage

- Lazy loading of cell data
- Viewport-based rendering
- Efficient data structures

### Large Files

- Stream processing for CSV
- Chunked loading for Excel
- Progressive rendering

## Future Architecture Improvements

### Planned

1. **Plugin System**: Load custom themes/exporters
2. **Configuration File**: User preferences persistence
3. **Macro System**: Record/replay actions
4. **Async Operations**: Background file loading

### Considered

1. **Database Backend**: For very large files
2. **Remote Files**: Load from URLs
3. **Collaboration**: Real-time multi-user
4. **Web Interface**: Browser-based UI

## Conclusion

The architecture is designed for:

- **Simplicity**: Easy to understand and modify
- **Maintainability**: Clear structure and documentation
- **Extensibility**: Easy to add features
- **Performance**: Efficient for large files
- **Testability**: Modular and mockable

For questions or suggestions, please open an issue or pull request.
