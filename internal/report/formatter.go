/*
Copyright © 2025 Eden Phillips
*/
package report

import (
	"fmt"
	"strings"
	"time"
)

func PrintReport(sessions []Session, days int) {
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
		stars := int(hours * 5 / 8)
		if stars > 5 {
			stars = 5
		}
		starStr := strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars)
		fmt.Printf("%s: %.1fh (%s)\n",
			stats.Date.Format("Jan 2"),
			hours,
			starStr)
	}

	fmt.Println()
}
