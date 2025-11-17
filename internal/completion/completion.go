// Package completion provides shell completion installation functionality.
// It supports automatic detection and installation of completions for bash, zsh, and fish shells.
package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Installer handles shell completion installation
type Installer struct {
	shellType string
}

// NewInstaller creates a new completion installer
func NewInstaller() (*Installer, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return nil, fmt.Errorf("could not detect shell from SHELL environment variable")
	}

	var shellType string
	switch {
	case strings.Contains(shell, "bash"):
		shellType = "bash"
	case strings.Contains(shell, "zsh"):
		shellType = "zsh"
	case strings.Contains(shell, "fish"):
		shellType = "fish"
	default:
		return nil, fmt.Errorf("unsupported shell: %s", shell)
	}

	return &Installer{shellType: shellType}, nil
}

// Install installs shell completion for the detected shell
func (i *Installer) Install(rootCmd *cobra.Command) error {
	switch i.shellType {
	case "bash":
		return installBashCompletion(rootCmd)
	case "zsh":
		return installZshCompletion(rootCmd)
	case "fish":
		return installFishCompletion(rootCmd)
	default:
		return fmt.Errorf("unsupported shell: %s", i.shellType)
	}
}

// InstallCompletion detects the current shell from the SHELL environment variable
// and installs the appropriate shell completions for the given Cobra command.
// Supported shells: bash, zsh, fish.
func InstallCompletion(rootCmd *cobra.Command) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		fmt.Println("Could not detect shell. Please set the SHELL environment variable.")
		return
	}
	var shellType string
	switch {
	case strings.Contains(shell, "bash"):
		shellType = "bash"
	case strings.Contains(shell, "zsh"):
		shellType = "zsh"
	case strings.Contains(shell, "fish"):
		shellType = "fish"
	default:
		fmt.Printf("Shell %s is not supported for automatic installation.\n", shell)
		return
	}

	switch shellType {
	case "bash":
		installBashCompletion(rootCmd)
	case "zsh":
		installZshCompletion(rootCmd)
	case "fish":
		installFishCompletion(rootCmd)
	}
}

// installBashCompletion generates and installs bash completion scripts.
// It attempts to install to /etc/bash_completion.d/ if writable, otherwise to the user's home directory.
func installBashCompletion(rootCmd *cobra.Command) error {
	etcPath := "/etc/bash_completion.d/codexgigantus"
	targetPath := ""
	if isWritable(filepath.Dir(etcPath)) {
		targetPath = etcPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error finding user home directory: %w", err)
		}
		targetPath = filepath.Join(homeDir, ".codexgigantus_completion")
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating bash completion file: %w", err)
	}
	defer f.Close()

	if err := rootCmd.GenBashCompletion(f); err != nil {
		return fmt.Errorf("error generating bash completion: %w", err)
	}

	// If installed in the home directory, append a source command to .bashrc if needed.
	if !strings.HasPrefix(targetPath, "/etc/") {
		bashrc := filepath.Join(os.Getenv("HOME"), ".bashrc")
		sourceLine := fmt.Sprintf("\n# CodexGigantus completion\nsource %s\n", targetPath)
		appendIfNotExists(bashrc, sourceLine)
	}
	fmt.Printf("Bash completions installed to %s. Restart your shell or run 'source %s' to activate.\n", targetPath, targetPath)
	return nil
}

// installZshCompletion generates and installs zsh completion scripts.
// It creates a completions directory in ~/.zsh/ and updates .zshrc if needed.
func installZshCompletion(rootCmd *cobra.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error finding user home directory: %w", err)
	}
	completionsDir := filepath.Join(homeDir, ".zsh", "completions")
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		return fmt.Errorf("error creating zsh completions directory: %w", err)
	}
	targetPath := filepath.Join(completionsDir, "_codexgigantus")
	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating zsh completion file: %w", err)
	}
	defer f.Close()

	if err := rootCmd.GenZshCompletion(f); err != nil {
		return fmt.Errorf("error generating zsh completion: %w", err)
	}

	// Ensure that .zshrc contains the necessary configuration.
	zshrc := filepath.Join(homeDir, ".zshrc")
	zshSetup := fmt.Sprintf("\n# CodexGigantus completion\nfpath=(%s $fpath)\nautoload -Uz compinit && compinit\n", completionsDir)
	appendIfNotExists(zshrc, zshSetup)
	fmt.Printf("Zsh completions installed to %s. Restart your shell to activate.\n", targetPath)
	return nil
}

// installFishCompletion generates and installs fish completion scripts.
// It creates the completion file in ~/.config/fish/completions/.
func installFishCompletion(rootCmd *cobra.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error finding user home directory: %w", err)
	}
	completionsDir := filepath.Join(homeDir, ".config", "fish", "completions")
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		return fmt.Errorf("error creating fish completions directory: %w", err)
	}
	targetPath := filepath.Join(completionsDir, "codexgigantus.fish")
	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating fish completion file: %w", err)
	}
	defer f.Close()

	if err := rootCmd.GenFishCompletion(f, true); err != nil {
		return fmt.Errorf("error generating fish completion: %w", err)
	}
	fmt.Printf("Fish completions installed to %s. Restart your shell to activate.\n", targetPath)
	return nil
}

// isWritable checks if a directory is writable by attempting to create a test file.
// It returns true if the directory is writable, false otherwise.
func isWritable(dir string) bool {
	testFile := filepath.Join(dir, ".writetest")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		return false
	}
	os.Remove(testFile)
	return true
}

// appendIfNotExists appends a line to a file if it isn't already present.
// If the file doesn't exist, it creates it with the given content.
func appendIfNotExists(filename, line string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// If the file doesn't exist, create it with the line.
		os.WriteFile(filename, []byte(line), 0644)
		return
	}
	if !strings.Contains(string(data), line) {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		f.WriteString(line)
	}
}
