// Package utils provides utility functions for output generation and file operations.
// It includes functions for formatting results, saving output, and extracting code information.
package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/baditaflorin/codexgigantus/pkg/config"
)

// FileResult represents a processed file with its path and content.
// It is the output of file processing operations.
type FileResult struct {
	// Path is the file system path to the file
	Path string
	// Content is the text content of the file
	Content string
}

// GenerateOutput formats the processed file results into a string.
// If ShowFuncs is enabled in the config, it extracts and displays only function signatures
// for Go files. Otherwise, it displays the full file content.
func GenerateOutput(results []FileResult, cfg *config.Config) string {
	var buffer bytes.Buffer

	for _, result := range results {
		if cfg.ShowFuncs && isGoFile(result.Path) {
			funcs := extractFunctions(result.Content)
			if len(funcs) > 0 {
				buffer.WriteString(fmt.Sprintf("File: %s\n", result.Path))
				for _, f := range funcs {
					buffer.WriteString(fmt.Sprintf("Function: %s\n", f))
				}
				buffer.WriteString("\n")
			}
		} else {
			buffer.WriteString(fmt.Sprintf("File: %s\n", result.Path))
			buffer.WriteString(result.Content)
			buffer.WriteString("\n\n")
		}
	}

	return buffer.String()
}

// SaveOutput writes the given output string to a file with the specified filename.
// It creates or overwrites the file with 0644 permissions.
func SaveOutput(output, filename string) error {
	return os.WriteFile(filename, []byte(output), 0644)
}

// isGoFile checks if a file path has a .go extension.
func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

// extractFunctions parses Go source code and extracts function signatures.
// It returns a slice of strings representing function names and their parameters.
func extractFunctions(content string) []string {
	var funcs []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return funcs
	}

	for _, f := range node.Decls {
		if fn, ok := f.(*ast.FuncDecl); ok {
			funcSignature := fn.Name.Name + "("
			params := []string{}
			for _, param := range fn.Type.Params.List {
				names := []string{}
				for _, name := range param.Names {
					names = append(names, name.Name)
				}
				paramType := fmt.Sprintf("%s", param.Type)
				params = append(params, strings.Join(names, ", ")+" "+paramType)
			}
			funcSignature += strings.Join(params, ", ") + ")"
			funcs = append(funcs, funcSignature)
		}
	}

	return funcs
}

// Debug prints a formatted debug message to stdout.
// Messages are prefixed with "DEBUG: " to distinguish them from regular output.
func Debug(format string, args ...interface{}) {
	fmt.Printf("DEBUG: "+format+"\n", args...)
}
