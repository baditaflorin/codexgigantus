package utils

import (
	"fmt"
	"os"

	"github.com/baditaflorin/codexgigantus/config"
)

func SaveOutput(filename, data string) error {
	if err := os.WriteFile(filename, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to save output to %s: %w", filename, err)
	}
	return nil
}

func DisplaySize(files []string) error {
	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", file, err)
		}
		totalSize += info.Size()
	}

	switch {
	case totalSize < 1024:
		fmt.Printf("Output size: %d bytes\n", totalSize)
	case totalSize < 1024*1024:
		fmt.Printf("Output size: %.2f KB\n", float64(totalSize)/1024)
	default:
		fmt.Printf("Output size: %.2f MB\n", float64(totalSize)/(1024*1024))
	}
	return nil
}

func PrintDebugInfo(cfg *config.Config) {
	fmt.Println("Debug mode enabled v_0.1")
	fmt.Printf("Directories: %s\n", cfg.Dirs)
	fmt.Printf("Ignore files: %s\n", cfg.IgnoreFiles)
	fmt.Printf("Ignore directories: %s\n", cfg.IgnoreDirs)
	fmt.Printf("Ignore extensions: %s\n", cfg.IgnoreExts)
	fmt.Printf("Recursive: %v\n", cfg.Recursive)
}
