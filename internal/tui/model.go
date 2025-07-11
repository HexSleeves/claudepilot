package tui

import (
	"fmt"
	"strings"

	"claude-session-manager/internal/session"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FocusedPane int

const (
	SessionListPane FocusedPane = iota
	OutputPane
	InputPane
)

type Model struct {
	// Window dimensions
	width  int
	height int

	// State
	focusedPane     FocusedPane
	sessionManager  *session.Manager
	selectedSession *session.Session
	sessionCursor   int
	showHelp        bool

	// Input handling
	inputValue   string
	inputHistory []string
	historyIndex int

	// Output scrolling
	outputScroll int

	// UI components
	styles *Styles

	// Panel boundaries for mouse interaction
	sessionListBounds struct{ x, y, width, height int }
	outputPaneBounds  struct{ x, y, width, height int }
	inputPaneBounds   struct{ x, y, width, height int }

	// Application state
	quitting bool
}

func NewModel() *Model {
	sessionManager := session.NewManager()

	// Create some demo sessions
	session1 := sessionManager.CreateSession("Main Session")
	session1.SetStatus(session.StatusRunning)
	session1.AddOutput("Welcome to ClaudePilot!")
	session1.AddOutput("This is your main Claude session.")
	session1.AddOutput("Type your commands in the input pane below.")

	session2 := sessionManager.CreateSession("Analysis Session")
	session2.SetStatus(session.StatusIdle)
	session2.AddOutput("Analysis session ready for data processing.")

	session3 := sessionManager.CreateSession("Debug Session")
	session3.SetStatus(session.StatusError)
	session3.AddOutput("Error: Connection failed to Claude API")
	session3.AddOutput("Retrying connection...")

	model := &Model{
		sessionManager:  sessionManager,
		selectedSession: session1,
		sessionCursor:   0,
		focusedPane:     SessionListPane,
		styles:          NewStyles(),
		inputHistory:    make([]string, 0),
		historyIndex:    -1,
	}

	// Initialize panel bounds for mouse interaction
	model.updatePanelBounds()

	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePanelBounds()
		return m, nil

	case tea.KeyMsg:
		if m.showHelp {
			return m.handleHelpKeys(msg)
		}
		return m.handleKeys(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m, nil
}

func (m *Model) handleHelpKeys(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "?", "esc", "q":
		m.showHelp = false
	}
	return m, nil
}

func (m *Model) handleKeys(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "?":
		m.showHelp = true
		return m, nil

	case "tab":
		m.focusedPane = (m.focusedPane + 1) % 3
		return m, nil

	case "shift+tab":
		if m.focusedPane == 0 {
			m.focusedPane = 2
		} else {
			m.focusedPane = m.focusedPane - 1
		}
		return m, nil
	}

	switch m.focusedPane {
	case SessionListPane:
		return m.handleSessionListKeys(msg)
	case OutputPane:
		return m.handleOutputKeys(msg)
	case InputPane:
		return m.handleInputKeys(msg)
	}

	return m, nil
}

func (m *Model) handleSessionListKeys(msg tea.KeyMsg) (*Model, tea.Cmd) {
	sessions := m.sessionManager.GetSessions()

	switch msg.String() {
	case "j", "down":
		if m.sessionCursor < len(sessions)-1 {
			m.sessionCursor++
			m.selectedSession = sessions[m.sessionCursor]
			m.outputScroll = 0
		}

	case "k", "up":
		if m.sessionCursor > 0 {
			m.sessionCursor--
			m.selectedSession = sessions[m.sessionCursor]
			m.outputScroll = 0
		}

	case "n":
		sessionName := fmt.Sprintf("Session %d", len(sessions)+1)
		newSession := m.sessionManager.CreateSession(sessionName)
		newSession.AddOutput(fmt.Sprintf("New session '%s' created", sessionName))
		m.selectedSession = newSession
		m.sessionCursor = len(sessions)
		m.outputScroll = 0

	case "d", "x":
		if len(sessions) > 0 && m.selectedSession != nil {
			m.sessionManager.RemoveSession(m.selectedSession.ID)
			sessions = m.sessionManager.GetSessions()
			if len(sessions) == 0 {
				m.selectedSession = nil
				m.sessionCursor = 0
			} else {
				if m.sessionCursor >= len(sessions) {
					m.sessionCursor = len(sessions) - 1
				}
				m.selectedSession = sessions[m.sessionCursor]
			}
			m.outputScroll = 0
		}

	case "s":
		if m.selectedSession != nil {
			currentStatus := m.selectedSession.GetStatus()
			if currentStatus == session.StatusRunning {
				m.selectedSession.SetStatus(session.StatusStopped)
				m.selectedSession.AddOutput("Session stopped by user")
			} else {
				m.selectedSession.SetStatus(session.StatusRunning)
				m.selectedSession.AddOutput("Session started")
			}
		}
	}

	return m, nil
}

func (m *Model) handleOutputKeys(msg tea.KeyMsg) (*Model, tea.Cmd) {
	if m.selectedSession == nil {
		return m, nil
	}

	output := m.selectedSession.GetOutput()
	maxScroll := len(output) - m.getOutputHeight() + 2
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch msg.String() {
	case "j", "down":
		if m.outputScroll < maxScroll {
			m.outputScroll++
		}

	case "k", "up":
		if m.outputScroll > 0 {
			m.outputScroll--
		}

	case "g":
		m.outputScroll = 0

	case "G":
		m.outputScroll = maxScroll
	}

	return m, nil
}

func (m *Model) handleInputKeys(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Enter creates a new line
		m.inputValue += "\n"

	case "ctrl+enter":
		// Ctrl+Enter submits the input
		if strings.TrimSpace(m.inputValue) != "" && m.selectedSession != nil {
			m.inputHistory = append(m.inputHistory, m.inputValue)
			m.historyIndex = -1

			// Display the user input with proper formatting
			inputLines := strings.Split(m.inputValue, "\n")
			for i, line := range inputLines {
				if i == 0 {
					m.selectedSession.AddOutput(fmt.Sprintf("> %s", line))
				} else {
					m.selectedSession.AddOutput(fmt.Sprintf("  %s", line))
				}
			}

			// Send to Claude session
			m.selectedSession.SendInput(m.inputValue)

			m.inputValue = ""
			m.outputScroll = len(m.selectedSession.GetOutput())
		}

	case "up":
		// Only navigate history if not at start of input or input is empty
		if len(m.inputHistory) > 0 && (m.inputValue == "" || !strings.Contains(m.inputValue, "\n")) {
			if m.historyIndex == -1 {
				m.historyIndex = len(m.inputHistory) - 1
			} else if m.historyIndex > 0 {
				m.historyIndex--
			}
			if m.historyIndex >= 0 && m.historyIndex < len(m.inputHistory) {
				m.inputValue = m.inputHistory[m.historyIndex]
			}
		}

	case "down":
		// Only navigate history if input doesn't contain newlines or is at end
		if m.historyIndex != -1 && !strings.Contains(m.inputValue, "\n") {
			if m.historyIndex < len(m.inputHistory)-1 {
				m.historyIndex++
				m.inputValue = m.inputHistory[m.historyIndex]
			} else {
				m.historyIndex = -1
				m.inputValue = ""
			}
		}

	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}

	case "ctrl+backspace":
		// Delete word backwards
		if len(m.inputValue) > 0 {
			words := strings.Fields(m.inputValue)
			if len(words) > 0 {
				lastWord := words[len(words)-1]
				m.inputValue = m.inputValue[:len(m.inputValue)-len(lastWord)]
				m.inputValue = strings.TrimRight(m.inputValue, " ")
			}
		}

	default:
		if len(msg.String()) == 1 {
			m.inputValue += msg.String()
		}
	}

	return m, nil
}

func (m *Model) updatePanelBounds() {
	if m.width < 60 || m.height < 15 {
		return
	}

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth - 2

	// Session list bounds (left panel)
	m.sessionListBounds.x = 0
	m.sessionListBounds.y = 2 // After title
	m.sessionListBounds.width = leftWidth
	m.sessionListBounds.height = m.height - 4 // Account for title and footer

	// Output pane bounds (top right)
	m.outputPaneBounds.x = leftWidth + 2
	m.outputPaneBounds.y = 2 // After title
	m.outputPaneBounds.width = rightWidth
	m.outputPaneBounds.height = m.height/2 - 2

	// Input pane bounds (bottom right)
	m.inputPaneBounds.x = leftWidth + 2
	m.inputPaneBounds.y = 2 + m.height/2 - 2
	m.inputPaneBounds.width = rightWidth
	m.inputPaneBounds.height = m.height - m.height/2 - 4
}

func (m *Model) handleMouse(msg tea.MouseMsg) (*Model, tea.Cmd) {
	if m.showHelp {
		return m, nil
	}

	switch msg.Type {
	case tea.MouseLeft:
		// Check which panel was clicked
		if m.isPointInBounds(msg.X, msg.Y, m.sessionListBounds) {
			m.focusedPane = SessionListPane
			// Handle session list clicks
			return m.handleSessionListClick(msg.X, msg.Y)
		} else if m.isPointInBounds(msg.X, msg.Y, m.outputPaneBounds) {
			m.focusedPane = OutputPane
		} else if m.isPointInBounds(msg.X, msg.Y, m.inputPaneBounds) {
			m.focusedPane = InputPane
		}

	case tea.MouseWheelUp:
		if m.isPointInBounds(msg.X, msg.Y, m.outputPaneBounds) && m.focusedPane == OutputPane {
			if m.outputScroll > 0 {
				m.outputScroll--
			}
		} else if m.isPointInBounds(msg.X, msg.Y, m.sessionListBounds) && m.focusedPane == SessionListPane {
			sessions := m.sessionManager.GetSessions()
			if m.sessionCursor > 0 {
				m.sessionCursor--
				m.selectedSession = sessions[m.sessionCursor]
				m.outputScroll = 0
			}
		}

	case tea.MouseWheelDown:
		if m.isPointInBounds(msg.X, msg.Y, m.outputPaneBounds) && m.focusedPane == OutputPane {
			if m.selectedSession != nil {
				output := m.selectedSession.GetOutput()
				maxScroll := len(output) - m.getOutputHeight() + 2
				if maxScroll < 0 {
					maxScroll = 0
				}
				if m.outputScroll < maxScroll {
					m.outputScroll++
				}
			}
		} else if m.isPointInBounds(msg.X, msg.Y, m.sessionListBounds) && m.focusedPane == SessionListPane {
			sessions := m.sessionManager.GetSessions()
			if m.sessionCursor < len(sessions)-1 {
				m.sessionCursor++
				m.selectedSession = sessions[m.sessionCursor]
				m.outputScroll = 0
			}
		}
	}

	return m, nil
}

func (m *Model) isPointInBounds(x, y int, bounds struct{ x, y, width, height int }) bool {
	return x >= bounds.x && x < bounds.x+bounds.width &&
		y >= bounds.y && y < bounds.y+bounds.height
}

func (m *Model) handleSessionListClick(x, y int) (*Model, tea.Cmd) {
	sessions := m.sessionManager.GetSessions()
	if len(sessions) == 0 {
		return m, nil
	}

	// Calculate which session was clicked
	// Account for title (1 line) and border
	relativeY := y - m.sessionListBounds.y - 2 // Subtract title and border

	// Each session takes up 2 lines (name + preview)
	sessionIndex := relativeY / 2

	if sessionIndex >= 0 && sessionIndex < len(sessions) {
		m.sessionCursor = sessionIndex
		m.selectedSession = sessions[sessionIndex]
		m.outputScroll = 0
	}

	return m, nil
}

func (m *Model) getOutputHeight() int {
	return m.height - 8 // Account for borders, input pane, and title
}

func (m *Model) View() string {
	if m.quitting {
		return m.styles.InfoText.Render("Thanks for using ClaudePilot! 👋")
	}

	if m.showHelp {
		return m.renderHelp()
	}

	return m.renderMain()
}

func (m *Model) renderMain() string {
	if m.width < 60 || m.height < 15 {
		return m.styles.ErrorText.Render("Terminal too small. Please resize to at least 60x15.")
	}

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth - 2

	sessionList := m.renderSessionList(leftWidth, m.height-4)
	outputPane := m.renderOutputPane(rightWidth, m.height/2-2)
	inputPane := m.renderInputPane(rightWidth, m.height-m.height/2-4)

	rightColumn := lipgloss.JoinVertical(lipgloss.Top, outputPane, inputPane)

	main := lipgloss.JoinHorizontal(lipgloss.Top, sessionList, rightColumn)

	title := m.styles.TitleStyle.Render("ClaudePilot - Claude Session Manager")
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Top, title, main, footer)
}

func (m *Model) renderSessionList(width, height int) string {
	sessions := m.sessionManager.GetSessions()

	var items []string
	for i, sess := range sessions {
		status := m.styles.StatusIndicator(sess.GetStatus().String())

		var style lipgloss.Style
		if i == m.sessionCursor && m.focusedPane == SessionListPane {
			style = m.styles.SessionActive
		} else {
			style = m.styles.SessionInactive
		}

		line := fmt.Sprintf("%s %s", status, sess.Name)
		if sess.LastMessage != "" && len(sess.LastMessage) > 20 {
			preview := sess.LastMessage[:17] + "..."
			line += fmt.Sprintf("\n  %s", m.styles.InfoText.Render(preview))
		}

		items = append(items, style.Render(line))
	}

	if len(items) == 0 {
		items = append(items, m.styles.InfoText.Render("No sessions. Press 'n' to create one."))
	}

	content := strings.Join(items, "\n")

	var borderStyle lipgloss.Style
	if m.focusedPane == SessionListPane {
		borderStyle = m.styles.ActiveBorder
	} else {
		borderStyle = m.styles.InactiveBorder
	}

	title := "Sessions"
	if m.focusedPane == SessionListPane {
		title = "● Sessions"
	}

	return borderStyle.
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.TitleStyle.Render(title),
			content,
		))
}

func (m *Model) renderOutputPane(width, height int) string {
	var content string

	if m.selectedSession == nil {
		content = m.styles.InfoText.Render("Select a session to view output")
	} else {
		output := m.selectedSession.GetOutput()

		startLine := m.outputScroll
		endLine := startLine + height - 3
		if endLine > len(output) {
			endLine = len(output)
		}

		var lines []string
		for i := startLine; i < endLine; i++ {
			lines = append(lines, output[i])
		}

		content = strings.Join(lines, "\n")
		if content == "" {
			content = m.styles.InfoText.Render("No output yet...")
		}
	}

	var borderStyle lipgloss.Style
	if m.focusedPane == OutputPane {
		borderStyle = m.styles.ActiveBorder
	} else {
		borderStyle = m.styles.InactiveBorder
	}

	title := "Output"
	if m.focusedPane == OutputPane {
		title = "● Output"
	}

	return borderStyle.
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.TitleStyle.Render(title),
			m.styles.OutputText.Render(content),
		))
}

func (m *Model) renderInputPane(width, height int) string {
	prompt := m.styles.InputPrompt.Render("➤ ")

	// Handle multiline input display
	inputText := m.inputValue
	if m.focusedPane == InputPane {
		inputText += "█" // Show cursor only when focused
	}

	// Split input into lines and format each line
	lines := strings.Split(inputText, "\n")
	var formattedLines []string

	for i, line := range lines {
		if i == 0 {
			// First line gets the prompt
			formattedLines = append(formattedLines, prompt+m.styles.InputField.Render(line))
		} else {
			// Subsequent lines are indented
			indentedPrompt := strings.Repeat(" ", len("➤ "))
			formattedLines = append(formattedLines, indentedPrompt+m.styles.InputField.Render(line))
		}
	}

	content := strings.Join(formattedLines, "\n")

	// Add scroll if content is too tall for the pane
	availableHeight := height - 3 // Account for title and borders
	if len(formattedLines) > availableHeight {
		// Show last lines that fit
		start := len(formattedLines) - availableHeight
		content = strings.Join(formattedLines[start:], "\n")
	}

	var borderStyle lipgloss.Style
	if m.focusedPane == InputPane {
		borderStyle = m.styles.ActiveBorder
	} else {
		borderStyle = m.styles.InactiveBorder
	}

	title := "Input"
	if m.focusedPane == InputPane {
		title = "● Input"
	}

	return borderStyle.
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.TitleStyle.Render(title),
			content,
		))
}

func (m *Model) renderFooter() string {
	keys := []string{
		"Tab: Switch panes",
		"Mouse: Click panels/scroll",
		"?: Help",
		"Ctrl+C: Quit",
	}

	if m.focusedPane == SessionListPane {
		keys = append(keys, "n: New", "d: Delete", "s: Start/Stop", "Click: Select session")
	} else if m.focusedPane == InputPane {
		keys = append(keys, "Enter: New line", "Ctrl+Enter: Send", "↑/↓: History")
	} else if m.focusedPane == OutputPane {
		keys = append(keys, "j/k: Scroll", "g/G: Top/Bottom", "Wheel: Scroll")
	}

	return m.styles.InfoText.Render(strings.Join(keys, "  |  "))
}

func (m *Model) renderHelp() string {
	help := []string{
		m.styles.HelpTitle.Render("ClaudePilot Help"),
		"",
		m.styles.HelpKey.Render("Global Keys:"),
		"  Tab / Shift+Tab    Switch between panes",
		"  ?                  Show/hide this help",
		"  Ctrl+C             Quit application",
		"",
		m.styles.HelpKey.Render("Mouse Controls:"),
		"  Click              Focus panel and select items",
		"  Scroll Wheel       Navigate lists and scroll output",
		"  Click sessions     Select different sessions",
		"",
		m.styles.HelpKey.Render("Session List (Left Pane):"),
		"  j / ↓              Move cursor down",
		"  k / ↑              Move cursor up",
		"  n                  Create new session",
		"  d / x              Delete selected session",
		"  s                  Start/stop selected session",
		"  Click session      Select session",
		"",
		m.styles.HelpKey.Render("Output Pane (Top Right):"),
		"  j / ↓              Scroll down",
		"  k / ↑              Scroll up",
		"  g                  Go to top",
		"  G                  Go to bottom",
		"  Scroll wheel       Scroll output",
		"",
		m.styles.HelpKey.Render("Input Pane (Bottom Right):"),
		"  Enter              Create new line",
		"  Ctrl+Enter         Send message to Claude",
		"  ↑ / ↓              Navigate command history",
		"  Backspace          Delete character",
		"  Ctrl+Backspace     Delete word backward",
		"",
		m.styles.InfoText.Render("Press '?' or 'Esc' to close this help"),
	}

	content := strings.Join(help, "\n")

	return m.styles.BorderStyle.
		Width(m.width - 4).
		Height(m.height - 4).
		Padding(2).
		Render(content)
}
