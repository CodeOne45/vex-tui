package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/odesaur/vex-tui/internal/theme"
	"github.com/odesaur/vex-tui/pkg/models"
)

const (
	MinCellWidth = 12
	MaxCellWidth = 40
)

// Styles holds all lipgloss styles for the UI
type Styles struct {
	Title           lipgloss.Style
	Header          lipgloss.Style
	HeaderHighlight lipgloss.Style
	Cell            lipgloss.Style
	SelectedCell    lipgloss.Style
	RowHighlight    lipgloss.Style
	ColHighlight    lipgloss.Style
	SearchMatch     lipgloss.Style
	RowNum          lipgloss.Style
	SelectedRowNum  lipgloss.Style
	StatusBar       lipgloss.Style
	SearchBar       lipgloss.Style
	SearchPrompt    lipgloss.Style
	Modal           lipgloss.Style
	ModalTitle      lipgloss.Style
	ModalContent    lipgloss.Style
	ModalKey        lipgloss.Style
	ModalValue      lipgloss.Style
	Help            lipgloss.Style
	FormulaBar      lipgloss.Style
	Separator       lipgloss.Style
}

// InitStyles creates and returns styles based on current theme
func InitStyles() *Styles {
	t := theme.GetCurrentTheme()

	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			Background(t.Background).
			Padding(0, 1),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Secondary).
			Background(t.Border).
			Align(lipgloss.Left).
			Width(MinCellWidth),

		HeaderHighlight: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Secondary).
			Background(t.ColHighlight).
			Align(lipgloss.Left).
			Width(MinCellWidth),

		Cell: lipgloss.NewStyle().
			Foreground(t.Text).
			Width(MinCellWidth),

		SelectedCell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(t.CellHighlight).
			Bold(true).
			Width(MinCellWidth),

		RowHighlight: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.RowHighlight).
			Width(MinCellWidth),

		ColHighlight: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.ColHighlight).
			Width(MinCellWidth),

		SearchMatch: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(t.SearchMatch).
			Bold(true).
			Width(MinCellWidth),

		RowNum: lipgloss.NewStyle().
			Foreground(t.DimText).
			Align(lipgloss.Right).
			Width(5),

		SelectedRowNum: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.RowHighlight).
			Bold(true).
			Align(lipgloss.Right).
			Width(5),

		Separator: lipgloss.NewStyle().
			Foreground(t.Border),

		StatusBar: lipgloss.NewStyle().
			Background(t.Border).
			Padding(0, 1),

		SearchBar: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Border).
			Padding(0, 1),

		SearchPrompt: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),

		Modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Background(t.Background).
			Padding(1, 2).
			Width(60),

		ModalTitle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			Underline(true).
			MarginBottom(1),

		ModalContent: lipgloss.NewStyle().
			Foreground(t.Text),

		ModalKey: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Bold(true),

		ModalValue: lipgloss.NewStyle().
			Foreground(t.Text),

		Help: lipgloss.NewStyle().
			Foreground(t.DimText).
			Padding(0, 1),

		FormulaBar: lipgloss.NewStyle().
			Foreground(t.Accent).
			Background(t.Border).
			Padding(0, 1),
	}
}

// Helper functions

// ColIndexToLetter converts a 0-indexed column number to Excel-style letter
func ColIndexToLetter(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
}

// Truncate truncates a string to maxLen with ellipsis
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// TruncateToWidth ensures text fits exactly in the cell width
func TruncateToWidth(s string, width int) string {
	// Remove any newlines or tabs that would break layout
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	runes := []rune(s)

	if len(runes) > width {
		if width <= 3 {
			return strings.Repeat(".", width)
		}
		return string(runes[:width-3]) + "..."
	} else if len(runes) < width {
		return s + strings.Repeat(" ", width-len(runes))
	}

	return s
}

// PadCenter centers text in a field of given width
func PadCenter(s string, width int) string {
	if len(s) >= width {
		return TruncateToWidth(s, width)
	}

	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad

	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// PadRight pads a string on the right to fit the desired width
func PadRight(s string, width int) string {
	if len(s) >= width {
		return TruncateToWidth(s, width)
	}
	return s + strings.Repeat(" ", width-len(s))
}

// WrapText wraps text to fit within a specified width
func WrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > width {
			result.WriteString("\n")
			lineLen = 0
		}
		if i > 0 && lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += wordLen
	}

	return result.String()
}

// GetCellType determines the type of a cell
func GetCellType(cell models.Cell) string {
	if cell.Formula != "" {
		return "Formula"
	}
	if cell.Value == "" {
		return "Empty"
	}
	if _, err := strconv.ParseFloat(cell.Value, 64); err == nil {
		return "Number"
	}
	return "Text"
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetStatusColor returns the color for a status type
func GetStatusColor(statusType string) lipgloss.Color {
	t := theme.GetCurrentTheme()
	switch statusType {
	case models.StatusSuccess:
		return t.Success
	case models.StatusError:
		return t.Error
	case models.StatusWarning:
		return t.Warning
	default:
		return t.Secondary
	}
}

// RenderModal centers a modal in the viewport
func RenderModal(width, height int, modal string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}
