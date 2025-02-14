package loader

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type PackageType string

const (
	TypeJSR   PackageType = "JSR"
	TypeNPM   PackageType = "NPM"
	TypeCDN   PackageType = "CDN"
	TypeLocal PackageType = "Local"
)

type ValidationResult struct {
	IsValid     bool
	PackageType PackageType
	Error       error
}

func ValidateURL(urlStr string) ValidationResult {
	// Handle empty input
	if urlStr == "" {
		return ValidationResult{
			IsValid: false,
			Error:   fmt.Errorf("empty URL provided"),
		}
	}

	// Check for package prefixes
	if strings.HasPrefix(urlStr, "npm:") {
		return ValidationResult{
			IsValid:     true,
			PackageType: TypeNPM,
		}
	}

	if strings.HasPrefix(urlStr, "jsr:") {
		return ValidationResult{
			IsValid:     true,
			PackageType: TypeJSR,
		}
	}

	// Check if it's a local file path
	if isLocalPath(urlStr) {
		return ValidationResult{
			IsValid:     true,
			PackageType: TypeLocal,
		}
	}

	// Parse URL for CDN validation
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ValidationResult{
			IsValid: false,
			Error:   fmt.Errorf("invalid URL format: %v", err),
		}
	}

	// Validate CDN URLs
	if isCDNURL(parsedURL) {
		return ValidationResult{
			IsValid:     true,
			PackageType: TypeCDN,
		}
	}

	return ValidationResult{
		IsValid: false,
		Error:   fmt.Errorf("URL does not match any supported package type"),
	}
}

func isLocalPath(path string) bool {
	// Check if the path is absolute
	if filepath.IsAbs(path) {
		return true
	}

	// Check if path starts with ./ or ../
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		return true
	}

	// Check if it's a Windows-style absolute path
	if len(path) >= 2 && path[1] == ':' {
		return true
	}

	return false
}

func isCDNURL(parsedURL *url.URL) bool {
	// List of known CDN domains
	cdnDomains := []string{
		"cdn.jsdelivr.net",
		"unpkg.com",
		"cdnjs.cloudflare.com",
	}

	for _, domain := range cdnDomains {
		if strings.Contains(parsedURL.Host, domain) {
			return true
		}
	}

	return false
}