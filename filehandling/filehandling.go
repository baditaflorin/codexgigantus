package filehandling

import (
	"os"
	"path/filepath"
	"strings"
)

func ValidateDirectory(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

func GatherIncludedFiles(dirs, ignoreFiles, ignoreDirs, ignoreExts string, debug bool) []string {
	var files []string
	for _, dir := range strings.Split(dirs, ",") {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && shouldSkipDir(path, ignoreDirs) {
				if debug {
					println("Skipping directory:", path)
				}
				return filepath.SkipDir
			}
			if !info.IsDir() && shouldIncludeFile(path, info, ignoreFiles, ignoreExts) {
				if debug {
					println("Including file:", path)
				}
				files = append(files, path)
			} else if debug {
				println("Excluding file:", path)
			}
			return nil
		})
	}
	return files
}

func shouldSkipDir(path, ignoreDirs string) bool {
	ignoreDirsList := strings.Split(ignoreDirs, ",")
	for _, dir := range ignoreDirsList {
		if strings.Contains(filepath.ToSlash(path), filepath.ToSlash(dir)) {
			return true
		}
	}
	return false
}

func shouldIncludeFile(path string, info os.FileInfo, ignoreFiles, ignoreExts string) bool {
	ignoreFilesList := strings.Split(ignoreFiles, ",")
	ignoreExtsList := strings.Split(ignoreExts, ",")
	return !contains(ignoreFilesList, info.Name()) && !containsExt(ignoreExtsList, filepath.Ext(info.Name()))
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
