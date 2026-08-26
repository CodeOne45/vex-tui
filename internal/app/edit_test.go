package app

import (
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
)

// TestInsertColumnShortRow reproduces the crash from issue #37: inserting a
// column at a cursor position past the end of a shorter row used to panic
// with "slice bounds out of range".
func TestInsertColumnShortRow(t *testing.T) {
	m := &Model{
		sheets: []models.Sheet{
			{
				Rows: [][]models.Cell{
					{{Row: 0, Col: 0, Value: "foo"}, {Row: 0, Col: 1, Value: "bar"}, {Row: 0, Col: 2, Value: "baz"}, {Row: 0, Col: 3, Value: "qux"}},
					{{Row: 1, Col: 0, Value: "foo"}, {Row: 1, Col: 1, Value: "bar"}},
					{{Row: 2, Col: 0, Value: "foo"}, {Row: 2, Col: 1, Value: "bar"}, {Row: 2, Col: 2, Value: "baz"}},
				},
				MaxCols: 4,
			},
		},
		cursorCol: 3,
	}

	m.insertColumn()

	for i, row := range m.sheets[0].Rows {
		if len(row) <= m.cursorCol {
			t.Fatalf("row %d: expected at least %d cells after insert, got %d", i, m.cursorCol+1, len(row))
		}
	}

	if m.sheets[0].MaxCols != 5 {
		t.Fatalf("expected MaxCols 5, got %d", m.sheets[0].MaxCols)
	}
}
