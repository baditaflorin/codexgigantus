package validation

import (
	"strings"
	"testing"
)

func TestValidateSQLIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldName string
		wantErr   bool
	}{
		{"valid identifier", "users", "table_name", false},
		{"valid with underscore", "user_accounts", "table_name", false},
		{"valid with numbers", "table123", "table_name", false},
		{"empty string", "", "table_name", true},
		{"too long", strings.Repeat("a", MaxTableNameLength+1), "table_name", true},
		{"starts with number", "123table", "table_name", true},
		{"contains spaces", "user table", "table_name", true},
		{"contains dash", "user-table", "table_name", true},
		{"SQL comment", "users--", "table_name", true},
		{"SQL injection attempt", "users; DROP TABLE", "table_name", true},
		{"contains quotes", "users'", "table_name", true},
		{"xp_ prefix", "xp_cmdshell", "table_name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSQLIdentifier(tt.input, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSQLIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeSQLIdentifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"clean input", "users", "users"},
		{"with spaces", "user table", "usertable"},
		{"with dashes", "user-table", "usertable"},
		{"with quotes", "user's", "users"},
		{"SQL injection", "users; DROP TABLE users--", "usersDROPTABLEusers"},
		{"special chars", "user@#$%table", "usertable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeSQLIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeSQLIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldName string
		wantErr   bool
	}{
		{"valid relative path", "dir/file.txt", "file_path", false},
		{"valid absolute path", "/home/user/file.txt", "file_path", false},
		{"current directory", ".", "file_path", false},
		{"empty string", "", "file_path", true},
		{"too long", strings.Repeat("a", MaxPathLength+1), "file_path", true},
		{"path traversal with ..", "../etc/passwd", "file_path", true},
		{"path traversal with ~", "~/secret", "file_path", true},
		{"contains pipe", "file|command", "file_path", true},
		{"contains semicolon", "file;rm -rf", "file_path", true},
		{"contains backtick", "file`whoami`", "file_path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.input, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name      string
		port      int
		fieldName string
		wantErr   bool
	}{
		{"valid port 80", 80, "port", false},
		{"valid port 443", 443, "port", false},
		{"valid port 8080", 8080, "port", false},
		{"valid port 0", 0, "port", false},
		{"valid port 65535", 65535, "port", false},
		{"negative port", -1, "port", true},
		{"port too high", 65536, "port", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		fieldName string
		wantErr   bool
	}{
		{"valid hostname", "localhost", "host", false},
		{"valid domain", "example.com", "host", false},
		{"valid IP", "192.168.1.1", "host", false},
		{"empty host", "", "host", true},
		{"too long", strings.Repeat("a", 256), "host", true},
		{"contains pipe", "localhost|whoami", "host", true},
		{"contains semicolon", "localhost;id", "host", true},
		{"contains backtick", "localhost`whoami`", "host", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHost(tt.host, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDatabaseType(t *testing.T) {
	tests := []struct {
		name      string
		dbType    string
		fieldName string
		wantErr   bool
	}{
		{"postgres", "postgres", "db_type", false},
		{"mysql", "mysql", "db_type", false},
		{"sqlite", "sqlite", "db_type", false},
		{"uppercase postgres", "POSTGRES", "db_type", false},
		{"empty", "", "db_type", true},
		{"invalid type", "mongodb", "db_type", true},
		{"invalid type", "oracle", "db_type", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDatabaseType(tt.dbType, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDatabaseType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSourceType(t *testing.T) {
	tests := []struct {
		name       string
		sourceType string
		fieldName  string
		wantErr    bool
	}{
		{"filesystem", "filesystem", "source_type", false},
		{"csv", "csv", "source_type", false},
		{"tsv", "tsv", "source_type", false},
		{"database", "database", "source_type", false},
		{"uppercase", "FILESYSTEM", "source_type", false},
		{"empty", "", "source_type", true},
		{"invalid", "json", "source_type", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSourceType(tt.sourceType, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCustomQuery(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		fieldName string
		wantErr   bool
	}{
		{"valid SELECT", "SELECT * FROM users", "query", false},
		{"valid with WHERE", "SELECT id, name FROM users WHERE active = true", "query", false},
		{"valid with JOIN", "SELECT u.id, u.name FROM users u JOIN accounts a ON u.id = a.user_id", "query", false},
		{"empty query", "", "query", false},
		{"lowercase select", "select * from users", "query", false},
		{"too long", "SELECT " + strings.Repeat("a", MaxQueryLength), "query", true},
		{"DROP TABLE", "SELECT * FROM users; DROP TABLE users", "query", true},
		{"DELETE", "DELETE FROM users", "query", true},
		{"UPDATE", "UPDATE users SET name = 'hacked'", "query", true},
		{"INSERT", "INSERT INTO users VALUES (1, 'hacker')", "query", true},
		{"EXEC", "EXEC xp_cmdshell 'whoami'", "query", true},
		{"xp_ stored proc", "SELECT * FROM users; EXEC xp_cmdshell", "query", true},
		{"INTO OUTFILE", "SELECT * FROM users INTO OUTFILE '/tmp/users.txt'", "query", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomQuery(tt.query, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCustomQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileExtension(t *testing.T) {
	tests := []struct {
		name      string
		ext       string
		fieldName string
		wantErr   bool
	}{
		{"go extension", "go", "extension", false},
		{"with dot", ".go", "extension", false},
		{"js extension", "js", "extension", false},
		{"empty", "", "extension", false},
		{"too long", "verylongext", "extension", true},
		{"with special char", "go!", "extension", true},
		{"with space", "go ", "extension", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileExtension(tt.ext, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileExtension() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfigName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldName string
		wantErr   bool
	}{
		{"valid name", "My Config", "config_name", false},
		{"with dash", "my-config", "config_name", false},
		{"with underscore", "my_config", "config_name", false},
		{"with numbers", "config123", "config_name", false},
		{"empty", "", "config_name", false},
		{"too long", strings.Repeat("a", MaxConfigNameLength+1), "config_name", true},
		{"with special char", "config@123", "config_name", true},
		{"with slash", "config/123", "config_name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigName(tt.input, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCSVDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		fieldName string
		wantErr   bool
	}{
		{"comma", ",", "delimiter", false},
		{"tab", "\t", "delimiter", false},
		{"semicolon", ";", "delimiter", false},
		{"pipe", "|", "delimiter", false},
		{"empty", "", "delimiter", true},
		{"invalid", ":", "delimiter", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCSVDelimiter(tt.delimiter, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCSVDelimiter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositiveInt(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		wantErr   bool
	}{
		{"positive", 10, "value", false},
		{"zero", 0, "value", false},
		{"negative", -1, "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveInt(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositiveInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Message: "test message",
	}

	expected := "validation error for test_field: test message"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}
