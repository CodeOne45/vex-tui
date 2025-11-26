package models

// Cell represents a single cell in the spreadsheet
type Cell struct {
	Value   string
	Formula string
	Row     int
	Col     int
}

// Sheet represents a worksheet with its data
type Sheet struct {
	Name    string
	Rows    [][]Cell
	MaxRows int
	MaxCols int
}

// Mode represents the current application mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeSearch
	ModeDetail
	ModeJump
	ModeExport
	ModeTheme
)

// StatusMsg represents a status message with type
type StatusMsg struct {
	Message string
	Type    string // "info", "success", "error", "warning"
}

// StatusType constants for type safety
const (
	StatusInfo    = "info"
	StatusSuccess = "success"
	StatusError   = "error"
	StatusWarning = "warning"
)
