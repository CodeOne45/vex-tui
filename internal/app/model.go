package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/odesaur/vex-tui/v2/internal/loader"
	"github.com/odesaur/vex-tui/v2/internal/theme"
	"github.com/odesaur/vex-tui/v2/internal/ui"
	"github.com/odesaur/vex-tui/v2/pkg/models"
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
	files         []string
	filteredFiles []string
	fileIndex     int
	fileFilter    textinput.Model
	searchInput   textinput.Model
	jumpInput     textinput.Model
	exportInput   textinput.Model
	searchQuery   string
	searchResults []models.Cell
	searchIndex   int
	showFormulas  bool
	status        models.StatusMsg
	help          help.Model
	keys          KeyMap
	filename      string
	themeName     string
	themeIndex    int
	styles        *ui.Styles

	// Chart visualization
	chartType   int
	selectStart [2]int // [row, col]
	selectEnd   [2]int // [row, col]
	isSelecting bool
}

// NewModel creates a new application model
func NewModel(filename string, sheets []models.Sheet, themeName string) Model {
	// Set theme
	if !theme.SetTheme(themeName) {
		theme.SetTheme("rose-pine")
		themeName = "rose-pine"
	}

	// Initialize styles
	styles := ui.InitStyles()

	// Create input fields
	fileFilter := textinput.New()
	fileFilter.Placeholder = "Filter files..."
	fileFilter.CharLimit = 200
	fileFilter.Width = 50

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

	model := Model{
		sheets:       sheets,
		currentSheet: 0,
		fileFilter:   fileFilter,
		searchInput:  searchInput,
		jumpInput:    jumpInput,
		exportInput:  exportInput,
		help:         help.New(),
		keys:         DefaultKeyMap(),
		filename:     filename,
		themeName:    themeName,
		themeIndex:   themeIndexByKey(themeName),
		styles:       styles,
		status: models.StatusMsg{
			Message: "Ready | " + theme.GetCurrentTheme().Name,
			Type:    models.StatusInfo,
		},
	}

	if filename == "" || len(sheets) == 0 {
		model.mode = models.ModeFilePicker
		model.status = models.StatusMsg{Message: "Select a file to open", Type: models.StatusInfo}
		model.fileFilter.Focus()
		model.initFilePicker()
	} else {
		model.mode = models.ModeNormal
	}

	return model
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

// filterFiles filters the file list based on the current filter value.
func (m *Model) filterFiles() {
	query := strings.ToLower(strings.TrimSpace(m.fileFilter.Value()))
	if query == "" {
		m.filteredFiles = m.files
		if m.fileIndex >= len(m.filteredFiles) {
			m.fileIndex = 0
		}
		return
	}

	filtered := make([]string, 0, len(m.files))
	for _, f := range m.files {
		if strings.Contains(strings.ToLower(f), query) {
			filtered = append(filtered, f)
		}
	}
	m.filteredFiles = filtered
	if m.fileIndex >= len(m.filteredFiles) {
		m.fileIndex = 0
	}
}

// adjustViewport adjusts the viewport to keep cursor visible
func (m *Model) adjustViewport() {
	visibleRows := ui.Max(1, m.height-9)
	visibleCols := ui.Max(1, (m.width-8)/(ui.MinCellWidth+2))
	if len(m.sheets) > 0 {
		visibleCols = ui.Min(visibleCols, m.sheets[m.currentSheet].MaxCols)
	}

	// Adjust vertical
	if m.cursorRow < m.offsetRow {
		m.offsetRow = m.cursorRow
	} else if m.cursorRow >= m.offsetRow+visibleRows {
		m.offsetRow = m.cursorRow - visibleRows + 1
	}

	// Adjust horizontal
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
		m.themeIndex = themeIndexByKey(name)
		m.styles = ui.InitStyles()
		m.status = models.StatusMsg{
			Message: "Theme: " + theme.GetCurrentTheme().Name,
			Type:    models.StatusSuccess,
		}
	}
}

// initFilePicker loads supported files from the current directory.
func (m *Model) initFilePicker() {
	files, err := loader.DiscoverFiles(".", 400)
	if err != nil {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("File scan failed: %v", err),
			Type:    models.StatusError,
		}
		return
	}
	m.files = files
	m.filteredFiles = files
	m.fileIndex = 0
}

// applyFileSelection loads the selected file and switches to normal mode.
func (m *Model) applyFileSelection(path string) {
	sheets, err := loader.LoadFile(path)
	if err != nil {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Load failed: %v", err),
			Type:    models.StatusError,
		}
		return
	}
	if len(sheets) == 0 {
		m.status = models.StatusMsg{
			Message: "No sheets found in file",
			Type:    models.StatusError,
		}
		return
	}

	m.sheets = sheets
	m.filename = path
	m.currentSheet = 0
	m.cursorRow, m.cursorCol = 0, 0
	m.offsetRow, m.offsetCol = 0, 0
	m.searchResults = nil
	m.searchQuery = ""
	m.searchIndex = 0
	m.fileFilter.Blur()
	m.mode = models.ModeNormal
	m.status = models.StatusMsg{
		Message: "Ready | " + theme.GetCurrentTheme().Name,
		Type:    models.StatusSuccess,
	}
}
