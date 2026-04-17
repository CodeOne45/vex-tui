package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/CodeOne45/vex-tui/internal/theme"
	"github.com/CodeOne45/vex-tui/internal/ui"
	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
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
	case models.ModeEdit:
		return m.renderEditMode()
	case models.ModeSaveAs:
		return ui.RenderModal(m.width, m.height, m.renderSaveAs())
	case models.ModeFilter:
		return m.renderFilterMode()
	case models.ModeDataProfile:
		return ui.RenderModal(m.width, m.height, m.renderDataProfile())
	default:
		return m.renderNormal()
	}
}

// renderEmpty renders the empty state
func (m Model) renderEmpty() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("📊 Excel TUI v2.0"))
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

	// Title bar — show modified indicator and freeze/sort state
	modMark := ""
	if m.modified {
		modMark = lipgloss.NewStyle().Foreground(theme.GetCurrentTheme().Warning).Bold(true).Render(" ●")
	}
	sheetInfo := sheet.Name
	if len(m.sheets) > 1 {
		sheetInfo = fmt.Sprintf("%s (%d/%d)", sheet.Name, m.currentSheet+1, len(m.sheets))
	}
	extras := ""
	if m.freezeHeader {
		extras += lipgloss.NewStyle().Foreground(theme.GetCurrentTheme().Accent).Render(" ⌶")
	}
	if m.filterActive {
		extras += lipgloss.NewStyle().Foreground(theme.GetCurrentTheme().Warning).
			Render(fmt.Sprintf(" ▾ %s", m.filterQuery))
	}
	if m.sortActive {
		dir := "↑"
		if !m.sortAsc {
			dir = "↓"
		}
		extras += lipgloss.NewStyle().Foreground(theme.GetCurrentTheme().Secondary).
			Render(fmt.Sprintf(" %s%s", ui.ColIndexToLetter(m.sortCol), dir))
	}
	title := fmt.Sprintf("  %s  •  %s%s", m.filename, sheetInfo, extras)
	b.WriteString(m.styles.Title.Render(title) + modMark)
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

	// Help
	b.WriteString("\n")
	b.WriteString(m.styles.Help.Render(m.help.View(m.keys)))

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
	t := theme.GetCurrentTheme()

	renderColHeaders := func() {
		b.WriteString(m.styles.RowNum.Render(""))
		b.WriteString(sep)
		for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
			width := ui.MinCellWidth
			if w, ok := sheet.ColWidths[col]; ok && w > 0 {
				width = w
			}
			label := ui.ColIndexToLetter(col)
			if m.sortActive && col == m.sortCol {
				if m.sortAsc {
					label += "↑"
				} else {
					label += "↓"
				}
			}
			headerStyle := m.styles.Header
			if col == m.cursorCol {
				headerStyle = m.styles.HeaderHighlight
			}
			b.WriteString(headerStyle.Width(width).Render(ui.PadCenter(label, width)))
			b.WriteString(sep)
		}
		b.WriteString("\n")
	}

	renderRow := func(row int, isHeaderPin bool) {
		if row == m.cursorRow {
			b.WriteString(m.styles.SelectedRowNum.Render(fmt.Sprintf("%d", row+1)))
		} else if isHeaderPin {
			b.WriteString(lipgloss.NewStyle().
				Foreground(t.Primary).Bold(true).Align(lipgloss.Right).Width(5).
				Render(fmt.Sprintf("%d", row+1)))
		} else {
			b.WriteString(m.styles.RowNum.Render(fmt.Sprintf("%d", row+1)))
		}
		b.WriteString(sep)

		if row < len(sheet.Rows) {
			for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
				width := ui.MinCellWidth
				if w, ok := sheet.ColWidths[col]; ok && w > 0 {
					width = w
				}

				rawVal := ""
				isFormula := false
				isEmpty := true
				if col < len(sheet.Rows[row]) {
					cell := sheet.Rows[row][col]
					isEmpty = cell.Value == "" && cell.Formula == ""
					if m.showFormulas && cell.Formula != "" {
						rawVal = "=" + cell.Formula
						isFormula = true
					} else {
						rawVal = cell.Value
						isFormula = cell.Formula != ""
					}
				}

				cellText := ui.TruncateToWidth(rawVal, width)

				var style lipgloss.Style
				switch {
				case row == m.cursorRow && col == m.cursorCol:
					style = m.styles.SelectedCell
				case m.isSelecting && m.isInSelection(row, col):
					style = m.styles.SelectionCell
				case m.isSearchMatch(row, col):
					style = m.styles.SearchMatch
				case isHeaderPin:
					style = m.styles.FrozenHeaderCell
				case row == m.cursorRow:
					style = m.styles.RowHighlight
				case col == m.cursorCol:
					// apply cell-type coloring within col highlight
					if isFormula {
						style = m.styles.ColHighlight.Foreground(t.Secondary)
					} else if !isEmpty {
						if _, err := strconv.ParseFloat(rawVal, 64); err == nil {
							style = m.styles.ColHighlight.Foreground(t.Accent)
						} else {
							style = m.styles.ColHighlight
						}
					} else {
						style = m.styles.ColHighlight
					}
				default:
					// Cell type coloring in normal cells
					if isFormula {
						style = m.styles.FormulaCell
					} else if isEmpty {
						style = m.styles.Cell
					} else if _, err := strconv.ParseFloat(rawVal, 64); err == nil {
						style = m.styles.NumberCell
					} else {
						style = m.styles.Cell
					}
				}

				b.WriteString(style.Width(width).Render(cellText))
				b.WriteString(sep)
			}
		} else {
			for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
				width := ui.MinCellWidth
				if w, ok := sheet.ColWidths[col]; ok && w > 0 {
					width = w
				}
				var style lipgloss.Style
				if row == m.cursorRow && col == m.cursorCol {
					style = m.styles.SelectedCell
				} else if row == m.cursorRow {
					style = m.styles.RowHighlight
				} else if col == m.cursorCol {
					style = m.styles.ColHighlight
				} else {
					style = m.styles.Cell
				}
				b.WriteString(style.Width(width).Render(strings.Repeat(" ", width)))
				b.WriteString(sep)
			}
		}
		b.WriteString("\n")
	}

	renderColHeaders()

	// Frozen header row — always pinned at top
	if m.freezeHeader && sheet.MaxRows > 0 {
		renderRow(0, true)
		b.WriteString(m.styles.FrozenDivider.Render(strings.Repeat("─", m.width-2)) + "\n")
	}

	// Determine scrollable viewport start row
	startRow := m.offsetRow
	if m.freezeHeader && startRow == 0 && sheet.MaxRows > 1 {
		startRow = 1
	}

	endRow := ui.Min(startRow+visibleRows, sheet.MaxRows)
	if m.freezeHeader {
		endRow = ui.Min(startRow+visibleRows-2, sheet.MaxRows)
	}
	for row := startRow; row < endRow; row++ {
		if m.freezeHeader && row == 0 {
			continue
		}
		renderRow(row, false)
	}

	return b.String()
}

// renderStatusBar renders the status bar at the bottom
func (m Model) renderStatusBar() string {
	sheet := m.sheets[m.currentSheet]
	t := theme.GetCurrentTheme()

	dim := lipgloss.NewStyle().Foreground(t.DimText)
	bold := lipgloss.NewStyle().Foreground(t.Secondary).Bold(true)
	accent := lipgloss.NewStyle().Foreground(t.Accent)

	cellRef := ui.ColIndexToLetter(m.cursorCol) + fmt.Sprintf("%d", m.cursorRow+1)

	// Left cluster: position + dimensions
	left := []string{
		bold.Render(cellRef),
		dim.Render(fmt.Sprintf("%d×%d", sheet.MaxRows, sheet.MaxCols)),
	}
	if m.showFormulas {
		left = append(left, accent.Render("ƒx"))
	}
	if m.freezeHeader {
		left = append(left, accent.Render("⌶"))
	}

	// Middle cluster: selection info or column stats
	middle := []string{}
	if m.isSelecting {
		rows := abs(m.selectEnd[0]-m.selectStart[0]) + 1
		cols := abs(m.selectEnd[1]-m.selectStart[1]) + 1
		middle = append(middle, lipgloss.NewStyle().Foreground(t.Primary).Bold(true).
			Render(fmt.Sprintf("▋ %d×%d", rows, cols)))
	} else if stats := m.computeColStats(); stats != nil {
		middle = append(middle,
			dim.Render("Σ ")+accent.Render(formatStatNum(stats.Sum)),
			dim.Render("avg ")+accent.Render(formatStatNum(stats.Avg)),
			dim.Render("min ")+accent.Render(formatStatNum(stats.Min)),
			dim.Render("max ")+accent.Render(formatStatNum(stats.Max)),
			dim.Render(fmt.Sprintf("n=%d", stats.Count)),
		)
	}

	// Right cluster: search + status message
	right := []string{}
	if len(m.searchResults) > 0 {
		right = append(right, lipgloss.NewStyle().Foreground(t.SearchMatch).Bold(true).
			Render(fmt.Sprintf("/%s  %d/%d", m.searchQuery, m.searchIndex+1, len(m.searchResults))))
	}
	if m.status.Message != "" {
		right = append(right, lipgloss.NewStyle().Foreground(ui.GetStatusColor(m.status.Type)).
			Render(m.status.Message))
	}

	sep := dim.Render("  │  ")
	var allParts []string
	if len(left) > 0 {
		allParts = append(allParts, strings.Join(left, "  "))
	}
	if len(middle) > 0 {
		allParts = append(allParts, strings.Join(middle, "  "))
	}
	if len(right) > 0 {
		allParts = append(allParts, strings.Join(right, "  "))
	}

	return m.styles.StatusBar.Render(strings.Join(allParts, sep))
}

// formatStatNum formats a float for compact status bar display.
func formatStatNum(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%.2f", v)
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

	content := m.styles.ModalTitle.Render("📊 Cell Details") + "\n\n"
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

	content := m.styles.ModalTitle.Render("🎯 Jump to Cell") + "\n\n"
	content += m.styles.ModalKey.Render("Enter cell reference:") + "\n"
	content += m.jumpInput.View() + "\n\n"
	content += lipgloss.NewStyle().Foreground(t.DimText).Render("Formats:\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  • A100   (column + row)\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  • 500    (row only)\n")
	content += lipgloss.NewStyle().Foreground(t.Text).Render("  • 10,5   (row,col)")

	return m.styles.Modal.Width(50).Render(content)
}

// renderExport renders the export modal
func (m Model) renderExport() string {
	t := theme.GetCurrentTheme()

	content := m.styles.ModalTitle.Render("💾 Export Sheet") + "\n\n"
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

	content := m.styles.ModalTitle.Render("🎨 Select Theme") + "\n\n"

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
		{"7", "Catppuccin Latte", "Light pastel theme"},
		{"8", "Solarized Light", "Balanced contrast"},
		{"9", "GitHub Light", "Clean and minimal"},
		{"0", "One Light", "Soft Atom colors"},
	}

	for _, theme := range themes {
		numStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(t.Text).Bold(true)
		descStyle := lipgloss.NewStyle().Foreground(t.DimText)

		current := ""
		if strings.Contains(strings.ToLower(theme.name), strings.ToLower(m.themeName)) ||
			strings.Contains(strings.ToLower(m.themeName), strings.ToLower(strings.ReplaceAll(theme.name, " ", "-"))) {
			current = lipgloss.NewStyle().Foreground(t.Accent).Render(" ✓")
		}

		content += numStyle.Render(theme.num) + "  " + nameStyle.Render(theme.name) + current + "\n"
		content += "   " + descStyle.Render(theme.desc) + "\n\n"
	}

	content += lipgloss.NewStyle().
		Foreground(t.DimText).
		Italic(true).
		Render("\nPress 1-9, 0 to select, Esc to cancel")

	modalW := ui.Max(50, ui.Min(m.width-4, 64))
	return m.styles.Modal.Width(modalW).Render(content)
}

// renderChart renders the chart visualization modal
func (m Model) renderChart() string {
	t := theme.GetCurrentTheme()

	modalW := ui.Max(60, ui.Min(m.width-4, 90))
	modalH := ui.Max(20, ui.Min(m.height-4, 36))
	barW := modalW - 26 // label(16) + " │ "(3) + value(7)
	if barW < 10 {
		barW = 10
	}
	chartH := modalH - 12 // reserve for type tabs + footer
	if chartH < 6 {
		chartH = 6
	}

	// Chart type tabs
	types := []struct{ num, name string }{
		{"1", "Bar"}, {"2", "Line"}, {"3", "Sparkline"}, {"4", "Pie"},
	}
	var tabs strings.Builder
	for i, typ := range types {
		if i == m.chartType {
			tabs.WriteString(lipgloss.NewStyle().
				Background(t.Accent).Foreground(lipgloss.Color("#000000")).
				Bold(true).Padding(0, 1).
				Render(typ.num+". "+typ.name) + " ")
		} else {
			tabs.WriteString(lipgloss.NewStyle().
				Foreground(t.DimText).Padding(0, 1).
				Render(typ.num+". "+typ.name) + " ")
		}
	}

	// Extract data
	startRow, startCol := m.selectStart[0], m.selectStart[1]
	endRow, endCol := m.selectEnd[0], m.selectEnd[1]
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}
	cd := extractChartData(m.sheets[m.currentSheet], startRow, startCol, endRow, endCol)

	divider := lipgloss.NewStyle().Foreground(t.Border).Render(strings.Repeat("─", modalW-4))

	var chartContent string
	colors := []lipgloss.Color{t.Accent, t.Primary, t.Secondary, t.Success, t.Warning}
	switch m.chartType {
	case 0:
		chartContent = renderBarChart(cd, barW, t.Accent, t.Text)
	case 1:
		chartContent = renderLineChart(cd, barW, chartH, t.Accent, t.Text)
	case 2:
		chartContent = "  " + renderSparkline(cd, t.Accent) + "\n\n" +
			renderBarChart(cd, barW, t.Accent, t.Text)
	case 3:
		chartContent = renderPieChart(cd, colors, t.Text)
	}

	if len(cd.Values) == 0 {
		chartContent = lipgloss.NewStyle().Foreground(t.DimText).
			Render("  No numeric data in selection.\n  Select a range with a label column + value column.")
	}

	content := m.styles.ModalTitle.Render("  Data Visualization") + "\n\n" +
		tabs.String() + "\n\n" +
		divider + "\n\n" +
		chartContent + "\n\n" +
		divider + "\n" +
		lipgloss.NewStyle().Foreground(t.DimText).Italic(true).
			Render("  1-4 switch type  •  p export as image  •  Esc close")

	return m.styles.Modal.Width(modalW).Height(modalH).Render(content)
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

// renderFilterMode renders the normal view with a filter input bar at the bottom.
func (m Model) renderFilterMode() string {
	base := m.renderNormal()
	t := theme.GetCurrentTheme()

	scope := fmt.Sprintf("col %s", ui.ColIndexToLetter(m.cursorCol))
	if m.filterColAll {
		scope = "all cols"
	}
	label := lipgloss.NewStyle().
		Background(t.Warning).Foreground(lipgloss.Color("#000000")).
		Bold(true).Padding(0, 1).
		Render(fmt.Sprintf("FILTER [%s]", scope))
	hint := lipgloss.NewStyle().Foreground(t.DimText).
		Render("  Tab=toggle scope  Enter=apply  Esc=cancel")
	bar := lipgloss.NewStyle().Background(t.Border).Padding(0, 1).
		Render(label + "  " + m.filterInput.View() + hint)

	return base + "\n" + bar
}

// renderDataProfile renders a per-column statistical summary.
func (m Model) renderDataProfile() string {
	sheet := m.sheets[m.currentSheet]
	t := theme.GetCurrentTheme()

	modalW := ui.Max(70, ui.Min(m.width-4, 90))

	title := m.styles.ModalTitle.Render("  Data Profile")
	summary := lipgloss.NewStyle().Foreground(t.DimText).
		Render(fmt.Sprintf("  %d rows  •  %d columns  •  %s",
			sheet.MaxRows, sheet.MaxCols, m.filename))

	header := lipgloss.NewStyle().Foreground(t.Secondary).Bold(true).
		Render(fmt.Sprintf("\n  %-4s %-16s %-8s %-8s %-8s %s\n",
			"Col", "Type", "Non-null", "Nulls", "Unique", "Stats"))
	divider := lipgloss.NewStyle().Foreground(t.Border).
		Render("  " + strings.Repeat("─", modalW-6) + "\n")

	var rows strings.Builder
	for col := 0; col < sheet.MaxCols; col++ {
		colLetter := ui.ColIndexToLetter(col)

		var vals []string
		nullCount := 0
		for _, row := range sheet.Rows {
			if col < len(row) {
				v := row[col].Value
				if v == "" {
					nullCount++
				} else {
					vals = append(vals, v)
				}
			} else {
				nullCount++
			}
		}

		nonNull := len(vals)

		// Detect type
		numericCount := 0
		var nums []float64
		for _, v := range vals {
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				numericCount++
				nums = append(nums, n)
			}
		}

		colType := "text"
		typeColor := t.Text
		statsStr := ""

		if numericCount == nonNull && nonNull > 0 {
			colType = "number"
			typeColor = t.Accent
			sum, mn, mx := nums[0], nums[0], nums[0]
			for _, n := range nums[1:] {
				sum += n
				if n < mn {
					mn = n
				}
				if n > mx {
					mx = n
				}
			}
			avg := sum / float64(len(nums))
			statsStr = fmt.Sprintf("Σ%-8s avg %-8s [%s – %s]",
				formatStatNum(sum), formatStatNum(avg),
				formatStatNum(mn), formatStatNum(mx))
		} else {
			// Unique count + sample values
			seen := make(map[string]struct{})
			for _, v := range vals {
				seen[v] = struct{}{}
			}
			unique := len(seen)
			sample := ""
			count := 0
			for v := range seen {
				if count >= 3 {
					break
				}
				if sample != "" {
					sample += ", "
				}
				if len(v) > 10 {
					v = v[:9] + "…"
				}
				sample += v
				count++
			}
			if unique > 3 {
				sample += "…"
			}
			statsStr = fmt.Sprintf("%d unique  %s", unique, sample)
			if numericCount > 0 && numericCount < nonNull {
				colType = "mixed"
				typeColor = t.Warning
			}
		}

		highlight := col == m.cursorCol
		rowStyle := lipgloss.NewStyle().Foreground(t.Text)
		if highlight {
			rowStyle = rowStyle.Background(t.RowHighlight)
		}

		line := fmt.Sprintf("  %-4s %-16s %-8d %-8d %-8s %s",
			colLetter,
			lipgloss.NewStyle().Foreground(typeColor).Render(colType),
			nonNull,
			nullCount,
			"",
			statsStr,
		)
		rows.WriteString(rowStyle.Render(line) + "\n")
	}

	content := title + "\n" + summary + "\n" + header + divider + rows.String() + "\n" +
		lipgloss.NewStyle().Foreground(t.DimText).Italic(true).
			Render("  Any key to close  •  highlighted column = cursor position")

	return m.styles.Modal.Width(modalW).Render(content)
}

// extractChartData extracts chart data from a selection.
// If the first column of the selection is non-numeric it becomes labels;
// otherwise row numbers are used as labels. The first numeric column is used
// as values.
func extractChartData(sheet models.Sheet, startRow, startCol, endRow, endCol int) chartData {
	data := chartData{
		Labels: make([]string, 0),
		Values: make([]float64, 0),
	}
	if startRow > endRow || startCol > endCol {
		return data
	}

	// Determine whether the first selected column is a label column (non-numeric)
	firstColNumeric := true
	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		if startCol < len(sheet.Rows[row]) {
			if _, err := strconv.ParseFloat(sheet.Rows[row][startCol].Value, 64); err != nil {
				if sheet.Rows[row][startCol].Value != "" {
					firstColNumeric = false
					break
				}
			}
		}
	}

	labelCol := -1
	valueCol := startCol
	if !firstColNumeric {
		labelCol = startCol
		// Find first numeric col after label col
		for col := startCol + 1; col <= endCol; col++ {
			for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
				if col < len(sheet.Rows[row]) {
					if _, err := strconv.ParseFloat(sheet.Rows[row][col].Value, 64); err == nil {
						valueCol = col
						goto foundValueCol
					}
				}
			}
		}
	}
foundValueCol:

	for row := startRow; row <= endRow && row < len(sheet.Rows); row++ {
		label := fmt.Sprintf("%d", row+1)
		if labelCol >= 0 && labelCol < len(sheet.Rows[row]) {
			if v := sheet.Rows[row][labelCol].Value; v != "" {
				label = v
			}
		}

		if valueCol < len(sheet.Rows[row]) {
			if val, err := strconv.ParseFloat(sheet.Rows[row][valueCol].Value, 64); err == nil {
				data.Labels = append(data.Labels, label)
				data.Values = append(data.Values, val)
			}
		}
	}
	return data
}

type chartData struct {
	Labels []string
	Values []float64
}

func renderBarChart(data chartData, barWidth int, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return "No data"
	}

	var b strings.Builder
	maxVal := maxFloat(data.Values)
	minVal := minFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}
	// Support negative values: shift so min is at 0
	offset := 0.0
	if minVal < 0 {
		offset = -minVal
	}
	effectiveMax := maxVal + offset

	const labelW = 14
	for i, val := range data.Values {
		if i >= len(data.Labels) {
			break
		}
		label := data.Labels[i]
		if len([]rune(label)) > labelW {
			label = string([]rune(label)[:labelW-1]) + "…"
		}
		labelStr := lipgloss.NewStyle().Foreground(textColor).Width(labelW).Render(label)

		filled := int(float64(barWidth) * ((val + offset) / effectiveMax))
		if filled < 0 {
			filled = 0
		}
		if filled > barWidth {
			filled = barWidth
		}
		empty := barWidth - filled

		barColor := accentColor
		if val < 0 {
			barColor = lipgloss.Color("#ff6b6b")
		}
		bar := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filled)) +
			strings.Repeat("░", empty)
		valStr := lipgloss.NewStyle().Foreground(textColor).Render(fmt.Sprintf(" %7s", formatStatNum(val)))

		b.WriteString("  " + labelStr + " │" + bar + valStr + "\n")
	}
	b.WriteString("  " + strings.Repeat(" ", labelW) + " └" + strings.Repeat("─", barWidth) + "\n")
	b.WriteString("  " + strings.Repeat(" ", labelW) + "   " +
		lipgloss.NewStyle().Foreground(textColor).Render(formatStatNum(minVal)) +
		strings.Repeat(" ", barWidth/2-4) +
		lipgloss.NewStyle().Foreground(textColor).Render(formatStatNum(maxVal)))

	return b.String()
}

func renderLineChart(data chartData, width, height int, accentColor, textColor lipgloss.Color) string {
	if len(data.Values) == 0 {
		return "No data"
	}
	if len(data.Values) == 1 {
		return renderBarChart(data, width, accentColor, textColor)
	}

	minVal := minFloat(data.Values)
	maxVal := maxFloat(data.Values)
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	// Build a grid of runes
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Map each value to (x,y) and draw connecting segments
	points := make([][2]int, len(data.Values))
	for i, val := range data.Values {
		x := int(float64(i) * float64(width-1) / float64(len(data.Values)-1))
		x = ui.Min(x, width-1)
		normalized := (val - minVal) / (maxVal - minVal)
		y := height - 1 - int(normalized*float64(height-1))
		y = ui.Max(0, ui.Min(y, height-1))
		points[i] = [2]int{x, y}
	}

	// Draw connecting lines between adjacent points using Bresenham
	for i := 0; i < len(points)-1; i++ {
		x0, y0 := points[i][0], points[i][1]
		x1, y1 := points[i+1][0], points[i+1][1]
		dx := x1 - x0
		if dx < 0 {
			dx = -dx
		}
		dy := y1 - y0
		if dy < 0 {
			dy = -dy
		}
		sx := 1
		if x0 > x1 {
			sx = -1
		}
		sy := 1
		if y0 > y1 {
			sy = -1
		}
		err := dx - dy
		for {
			if x0 >= 0 && x0 < width && y0 >= 0 && y0 < height {
				if grid[y0][x0] == ' ' {
					grid[y0][x0] = '·'
				}
			}
			if x0 == x1 && y0 == y1 {
				break
			}
			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x0 += sx
			}
			if e2 < dx {
				err += dx
				y0 += sy
			}
		}
	}
	// Draw data points on top
	for _, p := range points {
		if p[0] >= 0 && p[0] < width && p[1] >= 0 && p[1] < height {
			grid[p[1]][p[0]] = '●'
		}
	}

	var b strings.Builder
	for i, row := range grid {
		yVal := maxVal - (float64(i)/float64(height-1))*(maxVal-minVal)
		yLabel := lipgloss.NewStyle().Foreground(textColor).Width(8).Align(lipgloss.Right).
			Render(formatStatNum(yVal))
		b.WriteString("  " + yLabel + " │")
		for _, ch := range row {
			switch ch {
			case '●':
				b.WriteString(lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render("●"))
			case '·':
				b.WriteString(lipgloss.NewStyle().Foreground(accentColor).Render("·"))
			default:
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("  " + strings.Repeat(" ", 8) + " └" + strings.Repeat("─", width) + "\n")

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

func renderPieChart(data chartData, colors []lipgloss.Color, textColor lipgloss.Color) string {
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
