package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/katungi/edon/internal/runtime"
)

func main() {
	// Parse command line flags
	evalScript := flag.String("eval", "", "Script to evaluate")
	flag.Parse()

	// Initialize runtime
	rt, err := runtime.New()
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	// If script provided, run it
	if *evalScript != "" {
		if err := rt.Eval(*evalScript); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// If file provided as argument, execute it
	if len(flag.Args()) > 0 {
		if err := rt.ExecuteFile(flag.Args()[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Otherwise start REPL
	if err := rt.StartREPL(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
