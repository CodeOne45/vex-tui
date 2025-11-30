package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/odesaur/vex-tui/internal/theme"
	"github.com/odesaur/vex-tui/internal/ui"
	"github.com/odesaur/vex-tui/pkg/models"
)

// View renders the current state
func (m Model) View() string {
	// Wait for terminal size
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if len(m.sheets) == 0 {
		return m.renderEmpty()
	}

	switch m.mode {
	case models.ModeDetail:
		return ui.RenderModal(m.width, m.height, m.renderDetail())
	case models.ModeJump:
		return ui.RenderModal(m.width, m.height, m.renderJump())
	case models.ModeExport:
		return ui.RenderModal(m.width, m.height, m.renderExport())
	case models.ModeTheme:
		return ui.RenderModal(m.width, m.height, m.renderThemeSelector())
	case models.ModeChart:
		return ui.RenderModal(m.width, m.height, m.renderChart())
	case models.ModeSelectRange:
		return m.renderSelectRange()
	default:
		return m.renderNormal()
	}
}

// renderEmpty renders the empty state
func (m Model) renderEmpty() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("Excel TUI v2.0"))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(theme.GetCurrentTheme().DimText).
		Render("No data to display"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.Help.Render(m.help.View(m.keys)))
	return b.String()
}

// renderNormal renders the normal viewing mode
func (m Model) renderNormal() string {
	sheet := m.sheets[m.currentSheet]
	var b strings.Builder

	// Title bar
	title := fmt.Sprintf("File: %s", m.filename)
	if len(m.sheets) > 1 {
		title += fmt.Sprintf(" | %s (%d/%d)", sheet.Name, m.currentSheet+1, len(m.sheets))
	} else {
		title += fmt.Sprintf(" | %s", sheet.Name)
	}
	b.WriteString(m.styles.Title.Render(title))
	b.WriteString("\n")

	// Formula bar
	b.WriteString(m.renderFormulaBar())
	b.WriteString("\n\n")

	// Render table
	b.WriteString(m.renderTable())

	// Status bar
	b.WriteString("\n")
	b.WriteString(m.renderStatusBar())

	// Search bar (vim-style at bottom)
	if m.mode == models.ModeSearch || m.searchQuery != "" {
		b.WriteString("\n")
		b.WriteString(m.renderSearchBar())
	}

	// Help (short by default, full when toggled with '?')
	b.WriteString("\n")
	if m.help.ShowAll {
		b.WriteString(m.styles.Help.Render(m.help.View(m.keys)))
	} else {
		b.WriteString(m.styles.Help.Render(m.help.ShortHelpView(m.keys.ShortHelp())))
	}

	return b.String()
}

// renderFormulaBar renders the formula bar showing current cell info
func (m Model) renderFormulaBar() string {
	sheet := m.sheets[m.currentSheet]
	if m.cursorRow < len(sheet.Rows) && m.cursorCol < len(sheet.Rows[m.cursorRow]) {
		cell := sheet.Rows[m.cursorRow][m.cursorCol]
		cellRef := ui.ColIndexToLetter(m.cursorCol) + fmt.Sprintf("%d", m.cursorRow+1)

		t := theme.GetCurrentTheme()
		formulaText := lipgloss.NewStyle().
			Foreground(t.Secondary).
			Bold(true).
			Render(cellRef)

		if cell.Formula != "" {
			formulaText += lipgloss.NewStyle().
				Foreground(t.Text).
				Render(" = " + ui.Truncate(cell.Formula, 100))
		} else {
			formulaText += lipgloss.NewStyle().
				Foreground(t.DimText).
				Render(" " + ui.Truncate(cell.Value, 100))
		}
		return m.styles.FormulaBar.Render(formulaText)
	}
	return m.styles.FormulaBar.Render(" ")
}

// renderTable renders the spreadsheet table
func (m Model) renderTable() string {
	sheet := m.sheets[m.currentSheet]
	visibleRows := ui.Max(1, m.height-9)
	visibleCols := ui.Max(1, (m.width-8)/(ui.MinCellWidth+2))

	var b strings.Builder
	sep := m.styles.Separator.Render("│")
	rowNumberWidth := 5

	availableWidth := m.width - rowNumberWidth - (visibleCols + 1) // account for separators
	colWidth := ui.Max(ui.MinCellWidth, availableWidth/ui.Max(1, visibleCols))
	if colWidth > ui.MaxCellWidth {
		colWidth = ui.MaxCellWidth
	}

	header := m.styles.Header.Width(colWidth)
	headerHighlight := m.styles.HeaderHighlight.Width(colWidth)
	cellStyle := m.styles.Cell.Width(colWidth)
	selectedCell := m.styles.SelectedCell.Width(colWidth)
	rowHighlight := m.styles.RowHighlight.Width(colWidth)
	colHighlight := m.styles.ColHighlight.Width(colWidth)
	searchMatch := m.styles.SearchMatch.Width(colWidth)
	selectionStyle := lipgloss.NewStyle().
		Foreground(theme.GetCurrentTheme().Text).
		Background(theme.GetCurrentTheme().Secondary).
		Width(colWidth)

	// Column headers
	b.WriteString(m.styles.RowNum.Render(""))
	b.WriteString(sep)

	for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
		colLetter := ui.ColIndexToLetter(col)
		if col == m.cursorCol {
			b.WriteString(headerHighlight.Render(ui.PadCenter(colLetter, colWidth)))
		} else {
			b.WriteString(header.Render(ui.PadCenter(colLetter, colWidth)))
		}
		b.WriteString(sep)
	}
	b.WriteString("\n")

	// Data rows
	endRow := ui.Min(m.offsetRow+visibleRows, sheet.MaxRows)
	for row := m.offsetRow; row < endRow; row++ {
		// Row number
		if row == m.cursorRow {
			b.WriteString(m.styles.SelectedRowNum.Render(fmt.Sprintf("%d", row+1)))
		} else {
			b.WriteString(m.styles.RowNum.Render(fmt.Sprintf("%d", row+1)))
		}
		b.WriteString(sep)

		// Cells
		if row < len(sheet.Rows) {
			for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
				cellText := ""

				if col < len(sheet.Rows[row]) {
					cell := sheet.Rows[row][col]
					if m.showFormulas && cell.Formula != "" {
						cellText = "=" + cell.Formula
					} else {
						cellText = cell.Value
					}
				}

				cellText = ui.TruncateToWidth(cellText, colWidth)

				// Determine style
				var style lipgloss.Style
				if row == m.cursorRow && col == m.cursorCol {
					style = selectedCell
				} else if m.isSelecting && m.isInSelection(row, col) {
					// Highlight selection with different color
					style = selectionStyle
				} else if m.isSearchMatch(row, col) {
					style = searchMatch
				} else if row == m.cursorRow {
					style = rowHighlight
				} else if col == m.cursorCol {
					style = colHighlight
				} else {
					style = cellStyle
				}

				b.WriteString(style.Render(cellText))
				b.WriteString(sep)
			}
		} else {
			// Empty row
			for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
				var style lipgloss.Style
				if row == m.cursorRow && col == m.cursorCol {
					style = selectedCell
				} else if row == m.cursorRow {
					style = rowHighlight
				} else if col == m.cursorCol {
					style = colHighlight
				} else {
					style = cellStyle
				}
				b.WriteString(style.Render(strings.Repeat(" ", colWidth)))
				b.WriteString(sep)
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderStatusBar renders the status bar at the bottom
func (m Model) renderStatusBar() string {
	sheet := m.sheets[m.currentSheet]
	t := theme.GetCurrentTheme()

	parts := []string{
		lipgloss.NewStyle().Foreground(t.Secondary).Bold(true).Render("Rows:") +
			lipgloss.NewStyle().Foreground(t.Text).Render(fmt.Sprintf(" %d", sheet.MaxRows)),
		lipgloss.NewStyle().Foreground(t.Secondary).Bold(true).Render("Cols:") +
			lipgloss.NewStyle().Foreground(t.Text).Render(fmt.Sprintf(" %d", sheet.MaxCols)),
		lipgloss.NewStyle().Foreground(t.Secondary).Bold(true).Render("Pos:") +
			lipgloss.NewStyle().Foreground(t.Text).Render(fmt.Sprintf(" %s", ui.ColIndexToLetter(m.cursorCol)+fmt.Sprintf("%d", m.cursorRow+1))),
	}

	if m.showFormulas {
		parts = append(parts, lipgloss.NewStyle().Foreground(t.Accent).Render("Formulas"))
	}

	if len(m.searchResults) > 0 {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(t.SearchMatch).
			Bold(true).
			Render(fmt.Sprintf("Search %d/%d", m.searchIndex+1, len(m.searchResults))))
	}

	if m.status.Message != "" {
		statusColor := ui.GetStatusColor(m.status.Type)
		parts = append(parts, lipgloss.NewStyle().
			Foreground(statusColor).
			Render(m.status.Message))
	}

	return m.styles.StatusBar.Render(strings.Join(parts, " │ "))
}

// renderSearchBar renders the search bar
func (m Model) renderSearchBar() string {
	t := theme.GetCurrentTheme()

	if m.mode == models.ModeSearch {
		prompt := m.styles.SearchPrompt.Render("/")
		input := m.searchInput.View()
		return m.styles.SearchBar.Render(prompt + input)
	} else if m.searchQuery != "" {
		searchInfo := m.styles.SearchPrompt.Render("/") +
			lipgloss.NewStyle().Foreground(t.Text).Render(m.searchQuery)
		if len(m.searchResults) > 0 {
			searchInfo += lipgloss.NewStyle().
				Foreground(t.DimText).
				Render(fmt.Sprintf(" (%d results)", len(m.searchResults)))
		}
		return m.styles.SearchBar.Render(searchInfo)
	}
	return ""
}

// renderDetail renders the cell detail modal
func (m Model) renderDetail() string {
	sheet := m.sheets[m.currentSheet]
	if m.cursorRow >= len(sheet.Rows) || m.cursorCol >= len(sheet.Rows[m.cursorRow]) {
		return m.styles.Modal.Render(m.styles.ModalTitle.Render("Cell Details") + "\n\nNo data")
	}

	cell := sheet.Rows[m.cursorRow][m.cursorCol]
	cellRef := ui.ColIndexToLetter(m.cursorCol) + fmt.Sprintf("%d", m.cursorRow+1)
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("Cell Details") + "\n\n"
	content += m.styles.ModalKey.Render("Cell: ") + m.styles.ModalValue.Render(cellRef) + "\n\n"
	content += m.styles.ModalKey.Render("Value:\n") + m.styles.ModalValue.Render(ui.WrapText(cell.Value, 56)) + "\n\n"

	if cell.Formula != "" {
		content += m.styles.ModalKey.Render("Formula:\n") + m.styles.ModalValue.Render("="+ui.WrapText(cell.Formula, 55)) + "\n\n"
	}

	content += m.styles.ModalKey.Render("Type: ") + m.styles.ModalValue.Render(ui.GetCellType(cell)) + "\n"
	content += lipgloss.NewStyle().
		Foreground(t.DimText).
		Italic(true).
		Render("\nPress Enter or Esc to close")

	return m.styles.Modal.Render(content)
}

// renderJump renders the jump to cell modal
func (m Model) renderJump() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("Jump to Cell") + "\n\n"
	content += m.styles.ModalKey.Render("Enter cell reference:") + "\n"
	content += m.jumpInput.View() + "\n\n"
	content += lipgloss.NewStyle().Foreground(t.DimText).Render("Formats:\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  - A100   (column + row)\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  - 500    (row only)\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  - 10,5   (row,col)")

	return m.styles.Modal.Width(50).Render(content)
}

// renderExport renders the export modal
func (m Model) renderExport() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("Export Sheet") + "\n\n"
	content += m.styles.ModalKey.Render("Filename:") + "\n"
	content += m.exportInput.View() + "\n\n"
	content += lipgloss.NewStyle().
		Foreground(t.DimText).
		Render("Supported formats: .csv, .json")

	return m.styles.Modal.Width(50).Render(content)
}

// renderThemeSelector renders the theme selection modal
func (m Model) renderThemeSelector() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("Select Theme") + "\n\n"

	themes := []struct {
		num  string
		name string
		desc string
	}{
		{"1", "Catppuccin Mocha", "Soft pastels, gentle on the eyes"},
		{"2", "Nord", "Cool Arctic blues, minimal"},
		{"3", "Rosé Pine", "Elegant rose tones"},
		{"4", "Tokyo Night", "Vibrant cyberpunk vibes"},
		{"5", "Gruvbox", "Warm retro colors"},
		{"6", "Dracula", "Classic high contrast"},
	}

	for _, theme := range themes {
		numStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(t.Text).Bold(true)
		descStyle := lipgloss.NewStyle().Foreground(t.DimText)

		current := ""
		if strings.Contains(strings.ToLower(theme.name), strings.ToLower(m.themeName)) ||
			strings.Contains(strings.ToLower(m.themeName), strings.ToLower(strings.ReplaceAll(theme.name, " ", "-"))) {
			current = lipgloss.NewStyle().Foreground(t.Accent).Render(" *")
		}

		content += numStyle.Render(theme.num) + "  " + nameStyle.Render(theme.name) + current + "\n"
		content += "   " + descStyle.Render(theme.desc) + "\n\n"
	}

	content += lipgloss.NewStyle().
		Foreground(t.DimText).
		Italic(true).
		Render("\nPress 1-6 to select, Esc to cancel")

	return m.styles.Modal.Width(60).Render(content)
}

// renderChart renders the chart visualization modal
func (m Model) renderChart() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("Data Visualization") + "\n\n"

	// Show chart type selector
	types := []string{"1. Bar Chart", "2. Line Chart", "3. Sparkline", "4. Pie Chart"}
	for i, typ := range types {
		if i == m.chartType {
			content += lipgloss.NewStyle().Foreground(t.Accent).Bold(true).Render("-> "+typ) + "\n"
		} else {
			content += lipgloss.NewStyle().Foreground(t.Text).Render("  "+typ) + "\n"
		}
	}

	content += "\n" + strings.Repeat("─", 60) + "\n\n"

	// Extract data from selection
	startRow := m.selectStart[0]
	startCol := m.selectStart[1]
	endRow := m.selectEnd[0]
	endCol := m.selectEnd[1]

	// Normalize selection
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	// Import chart package
	chartData := extractChartData(m.sheets[m.currentSheet], startRow, startCol, endRow, endCol)

	// Render chart
	modalStyle := lipgloss.NewStyle()

	switch m.chartType {
	case 0: // Bar
		content += renderBarChart(chartData, modalStyle, t.Accent, t.Text)
	case 1: // Line
		content += renderLineChart(chartData, modalStyle, t.Accent, t.Text)
	case 2: // Sparkline
		content += "Sparkline: " + renderSparkline(chartData, t.Accent) + "\n"
		content += renderBarChart(chartData, modalStyle, t.Accent, t.Text)
	case 3: // Pie
		colors := []lipgloss.Color{t.Accent, t.Primary, t.Secondary, t.Success, t.Warning}
		content += renderPieChart(chartData, modalStyle, colors, t.Text)
	}

	content += "\n" + lipgloss.NewStyle().
		Foreground(t.DimText).
		Italic(true).
		Render("Press 1-4 to switch chart type, Esc to close")

	return m.styles.Modal.Width(70).Height(30).Render(content)
}

// renderSelectRange renders the selection mode overlay
func (m Model) renderSelectRange() string {
	base := m.renderNormal()

	// Add selection info overlay
	t := theme.GetCurrentTheme()
	info := lipgloss.NewStyle().
		Background(t.Border).
		Foreground(t.Accent).
		Padding(0, 2).
		Bold(true).
		Render(fmt.Sprintf("SELECTION MODE: %d×%d | Move with arrows | V to finish | Esc to cancel",
			abs(m.selectEnd[0]-m.selectStart[0])+1,
			abs(m.selectEnd[1]-m.selectStart[1])+1))

	return base + "\n" + info
}

// isInSelection checks if a cell is in the current selection
func (m Model) isInSelection(row, col int) bool {
	if !m.isSelecting {
		return false
	}

	startRow := m.selectStart[0]
	startCol := m.selectStart[1]
	endRow := m.selectEnd[0]
	endCol := m.selectEnd[1]

	// Normalize
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	return row >= startRow && row <= endRow && col >= startCol && col <= endCol
}

// Chart rendering helpers using internal/chart package
func extractChartData(sheet models.Sheet, startRow, startCol, endRow, endCol int) chartData {
	data := chartData{
		Labels: make([]string, 0),
		Values: make([]float64, 0),
	}

	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		if startCol < len(sheet.Rows[row]) {
			label := sheet.Rows[row][startCol].Value
			if label == "" {
				label = fmt.Sprintf("Row %d", row+1)
			}
			data.Labels = append(data.Labels, label)

			if startCol+1 <= endCol && startCol+1 < len(sheet.Rows[row]) {
				valStr := sheet.Rows[row][startCol+1].Value
				if val, err := strconv.ParseFloat(valStr, 64); err == nil {
					data.Values = append(data.Values, val)
				} else {
					data.Values = append(data.Values, 0)
				}
			}
		}
	}

	return data
}

type chartData struct {
	Labels []string
	Values []float64
}

func renderBarChart(data chartData, style lipgloss.Style, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return "No data"
	}

	var b strings.Builder
	maxVal := maxFloat(data.Values)
	if maxVal == 0 {
		maxVal = 1
	}

	maxLabelLen := 15
	barWidth := 40

	for i, val := range data.Values {
		if i >= len(data.Labels) {
			break
		}

		label := data.Labels[i]
		if len(label) > maxLabelLen {
			label = label[:maxLabelLen-2] + ".."
		}
		labelStr := lipgloss.NewStyle().Foreground(textColor).Width(maxLabelLen).Render(label)

		barLen := int(float64(barWidth) * (val / maxVal))
		if barLen < 0 {
			barLen = 0
		}
		bar := strings.Repeat("█", barLen)
		barStr := lipgloss.NewStyle().Foreground(accentColor).Render(bar)
		valStr := lipgloss.NewStyle().Foreground(textColor).Render(fmt.Sprintf(" %.1f", val))

		b.WriteString(labelStr + " │ " + barStr + valStr + "\n")
	}

	b.WriteString(strings.Repeat(" ", maxLabelLen) + " └" + strings.Repeat("─", barWidth+2) + "\n")
	b.WriteString(strings.Repeat(" ", maxLabelLen) + "  0" + strings.Repeat(" ", barWidth-10) + fmt.Sprintf("%.1f", maxVal))

	return b.String()
}

func renderLineChart(data chartData, style lipgloss.Style, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return "No data"
	}

	height, width := 12, 50
	var b strings.Builder

	minVal := minFloat(data.Values)
	maxVal := maxFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for i, val := range data.Values {
		x := int(float64(i) * float64(width-1) / float64(len(data.Values)-1))
		if x >= width {
			x = width - 1
		}
		normalized := (val - minVal) / (maxVal - minVal)
		y := height - 1 - int(normalized*float64(height-1))
		if y < 0 {
			y = 0
		}
		if y >= height {
			y = height - 1
		}
		grid[y][x] = '●'
	}

	for i, row := range grid {
		for _, char := range row {
			if char == '●' {
				b.WriteString(lipgloss.NewStyle().Foreground(accentColor).Render(string(char)))
			} else {
				b.WriteString(" ")
			}
		}
		if i%3 == 0 {
			val := maxVal - (float64(i)/float64(height-1))*(maxVal-minVal)
			b.WriteString(lipgloss.NewStyle().Foreground(textColor).Render(fmt.Sprintf(" %.1f", val)))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderSparkline(data chartData, accentColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return ""
	}

	chars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	minVal := minFloat(data.Values)
	maxVal := maxFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	var b strings.Builder
	for _, val := range data.Values {
		normalized := (val - minVal) / (maxVal - minVal)
		idx := int(normalized * float64(len(chars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(chars) {
			idx = len(chars) - 1
		}
		b.WriteRune(chars[idx])
	}

	return lipgloss.NewStyle().Foreground(accentColor).Render(b.String())
}

func renderPieChart(data chartData, style lipgloss.Style, colors []lipgloss.Color, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return "No data"
	}

	var b strings.Builder
	total := sumFloat(data.Values)
	if total == 0 {
		total = 1
	}

	radius := 8
	centerX, centerY := radius, radius

	grid := make([][]int, radius*2+1)
	for i := range grid {
		grid[i] = make([]int, radius*2+1)
		for j := range grid[i] {
			grid[i][j] = -1
		}
	}

	currentAngle := 0.0
	for i, val := range data.Values {
		percentage := val / total
		angle := percentage * 360

		for y := 0; y <= radius*2; y++ {
			for x := 0; x <= radius*2; x++ {
				dx := x - centerX
				dy := y - centerY
				dist := math.Sqrt(float64(dx*dx + dy*dy))

				if dist <= float64(radius) {
					pointAngle := math.Atan2(float64(dy), float64(dx)) * 180 / math.Pi
					if pointAngle < 0 {
						pointAngle += 360
					}
					if pointAngle >= currentAngle && pointAngle < currentAngle+angle {
						grid[y][x] = i
					}
				}
			}
		}
		currentAngle += angle
	}

	for _, row := range grid {
		for _, val := range row {
			if val >= 0 && val < len(colors) {
				b.WriteString(lipgloss.NewStyle().Foreground(colors[val%len(colors)]).Render("●"))
			} else {
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	for i, label := range data.Labels {
		if i >= len(data.Values) || i >= len(colors) {
			break
		}
		percentage := (data.Values[i] / total) * 100
		colorBox := lipgloss.NewStyle().Foreground(colors[i%len(colors)]).Render("■")
		labelStr := lipgloss.NewStyle().Foreground(textColor).Render(fmt.Sprintf(" %s: %.1f%%", label, percentage))
		b.WriteString(colorBox + labelStr + "\n")
	}

	return b.String()
}

func maxFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}

func minFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}

func sumFloat(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum
}
