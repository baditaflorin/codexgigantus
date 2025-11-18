// Package gui provides web GUI functionality for CodexGigantus
package gui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/baditaflorin/codexgigantus/pkg/configfile"
	"github.com/baditaflorin/codexgigantus/pkg/processor"
	"github.com/baditaflorin/codexgigantus/pkg/sources/csv"
	"github.com/baditaflorin/codexgigantus/pkg/sources/database"
	"github.com/baditaflorin/codexgigantus/pkg/utils"
	"github.com/baditaflorin/codexgigantus/pkg/validation"
)

const (
	// Security: Limit request body size to prevent DoS
	maxRequestBodySize = 10 * 1024 * 1024 // 10MB
	maxConfigFileSize  = 1 * 1024 * 1024  // 1MB for config files
)

// Server represents the web GUI server
type Server struct {
	templates *template.Template
	config    *configfile.AppConfig
}

// NewServer creates a new GUI server
func NewServer() (*Server, error) {
	// Parse templates
	tmpl, err := template.ParseGlob(filepath.Join("internal", "gui", "templates", "*.html"))
	if err != nil {
		// If templates don't exist, create embedded ones
		tmpl = template.New("embedded")
		tmpl, err = tmpl.Parse(indexTemplate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse templates: %w", err)
		}
	}

	return &Server{
		templates: tmpl,
		config:    configfile.NewDefault(),
	}, nil
}

// Start starts the web server with security middleware
func (s *Server) Start(host string, port int) error {
	// Wrap handlers with security middleware
	http.HandleFunc("/", s.withSecurityHeaders(s.handleIndex))
	http.HandleFunc("/api/config", s.withSecurityHeaders(s.handleConfig))
	http.HandleFunc("/api/config/load", s.withSecurityHeaders(s.handleLoadConfig))
	http.HandleFunc("/api/config/save", s.withSecurityHeaders(s.handleSaveConfig))
	http.HandleFunc("/api/process", s.withSecurityHeaders(s.handleProcess))
	http.HandleFunc("/api/test-db", s.withSecurityHeaders(s.handleTestDB))

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting web GUI on http://%s\n", addr)
	fmt.Println("Security features enabled: request size limits, input validation, secure headers")
	return http.ListenAndServe(addr, nil)
}

// withSecurityHeaders adds security headers to responses
func (s *Server) withSecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Note: CSP allows 'unsafe-inline' for scripts/styles which reduces XSS protection.
		// This is currently required for the application to function with inline scripts/styles.
		// Consider refactoring to use nonces or hashes for better security.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Limit request body size to prevent DoS
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

		next(w, r)
	}
}

// sendError sends a secure error response without leaking details
func sendError(w http.ResponseWriter, userMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": userMessage,
	})
}

// sendSuccess sends a success response
func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", s.config); err != nil {
		// Security: Don't leak template error details
		sendError(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// handleConfig handles GET/POST for configuration
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sendSuccess(w, s.config)

	case http.MethodPost:
		var config configfile.AppConfig

		// Limit JSON decoding
		decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodySize))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&config); err != nil {
			sendError(w, "Invalid configuration format", http.StatusBadRequest)
			return
		}

		config.SetDefaults()
		if err := config.Validate(); err != nil {
			// Security: Validation errors are safe to return
			sendError(w, fmt.Sprintf("Configuration validation failed: %v", err), http.StatusBadRequest)
			return
		}

		s.config = &config
		sendSuccess(w, map[string]string{"status": "success", "message": "Configuration updated"})

	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLoadConfig loads configuration from a file
func (s *Server) handleLoadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}

	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodySize))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Security: Validate file path to prevent path traversal
	if err := validation.ValidateFilePath(req.FilePath, "file_path"); err != nil {
		sendError(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Security: Only allow loading from configs directory
	cleanPath := filepath.Clean(req.FilePath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join("configs", cleanPath)
	}

	config, err := configfile.Load(cleanPath)
	if err != nil {
		// Security: Don't leak file system details
		sendError(w, "Failed to load configuration file", http.StatusBadRequest)
		return
	}

	config.SetDefaults()
	if err := config.Validate(); err != nil {
		sendError(w, fmt.Sprintf("Loaded configuration is invalid: %v", err), http.StatusBadRequest)
		return
	}

	s.config = config
	sendSuccess(w, map[string]interface{}{
		"status":  "success",
		"message": "Configuration loaded successfully",
		"config":  config,
	})
}

// handleSaveConfig saves current configuration to a file
func (s *Server) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}

	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodySize))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Security: Validate file path to prevent path traversal
	if err := validation.ValidateFilePath(req.FilePath, "file_path"); err != nil {
		sendError(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Security: Only allow saving to configs directory
	cleanPath := filepath.Clean(req.FilePath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join("configs", cleanPath)
	}

	// Security: Validate file extension
	ext := filepath.Ext(cleanPath)
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		sendError(w, "Invalid file format (must be .json, .yaml, or .yml)", http.StatusBadRequest)
		return
	}

	if err := configfile.Save(s.config, cleanPath); err != nil {
		// Security: Don't leak file system details
		sendError(w, "Failed to save configuration file", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, map[string]string{
		"status":  "success",
		"message": "Configuration saved successfully",
	})
}

// handleProcess processes files based on current configuration
func (s *Server) handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Security: Validate configuration before processing
	if err := s.config.Validate(); err != nil {
		sendError(w, fmt.Sprintf("Invalid configuration: %v", err), http.StatusBadRequest)
		return
	}

	var results []utils.FileResult
	var err error

	switch s.config.SourceType {
	case "filesystem":
		results, err = s.processFilesystem()
	case "csv", "tsv":
		results, err = s.processCSV()
	case "database":
		results, err = s.processDatabase()
	default:
		sendError(w, "Invalid source type", http.StatusBadRequest)
		return
	}

	if err != nil {
		// Security: Don't leak internal error details
		sendError(w, "Processing failed", http.StatusInternalServerError)
		return
	}

	// Generate output
	output := utils.GenerateOutput(results, s.config.ShowFuncs)

	// Save to file if configured
	if s.config.OutputFile != "" {
		// Security: Validate output file path
		if err := validation.ValidateFilePath(s.config.OutputFile, "output_file"); err != nil {
			sendError(w, "Invalid output file path", http.StatusBadRequest)
			return
		}

		if err := utils.SaveOutput(output, s.config.OutputFile); err != nil {
			sendError(w, "Failed to save output", http.StatusInternalServerError)
			return
		}
	}

	sendSuccess(w, map[string]interface{}{
		"status":      "success",
		"file_count":  len(results),
		"output_size": len(output),
		"output_file": s.config.OutputFile,
		"preview":     truncateOutput(output, 1000), // Security: Limit response size
	})
}

// truncateOutput limits output size for API responses
func truncateOutput(output string, maxLen int) string {
	if len(output) <= maxLen {
		return output
	}
	return output[:maxLen] + "... (truncated)"
}

// handleTestDB tests database connection
func (s *Server) handleTestDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.config.SourceType != "database" {
		sendError(w, "Configuration is not set for database source", http.StatusBadRequest)
		return
	}

	dbProc := database.NewProcessor(
		s.config.DBType,
		s.config.DBHost,
		s.config.DBPort,
		s.config.DBName,
		s.config.DBUser,
		s.config.DBPassword,
		s.config.DBSSLMode,
		false, // Security: Disable debug for connection test
	)
	dbProc.SetDefaults()

	// Security: Validate before attempting connection
	if err := dbProc.Validate(); err != nil {
		sendError(w, fmt.Sprintf("Invalid database configuration: %v", err), http.StatusBadRequest)
		return
	}

	if err := dbProc.TestConnection(); err != nil {
		// Security: Don't leak connection details
		sendError(w, "Database connection failed", http.StatusBadRequest)
		return
	}

	sendSuccess(w, map[string]string{
		"status":  "success",
		"message": "Database connection successful",
	})
}

// processFilesystem processes files from filesystem
func (s *Server) processFilesystem() ([]utils.FileResult, error) {
	config := processor.Config{
		Dirs:        s.config.Directories,
		IgnoreFiles: s.config.IgnoreFiles,
		IgnoreDirs:  s.config.IgnoreDirs,
		IgnoreExts:  s.config.ExcludeExtensions,
		IncludeExts: s.config.IncludeExtensions,
		Recursive:   s.config.Recursive,
		Debug:       s.config.Debug,
	}

	return processor.ProcessFiles(&config)
}

// processCSV processes files from CSV/TSV
func (s *Server) processCSV() ([]utils.FileResult, error) {
	delimiter := rune(',')
	if s.config.CSVDelimiter != "" {
		delimiter = rune(s.config.CSVDelimiter[0])
	}

	proc := csv.NewProcessor(
		s.config.CSVFilePath,
		delimiter,
		s.config.CSVPathColumn,
		s.config.CSVContentColumn,
		s.config.CSVHasHeader,
		s.config.Debug,
	)

	if err := proc.Validate(); err != nil {
		return nil, err
	}

	return proc.Process()
}

// processDatabase processes files from database
func (s *Server) processDatabase() ([]utils.FileResult, error) {
	dbProc := database.NewProcessor(
		s.config.DBType,
		s.config.DBHost,
		s.config.DBPort,
		s.config.DBName,
		s.config.DBUser,
		s.config.DBPassword,
		s.config.DBSSLMode,
		s.config.Debug,
	)

	dbProc.TableName = s.config.DBTableName
	dbProc.ColumnPath = s.config.DBColumnPath
	dbProc.ColumnContent = s.config.DBColumnContent
	dbProc.ColumnType = s.config.DBColumnType
	dbProc.ColumnSize = s.config.DBColumnSize
	dbProc.CustomQuery = s.config.DBQuery

	dbProc.SetDefaults()

	if err := dbProc.Validate(); err != nil {
		return nil, err
	}

	if err := dbProc.Connect(); err != nil {
		return nil, err
	}
	defer dbProc.Close()

	return dbProc.Process()
}
