/*
Copyright © 2025 Eden Phillips

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuItem int

const (
	menuStart menuItem = iota
	menuConfig
	menuExit
)

type menuModel struct {
	cursor   int
	selected menuItem
	quitting bool
}

var (
	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Bold(true).
			Align(lipgloss.Center).
			Margin(1, 0)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			PaddingLeft(2).
			Margin(1, 0)

	menuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7")).
				Bold(true).
				PaddingLeft(2).
				Margin(1, 0)

	menuHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(2).
			Align(lipgloss.Center)
)

func initialMenuModel() menuModel {
	return menuModel{
		cursor:   0,
		selected: menuStart,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.selected = menuExit
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = 2
			}
			m.selected = menuItem(m.cursor)

		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.selected = menuItem(m.cursor)

		case "enter", " ":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m menuModel) View() string {
	if m.quitting {
		return ""
	}

	var s string

	logo := `
     █████╗ ██╗  ████████╗██╗   ██╗███╗   ███╗
    ██╔══██╗██║  ╚══██╔══╝██║   ██║████╗ ████║
    ███████║██║     ██║   ██║   ██║██╔████╔██║
    ██╔══██║██║     ██║   ██║   ██║██║╚██╔╝██║
    ██║  ██║███████╗██║   ╚██████╔╝██║ ╚═╝ ██║
    ╚═╝  ╚═╝╚══════╝╚═╝    ╚═════╝ ╚═╝     ╚═╝
    
        Deep Work Companion
`
	s += logoStyle.Render(logo)
	s += "\n"

	items := []string{
		"Start Deep Work Session",
		"️Configure Settings",
		"Exit",
	}

	for i, item := range items {
		if m.cursor == i {
			s += menuItemSelectedStyle.Render(fmt.Sprintf("▶ %s", item))
		} else {
			s += menuItemStyle.Render(fmt.Sprintf("  %s", item))
		}
		s += "\n"
	}

	s += menuHelpStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit")

	return s
}

func runMenu() menuItem {
	m := initialMenuModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running menu: %v\n", err)
		return menuExit
	}

	return m.selected
}
