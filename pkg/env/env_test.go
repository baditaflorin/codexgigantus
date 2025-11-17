package env

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns default when not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "returns env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "returns default when not set",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "",
			expected:     42,
		},
		{
			name:         "returns env value when set",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "100",
			expected:     100,
		},
		{
			name:         "returns default when invalid",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "invalid",
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "returns default when not set",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "",
			expected:     true,
		},
		{
			name:         "returns false when set to false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "returns true when set to true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "returns default when invalid",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvBool(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue []string
		envValue     string
		expected     []string
	}{
		{
			name:         "returns default when not set",
			key:          "TEST_SLICE",
			defaultValue: []string{"a", "b"},
			envValue:     "",
			expected:     []string{"a", "b"},
		},
		{
			name:         "returns parsed slice",
			key:          "TEST_SLICE",
			defaultValue: []string{"a"},
			envValue:     "x,y,z",
			expected:     []string{"x", "y", "z"},
		},
		{
			name:         "handles whitespace",
			key:          "TEST_SLICE",
			defaultValue: []string{"a"},
			envValue:     "x , y , z ",
			expected:     []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvSlice(tt.key, tt.defaultValue)
			if len(result) != len(tt.expected) {
				t.Errorf("getEnvSlice() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("getEnvSlice()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create temporary .env file
	tmpFile := t.TempDir() + "/.env"
	content := `# Test env file
APP_MODE=test
WEB_PORT=9090
# Comment line
EMPTY_VALUE=

VALID_KEY=valid_value
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("APP_MODE")
	os.Unsetenv("WEB_PORT")
	os.Unsetenv("VALID_KEY")

	err := loadEnvFile(tmpFile)
	if err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	// Verify values were loaded
	if os.Getenv("APP_MODE") != "test" {
		t.Errorf("APP_MODE = %v, want test", os.Getenv("APP_MODE"))
	}
	if os.Getenv("WEB_PORT") != "9090" {
		t.Errorf("WEB_PORT = %v, want 9090", os.Getenv("WEB_PORT"))
	}
	if os.Getenv("VALID_KEY") != "valid_value" {
		t.Errorf("VALID_KEY = %v, want valid_value", os.Getenv("VALID_KEY"))
	}

	// Cleanup
	os.Unsetenv("APP_MODE")
	os.Unsetenv("WEB_PORT")
	os.Unsetenv("VALID_KEY")
}

func TestLoad(t *testing.T) {
	// This test verifies the Load function returns a valid config
	// We don't create a .env file, so it should use defaults
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Verify some defaults
	if cfg.WebPort <= 0 {
		t.Errorf("WebPort = %v, want positive integer", cfg.WebPort)
	}
	if cfg.DBType == "" {
		t.Error("DBType is empty")
	}
	if cfg.MaxFileSize <= 0 {
		t.Errorf("MaxFileSize = %v, want positive integer", cfg.MaxFileSize)
	}
}
