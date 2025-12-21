package chart

import (
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/charmbracelet/lipgloss"
)

func TestExtractChartData(t *testing.T) {
	sheet := models.Sheet{
		Name:    "Test",
		MaxRows: 5,
		MaxCols: 3,
		Rows: [][]models.Cell{
			{{Value: "Product"}, {Value: "Sales"}},
			{{Value: "Laptop"}, {Value: "100"}},
			{{Value: "Mouse"}, {Value: "50"}},
			{{Value: "Keyboard"}, {Value: "75"}},
		},
	}

	data := ExtractChartData(sheet, 1, 0, 3, 1)

	if len(data.Labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(data.Labels))
	}

	if len(data.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(data.Values))
	}

	expectedLabels := []string{"Laptop", "Mouse", "Keyboard"}
	for i, label := range expectedLabels {
		if data.Labels[i] != label {
			t.Errorf("Expected label '%s', got '%s'", label, data.Labels[i])
		}
	}

	expectedValues := []float64{100, 50, 75}
	for i, val := range expectedValues {
		if data.Values[i] != val {
			t.Errorf("Expected value %f, got %f", val, data.Values[i])
		}
	}
}

func TestRenderBarChart(t *testing.T) {
	data := ChartData{
		Labels: []string{"A", "B", "C"},
		Values: []float64{10, 20, 15},
		Title:  "Test Chart",
	}

	style := lipgloss.NewStyle()
	result := RenderBarChart(data, style, lipgloss.Color("#00FF00"), lipgloss.Color("#FFFFFF"))

	if result == "" {
		t.Error("Bar chart should not be empty")
	}
}

func TestRenderLineChart(t *testing.T) {
	data := ChartData{
		Labels: []string{"A", "B", "C"},
		Values: []float64{10, 20, 15},
	}

	style := lipgloss.NewStyle()
	result := RenderLineChart(data, style, lipgloss.Color("#00FF00"), lipgloss.Color("#FFFFFF"))

	if result == "" {
		t.Error("Line chart should not be empty")
	}
}

func TestRenderSparkline(t *testing.T) {
	data := ChartData{
		Values: []float64{1, 5, 3, 8, 2},
	}

	result := RenderSparkline(data, lipgloss.Color("#00FF00"))

	if result == "" {
		t.Error("Sparkline should not be empty")
	}

	// Should return 5 characters
	if len([]rune(result)) < 5 {
		t.Errorf("Expected at least 5 characters in sparkline, got %d", len([]rune(result)))
	}
}

func TestRenderPieChart(t *testing.T) {
	data := ChartData{
		Labels: []string{"A", "B"},
		Values: []float64{60, 40},
	}

	style := lipgloss.NewStyle()
	colors := []lipgloss.Color{lipgloss.Color("#FF0000"), lipgloss.Color("#00FF00")}
	result := RenderPieChart(data, style, colors, lipgloss.Color("#FFFFFF"))

	if result == "" {
		t.Error("Pie chart should not be empty")
	}
}

func TestMaxFloat(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 5, 3}, 5},
		{[]float64{-5, -2, -10}, -2},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := maxFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("maxFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}

func TestMinFloat(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 5, 3}, 1},
		{[]float64{-5, -2, -10}, -10},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := minFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("minFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}

func TestSumFloat(t *testing.T) {
	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{-1, 1}, 0},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := sumFloat(tt.vals)
		if result != tt.expected {
			t.Errorf("sumFloat(%v) = %f; want %f", tt.vals, result, tt.expected)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}