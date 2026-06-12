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

	"altum/internal/db"
	session "altum/internal/session"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a deep work session",
	Long: `Start a session for a deep work session. The session will run until you press Enter.
After stopping, you'll be prompted for a rating, interruptions, reflection and notes about the session.`,
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

		m := session.InitialModel(dbPath)
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
