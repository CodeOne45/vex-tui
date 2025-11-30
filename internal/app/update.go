package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/odesaur/vex-tui/internal/loader"
	"github.com/odesaur/vex-tui/internal/ui"
	"github.com/odesaur/vex-tui/pkg/models"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case models.ModeSearch:
			return m.updateSearch(msg)
		case models.ModeDetail:
			return m.updateDetail(msg)
		case models.ModeJump:
			return m.updateJump(msg)
		case models.ModeExport:
			return m.updateExport(msg)
		case models.ModeTheme:
			return m.updateTheme(msg)
		case models.ModeChart:
			return m.updateChart(msg)
		case models.ModeSelectRange:
			return m.updateSelectRange(msg)
		default:
			return m.updateNormal(msg)
		}
	}

	return m, nil
}

// updateNormal handles normal mode updates
func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.sheets) == 0 {
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
		return m, nil
	}

	sheet := m.sheets[m.currentSheet]

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Up):
		if m.cursorRow > 0 {
			m.cursorRow--
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Down):
		if m.cursorRow < sheet.MaxRows-1 {
			m.cursorRow++
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Left):
		if m.cursorCol > 0 {
			m.cursorCol--
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Right):
		if m.cursorCol < sheet.MaxCols-1 {
			m.cursorCol++
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.PageDown):
		visibleRows := ui.Max(1, m.height-9)
		m.cursorRow = ui.Min(m.cursorRow+visibleRows, sheet.MaxRows-1)
		m.adjustViewport()

	case key.Matches(msg, m.keys.PageUp):
		visibleRows := ui.Max(1, m.height-9)
		m.cursorRow = ui.Max(m.cursorRow-visibleRows, 0)
		m.adjustViewport()

	case key.Matches(msg, m.keys.Home):
		m.cursorCol = 0
		m.offsetCol = 0

	case key.Matches(msg, m.keys.End):
		m.cursorCol = sheet.MaxCols - 1
		m.adjustViewport()

	case key.Matches(msg, m.keys.FirstCol):
		m.cursorCol = 0
		m.offsetCol = 0

	case key.Matches(msg, m.keys.LastCol):
		m.cursorCol = sheet.MaxCols - 1
		m.adjustViewport()

	case key.Matches(msg, m.keys.NextSheet):
		if m.currentSheet < len(m.sheets)-1 {
			m.currentSheet++
			m.resetView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Next sheet: %s", m.sheets[m.currentSheet].Name),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.PrevSheet):
		if m.currentSheet > 0 {
			m.currentSheet--
			m.resetView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Prev sheet: %s", m.sheets[m.currentSheet].Name),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.Search):
		m.mode = models.ModeSearch
		m.searchInput.Focus()
		m.searchInput.SetValue(m.searchQuery)
		m.searchInput.CursorEnd()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.NextResult):
		if len(m.searchResults) > 0 {
			m.searchIndex = (m.searchIndex + 1) % len(m.searchResults)
			m.jumpToSearchResult()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Match %d/%d", m.searchIndex+1, len(m.searchResults)),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.PrevResult):
		if len(m.searchResults) > 0 {
			m.searchIndex = (m.searchIndex - 1 + len(m.searchResults)) % len(m.searchResults)
			m.jumpToSearchResult()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Match %d/%d", m.searchIndex+1, len(m.searchResults)),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.ClearSearch):
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.searchResults = nil
			m.searchIndex = 0
			m.status = models.StatusMsg{Message: "Search cleared", Type: models.StatusInfo}
		}

	case key.Matches(msg, m.keys.Detail):
		m.mode = models.ModeDetail
		return m, nil

	case key.Matches(msg, m.keys.Jump):
		m.mode = models.ModeJump
		m.jumpInput.Focus()
		m.jumpInput.SetValue("")
		return m, textinput.Blink

	case key.Matches(msg, m.keys.ToggleForm):
		m.showFormulas = !m.showFormulas
		if m.showFormulas {
			m.status = models.StatusMsg{Message: "Showing formulas", Type: models.StatusInfo}
		} else {
			m.status = models.StatusMsg{Message: "Showing values", Type: models.StatusInfo}
		}

	case key.Matches(msg, m.keys.Copy):
		m.copyCell()

	case key.Matches(msg, m.keys.CopyRow):
		m.copyRow()

	case key.Matches(msg, m.keys.Export):
		m.mode = models.ModeExport
		m.exportInput.Focus()
		m.exportInput.SetValue("")
		return m, textinput.Blink

	case key.Matches(msg, m.keys.Theme):
		m.mode = models.ModeTheme
		m.themeIndex = themeIndexByKey(m.themeName)
		return m, nil

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil

	case key.Matches(msg, m.keys.Visualize):
		if !m.isSelecting {
			m.status = models.StatusMsg{Message: "Select range first (V)", Type: models.StatusWarning}
		} else {
			m.mode = models.ModeChart
			m.chartType = 0
		}
		return m, nil

	case key.Matches(msg, m.keys.SelectRange):
		if !m.isSelecting {
			m.selectStart = [2]int{m.cursorRow, m.cursorCol}
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.isSelecting = true
			m.status = models.StatusMsg{Message: "Selection started - Move cursor, press V to finish", Type: models.StatusInfo}
		} else {
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Selected %dx%d range - Press v to visualize",
					abs(m.selectEnd[0]-m.selectStart[0])+1,
					abs(m.selectEnd[1]-m.selectStart[1])+1),
				Type: models.StatusSuccess,
			}
		}
		return m, nil
	}

	return m, nil
}

// updateSearch handles search mode updates
func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.searchInput.Blur()
		return m, nil

	case tea.KeyEnter:
		term := strings.TrimSpace(m.searchInput.Value())
		if term != "" {
			m.searchQuery = term
			sheet := m.sheets[m.currentSheet]
			m.searchResults = loader.SearchSheet(sheet, term)
			m.searchIndex = 0
			if len(m.searchResults) > 0 {
				m.jumpToSearchResult()
				m.status = models.StatusMsg{
					Message: fmt.Sprintf("Found %d results", len(m.searchResults)),
					Type:    models.StatusSuccess,
				}
			} else {
				m.status = models.StatusMsg{Message: "No results found", Type: models.StatusWarning}
			}
		}
		m.mode = models.ModeNormal
		m.searchInput.Blur()
		return m, nil
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

// updateDetail handles detail mode updates
func (m Model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEscape || msg.Type == tea.KeyEnter || msg.String() == "q" {
		m.mode = models.ModeNormal
	}
	return m, nil
}

// updateJump handles jump mode updates
func (m Model) updateJump(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.jumpInput.Blur()
		return m, nil

	case tea.KeyEnter:
		input := strings.TrimSpace(m.jumpInput.Value())
		if input != "" {
			m.jumpToCell(input)
		}
		m.mode = models.ModeNormal
		m.jumpInput.Blur()
		return m, nil
	}

	m.jumpInput, cmd = m.jumpInput.Update(msg)
	return m, cmd
}

// updateExport handles export mode updates
func (m Model) updateExport(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.exportInput.Blur()
		return m, nil

	case tea.KeyEnter:
		filename := strings.TrimSpace(m.exportInput.Value())
		if filename != "" {
			m.exportSheet(filename)
		}
		m.mode = models.ModeNormal
		m.exportInput.Blur()
		return m, nil
	}

	m.exportInput, cmd = m.exportInput.Update(msg)
	return m, cmd
}

// updateTheme handles theme selection mode updates
func (m Model) updateTheme(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = models.ModeNormal
	case "enter":
		m.mode = models.ModeNormal
	case "up", "k":
		if m.themeIndex > 0 {
			m.themeIndex--
		}
		m.applyTheme(themeOptions[m.themeIndex].Key)
	case "down", "j":
		if m.themeIndex < len(themeOptions)-1 {
			m.themeIndex++
		}
		m.applyTheme(themeOptions[m.themeIndex].Key)
	case "1", "2", "3", "4", "5", "6":
		if len(msg.Runes) > 0 {
			idx := int(msg.Runes[0] - '1')
			if idx >= 0 && idx < len(themeOptions) {
				m.themeIndex = idx
				m.applyTheme(themeOptions[idx].Key)
			}
		}
	}
	return m, nil
}

// jumpToSearchResult jumps to the current search result
func (m *Model) jumpToSearchResult() {
	if len(m.searchResults) == 0 {
		return
	}

	result := m.searchResults[m.searchIndex]
	m.cursorRow = result.Row
	m.cursorCol = result.Col
	m.centerView()
}

// jumpToCell jumps to a specific cell based on user input
func (m *Model) jumpToCell(input string) {
	sheet := m.sheets[m.currentSheet]
	input = strings.ToUpper(strings.TrimSpace(input))

	// Format: "A100" (column letter + row number)
	if len(input) > 0 && input[0] >= 'A' && input[0] <= 'Z' {
		col := 0
		row := 0
		i := 0

		// Parse column letters
		for i < len(input) && input[i] >= 'A' && input[i] <= 'Z' {
			col = col*26 + int(input[i]-'A') + 1
			i++
		}
		col-- // Convert to 0-indexed

		// Parse row number
		if i < len(input) {
			if r, err := strconv.Atoi(input[i:]); err == nil {
				row = r - 1 // Convert to 0-indexed
			}
		}

		if row >= 0 && row < sheet.MaxRows && col >= 0 && col < sheet.MaxCols {
			m.cursorRow = row
			m.cursorCol = col
			m.centerView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Position: %s", ui.ColIndexToLetter(col)+fmt.Sprintf("%d", row+1)),
				Type:    models.StatusSuccess,
			}
			return
		}
	}

	// Format: "500" (row number only)
	if row, err := strconv.Atoi(input); err == nil {
		row-- // Convert to 0-indexed
		if row >= 0 && row < sheet.MaxRows {
			m.cursorRow = row
			m.centerView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Row %d", row+1),
				Type:    models.StatusSuccess,
			}
			return
		}
	}

	// Format: "10,5" (row,col)
	if parts := strings.Split(input, ","); len(parts) == 2 {
		if row, err1 := strconv.Atoi(strings.TrimSpace(parts[0])); err1 == nil {
			if col, err2 := strconv.Atoi(strings.TrimSpace(parts[1])); err2 == nil {
				row-- // Convert to 0-indexed
				col-- // Convert to 0-indexed
				if row >= 0 && row < sheet.MaxRows && col >= 0 && col < sheet.MaxCols {
					m.cursorRow = row
					m.cursorCol = col
					m.centerView()
					m.status = models.StatusMsg{
						Message: fmt.Sprintf("Position: %d,%d", row+1, col+1),
						Type:    models.StatusSuccess,
					}
					return
				}
			}
		}
	}

	m.status = models.StatusMsg{Message: "Invalid cell reference", Type: models.StatusError}
}

// copyCell copies the current cell to clipboard
func (m *Model) copyCell() {
	sheet := m.sheets[m.currentSheet]
	if m.cursorRow < len(sheet.Rows) && m.cursorCol < len(sheet.Rows[m.cursorRow]) {
		cell := sheet.Rows[m.cursorRow][m.cursorCol]
		value := cell.Value
		if m.showFormulas && cell.Formula != "" {
			value = "=" + cell.Formula
		}
		if err := clipboard.WriteAll(value); err != nil {
			m.status = models.StatusMsg{Message: "Failed to copy", Type: models.StatusError}
		} else {
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Copied: %s", ui.Truncate(value, 30)),
				Type:    models.StatusSuccess,
			}
		}
	}
}

// copyRow copies the entire current row to clipboard
func (m *Model) copyRow() {
	sheet := m.sheets[m.currentSheet]
	if m.cursorRow < len(sheet.Rows) {
		row := sheet.Rows[m.cursorRow]
		values := make([]string, 0, len(row))
		for _, cell := range row {
			values = append(values, cell.Value)
		}
		rowText := strings.Join(values, "\t")
		if err := clipboard.WriteAll(rowText); err != nil {
			m.status = models.StatusMsg{Message: "Failed to copy row", Type: models.StatusError}
		} else {
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Copied row %d (%d cells)", m.cursorRow+1, len(values)),
				Type:    models.StatusSuccess,
			}
		}
	}
}

// exportSheet exports the current sheet to a file
func (m *Model) exportSheet(filename string) {
	sheet := m.sheets[m.currentSheet]
	var err error

	if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		err = loader.ExportToCSV(sheet, filename)
	} else if strings.HasSuffix(strings.ToLower(filename), ".json") {
		err = loader.ExportToJSON(sheet, filename)
	} else {
		m.status = models.StatusMsg{
			Message: "Use .csv or .json extension",
			Type:    models.StatusError,
		}
		return
	}

	if err != nil {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Export failed: %v", err),
			Type:    models.StatusError,
		}
	} else {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Exported to %s", filename),
			Type:    models.StatusSuccess,
		}
	}
}

// updateChart handles chart visualization mode
func (m Model) updateChart(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = models.ModeNormal
	case "1":
		m.chartType = 0 // Bar
	case "2":
		m.chartType = 1 // Line
	case "3":
		m.chartType = 2 // Sparkline
	case "4":
		m.chartType = 3 // Pie
	}
	return m, nil
}

// updateSelectRange handles range selection mode
func (m Model) updateSelectRange(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	sheet := m.sheets[m.currentSheet]

	switch {
	case key.Matches(msg, m.keys.Up):
		if m.cursorRow > 0 {
			m.cursorRow--
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.adjustViewport()
		}
	case key.Matches(msg, m.keys.Down):
		if m.cursorRow < sheet.MaxRows-1 {
			m.cursorRow++
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.adjustViewport()
		}
	case key.Matches(msg, m.keys.Left):
		if m.cursorCol > 0 {
			m.cursorCol--
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.adjustViewport()
		}
	case key.Matches(msg, m.keys.Right):
		if m.cursorCol < sheet.MaxCols-1 {
			m.cursorCol++
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.adjustViewport()
		}
	case key.Matches(msg, m.keys.SelectRange):
		m.mode = models.ModeNormal
		m.status = models.StatusMsg{
			Message: "Selected range - Press v to visualize",
			Type:    models.StatusSuccess,
		}
	case msg.Type == tea.KeyEscape:
		m.isSelecting = false
		m.mode = models.ModeNormal
		m.status = models.StatusMsg{Message: "Selection cancelled", Type: models.StatusInfo}
	}

	return m, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
