package loader

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/katungi/edon/internal/errors"
)

// NPMPackageManager handles NPM package installation and caching
type NPMPackageManager struct {
	cacheDir string
}

// NewNPMPackageManager creates a new instance of NPMPackageManager
func NewNPMPackageManager() (*NPMPackageManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(errors.ErrCacheDir, err.Error())
	}

	cacheDir := filepath.Join(homeDir, ".edon", "npm-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, errors.Wrap(errors.ErrCacheDir, err.Error())
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
		return "", errors.Wrap(errors.ErrPackageFetch, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.ErrPackageNotFound
	}

	// Create cache directory for package
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", errors.Wrap(errors.ErrCacheDir, err.Error())
	}

	// Download and extract package
	// TODO: Implement package download and extraction
	// For now, just create a placeholder file
	placeholder := filepath.Join(cachePath, "index.js")
	if err := os.WriteFile(placeholder, []byte("// TODO: Implement package content"), 0644); err != nil {
		return "", errors.Wrap(errors.ErrPackageInstall, err.Error())
	}

	return cachePath, nil
}
