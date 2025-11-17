// Package utils provides utility functions for file processing and output generation.
// It includes functions for reading files, extracting Go function signatures,
// and generating formatted output.
package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// FileResult represents a processed file with its path and content.
type FileResult struct {
	Path    string // Path to the file
	Content string // Content of the file
}

// GenerateOutput creates formatted output from file results.
// If showFuncs is true, it extracts only function signatures from Go files.
func GenerateOutput(results []FileResult, showFuncs bool) string {
	var output strings.Builder
	for _, result := range results {
		output.WriteString(fmt.Sprintf("File: %s\n", result.Path))
		if showFuncs && IsGoFile(result.Path) {
			funcs := ExtractFunctions(result.Content)
			for _, fn := range funcs {
				output.WriteString(fmt.Sprintf("Function: %s\n", fn))
			}
		} else {
			output.WriteString(result.Content)
			output.WriteString("\n\n")
		}
	}
	return output.String()
}

// SaveOutput writes the output string to a file.
func SaveOutput(output, filename string) error {
	return os.WriteFile(filename, []byte(output), 0644)
}

// IsGoFile checks if a file has a .go extension.
func IsGoFile(path string) bool {
	return filepath.Ext(path) == ".go"
}

// ExtractFunctions parses Go source code and extracts function signatures.
// Returns a slice of function names found in the source code.
func ExtractFunctions(content string) []string {
	var funcs []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return funcs
	}

	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcs = append(funcs, fn.Name.Name)
		}
	}

	return funcs
}

// Debug prints debug messages when debug mode is enabled.
// It's a simple helper function to avoid checking debug flags everywhere.
func Debug(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}
