package csv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	p := NewProcessor("test.csv", ',', 0, 1, true, false)

	if p.FilePath != "test.csv" {
		t.Errorf("FilePath = %v, want test.csv", p.FilePath)
	}
	if p.Delimiter != ',' {
		t.Errorf("Delimiter = %v, want ,", p.Delimiter)
	}
	if p.PathColumn != 0 {
		t.Errorf("PathColumn = %v, want 0", p.PathColumn)
	}
	if p.ContentColumn != 1 {
		t.Errorf("ContentColumn = %v, want 1", p.ContentColumn)
	}
	if p.HasHeader != true {
		t.Errorf("HasHeader = %v, want true", p.HasHeader)
	}
}

func TestProcessCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	// Create test CSV file
	content := `path,content
file1.go,package main
file2.py,import sys
file3.js,console.log("test")`

	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	// Process with header
	p := NewProcessor(csvFile, ',', 0, 1, true, false)
	results, err := p.Process()
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Process() returned %d results, want 3", len(results))
	}

	// Verify first result
	if results[0].Path != "file1.go" {
		t.Errorf("results[0].Path = %v, want file1.go", results[0].Path)
	}
	if results[0].Content != "package main" {
		t.Errorf("results[0].Content = %v, want 'package main'", results[0].Content)
	}
}

func TestProcessTSV(t *testing.T) {
	tmpDir := t.TempDir()
	tsvFile := filepath.Join(tmpDir, "test.tsv")

	// Create test TSV file (tab-separated)
	content := "file1.go\tpackage main\nfile2.py\timport sys"

	if err := os.WriteFile(tsvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test TSV: %v", err)
	}

	// Process without header
	p := NewProcessor(tsvFile, '\t', 0, 1, false, false)
	results, err := p.Process()
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Process() returned %d results, want 2", len(results))
	}

	// Verify results
	if results[0].Path != "file1.go" {
		t.Errorf("results[0].Path = %v, want file1.go", results[0].Path)
	}
}

func TestProcessEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "empty.csv")

	if err := os.WriteFile(csvFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	p := NewProcessor(csvFile, ',', 0, 1, false, false)
	_, err := p.Process()
	if err == nil {
		t.Error("Process() should fail for empty CSV")
	}
}

func TestProcessInvalidColumns(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "invalid.csv")

	// CSV with only 2 columns but we request column 5
	content := `col1,col2
val1,val2`

	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	p := NewProcessor(csvFile, ',', 5, 1, true, false)
	results, err := p.Process()
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Should skip invalid records and return empty
	if len(results) != 0 {
		t.Errorf("Process() returned %d results, want 0", len(results))
	}
}

func TestProcessWithEmptyPaths(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "empty_paths.csv")

	// CSV with some empty paths
	content := `path,content
,content1
file2.txt,content2
,content3`

	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	p := NewProcessor(csvFile, ',', 0, 1, true, false)
	results, err := p.Process()
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Should only return records with non-empty paths
	if len(results) != 1 {
		t.Errorf("Process() returned %d results, want 1", len(results))
	}
	if results[0].Path != "file2.txt" {
		t.Errorf("results[0].Path = %v, want file2.txt", results[0].Path)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		processor *Processor
		wantErr   bool
		skipFile  bool
	}{
		{
			name: "valid processor",
			processor: &Processor{
				FilePath:      "test.csv",
				PathColumn:    0,
				ContentColumn: 1,
			},
			wantErr:  true, // File doesn't exist, so validation should fail
			skipFile: false,
		},
		{
			name: "empty file path",
			processor: &Processor{
				FilePath:      "",
				PathColumn:    0,
				ContentColumn: 1,
			},
			wantErr: true,
		},
		{
			name: "negative path column",
			processor: &Processor{
				FilePath:      "test.csv",
				PathColumn:    -1,
				ContentColumn: 1,
			},
			wantErr: true,
		},
		{
			name: "negative content column",
			processor: &Processor{
				FilePath:      "test.csv",
				PathColumn:    0,
				ContentColumn: -1,
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

func TestValidateNonExistentFile(t *testing.T) {
	p := &Processor{
		FilePath:      "/nonexistent/file.csv",
		PathColumn:    0,
		ContentColumn: 1,
	}

	err := p.Validate()
	if err == nil {
		t.Error("Validate() should fail for non-existent file")
	}
}
