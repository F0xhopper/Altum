package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"altum/internal/config"
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
	Long:  `Set a configuration value. Available keys: db_path`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := config.SetConfigValue(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		configFile, err := config.GetConfigFilePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
			allValues := config.GetAllConfigValues()
			for key, value := range allValues {
				fmt.Printf("  %s: %s\n", key, value)
			}
		} else {
			key := args[0]
			value := config.GetConfigValue(key)
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
