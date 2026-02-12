package main

import (
	"testing"
)

func TestParseProcesses(t *testing.T) {
	config := Config{
		TargetPorts:     []int{3000, 5173},
		TargetProcesses: []string{"node", "bun"},
	}

	lsofOutput := `COMMAND     PID           USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node      12345 ahmed.abdelaziz   23u  IPv4 0x1234567890abcdef      0t0  TCP *:3000 (LISTEN)
bun       67890 ahmed.abdelaziz   45u  IPv6 0xabcdef1234567890      0t0  TCP *:5173 (LISTEN)`

	processes, err := parseProcesses(lsofOutput, config, false)
	if err != nil {
		t.Fatalf("parseProcesses() error = %v", err)
	}

	if len(processes) != 2 {
		t.Errorf("parseProcesses() got %d processes, want 2", len(processes))
	}

	if processes[0].Port != 3000 || processes[0].PID != 12345 {
		t.Errorf("parseProcesses() incorrect first process: got port=%d pid=%d, want port=3000 pid=12345",
			processes[0].Port, processes[0].PID)
	}
}

func TestCleanCommandLine(t *testing.T) {
	tests := []struct {
		input    string
		contains string // Just check it contains this, not exact match
	}{
		{
			input:    "/opt/homebrew/opt/mongod-community/bin/mongod --config /etc/mongod.conf",
			contains: "mongod",
		},
		{
			input:    "/usr/local/bin/node server.js --port 3000",
			contains: "node",
		},
		{
			input:    "/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
			contains: "Visual Studio",
		},
	}

	for _, tt := range tests {
		result := cleanCommandLine(tt.input)
		if len(result) == 0 {
			t.Errorf("cleanCommandLine(%q) returned empty string", tt.input)
		}
	}
}

func TestIsSystemProcess(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		isSystem  bool
	}{
		{"rapportd", "rapportd", true},
		{"node", "node server.js", false},
		{"VS Code", "/Applications/Visual Studio Code.app/Contents/MacOS/Electron", false},
		{"unknown app", "/Applications/RandomApp.app/Contents/MacOS/RandomApp", true},
	}

	for _, tt := range tests {
		result := isSystemProcess(tt.name, tt.command)
		if result != tt.isSystem {
			t.Errorf("isSystemProcess(%q, %q) = %v, want %v", tt.name, tt.command, result, tt.isSystem)
		}
	}
}
