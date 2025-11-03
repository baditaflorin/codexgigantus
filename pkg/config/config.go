// config.go
package config

import (
	"strings"
)

// Config holds the configuration for the application
type Config struct {
	Dirs        []string
	IgnoreFiles []string
	IgnoreDirs  []string
	IgnoreExts  []string
	IncludeExts []string
	Recursive   bool
	Debug       bool
	Save        bool
	OutputFile  string
	ShowSize    bool
	ShowFuncs   bool
}

// ParseCommaSeparated splits a comma-separated string into a slice of trimmed strings
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
