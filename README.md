# CodexGigantus

**CodexGigantus** is a powerful and flexible tool for aggregating code and text files from multiple sources (filesystem, CSV/TSV, databases) into a single output. Perfect for preparing codebases for Large Language Model (LLM) analysis, documentation generation, or code review workflows.

## ✨ Features

### Multiple Interfaces
- **CLI Mode**: Traditional command-line interface for scripts and automation
- **Web GUI**: User-friendly browser-based interface for configuration management
- **Docker Support**: Containerized deployment with Docker Compose

### Multiple Data Sources
- **Filesystem**: Process files from local or mounted directories
- **CSV/TSV**: Load file contents from tabular data
- **Database**: Query files from PostgreSQL, MySQL, or SQLite databases

### Advanced Configuration
- **JSON/YAML Config Files**: Save and load configurations for different projects
- **Environment Variables**: Flexible deployment configuration via `.env` files
- **Runtime Configuration**: Adjust settings on-the-fly via Web GUI

### Filtering & Processing
- Include/exclude by file extension
- Ignore specific files and directories
- Recursive directory traversal
- Extract function signatures from Go files
- Debug mode for troubleshooting

## 🚀 Quick Start

### Using Pre-built Binaries

```bash
# Clone the repository
git clone https://github.com/baditaflorin/codexgigantus.git
cd codexgigantus

# Build the project
make build

# Run CLI mode
./codexgigantus-cli --dir ./src --output output.txt

# Run Web GUI
./codexgigantus-web
# Open http://localhost:8080 in your browser
```

### Using Docker

```bash
# Start web GUI and PostgreSQL database
docker-compose up -d

# Access web interface at http://localhost:8080
# Access database admin at http://localhost:8081
```

### Using Makefile

```bash
# Show all available commands
make help

# Build both CLI and Web binaries
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Run CLI
make run-cli

# Run Web GUI
make run-web

# Build and run with Docker
make docker-build
make docker-up
```

## 📖 Usage

### CLI Mode

#### Basic Usage

```bash
# Process current directory
./codexgigantus-cli

# Process specific directories
./codexgigantus-cli --dir ./src,./pkg,./cmd

# Filter by extension (include only)
./codexgigantus-cli --include-ext go,py,js

# Exclude extensions
./codexgigantus-cli --exclude-ext log,tmp,md

# Ignore specific directories
./codexgigantus-cli --ignore-dir node_modules,.git,vendor

# Save output to file
./codexgigantus-cli --save --output mycode.txt

# Show function signatures only (Go files)
./codexgigantus-cli --show-funcs --include-ext go

# Enable debug mode
./codexgigantus-cli --debug
```

#### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | Comma-separated list of directories | `.` |
| `--recursive` | Search directories recursively | `true` |
| `--ignore-file` | Comma-separated files to ignore | - |
| `--ignore-dir` | Comma-separated directories to ignore | - |
| `--ignore-ext` | Comma-separated extensions to exclude | - |
| `--include-ext` | Comma-separated extensions to include (whitelist) | - |
| `--save` | Save output to file | `false` |
| `--output` | Output filename | `output.txt` |
| `--show-size` | Display output size | `false` |
| `--show-funcs` | Show only function signatures (Go files) | `false` |
| `--debug` | Enable debug output | `false` |

#### Shell Completion

```bash
# Install completions for your shell
./codexgigantus-cli install-completion

# Restart your shell to activate
```

### Web GUI Mode

Start the web server:

```bash
./codexgigantus-web
```

Then open http://localhost:8080 in your browser.

#### Web GUI Features

- **Visual Configuration**: Point-and-click interface for all settings
- **Save/Load Configs**: Manage multiple configuration profiles
- **Test Database Connections**: Verify database settings before processing
- **Live Processing**: Process files and view results instantly
- **Multi-Source Support**: Switch between filesystem, CSV, and database sources

### Configuration Files

#### Filesystem Configuration (`configs/filesystem.json`)

```json
{
  "name": "Filesystem Example",
  "source_type": "filesystem",
  "directories": ["./src", "./pkg"],
  "recursive": true,
  "ignore_dirs": ["node_modules", ".git"],
  "include_extensions": ["go", "py", "js"],
  "output_file": "output.txt"
}
```

#### Database Configuration (`configs/database.json`)

```json
{
  "name": "PostgreSQL Example",
  "source_type": "database",
  "db_type": "postgres",
  "db_host": "localhost",
  "db_port": 5432,
  "db_name": "codex",
  "db_user": "codex",
  "db_password": "secret",
  "db_table_name": "code_files",
  "db_column_path": "file_path",
  "db_column_content": "content",
  "output_file": "db_output.txt"
}
```

#### CSV Configuration (`configs/csv.json`)

```json
{
  "name": "CSV Example",
  "source_type": "csv",
  "csv_file_path": "./data/files.csv",
  "csv_delimiter": ",",
  "csv_path_column": 0,
  "csv_content_column": 1,
  "csv_has_header": true,
  "output_file": "csv_output.txt"
}
```

## 🔧 Environment Variables

Create a `.env` file to customize settings:

```bash
# Application Mode
APP_MODE=web                  # cli or web
WEB_PORT=8080
WEB_HOST=0.0.0.0

# Database Configuration
DB_TYPE=postgres              # postgres, mysql, sqlite
DB_HOST=localhost
DB_PORT=5432
DB_NAME=codex
DB_USER=postgres
DB_PASSWORD=postgres

# Database Schema
DB_TABLE_NAME=code_files
DB_COLUMN_PATH=file_path
DB_COLUMN_CONTENT=content
DB_COLUMN_TYPE=file_type
DB_COLUMN_SIZE=file_size

# Processing Defaults
DEFAULT_RECURSIVE=true
DEFAULT_DEBUG=false
DEFAULT_OUTPUT_FILE=output.txt
MAX_FILE_SIZE=10485760        # 10MB

# Logging
LOG_LEVEL=info                # debug, info, warn, error
LOG_FORMAT=text               # text or json
```

## 🐳 Docker Deployment

### Docker Compose

The provided `docker-compose.yml` includes:
- CodexGigantus Web GUI
- PostgreSQL database with sample schema
- Adminer for database management

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Clean everything (including volumes)
docker-compose down -v
```

### Database Setup

The PostgreSQL database is automatically initialized with:
- `code_files` table with proper schema
- Sample data for testing
- Indexes for performance

Connect to the database:
- **Host**: localhost
- **Port**: 5432
- **Database**: codex
- **User**: codex
- **Password**: codex_password

Or use Adminer at http://localhost:8081

## 🏗️ Architecture

```
codexgigantus/
├── cmd/
│   ├── cli/              # CLI application entry point
│   └── web/              # Web GUI entry point
├── internal/
│   ├── completion/       # Shell completion logic
│   └── gui/              # Web handlers and templates
├── pkg/
│   ├── config/           # Legacy CLI configuration
│   ├── configfile/       # JSON/YAML config management
│   ├── env/              # Environment variable handling
│   ├── processor/        # File system processing
│   ├── sources/
│   │   ├── csv/          # CSV/TSV processor
│   │   └── database/     # Database connector
│   └── utils/            # Shared utilities
├── configs/              # Example configurations
├── .env.example          # Example environment file
├── Dockerfile            # Container definition
├── docker-compose.yml    # Service orchestration
├── Makefile              # Build automation
└── init-db.sql           # Database initialization
```

### Design Principles

- **DRY (Don't Repeat Yourself)**: Shared logic in reusable packages
- **SOLID Principles**: Clear separation of concerns
- **Small Files**: Each file has a focused responsibility
- **Package Organization**: Clear boundaries between modules
- **Testability**: Comprehensive test coverage

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
# Opens coverage.html in browser

# Run tests with race detector
make test-race

# Run only unit tests
make test-unit
```

Current test coverage:
- ✅ Environment configuration
- ✅ Config file management (JSON/YAML)
- ✅ CSV/TSV processing
- ✅ Database connectivity
- ✅ File system processing
- ✅ Utility functions

## 📦 Building from Source

### Prerequisites

- Go 1.22 or later
- Make (optional, for Makefile targets)
- Docker (optional, for containerized deployment)

### Build Commands

```bash
# Install dependencies
make deps

# Build CLI only
make build-cli

# Build Web GUI only
make build-web

# Build both
make build

# Install to GOPATH/bin
make install

# Clean build artifacts
make clean
```

## 🔍 Use Cases

### 1. LLM Code Analysis

```bash
# Aggregate all Go code for GPT analysis
./codexgigantus-cli \
  --include-ext go \
  --ignore-dir vendor,node_modules \
  --save --output codebase.txt
```

### 2. Documentation Generation

```bash
# Extract function signatures for API docs
./codexgigantus-cli \
  --show-funcs \
  --include-ext go \
  --save --output api_functions.txt
```

### 3. Code Review Preparation

Load code from database for review:

1. Start Web GUI: `./codexgigantus-web`
2. Configure database connection
3. Run custom query to filter files
4. Export formatted output

### 4. Multi-Repository Analysis

Create a CSV with file paths and contents from multiple repos, then process:

```csv
file_path,content
repo1/main.go,package main...
repo2/app.py,import sys...
```

```bash
./codexgigantus-cli --source csv --csv-file repos.csv
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI parsing
- Database support via Go standard library `database/sql`
- YAML parsing with [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)

## 📧 Contact

- **Author**: Florin Badita
- **GitHub**: [@baditaflorin](https://github.com/baditaflorin)
- **Repository**: [github.com/baditaflorin/codexgigantus](https://github.com/baditaflorin/codexgigantus)

## 🗺️ Roadmap

- [ ] Add support for MongoDB
- [ ] REST API mode
- [ ] GraphQL query support
- [ ] Real-time file watching
- [ ] Incremental processing
- [ ] Output format templates (Markdown, HTML)
- [ ] Code metrics and statistics
- [ ] Integration with popular IDEs

---

**⭐ If you find this tool useful, please consider starring the repository!**
