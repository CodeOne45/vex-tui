package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CodeOne45/vex-tui/internal/loader"
	"github.com/CodeOne45/vex-tui/internal/theme"
	"github.com/CodeOne45/vex-tui/internal/ui"
	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updateEdit handles edit mode updates
func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.editInput.Blur()
		m.isEditing = false
		return m, nil

	case tea.KeyEnter:
		m.commitEdit()
		m.mode = models.ModeNormal
		m.editInput.Blur()
		m.isEditing = false
		return m, nil

	case tea.KeyTab:
		m.commitEdit()
		sheet := m.sheets[m.currentSheet]
		if m.cursorCol < sheet.MaxCols-1 {
			m.cursorCol++
		} else if m.cursorRow < sheet.MaxRows-1 {
			m.cursorCol = 0
			m.cursorRow++
		}
		m.adjustViewport()
		m.startEdit()
		return m, textinput.Blink

	case tea.KeyShiftTab:
		m.commitEdit()
		if m.cursorCol > 0 {
			m.cursorCol--
		} else if m.cursorRow > 0 {
			sheet := m.sheets[m.currentSheet]
			m.cursorCol = sheet.MaxCols - 1
			m.cursorRow--
		}
		m.adjustViewport()
		m.startEdit()
		return m, textinput.Blink

	case tea.KeyCtrlC:
		m.mode = models.ModeNormal
		m.editInput.Blur()
		m.isEditing = false
		return m, nil
	}

	m.editInput, cmd = m.editInput.Update(msg)
	return m, cmd
}

// startEdit initializes edit mode for current cell
func (m *Model) startEdit() {
	sheet := &m.sheets[m.currentSheet]
	if m.cursorRow >= len(sheet.Rows) {
		sheet.Rows = append(sheet.Rows, make([][]models.Cell, m.cursorRow-len(sheet.Rows)+1)...)
	}
	if m.cursorRow >= len(sheet.Rows) {
		for i := len(sheet.Rows); i <= m.cursorRow; i++ {
			sheet.Rows = append(sheet.Rows, make([]models.Cell, sheet.MaxCols))
		}
	}
	if m.cursorCol >= len(sheet.Rows[m.cursorRow]) {
		oldLen := len(sheet.Rows[m.cursorRow])
		newRow := make([]models.Cell, sheet.MaxCols)
		copy(newRow, sheet.Rows[m.cursorRow])
		for i := oldLen; i < len(newRow); i++ {
			newRow[i] = models.Cell{Row: m.cursorRow, Col: i}
		}
		sheet.Rows[m.cursorRow] = newRow
	}

	cell := sheet.Rows[m.cursorRow][m.cursorCol]
	value := cell.Value
	if m.showFormulas && cell.Formula != "" {
		value = "=" + cell.Formula
	}

	m.editInput.SetValue(value)
	m.editInput.CursorEnd()
	m.editInput.Focus()
	m.isEditing = true
	m.mode = models.ModeEdit
	m.modified = true
}

// commitEdit saves the current edit to the cell
func (m *Model) commitEdit() {
	sheet := &m.sheets[m.currentSheet]
	value := strings.TrimSpace(m.editInput.Value())

	if m.cursorRow >= len(sheet.Rows) {
		for i := len(sheet.Rows); i <= m.cursorRow; i++ {
			sheet.Rows = append(sheet.Rows, make([]models.Cell, sheet.MaxCols))
		}
	}
	if m.cursorCol >= len(sheet.Rows[m.cursorRow]) {
		oldLen := len(sheet.Rows[m.cursorRow])
		newRow := make([]models.Cell, sheet.MaxCols)
		copy(newRow, sheet.Rows[m.cursorRow])
		for i := oldLen; i < len(newRow); i++ {
			newRow[i] = models.Cell{Row: m.cursorRow, Col: i}
		}
		sheet.Rows[m.cursorRow] = newRow
	}

	cell := &sheet.Rows[m.cursorRow][m.cursorCol]
	cell.Row = m.cursorRow
	cell.Col = m.cursorCol

	if strings.HasPrefix(value, "=") {
		cell.Formula = value[1:]
		cell.Value = m.evaluateFormula(cell.Formula)
	} else {
		cell.Value = value
		cell.Formula = ""
	}

	m.modified = true
	m.recalculateFormulas()
}

// deleteCell deletes content of current cell
func (m *Model) deleteCell() {
	sheet := &m.sheets[m.currentSheet]
	if m.cursorRow < len(sheet.Rows) && m.cursorCol < len(sheet.Rows[m.cursorRow]) {
		cell := &sheet.Rows[m.cursorRow][m.cursorCol]
		cell.Value = ""
		cell.Formula = ""
		m.modified = true
		m.recalculateFormulas()
		m.status = models.StatusMsg{Message: "Cell cleared", Type: models.StatusSuccess}
	}
}

// deleteRow deletes entire current row
func (m *Model) deleteRow() {
	sheet := &m.sheets[m.currentSheet]
	if m.cursorRow < len(sheet.Rows) {
		sheet.Rows = append(sheet.Rows[:m.cursorRow], sheet.Rows[m.cursorRow+1:]...)
		sheet.MaxRows = len(sheet.Rows)
		if m.cursorRow >= sheet.MaxRows && sheet.MaxRows > 0 {
			m.cursorRow = sheet.MaxRows - 1
		}
		m.modified = true
		m.recalculateFormulas()
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Row %d deleted", m.cursorRow+1),
			Type:    models.StatusSuccess,
		}
	}
}

// deleteColumn deletes entire current column
func (m *Model) deleteColumn() {
	sheet := &m.sheets[m.currentSheet]
	for i := range sheet.Rows {
		if m.cursorCol < len(sheet.Rows[i]) {
			sheet.Rows[i] = append(sheet.Rows[i][:m.cursorCol], sheet.Rows[i][m.cursorCol+1:]...)
		}
	}
	sheet.MaxCols--
	if m.cursorCol >= sheet.MaxCols && sheet.MaxCols > 0 {
		m.cursorCol = sheet.MaxCols - 1
	}
	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Column %s deleted", ui.ColIndexToLetter(m.cursorCol)),
		Type:    models.StatusSuccess,
	}
}

// insertRow inserts a new row at current position
func (m *Model) insertRow() {
	sheet := &m.sheets[m.currentSheet]
	newRow := make([]models.Cell, sheet.MaxCols)
	for i := range newRow {
		newRow[i] = models.Cell{Row: m.cursorRow, Col: i}
	}
	sheet.Rows = append(sheet.Rows[:m.cursorRow], append([][]models.Cell{newRow}, sheet.Rows[m.cursorRow:]...)...)
	sheet.MaxRows = len(sheet.Rows)

	for i := m.cursorRow + 1; i < len(sheet.Rows); i++ {
		for j := range sheet.Rows[i] {
			sheet.Rows[i][j].Row = i
		}
	}

	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Row inserted at %d", m.cursorRow+1),
		Type:    models.StatusSuccess,
	}
}

// insertColumn inserts a new column at current position
func (m *Model) insertColumn() {
	sheet := &m.sheets[m.currentSheet]
	for i := range sheet.Rows {
		newCell := models.Cell{Row: i, Col: m.cursorCol}
		sheet.Rows[i] = append(sheet.Rows[i][:m.cursorCol], append([]models.Cell{newCell}, sheet.Rows[i][m.cursorCol:]...)...)

		for j := m.cursorCol + 1; j < len(sheet.Rows[i]); j++ {
			sheet.Rows[i][j].Col = j
		}
	}
	sheet.MaxCols++
	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Column inserted at %s", ui.ColIndexToLetter(m.cursorCol)),
		Type:    models.StatusSuccess,
	}
}

// pasteCell pastes clipboard content to current cell
func (m *Model) pasteCell() {
	content, err := clipboard.ReadAll()
	if err != nil {
		m.status = models.StatusMsg{Message: "Failed to paste", Type: models.StatusError}
		return
	}

	sheet := &m.sheets[m.currentSheet]
	if m.cursorRow >= len(sheet.Rows) {
		for i := len(sheet.Rows); i <= m.cursorRow; i++ {
			sheet.Rows = append(sheet.Rows, make([]models.Cell, sheet.MaxCols))
		}
	}
	if m.cursorCol >= len(sheet.Rows[m.cursorRow]) {
		oldLen := len(sheet.Rows[m.cursorRow])
		newRow := make([]models.Cell, sheet.MaxCols)
		copy(newRow, sheet.Rows[m.cursorRow])
		for i := oldLen; i < len(newRow); i++ {
			newRow[i] = models.Cell{Row: m.cursorRow, Col: i}
		}
		sheet.Rows[m.cursorRow] = newRow
	}

	lines := strings.Split(content, "\n")
	for rowOffset, line := range lines {
		targetRow := m.cursorRow + rowOffset
		if targetRow >= sheet.MaxRows {
			break
		}

		cells := strings.Split(line, "\t")
		for colOffset, cellValue := range cells {
			targetCol := m.cursorCol + colOffset
			if targetCol >= sheet.MaxCols {
				break
			}

			if targetRow >= len(sheet.Rows) {
				continue
			}
			if targetCol >= len(sheet.Rows[targetRow]) {
				continue
			}

			cell := &sheet.Rows[targetRow][targetCol]
			cellValue = strings.TrimSpace(cellValue)

			if strings.HasPrefix(cellValue, "=") {
				cell.Formula = cellValue[1:]
				cell.Value = m.evaluateFormula(cell.Formula)
			} else {
				cell.Value = cellValue
				cell.Formula = ""
			}
		}
	}

	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{Message: "Pasted", Type: models.StatusSuccess}
}

// saveFile saves the current workbook
func (m *Model) saveFile() {
	if m.fileFormat == "csv" {
		err := loader.SaveCSV(m.sheets[m.currentSheet], m.filename)
		if err != nil {
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Save failed: %v", err),
				Type:    models.StatusError,
			}
			return
		}
	} else {
		err := loader.SaveExcel(m.sheets, m.filename)
		if err != nil {
			m.status = models.StatusMsg{
				Message: fmt.Sprintf("Save failed: %v", err),
				Type:    models.StatusError,
			}
			return
		}
	}

	m.modified = false
	m.quitConfirm = false
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("✓ Saved to %s", m.filename),
		Type:    models.StatusSuccess,
	}
}

// updateSaveAs handles save as mode
func (m Model) updateSaveAs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEscape:
		m.mode = models.ModeNormal
		m.saveAsInput.Blur()
		return m, nil

	case tea.KeyEnter:
		filename := strings.TrimSpace(m.saveAsInput.Value())
		if filename != "" {
			oldFilename := m.filename
			oldFormat := m.fileFormat
			m.filename = filename

			ext := strings.ToLower(strings.TrimPrefix(filename[strings.LastIndex(filename, "."):], "."))
			if ext == "csv" {
				m.fileFormat = "csv"
			} else {
				m.fileFormat = "xlsx"
			}

			m.saveFile()
			if m.status.Type == models.StatusError {
				m.filename = oldFilename
				m.fileFormat = oldFormat
			}
		}
		m.mode = models.ModeNormal
		m.saveAsInput.Blur()
		return m, nil
	}

	m.saveAsInput, cmd = m.saveAsInput.Update(msg)
	return m, cmd
}

// fillDown fills current cell value down to selected range
func (m *Model) fillDown() {
	if !m.isSelecting {
		m.status = models.StatusMsg{Message: "Select range first (V)", Type: models.StatusWarning}
		return
	}

	sheet := &m.sheets[m.currentSheet]
	startRow := m.selectStart[0]
	endRow := m.selectEnd[0]
	col := m.selectStart[1]

	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}

	if startRow >= len(sheet.Rows) || col >= len(sheet.Rows[startRow]) {
		return
	}

	sourceCell := sheet.Rows[startRow][col]
	for row := startRow + 1; row <= endRow && row < len(sheet.Rows); row++ {
		if col < len(sheet.Rows[row]) {
			cell := &sheet.Rows[row][col]
			cell.Value = sourceCell.Value
			cell.Formula = sourceCell.Formula
		}
	}

	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Filled %d cells", endRow-startRow),
		Type:    models.StatusSuccess,
	}
}

// fillRight fills current cell value right to selected range
func (m *Model) fillRight() {
	if !m.isSelecting {
		m.status = models.StatusMsg{Message: "Select range first (V)", Type: models.StatusWarning}
		return
	}

	sheet := &m.sheets[m.currentSheet]
	row := m.selectStart[0]
	startCol := m.selectStart[1]
	endCol := m.selectEnd[1]

	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	if row >= len(sheet.Rows) || startCol >= len(sheet.Rows[row]) {
		return
	}

	sourceCell := sheet.Rows[row][startCol]
	for col := startCol + 1; col <= endCol && col < len(sheet.Rows[row]); col++ {
		cell := &sheet.Rows[row][col]
		cell.Value = sourceCell.Value
		cell.Formula = sourceCell.Formula
	}

	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Filled %d cells", endCol-startCol),
		Type:    models.StatusSuccess,
	}
}

// applyFormulaToRange applies current cell's formula to entire selected range
func (m *Model) applyFormulaToRange() {
	if !m.isSelecting {
		m.status = models.StatusMsg{Message: "Select range first (V)", Type: models.StatusWarning}
		return
	}

	sheet := &m.sheets[m.currentSheet]

	// Get source cell (current cursor position)
	if m.cursorRow >= len(sheet.Rows) || m.cursorCol >= len(sheet.Rows[m.cursorRow]) {
		m.status = models.StatusMsg{Message: "No formula in current cell", Type: models.StatusError}
		return
	}

	sourceCell := sheet.Rows[m.cursorRow][m.cursorCol]
	if sourceCell.Formula == "" {
		m.status = models.StatusMsg{Message: "Current cell has no formula", Type: models.StatusWarning}
		return
	}

	// Normalize selection
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

	// Calculate offset from source cell
	sourceRow := m.cursorRow
	sourceCol := m.cursorCol

	cellsUpdated := 0

	// Apply formula with relative references
	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		for col := startCol; col <= endCol && col < len(sheet.Rows[row]); col++ {
			if row == sourceRow && col == sourceCol {
				continue // Skip source cell
			}

			cell := &sheet.Rows[row][col]

			// Adjust formula relative to new position
			rowOffset := row - sourceRow
			colOffset := col - sourceCol

			adjustedFormula := m.adjustFormulaReferences(sourceCell.Formula, rowOffset, colOffset)

			cell.Formula = adjustedFormula
			cell.Value = m.evaluateFormula(adjustedFormula)
			cellsUpdated++
		}
	}

	m.modified = true
	m.recalculateFormulas()
	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Applied formula to %d cells", cellsUpdated),
		Type:    models.StatusSuccess,
	}
}

// adjustFormulaReferences adjusts cell references in formula by given offset
func (m *Model) adjustFormulaReferences(formula string, rowOffset, colOffset int) string {
	// This is a simplified version - it handles basic cell references
	// For production, you'd want a proper formula parser

	result := strings.ToUpper(formula)
	var adjusted strings.Builder

	i := 0
	for i < len(result) {
		if result[i] >= 'A' && result[i] <= 'Z' {
			// Found potential cell reference
			colStart := i
			col := 0

			// Parse column letters
			for i < len(result) && result[i] >= 'A' && result[i] <= 'Z' {
				col = col*26 + int(result[i]-'A') + 1
				i++
			}
			col--

			// Check if followed by number (row)
			if i < len(result) && result[i] >= '0' && result[i] <= '9' {
				// rowStart := i
				row := 0

				// Parse row number
				for i < len(result) && result[i] >= '0' && result[i] <= '9' {
					row = row*10 + int(result[i]-'0')
					i++
				}
				row--

				// Adjust the reference
				newRow := row + rowOffset
				newCol := col + colOffset

				if newRow >= 0 && newCol >= 0 {
					adjusted.WriteString(ui.ColIndexToLetter(newCol))
					adjusted.WriteString(strconv.Itoa(newRow + 1))
				} else {
					// Invalid reference, keep original
					adjusted.WriteString(result[colStart:i])
				}
			} else {
				// Not a cell reference, keep original
				adjusted.WriteString(result[colStart:i])
			}
		} else {
			adjusted.WriteByte(result[i])
			i++
		}
	}

	return adjusted.String()
}

// renderEditMode renders the edit mode overlay
func (m Model) renderEditMode() string {
	base := m.renderNormal()

	t := theme.GetCurrentTheme()
	editInfo := lipgloss.NewStyle().
		Background(t.Accent).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 2).
		Bold(true).
		Render("EDIT MODE")

	cellRef := ui.ColIndexToLetter(m.cursorCol) + fmt.Sprintf("%d", m.cursorRow+1)
	cellInfo := lipgloss.NewStyle().
		Background(t.Border).
		Foreground(t.Text).
		Padding(0, 2).
		Render(cellRef + ": " + m.editInput.View())

	modifiedIndicator := ""
	if m.modified {
		modifiedIndicator = lipgloss.NewStyle().
			Foreground(t.Warning).
			Bold(true).
			Render(" [Modified]")
	}

	info := lipgloss.JoinHorizontal(lipgloss.Top, editInfo, cellInfo, modifiedIndicator)

	return base + "\n" + info
}

// renderSaveAs renders the save as modal
func (m Model) renderSaveAs() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("💾 Save As") + "\n\n"
	content += m.styles.ModalKey.Render("Filename:") + "\n"
	content += m.saveAsInput.View() + "\n\n"
	content += lipgloss.NewStyle().
		Foreground(t.DimText).
		Render("Supported formats: .xlsx, .csv")

	return m.styles.Modal.Width(50).Render(content)
}
