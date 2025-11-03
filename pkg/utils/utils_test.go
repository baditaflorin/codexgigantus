package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/baditaflorin/codexgigantus/pkg/config"
)

func TestFileResult(t *testing.T) {
	fr := FileResult{
		Path:    "/path/to/file.txt",
		Content: "test content",
	}

	if fr.Path != "/path/to/file.txt" {
		t.Errorf("Expected path '/path/to/file.txt', got %s", fr.Path)
	}
	if fr.Content != "test content" {
		t.Errorf("Expected content 'test content', got %s", fr.Content)
	}
}

func TestGenerateOutput(t *testing.T) {
	tests := []struct {
		name     string
		results  []FileResult
		config   *config.Config
		contains []string
	}{
		{
			name: "basic output",
			results: []FileResult{
				{Path: "test.txt", Content: "hello world"},
			},
			config: &config.Config{
				ShowFuncs: false,
			},
			contains: []string{"File: test.txt", "hello world"},
		},
		{
			name: "multiple files",
			results: []FileResult{
				{Path: "file1.txt", Content: "content1"},
				{Path: "file2.txt", Content: "content2"},
			},
			config: &config.Config{
				ShowFuncs: false,
			},
			contains: []string{"File: file1.txt", "content1", "File: file2.txt", "content2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GenerateOutput(tt.results, tt.config)
			for _, s := range tt.contains {
				if !strings.Contains(output, s) {
					t.Errorf("Expected output to contain %q, but it didn't. Output: %s", s, output)
				}
			}
		})
	}
}

func TestSaveOutput(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_output.txt")

	content := "test content for save"
	err := SaveOutput(content, testFile)
	if err != nil {
		t.Fatalf("SaveOutput failed: %v", err)
	}

	// Read back the file
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Expected content %q, got %q", content, string(data))
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"test.go", true},
		{"test.txt", false},
		{"file.GO", false}, // case sensitive
		{"/path/to/file.go", true},
		{"no_extension", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isGoFile(tt.path)
			if result != tt.expected {
				t.Errorf("isGoFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestExtractFunctions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // number of functions expected
	}{
		{
			name: "simple function",
			content: `package main
func main() {
	println("hello")
}`,
			expected: 1,
		},
		{
			name: "multiple functions",
			content: `package main
func foo() {}
func bar(x int) {}
func baz(a string, b int) string { return "" }`,
			expected: 3,
		},
		{
			name:     "no functions",
			content:  `package main\nvar x = 5`,
			expected: 0,
		},
		{
			name:     "invalid syntax",
			content:  `this is not valid go code`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFunctions(tt.content)
			if len(result) != tt.expected {
				t.Errorf("extractFunctions() returned %d functions, expected %d. Functions: %v",
					len(result), tt.expected, result)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	// Debug function writes to stdout, so we just ensure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Debug panicked: %v", r)
		}
	}()

	Debug("test message")
	Debug("test with args: %s %d", "hello", 42)
}
