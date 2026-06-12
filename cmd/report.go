/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"altum/internal/db"
	"altum/internal/report"
)

var daysFlag int

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your deep work sessions",
	Long:  `Generate a report of your deep work sessions for the last N days. Shows statistics including total sessions, time spent, average ratings, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("db_path")
		if dbPath == "" {
			var err error
			dbPath, err = db.DefaultPath()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving database path: %v\n", err)
				os.Exit(1)
			}
		}

		sessions, err := report.ParseSessions(dbPath, daysFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading sessions: %v\n", err)
			os.Exit(1)
		}

		if len(sessions) == 0 {
			fmt.Printf("No sessions found in the last %d days.\n", daysFlag)
			return
		}

		report.PrintReport(sessions, daysFlag)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().IntVarP(&daysFlag, "days", "d", 7, "Number of days to include in the report")
}
