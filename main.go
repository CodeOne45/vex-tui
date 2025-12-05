package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CodeOne45/vex-tui/internal/app"
	"github.com/CodeOne45/vex-tui/internal/loader"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "1.1.1"

var (
	showVersion = flag.Bool("version", false, "Show version information")
	showHelp    = flag.Bool("help", false, "Show help information")
	themeName   = flag.String("theme", "catppuccin", "Set the color theme")
)

func main() {
	flag.StringVar(themeName, "t", "catppuccin", "Set the color theme (shorthand)")
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	filename := args[0]

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
	model := app.NewModel(filename, sheets, *themeName)
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

func printVersion() {
	fmt.Printf("vex version %s\n", version)
	fmt.Println("A beautiful terminal-based Excel and CSV viewer")
	fmt.Println("\nProject: https://github.com/CodeOne45/vex-tui")
}

func printHelp() {
	fmt.Printf("vex v%s - Terminal Excel Viewer\n\n", version)
	fmt.Println("USAGE:")
	fmt.Println("  vex [OPTIONS] <file>")
	fmt.Println("\nARGUMENTS:")
	fmt.Println("  <file>    Path to Excel (.xlsx, .xlsm, .xls) or CSV file")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  -t, --theme <name>    Set color theme (default: catppuccin)")
	fmt.Println("  --version             Show version information")
	fmt.Println("  --help                Show this help message")
	fmt.Println("\nAVAILABLE THEMES:")
	for _, name := range app.GetThemeNames() {
		fmt.Printf("  • %s\n", name)
	}
	fmt.Println("\nEXAMPLES:")
	fmt.Println("  vex data.xlsx")
	fmt.Println("  vex report.csv --theme nord")
	fmt.Println("  vex sales.xlsx -t tokyo-night")
	fmt.Println("\nKEYBOARD SHORTCUTS:")
	fmt.Println("  Navigation:  ↑↓←→ / hjkl, PgUp/PgDn, Home/End")
	fmt.Println("  Sheets:      Tab / Shift+Tab")
	fmt.Println("  Search:      / (search), n (next), N (prev)")
	fmt.Println("  Actions:     Enter (details), Ctrl+G (jump), c (copy)")
	fmt.Println("  Data viz:    V (select range), v (visualize)")
	fmt.Println("  Other:       e (export), t (theme), f (formulas), ? (help), q (quit)")
	fmt.Println("\nFor more information, visit: https://github.com/CodeOne45/vex-tui")
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "vex: missing file argument\n\n")
	fmt.Fprintf(os.Stderr, "Usage: vex [OPTIONS] <file>\n")
	fmt.Fprintf(os.Stderr, "Try 'vex --help' for more information.\n")
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
