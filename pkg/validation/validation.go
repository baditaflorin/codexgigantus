// Package validation provides comprehensive input validation and sanitization
package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Security constants
const (
	MaxPathLength       = 4096
	MaxQueryLength      = 10000
	MaxConfigNameLength = 255
	MaxTableNameLength  = 128
	MaxColumnNameLength = 128
)

var (
	// SQL identifier pattern: alphanumeric, underscore, starting with letter
	sqlIdentifierPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

	// Path traversal patterns to block
	pathTraversalPatterns = []string{
		"..",
		"~",
		"$",
		"|",
		";",
		"&",
		"`",
		"<",
		">",
	}

	// SQL injection patterns to detect
	sqlInjectionPatterns = []string{
		"--",
		"/*",
		"*/",
		";",
		"'",
		"\"",
		"xp_",
		"sp_",
		"DROP ",
		"INSERT ",
		"UPDATE ",
		"DELETE ",
		"EXEC",
		"UNION",
	}
)

// Error types for different validation failures
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// ValidateSQLIdentifier validates SQL table/column names to prevent injection
func ValidateSQLIdentifier(name, fieldName string) error {
	if name == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	if len(name) > MaxTableNameLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxTableNameLength),
		}
	}

	// Check for SQL identifier pattern
	if !sqlIdentifierPattern.MatchString(name) {
		return &ValidationError{
			Field:   fieldName,
			Message: "must contain only alphanumeric characters and underscores, starting with a letter",
		}
	}

	// Check for SQL injection patterns
	upperName := strings.ToUpper(name)
	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(upperName, pattern) {
			return &ValidationError{
				Field:   fieldName,
				Message: "contains potentially dangerous SQL characters or keywords",
			}
		}
	}

	return nil
}

// SanitizeSQLIdentifier sanitizes a SQL identifier by removing unsafe characters
func SanitizeSQLIdentifier(name string) string {
	// Remove all non-alphanumeric and non-underscore characters
	var result strings.Builder
	for i, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9' && i > 0) || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ValidateFilePath validates file paths to prevent path traversal attacks
func ValidateFilePath(path, fieldName string) error {
	if path == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	if len(path) > MaxPathLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxPathLength),
		}
	}

	// Check for path traversal patterns
	cleanPath := filepath.Clean(path)
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(path, pattern) {
			return &ValidationError{
				Field:   fieldName,
				Message: "contains potentially dangerous path traversal characters",
			}
		}
	}

	// Additional check for absolute path attempts when not expected
	if filepath.IsAbs(cleanPath) && !filepath.IsAbs(path) {
		return &ValidationError{
			Field:   fieldName,
			Message: "path normalization detected potential traversal attempt",
		}
	}

	return nil
}

// SanitizeFilePath sanitizes a file path
func SanitizeFilePath(path string) string {
	return filepath.Clean(path)
}

// ValidatePort validates port numbers
func ValidatePort(port int, fieldName string) error {
	if port < 0 || port > 65535 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be between 0 and 65535",
		}
	}
	return nil
}

// ValidateHost validates hostname or IP address
func ValidateHost(host, fieldName string) error {
	if host == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	if len(host) > 255 {
		return &ValidationError{
			Field:   fieldName,
			Message: "exceeds maximum length of 255 characters",
		}
	}

	// Check for dangerous characters
	for _, pattern := range []string{"|", ";", "&", "`", "$", "(", ")", "<", ">"} {
		if strings.Contains(host, pattern) {
			return &ValidationError{
				Field:   fieldName,
				Message: "contains potentially dangerous characters",
			}
		}
	}

	return nil
}

// ValidateDatabaseType validates database type
func ValidateDatabaseType(dbType, fieldName string) error {
	validTypes := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"sqlite":   true,
	}

	if dbType == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	if !validTypes[strings.ToLower(dbType)] {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be one of: postgres, mysql, sqlite",
		}
	}

	return nil
}

// ValidateSourceType validates data source type
func ValidateSourceType(sourceType, fieldName string) error {
	validTypes := map[string]bool{
		"filesystem": true,
		"csv":        true,
		"tsv":        true,
		"database":   true,
	}

	if sourceType == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	if !validTypes[strings.ToLower(sourceType)] {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be one of: filesystem, csv, tsv, database",
		}
	}

	return nil
}

// ValidateCustomQuery performs basic validation on custom SQL queries
func ValidateCustomQuery(query, fieldName string) error {
	if query == "" {
		return nil // Empty is allowed, will use default query
	}

	if len(query) > MaxQueryLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxQueryLength),
		}
	}

	// Query must be a SELECT statement
	upperQuery := strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(upperQuery, "SELECT") {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be a SELECT statement",
		}
	}

	// Block dangerous operations
	dangerousKeywords := []string{
		"DROP",
		"DELETE",
		"UPDATE",
		"INSERT",
		"ALTER",
		"CREATE",
		"EXEC",
		"EXECUTE",
		"xp_",
		"sp_",
		"INTO OUTFILE",
		"INTO DUMPFILE",
		"LOAD_FILE",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperQuery, keyword) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("contains forbidden keyword: %s", keyword),
			}
		}
	}

	return nil
}

// ValidateFileExtension validates file extensions
func ValidateFileExtension(ext, fieldName string) error {
	if ext == "" {
		return nil // Empty is allowed
	}

	// Remove leading dot if present
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	// Check length
	if len(ext) > 10 {
		return &ValidationError{
			Field:   fieldName,
			Message: "extension too long (max 10 characters)",
		}
	}

	// Check for valid characters
	for _, r := range ext {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return &ValidationError{
				Field:   fieldName,
				Message: "extension must contain only alphanumeric characters",
			}
		}
	}

	return nil
}

// ValidateConfigName validates configuration profile names
// Empty names are allowed for unnamed/default configurations
func ValidateConfigName(name, fieldName string) error {
	if name == "" {
		return nil // Empty is allowed for unnamed configs
	}

	if len(name) > MaxConfigNameLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxConfigNameLength),
		}
	}

	// Check for valid characters (alphanumeric, space, dash, underscore)
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '_') {
			return &ValidationError{
				Field:   fieldName,
				Message: "must contain only alphanumeric characters, spaces, dashes, and underscores",
			}
		}
	}

	return nil
}

// ValidateCSVDelimiter validates CSV delimiter
func ValidateCSVDelimiter(delimiter, fieldName string) error {
	if delimiter == "" {
		return &ValidationError{Field: fieldName, Message: "cannot be empty"}
	}

	validDelimiters := map[string]bool{
		",":  true,
		"\t": true,
		";":  true,
		"|":  true,
	}

	if !validDelimiters[delimiter] {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be one of: comma (,), tab (\\t), semicolon (;), or pipe (|)",
		}
	}

	return nil
}

// ValidatePositiveInt validates non-negative integers (zero or positive)
func ValidatePositiveInt(value int, fieldName string) error {
	if value < 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be zero or a positive integer",
		}
	}
	return nil
}

// ValidateNonNegativeInt validates non-negative integers
func ValidateNonNegativeInt(value int, fieldName string) error {
	if value < 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be zero or a positive integer",
		}
	}
	return nil
}
