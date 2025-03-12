package loader

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// ModuleCache represents a thread-safe cache for loaded modules
type ModuleCache struct {
	mu      sync.RWMutex
	modules map[string]*Module
}

// Module represents a loaded module with its content and metadata
type Module struct {
	URL     string
	Content string
	Type    PackageType
}

// ModuleLoader handles the loading of modules from various sources
type ModuleLoader struct {
	cache *ModuleCache
}

// NewModuleLoader creates a new instance of ModuleLoader
func NewModuleLoader() *ModuleLoader {
	return &ModuleLoader{
		cache: &ModuleCache{
			modules: make(map[string]*Module),
		},
	}
}

// LoadModule loads a module from the given URL, using cache if available
func (l *ModuleLoader) LoadModule(ctx context.Context, urlStr string) (*Module, error) {
	// Validate the URL first
	validation := ValidateURL(urlStr)
	if !validation.IsValid {
		return nil, validation.Error
	}

	// Check cache first
	if module := l.getFromCache(urlStr); module != nil {
		return module, nil
	}

	// Load module based on its type
	var module *Module
	var err error

	switch validation.PackageType {
	case TypeLocal:
		module, err = l.loadLocalModule(urlStr)
	case TypeCDN:
		module, err = l.loadCDNModule(ctx, urlStr)
	case TypeNPM:
		module, err = l.loadNPMModule(ctx, urlStr)
	case TypeJSR:
		module, err = l.loadJSRModule(ctx, urlStr)
	default:
		return nil, fmt.Errorf("unsupported module type")
	}

	if err != nil {
		return nil, err
	}

	// Cache the loaded module
	l.cache.mu.Lock()
	l.cache.modules[urlStr] = module
	l.cache.mu.Unlock()

	return module, nil
}

// getFromCache retrieves a module from the cache if it exists
func (l *ModuleLoader) getFromCache(url string) *Module {
	l.cache.mu.RLock()
	defer l.cache.mu.RUnlock()
	return l.cache.modules[url]
}

// loadLocalModule loads a module from the local filesystem
func (l *ModuleLoader) loadLocalModule(path string) (*Module, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %v", err)
	}

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local file: %v", err)
	}

	return &Module{
		URL:     path,
		Content: string(content),
		Type:    TypeLocal,
	}, nil
}

// loadCDNModule loads a module from a CDN
func (l *ModuleLoader) loadCDNModule(ctx context.Context, url string) (*Module, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from CDN: %v", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return &Module{
		URL:     url,
		Content: string(content),
		Type:    TypeCDN,
	}, nil
}

// loadNPMModule loads a module from NPM registry
func (l *ModuleLoader) loadNPMModule(ctx context.Context, url string) (*Module, error) {
	// Extract package name from npm: URL
	packageName := strings.TrimPrefix(url, "npm:")

	// Initialize NPM package manager
	pm, err := NewNPMPackageManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NPM package manager: %v", err)
	}

	// Install the package
	packagePath, err := pm.InstallPackage(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to install NPM package: %v", err)
	}

	// Read the package's main file
	mainFile := filepath.Join(packagePath, "index.js")
	content, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read package file: %v", err)
	}

	return &Module{
		URL:     url,
		Content: string(content),
		Type:    TypeNPM,
	}, nil
}

// loadJSRModule loads a module from JSR registry
func (l *ModuleLoader) loadJSRModule(ctx context.Context, url string) (*Module, error) {
	// TODO: Implement JSR module loading
	return nil, fmt.Errorf("JSR module loading not implemented yet")
}