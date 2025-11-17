package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	p := NewProcessor("postgres", "localhost", 5432, "testdb", "user", "pass", "disable", false)

	if p.DBType != "postgres" {
		t.Errorf("DBType = %v, want postgres", p.DBType)
	}
	if p.Host != "localhost" {
		t.Errorf("Host = %v, want localhost", p.Host)
	}
	if p.Port != 5432 {
		t.Errorf("Port = %v, want 5432", p.Port)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		processor *Processor
		wantErr   bool
	}{
		{
			name: "valid postgres config",
			processor: &Processor{
				DBType:        "postgres",
				Host:          "localhost",
				Port:          5432,
				DBName:        "test",
				User:          "user",
				TableName:     "files",
				ColumnPath:    "path",
				ColumnContent: "content",
			},
			wantErr: false,
		},
		{
			name: "valid sqlite config",
			processor: &Processor{
				DBType:        "sqlite",
				DBName:        "test.db",
				TableName:     "files",
				ColumnPath:    "path",
				ColumnContent: "content",
			},
			wantErr: false,
		},
		{
			name: "valid with custom query",
			processor: &Processor{
				DBType:      "postgres",
				Host:        "localhost",
				Port:        5432,
				DBName:      "test",
				User:        "user",
				CustomQuery: "SELECT path, content FROM files",
			},
			wantErr: false,
		},
		{
			name:      "missing db type",
			processor: &Processor{},
			wantErr:   true,
		},
		{
			name: "invalid db type",
			processor: &Processor{
				DBType: "invalid",
			},
			wantErr: true,
		},
		{
			name: "postgres missing host",
			processor: &Processor{
				DBType: "postgres",
				Port:   5432,
				DBName: "test",
				User:   "user",
			},
			wantErr: true,
		},
		{
			name: "postgres missing user",
			processor: &Processor{
				DBType: "postgres",
				Host:   "localhost",
				Port:   5432,
				DBName: "test",
			},
			wantErr: true,
		},
		{
			name: "missing table name without custom query",
			processor: &Processor{
				DBType:        "sqlite",
				DBName:        "test.db",
				ColumnPath:    "path",
				ColumnContent: "content",
			},
			wantErr: true,
		},
		{
			name: "missing columns without custom query",
			processor: &Processor{
				DBType:    "sqlite",
				DBName:    "test.db",
				TableName: "files",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.processor.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	tests := []struct {
		name   string
		before *Processor
		check  func(*Processor) error
	}{
		{
			name: "sets postgres defaults",
			before: &Processor{
				DBType: "postgres",
			},
			check: func(p *Processor) error {
				if p.Host != "localhost" {
					t.Errorf("Host = %v, want localhost", p.Host)
				}
				if p.Port != 5432 {
					t.Errorf("Port = %v, want 5432", p.Port)
				}
				if p.SSLMode != "disable" {
					t.Errorf("SSLMode = %v, want disable", p.SSLMode)
				}
				return nil
			},
		},
		{
			name: "sets mysql defaults",
			before: &Processor{
				DBType: "mysql",
			},
			check: func(p *Processor) error {
				if p.Port != 3306 {
					t.Errorf("Port = %v, want 3306", p.Port)
				}
				return nil
			},
		},
		{
			name: "sets column defaults",
			before: &Processor{
				DBType: "postgres",
			},
			check: func(p *Processor) error {
				if p.ColumnPath != "file_path" {
					t.Errorf("ColumnPath = %v, want file_path", p.ColumnPath)
				}
				if p.ColumnContent != "content" {
					t.Errorf("ColumnContent = %v, want content", p.ColumnContent)
				}
				return nil
			},
		},
		{
			name: "doesn't override existing values",
			before: &Processor{
				DBType:     "postgres",
				Host:       "custom.host",
				Port:       9999,
				ColumnPath: "custom_path",
			},
			check: func(p *Processor) error {
				if p.Host != "custom.host" {
					t.Errorf("Host = %v, want custom.host", p.Host)
				}
				if p.Port != 9999 {
					t.Errorf("Port = %v, want 9999", p.Port)
				}
				if p.ColumnPath != "custom_path" {
					t.Errorf("ColumnPath = %v, want custom_path", p.ColumnPath)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before.SetDefaults()
			tt.check(tt.before)
		})
	}
}

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name      string
		processor *Processor
		want      string
	}{
		{
			name: "basic query",
			processor: &Processor{
				TableName:     "files",
				ColumnPath:    "path",
				ColumnContent: "content",
			},
			want: "SELECT path, content FROM files",
		},
		{
			name: "query with type column",
			processor: &Processor{
				TableName:     "files",
				ColumnPath:    "path",
				ColumnContent: "content",
				ColumnType:    "file_type",
			},
			want: "SELECT path, content, file_type FROM files",
		},
		{
			name: "query with all columns",
			processor: &Processor{
				TableName:     "code_files",
				ColumnPath:    "file_path",
				ColumnContent: "file_content",
				ColumnType:    "type",
				ColumnSize:    "size",
			},
			want: "SELECT file_path, file_content, type, size FROM code_files",
		},
		{
			name: "custom query",
			processor: &Processor{
				CustomQuery: "SELECT * FROM files WHERE type = 'go'",
			},
			want: "SELECT * FROM files WHERE type = 'go'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.buildQuery()
			if got != tt.want {
				t.Errorf("buildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestProcessWithSQLite tests the Process method using an in-memory SQLite database
func TestProcessWithSQLite(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create and setup SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite database: %v", err)
	}

	// Create table
	_, err = db.Exec(`CREATE TABLE files (
		path TEXT,
		content TEXT
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO files (path, content) VALUES
		('file1.go', 'package main'),
		('file2.py', 'import sys'),
		('file3.js', 'console.log("test")')`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	db.Close()

	// Create processor
	p := &Processor{
		DBType:        "sqlite",
		DBName:        dbPath,
		TableName:     "files",
		ColumnPath:    "path",
		ColumnContent: "content",
		Debug:         false,
	}

	// Connect
	err = p.Connect()
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer p.Close()

	// Process
	results, err := p.Process()
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Verify results
	if len(results) != 3 {
		t.Errorf("Process() returned %d results, want 3", len(results))
	}

	if results[0].Path != "file1.go" {
		t.Errorf("results[0].Path = %v, want file1.go", results[0].Path)
	}
	if results[0].Content != "package main" {
		t.Errorf("results[0].Content = %v, want 'package main'", results[0].Content)
	}
}

func TestConnectInvalidDatabase(t *testing.T) {
	p := &Processor{
		DBType: "postgres",
		Host:   "invalid-host-12345",
		Port:   5432,
		DBName: "test",
		User:   "user",
		Password: "pass",
		SSLMode: "disable",
	}

	err := p.Connect()
	if err == nil {
		t.Error("Connect() should fail for invalid host")
		p.Close()
	}
}

func TestProcessWithoutConnection(t *testing.T) {
	p := &Processor{
		DBType:        "sqlite",
		DBName:        "test.db",
		TableName:     "files",
		ColumnPath:    "path",
		ColumnContent: "content",
	}

	_, err := p.Process()
	if err == nil {
		t.Error("Process() should fail when database is not connected")
	}
}

func TestTestConnection(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping database test in CI environment")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a simple SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	db.Close()

	p := &Processor{
		DBType: "sqlite",
		DBName: dbPath,
	}

	err = p.TestConnection()
	if err != nil {
		t.Errorf("TestConnection() error = %v", err)
	}

	// Verify connection is closed
	if p.db != nil {
		t.Error("Connection should be closed after TestConnection()")
	}
}
