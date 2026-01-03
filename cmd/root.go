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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"altum/internal/tui/menu"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "altum",
	Short: "A minimalist CLI deep work companion for focused creators and knowledge workers",
	Long: `Altum is a minimalist CLI deep work companion for the terminal, built for focused 
creators and knowledge workers.

Track your deep work sessions with an elegant terminal interface. Altum helps you:
  • Time and monitor your focused work sessions
  • Capture milestones, reflections, and interruptions after each session
  • Automatically log sessions to your daily notes (Obsidian-compatible)
  • Build awareness of your deep work patterns and habits

The name "Altum" comes from the Latin word meaning "deep" — a fitting name for a tool 
designed to help you achieve deeper, more meaningful work.`,
	Run: func(cmd *cobra.Command, args []string) {
		selected := menu.RunMenu()
		switch selected {
		case menu.MenuStart:
			startCmd.Run(startCmd, []string{})
		case menu.MenuConfig:
			configCmd.Run(configCmd, []string{})
		case menu.MenuExit:
			os.Exit(0)
		}
	}}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.altum.yaml)")

	rootCmd.PersistentFlags().String("daily_notes_folder_path", "", "Path to the daily notes folder (required)")
	rootCmd.PersistentFlags().String("date_format", "2006-01-02", "Date format for notes (Obsidian format)")

	viper.BindPFlag("daily_notes_folder_path", rootCmd.PersistentFlags().Lookup("daily_notes_folder_path"))
	viper.BindPFlag("date_format", rootCmd.PersistentFlags().Lookup("date_format"))

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configHome := os.ExpandEnv("$HOME/.config")
		if configHome == "$HOME/.config" {
			home, _ := os.UserHomeDir()
			configHome = filepath.Join(home, ".config")
		}
		altumConfigDir := filepath.Join(configHome, "altum")
		viper.AddConfigPath(altumConfigDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("ALTUM")
	viper.AutomaticEnv()

	viper.SetDefault("date_format", "2006-01-02")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
