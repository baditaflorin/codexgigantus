package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	dir         string
	ignoreFiles string
	ignoreDirs  string
	ignoreExts  string
	recursive   bool
	debug       bool
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory to search")
	flag.StringVar(&ignoreFiles, "ignore-file", "", "Comma-separated list of files to ignore")
	flag.StringVar(&ignoreDirs, "ignore-dir", "", "Comma-separated list of directories to ignore")
	flag.StringVar(&ignoreExts, "ignore-ext", "", "Comma-separated list of file extensions to ignore")
	flag.BoolVar(&recursive, "recursive", true, "Recursively search directories")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
}

func main() {
	flag.Parse()

	if debug {
		printDebugInfo()
	}

	ignoreFilesList := strings.Split(ignoreFiles, ",")
	ignoreDirsList := strings.Split(ignoreDirs, ",")
	ignoreExtsList := strings.Split(ignoreExts, ",")

	err := filepath.Walk(dir, createWalkFunc(ignoreFilesList, ignoreDirsList, ignoreExtsList))
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", dir, err)
	}
}

func printDebugInfo() {
	fmt.Println("Debug mode enabled")
	fmt.Printf("Directory: %s\n", dir)
	fmt.Printf("Ignore files: %s\n", ignoreFiles)
	fmt.Printf("Ignore directories: %s\n", ignoreDirs)
	fmt.Printf("Ignore extensions: %s\n", ignoreExts)
	fmt.Printf("Recursive: %v\n", recursive)
}

func createWalkFunc(ignoreFilesList, ignoreDirsList, ignoreExtsList []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return handleError(err, path)
		}

		if info.IsDir() {
			return handleDirectory(path, info, ignoreDirsList)
		}

		return handleFile(path, info, ignoreFilesList, ignoreExtsList)
	}
}

func handleError(err error, path string) error {
	fmt.Printf("Error accessing path %q: %v\n", path, err)
	return err
}

func handleDirectory(path string, info os.FileInfo, ignoreDirsList []string) error {
	if contains(ignoreDirsList, info.Name()) {
		if debug {
			fmt.Printf("Skipping directory: %s\n", path)
		}
		return filepath.SkipDir
	}
	if !recursive && path != dir {
		return filepath.SkipDir
	}
	return nil
}

func handleFile(path string, info os.FileInfo, ignoreFilesList, ignoreExtsList []string) error {
	if contains(ignoreFilesList, info.Name()) {
		if debug {
			fmt.Printf("Skipping file: %s\n", path)
		}
		return nil
	}
	if containsExt(ignoreExtsList, filepath.Ext(info.Name())) {
		if debug {
			fmt.Printf("Skipping file with extension: %s\n", path)
		}
		return nil
	}
	printFileContent(path)
	return nil
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

func printFileContent(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}
	fmt.Printf("________\nPath: %s\n", path)
	fmt.Printf("Content:\n%s\n", content)
}
