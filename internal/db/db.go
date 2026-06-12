/*
Copyright © 2025 Eden Phillips

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS sessions (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    date             TEXT    NOT NULL,
    start_time       TEXT    NOT NULL,
    end_time         TEXT    NOT NULL,
    duration_seconds INTEGER NOT NULL,
    milestone        TEXT    NOT NULL,
    focus_quality    INTEGER NOT NULL DEFAULT 3,
    interruptions    TEXT,
    reflection       TEXT,
    created_at       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
CREATE INDEX IF NOT EXISTS idx_sessions_date ON sessions(date);
`

// DefaultPath returns the platform-appropriate default database path.
//
// Linux/BSD: $XDG_DATA_HOME/altum/altum.db  (~/.local/share/altum/altum.db)
// macOS:     ~/Library/Application Support/altum/altum.db
// Windows:   %AppData%\altum\altum.db
func DefaultPath() (string, error) {
	dataDir := userDataDir()
	return filepath.Join(dataDir, "altum", "altum.db"), nil
}

// userDataDir returns the OS-appropriate user data directory.
func userDataDir() string {
	// XDG_DATA_HOME takes priority on Linux/BSD.
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".local", "share")
	}
	return filepath.Join(home, ".local", "share")
}

// Open opens (or creates) the SQLite database at path and applies the schema.
func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialise schema: %w", err)
	}
	return db, nil
}

// SessionRecord holds all fields for a single deep-work session.
type SessionRecord struct {
	Date            string
	StartTime       string
	EndTime         string
	DurationSeconds int
	Milestone       string
	FocusQuality    int
	Interruptions   string
	Reflection      string
}

// Save inserts a new session into the database.
func Save(db *sql.DB, s SessionRecord) error {
	_, err := db.Exec(`
		INSERT INTO sessions
		    (date, start_time, end_time, duration_seconds, milestone, focus_quality, interruptions, reflection)
		VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''))`,
		s.Date, s.StartTime, s.EndTime, s.DurationSeconds,
		s.Milestone, s.FocusQuality, s.Interruptions, s.Reflection,
	)
	return err
}

// CountByDate returns the number of sessions recorded on a given date (YYYY-MM-DD).
func CountByDate(db *sql.DB, date string) (int, error) {
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM sessions WHERE date = ?`, date).Scan(&n)
	return n, err
}

// Query returns sessions whose date falls within the last `days` days (inclusive of today).
func Query(db *sql.DB, days int) ([]SessionRecord, error) {
	cutoff := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")
	rows, err := db.Query(`
		SELECT date, start_time, end_time, duration_seconds, milestone,
		       focus_quality, COALESCE(interruptions, ''), COALESCE(reflection, '')
		FROM sessions
		WHERE date >= ?
		ORDER BY date ASC, start_time ASC`,
		cutoff,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []SessionRecord
	for rows.Next() {
		var s SessionRecord
		if err := rows.Scan(
			&s.Date, &s.StartTime, &s.EndTime, &s.DurationSeconds,
			&s.Milestone, &s.FocusQuality, &s.Interruptions, &s.Reflection,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}
