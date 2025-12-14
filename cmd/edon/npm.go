package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/katungi/edon/internal/modules/loader"
)

var (
	InstallCmd = flag.NewFlagSet("install", flag.ExitOnError)
)

func HandleInstall() error {
	if InstallCmd.NArg() < 1 {
		return fmt.Errorf("package name is required")
	}

	pm, err := loader.NewNPMPackageManager()
	if err != nil {
		return fmt.Errorf("failed to initialize NPM package manager: %v", err)
	}

	for _, pkg := range InstallCmd.Args() {
		fmt.Printf("Installing %s...\n", pkg)
		path, err := pm.InstallPackage(context.Background(), pkg)
		if err != nil {
			return fmt.Errorf("failed to install %s: %v", pkg, err)
		}
		fmt.Printf("Successfully installed %s at %s\n", pkg, path)
	}

	return nil
}
