package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vex/internal/theme"
	"github.com/vex/internal/ui"
	"github.com/vex/pkg/models"
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

	// Title bar
	title := fmt.Sprintf("📊 %s", m.filename)
	if len(m.sheets) > 1 {
		title += fmt.Sprintf(" • %s (%d/%d)", sheet.Name, m.currentSheet+1, len(m.sheets))
	} else {
		title += fmt.Sprintf(" • %s", sheet.Name)
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

	// Help
	b.WriteString("\n")
	b.WriteString(m.styles.Help.Render(m.help.ShortHelpView(m.keys.ShortHelp())))

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

	// Column headers
	b.WriteString(m.styles.RowNum.Render(""))
	b.WriteString(sep)

	for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
		colLetter := ui.ColIndexToLetter(col)
		if col == m.cursorCol {
			b.WriteString(m.styles.HeaderHighlight.Render(ui.PadCenter(colLetter, ui.MinCellWidth)))
		} else {
			b.WriteString(m.styles.Header.Render(ui.PadCenter(colLetter, ui.MinCellWidth)))
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

				cellText = ui.TruncateToWidth(cellText, ui.MinCellWidth)

				// Determine style
				var style lipgloss.Style
				if row == m.cursorRow && col == m.cursorCol {
					style = m.styles.SelectedCell
				} else if m.isSearchMatch(row, col) {
					style = m.styles.SearchMatch
				} else if row == m.cursorRow {
					style = m.styles.RowHighlight
				} else if col == m.cursorCol {
					style = m.styles.ColHighlight
				} else {
					style = m.styles.Cell
				}

				b.WriteString(style.Render(cellText))
				b.WriteString(sep)
			}
		} else {
			// Empty row
			for col := m.offsetCol; col < ui.Min(m.offsetCol+visibleCols, sheet.MaxCols); col++ {
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
				b.WriteString(style.Render(strings.Repeat(" ", ui.MinCellWidth)))
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
			Render(fmt.Sprintf("🔍 %d/%d", m.searchIndex+1, len(m.searchResults))))
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
		Render("\nPress 1-6 to select, Esc to cancel")

	return m.styles.Modal.Width(60).Render(content)
}
