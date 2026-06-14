package app

import (
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
)

// blankSheet mimics what the loader returns for a 0-byte CSV or an empty
// Excel sheet: no rows and zero dimensions. This is loader-agnostic, so it
// covers both the CSV and Excel empty-file cases at the layer where the bug
// actually lived.
func blankSheet() models.Sheet {
	return models.Sheet{Name: "blank.csv"}
}

// TestStartEditBlankSheetIsNonMutating verifies that pressing `i` on a blank
// file does not panic and does not mutate the sheet. Cancelling with Esc must
// leave an empty file empty so it round-trips byte-for-byte.
func TestStartEditBlankSheetIsNonMutating(t *testing.T) {
	m := NewModel("blank.csv", []models.Sheet{blankSheet()}, "catppuccin")

	m.startEdit() // previously panicked: index out of range on Rows[0][0]

	if got := len(m.sheets[0].Rows); got != 0 {
		t.Fatalf("startEdit grew an empty sheet (rows=%d); it should stay empty until commit", got)
	}
	if m.modified {
		t.Fatal("startEdit marked a clean sheet as modified before any commit")
	}
}

// TestCommitEditBlankSheetGrowsGrid verifies that committing a value on a
// blank sheet grows it to a usable shape and stores the value at A1.
func TestCommitEditBlankSheetGrowsGrid(t *testing.T) {
	m := NewModel("blank.csv", []models.Sheet{blankSheet()}, "catppuccin")

	m.startEdit()
	m.editInput.SetValue("hello")
	m.commitEdit()

	s := m.sheets[0]
	if s.MaxRows < 1 || s.MaxCols < 1 {
		t.Fatalf("commitEdit left unusable dimensions: rows=%d cols=%d", s.MaxRows, s.MaxCols)
	}
	if len(s.Rows) < 1 || len(s.Rows[0]) < 1 {
		t.Fatalf("commitEdit left unusable grid: rows=%d", len(s.Rows))
	}
	if got := s.Rows[0][0].Value; got != "hello" {
		t.Fatalf("A1 = %q, want %q", got, "hello")
	}
	if !m.modified {
		t.Fatal("commitEdit did not mark the sheet as modified")
	}
}
