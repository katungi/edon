package loader

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NPMPackageManager handles NPM package installation and caching
type NPMPackageManager struct {
	cacheDir string
}

// NewNPMPackageManager creates a new instance of NPMPackageManager
func NewNPMPackageManager() (*NPMPackageManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	cacheDir := filepath.Join(homeDir, ".edon", "npm-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %v", err)
	}

	return &NPMPackageManager{
		cacheDir: cacheDir,
	}, nil
}

// InstallPackage installs an NPM package and returns its local path
func (pm *NPMPackageManager) InstallPackage(ctx context.Context, packageName string) (string, error) {
	// Parse package name and version
	parts := strings.Split(packageName, "@")
	name := parts[0]
	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	// Check if package is already cached
	cachePath := filepath.Join(pm.cacheDir, name, version)
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	// Fetch package metadata from NPM registry
	registryURL := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	resp, err := http.Get(registryURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch package metadata: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch package: HTTP %d", resp.StatusCode)
	}

	// Create cache directory for package
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create package cache directory: %v", err)
	}

	// Download and extract package
	// TODO: Implement package download and extraction
	// For now, just create a placeholder file
	placeholder := filepath.Join(cachePath, "index.js")
	if err := ioutil.WriteFile(placeholder, []byte("// TODO: Implement package content"), 0644); err != nil {
		return "", fmt.Errorf("failed to write package file: %v", err)
	}

	return cachePath, nil
}