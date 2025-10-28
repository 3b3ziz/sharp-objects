package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigFile represents the YAML config file structure
type ConfigFile struct {
	Ports     []int    `yaml:"ports"`
	Processes []string `yaml:"processes"`
}

// LoadConfig loads configuration from file or returns defaults
func LoadConfig() Config {
	config := DefaultConfig()

	// Try to load from config file
	configPath := getConfigPath()
	if configPath == "" {
		return config
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file doesn't exist or can't be read, use defaults
		return config
	}

	var configFile ConfigFile
	err = yaml.Unmarshal(data, &configFile)
	if err != nil {
		// Invalid YAML, use defaults
		return config
	}

	// Override defaults with config file values
	if len(configFile.Ports) > 0 {
		config.TargetPorts = configFile.Ports
	}
	if len(configFile.Processes) > 0 {
		config.TargetProcesses = configFile.Processes
	}

	return config
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	// Try XDG_CONFIG_HOME first
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		// Fall back to ~/.config
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(home, ".config")
	}

	configPath := filepath.Join(configDir, "sharp-objects", "config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		return ""
	}

	return configPath
}
