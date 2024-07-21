package main

import (
	"fmt"
	"strings"

	"github.com/baditaflorin/codexgigantus/config"
	"github.com/baditaflorin/codexgigantus/filehandling"
	"github.com/baditaflorin/codexgigantus/processing"
	"github.com/baditaflorin/codexgigantus/utils"
)

func main() {
	cfg := config.ParseConfig()

	if cfg.Debug {
		utils.PrintDebugInfo(cfg)
	}

	output := processFiles(cfg, filehandling.GatherIncludedFiles, processing.ProcessFiles, processing.DefaultProcessFunc)

	if cfg.Save {
		if err := utils.SaveOutput(cfg.OutputFile, output.String()); err != nil {
			fmt.Printf("Error saving output to file %q: %v\n", cfg.OutputFile, err)
		}
	}

	if cfg.ShowSize {
		if err := utils.DisplaySize(strings.Split(output.String(), "\n")); err != nil {
			fmt.Printf("Error displaying size: %v\n", err)
		}
	} else {
		fmt.Print(output.String())
	}
}

func processFiles(cfg *config.Config, gatherFunc func(*config.Config) ([]string, error), processFunc func([]string, *strings.Builder, func(string) ([]byte, error)), fileProcessFunc func(string) ([]byte, error)) strings.Builder {
	var output strings.Builder

	if cfg.ShowFuncs {
		filehandling.ProcessDirectories(cfg.Dirs, processing.ShowFunctions, cfg)
		return output
	}

	includedFiles, err := gatherFunc(cfg)
	if err != nil {
		fmt.Printf("Error gathering included files: %v\n", err)
		return output
	}

	if cfg.Debug {
		fmt.Println("Included files:")
		for _, file := range includedFiles {
			fmt.Println(file)
		}
	}

	processFunc(includedFiles, &output, fileProcessFunc)
	return output
}
