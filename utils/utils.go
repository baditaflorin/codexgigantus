package utils

import (
	"fmt"
	"os"
)

func SaveOutput(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func DisplaySize(files []string) {
	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err == nil {
			totalSize += info.Size()
		}
	}

	switch {
	case totalSize < 1024:
		fmt.Printf("Output size: %d bytes\n", totalSize)
	case totalSize < 1024*1024:
		fmt.Printf("Output size: %.2f KB\n", float64(totalSize)/1024)
	default:
		fmt.Printf("Output size: %.2f MB\n", float64(totalSize)/(1024*1024))
	}
}
