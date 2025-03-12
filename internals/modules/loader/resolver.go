package loader

import (
	"fmt"
	"sync"
)

// DependencyGraph represents a directed graph of module dependencies
type DependencyGraph struct {
	mu    sync.RWMutex
	edges map[string][]string // maps module URL to its dependencies
}

// NewDependencyGraph creates a new instance of DependencyGraph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		edges: make(map[string][]string),
	}
}

// AddDependency adds a dependency edge from parent to child module
func (g *DependencyGraph) AddDependency(parent, child string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check for circular dependency before adding
	if g.wouldCreateCycle(parent, child, make(map[string]bool)) {
		return fmt.Errorf("circular dependency detected: %s -> %s", parent, child)
	}

	// Add the dependency
	g.edges[parent] = append(g.edges[parent], child)
	return nil
}

// wouldCreateCycle checks if adding a new dependency would create a cycle
func (g *DependencyGraph) wouldCreateCycle(start, current string, visited map[string]bool) bool {
	if start == current {
		return true
	}

	if visited[current] {
		return false
	}
	visited[current] = true

	for _, dep := range g.edges[current] {
		if g.wouldCreateCycle(start, dep, visited) {
			return true
		}
	}

	return false
}

// GetDependencies returns all dependencies for a given module
func (g *DependencyGraph) GetDependencies(moduleURL string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.edges[moduleURL]
}

// ResolveDependencies returns a topologically sorted list of modules to load
func (g *DependencyGraph) ResolveDependencies(moduleURL string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	result := make([]string, 0)

	var visit func(string) error
	visit = func(current string) error {
		if visited[current] {
			return nil
		}

		visited[current] = true

		// Visit all dependencies first
		for _, dep := range g.edges[current] {
			if err := visit(dep); err != nil {
				return err
			}
		}

		result = append(result, current)
		return nil
	}

	if err := visit(moduleURL); err != nil {
		return nil, err
	}

	return result, nil
}