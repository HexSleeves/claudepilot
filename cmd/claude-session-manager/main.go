package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-session-manager",
	Short: "ClaudePilot - Terminal-based TUI for managing multiple Claude AI sessions",
	Long: `ClaudePilot is a terminal-based TUI application that manages and orchestrates 
multiple Claude AI sessions simultaneously. It enables developers to spawn, 
monitor, and facilitate complex interactions between multiple AI sessions 
from a single terminal interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := startTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ClaudePilot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ClaudePilot v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type model struct {
	quitting bool
	width    int
	height   int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Thanks for using ClaudePilot!\n"
	}

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 1)

	welcomeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Margin(1, 0)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Margin(1, 0, 0, 0)

	s.WriteString(titleStyle.Render("ClaudePilot - Claude Session Manager"))
	s.WriteString("\n\n")
	s.WriteString(welcomeStyle.Render("Welcome to ClaudePilot TUI!"))
	s.WriteString("\n")
	s.WriteString(welcomeStyle.Render("This application manages and orchestrates multiple Claude AI sessions."))
	s.WriteString("\n\n")
	s.WriteString("┌─ Session List ─────────────┐  ┌─ Output Display ──────────────────┐\n")
	s.WriteString("│                            │  │                                   │\n")
	s.WriteString("│  No sessions running       │  │  Welcome to ClaudePilot!          │\n")
	s.WriteString("│                            │  │                                   │\n")
	s.WriteString("│  Press 'n' to create new  │  │  Your Claude sessions will        │\n")
	s.WriteString("│  session                   │  │  appear in the left panel.       │\n")
	s.WriteString("│                            │  │                                   │\n")
	s.WriteString("│                            │  │  Output from selected session    │\n")
	s.WriteString("│                            │  │  will be displayed here.          │\n")
	s.WriteString("│                            │  │                                   │\n")
	s.WriteString("└────────────────────────────┘  └───────────────────────────────────┘\n")
	s.WriteString("                                ┌─ Input ───────────────────────────┐\n")
	s.WriteString("                                │                                   │\n")
	s.WriteString("                                │  Enter commands here...           │\n")
	s.WriteString("                                │                                   │\n")
	s.WriteString("                                └───────────────────────────────────┘\n")
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Controls: Ctrl+C or 'q' to quit  |  '?' for help  |  Tab to cycle panes"))

	return s.String()
}

func startTUI() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
