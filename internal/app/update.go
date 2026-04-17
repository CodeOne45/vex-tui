package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/CodeOne45/vex-tui/internal/loader"
	"github.com/CodeOne45/vex-tui/internal/ui"
	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
		case models.ModeEdit:
			return m.updateEdit(msg)
		case models.ModeSaveAs:
			return m.updateSaveAs(msg)
		case models.ModeFilter:
			return m.updateFilter(msg)
		case models.ModeDataProfile:
			m.mode = models.ModeNormal
			return m, nil
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

	// Handle pending key sequences (dd, dc, gg)
	if m.pendingKey != "" {
		pending := m.pendingKey
		m.pendingKey = ""
		k := msg.String()

		switch pending {
		case "d":
			switch k {
			case "d":
				m.quitConfirm = false
				m.deleteRow()
				return m, nil
			case "c":
				m.quitConfirm = false
				m.deleteColumn()
				return m, nil
			}
			// Not a valid sequence, ignore
			return m, nil
		case "g":
			if k == "g" {
				m.quitConfirm = false
				m.cursorRow = 0
				m.cursorCol = 0
				m.offsetRow = 0
				m.offsetCol = 0
				m.status = models.StatusMsg{Message: "Top of file", Type: models.StatusInfo}
				return m, nil
			}
			// Not a valid sequence, ignore
			return m, nil
		}
		return m, nil
	}

	// Check for pending key triggers
	switch msg.String() {
	case "d":
		m.pendingKey = "d"
		return m, nil
	case "g":
		m.pendingKey = "g"
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		if m.modified && !m.quitConfirm {
			m.quitConfirm = true
			m.status = models.StatusMsg{
				Message: "Unsaved changes! Press q again to quit without saving, or Ctrl+S to save",
				Type:    models.StatusWarning,
			}
			return m, nil
		}
		return m, tea.Quit

	case key.Matches(msg, m.keys.Up):
		m.quitConfirm = false
		if m.cursorRow > 0 {
			m.cursorRow--
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Down):
		m.quitConfirm = false
		if m.cursorRow < sheet.MaxRows-1 {
			m.cursorRow++
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Left):
		m.quitConfirm = false
		if m.cursorCol > 0 {
			m.cursorCol--
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.Right):
		m.quitConfirm = false
		if m.cursorCol < sheet.MaxCols-1 {
			m.cursorCol++
			m.adjustViewport()
		}

	case key.Matches(msg, m.keys.PageDown):
		m.quitConfirm = false
		visibleRows := ui.Max(1, m.height-9)
		m.cursorRow = ui.Min(m.cursorRow+visibleRows, sheet.MaxRows-1)
		m.adjustViewport()

	case key.Matches(msg, m.keys.PageUp):
		m.quitConfirm = false
		visibleRows := ui.Max(1, m.height-9)
		m.cursorRow = ui.Max(m.cursorRow-visibleRows, 0)
		m.adjustViewport()

	case key.Matches(msg, m.keys.Home):
		m.quitConfirm = false
		m.cursorCol = 0
		m.offsetCol = 0

	case key.Matches(msg, m.keys.End):
		m.quitConfirm = false
		m.cursorCol = sheet.MaxCols - 1
		m.adjustViewport()

	case key.Matches(msg, m.keys.GotoBottom):
		m.quitConfirm = false
		m.cursorRow = sheet.MaxRows - 1
		m.adjustViewport()
		m.status = models.StatusMsg{Message: "Bottom of file", Type: models.StatusInfo}

	case key.Matches(msg, m.keys.NextSheet):
		m.quitConfirm = false
		if m.currentSheet < len(m.sheets)-1 {
			m.currentSheet++
			m.resetView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("→ %s", m.sheets[m.currentSheet].Name),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.PrevSheet):
		m.quitConfirm = false
		if m.currentSheet > 0 {
			m.currentSheet--
			m.resetView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("← %s", m.sheets[m.currentSheet].Name),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.Search):
		m.quitConfirm = false
		m.mode = models.ModeSearch
		m.searchInput.Focus()
		m.searchInput.SetValue(m.searchQuery)
		m.searchInput.CursorEnd()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.NextResult):
		m.quitConfirm = false
		if len(m.searchResults) > 0 {
			m.searchIndex = (m.searchIndex + 1) % len(m.searchResults)
			m.jumpToSearchResult()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Match %d/%d", m.searchIndex+1, len(m.searchResults)),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.PrevResult):
		m.quitConfirm = false
		if len(m.searchResults) > 0 {
			m.searchIndex = (m.searchIndex - 1 + len(m.searchResults)) % len(m.searchResults)
			m.jumpToSearchResult()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Match %d/%d", m.searchIndex+1, len(m.searchResults)),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.ClearSearch):
		m.quitConfirm = false
		if m.isSelecting {
			m.isSelecting = false
			m.status = models.StatusMsg{Message: "Selection cancelled", Type: models.StatusInfo}
		} else if m.filterActive {
			m.clearFilter()
		} else if m.searchQuery != "" {
			m.searchQuery = ""
			m.searchResults = nil
			m.searchIndex = 0
			m.status = models.StatusMsg{Message: "Search cleared", Type: models.StatusInfo}
		}

	case key.Matches(msg, m.keys.Detail):
		m.quitConfirm = false
		m.mode = models.ModeDetail
		return m, nil

	case key.Matches(msg, m.keys.Jump):
		m.quitConfirm = false
		m.mode = models.ModeJump
		m.jumpInput.Focus()
		m.jumpInput.SetValue("")
		return m, textinput.Blink

	case key.Matches(msg, m.keys.ToggleForm):
		m.quitConfirm = false
		m.showFormulas = !m.showFormulas
		if m.showFormulas {
			m.status = models.StatusMsg{Message: "Showing formulas", Type: models.StatusInfo}
		} else {
			m.status = models.StatusMsg{Message: "Showing values", Type: models.StatusInfo}
		}

	case key.Matches(msg, m.keys.Copy):
		m.quitConfirm = false
		if m.isSelecting {
			m.copyRange()
		} else {
			m.copyCell()
		}

	case key.Matches(msg, m.keys.CopyRow):
		m.quitConfirm = false
		m.copyRow()

	case key.Matches(msg, m.keys.Export):
		m.quitConfirm = false
		m.mode = models.ModeExport
		m.exportInput.Focus()
		m.exportInput.SetValue("")
		return m, textinput.Blink

	case key.Matches(msg, m.keys.Theme):
		m.quitConfirm = false
		m.mode = models.ModeTheme
		return m, nil

	case key.Matches(msg, m.keys.Help):
		m.quitConfirm = false
		m.help.ShowAll = !m.help.ShowAll
		return m, nil

	case key.Matches(msg, m.keys.Visualize):
		m.quitConfirm = false
		if !m.isSelecting {
			m.status = models.StatusMsg{Message: "Select range first (V)", Type: models.StatusWarning}
		} else {
			m.mode = models.ModeChart
			m.chartType = 0
		}
		return m, nil

	case key.Matches(msg, m.keys.SelectRange):
		m.quitConfirm = false
		if !m.isSelecting {
			m.selectStart = [2]int{m.cursorRow, m.cursorCol}
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.isSelecting = true
			m.status = models.StatusMsg{Message: "Selection started - Move cursor, press V to finish", Type: models.StatusInfo}
		} else {
			m.selectEnd = [2]int{m.cursorRow, m.cursorCol}
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Selected %dx%d range - Press v to visualize or Ctrl+A to apply formula",
					abs(m.selectEnd[0]-m.selectStart[0])+1,
					abs(m.selectEnd[1]-m.selectStart[1])+1),
				Type: models.StatusSuccess,
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		m.quitConfirm = false
		m.startEdit()
		return m, textarea.Blink

	case key.Matches(msg, m.keys.Delete):
		m.quitConfirm = false
		m.deleteCell()

	case key.Matches(msg, m.keys.InsertRow):
		m.quitConfirm = false
		m.insertRow()

	case key.Matches(msg, m.keys.InsertCol):
		m.quitConfirm = false
		m.insertColumn()

	case key.Matches(msg, m.keys.Paste):
		m.quitConfirm = false
		m.pasteCell()

	case key.Matches(msg, m.keys.Save):
		m.quitConfirm = false
		m.saveFile()

	case key.Matches(msg, m.keys.SaveAs):
		m.quitConfirm = false
		m.mode = models.ModeSaveAs
		m.saveAsInput.Focus()
		m.saveAsInput.SetValue(m.filename)
		return m, textinput.Blink

	case key.Matches(msg, m.keys.FillDown):
		m.quitConfirm = false
		m.fillDown()

	case key.Matches(msg, m.keys.FillRight):
		m.quitConfirm = false
		m.fillRight()

	case key.Matches(msg, m.keys.ApplyFormula):
		m.quitConfirm = false
		m.applyFormulaToRange()

	case key.Matches(msg, m.keys.SortAsc):
		m.quitConfirm = false
		m.sortByColumn(m.cursorCol, true)

	case key.Matches(msg, m.keys.SortDesc):
		m.quitConfirm = false
		m.sortByColumn(m.cursorCol, false)

	case key.Matches(msg, m.keys.SortReset):
		m.quitConfirm = false
		m.resetSort()

	case key.Matches(msg, m.keys.Filter):
		m.quitConfirm = false
		m.mode = models.ModeFilter
		m.filterColAll = false
		m.filterInput.Placeholder = fmt.Sprintf("Filter col %s — text, >100, <50, =val  (Tab=all cols)",
			ui.ColIndexToLetter(m.cursorCol))
		m.filterInput.SetValue("")
		m.filterInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.DataProfile):
		m.quitConfirm = false
		m.mode = models.ModeDataProfile
		return m, nil

	case key.Matches(msg, m.keys.FreezeHeader):
		m.quitConfirm = false
		m.freezeHeader = !m.freezeHeader
		if m.freezeHeader {
			if m.cursorRow == 0 && sheet.MaxRows > 1 {
				m.cursorRow = 1
			}
			m.status = models.StatusMsg{Message: "Header frozen — row 1 always visible", Type: models.StatusInfo}
		} else {
			m.status = models.StatusMsg{Message: "Header unfrozen", Type: models.StatusInfo}
		}

	case key.Matches(msg, m.keys.ColWidthInc):
		m.quitConfirm = false
		if sheet.ColWidths == nil {
			sheet.ColWidths = make(map[int]int)
			m.sheets[m.currentSheet] = sheet
		}
		width := sheet.ColWidths[m.cursorCol]
		if width == 0 {
			width = ui.MinCellWidth
		}
		if width < ui.MaxCellWidth {
			sheet.ColWidths[m.cursorCol] = width + 1
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Width: %d", width+1),
				Type:    models.StatusInfo,
			}
		}

	case key.Matches(msg, m.keys.ColWidthDec):
		m.quitConfirm = false
		if sheet.ColWidths == nil {
			sheet.ColWidths = make(map[int]int)
			m.sheets[m.currentSheet] = sheet
		}
		width := sheet.ColWidths[m.cursorCol]
		if width == 0 {
			width = ui.MinCellWidth
		}
		if width > 2 {
			sheet.ColWidths[m.cursorCol] = width - 1
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Width: %d", width-1),
				Type:    models.StatusInfo,
			}
		}
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
	case "1":
		m.applyTheme("catppuccin")
	case "2":
		m.applyTheme("nord")
	case "3":
		m.applyTheme("rose-pine")
	case "4":
		m.applyTheme("tokyo-night")
	case "5":
		m.applyTheme("gruvbox")
	case "6":
		m.applyTheme("dracula")
	case "7":
		m.applyTheme("catppuccin-latte")
	case "8":
		m.applyTheme("solarized-light")
	case "9":
		m.applyTheme("github-light")
	case "0":
		m.applyTheme("one-light")
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

	if len(input) > 0 && input[0] >= 'A' && input[0] <= 'Z' {
		col := 0
		row := 0
		i := 0

		for i < len(input) && input[i] >= 'A' && input[i] <= 'Z' {
			col = col*26 + int(input[i]-'A') + 1
			i++
		}
		col--

		if i < len(input) {
			if r, err := strconv.Atoi(input[i:]); err == nil {
				row = r - 1
			}
		}

		if row >= 0 && row < sheet.MaxRows && col >= 0 && col < sheet.MaxCols {
			m.cursorRow = row
			m.cursorCol = col
			m.centerView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("→ %s", ui.ColIndexToLetter(col)+fmt.Sprintf("%d", row+1)),
				Type:    models.StatusSuccess,
			}
			return
		}
	}

	if row, err := strconv.Atoi(input); err == nil {
		row--
		if row >= 0 && row < sheet.MaxRows {
			m.cursorRow = row
			m.centerView()
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("→ Row %d", row+1),
				Type:    models.StatusSuccess,
			}
			return
		}
	}

	if parts := strings.Split(input, ","); len(parts) == 2 {
		if row, err1 := strconv.Atoi(strings.TrimSpace(parts[0])); err1 == nil {
			if col, err2 := strconv.Atoi(strings.TrimSpace(parts[1])); err2 == nil {
				row--
				col--
				if row >= 0 && row < sheet.MaxRows && col >= 0 && col < sheet.MaxCols {
					m.cursorRow = row
					m.cursorCol = col
					m.centerView()
					m.status = models.StatusMsg{
						Message: fmt.Sprintf("→ %d,%d", row+1, col+1),
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

// copyRange copies the selected range to clipboard (tab-separated columns, newline-separated rows)
func (m *Model) copyRange() {
	sheet := m.sheets[m.currentSheet]

	startRow := m.selectStart[0]
	endRow := m.selectEnd[0]
	startCol := m.selectStart[1]
	endCol := m.selectEnd[1]

	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	var rows []string
	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		var cells []string
		for col := startCol; col <= endCol; col++ {
			if col < len(sheet.Rows[row]) {
				cells = append(cells, sheet.Rows[row][col].Value)
			} else {
				cells = append(cells, "")
			}
		}
		rows = append(rows, strings.Join(cells, "\t"))
	}

	rangeText := strings.Join(rows, "\n")
	if err := clipboard.WriteAll(rangeText); err != nil {
		m.status = models.StatusMsg{Message: "Failed to copy range", Type: models.StatusError}
	} else {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Copied %dx%d range", endRow-startRow+1, endCol-startCol+1),
			Type:    models.StatusSuccess,
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
			Message: fmt.Sprintf("✓ Exported to %s", filename),
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
		m.chartType = 0
	case "2":
		m.chartType = 1
	case "3":
		m.chartType = 2
	case "4":
		m.chartType = 3
	case "p", "P":
		m.exportChartImage()
	}
	return m, nil
}

// freezeLookup finds the freeze binary across macOS and Linux.
// Search order:
//  1. PATH (covers correctly-linked installs everywhere)
//  2. Fixed well-known paths (Homebrew symlinks, Linuxbrew, user-local)
//  3. Homebrew / Linuxbrew Cellar walk (catches broken symlinks)
func freezeLookup() (string, bool) {
	// 1. Respect PATH first — works for any correctly installed binary.
	if p, err := exec.LookPath("freeze"); err == nil {
		return p, true
	}

	// 2. Fixed paths.
	home, _ := os.UserHomeDir()
	fixed := []string{
		"/opt/homebrew/bin/freeze",                      // macOS Apple Silicon (Homebrew symlink)
		"/opt/homebrew/opt/freeze/bin/freeze",           // macOS Apple Silicon (Homebrew opt)
		"/usr/local/bin/freeze",                         // macOS Intel (Homebrew symlink) + common Linux
		"/home/linuxbrew/.linuxbrew/bin/freeze",         // Linuxbrew
		filepath.Join(home, ".local", "bin", "freeze"),  // Linux non-root installs
		filepath.Join(home, "go", "bin", "freeze"),      // go install
		filepath.Join(home, ".brew", "bin", "freeze"),   // custom Homebrew prefix
	}
	for _, p := range fixed {
		if p == "" {
			continue
		}
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}

	// 3. Cellar walk — catches Homebrew installs whose bin symlink is missing.
	cellarRoots := []string{
		"/opt/homebrew/Cellar/freeze",              // Apple Silicon
		"/usr/local/Cellar/freeze",                 // Intel macOS
		"/home/linuxbrew/.linuxbrew/Cellar/freeze", // Linuxbrew
	}
	for _, root := range cellarRoots {
		versions, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, v := range versions {
			p := filepath.Join(root, v.Name(), "bin", "freeze")
			if info, err := os.Stat(p); err == nil && !info.IsDir() {
				return p, true
			}
		}
	}

	return "", false
}

// exportChartImage renders the current chart and saves it as a PNG via the
// `freeze` CLI (github.com/charmbracelet/freeze). Falls back to a clear
// install hint if freeze is not found.
func (m *Model) exportChartImage() {
	freezePath, found := freezeLookup()
	if !found {
		m.status = models.StatusMsg{
			Message: "freeze not found — install: brew install charmbracelet/tap/freeze",
			Type:    models.StatusWarning,
		}
		return
	}

	// Build the output filename next to the source file.
	base := strings.TrimSuffix(filepath.Base(m.filename), filepath.Ext(m.filename))
	outFile := base + "-chart.png"

	// Render the chart modal (returns an ANSI string from lipgloss).
	chartANSI := m.renderChart()

	cmd := exec.Command(freezePath,
		"--output", outFile,
		"--padding", "20,30",
		"--border.radius", "8",
		"--shadow",
	)
	cmd.Stdin = strings.NewReader(chartANSI)

	if err := cmd.Run(); err != nil {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Export failed: %v", err),
			Type:    models.StatusError,
		}
		return
	}

	m.status = models.StatusMsg{
		Message: fmt.Sprintf("✓ Saved %s", outFile),
		Type:    models.StatusSuccess,
	}
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
			Message: "Selected range - Press v to visualize or Ctrl+A to apply formula",
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

// resetSort restores the pre-sort row order.
func (m *Model) resetSort() {
	if !m.preSortSaved {
		m.status = models.StatusMsg{Message: "No sort to undo", Type: models.StatusInfo}
		return
	}
	sheet := &m.sheets[m.currentSheet]
	sheet.Rows = m.preSortRows
	sheet.MaxRows = m.preSortMaxRows
	m.preSortSaved = false
	m.preSortRows = nil
	m.sortActive = false
	m.cursorRow = 0
	m.offsetRow = 0
	m.status = models.StatusMsg{Message: "Sort cleared — original order restored", Type: models.StatusInfo}
}

// updateFilter handles the filter input mode.
func (m Model) updateFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.filterInput.Blur()
		return m, nil

	case tea.KeyTab:
		m.filterColAll = !m.filterColAll
		if m.filterColAll {
			m.filterInput.Placeholder = "Filter ALL columns — text, >100, <50, =val"
		} else {
			m.filterInput.Placeholder = fmt.Sprintf("Filter col %s — text, >100, <50, =val  (Tab=all cols)",
				ui.ColIndexToLetter(m.cursorCol))
		}
		return m, nil

	case tea.KeyEnter:
		expr := strings.TrimSpace(m.filterInput.Value())
		m.filterInput.Blur()
		m.mode = models.ModeNormal
		if expr == "" {
			m.clearFilter()
		} else {
			m.applyFilter(expr)
		}
		return m, nil
	}

	m.filterInput, cmd = m.filterInput.Update(msg)
	return m, cmd
}

// applyFilter filters the current sheet rows by the given expression.
// Expression formats: plain text (contains), >N, <N, >=N, <=N, =exact.
func (m *Model) applyFilter(expr string) {
	sheet := &m.sheets[m.currentSheet]

	// Save original rows if not already filtered
	if !m.filterActive {
		rows := make([][]models.Cell, len(sheet.Rows))
		copy(rows, sheet.Rows)
		m.preFilterRows = rows
		m.preFilterMax = sheet.MaxRows
	}

	col := m.cursorCol
	matchRow := func(row []models.Cell) bool {
		if m.filterColAll {
			for _, cell := range row {
				if matchExpr(cell.Value, expr) {
					return true
				}
			}
			return false
		}
		if col < len(row) {
			return matchExpr(row[col].Value, expr)
		}
		return false
	}

	var filtered [][]models.Cell
	for _, row := range m.preFilterRows {
		if matchRow(row) {
			filtered = append(filtered, row)
		}
	}

	if len(filtered) == 0 {
		m.status = models.StatusMsg{Message: "No rows match — filter not applied", Type: models.StatusWarning}
		return
	}

	sheet.Rows = filtered
	sheet.MaxRows = len(filtered)
	m.filterActive = true
	m.filterQuery = expr
	m.cursorRow = 0
	m.offsetRow = 0

	scope := fmt.Sprintf("col %s", ui.ColIndexToLetter(col))
	if m.filterColAll {
		scope = "all cols"
	}
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Filter [%s] on %s — %d/%d rows  (Esc or f + Enter to clear)",
			expr, scope, len(filtered), m.preFilterMax),
		Type: models.StatusSuccess,
	}
}

// clearFilter restores the pre-filter rows.
func (m *Model) clearFilter() {
	if !m.filterActive {
		return
	}
	sheet := &m.sheets[m.currentSheet]
	sheet.Rows = m.preFilterRows
	sheet.MaxRows = m.preFilterMax
	m.filterActive = false
	m.filterQuery = ""
	m.preFilterRows = nil
	m.cursorRow = 0
	m.offsetRow = 0
	m.status = models.StatusMsg{Message: "Filter cleared", Type: models.StatusInfo}
}

// matchExpr tests whether a cell value satisfies the filter expression.
func matchExpr(value, expr string) bool {
	// Numeric comparisons
	for _, prefix := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(expr, prefix) {
			operand := strings.TrimSpace(expr[len(prefix):])
			target, err1 := strconv.ParseFloat(operand, 64)
			actual, err2 := strconv.ParseFloat(value, 64)
			if err1 != nil || err2 != nil {
				if prefix == "=" {
					return strings.EqualFold(value, operand)
				}
				return false
			}
			switch prefix {
			case ">=":
				return actual >= target
			case "<=":
				return actual <= target
			case ">":
				return actual > target
			case "<":
				return actual < target
			case "=":
				return actual == target
			}
		}
	}
	// Default: case-insensitive contains
	return strings.Contains(strings.ToLower(value), strings.ToLower(expr))
}

// sortByColumn sorts the sheet rows by the given column. If freezeHeader is
// enabled the first row is treated as a header and kept in place.
func (m *Model) sortByColumn(col int, ascending bool) {
	sheet := &m.sheets[m.currentSheet]
	if len(sheet.Rows) == 0 {
		return
	}

	// Save original order before first sort
	if !m.preSortSaved {
		saved := make([][]models.Cell, len(sheet.Rows))
		copy(saved, sheet.Rows)
		m.preSortRows = saved
		m.preSortMaxRows = sheet.MaxRows
		m.preSortSaved = true
	}

	startRow := 0
	if m.freezeHeader && len(sheet.Rows) > 1 {
		startRow = 1
	}

	data := sheet.Rows[startRow:]
	sort.SliceStable(data, func(i, j int) bool {
		vi, vj := "", ""
		if col < len(data[i]) {
			vi = data[i][col].Value
		}
		if col < len(data[j]) {
			vj = data[j][col].Value
		}
		ni, erri := strconv.ParseFloat(vi, 64)
		nj, errj := strconv.ParseFloat(vj, 64)
		if erri == nil && errj == nil {
			if ascending {
				return ni < nj
			}
			return ni > nj
		}
		if ascending {
			return strings.ToLower(vi) < strings.ToLower(vj)
		}
		return strings.ToLower(vi) > strings.ToLower(vj)
	})

	// Fix row indices after sort
	for i := range sheet.Rows {
		for j := range sheet.Rows[i] {
			sheet.Rows[i][j].Row = i
		}
	}

	m.modified = true
	m.sortCol = col
	m.sortAsc = ascending
	m.sortActive = true
	dir := "↑"
	if !ascending {
		dir = "↓"
	}
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Sorted by %s %s", ui.ColIndexToLetter(col), dir),
		Type:    models.StatusSuccess,
	}
}
