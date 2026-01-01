/*
Copyright © 2025 Eden Phillips
*/
package timer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type saveSuccessMsg struct {
	noteFilePath string
	sessionCount int
}

type saveErrorMsg struct {
	err error
}

func (m *model) saveSession() tea.Cmd {
	return func() tea.Msg {
		today := time.Now().Format(m.dateFormat)
		noteFileName := fmt.Sprintf("%s.md", today)
		noteFilePath := filepath.Join(m.dailyNotesPath, noteFileName)

		file, err := os.OpenFile(noteFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return saveErrorMsg{err: err}
		}
		defer file.Close()

		titleRe := regexp.MustCompile(`^## Altum Work Sessions$`)
		sessionTitleRe := regexp.MustCompile(`^#### Session \d+`)

		isTitleLineFound := false
		sessionCount := 1

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !isTitleLineFound && titleRe.MatchString(line) {
				isTitleLineFound = true
				continue
			}
			if isTitleLineFound && sessionTitleRe.MatchString(line) {
				sessionCount++
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			return saveErrorMsg{err: err}
		}

		entry := ""
		if !isTitleLineFound {
			entry += "\n## Altum Work Sessions\n"
		}
		sessionStartTime := m.startTime.Format("15:04:05")
		sessionEndTime := time.Now().Format("15:04:05")
		entry += fmt.Sprintf("\n#### Session %d\n", sessionCount)
		entry += fmt.Sprintf("- Time: %s - %s\n", sessionStartTime, sessionEndTime)
		minutes := int(m.duration.Minutes())
		seconds := int(m.duration.Seconds()) % 60
		entry += fmt.Sprintf("- Duration: %d minutes %d seconds\n", minutes, seconds)
		entry += fmt.Sprintf("- Rate: %s/10\n", m.rating)
		if m.notes != "" {
			entry += fmt.Sprintf("- Notes: %s\n", m.notes)
		}

		if _, err := file.WriteString(entry); err != nil {
			return saveErrorMsg{err: err}
		}

		return saveSuccessMsg{
			noteFilePath: noteFilePath,
			sessionCount: sessionCount,
		}
	}
}

func (m model) handleSaveSuccess(msg saveSuccessMsg) model {
	m.noteFilePath = msg.noteFilePath
	m.sessionCount = msg.sessionCount
	m.state = stateDone
	return m
}

func (m model) handleSaveError(msg saveErrorMsg) model {
	m.err = msg.err
	m.state = stateDone
	return m
}

