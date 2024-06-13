package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dirs          string
	ignoreFiles   string
	ignoreDirs    string
	ignoreExts    string
	recursive     bool
	debug         bool
	save          bool
	outputFile    string
	showSize      bool
	output        strings.Builder
	showFuncs     bool
	includedFiles []string
)

func init() {
	flag.StringVar(&dirs, "dir", ".", "Comma-separated list of directories to search")
	flag.StringVar(&ignoreFiles, "ignore-file", "", "Comma-separated list of files to ignore")
	flag.StringVar(&ignoreDirs, "ignore-dir", "", "Comma-separated list of directories to ignore")
	flag.StringVar(&ignoreExts, "ignore-ext", "", "Comma-separated list of file extensions to ignore")
	flag.BoolVar(&recursive, "recursive", true, "Recursively search directories")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.BoolVar(&save, "save", false, "Save the output to a file")
	flag.StringVar(&outputFile, "output-file", "output.txt", "Specify the output file name")
	flag.BoolVar(&showSize, "show-size", false, "Show the size of the result in bytes")
	flag.BoolVar(&showFuncs, "show-funcs", false, "Show only functions and their parameters")
}

func main() {
	flag.Parse()

	if debug {
		printDebugInfo()
	}

	if showFuncs {
		for _, dir := range strings.Split(dirs, ",") {
			err := filepath.Walk(dir, showFunctions)
			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", dir, err)
			}
		}
		return
	}

	includedFiles = gatherIncludedFiles(strings.Split(dirs, ","), ignoreFiles, ignoreDirs, ignoreExts)

	if debug {
		fmt.Println("Included files:")
		for _, file := range includedFiles {
			fmt.Println(file)
		}
	}

	if save {
		for _, file := range includedFiles {
			addFileContentToOutput(file)
		}
		err := saveOutput(outputFile, output.String())
		if err != nil {
			fmt.Printf("Error saving output to file %q: %v\n", outputFile, err)
		}
	}

	if showSize {
		displaySize(includedFiles)
	} else {
		fmt.Print(output.String())
	}
}

func printDebugInfo() {
	fmt.Println("Debug mode enabled v_0.1")
	fmt.Printf("Directories: %s\n", dirs)
	fmt.Printf("Ignore files: %s\n", ignoreFiles)
	fmt.Printf("Ignore directories: %s\n", ignoreDirs)
	fmt.Printf("Ignore extensions: %s\n", ignoreExts)
	fmt.Printf("Recursive: %v\n", recursive)
}

func shouldSkipDir(path string, ignoreDirs string) bool {
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
	if contains(ignoreFilesList, info.Name()) {
		return false
	}
	if containsExt(ignoreExtsList, filepath.Ext(info.Name())) {
		return false
	}
	return true
}

func gatherIncludedFiles(dirs []string, ignoreFiles, ignoreDirs, ignoreExts string) []string {
	var files []string
	for _, dir := range dirs {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return handleError(err, path)
			}

			if info.IsDir() {
				if shouldSkipDir(path, ignoreDirs) {
					if debug {
						fmt.Printf("Skipping directory: %s\n", path)
					}
					return filepath.SkipDir
				}
				return nil
			}

			if shouldIncludeFile(path, info, ignoreFiles, ignoreExts) {
				if debug {
					fmt.Printf("Including file: %s\n", path)
				}
				files = append(files, path)
			} else {
				if debug {
					fmt.Printf("Excluding file: %s\n", path)
				}
			}

			return nil
		})
	}
	return files
}

func showFunctions(path string, info os.FileInfo, err error) error {
	if err != nil {
		return handleError(err, path)
	}

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}

		functions := extractFunctions(fileContent)
		if len(functions) > 0 {
			fmt.Printf("Functions in file %s:\n", path)
			for _, f := range functions {
				fmt.Println(f)
			}
		}
	}

	return nil
}

func extractFunctions(content []byte) []string {
	var functions []string
	re := regexp.MustCompile(`func\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)`)
	matches := re.FindAllSubmatch(content, -1)
	for _, match := range matches {
		functionName := string(match[1])
		parameters := string(match[2])
		functions = append(functions, fmt.Sprintf("%s(%s)", functionName, parameters))
	}
	return functions
}

func createWalkFunc(ignoreFiles, ignoreDirs, ignoreExts string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return handleError(err, path)
		}

		if info.IsDir() {
			return handleDirectory(path, info, ignoreDirs)
		}

		return handleFile(path, info, ignoreFiles, ignoreExts)
	}
}

func handleError(err error, path string) error {
	fmt.Printf("Error accessing path %q: %v\n", path, err)
	return err
}

func handleDirectory(path string, info os.FileInfo, ignoreDirs string) error {
	ignoreDirsList := strings.Split(ignoreDirs, ",")
	for _, dir := range ignoreDirsList {
		if strings.Contains(filepath.ToSlash(path), filepath.ToSlash(dir)) {
			if debug {
				fmt.Printf("Skipping directory: %s\n", path)
			}
			return filepath.SkipDir
		}
	}
	return nil
}

func handleFile(path string, info os.FileInfo, ignoreFiles, ignoreExts string) error {
	ignoreFilesList := strings.Split(ignoreFiles, ",")
	ignoreExtsList := strings.Split(ignoreExts, ",")
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
	addFileContentToOutput(path)
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

func addFileContentToOutput(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}
	output.WriteString(fmt.Sprintf("________\nPath: %s\n", path))
	output.WriteString(fmt.Sprintf("Content:\n%s\n", content))
}

func saveOutput(filename, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func displaySize(files []string) {
	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err == nil {
			totalSize += info.Size()
		}
	}

	if totalSize < 1024 {
		fmt.Printf("Output size: %d bytes\n", totalSize)
	} else if totalSize < 1024*1024 {
		fmt.Printf("Output size: %.2f KB\n", float64(totalSize)/1024)
	} else {
		fmt.Printf("Output size: %.2f MB\n", float64(totalSize)/(1024*1024))
	}
}
