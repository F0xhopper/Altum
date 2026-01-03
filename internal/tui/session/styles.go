/*
Copyright © 2025 Eden Phillips
*/
package session

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("7")).Padding(1, 2)
	SessionTimerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("7")).Align(lipgloss.Center).Padding(1)
	SuccessStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true)
	ErrorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Bold(true)
	InputStyle        = lipgloss.NewStyle().BorderForeground(lipgloss.Color("8")).BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
	FocusedStyle      = lipgloss.NewStyle().BorderForeground(lipgloss.Color("7")).BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
)
