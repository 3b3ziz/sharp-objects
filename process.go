package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Process represents a running process with port information
type Process struct {
	Port       int
	PID        int
	Name       string
	WorkingDir string
}

// Config holds the application configuration
type Config struct {
	TargetPorts     []int
	TargetProcesses []string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		// Ports based on actual usage from shell history
		TargetPorts: []int{
			// Common JS dev servers
			3000, 3001, 3002,
			// Vite
			5173, 5174,
			// Astro
			4321,
			// General dev
			5000, 8000, 8001,
			// Cloudflare Wrangler
			8787, 8788,
			// Other common
			8080, 8081,
			// Node debugger
			9229,
			// PostgreSQL
			5432,
			// Hugo
			1313,
			// Other observed
			7878, 7272,
		},
		TargetProcesses: []string{"node", "bun", "deno", "wrangler"},
	}
}

// GetProcesses retrieves all processes listening on configured ports
func GetProcesses(config Config) ([]Process, error) {
	// Run lsof to get listening processes
	// -i TCP -s TCP:LISTEN gets listening TCP connections
	// -n prevents hostname resolution (faster)
	// -P prevents port name resolution (shows numbers)
	cmd := exec.Command("lsof", "-i", "TCP", "-s", "TCP:LISTEN", "-n", "-P")
	output, err := cmd.Output()
	if err != nil {
		// If lsof fails (e.g., no processes found), return empty list
		return []Process{}, nil
	}

	return parseProcesses(string(output), config)
}

// parseProcesses parses lsof output and filters for target processes/ports
func parseProcesses(output string, config Config) ([]Process, error) {
	lines := strings.Split(output, "\n")
	processes := make([]Process, 0)
	seen := make(map[string]bool) // Deduplicate by PID+Port

	// Create a map for quick port lookup
	portMap := make(map[int]bool)
	for _, port := range config.TargetPorts {
		portMap[port] = true
	}

	// Create a map for quick process name lookup
	processMap := make(map[string]bool)
	for _, name := range config.TargetProcesses {
		processMap[strings.ToLower(name)] = true
	}

	// Regex to extract port from format like *:3000 or 127.0.0.1:8080
	portRegex := regexp.MustCompile(`:(\d+)`)

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // Skip header and empty lines
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		processName := fields[0]
		pidStr := fields[1]
		address := fields[8] // Format: *:PORT or IP:PORT

		// Extract port from address
		matches := portRegex.FindStringSubmatch(address)
		if len(matches) < 2 {
			continue
		}

		port, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		// Filter: check if process name matches OR port matches
		processNameLower := strings.ToLower(processName)
		matchesProcess := processMap[processNameLower]
		matchesPort := portMap[port]

		if !matchesProcess && !matchesPort {
			continue
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Deduplicate by PID+Port
		key := fmt.Sprintf("%d-%d", pid, port)
		if seen[key] {
			continue
		}
		seen[key] = true

		// Get working directory for the process
		workingDir := getWorkingDir(pid)

		processes = append(processes, Process{
			Port:       port,
			PID:        pid,
			Name:       processName,
			WorkingDir: workingDir,
		})
	}

	return processes, nil
}

// getWorkingDir retrieves the working directory of a process
func getWorkingDir(pid int) string {
	// Use lsof to get the cwd (current working directory)
	cmd := exec.Command("lsof", "-a", "-d", "cwd", "-p", strconv.Itoa(pid), "-Fn")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Parse lsof output format (lines starting with 'n' contain the path)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "n") {
			path := strings.TrimPrefix(line, "n")

			// Shorten home directory to ~
			home, _ := os.UserHomeDir()
			if home != "" && strings.HasPrefix(path, home) {
				path = "~" + strings.TrimPrefix(path, home)
			}

			return path
		}
	}

	return "unknown"
}

// KillProcess terminates a process by PID
func KillProcess(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	return cmd.Run()
}
