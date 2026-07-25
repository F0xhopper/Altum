/*
Copyright © 2025 Eden Phillips
*/
package session

import (
	"time"

	"github.com/F0xhopper/Altum/internal/db"

	tea "github.com/charmbracelet/bubbletea"
)

type saveSuccessMsg struct {
	sessionCount int
}

type saveErrorMsg struct {
	err error
}

func (m *model) saveSession() tea.Cmd {
	return func() tea.Msg {
		conn, err := db.Open(m.dbPath)
		if err != nil {
			return saveErrorMsg{err: err}
		}
		defer conn.Close()

		today := time.Now().Format("2006-01-02")

		count, err := db.CountByDate(conn, today)
		if err != nil {
			return saveErrorMsg{err: err}
		}
		sessionCount := count + 1

		rec := db.SessionRecord{
			Date:            today,
			StartTime:       m.startTime.Format("15:04:05"),
			EndTime:         time.Now().Format("15:04:05"),
			DurationSeconds: int(m.duration.Seconds()),
			Milestone:       m.milestone,
			FocusQuality:    focusQualityInt(m.focusQuality),
			Interruptions:   m.interruptions,
			Reflection:      m.reflection,
		}

		if err := db.Save(conn, rec); err != nil {
			return saveErrorMsg{err: err}
		}

		return saveSuccessMsg{sessionCount: sessionCount}
	}
}

func focusQualityInt(s string) int {
	switch s {
	case "1":
		return 1
	case "2":
		return 2
	case "4":
		return 4
	case "5":
		return 5
	default:
		return 3
	}
}

func (m model) handleSaveSuccess(msg saveSuccessMsg) model {
	m.sessionCount = msg.sessionCount
	m.state = stateDone
	return m
}

func (m model) handleSaveError(msg saveErrorMsg) model {
	m.err = msg.err
	m.state = stateDone
	return m
}
