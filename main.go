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
	// Only keep the version flag at the top level
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Print version information")
	flag.BoolVar(&versionFlag, "v", false, "Print version information (shorthand)")

	// Parse top-level flags
	flag.Parse()

	// Check if version flag was provided before doing anything else
	if versionFlag {
		fmt.Println(version.VersionInfo())
		return
	}

	args := flag.Args()

	// If no arguments or first argument is not "run", show help
	if len(args) == 0 || args[0] != "run" {
		fmt.Println("Usage: citizenship-tracker-cli run [options]")
		fmt.Println("  or:  citizenship-tracker-cli -v/--version")
		fmt.Println("\nRun citizenship-tracker-cli run -h for available options")
		return
	}

	// Handle 'run' command with flags
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	// Define flags specific to 'run' command
	headlessFlag := runCmd.Bool("headless", false, "Run program in headless mode")
	runCmd.BoolVar(headlessFlag, "l", false, "Run program in headless mode (shorthand)")

	statusFlag := runCmd.Bool("status", false, "Show status")
	runCmd.BoolVar(statusFlag, "s", false, "Show status (shorthand)")

	// Parse the remaining arguments
	runCmd.Parse(args[1:])

	// fmt.Printf("Run command flags: headless=%v, status=%v\n", *headlessFlag, *statusFlag)

	// Handle run command with flags
	if *headlessFlag {
		fmt.Println("Running in headless mode")
		app.RunHeadless()
		return
	}

	// fmt.Println("Running with status")
	if _, err := tea.NewProgram(app.InitialTeaModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
