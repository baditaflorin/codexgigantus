// Package config provides configuration management for CodexGigantus.
// It handles parsing and storing application settings.
package config

import (
	"strings"
)

// Config holds the configuration for the application.
// It defines all the parameters that control how files are processed.
type Config struct {
	// Dirs is a list of directories to search for files
	Dirs []string
	// IgnoreFiles is a list of specific file names to ignore
	IgnoreFiles []string
	// IgnoreDirs is a list of directory names to ignore during traversal
	IgnoreDirs []string
	// IgnoreExts is a list of file extensions to ignore (without the dot)
	IgnoreExts []string
	// IncludeExts is a list of file extensions to include (without the dot).
	// If specified, only files with these extensions will be processed
	IncludeExts []string
	// Recursive determines whether to search directories recursively
	Recursive bool
	// Debug enables debug output when true
	Debug bool
	// Save determines whether to save output to a file
	Save bool
	// OutputFile is the path where output should be saved if Save is true
	OutputFile string
	// ShowSize determines whether to display the size of the result
	ShowSize bool
	// ShowFuncs determines whether to show only function signatures
	// (only applicable for Go files)
	ShowFuncs bool
}

// ParseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
// Empty strings result in an empty slice. Leading and trailing whitespace is removed
// from each element.
func ParseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
