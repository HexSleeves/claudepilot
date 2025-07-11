.PHONY: build test lint fmt clean install-tools

# Build the application
build:
	go build -o bin/claude-session-manager ./cmd/claude-session-manager

# Run the application in development mode
run-dev:
	go run ./cmd/claude-session-manager/main.go

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -rf bin/

# Install development tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

# Run all checks
check: fmt lint test

# Development workflow
dev: fmt lint test build
