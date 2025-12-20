/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"fmt"
	"os"

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
		dailyNotesFolderPath := viper.GetString("daily_notes_folder_path")
		dateFormat := viper.GetString("date_format")
		
		if dailyNotesFolderPath == "" {
			fmt.Fprintf(os.Stderr, "Error: daily_notes_folder_path is required. Please set it using:\n")
			fmt.Fprintf(os.Stderr, "  altum config set daily_notes_folder_path <folder_path>\n")
			fmt.Fprintf(os.Stderr, "  or use --daily_notes_folder_path flag\n")
			os.Exit(1)
		}
		
		fmt.Println("Starting application with configuration:")
		fmt.Printf("  Daily Notes Folder: %s\n", dailyNotesFolderPath)
		fmt.Printf("  Date Format: %s\n", dateFormat)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)


}
