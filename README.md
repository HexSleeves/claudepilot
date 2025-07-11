# Claude Session Manager TUI

A command-line tool that provides a Text-based User Interface (TUI) for managing multiple Claude sessions.

## Features

- Spawn and manage multiple Claude sessions
- Real-time monitoring and interaction
- Session-to-session communication
- Intuitive TUI interface

## Development

### Prerequisites

- Go 1.21 or later

### Building

```bash
go build -o bin/claude-session-manager ./cmd/claude-session-manager
```

### Running

```bash
./bin/claude-session-manager
```

### Development Commands

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...
```