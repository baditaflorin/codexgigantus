// Package configfile provides configuration file management (save/load JSON/YAML)
package configfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppConfig represents the application configuration that can be saved/loaded
type AppConfig struct {
	// Source configuration
	SourceType string `json:"source_type" yaml:"source_type"` // filesystem, csv, tsv, database

	// Filesystem source settings
	Directories      []string `json:"directories,omitempty" yaml:"directories,omitempty"`
	Recursive        bool     `json:"recursive" yaml:"recursive"`
	IgnoreFiles      []string `json:"ignore_files,omitempty" yaml:"ignore_files,omitempty"`
	IgnoreDirs       []string `json:"ignore_dirs,omitempty" yaml:"ignore_dirs,omitempty"`
	ExcludeExtensions []string `json:"exclude_extensions,omitempty" yaml:"exclude_extensions,omitempty"`
	IncludeExtensions []string `json:"include_extensions,omitempty" yaml:"include_extensions,omitempty"`

	// CSV/TSV source settings
	CSVFilePath      string `json:"csv_file_path,omitempty" yaml:"csv_file_path,omitempty"`
	CSVDelimiter     string `json:"csv_delimiter,omitempty" yaml:"csv_delimiter,omitempty"` // "," for CSV, "\t" for TSV
	CSVPathColumn    int    `json:"csv_path_column,omitempty" yaml:"csv_path_column,omitempty"`
	CSVContentColumn int    `json:"csv_content_column,omitempty" yaml:"csv_content_column,omitempty"`
	CSVHasHeader     bool   `json:"csv_has_header" yaml:"csv_has_header"`

	// Database source settings
	DBType          string `json:"db_type,omitempty" yaml:"db_type,omitempty"`           // postgres, mysql, sqlite
	DBHost          string `json:"db_host,omitempty" yaml:"db_host,omitempty"`
	DBPort          int    `json:"db_port,omitempty" yaml:"db_port,omitempty"`
	DBName          string `json:"db_name,omitempty" yaml:"db_name,omitempty"`
	DBUser          string `json:"db_user,omitempty" yaml:"db_user,omitempty"`
	DBPassword      string `json:"db_password,omitempty" yaml:"db_password,omitempty"`
	DBSSLMode       string `json:"db_ssl_mode,omitempty" yaml:"db_ssl_mode,omitempty"`
	DBTableName     string `json:"db_table_name,omitempty" yaml:"db_table_name,omitempty"`
	DBColumnPath    string `json:"db_column_path,omitempty" yaml:"db_column_path,omitempty"`
	DBColumnContent string `json:"db_column_content,omitempty" yaml:"db_column_content,omitempty"`
	DBColumnType    string `json:"db_column_type,omitempty" yaml:"db_column_type,omitempty"`
	DBColumnSize    string `json:"db_column_size,omitempty" yaml:"db_column_size,omitempty"`
	DBQuery         string `json:"db_query,omitempty" yaml:"db_query,omitempty"` // Optional custom query

	// Output settings
	OutputFile string `json:"output_file" yaml:"output_file"`
	ShowSize   bool   `json:"show_size" yaml:"show_size"`
	ShowFuncs  bool   `json:"show_funcs" yaml:"show_funcs"`
	Debug      bool   `json:"debug" yaml:"debug"`

	// Metadata
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`               // Config profile name
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // Config description
}

// SaveJSON saves the configuration to a JSON file
func SaveJSON(config *AppConfig, filepath string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// LoadJSON loads configuration from a JSON file
func LoadJSON(filepath string) (*AppConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &config, nil
}

// SaveYAML saves the configuration to a YAML file
func SaveYAML(config *AppConfig, filepath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// LoadYAML loads configuration from a YAML file
func LoadYAML(filepath string) (*AppConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to a file (auto-detects format from extension)
func Save(config *AppConfig, path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return SaveJSON(config, path)
	case ".yaml", ".yml":
		return SaveYAML(config, path)
	default:
		return fmt.Errorf("unsupported file format: %s (use .json, .yaml, or .yml)", ext)
	}
}

// Load loads configuration from a file (auto-detects format from extension)
func Load(path string) (*AppConfig, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return LoadJSON(path)
	case ".yaml", ".yml":
		return LoadYAML(path)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (use .json, .yaml, or .yml)", ext)
	}
}

// Validate validates the configuration
func (c *AppConfig) Validate() error {
	if c.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}

	switch c.SourceType {
	case "filesystem":
		if len(c.Directories) == 0 {
			return fmt.Errorf("directories are required for filesystem source")
		}
	case "csv", "tsv":
		if c.CSVFilePath == "" {
			return fmt.Errorf("csv_file_path is required for CSV/TSV source")
		}
		if c.SourceType == "csv" && c.CSVDelimiter == "" {
			c.CSVDelimiter = ","
		}
		if c.SourceType == "tsv" && c.CSVDelimiter == "" {
			c.CSVDelimiter = "\t"
		}
	case "database":
		if c.DBType == "" {
			return fmt.Errorf("db_type is required for database source")
		}
		if c.DBTableName == "" && c.DBQuery == "" {
			return fmt.Errorf("either db_table_name or db_query is required for database source")
		}
	default:
		return fmt.Errorf("invalid source_type: %s (must be filesystem, csv, tsv, or database)", c.SourceType)
	}

	return nil
}

// SetDefaults sets default values for optional fields
func (c *AppConfig) SetDefaults() {
	if c.OutputFile == "" {
		c.OutputFile = "output.txt"
	}

	if c.SourceType == "filesystem" && c.Directories == nil {
		c.Directories = []string{"."}
	}

	if c.SourceType == "csv" && c.CSVDelimiter == "" {
		c.CSVDelimiter = ","
	}

	if c.SourceType == "tsv" && c.CSVDelimiter == "" {
		c.CSVDelimiter = "\t"
	}

	if c.DBType != "" {
		if c.DBHost == "" {
			c.DBHost = "localhost"
		}
		if c.DBPort == 0 {
			switch c.DBType {
			case "postgres":
				c.DBPort = 5432
			case "mysql":
				c.DBPort = 3306
			case "sqlite":
				c.DBPort = 0 // SQLite doesn't use ports
			}
		}
		if c.DBSSLMode == "" {
			c.DBSSLMode = "disable"
		}
	}
}

// NewDefault creates a new AppConfig with default values
func NewDefault() *AppConfig {
	config := &AppConfig{
		SourceType:  "filesystem",
		Directories: []string{"."},
		Recursive:   true,
		OutputFile:  "output.txt",
	}
	return config
}
