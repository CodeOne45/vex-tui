package app

import (
	"github.com/CodeOne45/vex-tui/internal/theme"
	"github.com/CodeOne45/vex-tui/internal/ui"
	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
)

// Model represents the application state
type Model struct {
	sheets        []models.Sheet
	currentSheet  int
	cursorRow     int
	cursorCol     int
	offsetRow     int
	offsetCol     int
	width         int
	height        int
	mode          models.Mode
	searchInput   textinput.Model
	jumpInput     textinput.Model
	exportInput   textinput.Model
	editInput     textinput.Model
	saveAsInput   textinput.Model
	searchQuery   string
	searchResults []models.Cell
	searchIndex   int
	showFormulas  bool
	status        models.StatusMsg
	help          help.Model
	keys          KeyMap
	filename      string
	themeName     string
	styles        *ui.Styles

	// Chart visualization
	chartType   int
	selectStart [2]int
	selectEnd   [2]int
	isSelecting bool

	// Edit mode
	isEditing   bool
	modified    bool
	fileFormat  string
	quitConfirm bool
}

// NewModel creates a new application model
func NewModel(filename string, sheets []models.Sheet, themeName string) Model {
	if !theme.SetTheme(themeName) {
		theme.SetTheme("catppuccin")
		themeName = "catppuccin"
	}

	styles := ui.InitStyles()

	searchInput := textinput.New()
	searchInput.Placeholder = "search..."
	searchInput.CharLimit = 100
	searchInput.Width = 50

	jumpInput := textinput.New()
	jumpInput.Placeholder = "A100, 500, or 10,5"
	jumpInput.CharLimit = 50
	jumpInput.Width = 30

	exportInput := textinput.New()
	exportInput.Placeholder = "filename.csv or .json"
	exportInput.CharLimit = 100
	exportInput.Width = 40

	editInput := textinput.New()
	editInput.Placeholder = "Enter value or formula"
	editInput.CharLimit = 1000
	editInput.Width = 80

	saveAsInput := textinput.New()
	saveAsInput.Placeholder = "filename.xlsx or .csv"
	saveAsInput.CharLimit = 200
	saveAsInput.Width = 40

	fileFormat := "xlsx"
	if len(filename) > 4 && filename[len(filename)-4:] == ".csv" {
		fileFormat = "csv"
	}

	return Model{
		sheets:       sheets,
		currentSheet: 0,
		searchInput:  searchInput,
		jumpInput:    jumpInput,
		exportInput:  exportInput,
		editInput:    editInput,
		saveAsInput:  saveAsInput,
		help:         help.New(),
		keys:         DefaultKeyMap(),
		filename:     filename,
		themeName:    themeName,
		styles:       styles,
		fileFormat:   fileFormat,
		status: models.StatusMsg{
			Message: "Ready • " + theme.GetCurrentTheme().Name,
			Type:    models.StatusInfo,
		},
	}
}

// GetThemeNames returns available theme names
func GetThemeNames() []string {
	return theme.GetThemeNames()
}

// resetView resets cursor and viewport to initial state
func (m *Model) resetView() {
	m.cursorRow = 0
	m.cursorCol = 0
	m.offsetRow = 0
	m.offsetCol = 0
}

// adjustViewport adjusts the viewport to keep cursor visible
func (m *Model) adjustViewport() {
	visibleRows := ui.Max(1, m.height-9)
	visibleCols := ui.Max(1, (m.width-8)/(ui.MinCellWidth+2))

	if m.cursorRow < m.offsetRow {
		m.offsetRow = m.cursorRow
	} else if m.cursorRow >= m.offsetRow+visibleRows {
		m.offsetRow = m.cursorRow - visibleRows + 1
	}

	if m.cursorCol < m.offsetCol {
		m.offsetCol = m.cursorCol
	} else if m.cursorCol >= m.offsetCol+visibleCols {
		m.offsetCol = m.cursorCol - visibleCols + 1
	}
}

// centerView centers the viewport on the current cursor
func (m *Model) centerView() {
	visibleRows := ui.Max(1, m.height-9)
	visibleCols := ui.Max(1, (m.width-8)/(ui.MinCellWidth+2))

	m.offsetRow = ui.Max(0, m.cursorRow-visibleRows/2)
	m.offsetCol = ui.Max(0, m.cursorCol-visibleCols/2)
}

// isSearchMatch checks if a cell is a search match
func (m *Model) isSearchMatch(row, col int) bool {
	for _, result := range m.searchResults {
		if result.Row == row && result.Col == col {
			return true
		}
	}
	return false
}

// applyTheme applies a new theme and reinitializes styles
func (m *Model) applyTheme(name string) {
	if theme.SetTheme(name) {
		m.themeName = name
		m.styles = ui.InitStyles()
		m.mode = models.ModeNormal
		m.status = models.StatusMsg{
			Message: "Theme: " + theme.GetCurrentTheme().Name,
			Type:    models.StatusSuccess,
		}
	}
}
