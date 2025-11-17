package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileResult(t *testing.T) {
	result := FileResult{
		Path:    "/path/to/file.go",
		Content: "package main",
	}

	if result.Path != "/path/to/file.go" {
		t.Errorf("Path = %v, want /path/to/file.go", result.Path)
	}
	if result.Content != "package main" {
		t.Errorf("Content = %v, want 'package main'", result.Content)
	}
}

func TestGenerateOutput(t *testing.T) {
	tests := []struct {
		name      string
		results   []FileResult
		showFuncs bool
		want      string
	}{
		{
			name: "basic output",
			results: []FileResult{
				{Path: "test.txt", Content: "hello"},
			},
			showFuncs: false,
			want:      "File: test.txt\nhello\n\n",
		},
		{
			name: "multiple files",
			results: []FileResult{
				{Path: "file1.txt", Content: "content1"},
				{Path: "file2.txt", Content: "content2"},
			},
			showFuncs: false,
			want:      "File: file1.txt\ncontent1\n\nFile: file2.txt\ncontent2\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateOutput(tt.results, tt.showFuncs)
			if got != tt.want {
				t.Errorf("GenerateOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveOutput(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.txt")
	content := "test content"

	err := SaveOutput(content, testFile)
	if err != nil {
		t.Fatalf("SaveOutput() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Verify content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(data) != content {
		t.Errorf("File content = %v, want %v", string(data), content)
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"main.go", true},
		{"test.py", false},
		{"file.txt", false},
		{"/path/to/file.go", true},
		{"file.GO", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := IsGoFile(tt.path)
			if got != tt.want {
				t.Errorf("IsGoFile(%v) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestExtractFunctions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "simple function",
			content: "package main\nfunc hello() {}",
			want:    []string{"hello"},
		},
		{
			name:    "multiple functions",
			content: "package main\nfunc foo() {}\nfunc bar() {}",
			want:    []string{"foo", "bar"},
		},
		{
			name:    "no functions",
			content: "package main\nvar x = 1",
			want:    []string{},
		},
		{
			name:    "invalid syntax",
			content: "invalid go code",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFunctions(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractFunctions() returned %d functions, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ExtractFunctions()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDebug(t *testing.T) {
	// This test just ensures Debug doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Debug() panicked: %v", r)
		}
	}()

	Debug("test message")
	Debug("test with args: %s %d", "hello", 42)
}

func TestGenerateOutputWithShowFuncs(t *testing.T) {
	goContent := `package main

func main() {
	println("hello")
}

func helper() {
	println("helper")
}`

	results := []FileResult{
		{Path: "main.go", Content: goContent},
		{Path: "test.txt", Content: "not go"},
	}

	output := GenerateOutput(results, true)

	// Should contain function names for .go file
	if !strings.Contains(output, "main") || !strings.Contains(output, "helper") {
		t.Error("Output should contain function names for Go files")
	}

	// Should contain full content for non-Go files
	if !strings.Contains(output, "not go") {
		t.Error("Output should contain full content for non-Go files")
	}
}
