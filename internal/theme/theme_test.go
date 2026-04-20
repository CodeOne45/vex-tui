package theme

import "testing"

func TestGetThemeNames(t *testing.T) {
	names := GetThemeNames()
	expectedCount := 6

	if len(names) != expectedCount {
		t.Errorf("Expected %d themes, got %d", expectedCount, len(names))
	}

	expectedThemes := []string{"catppuccin", "nord", "rose-pine", "tokyo-night", "gruvbox", "dracula"}
	for _, expected := range expectedThemes {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected theme '%s' not found", expected)
		}
	}
}

func TestSetTheme(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"catppuccin", true},
		{"nord", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		result := SetTheme(tt.name)
		if result != tt.expected {
			t.Errorf("SetTheme(%q) = %v; want %v", tt.name, result, tt.expected)
		}
	}
}

func TestGetCurrentTheme(t *testing.T) {
	SetTheme("nord")
	theme := GetCurrentTheme()

	if theme.Name != "Nord" {
		t.Errorf("Expected theme name 'Nord', got '%s'", theme.Name)
	}

	SetTheme("catppuccin")
	theme = GetCurrentTheme()

	if theme.Name != "Catppuccin Mocha" {
		t.Errorf("Expected theme name 'Catppuccin Mocha', got '%s'", theme.Name)
	}
}