package database

import (
	"strings"
	"testing"
)

// TestDatabase_ConnectionFailures tests various database connection failure scenarios
func TestDatabase_ConnectionFailures(t *testing.T) {
	tests := []struct {
		name     string
		dbType   string
		host     string
		port     int
		dbName   string
		user     string
		password string
		sslMode  string
		wantErr  bool
	}{
		{
			name:     "Invalid host",
			dbType:   "postgres",
			host:     "nonexistent.invalid.host",
			port:     5432,
			dbName:   "testdb",
			user:     "testuser",
			password: "testpass",
			sslMode:  "disable",
			wantErr:  true,
		},
		{
			name:     "Invalid port",
			dbType:   "postgres",
			host:     "localhost",
			port:     99999, // Should fail validation
			dbName:   "testdb",
			user:     "testuser",
			password: "testpass",
			sslMode:  "disable",
			wantErr:  true,
		},
		{
			name:     "Wrong credentials",
			dbType:   "postgres",
			host:     "localhost",
			port:     5432,
			dbName:   "postgres",
			user:     "wrong_user",
			password: "wrong_password",
			sslMode:  "disable",
			wantErr:  true,
		},
		{
			name:     "Empty database name",
			dbType:   "postgres",
			host:     "localhost",
			port:     5432,
			dbName:   "",
			user:     "testuser",
			password: "testpass",
			sslMode:  "disable",
			wantErr:  true,
		},
		{
			name:     "Invalid database type",
			dbType:   "mongodb", // Unsupported
			host:     "localhost",
			port:     27017,
			dbName:   "testdb",
			user:     "testuser",
			password: "testpass",
			sslMode:  "disable",
			wantErr:  true,
		},
		{
			name:     "SQL injection in database name",
			dbType:   "postgres",
			host:     "localhost",
			port:     5432,
			dbName:   "testdb'; DROP TABLE users--",
			user:     "testuser",
			password: "testpass",
			sslMode:  "disable",
			wantErr:  true, // Should fail during connection or validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := NewProcessor(
				tt.dbType,
				tt.host,
				tt.port,
				tt.dbName,
				tt.user,
				tt.password,
				tt.sslMode,
				false,
			)

			// First test validation
			err := proc.Validate()
			if tt.wantErr && err == nil {
				// If we expect an error, validation might catch it
				// Try connection
				err = proc.Connect()
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got error: %v", tt.wantErr, err)
			}

			// Clean up if connection succeeded
			if err == nil {
				proc.Close()
			}
		})
	}
}

// TestDatabase_MaliciousTableNames tests SQL injection in table names
func TestDatabase_MaliciousTableNames(t *testing.T) {
	maliciousNames := []string{
		"users; DROP TABLE code_files--",
		"users' OR '1'='1",
		"users UNION SELECT * FROM passwords",
		"users/**/OR/**/1=1",
		"users--",
		"users;--",
		"users'; DROP TABLE users CASCADE;--",
	}

	for _, tableName := range maliciousNames {
		t.Run("MaliciousTable_"+tableName[:15], func(t *testing.T) {
			proc := NewProcessor(
				"postgres",
				"localhost",
				5432,
				"testdb",
				"testuser",
				"testpass",
				"disable",
				false,
			)

			proc.TableName = tableName
			proc.ColumnPath = "file_path"
			proc.ColumnContent = "content"

			// Should fail validation
			err := proc.Validate()
			if err == nil {
				t.Errorf("Malicious table name should have been rejected: %s", tableName)
			}
		})
	}
}

// TestDatabase_MaliciousColumnNames tests SQL injection in column names
func TestDatabase_MaliciousColumnNames(t *testing.T) {
	maliciousNames := []string{
		"file_path; DROP TABLE users--",
		"content' OR '1'='1",
		"file_path UNION SELECT password FROM users",
		"content--",
	}

	for _, columnName := range maliciousNames {
		t.Run("MaliciousColumn", func(t *testing.T) {
			proc := NewProcessor(
				"postgres",
				"localhost",
				5432,
				"testdb",
				"testuser",
				"testpass",
				"disable",
				false,
			)

			proc.TableName = "code_files"
			proc.ColumnPath = columnName
			proc.ColumnContent = "content"

			// Should fail validation
			err := proc.Validate()
			if err == nil {
				t.Errorf("Malicious column name should have been rejected: %s", columnName)
			}
		})
	}
}

// TestDatabase_QueryBuilding tests that query building produces safe queries
func TestDatabase_QueryBuilding(t *testing.T) {
	tests := []struct {
		name           string
		tableName      string
		columnPath     string
		columnContent  string
		columnType     string
		columnSize     string
		expectError    bool
		expectContains string
	}{
		{
			name:           "Valid basic query",
			tableName:      "code_files",
			columnPath:     "file_path",
			columnContent:  "content",
			columnType:     "",
			columnSize:     "",
			expectError:    false,
			expectContains: "SELECT file_path, content FROM code_files",
		},
		{
			name:           "Valid query with optional columns",
			tableName:      "code_files",
			columnPath:     "file_path",
			columnContent:  "content",
			columnType:     "file_type",
			columnSize:     "file_size",
			expectError:    false,
			expectContains: "SELECT file_path, content, file_type, file_size FROM code_files",
		},
		{
			name:          "Invalid table name with SQL injection",
			tableName:     "code_files; DROP TABLE users",
			columnPath:    "file_path",
			columnContent: "content",
			expectError:   true,
		},
		{
			name:          "Invalid column name with SQL injection",
			tableName:     "code_files",
			columnPath:    "file_path; DELETE FROM code_files",
			columnContent: "content",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := &Processor{
				DBType:        "postgres",
				TableName:     tt.tableName,
				ColumnPath:    tt.columnPath,
				ColumnContent: tt.columnContent,
				ColumnType:    tt.columnType,
				ColumnSize:    tt.columnSize,
			}

			query, err := proc.buildQuery()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. Query: %s", query)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if query != tt.expectContains {
					t.Errorf("Expected query to contain '%s', got '%s'", tt.expectContains, query)
				}
			}
		})
	}
}

// TestDatabase_CustomQueryValidation tests custom query validation
func TestDatabase_CustomQueryValidation(t *testing.T) {
	proc := NewProcessor(
		"postgres",
		"localhost",
		5432,
		"testdb",
		"testuser",
		"testpass",
		"disable",
		false,
	)

	dangerousQueries := []string{
		"DROP TABLE code_files",
		"DELETE FROM code_files",
		"UPDATE code_files SET content = 'hacked'",
		"INSERT INTO code_files VALUES ('hack', 'hack')",
		"SELECT * FROM code_files; DROP TABLE users",
		"SELECT * FROM code_files INTO OUTFILE '/tmp/dump.txt'",
		"SELECT * FROM code_files; EXEC xp_cmdshell 'whoami'",
	}

	for _, query := range dangerousQueries {
		t.Run("DangerousQuery", func(t *testing.T) {
			proc.CustomQuery = query
			proc.TableName = "" // Clear table name since we're using custom query

			err := proc.Validate()
			if err == nil {
				t.Errorf("Dangerous query should have been rejected: %s", query)
			}
		})
	}

	// Test valid SELECT queries
	validQueries := []string{
		"SELECT file_path, content FROM code_files",
		"SELECT * FROM code_files WHERE file_type = 'go'",
		"SELECT file_path, content FROM code_files LIMIT 100",
		"SELECT f.file_path, f.content FROM code_files f JOIN metadata m ON f.id = m.file_id",
	}

	for _, query := range validQueries {
		t.Run("ValidQuery", func(t *testing.T) {
			proc.CustomQuery = query
			err := proc.Validate()
			if err != nil {
				t.Errorf("Valid query should have been accepted: %s. Error: %v", query, err)
			}
		})
	}
}

// TestDatabase_ConnectionPoolLimits tests that connection pools are properly limited
func TestDatabase_ConnectionPoolLimits(t *testing.T) {
	proc := NewProcessor(
		"sqlite",
		"",
		0,
		":memory:",
		"",
		"",
		"",
		false,
	)

	proc.SetDefaults()

	err := proc.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to in-memory SQLite: %v", err)
	}
	defer proc.Close()

	// Check that connection pool limits are set
	stats := proc.db.Stats()
	if stats.MaxOpenConnections != 25 {
		t.Errorf("Expected MaxOpenConnections to be 25, got %d", stats.MaxOpenConnections)
	}
}

// TestDatabase_ErrorMessages tests that error messages don't leak sensitive information
func TestDatabase_ErrorMessages(t *testing.T) {
	proc := NewProcessor(
		"postgres",
		"nonexistent.host",
		5432,
		"testdb",
		"testuser",
		"testpass",
		"disable",
		false,
	)

	err := proc.Connect()
	if err == nil {
		t.Fatal("Expected connection to fail")
	}

	// Error message should not contain password
	errorMsg := err.Error()
	if contains(errorMsg, "testpass") {
		t.Errorf("Error message contains password: %s", errorMsg)
	}

	// Error message should be generic
	if !contains(errorMsg, "failed to establish database connection") && !contains(errorMsg, "database connection test failed") {
		t.Errorf("Error message should be generic, got: %s", errorMsg)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
