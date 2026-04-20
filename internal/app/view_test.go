package app

import (
	"strings"
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
)

func TestView(t *testing.T) {
	sheets := []models.Sheet{
		{
			Name:    "Test",
			MaxRows: 3,
			MaxCols: 3,
			Rows: [][]models.Cell{
				{{Value: "A1"}, {Value: "B1"}, {Value: "C1"}},
				{{Value: "A2"}, {Value: "B2"}, {Value: "C2"}},
			},
		},
	}
	model := NewModel("test.csv", sheets, "catppuccin")
	model.width = 80
	model.height = 24

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !strings.Contains(view, "test.csv") {
		t.Error("View should contain filename")
	}
}

func TestRenderEmpty(t *testing.T) {
	model := NewModel("test.csv", []models.Sheet{}, "catppuccin")
	model.width = 80
	model.height = 24

	view := model.renderEmpty()
	if !strings.Contains(view, "No data") {
		t.Error("Empty view should contain 'No data' message")
	}
}

func TestIsInSelection(t *testing.T) {
	sheets := []models.Sheet{{Name: "Test", MaxRows: 10, MaxCols: 10}}
	model := NewModel("test.csv", sheets, "catppuccin")

	model.isSelecting = true
	model.selectStart = [2]int{2, 3}
	model.selectEnd = [2]int{5, 7}

	tests := []struct {
		row      int
		col      int
		expected bool
	}{
		{2, 3, true},
		{3, 5, true},
		{5, 7, true},
		{1, 1, false},
		{6, 6, false},
		{4, 2, false},
	}

	for _, tt := range tests {
		result := model.isInSelection(tt.row, tt.col)
		if result != tt.expected {
			t.Errorf("isInSelection(%d, %d) = %v; want %v", tt.row, tt.col, result, tt.expected)
		}
	}
}

func TestMaxFloat_View(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, 3},
		{[]float64{-1, -2, -3}, -1},
		{[]float64{5.5, 2.3, 8.1}, 8.1},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := maxFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("maxFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}

func TestMinFloat_(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, 1},
		{[]float64{-1, -2, -3}, -3},
		{[]float64{5.5, 2.3, 8.1}, 2.3},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := minFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("minFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}

func TestSumFloat_(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{-1, 1}, 0},
		{[]float64{}, 0},
		{[]float64{2.5, 2.5}, 5},
	}

	for _, tt := range tests {
		result := sumFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("sumFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}