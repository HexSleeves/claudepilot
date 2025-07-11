package main

import (
	"fmt"
	"os"

	"claude-session-manager/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
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

func startTUI() error {
	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
