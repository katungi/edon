package loader

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/katungi/edon/internal/errors"
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
	cache      *ModuleCache
	httpClient *http.Client
}

// NewModuleLoader creates a new instance of ModuleLoader
func NewModuleLoader() *ModuleLoader {
	// #81: Don't use default HTTP client - configure timeouts
	return &ModuleLoader{
		cache: &ModuleCache{
			modules: make(map[string]*Module),
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
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
		return nil, errors.ErrUnsupportedModule
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
		return nil, errors.Wrap(errors.ErrModuleNotFound, err.Error())
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, errors.Wrap(errors.ErrFileRead, err.Error())
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
		return nil, errors.Wrap(errors.ErrModuleNotFound, err.Error())
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrModuleNotFound, err.Error())
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.ErrFileRead, err.Error())
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
		return nil, errors.Wrap(errors.ErrPackageInstall, err.Error())
	}

	// Install the package
	packagePath, err := pm.InstallPackage(ctx, packageName)
	if err != nil {
		return nil, errors.Wrap(errors.ErrPackageInstall, err.Error())
	}

	// Read the package's main file
	mainFile := filepath.Join(packagePath, "index.js")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return nil, errors.Wrap(errors.ErrFileRead, err.Error())
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
	return nil, errors.ErrJSRNotImplemented
}
