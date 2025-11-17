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

// Connect establishes a database connection
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
		return fmt.Errorf("unsupported database type: %s", p.DBType)
	}

	if p.Debug {
		// Mask password in debug output
		safeDSN := strings.ReplaceAll(dsn, p.Password, "****")
		fmt.Printf("Connecting to database: %s\n", safeDSN)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
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
		return nil, fmt.Errorf("database connection not established, call Connect() first")
	}

	query := p.buildQuery()

	if p.Debug {
		fmt.Printf("Executing query: %s\n", query)
	}

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
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
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if p.Debug {
		fmt.Printf("Retrieved %d records from database\n", len(results))
	}

	return results, nil
}

// buildQuery constructs the SQL query
func (p *Processor) buildQuery() string {
	// Use custom query if provided
	if p.CustomQuery != "" {
		return p.CustomQuery
	}

	// Build query from table and column names
	columns := []string{p.ColumnPath, p.ColumnContent}

	if p.ColumnType != "" {
		columns = append(columns, p.ColumnType)
	}
	if p.ColumnSize != "" {
		columns = append(columns, p.ColumnSize)
	}

	return fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(columns, ", "),
		p.TableName)
}

// Validate validates the processor configuration
func (p *Processor) Validate() error {
	if p.DBType == "" {
		return fmt.Errorf("database type is required")
	}

	validTypes := map[string]bool{"postgres": true, "mysql": true, "sqlite": true}
	if !validTypes[p.DBType] {
		return fmt.Errorf("invalid database type: %s (must be postgres, mysql, or sqlite)", p.DBType)
	}

	if p.DBType != "sqlite" {
		if p.Host == "" {
			return fmt.Errorf("host is required for %s", p.DBType)
		}
		if p.Port <= 0 {
			return fmt.Errorf("port must be > 0")
		}
		if p.User == "" {
			return fmt.Errorf("user is required for %s", p.DBType)
		}
	}

	if p.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate query configuration
	if p.CustomQuery == "" {
		// If no custom query, table and columns are required
		if p.TableName == "" {
			return fmt.Errorf("table_name is required when custom_query is not provided")
		}
		if p.ColumnPath == "" {
			return fmt.Errorf("column_path is required when custom_query is not provided")
		}
		if p.ColumnContent == "" {
			return fmt.Errorf("column_content is required when custom_query is not provided")
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
