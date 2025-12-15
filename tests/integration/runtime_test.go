package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/katungi/edon/internal/runtime"
)

func TestFileExecution(t *testing.T) {
	// Create a temporary test file
	content := `
        let x = 1;
        let y = 2;
        console.log(x + y);
    `
	tmpfile, err := os.CreateTemp("", "test-*.js")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test file execution
	rt, err := runtime.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	if err := rt.ExecuteFile(tmpfile.Name()); err != nil {
		t.Errorf("ExecuteFile() error = %v", err)
	}
}

func TestModuleLoading(t *testing.T) {
	// Test module loading
	moduleContent := `
        export function add(a, b) {
            return a + b;
        }
    `
	mainContent := `
        import { add } from './math.js';
        console.log(add(2, 3));
    `

	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "test-modules")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create module file
	if err := os.WriteFile(
		filepath.Join(tmpDir, "math.js"),
		[]byte(moduleContent),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.js")
	if err := os.WriteFile(
		mainFile,
		[]byte(mainContent),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	// Test execution
	rt, err := runtime.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	if err := rt.ExecuteFile(mainFile); err != nil {
		t.Errorf("Module loading test failed: %v", err)
	}
}
