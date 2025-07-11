package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	// Base styles
	BorderStyle    lipgloss.Style
	TitleStyle     lipgloss.Style
	ActiveBorder   lipgloss.Style
	InactiveBorder lipgloss.Style

	// Session list styles
	SessionActive   lipgloss.Style
	SessionInactive lipgloss.Style
	SessionRunning  lipgloss.Style
	SessionError    lipgloss.Style
	SessionIdle     lipgloss.Style
	SessionStopped  lipgloss.Style

	// Output styles
	OutputText lipgloss.Style
	ErrorText  lipgloss.Style
	InfoText   lipgloss.Style

	// Input styles
	InputField  lipgloss.Style
	InputPrompt lipgloss.Style

	// Help styles
	HelpKey   lipgloss.Style
	HelpDesc  lipgloss.Style
	HelpTitle lipgloss.Style

	// Status indicators
	StatusRunning    lipgloss.Style
	StatusIdle       lipgloss.Style
	StatusError      lipgloss.Style
	StatusStopped    lipgloss.Style
	StatusConnecting lipgloss.Style
}

func NewStyles() *Styles {
	// Color palette - vibrant and modern
	primary := lipgloss.Color("#7C3AED")   // Vibrant purple
	secondary := lipgloss.Color("#06B6D4") // Cyan
	success := lipgloss.Color("#10B981")   // Green
	warning := lipgloss.Color("#F59E0B")   // Amber
	danger := lipgloss.Color("#EF4444")    // Red
	muted := lipgloss.Color("#6B7280")     // Gray
	surface := lipgloss.Color("#374151")   // Surface color
	text := lipgloss.Color("#F9FAFB")      // Light text

	return &Styles{
		// Base styles
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(muted),

		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(primary).
			Padding(0, 1),

		ActiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),

		InactiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(muted).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),

		// Session list styles
		SessionActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(text).
			Background(primary).
			Padding(0, 1).
			Margin(0, 0, 0, 1),

		SessionInactive: lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 1).
			Margin(0, 0, 0, 1),

		SessionRunning: lipgloss.NewStyle().
			Foreground(success).
			Bold(true),

		SessionError: lipgloss.NewStyle().
			Foreground(danger).
			Bold(true),

		SessionIdle: lipgloss.NewStyle().
			Foreground(muted),

		SessionStopped: lipgloss.NewStyle().
			Foreground(warning),

		// Output styles
		OutputText: lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 1),

		ErrorText: lipgloss.NewStyle().
			Foreground(danger).
			Bold(true),

		InfoText: lipgloss.NewStyle().
			Foreground(secondary).
			Italic(true),

		// Input styles
		InputField: lipgloss.NewStyle().
			Foreground(text).
			Background(surface).
			Padding(0, 1),

		InputPrompt: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true),

		// Help styles
		HelpKey: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Padding(0, 1),

		HelpDesc: lipgloss.NewStyle().
			Foreground(text),

		HelpTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(secondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(0, 2).
			Margin(1, 0),

		// Status indicators
		StatusRunning: lipgloss.NewStyle().
			Foreground(success).
			SetString("●"),

		StatusIdle: lipgloss.NewStyle().
			Foreground(muted).
			SetString("○"),

		StatusError: lipgloss.NewStyle().
			Foreground(danger).
			SetString("✗"),

		StatusStopped: lipgloss.NewStyle().
			Foreground(warning).
			SetString("■"),

		StatusConnecting: lipgloss.NewStyle().
			Foreground(secondary).
			SetString("⟳"),
	}
}

func (s *Styles) StatusIndicator(status string) string {
	switch status {
	case "running":
		return s.StatusRunning.Render()
	case "idle":
		return s.StatusIdle.Render()
	case "error":
		return s.StatusError.Render()
	case "stopped":
		return s.StatusStopped.Render()
	case "connecting":
		return s.StatusConnecting.Render()
	default:
		return s.StatusIdle.Render()
	}
}
