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

var rootCmd = &cobra.Command{
	Use:   "codexgigantus",
	Short: "Process files in a directory based on given criteria",
	Long: `CodexGigantus is a command-line tool that processes files from specified directories.
It supports ignoring directories, filtering by file extensions, and more.
Now using Cobra for robust CLI parsing and automatic shell completions installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags using the config package
		cfg := config.ParseFlags()

		fmt.Println("Running CodexGigantus with the following configuration:")
		fmt.Printf("  Directory: %v\n", cfg.Dirs)
		fmt.Printf("  Ignore Files: %v\n", cfg.IgnoreFiles)
		fmt.Printf("  Ignore Dirs: %v\n", cfg.IgnoreDirs)
		fmt.Printf("  Ignore Ext: %v\n", cfg.IgnoreExts)
		fmt.Printf("  Include Ext: %v\n", cfg.IncludeExts)
		fmt.Printf("  Recursive: %v\n", cfg.Recursive)
		fmt.Printf("  Debug: %v\n", cfg.Debug)
		fmt.Printf("  Save: %v\n", cfg.Save)
		fmt.Printf("  Output File: %s\n", cfg.OutputFile)
		fmt.Printf("  Show Size: %v\n", cfg.ShowSize)
		fmt.Printf("  Show Funcs: %v\n", cfg.ShowFuncs)

		results, err := processor.ProcessFiles(cfg)
		if err != nil {
			fmt.Println("Error processing files:", err)
			os.Exit(1)
		}

		output := utils.GenerateOutput(results, cfg)
		fmt.Println(output)

		if cfg.Save {
			err = utils.SaveOutput(output, cfg.OutputFile)
			if err != nil {
				fmt.Println("Error saving output:", err)
			} else {
				fmt.Printf("Output saved to %s\n", cfg.OutputFile)
			}
		}
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `Generate the auto-completion script for your shell.

If no shell is specified, the available options are:
  bash
  zsh
  fish
  powershell

Usage:
  ./codexgigantus completion [shell]
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Available shells:")
			fmt.Println("  bash")
			fmt.Println("  zsh")
			fmt.Println("  fish")
			fmt.Println("  powershell")
			fmt.Println("\nUsage: ./codexgigantus completion [shell]")
			return
		}
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			fmt.Printf("Shell %s is not supported.\n", args[0])
		}
	},
}

var installCompletionCmd = &cobra.Command{
	Use:   "install-completion",
	Short: "Automatically install shell completions",
	Long:  `Automatically install shell completions for your shell.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Delegate the installation to the completion package.
		completion.InstallCompletion(rootCmd)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(installCompletionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
