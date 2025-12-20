/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dailyNotesFolder := viper.GetString("daily_notes_folder")
		dateFormat := viper.GetString("date_format")

		fmt.Println("Starting application with configuration:")
		fmt.Printf("  Daily Notes Folder: %s\n", dailyNotesFolder)
		fmt.Printf("  Date Format: %s\n", dateFormat)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)


}
