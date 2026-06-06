/*
Copyright © 2025 Eden Phillips
*/
package report

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseSessions(dailyNotesPath, dateFormat string, days int) ([]Session, error) {
	var sessions []Session
	now := time.Now()

	dateMap := make(map[string]bool)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format(dateFormat)
		dateMap[dateStr] = true
	}

	files, err := os.ReadDir(dailyNotesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read daily notes directory: %w", err)
	}

	sessionRe := regexp.MustCompile(`^#### Session \d+$`)
	durationRe := regexp.MustCompile(`^- Duration: (\d+) minutes (\d+) seconds$`)
	focusQualityRe := regexp.MustCompile(`^- Focus Quality: (\d+)/5$`)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filename := strings.TrimSuffix(file.Name(), ".md")

		fileDate, err := time.Parse(dateFormat, filename)
		if err != nil {
			continue
		}

		dateStr := fileDate.Format(dateFormat)
		if !dateMap[dateStr] {
			continue
		}

		filePath := filepath.Join(dailyNotesPath, file.Name())
		fileSessions, err := parseFileSessions(filePath, fileDate, sessionRe, durationRe, focusQualityRe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", file.Name(), err)
			continue
		}

		sessions = append(sessions, fileSessions...)
	}

	return sessions, nil
}

func parseFileSessions(filePath string, fileDate time.Time, sessionRe, durationRe, focusQualityRe *regexp.Regexp) ([]Session, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sessions []Session
	scanner := bufio.NewScanner(file)

	var currentSession *Session
	inSessionsSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "## Altum Work Sessions") {
			inSessionsSection = true
			continue
		}

		if !inSessionsSection {
			continue
		}

		if sessionRe.MatchString(line) {
			if currentSession != nil {
				sessions = append(sessions, *currentSession)
			}
			currentSession = &Session{
				Date: fileDate,
			}
			continue
		}

		if currentSession == nil {
			continue
		}

		if matches := durationRe.FindStringSubmatch(line); matches != nil {
			minutes, _ := strconv.Atoi(matches[1])
			seconds, _ := strconv.Atoi(matches[2])
			currentSession.Duration = time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
			continue
		}

		if matches := focusQualityRe.FindStringSubmatch(line); matches != nil {
			quality, _ := strconv.Atoi(matches[1])
			currentSession.FocusQuality = quality
			continue
		}

		if strings.HasPrefix(line, "- Milestone: ") {
			currentSession.Milestone = strings.TrimPrefix(line, "- Milestone: ")
			continue
		}
	}

	if currentSession != nil {
		sessions = append(sessions, *currentSession)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
