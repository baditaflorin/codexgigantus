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

func ProcessDirectories(dirs string, processFunc filepath.WalkFunc, cfg *config.Config) error {
	dirList := strings.Split(dirs, ",")
	for _, dir := range dirList {
		if err := filepath.Walk(dir, processFunc); err != nil {
			return err
		}
	}
	return nil
}


func GatherIncludedFiles(cfg *config.Config) ([]string, error) {
	var files []string
	dirList := strings.Split(cfg.Dirs, ",")
	for _, dir := range dirList {
		if err := filepath.Walk(dir, createWalkFunc(cfg, &files)); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func createWalkFunc(cfg *config.Config, files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && shouldSkipDir(path, cfg.IgnoreDirs) {
			if cfg.Debug {
				println("Skipping directory:", path)
			}
			return filepath.SkipDir
		}
		if !info.IsDir() && shouldIncludeFile(path, info, cfg) {
			if cfg.Debug {
				println("Including file:", path)
			}
			*files = append(*files, path)
		} else if cfg.Debug {
			println("Excluding file:", path)
		}
		return nil
	}
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

func shouldIncludeFile(path string, info os.FileInfo, cfg *config.Config) bool {
	return !contains(strings.Split(cfg.IgnoreFiles, ","), info.Name()) &&
		!containsExt(strings.Split(cfg.IgnoreExts, ","), filepath.Ext(info.Name())) &&
		!containsSuffix(strings.Split(cfg.IgnoreSuffix, ","), info.Name())
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
