package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	// Create temp CSV file
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := "Name,Age,City\nAlice,30,NYC\nBob,25,LA\n"
	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	sheets, err := loadCSV(csvFile)
	if err != nil {
		t.Fatalf("loadCSV failed: %v", err)
	}

	if len(sheets) != 1 {
		t.Errorf("Expected 1 sheet, got %d", len(sheets))
	}

	sheet := sheets[0]
	if sheet.MaxRows != 3 {
		t.Errorf("Expected 3 rows, got %d", sheet.MaxRows)
	}
	if sheet.MaxCols != 3 {
		t.Errorf("Expected 3 cols, got %d", sheet.MaxCols)
	}
}

func TestLoadFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     string
		shouldError bool
	}{
		{"valid CSV", "test.csv", "A,B\n1,2\n", false},
		{"invalid extension", "test.txt", "data", true},
		{"nonexistent", "missing.csv", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.content != "" {
				path := filepath.Join(tmpDir, tt.filename)
				if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				tt.filename = path
			}

			_, err := LoadFile(tt.filename)
			if tt.shouldError && err == nil {
				t.Error("Expected error, got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestExportToCSV(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.csv")

	// Create temp CSV and load it
	inputFile := filepath.Join(tmpDir, "input.csv")
	content := "Name,Value\nTest,123\n"
	if err := os.WriteFile(inputFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	sheets, err := loadCSV(inputFile)
	if err != nil {
		t.Fatalf("Failed to load CSV: %v", err)
	}

	if err := ExportToCSV(sheets[0], outputFile); err != nil {
		t.Fatalf("ExportToCSV failed: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestSearchSheet(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := "Name,Age\nAlice,30\nBob,25\nCharlie,30\n"
	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	sheets, err := loadCSV(csvFile)
	if err != nil {
		t.Fatalf("loadCSV failed: %v", err)
	}

	results := SearchSheet(sheets[0], "30")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for '30', got %d", len(results))
	}

	results = SearchSheet(sheets[0], "alice")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'alice', got %d", len(results))
	}

	results = SearchSheet(sheets[0], "xyz")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'xyz', got %d", len(results))
	}
}