package models

import "testing"

func TestStatusConstants(t *testing.T) {
	if StatusInfo != "info" {
		t.Errorf("StatusInfo = %s; want 'info'", StatusInfo)
	}
	if StatusSuccess != "success" {
		t.Errorf("StatusSuccess = %s; want 'success'", StatusSuccess)
	}
	if StatusError != "error" {
		t.Errorf("StatusError = %s; want 'error'", StatusError)
	}
	if StatusWarning != "warning" {
		t.Errorf("StatusWarning = %s; want 'warning'", StatusWarning)
	}
}

func TestModeConstants(t *testing.T) {
	modes := []Mode{ModeNormal, ModeSearch, ModeDetail, ModeJump, ModeExport, ModeTheme, ModeChart, ModeSelectRange}

	if len(modes) != 8 {
		t.Errorf("Expected 8 modes, got %d", len(modes))
	}

	// Ensure modes are distinct
	seen := make(map[Mode]bool)
	for _, mode := range modes {
		if seen[mode] {
			t.Errorf("Duplicate mode value: %d", mode)
		}
		seen[mode] = true
	}
}