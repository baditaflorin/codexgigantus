// utils.go
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

// FileResult represents a processed file with its path and content
type FileResult struct {
	Path    string
	Content string
}

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

func SaveOutput(output, filename string) error {
	return os.WriteFile(filename, []byte(output), 0644)
}

func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

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

func Debug(format string, args ...interface{}) {
	fmt.Printf("DEBUG: "+format+"\n", args...)
}
