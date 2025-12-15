package errors

import (
	"errors"
	"fmt"
)

// REPL errors
var (
	ErrInterrupt = errors.New("interrupted")
	ErrExit      = errors.New("exit")
)

// Runtime errors
var (
	ErrRuntimeInit   = errors.New("failed to initialize runtime")
	ErrBuiltinInit   = errors.New("failed to initialize builtins")
	ErrConsoleInit   = errors.New("failed to initialize console")
	ErrEvalFailed    = errors.New("evaluation failed")
	ErrFileNotFound  = errors.New("file not found")
	ErrFileRead      = errors.New("failed to read file")
	ErrInvalidScript = errors.New("invalid script")
)

// Module loader errors
var (
	ErrEmptyURL           = errors.New("empty URL provided")
	ErrInvalidURL         = errors.New("invalid URL format")
	ErrUnsupportedModule  = errors.New("unsupported module type")
	ErrModuleNotFound     = errors.New("module not found")
	ErrCircularDependency = errors.New("circular dependency detected")
	ErrJSRNotImplemented  = errors.New("JSR module loading not implemented yet")
)

// NPM errors
var (
	ErrPackageRequired = errors.New("package name is required")
	ErrPackageNotFound = errors.New("package not found")
	ErrPackageInstall  = errors.New("failed to install package")
	ErrPackageFetch    = errors.New("failed to fetch package metadata")
	ErrCacheDir        = errors.New("failed to create cache directory")
)

// Server errors
var (
	ErrServerInit     = errors.New("failed to initialize server")
	ErrServerStart    = errors.New("failed to start server")
	ErrInvalidRequest = errors.New("invalid request")
)

// Wrap wraps an error with additional context
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapWith joins a sentinel/category error with a cause (both discoverable via errors.Is)
func WrapWith(sentinel, cause error, msg string) error {
	if sentinel == nil && cause == nil {
		return nil
	}
	joined := errors.Join(sentinel, cause)
	if msg == "" {
		return joined
	}
	return fmt.Errorf("%s: %w", msg, joined)
}

// Is checks if an error matches a target error
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target any) bool {
	return errors.As(err, target)
}
