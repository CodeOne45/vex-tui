package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete color scheme
type Theme struct {
	Name          string
	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Accent        lipgloss.Color
	Text          lipgloss.Color
	DimText       lipgloss.Color
	Background    lipgloss.Color
	Border        lipgloss.Color
	RowHighlight  lipgloss.Color
	ColHighlight  lipgloss.Color
	CellHighlight lipgloss.Color
	SearchMatch   lipgloss.Color
	Success       lipgloss.Color
	Error         lipgloss.Color
	Warning       lipgloss.Color
}

var (
	// Available themes inspired by modern design systems
	themes = map[string]Theme{
		"catppuccin": {
			Name:          "Catppuccin Mocha",
			Primary:       lipgloss.Color("#CBA6F7"),
			Secondary:     lipgloss.Color("#89DCEB"),
			Accent:        lipgloss.Color("#A6E3A1"),
			Text:          lipgloss.Color("#CDD6F4"),
			DimText:       lipgloss.Color("#6C7086"),
			Background:    lipgloss.Color("#1E1E2E"),
			Border:        lipgloss.Color("#313244"),
			RowHighlight:  lipgloss.Color("#181825"),
			ColHighlight:  lipgloss.Color("#313244"),
			CellHighlight: lipgloss.Color("#B4BEFE"),
			SearchMatch:   lipgloss.Color("#F9E2AF"),
			Success:       lipgloss.Color("#A6E3A1"),
			Error:         lipgloss.Color("#F38BA8"),
			Warning:       lipgloss.Color("#FAB387"),
		},
		"nord": {
			Name:          "Nord",
			Primary:       lipgloss.Color("#88C0D0"),
			Secondary:     lipgloss.Color("#81A1C1"),
			Accent:        lipgloss.Color("#A3BE8C"),
			Text:          lipgloss.Color("#ECEFF4"),
			DimText:       lipgloss.Color("#4C566A"),
			Background:    lipgloss.Color("#2E3440"),
			Border:        lipgloss.Color("#3B4252"),
			RowHighlight:  lipgloss.Color("#242933"),
			ColHighlight:  lipgloss.Color("#3B4252"),
			CellHighlight: lipgloss.Color("#8FBCBB"),
			SearchMatch:   lipgloss.Color("#EBCB8B"),
			Success:       lipgloss.Color("#A3BE8C"),
			Error:         lipgloss.Color("#BF616A"),
			Warning:       lipgloss.Color("#D08770"),
		},
		"rose-pine": {
			Name:          "Rosé Pine",
			Primary:       lipgloss.Color("#EBBCBA"),
			Secondary:     lipgloss.Color("#9CCFD8"),
			Accent:        lipgloss.Color("#F6C177"),
			Text:          lipgloss.Color("#E0DEF4"),
			DimText:       lipgloss.Color("#6E6A86"),
			Background:    lipgloss.Color("#191724"),
			Border:        lipgloss.Color("#26233A"),
			RowHighlight:  lipgloss.Color("#1F1D2E"),
			ColHighlight:  lipgloss.Color("#26233A"),
			CellHighlight: lipgloss.Color("#C4A7E7"),
			SearchMatch:   lipgloss.Color("#F6C177"),
			Success:       lipgloss.Color("#9CCFD8"),
			Error:         lipgloss.Color("#EB6F92"),
			Warning:       lipgloss.Color("#F6C177"),
		},
		"tokyo-night": {
			Name:          "Tokyo Night",
			Primary:       lipgloss.Color("#BB9AF7"),
			Secondary:     lipgloss.Color("#7DCFFF"),
			Accent:        lipgloss.Color("#9ECE6A"),
			Text:          lipgloss.Color("#C0CAF5"),
			DimText:       lipgloss.Color("#565F89"),
			Background:    lipgloss.Color("#1A1B26"),
			Border:        lipgloss.Color("#24283B"),
			RowHighlight:  lipgloss.Color("#16161E"),
			ColHighlight:  lipgloss.Color("#24283B"),
			CellHighlight: lipgloss.Color("#7AA2F7"),
			SearchMatch:   lipgloss.Color("#E0AF68"),
			Success:       lipgloss.Color("#9ECE6A"),
			Error:         lipgloss.Color("#F7768E"),
			Warning:       lipgloss.Color("#FF9E64"),
		},
		"gruvbox": {
			Name:          "Gruvbox Dark",
			Primary:       lipgloss.Color("#D3869B"),
			Secondary:     lipgloss.Color("#83A598"),
			Accent:        lipgloss.Color("#B8BB26"),
			Text:          lipgloss.Color("#EBDBB2"),
			DimText:       lipgloss.Color("#928374"),
			Background:    lipgloss.Color("#282828"),
			Border:        lipgloss.Color("#3C3836"),
			RowHighlight:  lipgloss.Color("#1D2021"),
			ColHighlight:  lipgloss.Color("#3C3836"),
			CellHighlight: lipgloss.Color("#FABD2F"),
			SearchMatch:   lipgloss.Color("#FE8019"),
			Success:       lipgloss.Color("#B8BB26"),
			Error:         lipgloss.Color("#FB4934"),
			Warning:       lipgloss.Color("#FABD2F"),
		},
		"dracula": {
			Name:          "Dracula",
			Primary:       lipgloss.Color("#BD93F9"),
			Secondary:     lipgloss.Color("#8BE9FD"),
			Accent:        lipgloss.Color("#50FA7B"),
			Text:          lipgloss.Color("#F8F8F2"),
			DimText:       lipgloss.Color("#6272A4"),
			Background:    lipgloss.Color("#282A36"),
			Border:        lipgloss.Color("#44475A"),
			RowHighlight:  lipgloss.Color("#21222C"),
			ColHighlight:  lipgloss.Color("#44475A"),
			CellHighlight: lipgloss.Color("#FFB86C"),
			SearchMatch:   lipgloss.Color("#F1FA8C"),
			Success:       lipgloss.Color("#50FA7B"),
			Error:         lipgloss.Color("#FF5555"),
			Warning:       lipgloss.Color("#FFB86C"),
		},
		// Light themes
		"catppuccin-latte": {
			Name:          "Catppuccin Latte",
			Primary:       lipgloss.Color("#8839EF"),
			Secondary:     lipgloss.Color("#04A5E5"),
			Accent:        lipgloss.Color("#40A02B"),
			Text:          lipgloss.Color("#4C4F69"),
			DimText:       lipgloss.Color("#9CA0B0"),
			Background:    lipgloss.Color("#EFF1F5"),
			Border:        lipgloss.Color("#DCE0E8"),
			RowHighlight:  lipgloss.Color("#E6E9EF"),
			ColHighlight:  lipgloss.Color("#DCE0E8"),
			CellHighlight: lipgloss.Color("#7287FD"),
			SearchMatch:   lipgloss.Color("#DF8E1D"),
			Success:       lipgloss.Color("#40A02B"),
			Error:         lipgloss.Color("#D20F39"),
			Warning:       lipgloss.Color("#FE640B"),
		},
		"solarized-light": {
			Name:          "Solarized Light",
			Primary:       lipgloss.Color("#268BD2"),
			Secondary:     lipgloss.Color("#2AA198"),
			Accent:        lipgloss.Color("#859900"),
			Text:          lipgloss.Color("#657B83"),
			DimText:       lipgloss.Color("#93A1A1"),
			Background:    lipgloss.Color("#FDF6E3"),
			Border:        lipgloss.Color("#EEE8D5"),
			RowHighlight:  lipgloss.Color("#EEE8D5"),
			ColHighlight:  lipgloss.Color("#EEE8D5"),
			CellHighlight: lipgloss.Color("#268BD2"),
			SearchMatch:   lipgloss.Color("#B58900"),
			Success:       lipgloss.Color("#859900"),
			Error:         lipgloss.Color("#DC322F"),
			Warning:       lipgloss.Color("#CB4B16"),
		},
		"github-light": {
			Name:          "GitHub Light",
			Primary:       lipgloss.Color("#0969DA"),
			Secondary:     lipgloss.Color("#8250DF"),
			Accent:        lipgloss.Color("#1A7F37"),
			Text:          lipgloss.Color("#1F2328"),
			DimText:       lipgloss.Color("#656D76"),
			Background:    lipgloss.Color("#FFFFFF"),
			Border:        lipgloss.Color("#D0D7DE"),
			RowHighlight:  lipgloss.Color("#F6F8FA"),
			ColHighlight:  lipgloss.Color("#EAEEF2"),
			CellHighlight: lipgloss.Color("#0969DA"),
			SearchMatch:   lipgloss.Color("#9A6700"),
			Success:       lipgloss.Color("#1A7F37"),
			Error:         lipgloss.Color("#CF222E"),
			Warning:       lipgloss.Color("#9A6700"),
		},
		"one-light": {
			Name:          "One Light",
			Primary:       lipgloss.Color("#A626A4"),
			Secondary:     lipgloss.Color("#0184BC"),
			Accent:        lipgloss.Color("#50A14F"),
			Text:          lipgloss.Color("#383A42"),
			DimText:       lipgloss.Color("#A0A1A7"),
			Background:    lipgloss.Color("#FAFAFA"),
			Border:        lipgloss.Color("#DBDBDC"),
			RowHighlight:  lipgloss.Color("#F0F0F0"),
			ColHighlight:  lipgloss.Color("#E5E5E6"),
			CellHighlight: lipgloss.Color("#4078F2"),
			SearchMatch:   lipgloss.Color("#C18401"),
			Success:       lipgloss.Color("#50A14F"),
			Error:         lipgloss.Color("#E45649"),
			Warning:       lipgloss.Color("#C18401"),
		},
	}

	// currentTheme is the active theme
	currentTheme = themes["catppuccin"]
)

// GetThemeNames returns all available theme names
func GetThemeNames() []string {
	return []string{"catppuccin", "nord", "rose-pine", "tokyo-night", "gruvbox", "dracula", "catppuccin-latte", "solarized-light", "github-light", "one-light"}
}

// SetTheme changes the current theme
func SetTheme(name string) bool {
	if t, ok := themes[name]; ok {
		currentTheme = t
		return true
	}
	return false
}

// GetCurrentTheme returns the current active theme
func GetCurrentTheme() Theme {
	return currentTheme
}
