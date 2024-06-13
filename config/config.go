package config

import "flag"

type Config struct {
	Dirs        string
	IgnoreFiles string
	IgnoreDirs  string
	IgnoreExts  string
	Recursive   bool
	Debug       bool
	Save        bool
	OutputFile  string
	ShowSize    bool
	ShowFuncs   bool
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.Dirs, "dir", ".", "Comma-separated list of directories to search")
	flag.StringVar(&c.IgnoreFiles, "ignore-file", "", "Comma-separated list of files to ignore")
	flag.StringVar(&c.IgnoreDirs, "ignore-dir", "", "Comma-separated list of directories to ignore")
	flag.StringVar(&c.IgnoreExts, "ignore-ext", "", "Comma-separated list of file extensions to ignore")
	flag.BoolVar(&c.Recursive, "recursive", true, "Recursively search directories")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debug output")
	flag.BoolVar(&c.Save, "save", false, "Save the output to a file")
	flag.StringVar(&c.OutputFile, "output-file", "output.txt", "Specify the output file name")
	flag.BoolVar(&c.ShowSize, "show-size", false, "Show the size of the result in bytes")
	flag.BoolVar(&c.ShowFuncs, "show-funcs", false, "Show only functions and their parameters")
	flag.Parse()
}
