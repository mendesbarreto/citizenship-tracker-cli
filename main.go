package main

import (
	"citizenship-tracker-cli/pkg/app"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(app.InitialTeaModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
