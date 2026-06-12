/*
Copyright © 2025 Eden Phillips
*/
package report

import (
	"fmt"
	"time"

	"altum/internal/db"
)

// ParseSessions opens the database at dbPath and returns all sessions
// recorded within the last `days` days (inclusive of today).
func ParseSessions(dbPath string, days int) ([]Session, error) {
	conn, err := db.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer conn.Close()

	records, err := db.Query(conn, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}

	sessions := make([]Session, 0, len(records))
	for _, r := range records {
		date, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			continue
		}
		sessions = append(sessions, Session{
			Date:         date,
			Duration:     time.Duration(r.DurationSeconds) * time.Second,
			FocusQuality: r.FocusQuality,
			Milestone:    r.Milestone,
		})
	}
	return sessions, nil
}
