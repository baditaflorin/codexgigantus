// Package gui provides web GUI functionality for CodexGigantus
package gui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/baditaflorin/codexgigantus/pkg/configfile"
	"github.com/baditaflorin/codexgigantus/pkg/processor"
	"github.com/baditaflorin/codexgigantus/pkg/sources/csv"
	"github.com/baditaflorin/codexgigantus/pkg/sources/database"
	"github.com/baditaflorin/codexgigantus/pkg/utils"
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

// Start starts the web server
func (s *Server) Start(host string, port int) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/config", s.handleConfig)
	http.HandleFunc("/api/config/load", s.handleLoadConfig)
	http.HandleFunc("/api/config/save", s.handleSaveConfig)
	http.HandleFunc("/api/process", s.handleProcess)
	http.HandleFunc("/api/test-db", s.handleTestDB)

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting web GUI on http://%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", s.config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleConfig handles GET/POST for configuration
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.config)

	case http.MethodPost:
		var config configfile.AppConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		config.SetDefaults()
		if err := config.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("Invalid configuration: %v", err), http.StatusBadRequest)
			return
		}

		s.config = &config
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLoadConfig loads configuration from a file
func (s *Server) handleLoadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	config, err := configfile.Load(req.FilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load config: %v", err), http.StatusBadRequest)
		return
	}

	config.SetDefaults()
	s.config = config

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"config": config,
	})
}

// handleSaveConfig saves current configuration to a file
func (s *Server) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := configfile.Save(s.config, req.FilePath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Configuration saved to %s", req.FilePath),
	})
}

// handleProcess processes files based on current configuration
func (s *Server) handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, fmt.Sprintf("Unknown source type: %s", s.config.SourceType), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate output
	output := utils.GenerateOutput(results, s.config.ShowFuncs)

	// Save to file if configured
	if s.config.OutputFile != "" {
		if err := utils.SaveOutput(output, s.config.OutputFile); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save output: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"file_count":  len(results),
		"output_size": len(output),
		"output_file": s.config.OutputFile,
		"output":      output,
	})
}

// handleTestDB tests database connection
func (s *Server) handleTestDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.config.SourceType != "database" {
		http.Error(w, "Current configuration is not for database source", http.StatusBadRequest)
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
		s.config.Debug,
	)
	dbProc.SetDefaults()

	if err := dbProc.TestConnection(); err != nil {
		http.Error(w, fmt.Sprintf("Connection failed: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
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
