package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/baditaflorin/codexgigantus/pkg/config"
	"github.com/baditaflorin/codexgigantus/pkg/utils"
)

func TestShouldIgnoreDir(t *testing.T) {
	cfg := &config.Config{
		IgnoreDirs: []string{".git", "node_modules", "vendor"},
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"/path/to/.git", true},
		{"/path/to/node_modules", true},
		{"/path/to/vendor/lib", true},
		{"/path/to/src", false},
		{"/normal/directory", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldIgnoreDir(tt.path, cfg)
			if result != tt.expected {
				t.Errorf("shouldIgnoreDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		cfg      *config.Config
		expected bool
	}{
		{
			name: "ignore by filename",
			path: "/path/to/.DS_Store",
			cfg: &config.Config{
				IgnoreFiles: []string{".DS_Store"},
			},
			expected: true,
		},
		{
			name: "ignore by extension",
			path: "/path/to/file.log",
			cfg: &config.Config{
				IgnoreExts: []string{"log", "tmp"},
			},
			expected: true,
		},
		{
			name: "include only specific extensions - match",
			path: "/path/to/file.go",
			cfg: &config.Config{
				IncludeExts: []string{"go", "md"},
			},
			expected: false,
		},
		{
			name: "include only specific extensions - no match",
			path: "/path/to/file.txt",
			cfg: &config.Config{
				IncludeExts: []string{"go", "md"},
			},
			expected: true,
		},
		{
			name: "no filters",
			path: "/path/to/file.txt",
			cfg: &config.Config{
				IgnoreFiles: []string{},
				IgnoreExts:  []string{},
				IncludeExts: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIgnoreFile(tt.path, tt.cfg)
			if result != tt.expected {
				t.Errorf("shouldIgnoreFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestProcessFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.go":  "package main",
		"file3.log": "log content",
	}

	for name, content := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	// Create a subdirectory with a file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatalf("Failed to create subfile: %v", err)
	}

	tests := []struct {
		name          string
		cfg           *config.Config
		expectedCount int
		checkContent  func(t *testing.T, results []utils.FileResult)
	}{
		{
			name: "process all files recursively",
			cfg: &config.Config{
				Dirs:      []string{tmpDir},
				Recursive: true,
			},
			expectedCount: 4, // all files including subdirectory
		},
		{
			name: "ignore log files",
			cfg: &config.Config{
				Dirs:       []string{tmpDir},
				Recursive:  true,
				IgnoreExts: []string{"log"},
			},
			expectedCount: 3, // all except .log
		},
		{
			name: "include only go files",
			cfg: &config.Config{
				Dirs:        []string{tmpDir},
				Recursive:   true,
				IncludeExts: []string{"go"},
			},
			expectedCount: 1, // only .go file
			checkContent: func(t *testing.T, results []utils.FileResult) {
				if len(results) > 0 && results[0].Content != "package main" {
					t.Errorf("Expected Go file content 'package main', got %q", results[0].Content)
				}
			},
		},
		{
			name: "non-recursive",
			cfg: &config.Config{
				Dirs:      []string{tmpDir},
				Recursive: false,
			},
			expectedCount: 3, // only files in root, not subdirectory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ProcessFiles(tt.cfg)
			if err != nil {
				t.Fatalf("ProcessFiles failed: %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d files, got %d", tt.expectedCount, len(results))
			}

			if tt.checkContent != nil {
				tt.checkContent(t, results)
			}
		})
	}
}

func TestProcessFilesError(t *testing.T) {
	// Test with non-existent directory
	cfg := &config.Config{
		Dirs:      []string{"/this/path/does/not/exist"},
		Recursive: true,
	}

	_, err := ProcessFiles(cfg)
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}
