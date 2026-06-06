/*
Copyright © 2025 Eden Phillips
*/
package menu

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem int

const (
	MenuStart MenuItem = iota
	MenuConfig
	MenuExit
)

type model struct {
	cursor   int
	selected MenuItem
	quitting bool
	keyMap   KeyMap
}

func InitialModel() model {
	return model{
		cursor:   0,
		selected: MenuStart,
		keyMap:   DefaultKeyMap,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c" || msg.String() == "q":
			m.quitting = true
			m.selected = MenuExit
			return m, tea.Quit

		case msg.String() == "up" || msg.String() == "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = 2
			}
			m.selected = MenuItem(m.cursor)

		case msg.String() == "down" || msg.String() == "j":
			if m.cursor < 2 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.selected = MenuItem(m.cursor)

		case msg.String() == "enter" || msg.String() == " ":
			m.quitting = true
			m.selected = MenuItem(m.cursor)
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
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
	s += LogoStyle.Render(logo)
	s += "\n"

	items := []string{
		"Start Deep Work Session",
		"️Configure Settings",
		"Exit",
	}

	for i, item := range items {
		if m.cursor == i {
			s += MenuItemSelectedStyle.Render(fmt.Sprintf("▶ %s", item))
		} else {
			s += MenuItemStyle.Render(fmt.Sprintf("  %s", item))
		}
		s += "\n"
	}

	s += MenuHelpStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit")

	return s
}
