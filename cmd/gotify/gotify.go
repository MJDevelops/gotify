package main

import (
	"fmt"
	"os"

	"github.com/MJDevelops/gotify/internal/app/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialAuthSelect())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}
