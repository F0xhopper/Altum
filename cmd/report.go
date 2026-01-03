/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var daysFlag int

type Session struct {
	Date         time.Time
	Duration     time.Duration
	FocusQuality int
	Milestone    string
}

type DayStats struct {
	Date     time.Time
	Sessions int
	Duration time.Duration
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your deep work sessions",
	Long:  `Generate a report of your deep work sessions for the last N days. Shows statistics including total sessions, time spent, average ratings, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		dailyNotesFolderPath := viper.GetString("daily_notes_folder_path")
		dateFormat := viper.GetString("date_format")

		if dailyNotesFolderPath == "" {
			fmt.Fprintf(os.Stderr, "Error: daily_notes_folder_path is required. Please set it using:\n")
			fmt.Fprintf(os.Stderr, "  altum config set daily_notes_folder_path <folder_path>\n")
			fmt.Fprintf(os.Stderr, "  or use --daily_notes_folder_path flag\n")
			os.Exit(1)
		}

		if dateFormat == "" {
			dateFormat = "2006-01-02"
		}

		sessions, err := parseSessions(dailyNotesFolderPath, dateFormat, daysFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing sessions: %v\n", err)
			os.Exit(1)
		}

		if len(sessions) == 0 {
			fmt.Printf("No sessions found in the last %d days.\n", daysFlag)
			return
		}

		printReport(sessions, daysFlag)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().IntVarP(&daysFlag, "days", "d", 7, "Number of days to include in the report")
}

func parseSessions(dailyNotesPath, dateFormat string, days int) ([]Session, error) {
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

func printReport(sessions []Session, days int) {
	if len(sessions) == 0 {
		return
	}

	var totalDuration time.Duration
	var totalFocusQuality int
	var focusQualityCount int
	longestSession := sessions[0]

	dayStatsMap := make(map[string]*DayStats)

	for _, session := range sessions {
		totalDuration += session.Duration

		if session.FocusQuality > 0 {
			totalFocusQuality += session.FocusQuality
			focusQualityCount++
		}

		if session.Duration > longestSession.Duration {
			longestSession = session
		}

		dateStr := session.Date.Format("2006-01-02")
		if dayStatsMap[dateStr] == nil {
			dayStatsMap[dateStr] = &DayStats{
				Date: session.Date,
			}
		}
		dayStatsMap[dateStr].Sessions++
		dayStatsMap[dateStr].Duration += session.Duration
	}

	var dayStats []*DayStats
	for _, stats := range dayStatsMap {
		dayStats = append(dayStats, stats)
	}

	var bestDay *DayStats
	if len(dayStats) > 0 {
		bestDay = dayStats[0]
		for _, stats := range dayStats {
			if stats.Duration > bestDay.Duration {
				bestDay = stats
			}
		}
	}

	avgDuration := totalDuration / time.Duration(len(sessions))
	avgFocusQuality := 0.0
	if focusQualityCount > 0 {
		avgFocusQuality = float64(totalFocusQuality) / float64(focusQualityCount)
	}

	totalHours := totalDuration.Hours()
	totalMinutes := int(totalDuration.Minutes())
	avgMinutes := int(avgDuration.Minutes())

	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1)

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  Deep Work Report: %s - %s\n",
		startDate.Format("Jan 2, 2006"),
		now.Format("Jan 2, 2006"))
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("Total sessions: %d\n", len(sessions))
	fmt.Printf("Total deep work: %.1f hours (%d minutes)\n", totalHours, totalMinutes)
	fmt.Printf("Average session: %d minutes\n", avgMinutes)

	if focusQualityCount > 0 {
		fmt.Printf("Average rating: %.1f / 5\n", avgFocusQuality)
	}

	if bestDay != nil {
		bestHours := bestDay.Duration.Hours()
		fmt.Printf("Best day: %s – %.1f hours (%d sessions)\n",
			bestDay.Date.Format("Jan 2"),
			bestHours,
			bestDay.Sessions)
	}

	longestMinutes := int(longestSession.Duration.Minutes())
	fmt.Printf("Longest session: %d minutes (%s)\n",
		longestMinutes,
		longestSession.Date.Format("Jan 2"))

	daysWithWork := len(dayStats)
	fmt.Printf("Days with deep work: %d / %d (%.0f%%)\n",
		daysWithWork,
		days,
		float64(daysWithWork)/float64(days)*100)

	if focusQualityCount > 0 {
		fmt.Printf("Total rating points: %d\n", totalFocusQuality)
	}

	fmt.Println()
	fmt.Println("Top performing days:")

	for i := 0; i < len(dayStats)-1; i++ {
		for j := 0; j < len(dayStats)-i-1; j++ {
			if dayStats[j].Duration < dayStats[j+1].Duration {
				dayStats[j], dayStats[j+1] = dayStats[j+1], dayStats[j]
			}
		}
	}

	topDays := 5
	if len(dayStats) < topDays {
		topDays = len(dayStats)
	}

	for i := 0; i < topDays; i++ {
		stats := dayStats[i]
		hours := stats.Duration.Hours()
		fmt.Printf("%s: %.1fh (%s)\n",
			stats.Date.Format("Jan 2"),
			hours)
	}

	fmt.Println()
}
