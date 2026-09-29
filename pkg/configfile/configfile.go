// Package configfile provides configuration file management (save/load JSON/YAML)
package configfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/baditaflorin/codexgigantus/pkg/validation"
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

// Validate validates the configuration with security checks
func (c *AppConfig) Validate() error {
	// Validate source type
	if err := validation.ValidateSourceType(c.SourceType, "source_type"); err != nil {
		return err
	}

	// Validate config name if provided
	if err := validation.ValidateConfigName(c.Name, "name"); err != nil {
		return err
	}

	switch c.SourceType {
	case "filesystem":
		if len(c.Directories) == 0 {
			return fmt.Errorf("directories are required for filesystem source")
		}
		// Validate each directory path
		for i, dir := range c.Directories {
			if err := validation.ValidateFilePath(dir, fmt.Sprintf("directories[%d]", i)); err != nil {
				return fmt.Errorf("invalid directory path: %w", err)
			}
		}
		// Validate file extensions
		for i, ext := range c.IncludeExtensions {
			if err := validation.ValidateFileExtension(ext, fmt.Sprintf("include_extensions[%d]", i)); err != nil {
				return err
			}
		}
		for i, ext := range c.ExcludeExtensions {
			if err := validation.ValidateFileExtension(ext, fmt.Sprintf("exclude_extensions[%d]", i)); err != nil {
				return err
			}
		}

	case "csv", "tsv":
		if c.CSVFilePath == "" {
			return fmt.Errorf("csv_file_path is required for CSV/TSV source")
		}
		if err := validation.ValidateFilePath(c.CSVFilePath, "csv_file_path"); err != nil {
			return fmt.Errorf("invalid CSV file path: %w", err)
		}
		// Set default delimiter
		if c.SourceType == "csv" && c.CSVDelimiter == "" {
			c.CSVDelimiter = ","
		}
		if c.SourceType == "tsv" && c.CSVDelimiter == "" {
			c.CSVDelimiter = "\t"
		}
		// Validate delimiter
		if err := validation.ValidateCSVDelimiter(c.CSVDelimiter, "csv_delimiter"); err != nil {
			return err
		}
		// Validate column indices
		if err := validation.ValidateNonNegativeInt(c.CSVPathColumn, "csv_path_column"); err != nil {
			return err
		}
		if err := validation.ValidateNonNegativeInt(c.CSVContentColumn, "csv_content_column"); err != nil {
			return err
		}

	case "database":
		// Validate database type
		if err := validation.ValidateDatabaseType(c.DBType, "db_type"); err != nil {
			return err
		}
		// Validate host and port for non-SQLite databases
		if c.DBType != "sqlite" {
			if err := validation.ValidateHost(c.DBHost, "db_host"); err != nil {
				return err
			}
			if err := validation.ValidatePort(c.DBPort, "db_port"); err != nil {
				return err
			}
		}
		// Validate custom query or table/column names
		if c.DBQuery != "" {
			if err := validation.ValidateCustomQuery(c.DBQuery, "db_query"); err != nil {
				return err
			}
		} else {
			if c.DBTableName == "" {
				return fmt.Errorf("db_table_name is required when db_query is not provided")
			}
			if err := validation.ValidateSQLIdentifier(c.DBTableName, "db_table_name"); err != nil {
				return err
			}
			if err := validation.ValidateSQLIdentifier(c.DBColumnPath, "db_column_path"); err != nil {
				return err
			}
			if err := validation.ValidateSQLIdentifier(c.DBColumnContent, "db_column_content"); err != nil {
				return err
			}
			if c.DBColumnType != "" {
				if err := validation.ValidateSQLIdentifier(c.DBColumnType, "db_column_type"); err != nil {
					return err
				}
			}
			if c.DBColumnSize != "" {
				if err := validation.ValidateSQLIdentifier(c.DBColumnSize, "db_column_size"); err != nil {
					return err
				}
			}
		}
	}

	// Validate output file path if provided
	if c.OutputFile != "" {
		if err := validation.ValidateFilePath(c.OutputFile, "output_file"); err != nil {
			return fmt.Errorf("invalid output file path: %w", err)
		}
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
