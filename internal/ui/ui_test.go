package ui

import (
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
)

func TestColIndexToLetter(t *testing.T) {
	tests := []struct {
		index    int
		expected string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{701, "ZZ"},
		{702, "AAA"},
	}

	for _, tt := range tests {
		result := ColIndexToLetter(tt.index)
		if result != tt.expected {
			t.Errorf("ColIndexToLetter(%d) = %s; want %s", tt.index, result, tt.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"hi", 5, "hi"},
		{"test", 3, "tes"},
	}

	for _, tt := range tests {
		result := Truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q; want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected int
	}{
		{"hello", 10, 10},
		{"hello world", 5, 5},
		{"hi", 5, 5},
	}

	for _, tt := range tests {
		result := TruncateToWidth(tt.input, tt.width)
		if len(result) != tt.expected {
			t.Errorf("TruncateToWidth(%q, %d) length = %d; want %d", tt.input, tt.width, len(result), tt.expected)
		}
	}
}

func TestPadCenter(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected int
	}{
		{"hi", 10, 10},
		{"hello", 5, 5},
		{"test", 8, 8},
	}

	for _, tt := range tests {
		result := PadCenter(tt.input, tt.width)
		if len(result) != tt.expected {
			t.Errorf("PadCenter(%q, %d) length = %d; want %d", tt.input, tt.width, len(result), tt.expected)
		}
	}
}

func TestGetCellType(t *testing.T) {
	tests := []struct {
		cell     models.Cell
		expected string
	}{
		{models.Cell{Value: "", Formula: ""}, "Empty"},
		{models.Cell{Value: "test", Formula: ""}, "Text"},
		{models.Cell{Value: "123", Formula: ""}, "Number"},
		{models.Cell{Value: "123.45", Formula: ""}, "Number"},
		{models.Cell{Value: "100", Formula: "A1+A2"}, "Formula"},
	}

	for _, tt := range tests {
		result := GetCellType(tt.cell)
		if result != tt.expected {
			t.Errorf("GetCellType(%+v) = %s; want %s", tt.cell, result, tt.expected)
		}
	}
}

func TestMaxMin(t *testing.T) {
	if Max(5, 10) != 10 {
		t.Errorf("Max(5, 10) = %d; want 10", Max(5, 10))
	}
	if Max(10, 5) != 10 {
		t.Errorf("Max(10, 5) = %d; want 10", Max(10, 5))
	}
	if Min(5, 10) != 5 {
		t.Errorf("Min(5, 10) = %d; want 5", Min(5, 10))
	}
	if Min(10, 5) != 5 {
		t.Errorf("Min(10, 5) = %d; want 5", Min(10, 5))
	}
}