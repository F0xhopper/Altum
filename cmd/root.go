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
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "altum",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		selected := runMenu()
		switch selected {
		case menuStart:
			startCmd.Run(startCmd, []string{})
		case menuConfig:
			configCmd.Run(configCmd, []string{})
		case menuExit:
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
