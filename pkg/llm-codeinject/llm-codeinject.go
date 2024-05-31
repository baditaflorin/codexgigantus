package llm_codeinject

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ValidateDirectory checks if a directory exists
func ValidateDirectory(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// FilterFiles filters files by directory and extensions
func FilterFiles(root string, ignoreDirs []string, ignoreExts []string, includeExts []string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a directory we need to ignore
		for _, dir := range ignoreDirs {
			if info.IsDir() && strings.Contains(path, dir) {
				return filepath.SkipDir
			}
		}

		// If it's a file, apply extension filters
		if !info.IsDir() {
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if len(ignoreExts) > 0 && contains(ignoreExts, ext) {
				return nil
			}
			if len(includeExts) == 0 || contains(includeExts, ext) {
				files = append(files, path)
			}
		}

		return nil
	})
	return files, err
}

// ProcessFiles processes each file found
func ProcessFiles(files []string) {
	for _, file := range files {
		if file != "" {
			println("Processing file:", file)
			content, err := ioutil.ReadFile(file)
			if err != nil {
				println("Error reading file:", err)
				continue
			}
			println("Contents of file:", file)
			println(string(content))
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
