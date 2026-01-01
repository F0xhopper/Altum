/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage configuration settings for Altum.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config command used")
	},
}
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: daily_notes_folder_path, date_format`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		validKeys := map[string]bool{
			"daily_notes_folder_path": true,
			"date_format":             true,
		}
		if !validKeys[key] {
			fmt.Fprintf(os.Stderr, "Error: Invalid key '%s'. Valid keys are: daily_notes_folder_path, date_format\n", key)
			os.Exit(1)
		}

		configHome := os.ExpandEnv("$HOME/.config")
		if configHome == "$HOME/.config" {
			home, _ := os.UserHomeDir()
			configHome = filepath.Join(home, ".config")
		}
		altumConfigDir := filepath.Join(configHome, "altum")
		configFile := filepath.Join(altumConfigDir, "config.yaml")

		if err := os.MkdirAll(altumConfigDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create config directory: %v\n", err)
			os.Exit(1)
		}

		viper.SetConfigFile(configFile)
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
		}

		viper.Set(key, value)

		if err := viper.WriteConfigAs(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		fmt.Printf("Configuration saved to: %s\n", configFile)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value. Shows all values if no key is provided.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Current configuration:")
			fmt.Printf("  daily_notes_folder_path: %s\n", viper.GetString("daily_notes_folder_path"))
			fmt.Printf("  date_format: %s\n", viper.GetString("date_format"))
		} else {
			key := args[0]
			value := viper.GetString(key)
			if value == "" {
				fmt.Printf("%s is not set\n", key)
			} else {
				fmt.Printf("%s = %s\n", key, value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}
