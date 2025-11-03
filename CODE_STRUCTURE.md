# Code Structure

This document describes the architecture and organization of the CodexGigantus codebase.

## Overview

CodexGigantus is organized following Go best practices with clear separation of concerns:

- **main.go**: Application entry point and CLI setup
- **internal/**: Private packages not intended for external use
- **pkg/**: Public packages that could be imported by other projects

## Directory Structure

```
codexgigantus/
├── main.go                          # Application entry point, CLI setup
├── build.sh                         # Build script
├── go.mod                           # Go module definition
├── go.sum                           # Go module checksums
├── .gitignore                       # Git ignore patterns
├── README.md                        # User documentation
├── CONTRIBUTING.md                  # Contributor guidelines
├── CODE_STRUCTURE.md               # This file
├── internal/                        # Internal packages
│   └── completion/                 # Shell completion functionality
│       └── completion.go           # Completion installation logic
└── pkg/                            # Public packages
    ├── config/                     # Configuration management
    │   ├── config.go              # Config types and parsing
    │   └── config_test.go         # Config tests
    ├── processor/                  # File processing
    │   ├── file_processor.go      # File traversal and filtering
    │   └── file_processor_test.go # Processor tests
    └── utils/                      # Utility functions
        ├── utils.go               # Output generation, file ops
        └── utils_test.go          # Utils tests
```

## Package Descriptions

### main (root)

**File**: `main.go`

The entry point of the application. Contains:
- Cobra command definitions (root command, completion, install-completion)
- Flag registration and parsing
- Application initialization and execution flow

**Key responsibilities**:
- Initialize the CLI framework
- Register command-line flags
- Wire together the various packages
- Handle top-level errors

### internal/completion

**Package**: `github.com/baditaflorin/codexgigantus/internal/completion`

Handles shell completion installation for bash, zsh, and fish shells.

**Key components**:
- `InstallCompletion()`: Detects shell and installs appropriate completions
- `installBashCompletion()`: Bash-specific installation
- `installZshCompletion()`: Zsh-specific installation
- `installFishCompletion()`: Fish-specific installation
- Helper functions for file system operations

**Why internal?**
This package is specific to CodexGigantus and not intended for reuse in other projects.

### pkg/config

**Package**: `github.com/baditaflorin/codexgigantus/pkg/config`

Manages application configuration and command-line flag parsing.

**Key types**:
- `Config`: Struct holding all configuration options
  - Directory settings (Dirs, IgnoreDirs)
  - File filtering (IgnoreFiles, IgnoreExts, IncludeExts)
  - Behavior flags (Recursive, Debug, Save, etc.)

**Key functions**:
- `ParseCommaSeparated()`: Utility for parsing comma-separated flag values

**Design decisions**:
- Centralized configuration makes it easy to pass settings through the application
- All fields are exported for easy testing and potential external use

### pkg/processor

**Package**: `github.com/baditaflorin/codexgigantus/pkg/processor`

Handles file system traversal and file processing logic.

**Key types**:
- (Uses `utils.FileResult` for return values)

**Key functions**:
- `ProcessFiles()`: Main entry point - walks directories and processes files
- `shouldIgnoreDir()`: Determines if a directory should be skipped
- `shouldIgnoreFile()`: Determines if a file should be skipped

**Design decisions**:
- Uses `filepath.Walk` for directory traversal
- Applies filters during traversal for efficiency
- Reads entire file contents into memory (suitable for code files)

**Dependencies**:
- `pkg/config`: For accessing configuration settings
- `pkg/utils`: For FileResult type and Debug function

### pkg/utils

**Package**: `github.com/baditaflorin/codexgigantus/pkg/utils`

Provides utility functions for output generation and file operations.

**Key types**:
- `FileResult`: Represents a processed file (Path + Content)

**Key functions**:
- `GenerateOutput()`: Formats file results into output string
- `SaveOutput()`: Writes output to a file
- `Debug()`: Prints debug messages
- `extractFunctions()`: Extracts function signatures from Go code
- `isGoFile()`: Checks if a file is a Go source file

**Design decisions**:
- FileResult is defined here to avoid circular dependencies
- Uses Go's AST parser for extracting function information
- Handles both regular content display and function-only mode

**Dependencies**:
- `pkg/config`: For accessing configuration settings
- Standard library packages for AST parsing

## Data Flow

```
1. main.go
   ↓ parses flags into
2. config.Config
   ↓ passed to
3. processor.ProcessFiles()
   ↓ returns
4. []utils.FileResult
   ↓ passed to
5. utils.GenerateOutput()
   ↓ returns formatted string
6. main.go (displays or saves)
```

## Design Principles

### Separation of Concerns

Each package has a single, well-defined responsibility:
- **config**: Configuration management
- **processor**: File system operations
- **utils**: Output formatting and helpers
- **completion**: Shell integration

### Testability

- All packages are independently testable
- Functions accept configuration as parameters (dependency injection)
- No global state (except CLI flags in main)
- Table-driven tests for comprehensive coverage

### Modularity

- Packages can be used independently
- Clear interfaces between components
- Minimal coupling between packages

### Go Best Practices

- Exported items have godoc comments
- Error handling at all levels
- Use of standard library where possible
- idiomatic Go naming and structure

## Adding New Features

### Adding a New Filter Type

1. Add field to `config.Config`
2. Register flag in `main.go`
3. Implement filtering logic in `processor.shouldIgnoreFile()` or `processor.shouldIgnoreDir()`
4. Add tests in `pkg/processor/file_processor_test.go`

### Adding a New Output Format

1. Add flag to `config.Config`
2. Implement formatting logic in `utils.GenerateOutput()`
3. Add tests in `pkg/utils/utils_test.go`

### Adding a New Command

1. Create command in `main.go` using Cobra
2. Implement logic (possibly in new package)
3. Add to `init()` function with `rootCmd.AddCommand()`

## Testing Strategy

### Unit Tests

Each package has comprehensive unit tests:
- `config_test.go`: Tests configuration parsing
- `file_processor_test.go`: Tests file filtering and processing
- `utils_test.go`: Tests output generation and utilities

### Test Coverage

Run tests with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Target: >80% code coverage for all packages

### Integration Testing

Manual testing checklist:
1. Process directory with various filters
2. Save output to file
3. Test recursive vs non-recursive
4. Test function extraction mode
5. Test completion installation

## Future Improvements

Potential areas for enhancement:

1. **Performance**
   - Parallel file processing
   - Streaming output for large codebases

2. **Features**
   - Additional output formats (JSON, XML)
   - File content filtering (e.g., by regex)
   - Git integration (process only tracked files)

3. **Architecture**
   - Plugin system for custom filters
   - Configuration file support
   - Progress reporting for large operations

## Questions?

See [CONTRIBUTING.md](CONTRIBUTING.md) for information on contributing to the project.
