package app

import (
	"testing"

	"github.com/CodeOne45/vex-tui/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
)

func TestInit(t *testing.T) {
	sheets := []models.Sheet{{Name: "Test", MaxRows: 10, MaxCols: 5}}
	model := NewModel("test.csv", sheets, "catppuccin")

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	sheets := []models.Sheet{{Name: "Test", MaxRows: 10, MaxCols: 5}}
	model := NewModel("test.csv", sheets, "catppuccin")

	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updatedModel, _ := model.Update(msg)

	m := updatedModel.(Model)
	if m.width != 100 || m.height != 40 {
		t.Errorf("Expected size (100,40), got (%d,%d)", m.width, m.height)
	}
}

func TestAbsUpdate(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}