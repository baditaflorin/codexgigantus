package filehandling

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/baditaflorin/codexgigantus/config"
)

func ValidateDirectory(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

func GatherIncludedFiles(dirs, ignoreFiles, ignoreDirs, ignoreExts, ignoreSuffixes string, debug bool) ([]string, error) {
	var files []string
	dirList := strings.Split(dirs, ",")
	for _, dir := range dirList {
		if err := filepath.Walk(dir, createWalkFunc(ignoreFiles, ignoreDirs, ignoreExts, ignoreSuffixes, &files, debug)); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func ProcessDirectories(dirs string, processFunc filepath.WalkFunc, cfg *config.Config) error {
	dirList := strings.Split(dirs, ",")
	for _, dir := range dirList {
		if err := filepath.Walk(dir, processFunc); err != nil {
			return err
		}
	}
	return nil
}

func createWalkFunc(ignoreFiles, ignoreDirs, ignoreExts, ignoreSuffixes string, files *[]string, debug bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Always include the file if no ignore flags are set
		if ignoreFiles == "" && ignoreDirs == "" && ignoreExts == "" && ignoreSuffixes == "" {
			if !info.IsDir() {
				if debug {
					println("Including file:", path)
				}
				*files = append(*files, path)
			}
			return nil
		}

		if info.IsDir() && shouldSkipDir(path, ignoreDirs) {
			if debug {
				println("Skipping directory:", path)
			}
			return filepath.SkipDir
		}

		if !info.IsDir() && shouldIncludeFile(path, info, ignoreFiles, ignoreExts, ignoreSuffixes) {
			if debug {
				println("Including file:", path)
			}
			*files = append(*files, path)
		} else if debug {
			println("Excluding file:", path)
		}
		return nil
	}
}

func shouldSkipDir(path, ignoreDirs string) bool {
	ignoreDirsList := strings.Split(ignoreDirs, ",")
	for _, dir := range ignoreDirsList {
		if filepath.Base(path) == dir {
			return true
		}
	}
	return false
}

func shouldIncludeFile(path string, info os.FileInfo, ignoreFiles, ignoreExts, ignoreSuffixes string) bool {
	// If no ignore flags are set, include all files
	if ignoreFiles == "" && ignoreExts == "" && ignoreSuffixes == "" {
		return true
	}

	ignoreFilesList := strings.Split(ignoreFiles, ",")
	ignoreExtsList := strings.Split(ignoreExts, ",")
	ignoreSuffixesList := strings.Split(ignoreSuffixes, ",")

	if contains(ignoreFilesList, info.Name()) {
		return false
	}
	if containsExt(ignoreExtsList, filepath.Ext(info.Name())) {
		return false
	}
	if containsSuffix(ignoreSuffixesList, info.Name()) {
		return false
	}
	return true
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func containsExt(list []string, ext string) bool {
	ext = strings.TrimPrefix(ext, ".")
	for _, e := range list {
		if e == ext {
			return true
		}
	}
	return false
}

func containsSuffix(list []string, name string) bool {
	for _, suffix := range list {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
