package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/katungi/edon/internals/runtime"
)

var (
	// CLI flags
	evalScript  = flag.String("eval", "", "Evaluate a JavaScript expression")
	showVersion = flag.Bool("version", false, "Show version information")
	showHelp    = flag.Bool("help", false, "Show help information")
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "install" {
		installCmd.Parse(os.Args[2:])
		if err := handleInstall(); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		return
	}

	flag.Usage = printHelp
	flag.Parse()

	if err := run(); err != nil {
		if err != runtime.ErrExit && err != runtime.ErrInterrupt {
			color.Red("Error: %v", err)
		} 

		os.Exit(1)
	}
}

// Version information
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func run() error {
	// Handle version flag
	if *showVersion {
		printVersion()
		return nil
	}

	// Handle help flag
	if *showHelp {
		printHelp()
		return nil
	}

	// Create new runtime instance
	rt, err := runtime.New()
	if err != nil {
		return fmt.Errorf("failed to initialize runtime: %w", err)
	}
	defer rt.Close()

	// Handle -eval flag
	if *evalScript != "" {
		return rt.Eval(*evalScript)
	}

	// Handle file argument
	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filename)
		}
		return rt.ExecuteFile(filename)
	}

	// No file or eval provided, start REPL
	return rt.StartREPL()
}

func printVersion() {
	info, ok := debug.ReadBuildInfo()

	fmt.Printf("Halo JavaScript Runtime %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Build date: %s\n", date)

	if ok {
		fmt.Println("\nDependencies:")
		for _, dep := range info.Deps {
			fmt.Printf("- %s %s\n", dep.Path, dep.Version)
		}
	}
}

func printHelp() {
	exe := filepath.Base(os.Args[0])
	help := `
Halo JavaScript Runtime

Usage:
  %s [options] [file]

Options:
  -eval string    Execute a JavaScript expression
  -version        Show version information
  -help           Show this help message

Examples:
  # Start REPL
  %s

  # Execute a file
  %s script.js

  # Evaluate expression
  %s -eval "console.log('Hello, World!')"
`
	fmt.Printf(help, exe, exe, exe, exe)
}
