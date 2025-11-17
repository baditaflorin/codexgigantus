package processor

// Config is a re-export of the legacy config.Config for backward compatibility
// while allowing processor package to be used independently
type Config struct {
	// Dirs is a list of directories to search for files
	Dirs []string
	// IgnoreFiles is a list of specific file names to ignore
	IgnoreFiles []string
	// IgnoreDirs is a list of directory names to ignore during traversal
	IgnoreDirs []string
	// IgnoreExts is a list of file extensions to ignore (without the dot)
	IgnoreExts []string
	// IncludeExts is a list of file extensions to include (without the dot).
	// If specified, only files with these extensions will be processed
	IncludeExts []string
	// Recursive determines whether to search directories recursively
	Recursive bool
	// Debug enables debug output when true
	Debug bool
}
