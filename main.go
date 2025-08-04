package main

import (
	"fmt"
	"os"
	"text-editor/internal/reader"
	"text-editor/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: editor <file>")
		os.Exit(1)
	}
	capacity := 1024
	tb, filePath, err := reader.ReadFile(os.Args[1], capacity)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New(tb, filePath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
