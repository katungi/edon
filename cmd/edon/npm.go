package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/katungi/edon/internals/modules/loader"
)

var (
	installCmd = flag.NewFlagSet("install", flag.ExitOnError)
)

func handleInstall() error {
	if installCmd.NArg() < 1 {
		return fmt.Errorf("package name is required")
	}

	pm, err := loader.NewNPMPackageManager()
	if err != nil {
		return fmt.Errorf("failed to initialize NPM package manager: %v", err)
	}

	for _, pkg := range installCmd.Args() {
		fmt.Printf("Installing %s...\n", pkg)
		path, err := pm.InstallPackage(context.Background(), pkg)
		if err != nil {
			return fmt.Errorf("failed to install %s: %v", pkg, err)
		}
		fmt.Printf("Successfully installed %s at %s\n", pkg, path)
	}

	return nil
}