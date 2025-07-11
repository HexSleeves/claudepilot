package session

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Status int

const (
	StatusIdle Status = iota
	StatusRunning
	StatusConnecting
	StatusError
	StatusStopped
)

func (s Status) String() string {
	switch s {
	case StatusIdle:
		return "idle"
	case StatusRunning:
		return "running"
	case StatusConnecting:
		return "connecting"
	case StatusError:
		return "error"
	case StatusStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Session struct {
	ID          string
	Name        string
	Status      Status
	Output      []string
	LastMessage string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	mu          sync.RWMutex
}

func NewSession(name string) *Session {
	return &Session{
		ID:        generateID(),
		Name:      name,
		Status:    StatusIdle,
		Output:    make([]string, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Session) AddOutput(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Output = append(s.Output, text)
	s.LastMessage = text
	s.UpdatedAt = time.Now()
}

func (s *Session) GetOutput() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	output := make([]string, len(s.Output))
	copy(output, s.Output)
	return output
}

func (s *Session) SetStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = status
	s.UpdatedAt = time.Now()
}

func (s *Session) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status
}

func (s *Session) SendInput(input string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// For now, this is a placeholder that simulates Claude processing
	// In a real implementation, this would send the input to Claude API
	s.AddOutput("🤖 Processing your request...")
	
	// Simulate a Claude response
	response := fmt.Sprintf("Claude response to: %s", strings.ReplaceAll(strings.TrimSpace(input), "\n", " "))
	s.AddOutput(response)
	
	s.UpdatedAt = time.Now()
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

type Manager struct {
	sessions []*Session
	mu       sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		sessions: make([]*Session, 0),
	}
}

func (m *Manager) CreateSession(name string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := NewSession(name)
	m.sessions = append(m.sessions, session)
	return session
}

func (m *Manager) GetSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessions := make([]*Session, len(m.sessions))
	copy(sessions, m.sessions)
	return sessions
}

func (m *Manager) GetSession(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.sessions {
		if session.ID == id {
			return session
		}
	}
	return nil
}

func (m *Manager) RemoveSession(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, session := range m.sessions {
		if session.ID == id {
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
			return true
		}
	}
	return false
}
