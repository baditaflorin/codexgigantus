package validation

import (
	"fmt"
	"testing"
)

// TestSQL_InjectionAttempts tests various SQL injection attack patterns
func TestSQL_InjectionAttempts(t *testing.T) {
	injectionAttempts := []string{
		"users; DROP TABLE users--",
		"users' OR '1'='1",
		"users\"; DROP TABLE users--",
		"1' OR '1'='1'; --",
		"admin'--",
		"' OR 1=1--",
		"'; EXEC xp_cmdshell('dir'); --",
		"1'; DROP TABLE code_files--",
		"users UNION SELECT * FROM passwords",
		"users/**/OR/**/1=1",
		"users%27%20OR%201=1--",
		"1' AND (SELECT * FROM (SELECT(SLEEP(5)))a)--",
		"' OR 'x'='x",
		"1'; WAITFOR DELAY '00:00:05'--",
		"users' AND 1=(SELECT COUNT(*) FROM tabname)--",
	}

	for i, attempt := range injectionAttempts {
		t.Run(fmt.Sprintf("SQLInjection_%d", i), func(t *testing.T) {
			err := ValidateSQLIdentifier(attempt, "test_field")
			if err == nil {
				t.Errorf("SQL injection attempt should have been blocked: %s", attempt)
			}
		})
	}
}

// TestPathTraversalAttempts tests various path traversal attack patterns
func TestPathTraversalAttempts(t *testing.T) {
	pathTraversalAttempts := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"....//....//....//etc/passwd",
		"..%2F..%2F..%2Fetc%2Fpasswd",
		"....\\\\....\\\\....\\\\windows\\\\system32",
		"/etc/passwd",
		"../../../../../../etc/shadow",
		"../../../../../../../../../../../etc/passwd",
		"file.txt;rm -rf /",
		"file.txt|cat /etc/passwd",
		"file.txt`whoami`",
		"file.txt$(whoami)",
		"file.txt&&whoami",
		"~/../../etc/passwd",
		"$HOME/../../etc/passwd",
	}

	for i, attempt := range pathTraversalAttempts {
		t.Run(fmt.Sprintf("PathTraversal_%d", i), func(t *testing.T) {
			err := ValidateFilePath(attempt, "test_path")
			if err == nil {
				t.Errorf("Path traversal attempt should have been blocked: %s", attempt)
			}
		})
	}
}

// TestCommandInjectionInHosts tests command injection in host fields
func TestCommandInjection_InHosts(t *testing.T) {
	commandInjectionAttempts := []string{
		"localhost; whoami",
		"localhost | cat /etc/passwd",
		"localhost && whoami",
		"localhost `whoami`",
		"localhost $(whoami)",
		"localhost;rm -rf /",
		"localhost||whoami",
		"127.0.0.1;id",
		"127.0.0.1|ls",
	}

	for i, attempt := range commandInjectionAttempts {
		t.Run(fmt.Sprintf("CommandInjection_%d", i), func(t *testing.T) {
			err := ValidateHost(attempt, "host")
			if err == nil {
				t.Errorf("Command injection attempt should have been blocked: %s", attempt)
			}
		})
	}
}

// TestMalformedDatabaseQueries tests malformed and malicious custom queries
func TestMalformedDatabaseQueries(t *testing.T) {
	maliciousQueries := []string{
		"DROP TABLE code_files",
		"DELETE FROM code_files",
		"UPDATE code_files SET content = 'hacked'",
		"INSERT INTO code_files VALUES ('hack', 'hack')",
		"SELECT * FROM code_files; DROP TABLE code_files",
		"SELECT * FROM code_files WHERE 1=1; DELETE FROM code_files",
		"SELECT * FROM code_files UNION SELECT * FROM users",
		"SELECT * FROM code_files INTO OUTFILE '/tmp/hack.txt'",
		"SELECT LOAD_FILE('/etc/passwd')",
		"EXEC xp_cmdshell('whoami')",
		"SELECT * FROM code_files; EXEC xp_cmdshell 'dir'",
		"SELECT * FROM code_files WHERE id = 1 OR 1=1",
		"SELECT pg_sleep(10)",
		"ALTER TABLE code_files ADD COLUMN hacked VARCHAR(100)",
		"CREATE TABLE hacked (id INT)",
	}

	for _, query := range maliciousQueries {
		t.Run("MaliciousQuery", func(t *testing.T) {
			err := ValidateCustomQuery(query, "custom_query")
			if err == nil {
				t.Errorf("Malicious query should have been blocked: %s", query)
			}
		})
	}
}

// TestOversizedInputs tests validation of oversized inputs
func TestOversizedInputs(t *testing.T) {
	tests := []struct {
		name      string
		validator func(string, string) error
		field     string
	}{
		{
			name: "Oversized table name",
			validator: func(input, field string) error {
				return ValidateSQLIdentifier(input, field)
			},
			field: "table_name",
		},
		{
			name: "Oversized host name",
			validator: func(input, field string) error {
				return ValidateHost(input, field)
			},
			field: "host",
		},
		{
			name: "Oversized config name",
			validator: func(input, field string) error {
				return ValidateConfigName(input, field)
			},
			field: "config_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create oversized input (10KB of 'a's)
			bytes := make([]byte, 10*1024)
			for i := range bytes {
				bytes[i] = 'a'
			}
			oversizedInput := string(bytes)

			err := tt.validator(oversizedInput, tt.field)
			if err == nil {
				t.Errorf("%s: oversized input should have been rejected", tt.name)
			}
		})
	}
}

// TestUnicodeAndSpecialCharacters tests handling of unicode and special characters
func TestUnicodeAndSpecialCharacters(t *testing.T) {
	specialInputs := []string{
		"table\x00name", // Null byte
		"table\nname",   // Newline
		"table\rname",   // Carriage return
		"table\tname",   // Tab
		"table\u0000name", // Unicode null
		"table😀name",     // Emoji
		"table中文name",     // Chinese characters
		"table<script>alert(1)</script>", // XSS attempt
		"table&lt;script&gt;",            // HTML entities
		"table%3Cscript%3E",              // URL encoded
	}

	for i, input := range specialInputs {
		t.Run(fmt.Sprintf("SpecialChar_%d", i), func(t *testing.T) {
			err := ValidateSQLIdentifier(input, "test_field")
			// Should fail because they contain non-alphanumeric characters
			if err == nil {
				t.Errorf("Special character input should have been rejected: %s", input)
			}
		})
	}
}

// TestBoundaryValues tests boundary conditions
func TestBoundaryValues(t *testing.T) {
	t.Run("Port boundary values", func(t *testing.T) {
		tests := []struct {
			port    int
			wantErr bool
		}{
			{-1, true},
			{0, false},
			{1, false},
			{65535, false},
			{65536, true},
			{99999, true},
		}

		for _, tt := range tests {
			err := ValidatePort(tt.port, "port")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		}
	})

	t.Run("String length boundaries", func(t *testing.T) {
		// Create a valid identifier at maximum length
		bytes := make([]byte, MaxTableNameLength)
		bytes[0] = 'a' // Must start with letter
		for i := 1; i < len(bytes); i++ {
			bytes[i] = 'a'
		}
		validMaxLength := string(bytes)

		err := ValidateSQLIdentifier(validMaxLength, "table_name")
		if err != nil {
			t.Errorf("Max length valid identifier should be accepted, got error: %v", err)
		}

		// One char over should fail
		tooLong := validMaxLength + "a"
		err = ValidateSQLIdentifier(tooLong, "table_name")
		if err == nil {
			t.Errorf("Identifier longer than max should be rejected")
		}
	})
}

// TestCSRFAndXSSVectors tests for cross-site scripting patterns
func TestCSRFAndXSSVectors(t *testing.T) {
	xssVectors := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg/onload=alert('XSS')>",
		"<iframe src=javascript:alert('XSS')>",
		"<body onload=alert('XSS')>",
		"<input onfocus=alert('XSS') autofocus>",
		"<select onfocus=alert('XSS') autofocus>",
		"<textarea onfocus=alert('XSS') autofocus>",
		"<marquee onstart=alert('XSS')>",
	}

	for _, vector := range xssVectors {
		t.Run("XSS_Vector", func(t *testing.T) {
			// These should fail validation for config names
			err := ValidateConfigName(vector, "config_name")
			if err == nil {
				t.Errorf("XSS vector should have been rejected: %s", vector)
			}

			// Should also fail for SQL identifiers
			err = ValidateSQLIdentifier(vector, "table_name")
			if err == nil {
				t.Errorf("XSS vector should have been rejected in SQL identifier: %s", vector)
			}
		})
	}
}

// TestLDAPInjection tests LDAP injection patterns (relevant for future auth)
func TestLDAPInjection(t *testing.T) {
	ldapInjections := []string{
		"*",
		"*)(uid=*",
		"admin)(&(password=*",
		"*)(objectClass=*",
		"admin)(|(password=*",
	}

	for _, injection := range ldapInjections {
		t.Run("LDAP_Injection", func(t *testing.T) {
			// Should fail validation for identifiers
			err := ValidateSQLIdentifier(injection, "username")
			if err == nil {
				t.Errorf("LDAP injection should have been rejected: %s", injection)
			}
		})
	}
}

// BenchmarkValidation benchmarks the validation functions
func BenchmarkValidateSQLIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateSQLIdentifier("users", "table_name")
	}
}

func BenchmarkValidateFilePath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateFilePath("/home/user/file.txt", "file_path")
	}
}

func BenchmarkValidateHost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateHost("localhost", "host")
	}
}
