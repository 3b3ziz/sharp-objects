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
	Command    string
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
func GetProcesses(config Config, includeSystem bool) ([]Process, error) {
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

	return parseProcesses(string(output), config, includeSystem)
}

// parseProcesses parses lsof output and filters for target processes/ports
func parseProcesses(output string, config Config, includeSystem bool) ([]Process, error) {
	lines := strings.Split(output, "\n")
	processes := make([]Process, 0)
	seen := make(map[string]bool) // Deduplicate by PID+Port

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

		// Get command line and working directory for the process
		command := getCommandLine(pid)
		workingDir := getWorkingDir(pid)

		// Filter system processes if requested
		if !includeSystem && isSystemProcess(processName, command) {
			continue
		}

		processes = append(processes, Process{
			Port:       port,
			PID:        pid,
			Name:       processName,
			Command:    command,
			WorkingDir: workingDir,
		})
	}

	return processes, nil
}

// getCommandLine retrieves the full command line of a process
func getCommandLine(pid int) string {
	// Use ps to get the full command line
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	cmdLine := strings.TrimSpace(string(output))

	// Clean up the command line for better readability
	cmdLine = cleanCommandLine(cmdLine)

	// Truncate very long commands
	if len(cmdLine) > 100 {
		cmdLine = cmdLine[:97] + "..."
	}

	return cmdLine
}

// isSystemProcess determines if a process is a system service (not dev-related)
func isSystemProcess(processName, command string) bool {
	// Common system process names to filter out
	systemProcesses := []string{
		"rapportd",
		"ControlCe",
		"ControlCenter",
		"NotificationCenter",
		"systemstats",
		"sharingd",
		"cloudd",
		"bird",
		"trustd",
		"UserEventAgent",
		"WiFiAgent",
	}

	processNameLower := strings.ToLower(processName)
	for _, sysProc := range systemProcesses {
		if strings.ToLower(sysProc) == processNameLower {
			return true
		}
	}

	// Filter out most .app bundles except dev tools
	if strings.Contains(command, ".app/Contents/") {
		// Keep these dev tools
		devTools := []string{
			"Visual Studio Code",
			"VSCode",
			"Cursor",
			"Sublime",
			"Atom",
			"WebStorm",
			"IntelliJ",
			"PyCharm",
			"Terminal",
			"iTerm",
			"Warp",
			"Postman",
			"Docker",
		}

		commandLower := strings.ToLower(command)
		for _, tool := range devTools {
			if strings.Contains(commandLower, strings.ToLower(tool)) {
				return false // Keep dev tools
			}
		}

		return true // Filter out other .app bundles
	}

	return false
}

// cleanCommandLine simplifies command paths for better readability
func cleanCommandLine(cmdLine string) string {
	// Split into parts
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		return cmdLine
	}

	// Clean the executable path
	executable := parts[0]

	// Strip common prefixes for cleaner display
	cleanPrefixes := []string{
		"/opt/homebrew/opt/",
		"/opt/homebrew/bin/",
		"/opt/homebrew/Cellar/",
		"/usr/local/bin/",
		"/usr/bin/",
		"/bin/",
		"/System/Library/",
		"/Applications/",
	}

	for _, prefix := range cleanPrefixes {
		if strings.HasPrefix(executable, prefix) {
			// For homebrew paths like /opt/homebrew/opt/mongod-community/bin/mongod
			// Extract just the binary name
			executable = strings.TrimPrefix(executable, prefix)

			// If it's a homebrew formula path, get just the binary at the end
			if strings.Contains(executable, "/bin/") {
				pathParts := strings.Split(executable, "/")
				executable = pathParts[len(pathParts)-1]
			}
			break
		}
	}

	// For .app bundles, show just the app name
	if strings.Contains(executable, ".app/Contents/") {
		appParts := strings.Split(executable, ".app/")
		if len(appParts) > 0 {
			appName := appParts[0]
			// Get just the app name, not the full path
			if idx := strings.LastIndex(appName, "/"); idx != -1 {
				appName = appName[idx+1:]
			}
			executable = appName + ".app"
		}
	}

	// Rebuild command with cleaned executable and important args
	result := executable

	// Add important arguments (skip very verbose ones)
	for i := 1; i < len(parts); i++ {
		arg := parts[i]

		// Skip very long path arguments
		if len(arg) > 50 && strings.Contains(arg, "/") {
			continue
		}

		// Keep config flags, ports, and short arguments
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") ||
		   len(arg) < 30 || strings.Contains(arg, "port") {
			result += " " + arg
		}
	}

	return result
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
