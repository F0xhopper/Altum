/*
Copyright © 2025 Eden Phillips
*/
package session

import (
	"fmt"
	"strconv"
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
	dbPath             string
	sessionCount       int
	err                error
}

func InitialModel(dbPath string) model {
	s := spinner.New()

	sw := stopwatch.NewWithInterval(time.Second)
	sw.Start()

	milestoneInput := textinput.New()
	milestoneInput.Placeholder = "Describe what you accomplished..."
	milestoneInput.CharLimit = 200
	milestoneInput.Width = 80

	focusQualityInput := textinput.New()
	focusQualityInput.Placeholder = "1–5 (default 3)"
	focusQualityInput.CharLimit = 1
	focusQualityInput.Width = 80
	focusQualityInput.Validate = func(s string) error {
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 5 {
			return fmt.Errorf("must be 1–5")
		}
		return nil
	}

	interruptionsInput := textinput.New()
	interruptionsInput.Placeholder = "e.g. Slack notifications, phone call... (optional)"
	interruptionsInput.CharLimit = 200
	interruptionsInput.Width = 80

	reflectionInput := textinput.New()
	reflectionInput.Placeholder = "What went well? What to improve? (optional)"
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
		dbPath:             dbPath,
		focusQuality:       "3",
	}
}

func (m model) stepInfo() (int, int) {
	steps := map[sessionState]int{
		stateMilestone:     1,
		stateFocusQuality:  2,
		stateInterruptions: 3,
		stateReflection:    4,
	}
	if step, ok := steps[m.state]; ok {
		return step, 4
	}
	return 0, 0
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

	case stateSaving:
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
		hours := int(elapsed.Hours())
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		var sessionTimerDisplay string
		if hours == 0 {
			sessionTimerDisplay = fmt.Sprintf("%02d:%02d", minutes, seconds)
		} else {
			sessionTimerDisplay = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		s += TitleStyle.Render("Deep Work Session")
		s += "\n\n"
		s += SessionTimerStyle.Render(m.spinner.View() + " " + sessionTimerDisplay)
		s += "\n\n"
		s += m.help.View(m.keyMap.sessionKeyMap())

	case stateMilestone:
		step, total := m.stepInfo()
		s += TitleStyle.Render("Session Milestone")
		s += StepStyle.Render(fmt.Sprintf("  step %d of %d", step, total))
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
		step, total := m.stepInfo()
		s += TitleStyle.Render("Focus Quality")
		s += StepStyle.Render(fmt.Sprintf("  step %d of %d", step, total))
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
		step, total := m.stepInfo()
		s += TitleStyle.Render("Interruptions")
		s += StepStyle.Render(fmt.Sprintf("  step %d of %d", step, total))
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
		step, total := m.stepInfo()
		s += TitleStyle.Render("Reflection")
		s += StepStyle.Render(fmt.Sprintf("  step %d of %d", step, total))
		s += "\n\n"
		s += "Quick reflection / what went well or to improve?\n"
		s += "(optional)\n\n"
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
		s += m.spinner.View() + " Saving your session..."
		s += "\n\n"
		s += m.help.View(m.keyMap.SavingKeyMap())

	case stateDone:
		s += TitleStyle.Render("Session Complete")
		s += "\n\n"
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error saving session: %v", m.err))
		} else {
			s += SuccessStyle.Render(fmt.Sprintf("Session #%d saved to database", m.sessionCount))
			s += "\n\n"
			durationHours := int(m.duration.Hours())
			durationMins := int(m.duration.Minutes()) % 60
			durationSecs := int(m.duration.Seconds()) % 60
			var durationStr string
			if durationHours > 0 {
				durationStr = fmt.Sprintf("%dh %dm %ds", durationHours, durationMins, durationSecs)
			} else {
				durationStr = fmt.Sprintf("%dm %ds", durationMins, durationSecs)
			}
			s += fmt.Sprintf("Duration: %s\n", durationStr)
		}
		s += "\n"
		s += m.help.View(m.keyMap.DoneKeyMap())
	}

	return s
}
