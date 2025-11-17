package configfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	// Create test config
	originalConfig := &AppConfig{
		SourceType:  "filesystem",
		Directories: []string{"/tmp", "/var"},
		Recursive:   true,
		OutputFile:  "test_output.txt",
		ShowSize:    true,
		Debug:       false,
		Name:        "Test Config",
		Description: "A test configuration",
	}

	// Save
	err := SaveJSON(originalConfig, testFile)
	if err != nil {
		t.Fatalf("SaveJSON() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("JSON file was not created")
	}

	// Load
	loadedConfig, err := LoadJSON(testFile)
	if err != nil {
		t.Fatalf("LoadJSON() error = %v", err)
	}

	// Compare
	if loadedConfig.SourceType != originalConfig.SourceType {
		t.Errorf("SourceType = %v, want %v", loadedConfig.SourceType, originalConfig.SourceType)
	}
	if len(loadedConfig.Directories) != len(originalConfig.Directories) {
		t.Errorf("Directories length = %v, want %v", len(loadedConfig.Directories), len(originalConfig.Directories))
	}
	if loadedConfig.Recursive != originalConfig.Recursive {
		t.Errorf("Recursive = %v, want %v", loadedConfig.Recursive, originalConfig.Recursive)
	}
	if loadedConfig.Name != originalConfig.Name {
		t.Errorf("Name = %v, want %v", loadedConfig.Name, originalConfig.Name)
	}
}

func TestSaveLoadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	originalConfig := &AppConfig{
		SourceType:        "database",
		DBType:            "postgres",
		DBHost:            "localhost",
		DBPort:            5432,
		DBName:            "testdb",
		DBUser:            "testuser",
		DBPassword:        "testpass",
		DBTableName:       "files",
		DBColumnPath:      "path",
		DBColumnContent:   "content",
		OutputFile:        "db_output.txt",
		Name:              "Database Config",
	}

	// Save
	err := SaveYAML(originalConfig, testFile)
	if err != nil {
		t.Fatalf("SaveYAML() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("YAML file was not created")
	}

	// Load
	loadedConfig, err := LoadYAML(testFile)
	if err != nil {
		t.Fatalf("LoadYAML() error = %v", err)
	}

	// Compare
	if loadedConfig.SourceType != originalConfig.SourceType {
		t.Errorf("SourceType = %v, want %v", loadedConfig.SourceType, originalConfig.SourceType)
	}
	if loadedConfig.DBHost != originalConfig.DBHost {
		t.Errorf("DBHost = %v, want %v", loadedConfig.DBHost, originalConfig.DBHost)
	}
	if loadedConfig.DBPort != originalConfig.DBPort {
		t.Errorf("DBPort = %v, want %v", loadedConfig.DBPort, originalConfig.DBPort)
	}
}

func TestSaveLoadAutoDetect(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
	}{
		{"JSON file", "config.json"},
		{"YAML file", "config.yaml"},
		{"YML file", "config.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tt.filename)

			config := &AppConfig{
				SourceType:  "csv",
				CSVFilePath: "/tmp/data.csv",
				CSVDelimiter: ",",
				OutputFile:  "output.txt",
			}

			// Save with auto-detect
			err := Save(config, testFile)
			if err != nil {
				t.Fatalf("Save() error = %v", err)
			}

			// Load with auto-detect
			loadedConfig, err := Load(testFile)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if loadedConfig.SourceType != config.SourceType {
				t.Errorf("SourceType = %v, want %v", loadedConfig.SourceType, config.SourceType)
			}
			if loadedConfig.CSVFilePath != config.CSVFilePath {
				t.Errorf("CSVFilePath = %v, want %v", loadedConfig.CSVFilePath, config.CSVFilePath)
			}
		})
	}
}

func TestSaveUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.txt")

	config := NewDefault()
	err := Save(config, testFile)
	if err == nil {
		t.Error("Save() should fail for unsupported format")
	}
}

func TestLoadUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	_, err := Load(testFile)
	if err == nil {
		t.Error("Load() should fail for unsupported format")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AppConfig
		wantErr bool
	}{
		{
			name: "valid filesystem config",
			config: &AppConfig{
				SourceType:  "filesystem",
				Directories: []string{"/tmp"},
			},
			wantErr: false,
		},
		{
			name: "valid csv config",
			config: &AppConfig{
				SourceType:  "csv",
				CSVFilePath: "/tmp/data.csv",
			},
			wantErr: false,
		},
		{
			name: "valid database config with table",
			config: &AppConfig{
				SourceType:  "database",
				DBType:      "postgres",
				DBTableName: "files",
			},
			wantErr: false,
		},
		{
			name: "valid database config with query",
			config: &AppConfig{
				SourceType: "database",
				DBType:     "postgres",
				DBQuery:    "SELECT * FROM files",
			},
			wantErr: false,
		},
		{
			name:    "missing source_type",
			config:  &AppConfig{},
			wantErr: true,
		},
		{
			name: "filesystem without directories",
			config: &AppConfig{
				SourceType: "filesystem",
			},
			wantErr: true,
		},
		{
			name: "csv without file path",
			config: &AppConfig{
				SourceType: "csv",
			},
			wantErr: true,
		},
		{
			name: "database without type",
			config: &AppConfig{
				SourceType:  "database",
				DBTableName: "files",
			},
			wantErr: true,
		},
		{
			name: "database without table or query",
			config: &AppConfig{
				SourceType: "database",
				DBType:     "postgres",
			},
			wantErr: true,
		},
		{
			name: "invalid source_type",
			config: &AppConfig{
				SourceType: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config *AppConfig
		check  func(*AppConfig) error
	}{
		{
			name:   "sets default output file",
			config: &AppConfig{SourceType: "filesystem"},
			check: func(c *AppConfig) error {
				if c.OutputFile != "output.txt" {
					t.Errorf("OutputFile = %v, want output.txt", c.OutputFile)
				}
				return nil
			},
		},
		{
			name:   "sets default directories",
			config: &AppConfig{SourceType: "filesystem"},
			check: func(c *AppConfig) error {
				if len(c.Directories) != 1 || c.Directories[0] != "." {
					t.Errorf("Directories = %v, want [.]", c.Directories)
				}
				return nil
			},
		},
		{
			name:   "sets csv delimiter",
			config: &AppConfig{SourceType: "csv"},
			check: func(c *AppConfig) error {
				if c.CSVDelimiter != "," {
					t.Errorf("CSVDelimiter = %v, want ,", c.CSVDelimiter)
				}
				return nil
			},
		},
		{
			name:   "sets tsv delimiter",
			config: &AppConfig{SourceType: "tsv"},
			check: func(c *AppConfig) error {
				if c.CSVDelimiter != "\t" {
					t.Errorf("CSVDelimiter = %v, want \\t", c.CSVDelimiter)
				}
				return nil
			},
		},
		{
			name:   "sets postgres port",
			config: &AppConfig{SourceType: "database", DBType: "postgres"},
			check: func(c *AppConfig) error {
				if c.DBPort != 5432 {
					t.Errorf("DBPort = %v, want 5432", c.DBPort)
				}
				return nil
			},
		},
		{
			name:   "sets mysql port",
			config: &AppConfig{SourceType: "database", DBType: "mysql"},
			check: func(c *AppConfig) error {
				if c.DBPort != 3306 {
					t.Errorf("DBPort = %v, want 3306", c.DBPort)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			tt.check(tt.config)
		})
	}
}

func TestNewDefault(t *testing.T) {
	config := NewDefault()

	if config.SourceType != "filesystem" {
		t.Errorf("SourceType = %v, want filesystem", config.SourceType)
	}
	if len(config.Directories) != 1 || config.Directories[0] != "." {
		t.Errorf("Directories = %v, want [.]", config.Directories)
	}
	if config.Recursive != true {
		t.Errorf("Recursive = %v, want true", config.Recursive)
	}
	if config.OutputFile != "output.txt" {
		t.Errorf("OutputFile = %v, want output.txt", config.OutputFile)
	}
}
