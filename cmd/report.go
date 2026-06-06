package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"altum/internal/report"
)

var daysFlag int

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

		sessions, err := report.ParseSessions(dailyNotesFolderPath, dateFormat, daysFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing sessions: %v\n", err)
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
