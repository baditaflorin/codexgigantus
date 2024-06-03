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
	dirs        string
	ignoreFiles string
	ignoreDirs  string
	ignoreExts  string
	recursive   bool
	debug       bool
	save        bool
	outputFile  string
	showSize    bool
	output      strings.Builder
	showFuncs   bool
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

	for _, dir := range strings.Split(dirs, ",") {
		err := filepath.Walk(dir, createWalkFunc(ignoreFiles, ignoreDirs, ignoreExts))
		if err != nil {
			fmt.Printf("Error walking the path %q: %v\n", dir, err)
		}
	}

	if save {
		err := saveOutput(outputFile, output.String())
		if err != nil {
			fmt.Printf("Error saving output to file %q: %v\n", outputFile, err)
		}
	}

	if showSize {
		displaySize(output.String())
	} else {
		fmt.Print(output.String())
	}
}

func printDebugInfo() {
	fmt.Println("Debug mode enabled")
	fmt.Printf("Directories: %s\n", dirs)
	fmt.Printf("Ignore files: %s\n", ignoreFiles)
	fmt.Printf("Ignore directories: %s\n", ignoreDirs)
	fmt.Printf("Ignore extensions: %s\n", ignoreExts)
	fmt.Printf("Recursive: %v\n", recursive)
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
	if contains(ignoreDirsList, info.Name()) {
		if debug {
			fmt.Printf("Skipping directory: %s\n", path)
		}
		return filepath.SkipDir
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

func displaySize(data string) {
	size := len(data)
	if size < 1024 {
		fmt.Printf("Output size: %d bytes\n", size)
	} else if size < 1024*1024 {
		fmt.Printf("Output size: %.2f KB\n", float64(size)/1024)
	} else {
		fmt.Printf("Output size: %.2f MB\n", float64(size)/(1024*1024))
	}
}
