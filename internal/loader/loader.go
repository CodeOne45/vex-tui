package loader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/odesaur/vex-tui/pkg/models"
	"github.com/xuri/excelize/v2"
)

// LoadFile loads an Excel or CSV file and returns the sheets
func LoadFile(filename string) ([]models.Sheet, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".xlsx", ".xlsm", ".xls":
		return loadExcel(filename)
	case ".csv":
		return loadCSV(filename)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (supported: .xlsx, .xlsm, .xls, .csv)", ext)
	}
}

// loadExcel loads an Excel file
func loadExcel(filename string) ([]models.Sheet, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// Log error but don't override return error
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
		}
	}()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	sheets := make([]models.Sheet, 0, len(sheetList))

	for _, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			// Skip sheets that can't be read
			continue
		}

		sheet := models.Sheet{
			Name:    sheetName,
			MaxRows: len(rows),
			Rows:    make([][]models.Cell, 0, len(rows)),
		}

		// Convert to Cell structure and track max columns
		for rowIdx, row := range rows {
			cellRow := make([]models.Cell, 0, len(row))
			for colIdx, cellValue := range row {
				// Get cell reference
				cellRef, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)

				// Try to get formula (ignore error as not all cells have formulas)
				formula, _ := f.GetCellFormula(sheetName, cellRef)

				cell := models.Cell{
					Value:   cellValue,
					Formula: formula,
					Row:     rowIdx,
					Col:     colIdx,
				}
				cellRow = append(cellRow, cell)

				if colIdx+1 > sheet.MaxCols {
					sheet.MaxCols = colIdx + 1
				}
			}
			sheet.Rows = append(sheet.Rows, cellRow)
		}

		sheets = append(sheets, sheet)
	}

	return sheets, nil
}

// loadCSV loads a CSV file
func loadCSV(filename string) ([]models.Sheet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
		}
	}()

	reader := csv.NewReader(file)
	reader.ReuseRecord = true // Performance optimization
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	sheet := models.Sheet{
		Name:    filepath.Base(filename),
		MaxRows: len(records),
		Rows:    make([][]models.Cell, 0, len(records)),
	}

	for rowIdx, record := range records {
		cellRow := make([]models.Cell, 0, len(record))
		for colIdx, value := range record {
			cell := models.Cell{
				Value: value,
				Row:   rowIdx,
				Col:   colIdx,
			}
			cellRow = append(cellRow, cell)

			if colIdx+1 > sheet.MaxCols {
				sheet.MaxCols = colIdx + 1
			}
		}
		sheet.Rows = append(sheet.Rows, cellRow)
	}

	return []models.Sheet{sheet}, nil
}

// ExportToCSV exports a sheet to CSV format
func ExportToCSV(sheet models.Sheet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range sheet.Rows {
		record := make([]string, 0, len(row))
		for _, cell := range row {
			record = append(record, cell.Value)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	return nil
}

// ExportToJSON exports a sheet to JSON format
func ExportToJSON(sheet models.Sheet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
		}
	}()

	// Use first row as headers if available
	if len(sheet.Rows) == 0 {
		return fmt.Errorf("sheet is empty")
	}

	data := make([]map[string]string, 0, len(sheet.Rows)-1)
	headers := sheet.Rows[0]

	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		record := make(map[string]string)

		for j, cell := range row {
			headerKey := fmt.Sprintf("col_%d", j)
			if j < len(headers) && headers[j].Value != "" {
				headerKey = headers[j].Value
			}
			record[headerKey] = cell.Value
		}
		data = append(data, record)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// SearchSheet searches for a term in the sheet
func SearchSheet(sheet models.Sheet, term string) []models.Cell {
	if term == "" {
		return nil
	}

	term = strings.ToLower(term)
	results := make([]models.Cell, 0)

	for _, row := range sheet.Rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell.Value), term) ||
				strings.Contains(strings.ToLower(cell.Formula), term) {
				results = append(results, cell)
			}
		}
	}

	return results
}
