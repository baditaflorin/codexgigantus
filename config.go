// config.go
package main

import (
	"flag"
	"strings"
)

type Config struct {
	Dirs        []string
	IgnoreFiles []string
	IgnoreDirs  []string
	IgnoreExts  []string
	IncludeExts []string
	Recursive   bool
	Debug       bool
	Save        bool
	OutputFile  string
	ShowSize    bool
	ShowFuncs   bool
}

func ParseFlags() *Config {
	config := &Config{}

	dirFlag := flag.String("dir", ".", "Comma-separated list of directories to search (default: current directory)")
	ignoreFileFlag := flag.String("ignore-file", "", "Comma-separated list of files to ignore")
	ignoreDirFlag := flag.String("ignore-dir", "", "Comma-separated list of directories to ignore")
	ignoreExtFlag := flag.String("ignore-ext", "", "Comma-separated list of file extensions to ignore")
	includeExtFlag := flag.String("include-ext", "", "Comma-separated list of file extensions to include")
	recursiveFlag := flag.Bool("recursive", true, "Recursively search directories (default: true)")
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	saveFlag := flag.Bool("save", false, "Save the output to a file")
	outputFileFlag := flag.String("output-file", "output.txt", "Specify the output file name (default: output.txt)")
	showSizeFlag := flag.Bool("show-size", false, "Show the size of the result in bytes")
	showFuncsFlag := flag.Bool("show-funcs", false, "Show only functions and their parameters")

	flag.Parse()

	config.Dirs = parseCommaSeparated(*dirFlag)
	config.IgnoreFiles = parseCommaSeparated(*ignoreFileFlag)
	config.IgnoreDirs = parseCommaSeparated(*ignoreDirFlag)
	config.IgnoreExts = parseCommaSeparated(*ignoreExtFlag)
	config.IncludeExts = parseCommaSeparated(*includeExtFlag)
	config.Recursive = *recursiveFlag
	config.Debug = *debugFlag
	config.Save = *saveFlag
	config.OutputFile = *outputFileFlag
	config.ShowSize = *showSizeFlag
	config.ShowFuncs = *showFuncsFlag

	return config
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
