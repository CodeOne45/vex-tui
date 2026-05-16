package loader

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadCSV_Empty verifies that loading a zero-byte CSV produces a sheet
// that is safe to navigate and edit (no panics when accessing cell A1).
// See https://github.com/CodeOne45/vex-tui — pressing `i` on a blank file
// previously panicked with index out of range.
func TestLoadCSV_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "blank.csv")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write blank.csv: %v", err)
	}

	sheets, err := LoadFile(path, 0)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if len(sheets) != 1 {
		t.Fatalf("got %d sheets, want 1", len(sheets))
	}
	s := sheets[0]
	if s.MaxRows < 1 || s.MaxCols < 1 {
		t.Fatalf("empty CSV produced unusable sheet: rows=%d cols=%d", s.MaxRows, s.MaxCols)
	}
	if len(s.Rows) < 1 || len(s.Rows[0]) < 1 {
		t.Fatalf("empty CSV produced unusable rows: rows=%d row0len=%d",
			len(s.Rows), func() int {
				if len(s.Rows) == 0 {
					return 0
				}
				return len(s.Rows[0])
			}())
	}
}
