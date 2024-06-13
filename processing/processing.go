package processing

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ProcessFiles(files []string, output *strings.Builder, processFunc func(string) ([]byte, error)) {
	for _, file := range files {
		content, err := processFunc(file)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", file, err)
			continue
		}
		output.WriteString(fmt.Sprintf("________\nPath: %s\nContent:\n%s\n", file, content))
	}
}

func DefaultProcessFunc(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func ShowFunctions(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
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
