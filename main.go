// main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/baditaflorin/codexgigantus/internal/completion"
)

var (
	// Global flags.
	dir         string
	ignoreFile  string
	ignoreDir   string
	ignoreExt   string
	includeExt  string
	recursive   bool
	debug       bool
	save        bool
	outputFile  string
	showSize    bool
	showFuncs   bool
)

var rootCmd = &cobra.Command{
	Use:   "codexgigantus",
	Short: "Process files in a directory based on given criteria",
	Long: `CodexGigantus is a command-line tool that processes files from specified directories.
It supports ignoring directories, filtering by file extensions, and more.
Now using Cobra for robust CLI parsing and automatic shell completions installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running CodexGigantus with the following configuration:")
		fmt.Printf("  Directory: %s\n", dir)
		fmt.Printf("  Ignore Files: %s\n", ignoreFile)
		fmt.Printf("  Ignore Dirs: %s\n", ignoreDir)
		fmt.Printf("  Ignore Ext: %s\n", ignoreExt)
		fmt.Printf("  Include Ext: %s\n", includeExt)
		fmt.Printf("  Recursive: %v\n", recursive)
		fmt.Printf("  Debug: %v\n", debug)
		fmt.Printf("  Save: %v\n", save)
		fmt.Printf("  Output File: %s\n", outputFile)
		fmt.Printf("  Show Size: %v\n", showSize)
		fmt.Printf("  Show Funcs: %v\n", showFuncs)
		// Insert processing logic here.
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `Generate the auto-completion script for your shell.

Examples:

Bash:
  $ codexgigantus completion bash > /etc/bash_completion.d/codexgigantus

Zsh:
  $ codexgigantus completion zsh > "${fpath[1]}/_codexgigantus"

Fish:
  $ codexgigantus completion fish > ~/.config/fish/completions/codexgigantus.fish

PowerShell:
  PS> codexgigantus completion powershell | Out-String | Invoke-Expression
`,
	Args: cobra.ExactValidArgs(1),
	ValidArgs: []string{
		"bash",
		"zsh",
		"fish",
		"powershell",
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
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
	cobra.OnInitialize(initConfig)

	// Global flags.
	rootCmd.PersistentFlags().StringVar(&dir, "dir", ".", "Comma-separated list of directories to search (default: current directory)")
	rootCmd.PersistentFlags().StringVar(&ignoreFile, "ignore-file", "", "Comma-separated list of files to ignore")
	rootCmd.PersistentFlags().StringVar(&ignoreDir, "ignore-dir", "", "Comma-separated list of directories to ignore")
	rootCmd.PersistentFlags().StringVar(&ignoreExt, "ignore-ext", "", "Comma-separated list of file extensions to ignore")
	rootCmd.PersistentFlags().StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include")
	rootCmd.PersistentFlags().BoolVar(&recursive, "recursive", true, "Recursively search directories (default: true)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&save, "save", false, "Save the output to a file")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output-file", "output.txt", "Specify the output file name (default: output.txt)")
	rootCmd.PersistentFlags().BoolVar(&showSize, "show-size", false, "Show the size of the result in bytes")
	rootCmd.PersistentFlags().BoolVar(&showFuncs, "show-funcs", false, "Show only functions and their parameters")

	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(installCompletionCmd)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
