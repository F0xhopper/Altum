/*
Copyright © 2025 Eden Phillips
*/
package timer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	stateTimer sessionState = iota
	stateRating
	stateNotes
	stateSaving
	stateDone
)

type model struct {
	state          sessionState
	stopwatch      stopwatch.Model
	ratingInput    textinput.Model
	notesInput     textarea.Model
	help           help.Model
	keyMap         KeyMap
	spinner        spinner.Model
	startTime      time.Time
	duration       time.Duration
	rating         string
	notes          string
	dailyNotesPath string
	dateFormat     string
	sessionCount   int
	noteFilePath   string
	err            error
}

func InitialModel(dailyNotesPath, dateFormat string) model {
	s := spinner.New()
	s.Spinner.FPS = 8

	sw := stopwatch.NewWithInterval(time.Second)
	sw.Start()
	ratingInput := textinput.New()
	ratingInput.Placeholder = "Enter rating (1-10)"
	ratingInput.Focus()
	ratingInput.CharLimit = 2
	ratingInput.Width = 20

	notesInput := textarea.New()
	notesInput.Placeholder = "Enter notes about your session..."
	notesInput.SetWidth(60)
	notesInput.SetHeight(8)
	notesInput.Blur()

	h := help.New()
	h.Width = 80

	return model{
		state:          stateTimer,
		stopwatch:      sw,
		spinner:        s,
		ratingInput:    ratingInput,
		notesInput:     notesInput,
		help:           h,
		keyMap:         DefaultKeyMap,
		startTime:      time.Now(),
		dailyNotesPath: dailyNotesPath,
		dateFormat:     dateFormat,
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
		case stateTimer:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.StopTimer):
				m.state = stateRating
				m.duration = m.stopwatch.Elapsed()
				m.stopwatch.Stop()
				m.ratingInput.Focus()
				return m, nil
			}

		case stateRating:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Continue):
				m.rating = m.ratingInput.Value()
				if m.rating == "" {
					m.rating = "0"
				}
				m.state = stateNotes
				m.ratingInput.Blur()
				m.notesInput.Focus()
				return m, nil
			case key.Matches(msg, m.keyMap.Skip):
				m.state = stateNotes
				m.notesInput.Focus()
				return m, nil
			}

		case stateNotes:
			switch {
			case key.Matches(msg, m.keyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Save):
				m.notes = m.notesInput.Value()
				m.state = stateSaving
				return m, m.saveSession()
			case key.Matches(msg, m.keyMap.Back):
				m.state = stateRating
				m.ratingInput.Focus()
				m.notesInput.Blur()
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
	case stateTimer:
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		cmds = append(cmds, cmd)
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case stateRating:
		m.ratingInput, cmd = m.ratingInput.Update(msg)
		cmds = append(cmds, cmd)

	case stateNotes:
		m.notesInput, cmd = m.notesInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	switch m.state {
	case stateTimer:
		elapsed := m.stopwatch.Elapsed()
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60
		hours := int(elapsed.Hours())

		timerDisplay := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		if hours == 0 {
			timerDisplay = fmt.Sprintf("%02d:%02d", minutes, seconds)
		}

		s += TitleStyle.Render("🧠 Deep Work Session")
		s += "\n\n"
		s += TimerStyle.Render(m.spinner.View() + " " + fmt.Sprintf("%s", timerDisplay))
		s += "\n\n"
		s += m.help.View(m.keyMap.TimerKeyMap())

	case stateRating:
		s += TitleStyle.Render("📊 Rate Your Session")
		s += "\n\n"
		s += "Rate this session (1-10):\n"
		if m.ratingInput.Focused() {
			s += FocusedStyle.Render(m.ratingInput.View())
		} else {
			s += InputStyle.Render(m.ratingInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.RatingKeyMap())

	case stateNotes:
		s += TitleStyle.Render("📝 Session Notes")
		s += "\n\n"
		s += "Notes (optional):\n"
		if m.notesInput.Focused() {
			s += FocusedStyle.Render(m.notesInput.View())
		} else {
			s += InputStyle.Render(m.notesInput.View())
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.NotesKeyMap())

	case stateSaving:
		s += TitleStyle.Render("💾 Saving Session...")
		s += "\n\n"
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		} else {
			s += SuccessStyle.Render("Session saved successfully!")
		}
		s += "\n\n"
		s += m.help.View(m.keyMap.SavingKeyMap())

	case stateDone:
		s += TitleStyle.Render("✅ Session Complete")
		s += "\n\n"
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error saving session: %v", m.err))
		} else {
			s += SuccessStyle.Render(fmt.Sprintf("Session logged to: %s", m.noteFilePath))
			s += "\n\n"
			minutes := int(m.duration.Minutes())
			seconds := int(m.duration.Seconds()) % 60
			s += fmt.Sprintf("Duration: %d minutes %d seconds\n", minutes, seconds)
			s += fmt.Sprintf("Rating: %s/10\n", m.rating)
		}
		s += "\n"
		s += m.help.View(m.keyMap.DoneKeyMap())
	}

	return s
}
