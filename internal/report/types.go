/*
Copyright © 2025 Eden Phillips
*/
package report

import "time"

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
