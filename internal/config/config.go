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
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ValidConfigKeys is the set of keys accepted by `altum config set`.
// db_path overrides the default SQLite database location.
var ValidConfigKeys = map[string]bool{
	"db_path": true,
}

func IsValidKey(key string) bool {
	return ValidConfigKeys[key]
}

func GetConfigDir() (string, error) {
	configHome := os.ExpandEnv("$HOME/.config")
	if configHome == "$HOME/.config" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "altum"), nil
}

func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func SetConfigValue(key, value string) error {
	if !IsValidKey(key) {
		return fmt.Errorf("invalid key '%s'. Valid keys are: db_path", key)
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configFile, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig()

	viper.Set(key, value)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func GetAllConfigValues() map[string]string {
	return map[string]string{
		"db_path": viper.GetString("db_path"),
	}
}
