/*
Copyright © 2025 Eden Phillips
*/
package session

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("7")).Padding(1, 2)
	SessionTimerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("7")).Width(30).Align(lipgloss.Center).Padding(1)
	SuccessStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	ErrorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	InputStyle        = lipgloss.NewStyle().BorderForeground(lipgloss.Color("8")).BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
	FocusedStyle      = lipgloss.NewStyle().BorderForeground(lipgloss.Color("7")).BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
	StepStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
