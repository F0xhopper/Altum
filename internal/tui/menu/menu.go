/*
Copyright © 2025 Eden Phillips
*/
package menu

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RunMenu() MenuItem {
	m := InitialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running menu: %v\n", err)
		return MenuExit
	}

	return finalModel.(model).selected
}
