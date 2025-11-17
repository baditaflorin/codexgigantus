package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/baditaflorin/codexgigantus/internal/completion"
	"github.com/baditaflorin/codexgigantus/pkg/config"
	"github.com/baditaflorin/codexgigantus/pkg/processor"
	"github.com/baditaflorin/codexgigantus/pkg/utils"
)

var (
	dirFlag        string
	ignoreFileFlag string
	ignoreDirFlag  string
	ignoreExtFlag  string
	includeExtFlag string
	recursiveFlag  bool
	debugFlag      bool
	saveFlag       bool
	outputFileFlag string
	showSizeFlag   bool
	showFuncsFlag  bool
)

var rootCmd = &cobra.Command{
	Use:   "codexgigantus",
	Short: "Process files in a directory based on given criteria",
	Long: `CodexGigantus is a command-line tool that processes files from specified directories.
It supports ignoring directories, filtering by file extensions, and more.
Now using Cobra for robust CLI parsing and automatic shell completions installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build config from flags
		cfg := &processor.Config{
			Dirs:        config.ParseCommaSeparated(dirFlag),
			IgnoreFiles: config.ParseCommaSeparated(ignoreFileFlag),
			IgnoreDirs:  config.ParseCommaSeparated(ignoreDirFlag),
			IgnoreExts:  config.ParseCommaSeparated(ignoreExtFlag),
			IncludeExts: config.ParseCommaSeparated(includeExtFlag),
			Recursive:   recursiveFlag,
			Debug:       debugFlag,
		}

		fmt.Println("Running CodexGigantus with the following configuration:")
		fmt.Printf("  Directory: %v\n", cfg.Dirs)
		fmt.Printf("  Ignore Files: %v\n", cfg.IgnoreFiles)
		fmt.Printf("  Ignore Dirs: %v\n", cfg.IgnoreDirs)
		fmt.Printf("  Ignore Ext: %v\n", cfg.IgnoreExts)
		fmt.Printf("  Include Ext: %v\n", cfg.IncludeExts)
		fmt.Printf("  Recursive: %v\n", cfg.Recursive)
		fmt.Printf("  Debug: %v\n", cfg.Debug)
		fmt.Printf("  Save: %v\n", saveFlag)
		fmt.Printf("  Output File: %s\n", outputFileFlag)
		fmt.Printf("  Show Size: %v\n", showSizeFlag)
		fmt.Printf("  Show Funcs: %v\n", showFuncsFlag)

		results, err := processor.ProcessFiles(cfg)
		if err != nil {
			fmt.Println("Error processing files:", err)
			os.Exit(1)
		}

		output := utils.GenerateOutput(results, showFuncsFlag)

		if showSizeFlag {
			fmt.Printf("\nTotal output size: %d bytes\n", len(output))
		}

		fmt.Println(output)

		if saveFlag {
			err = utils.SaveOutput(output, outputFileFlag)
			if err != nil {
				fmt.Println("Error saving output:", err)
			} else {
				fmt.Printf("Output saved to %s\n", outputFileFlag)
			}
		}
	},
}

var installCompletionCmd = &cobra.Command{
	Use:   "install-completion",
	Short: "Install shell completion",
	Long:  `Automatically install shell completions for your current shell (bash, zsh, fish, or PowerShell)`,
	Run: func(cmd *cobra.Command, args []string) {
		installer, err := completion.NewInstaller()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := installer.Install(rootCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Shell completion installed successfully!")
		fmt.Println("Please restart your shell or source your shell configuration file.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dirFlag, "dir", ".", "Comma-separated list of directories to search")
	rootCmd.PersistentFlags().StringVar(&ignoreFileFlag, "ignore-file", "", "Comma-separated list of file names to ignore")
	rootCmd.PersistentFlags().StringVar(&ignoreDirFlag, "ignore-dir", "", "Comma-separated list of directory names to ignore")
	rootCmd.PersistentFlags().StringVar(&ignoreExtFlag, "ignore-ext", "", "Comma-separated list of file extensions to ignore (without dot)")
	rootCmd.PersistentFlags().StringVar(&includeExtFlag, "include-ext", "", "Comma-separated list of file extensions to include (without dot)")
	rootCmd.PersistentFlags().BoolVar(&recursiveFlag, "recursive", true, "Search directories recursively")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&saveFlag, "save", false, "Save output to file")
	rootCmd.PersistentFlags().StringVar(&outputFileFlag, "output", "output.txt", "Output filename")
	rootCmd.PersistentFlags().BoolVar(&showSizeFlag, "show-size", false, "Show the size of the output")
	rootCmd.PersistentFlags().BoolVar(&showFuncsFlag, "show-funcs", false, "Show only function signatures from Go files")

	rootCmd.AddCommand(installCompletionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
