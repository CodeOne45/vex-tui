package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/odesaur/vex-tui/internal/app"
	"github.com/odesaur/vex-tui/internal/loader"
)

var version = "2.0.0"

const binaryName = "vex-tui"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[1]
	themeName := parseThemeFlag()

	// Validate file exists
	if err := validateFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load file
	sheets, err := loader.LoadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading file: %v\n", err)
		os.Exit(1)
	}

	if len(sheets) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No sheets found in file")
		os.Exit(1)
	}

	// Create and run application
	model := app.NewModel(filename, sheets, themeName)
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Excel TUI v%s - Modern Terminal Excel Viewer\n\n", version)
	fmt.Printf("Usage: %s <file> [--theme <name>]\n", binaryName)
	fmt.Println("\nAvailable themes:")
	for _, name := range app.GetThemeNames() {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println("\nExample:")
	fmt.Printf("  %s data.xlsx\n", binaryName)
	fmt.Printf("  %s report.csv --theme nord\n", binaryName)
}

func parseThemeFlag() string {
	themeName := "rose-pine" // default
	for i := 2; i < len(os.Args)-1; i++ {
		if os.Args[i] == "--theme" || os.Args[i] == "-t" {
			if i+1 < len(os.Args) {
				themeName = os.Args[i+1]
			}
		}
	}
	return themeName
}

func validateFile(filename string) error {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist", filename)
	}
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filename)
	}
	return nil
}
