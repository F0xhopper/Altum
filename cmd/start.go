/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	session "altum/internal/tui/session"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a deep work session",
	Long: `Start a session for a deep work session. The session will run until you press Enter.
After stopping, you'll be prompted for a rating, interruptions, reflection and notes about the session.`,
	Run: func(cmd *cobra.Command, args []string) {
		dailyNotesFolderPath := viper.GetString("daily_notes_folder_path")
		dateFormat := viper.GetString("date_format")

		if dailyNotesFolderPath == "" {
			fmt.Fprintf(os.Stderr, "Error: daily_notes_folder_path is required. Please set it using:\n")
			fmt.Fprintf(os.Stderr, "  altum config set daily_notes_folder_path <folder_path>\n")
			fmt.Fprintf(os.Stderr, "  or use --daily_notes_folder_path flag\n")
			os.Exit(1)
		}

		m := session.InitialModel(dailyNotesFolderPath, dateFormat)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
