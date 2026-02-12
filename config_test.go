package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.TargetPorts) == 0 {
		t.Error("DefaultConfig() should have at least one target port")
	}

	if len(config.TargetProcesses) == 0 {
		t.Error("DefaultConfig() should have at least one target process")
	}
}

func TestLoadConfig_NoConfigFile(t *testing.T) {
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	config := LoadConfig()

	// Should return defaults when no config file exists
	if len(config.TargetPorts) == 0 || len(config.TargetProcesses) == 0 {
		t.Error("LoadConfig() should return defaults when no config file exists")
	}
}

func TestLoadConfig_ValidYAML(t *testing.T) {
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	configDir := filepath.Join(tempDir, "sharp-objects")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")
	configContent := `processes:
  - node
  - python
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config := LoadConfig()

	// Verify custom processes loaded
	if len(config.TargetProcesses) != 2 || config.TargetProcesses[0] != "node" || config.TargetProcesses[1] != "python" {
		t.Errorf("LoadConfig() failed to load custom processes correctly")
	}
}
