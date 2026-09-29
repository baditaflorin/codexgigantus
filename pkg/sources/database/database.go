// Package database provides database connectivity and querying functionality
package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/baditaflorin/codexgigantus/pkg/utils"
	"github.com/baditaflorin/codexgigantus/pkg/validation"
)

// Processor handles database operations
type Processor struct {
	DBType         string
	Host           string
	Port           int
	DBName         string
	User           string
	Password       string
	SSLMode        string
	TableName      string
	ColumnPath     string
	ColumnContent  string
	ColumnType     string
	ColumnSize     string
	CustomQuery    string
	Debug          bool
	db             *sql.DB
}

// NewProcessor creates a new database processor
func NewProcessor(dbType, host string, port int, dbName, user, password, sslMode string, debug bool) *Processor {
	return &Processor{
		DBType:   dbType,
		Host:     host,
		Port:     port,
		DBName:   dbName,
		User:     user,
		Password: password,
		SSLMode:  sslMode,
		Debug:    debug,
	}
}

// Connect establishes a database connection with secure error handling
func (p *Processor) Connect() error {
	var dsn string
	var driverName string

	switch p.DBType {
	case "postgres":
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)

	case "mysql":
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			p.User, p.Password, p.Host, p.Port, p.DBName)

	case "sqlite":
		driverName = "sqlite3"
		dsn = p.DBName // For SQLite, DBName is the file path

	default:
		return fmt.Errorf("unsupported database type")
	}

	if p.Debug {
		// Never log passwords or connection strings in production
		fmt.Printf("Connecting to %s database at %s\n", p.DBType, p.Host)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		// Don't expose internal error details
		return fmt.Errorf("failed to establish database connection")
	}

	// Set connection pool limits for security
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		// Don't expose internal error details
		return fmt.Errorf("database connection test failed")
	}

	p.db = db

	if p.Debug {
		fmt.Printf("Successfully connected to %s database\n", p.DBType)
	}

	return nil
}

// Close closes the database connection
func (p *Processor) Close() error {
	if p.db != nil {
		err := p.db.Close()
		p.db = nil // Set to nil after closing
		return err
	}
	return nil
}

// Process executes the query and returns file results
func (p *Processor) Process() ([]utils.FileResult, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	query, err := p.buildQuery()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	if p.Debug {
		// Log query without sensitive data
		fmt.Printf("Executing query on table: %s\n", p.TableName)
	}

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed")
	}
	defer rows.Close()

	var results []utils.FileResult

	for rows.Next() {
		var filePath, content string
		var fileType, fileSize sql.NullString

		// Scan based on whether optional columns are included
		if p.ColumnType != "" && p.ColumnSize != "" {
			err = rows.Scan(&filePath, &content, &fileType, &fileSize)
		} else if p.ColumnType != "" {
			err = rows.Scan(&filePath, &content, &fileType)
		} else if p.ColumnSize != "" {
			err = rows.Scan(&filePath, &content, &fileSize)
		} else {
			err = rows.Scan(&filePath, &content)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read database row")
		}

		results = append(results, utils.FileResult{
			Path:    filePath,
			Content: content,
		})

		if p.Debug {
			fmt.Printf("Retrieved: %s (%d bytes)\n", filePath, len(content))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database results")
	}

	if p.Debug {
		fmt.Printf("Retrieved %d records from database\n", len(results))
	}

	return results, nil
}

// buildQuery constructs the SQL query with validation to prevent SQL injection
func (p *Processor) buildQuery() (string, error) {
	// Use custom query if provided (already validated)
	if p.CustomQuery != "" {
		return p.CustomQuery, nil
	}

	// Validate all SQL identifiers to prevent SQL injection
	if err := validation.ValidateSQLIdentifier(p.TableName, "table_name"); err != nil {
		return "", fmt.Errorf("invalid table name: %w", err)
	}

	if err := validation.ValidateSQLIdentifier(p.ColumnPath, "column_path"); err != nil {
		return "", fmt.Errorf("invalid path column: %w", err)
	}

	if err := validation.ValidateSQLIdentifier(p.ColumnContent, "column_content"); err != nil {
		return "", fmt.Errorf("invalid content column: %w", err)
	}

	// Build column list with validated identifiers
	columns := []string{p.ColumnPath, p.ColumnContent}

	if p.ColumnType != "" {
		if err := validation.ValidateSQLIdentifier(p.ColumnType, "column_type"); err != nil {
			return "", fmt.Errorf("invalid type column: %w", err)
		}
		columns = append(columns, p.ColumnType)
	}

	if p.ColumnSize != "" {
		if err := validation.ValidateSQLIdentifier(p.ColumnSize, "column_size"); err != nil {
			return "", fmt.Errorf("invalid size column: %w", err)
		}
		columns = append(columns, p.ColumnSize)
	}

	// Safe to use fmt.Sprintf here because all identifiers have been validated
	return fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(columns, ", "),
		p.TableName), nil
}

// Validate validates the processor configuration with security checks
func (p *Processor) Validate() error {
	// Validate database type
	if err := validation.ValidateDatabaseType(p.DBType, "db_type"); err != nil {
		return err
	}

	// Validate host and port for non-SQLite databases
	if p.DBType != "sqlite" {
		if err := validation.ValidateHost(p.Host, "host"); err != nil {
			return err
		}
		if err := validation.ValidatePort(p.Port, "port"); err != nil {
			return err
		}
		if p.User == "" {
			return fmt.Errorf("user is required for %s", p.DBType)
		}
	}

	if p.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate query configuration
	if p.CustomQuery != "" {
		// Validate custom query for SQL injection
		if err := validation.ValidateCustomQuery(p.CustomQuery, "custom_query"); err != nil {
			return err
		}
	} else {
		// Validate table and column names
		if err := validation.ValidateSQLIdentifier(p.TableName, "table_name"); err != nil {
			return err
		}
		if err := validation.ValidateSQLIdentifier(p.ColumnPath, "column_path"); err != nil {
			return err
		}
		if err := validation.ValidateSQLIdentifier(p.ColumnContent, "column_content"); err != nil {
			return err
		}

		// Validate optional columns if provided
		if p.ColumnType != "" {
			if err := validation.ValidateSQLIdentifier(p.ColumnType, "column_type"); err != nil {
				return err
			}
		}
		if p.ColumnSize != "" {
			if err := validation.ValidateSQLIdentifier(p.ColumnSize, "column_size"); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetDefaults sets default values for optional fields
func (p *Processor) SetDefaults() {
	if p.Host == "" && p.DBType != "sqlite" {
		p.Host = "localhost"
	}

	if p.Port == 0 {
		switch p.DBType {
		case "postgres":
			p.Port = 5432
		case "mysql":
			p.Port = 3306
		}
	}

	if p.SSLMode == "" && p.DBType == "postgres" {
		p.SSLMode = "disable"
	}

	if p.ColumnPath == "" && p.CustomQuery == "" {
		p.ColumnPath = "file_path"
	}

	if p.ColumnContent == "" && p.CustomQuery == "" {
		p.ColumnContent = "content"
	}
}

// TestConnection tests the database connection without keeping it open
func (p *Processor) TestConnection() error {
	if err := p.Connect(); err != nil {
		return err
	}
	return p.Close()
}
