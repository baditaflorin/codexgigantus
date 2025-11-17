// Package env provides environment variable management and configuration loading
package env

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all environment-based configuration
type Config struct {
	// Application settings
	AppMode  string
	WebPort  int
	WebHost  string

	// Database settings
	DBType     string
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// Database schema
	DBTableName      string
	DBColumnPath     string
	DBColumnContent  string
	DBColumnType     string
	DBColumnSize     string

	// Processing defaults
	DefaultRecursive  bool
	DefaultDebug      bool
	DefaultOutputFile string
	DefaultShowSize   bool
	DefaultShowFuncs  bool

	// File processing
	MaxFileSize     int64
	DefaultEncoding string

	// Shell completion paths
	BashCompletionDir string
	BashRCPath        string
	ZshCompletionDir  string
	ZshRCPath         string
	FishCompletionDir string

	// Logging
	LogLevel  string
	LogFormat string

	// Security
	AllowedExtensions  []string
	MaxConcurrentFiles int
}

// Load loads environment configuration from .env file and environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := loadEnvFile(".env"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		AppMode:            getEnv("APP_MODE", "cli"),
		WebPort:            getEnvInt("WEB_PORT", 8080),
		WebHost:            getEnv("WEB_HOST", "0.0.0.0"),
		DBType:             getEnv("DB_TYPE", "postgres"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnvInt("DB_PORT", 5432),
		DBName:             getEnv("DB_NAME", "codex"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:          getEnv("DB_SSL_MODE", "disable"),
		DBTableName:        getEnv("DB_TABLE_NAME", "code_files"),
		DBColumnPath:       getEnv("DB_COLUMN_PATH", "file_path"),
		DBColumnContent:    getEnv("DB_COLUMN_CONTENT", "content"),
		DBColumnType:       getEnv("DB_COLUMN_TYPE", "file_type"),
		DBColumnSize:       getEnv("DB_COLUMN_SIZE", "file_size"),
		DefaultRecursive:   getEnvBool("DEFAULT_RECURSIVE", true),
		DefaultDebug:       getEnvBool("DEFAULT_DEBUG", false),
		DefaultOutputFile:  getEnv("DEFAULT_OUTPUT_FILE", "output.txt"),
		DefaultShowSize:    getEnvBool("DEFAULT_SHOW_SIZE", false),
		DefaultShowFuncs:   getEnvBool("DEFAULT_SHOW_FUNCS", false),
		MaxFileSize:        getEnvInt64("MAX_FILE_SIZE", 10485760),
		DefaultEncoding:    getEnv("DEFAULT_ENCODING", "utf-8"),
		BashCompletionDir:  getEnv("BASH_COMPLETION_DIR", "/etc/bash_completion.d"),
		BashRCPath:         getEnv("BASH_RC_PATH", "~/.bashrc"),
		ZshCompletionDir:   getEnv("ZSH_COMPLETION_DIR", "~/.zsh/completions"),
		ZshRCPath:          getEnv("ZSH_RC_PATH", "~/.zshrc"),
		FishCompletionDir:  getEnv("FISH_COMPLETION_DIR", "~/.config/fish/completions"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", "text"),
		AllowedExtensions:  getEnvSlice("ALLOWED_EXTENSIONS", []string{".go", ".py", ".js", ".java"}),
		MaxConcurrentFiles: getEnvInt("MAX_CONCURRENT_FILES", 100),
	}

	return cfg, nil
}

// loadEnvFile loads environment variables from a file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set in environment
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvInt64 gets an int64 environment variable with a default value
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getEnvSlice gets a comma-separated environment variable as a slice
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}
