package main

import (
	"citizenship-tracker-cli/pkg/app"
	"citizenship-tracker-cli/pkg/version"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define version flag
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.BoolVar(versionFlag, "v", false, "Print version information (shorthand)")

	// Parse flags
	flag.Parse()

	// Check if version flag was provided
	if *versionFlag {
		fmt.Println(version.VersionInfo())
		return
	}

	// Run the application UI
	if _, err := tea.NewProgram(app.InitialTeaModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
