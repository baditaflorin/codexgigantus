// file_processor.go
package processor

import (
	"os"
	"path/filepath"
	"strings"
)

func ProcessFiles(config *Config) ([]FileResult, error) {
	var results []FileResult

	for _, dir := range config.Dirs {
		if config.Debug {
			Debug("Processing directory: %s", dir)
		}
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Handle directories
			if info.IsDir() {
				if shouldIgnoreDir(path, config) {
					if config.Debug {
						Debug("Ignoring directory: %s", path)
					}
					return filepath.SkipDir
				}
				if !config.Recursive && path != dir {
					return filepath.SkipDir
				}
				return nil
			}

			// Handle files
			if shouldIgnoreFile(path, config) {
				if config.Debug {
					Debug("Ignoring file: %s", path)
				}
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			results = append(results, FileResult{
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

func shouldIgnoreDir(path string, config *Config) bool {
	for _, ignoreDir := range config.IgnoreDirs {
		if strings.Contains(path, ignoreDir) {
			return true
		}
	}
	return false
}

func shouldIgnoreFile(path string, config *Config) bool {
	filename := filepath.Base(path)
	ext := strings.TrimPrefix(filepath.Ext(path), ".")

	for _, ignoreFile := range config.IgnoreFiles {
		if filename == ignoreFile {
			return true
		}
	}

	if len(config.IncludeExts) > 0 {
		include := false
		for _, includeExt := range config.IncludeExts {
			if ext == includeExt {
				include = true
				break
			}
		}
		if !include {
			return true
		}
	}

	for _, ignoreExt := range config.IgnoreExts {
		if ext == ignoreExt {
			return true
		}
	}

	return false
}

type FileResult struct {
	Path    string
	Content string
}
