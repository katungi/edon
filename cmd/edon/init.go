package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

var InitCmd = flag.NewFlagSet("init", flag.ExitOnError)

func HandleInit() error {
	// Get current directory or use the provided path
	dir := InitCmd.Arg(0)
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create package.json
	packageJSON := map[string]interface{}{
		"name":        filepath.Base(dir),
		"version":     "1.0.0",
		"description": "A new Edon project",
		"main":        "index.js",
		"scripts": map[string]string{
			"start": "edon index.js",
		},
	}

	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create package.json: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "package.json"), packageJSONBytes, 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	// Create index.js
	indexJSContent := "console.log('Hello from Edon!');"
	if err := os.WriteFile(filepath.Join(dir, "index.js"), []byte(indexJSContent), 0644); err != nil {
		return fmt.Errorf("failed to write index.js: %w", err)
	}

	color.Green("✓ Successfully initialized new Edon project in %s", dir)
	color.Green("✓ Created package.json")
	color.Green("✓ Created index.js")

	return nil
}