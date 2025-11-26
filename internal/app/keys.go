package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keyboard shortcuts
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Home        key.Binding
	End         key.Binding
	FirstCol    key.Binding
	LastCol     key.Binding
	NextSheet   key.Binding
	PrevSheet   key.Binding
	Search      key.Binding
	NextResult  key.Binding
	PrevResult  key.Binding
	ClearSearch key.Binding
	Detail      key.Binding
	Jump        key.Binding
	ToggleForm  key.Binding
	Copy        key.Binding
	CopyRow     key.Binding
	Export      key.Binding
	Theme       key.Binding
	Help        key.Binding
	Quit        key.Binding
}

// ShortHelp returns key bindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Jump, k.Detail, k.Theme, k.Help, k.Quit}
}

// FullHelp returns key bindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.PageUp, k.PageDown, k.FirstCol, k.LastCol},
		{k.Home, k.End, k.NextSheet, k.PrevSheet},
		{k.Search, k.NextResult, k.PrevResult, k.ClearSearch},
		{k.Detail, k.Jump, k.ToggleForm},
		{k.Copy, k.CopyRow, k.Export, k.Theme},
		{k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:        key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:       key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		PageUp:      key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup/^u", "page up")),
		PageDown:    key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn/^d", "page down")),
		Home:        key.NewBinding(key.WithKeys("home", "0"), key.WithHelp("home/0", "row start")),
		End:         key.NewBinding(key.WithKeys("end", "$"), key.WithHelp("end/$", "row end")),
		FirstCol:    key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "first col")),
		LastCol:     key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "last col")),
		NextSheet:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next sheet")),
		PrevSheet:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("⇧tab", "prev sheet")),
		Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		NextResult:  key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next")),
		PrevResult:  key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "prev")),
		ClearSearch: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear")),
		Detail:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detail")),
		Jump:        key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("^g", "jump")),
		ToggleForm:  key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "formulas")),
		Copy:        key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
		CopyRow:     key.NewBinding(key.WithKeys("C"), key.WithHelp("C", "copy row")),
		Export:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "export")),
		Theme:       key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}
