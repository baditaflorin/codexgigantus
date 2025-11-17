// Package csv provides CSV and TSV file processing functionality
package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/baditaflorin/codexgigantus/pkg/utils"
)

// Processor handles CSV/TSV file processing
type Processor struct {
	FilePath      string
	Delimiter     rune
	PathColumn    int
	ContentColumn int
	HasHeader     bool
	Debug         bool
}

// NewProcessor creates a new CSV/TSV processor
func NewProcessor(filePath string, delimiter rune, pathCol, contentCol int, hasHeader, debug bool) *Processor {
	return &Processor{
		FilePath:      filePath,
		Delimiter:     delimiter,
		PathColumn:    pathCol,
		ContentColumn: contentCol,
		HasHeader:     hasHeader,
		Debug:         debug,
	}
}

// Process reads the CSV/TSV file and returns file results
func (p *Processor) Process() ([]utils.FileResult, error) {
	file, err := os.Open(p.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = p.Delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Skip header if present
	startIndex := 0
	if p.HasHeader {
		startIndex = 1
	}

	var results []utils.FileResult

	for i := startIndex; i < len(records); i++ {
		record := records[i]

		// Validate column indices
		if p.PathColumn >= len(record) {
			if p.Debug {
				fmt.Printf("Warning: Path column %d out of range for record %d (has %d columns)\n",
					p.PathColumn, i, len(record))
			}
			continue
		}

		if p.ContentColumn >= len(record) {
			if p.Debug {
				fmt.Printf("Warning: Content column %d out of range for record %d (has %d columns)\n",
					p.ContentColumn, i, len(record))
			}
			continue
		}

		filePath := record[p.PathColumn]
		content := record[p.ContentColumn]

		if filePath == "" {
			if p.Debug {
				fmt.Printf("Warning: Empty file path in record %d\n", i)
			}
			continue
		}

		results = append(results, utils.FileResult{
			Path:    filePath,
			Content: content,
		})

		if p.Debug {
			fmt.Printf("Processed CSV record %d: %s (%d bytes)\n", i, filePath, len(content))
		}
	}

	if p.Debug {
		fmt.Printf("Processed %d records from CSV file\n", len(results))
	}

	return results, nil
}

// Validate validates the processor configuration
func (p *Processor) Validate() error {
	if p.FilePath == "" {
		return fmt.Errorf("file path is required")
	}

	if _, err := os.Stat(p.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", p.FilePath)
	}

	if p.PathColumn < 0 {
		return fmt.Errorf("path column index must be >= 0")
	}

	if p.ContentColumn < 0 {
		return fmt.Errorf("content column index must be >= 0")
	}

	return nil
}
