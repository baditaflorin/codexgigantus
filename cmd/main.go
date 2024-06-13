package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"llm-codeinject/config"
	"llm-codeinject/filehandling"
	"llm-codeinject/processing"
	"llm-codeinject/utils"
)

var output strings.Builder

func main() {
	cfg := config.NewConfig()
	cfg.ParseFlags()

	if cfg.Debug {
		printDebugInfo(cfg)
	}

	if cfg.ShowFuncs {
		for _, dir := range strings.Split(cfg.Dirs, ",") {
			err := filepath.Walk(dir, processing.ShowFunctions)
			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", dir, err)
			}
		}
		return
	}

	includedFiles := filehandling.GatherIncludedFiles(cfg.Dirs, cfg.IgnoreFiles, cfg.IgnoreDirs, cfg.IgnoreExts, cfg.Debug)

	if cfg.Debug {
		fmt.Println("Included files:")
		for _, file := range includedFiles {
			fmt.Println(file)
		}
	}

	if cfg.Save {
		processing.ProcessFiles(includedFiles, &output)
		err := utils.SaveOutput(cfg.OutputFile, output.String())
		if err != nil {
			fmt.Printf("Error saving output to file %q: %v\n", cfg.OutputFile, err)
		}
	}

	if cfg.ShowSize {
		utils.DisplaySize(includedFiles)
	} else {
		fmt.Print(output.String())
	}
}

func printDebugInfo(cfg *config.Config) {
	fmt.Println("Debug mode enabled v_0.1")
	fmt.Printf("Directories: %s\n", cfg.Dirs)
	fmt.Printf("Ignore files: %s\n", cfg.IgnoreFiles)
	fmt.Printf("Ignore directories: %s\n", cfg.IgnoreDirs)
	fmt.Printf("Ignore extensions: %s\n", cfg.IgnoreExts)
	fmt.Printf("Recursive: %v\n", cfg.Recursive)
}
