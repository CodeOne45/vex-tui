package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keyboard shortcuts
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	Home         key.Binding
	End          key.Binding
	FirstCol     key.Binding
	LastCol      key.Binding
	NextSheet    key.Binding
	PrevSheet    key.Binding
	Search       key.Binding
	NextResult   key.Binding
	PrevResult   key.Binding
	ClearSearch  key.Binding
	Detail       key.Binding
	Jump         key.Binding
	ToggleForm   key.Binding
	Copy         key.Binding
	CopyRow      key.Binding
	Export       key.Binding
	Theme        key.Binding
	Help         key.Binding
	Quit         key.Binding
	Visualize    key.Binding
	SelectRange  key.Binding
	Edit         key.Binding
	Delete       key.Binding
	DeleteRow    key.Binding
	DeleteCol    key.Binding
	InsertRow    key.Binding
	InsertCol    key.Binding
	Paste        key.Binding
	Save         key.Binding
	SaveAs       key.Binding
	FillDown     key.Binding
	FillRight    key.Binding
	ApplyFormula key.Binding
}

// ShortHelp returns key bindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Edit, k.Save, k.Search, k.Jump, k.Help, k.Quit}
}

// FullHelp returns key bindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.PageUp, k.PageDown, k.FirstCol, k.LastCol},
		{k.Home, k.End, k.NextSheet, k.PrevSheet},
		{k.Edit, k.Delete, k.Copy, k.Paste},
		{k.InsertRow, k.InsertCol, k.DeleteRow, k.DeleteCol},
		{k.FillDown, k.FillRight, k.ApplyFormula, k.ToggleForm},
		{k.Search, k.NextResult, k.PrevResult, k.ClearSearch},
		{k.Detail, k.Jump, k.Export, k.Theme},
		{k.Save, k.SaveAs, k.Visualize, k.SelectRange},
		{k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:           key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:         key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:         key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:        key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		PageUp:       key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup/^u", "page up")),
		PageDown:     key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn/^d", "page down")),
		Home:         key.NewBinding(key.WithKeys("home", "0"), key.WithHelp("home/0", "row start")),
		End:          key.NewBinding(key.WithKeys("end", "$"), key.WithHelp("end/$", "row end")),
		FirstCol:     key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "first col")),
		LastCol:      key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "last col")),
		NextSheet:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next sheet")),
		PrevSheet:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("⇧tab", "prev sheet")),
		Search:       key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		NextResult:   key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next")),
		PrevResult:   key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "prev")),
		ClearSearch:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear")),
		Detail:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detail/edit")),
		Jump:         key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("^g", "jump")),
		ToggleForm:   key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "formulas")),
		Copy:         key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
		CopyRow:      key.NewBinding(key.WithKeys("C"), key.WithHelp("C", "copy row")),
		Export:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "export")),
		Theme:        key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:         key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Visualize:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "visualize")),
		SelectRange:  key.NewBinding(key.WithKeys("V"), key.WithHelp("V", "select")),
		Edit:         key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "edit")),
		Delete:       key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete cell")),
		DeleteRow:    key.NewBinding(key.WithKeys("d", "d"), key.WithHelp("dd", "delete row")),
		DeleteCol:    key.NewBinding(key.WithKeys("d", "c"), key.WithHelp("dc", "delete col")),
		InsertRow:    key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "insert row")),
		InsertCol:    key.NewBinding(key.WithKeys("O"), key.WithHelp("O", "insert col")),
		Paste:        key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "paste")),
		Save:         key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("^s", "save")),
		SaveAs:       key.NewBinding(key.WithKeys("ctrl+shift+s"), key.WithHelp("^⇧s", "save as")),
		FillDown:     key.NewBinding(key.WithKeys("ctrl+j"), key.WithHelp("^j", "fill down")),
		FillRight:    key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("^l", "fill right")),
		ApplyFormula: key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("^a", "apply formula")),
	}
}
