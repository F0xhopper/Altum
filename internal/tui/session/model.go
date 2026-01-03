/*
Copyright © 2025 Eden Phillips
*/
package session

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	stateSession sessionState = iota
	stateMilestone
	stateFocusQuality
	stateInterruptions
	stateReflection
	stateSaving
	stateDone
)

type model struct {
	state              sessionState
	stopwatch          stopwatch.Model
	milestoneInput     textinput.Model
	focusQualityInput  textinput.Model
	interruptionsInput textinput.Model
	reflectionInput    textinput.Model
	help               help.Model
	keyMap             KeyMap
	spinner            spinner.Model
	startTime          time.Time
	duration           time.Duration
	milestone          string
	focusQuality       string
	interruptions      string
	reflection         string
	dailyNotesPath     string
	dateFormat         string
	sessionCount       int
	noteFilePath       string
	err                error
}

func InitialModel(dailyNotesPath, dateFormat string) model {
	s := spinner.New()

	sw := stopwatch.NewWithInterval(time.Second)
	sw.Start()

	milestoneInput := textinput.New()
	milestoneInput.Placeholder = "What concrete outcome or milestone did you achieve?"
	milestoneInput.CharLimit = 200
	milestoneInput.Width = 80

	focusQualityInput := textinput.New()
	focusQualityInput.Placeholder = "Rate focus quality (1-5, optional, default 3)"
	focusQualityInput.CharLimit = 1
	focusQualityInput.Width = 80

	interruptionsInput := textinput.New()
	interruptionsInput.Placeholder = "Any interruptions or distractions worth noting? (optional)"
	interruptionsInput.CharLimit = 200
	interruptionsInput.Width = 80

	reflectionInput := textinput.New()
	reflectionInput.Placeholder = "Quick reflection / what went well or to improve? (optional)"
	reflectionInput.CharLimit = 200
	reflectionInput.Width = 80

	h := help.New()
	h.Width = 80

	return model{
		state:              stateSession,
		stopwatch:          sw,
		spinner:            s,
		milestoneInput:     milestoneInput,
		focusQualityInput:  focusQualityInput,
		interruptionsInput: interruptionsInput,
		reflectionInput:    reflectionInput,
		help:               h,
		keyMap:             DefaultKeyMap,
		startTime:          time.Now(),
		dailyNotesPath:     dailyNotesPath,
		dateFormat:         dateFormat,
		focusQuality:       "3",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.stopwatch.Init(),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case saveSuccessMsg:
		m = m.handleSaveSuccess(msg)
		return m, nil

	case saveErrorMsg:
		m = m.handleSaveError(msg)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateSession:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.stopSession):
				m.state = stateMilestone
				m.duration = m.stopwatch.Elapsed()
				m.stopwatch.Stop()
				m.milestoneInput.Focus()
				return m, nil
			}

		case stateMilestone:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Continue):
				m.milestone = m.milestoneInput.Value()
				if m.milestone == "" {
					return m, nil
				}
				m.state = stateFocusQuality
				m.milestoneInput.Blur()
				m.focusQualityInput.Focus()
				return m, nil
			}

		case stateFocusQuality:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Continue):
				value := m.focusQualityInput.Value()
				if value == "" {
					m.focusQuality = "3"
				} else {
					m.focusQuality = value
				}
				m.state = stateInterruptions
				m.focusQualityInput.Blur()
				m.interruptionsInput.Focus()
				return m, nil
			case key.Matches(msg, m.keyMap.Skip):
				m.state = stateInterruptions
				m.focusQualityInput.Blur()
				m.interruptionsInput.Focus()
				return m, nil
			}

		case stateInterruptions:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Continue):
				m.interruptions = m.interruptionsInput.Value()
				m.state = stateReflection
				m.interruptionsInput.Blur()
				m.reflectionInput.Focus()
				return m, nil
			case key.Matches(msg, m.keyMap.Skip):
				m.state = stateReflection
				m.interruptionsInput.Blur()
				m.reflectionInput.Focus()
				return m, nil
			case key.Matches(msg, m.keyMap.Back):
				m.state = stateFocusQuality
				m.interruptionsInput.Blur()
				m.focusQualityInput.Focus()
				return m, nil
			}

		case stateReflection:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Save):
				m.reflection = m.reflectionInput.Value()
				m.state = stateSaving
				return m, m.saveSession()
			case key.Matches(msg, m.keyMap.Back):
				m.state = stateInterruptions
				m.reflectionInput.Blur()
				m.interruptionsInput.Focus()
				return m, nil
			}

		case stateDone:
			switch {
			case key.Matches(msg, m.keyMap.Quit), key.Matches(msg, m.keyMap.Exit):
				return m, tea.Quit
			}
		}
	}

	switch m.state {
	case stateSession:
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		cmds = append(cmds, cmd)
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case stateMilestone:
		m.milestoneInput, cmd = m.milestoneInput.Update(msg)
		cmds = append(cmds, cmd)

	case stateFocusQuality:
		m.focusQualityInput, cmd = m.focusQualityInput.Update(msg)
		cmds = append(cmds, cmd)

	case stateInterruptions:
		m.interruptionsInput, cmd = m.interruptionsInput.Update(msg)
		cmds = append(cmds, cmd)

	case stateReflection:
		m.reflectionInput, cmd = m.reflectionInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	switch m.state {
	case stateSession:
		elapsed := m.stopwatch.Elapsed()
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60
		hours := int(elapsed.Hours())

		sessionTimerDisplay := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		if hours == 0 {
			sessionTimerDisplay = fmt.Sprintf("%02d:%02d", minutes, seconds)
		}

		s += TitleStyle.Render("Deep Work Session")
		s += "\n\n"
		s += SessionTimerStyle.Render(m.spinner.View() + " " + fmt.Sprintf("%s", sessionTimerDisplay))
		s += "\n\n"
		s += m.help.View(m.keyMap.sessionKeyMap())

	case stateMilestone:
		s += TitleStyle.Render("Session Milestone")
		s += "\n\n"
		s += "What concrete outcome or milestone did you achieve?\n"
		if m.milestoneInput.Focused() {
			s += FocusedStyle.Render(m.milestoneInput.View())
		} else {
			s += InputStyle.Render(m.milestoneInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.MilestoneKeyMap())

	case stateFocusQuality:
		s += TitleStyle.Render("Focus Quality")
		s += "\n\n"
		s += "How would you rate your focus quality? (1–5)\n"
		s += "(optional, default 3 if skipped)\n\n"
		if m.focusQualityInput.Focused() {
			s += FocusedStyle.Render(m.focusQualityInput.View())
		} else {
			s += InputStyle.Render(m.focusQualityInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.FocusQualityKeyMap())

	case stateInterruptions:
		s += TitleStyle.Render("Interruptions")
		s += "\n\n"
		s += "Any interruptions or distractions worth noting?\n"
		s += "(optional)\n\n"
		if m.interruptionsInput.Focused() {
			s += FocusedStyle.Render(m.interruptionsInput.View())
		} else {
			s += InputStyle.Render(m.interruptionsInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.InterruptionsKeyMap())

	case stateReflection:
		s += TitleStyle.Render("Reflection")
		s += "\n\n"
		s += "Quick reflection / what went well or to improve?\n"
		s += "(optional, free text)\n\n"
		if m.reflectionInput.Focused() {
			s += FocusedStyle.Render(m.reflectionInput.View())
		} else {
			s += InputStyle.Render(m.reflectionInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.ReflectionKeyMap())

	case stateSaving:
		s += TitleStyle.Render("Saving Session...")
		s += "\n\n"
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		} else {
			s += SuccessStyle.Render("Session saved successfully!")
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.SavingKeyMap())

	case stateDone:
		s += TitleStyle.Render("Session Complete")
		s += "\n\n"
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error saving session: %v", m.err))
		} else {
			s += SuccessStyle.Render(fmt.Sprintf("Session logged to: %s", m.noteFilePath))
			s += "\n\n"
			minutes := int(m.duration.Minutes())
			seconds := int(m.duration.Seconds()) % 60
			s += fmt.Sprintf("Duration: %d minutes %d seconds\n", minutes, seconds)
		}
		s += "\n"
		s += m.help.View(m.keyMap.DoneKeyMap())
	}

	return s
}
