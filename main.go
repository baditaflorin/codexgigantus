// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	config := ParseFlags()

	if config.Debug {
		fmt.Println("Debug mode enabled")
		fmt.Printf("Configuration: %+v\n", config)
	}

	results, err := ProcessFiles(config)
	if err != nil {
		fmt.Println("Error processing files:", err)
		os.Exit(1)
	}

	output := GenerateOutput(results, config)

	if config.Save {
		err = SaveOutput(output, config.OutputFile)
		if err != nil {
			fmt.Println("Error saving output:", err)
			os.Exit(1)
		}
		fmt.Println("Output saved to", config.OutputFile)
	} else {
		fmt.Println(output)
	}

	if config.ShowSize {
		fmt.Printf("Total size: %d bytes\n", len(output))
	}
}
