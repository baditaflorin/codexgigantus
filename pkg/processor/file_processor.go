// Package processor handles file system traversal and file processing.
// It walks through directories, applies filters, and reads file contents.
package processor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/baditaflorin/codexgigantus/pkg/config"
	"github.com/baditaflorin/codexgigantus/pkg/utils"
)

// ProcessFiles walks through directories specified in the configuration,
// applies filters (ignore/include rules), and reads the contents of matching files.
// It returns a slice of FileResult containing the path and content of each processed file.
func ProcessFiles(cfg *config.Config) ([]utils.FileResult, error) {
	var results []utils.FileResult

	for _, dir := range cfg.Dirs {
		if cfg.Debug {
			utils.Debug("Processing directory: %s", dir)
		}
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Handle directories
			if info.IsDir() {
				if shouldIgnoreDir(path, cfg) {
					if cfg.Debug {
						utils.Debug("Ignoring directory: %s", path)
					}
					return filepath.SkipDir
				}
				if !cfg.Recursive && path != dir {
					return filepath.SkipDir
				}
				return nil
			}

			// Handle files
			if shouldIgnoreFile(path, cfg) {
				if cfg.Debug {
					utils.Debug("Ignoring file: %s", path)
				}
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			results = append(results, utils.FileResult{
				Path:    path,
				Content: string(content),
			})

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// shouldIgnoreDir checks if a directory should be ignored based on the configuration.
// It returns true if the directory path contains any of the ignore patterns.
func shouldIgnoreDir(path string, cfg *config.Config) bool {
	for _, ignoreDir := range cfg.IgnoreDirs {
		if strings.Contains(path, ignoreDir) {
			return true
		}
	}
	return false
}

// shouldIgnoreFile determines if a file should be ignored based on the configuration.
// It checks the filename, extension, and include/exclude rules.
// Returns true if the file should be skipped.
func shouldIgnoreFile(path string, cfg *config.Config) bool {
	filename := filepath.Base(path)
	ext := strings.TrimPrefix(filepath.Ext(path), ".")

	for _, ignoreFile := range cfg.IgnoreFiles {
		if filename == ignoreFile {
			return true
		}
	}

	if len(cfg.IncludeExts) > 0 {
		include := false
		for _, includeExt := range cfg.IncludeExts {
			if ext == includeExt {
				include = true
				break
			}
		}
		if !include {
			return true
		}
	}

	for _, ignoreExt := range cfg.IgnoreExts {
		if ext == ignoreExt {
			return true
		}
	}

	return false
}
