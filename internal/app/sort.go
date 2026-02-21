package app

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/CodeOne45/vex-tui/internal/ui"
	"github.com/CodeOne45/vex-tui/pkg/models"
)

type sortableValue struct {
	empty    bool
	isNumber bool
	number   float64
	text     string
}

// sortByCurrentColumn sorts rows on the active column.
// If cursor is on the first row, row 1 is treated as header and kept in place.
func (m *Model) sortByCurrentColumn(ascending bool) {
	sheet := &m.sheets[m.currentSheet]
	if sheet.MaxRows <= 1 || len(sheet.Rows) <= 1 {
		m.status = models.StatusMsg{Message: "Not enough rows to sort", Type: models.StatusWarning}
		return
	}

	startRow := 0
	if m.cursorRow == 0 {
		startRow = 1
	}
	if startRow >= len(sheet.Rows)-1 {
		m.status = models.StatusMsg{Message: "Not enough rows to sort", Type: models.StatusWarning}
		return
	}

	rows := sheet.Rows[startRow:]
	col := m.cursorCol

	sort.SliceStable(rows, func(i, j int) bool {
		a := cellSortValue(rows[i], col)
		b := cellSortValue(rows[j], col)
		return compareSortableValues(a, b, ascending)
	})

	for rowIdx := range sheet.Rows {
		for colIdx := range sheet.Rows[rowIdx] {
			sheet.Rows[rowIdx][colIdx].Row = rowIdx
			sheet.Rows[rowIdx][colIdx].Col = colIdx
		}
	}

	m.modified = true
	m.recalculateFormulas()

	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	colRef := ui.ColIndexToLetter(col)
	if startRow == 1 {
		m.status = models.StatusMsg{
			Message: fmt.Sprintf("Sorted column %s (%s), header kept", colRef, direction),
			Type:    models.StatusSuccess,
		}
		return
	}

	m.status = models.StatusMsg{
		Message: fmt.Sprintf("Sorted column %s (%s)", colRef, direction),
		Type:    models.StatusSuccess,
	}
}

func cellSortValue(row []models.Cell, col int) sortableValue {
	if col >= len(row) {
		return sortableValue{empty: true}
	}

	raw := strings.TrimSpace(row[col].Value)
	if raw == "" {
		return sortableValue{empty: true}
	}

	if n, err := strconv.ParseFloat(raw, 64); err == nil {
		return sortableValue{isNumber: true, number: n}
	}

	return sortableValue{text: strings.ToLower(raw)}
}

func compareSortableValues(a, b sortableValue, ascending bool) bool {
	if a.empty != b.empty {
		return !a.empty
	}

	if a.isNumber && b.isNumber {
		if a.number == b.number {
			return false
		}
		if ascending {
			return a.number < b.number
		}
		return a.number > b.number
	}

	if a.isNumber != b.isNumber {
		return a.isNumber
	}

	if a.text == b.text {
		return false
	}
	if ascending {
		return a.text < b.text
	}
	return a.text > b.text
}
