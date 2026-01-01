/*
Copyright © 2025 Eden Phillips
*/
package menu

import "github.com/charmbracelet/lipgloss"

var (
	LogoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Bold(true).
			Align(lipgloss.Center).
			Margin(1, 0)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			PaddingLeft(2).
			Margin(1, 0)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7")).
				Bold(true).
				PaddingLeft(2).
				Margin(1, 0)

	MenuHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(2).
			Align(lipgloss.Center)
)
