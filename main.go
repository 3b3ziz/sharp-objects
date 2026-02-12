package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("15"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

type model struct {
	processes          []Process
	filteredProcesses  []Process
	config             Config
	cursor             int
	pendingKill        int  // -1 means no pending kill, otherwise index of process to kill
	showSystemServices bool // Toggle to show/hide system services
	searchMode         bool
	searchQuery        string
	error              string
	successMsg         string
}

type refreshMsg struct{}
type killMsg struct {
	err error
}
type clearMsg struct{}

func initialModel() model {
	config := LoadConfig()
	showSystemServices := false // Default to hiding system services
	processes, _ := GetProcesses(config, showSystemServices)

	return model{
		processes:          processes,
		filteredProcesses:  processes,
		config:             config,
		cursor:             0,
		pendingKill:        -1,
		showSystemServices: showSystemServices,
		searchMode:         false,
		searchQuery:        "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// filterProcesses filters the process list based on the search query
func (m *model) filterProcesses() {
	if m.searchQuery == "" {
		m.filteredProcesses = m.processes
		return
	}

	query := strings.ToLower(m.searchQuery)
	filtered := make([]Process, 0)

	for _, p := range m.processes {
		// Search across port, PID, command, and working directory
		if strings.Contains(strings.ToLower(fmt.Sprintf("%d", p.Port)), query) ||
			strings.Contains(strings.ToLower(fmt.Sprintf("%d", p.PID)), query) ||
			strings.Contains(strings.ToLower(p.Command), query) ||
			strings.Contains(strings.ToLower(p.WorkingDir), query) ||
			strings.Contains(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}

	m.filteredProcesses = filtered
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			switch msg.String() {
			case "esc":
				// Exit search mode and clear search
				m.searchMode = false
				m.searchQuery = ""
				m.filterProcesses()
				m.cursor = 0
				m.pendingKill = -1
				return m, nil

			case "enter":
				// Exit search mode but keep the filter
				m.searchMode = false
				m.pendingKill = -1
				return m, nil

			case "backspace":
				// Delete last character
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.filterProcesses()
					m.cursor = 0
					m.pendingKill = -1
				}
				return m, nil

			default:
				// Add typed character to search query
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.filterProcesses()
					m.cursor = 0
					m.pendingKill = -1
				}
				return m, nil
			}
		}

		// Normal mode key handling
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "/":
			// Enter search mode
			m.searchMode = true
			m.searchQuery = ""
			m.pendingKill = -1
			return m, nil

		case "esc":
			// Clear search if not in search mode
			if m.searchQuery != "" {
				m.searchQuery = ""
				m.filterProcesses()
				m.cursor = 0
				m.pendingKill = -1
			}
			return m, nil

		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.pendingKill = -1 // Cancel pending kill on navigation
			}

		case "down":
			if m.cursor < len(m.filteredProcesses)-1 {
				m.cursor++
				m.pendingKill = -1 // Cancel pending kill on navigation
			}

		case "k":
			// Double-tap to kill
			if len(m.filteredProcesses) > 0 {
				if m.pendingKill == m.cursor {
					// Second press - kill the process
					pid := m.filteredProcesses[m.cursor].PID
					m.pendingKill = -1
					return m, func() tea.Msg {
						err := KillProcess(pid)
						return killMsg{err: err}
					}
				} else {
					// First press - mark for pending kill
					m.pendingKill = m.cursor
					m.error = ""
					m.successMsg = ""
				}
			}

		case "r":
			// Refresh process list
			m.pendingKill = -1 // Cancel pending kill on refresh
			return m, func() tea.Msg {
				return refreshMsg{}
			}

		case "f":
			// Toggle system services filter
			m.showSystemServices = !m.showSystemServices
			m.pendingKill = -1
			return m, func() tea.Msg {
				return refreshMsg{}
			}

		default:
			// Cancel pending kill on any other key
			m.pendingKill = -1
		}

	case refreshMsg:
		processes, _ := GetProcesses(m.config, m.showSystemServices)
		m.processes = processes
		m.filterProcesses()
		m.error = ""
		m.successMsg = ""
		m.pendingKill = -1
		// Adjust cursor if it's out of bounds
		if m.cursor >= len(m.filteredProcesses) {
			m.cursor = len(m.filteredProcesses) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}

	case killMsg:
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to kill process: %v", msg.err)
			m.successMsg = ""
		} else {
			m.error = ""
			m.successMsg = "I never loved you. I hope that brings you some comfort."
			// Refresh process list immediately after successful kill
			processes, _ := GetProcesses(m.config, m.showSystemServices)
			m.processes = processes
			m.filterProcesses()
			// Adjust cursor if it's out of bounds
			if m.cursor >= len(m.filteredProcesses) {
				m.cursor = len(m.filteredProcesses) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			// Set up message clearing
			return m, func() tea.Msg {
				time.Sleep(3 * time.Second)
				return clearMsg{}
			}
		}

	case clearMsg:
		m.successMsg = ""
		m.error = ""
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Sharp Objects"))
	b.WriteString("\n\n")

	// Show search bar
	if m.searchMode || m.searchQuery != "" {
		searchPrefix := "Search: "
		searchText := m.searchQuery
		if m.searchMode {
			searchText += "_" // Show cursor in search mode
		}
		resultCount := fmt.Sprintf(" (%d results)", len(m.filteredProcesses))
		b.WriteString(headerStyle.Render(searchPrefix + searchText + resultCount))
		b.WriteString("\n\n")
	}

	// Show error if any
	if m.error != "" {
		b.WriteString(errorStyle.Render(m.error))
		b.WriteString("\n\n")
	}

	// Show success message if any
	if m.successMsg != "" {
		b.WriteString(successStyle.Render(m.successMsg))
		b.WriteString("\n\n")
	}

	// Header
	header := fmt.Sprintf("%-8s %-10s %-50s %s", "PORT", "PID", "COMMAND", "WORKING DIR")
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 120))
	b.WriteString("\n")

	// Rows
	if len(m.filteredProcesses) == 0 {
		if m.searchQuery != "" {
			b.WriteString(helpStyle.Render("No processes match your search."))
		} else {
			b.WriteString(helpStyle.Render("No processes found on monitored ports."))
		}
		b.WriteString("\n\n")
	} else {
		for i, p := range m.filteredProcesses {
			// Truncate working dir if too long
			workingDir := p.WorkingDir
			if len(workingDir) > 30 {
				workingDir = "..." + workingDir[len(workingDir)-27:]
			}

			// Truncate command if too long
			command := p.Command
			if len(command) > 50 {
				command = command[:47] + "..."
			}
			if command == "" {
				command = p.Name // Fallback to process name if command is empty
			}

			row := fmt.Sprintf("%-8d %-10d %-50s %s", p.Port, p.PID, command, workingDir)

			// Add pending kill indicator
			if i == m.pendingKill {
				row += "  " + errorStyle.Render("(press k again to kill)")
			}

			if i == m.cursor {
				b.WriteString(selectedStyle.Render(row))
			} else {
				b.WriteString(normalStyle.Render(row))
			}
			b.WriteString("\n")
		}
	}

	// Help text
	b.WriteString("\n")
	filterStatus := "dev only"
	if m.showSystemServices {
		filterStatus = "all"
	}
	var helpText string
	if m.searchMode {
		helpText = "Type to search  •  Enter: confirm  •  Esc: clear search"
	} else {
		helpText = fmt.Sprintf("↑/↓: navigate  •  k: kill  •  /: search  •  r: refresh  •  f: filter (%s)  •  Esc: clear  •  q: quit", filterStatus)
	}
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
