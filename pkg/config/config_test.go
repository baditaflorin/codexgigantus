package config

import (
	"reflect"
	"testing"
)

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single value",
			input:    "test",
			expected: []string{"test"},
		},
		{
			name:     "multiple values",
			input:    "foo,bar,baz",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			name:     "values with spaces",
			input:    "foo, bar , baz",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			name:     "values with extra spaces",
			input:    "  foo  ,  bar  ,  baz  ",
			expected: []string{"foo", "bar", "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCommaSeparated(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCommaSeparated(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConfigStruct(t *testing.T) {
	// Test that Config struct can be instantiated with expected fields
	cfg := &Config{
		Dirs:        []string{"dir1", "dir2"},
		IgnoreFiles: []string{"file1.txt"},
		IgnoreDirs:  []string{".git"},
		IgnoreExts:  []string{"log"},
		IncludeExts: []string{"go", "md"},
		Recursive:   true,
		Debug:       false,
		Save:        true,
		OutputFile:  "output.txt",
		ShowSize:    false,
		ShowFuncs:   true,
	}

	if len(cfg.Dirs) != 2 {
		t.Errorf("Expected 2 dirs, got %d", len(cfg.Dirs))
	}
	if cfg.Dirs[0] != "dir1" {
		t.Errorf("Expected first dir to be 'dir1', got %s", cfg.Dirs[0])
	}
	if !cfg.Recursive {
		t.Error("Expected Recursive to be true")
	}
	if cfg.Debug {
		t.Error("Expected Debug to be false")
	}
}
