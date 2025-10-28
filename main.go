package main

import (
	"fmt"
	"os"
	"strings"

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

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Padding(1, 0)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))
)

type model struct {
	processes      []Process
	config         Config
	cursor         int
	confirmingKill bool
	error          string
}

type refreshMsg struct{}
type killMsg struct {
	err error
}

func initialModel() model {
	config := LoadConfig()
	processes, _ := GetProcesses(config)

	return model{
		processes:      processes,
		config:         config,
		cursor:         0,
		confirmingKill: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If confirming kill
		if m.confirmingKill {
			switch msg.String() {
			case "y", "Y":
				if m.cursor < len(m.processes) {
					pid := m.processes[m.cursor].PID
					return m, func() tea.Msg {
						err := KillProcess(pid)
						return killMsg{err: err}
					}
				}
				m.confirmingKill = false
			case "n", "N", "esc":
				m.confirmingKill = false
			}
			return m, nil
		}

		// Normal key handling
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.processes)-1 {
				m.cursor++
			}

		case "K":
			// Kill with confirmation
			if len(m.processes) > 0 {
				m.confirmingKill = true
				m.error = ""
			}

		case "r":
			// Refresh process list
			return m, func() tea.Msg {
				return refreshMsg{}
			}
		}

	case refreshMsg:
		processes, _ := GetProcesses(m.config)
		m.processes = processes
		m.error = ""
		// Adjust cursor if it's out of bounds
		if m.cursor >= len(m.processes) {
			m.cursor = len(m.processes) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}

	case killMsg:
		m.confirmingKill = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to kill process: %v", msg.err)
		} else {
			m.error = ""
			// Refresh after kill
			return m, func() tea.Msg {
				return refreshMsg{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🔌 Dev Process Monitor"))
	b.WriteString("\n\n")

	// Show confirmation dialog
	if m.confirmingKill && len(m.processes) > 0 {
		selected := m.processes[m.cursor]
		b.WriteString(confirmStyle.Render(
			fmt.Sprintf("Kill process %s (PID: %d) on port %d? (y/n)",
				selected.Name, selected.PID, selected.Port),
		))
		return b.String()
	}

	// Show error if any
	if m.error != "" {
		b.WriteString(errorStyle.Render(m.error))
		b.WriteString("\n\n")
	}

	// Header
	header := fmt.Sprintf("%-8s %-20s %-10s %s", "PORT", "PROCESS", "PID", "WORKING DIR")
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 80))
	b.WriteString("\n")

	// Rows
	if len(m.processes) == 0 {
		b.WriteString(helpStyle.Render("No processes found on monitored ports."))
		b.WriteString("\n\n")
	} else {
		for i, p := range m.processes {
			// Truncate working dir if too long
			workingDir := p.WorkingDir
			if len(workingDir) > 35 {
				workingDir = "..." + workingDir[len(workingDir)-32:]
			}

			row := fmt.Sprintf("%-8d %-20s %-10d %s", p.Port, p.Name, p.PID, workingDir)

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
	b.WriteString(helpStyle.Render("↑/↓ or j/k: navigate  •  K: kill  •  r: refresh  •  q: quit"))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
