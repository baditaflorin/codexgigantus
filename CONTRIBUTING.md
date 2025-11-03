# Contributing to CodexGigantus

Thank you for your interest in contributing to CodexGigantus! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/codexgigantus.git
   cd codexgigantus
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/baditaflorin/codexgigantus.git
   ```

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git

### Building the Project

```bash
# Build the binary
go build -o CodexGigantus

# Or use the build script
./build.sh
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Project Structure

See [CODE_STRUCTURE.md](CODE_STRUCTURE.md) for detailed information about the project structure.

```
codexgigantus/
├── main.go                 # Application entry point
├── internal/               # Internal packages (not importable by other projects)
│   └── completion/        # Shell completion installation
├── pkg/                    # Public packages
│   ├── config/            # Configuration management
│   ├── processor/         # File processing logic
│   └── utils/             # Utility functions
└── *_test.go              # Test files
```

## Making Changes

### Branching Strategy

- Create a feature branch from `main`:
  ```bash
  git checkout -b feature/your-feature-name
  ```
- Use descriptive branch names:
  - `feature/` for new features
  - `fix/` for bug fixes
  - `docs/` for documentation changes
  - `refactor/` for code refactoring

### Writing Code

1. **Follow Go conventions**: Use `gofmt` to format your code
2. **Add documentation**: All exported functions, types, and packages should have godoc comments
3. **Write tests**: All new code should include unit tests
4. **Keep changes focused**: Each PR should address a single concern

### Example: Adding a New Feature

```go
// Package example demonstrates proper documentation
package example

// ExampleFunction does something useful.
// It takes a parameter and returns a result.
func ExampleFunction(param string) string {
    // Implementation
    return param
}
```

## Testing

### Writing Tests

- Place test files alongside the code they test (e.g., `config.go` → `config_test.go`)
- Use table-driven tests for multiple test cases
- Test both success and error cases
- Use descriptive test names

Example test structure:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "input1", "output1"},
        {"case 2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Something(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Running Specific Tests

```bash
# Run tests for a specific package
go test ./pkg/config

# Run a specific test
go test -run TestParseCommaSeparated ./pkg/config
```

## Submitting Changes

### Before Submitting

1. Ensure all tests pass:
   ```bash
   go test ./...
   ```

2. Format your code:
   ```bash
   go fmt ./...
   ```

3. Run linters (if available):
   ```bash
   go vet ./...
   ```

4. Update documentation if needed

### Pull Request Process

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request on GitHub

3. Fill out the PR template with:
   - Clear description of changes
   - Link to related issues
   - Testing performed
   - Screenshots (if UI changes)

4. Wait for review and address feedback

### Commit Messages

Write clear, descriptive commit messages:

```
Short summary (50 chars or less)

Detailed explanation of what changed and why.
Can span multiple lines.

Fixes #123
```

## Code Style

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Keep functions small and focused
- Use meaningful variable names
- Add comments for complex logic

### Documentation Style

- Use godoc format for all exported items
- First sentence should be a complete, concise summary
- Start with the name of the item being documented
- Example:
  ```go
  // ParseCommaSeparated splits a comma-separated string into a slice.
  // It trims whitespace from each element and returns an empty slice
  // if the input is empty.
  func ParseCommaSeparated(s string) []string
  ```

## Questions?

If you have questions or need help, please:
1. Check existing issues and documentation
2. Open a new issue with the "question" label
3. Be patient and respectful

Thank you for contributing to CodexGigantus! 🎉
