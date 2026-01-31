package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Loader loads Cedar policies from various sources.
type Loader struct {
	engine *Engine
}

// NewLoader creates a new policy loader.
func NewLoader(engine *Engine) *Loader {
	return &Loader{engine: engine}
}

// LoadFromDirectory loads all .cedar files from a directory.
func (l *Loader) LoadFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively load from subdirectories
			subdir := filepath.Join(dir, entry.Name())
			if err := l.LoadFromDirectory(subdir); err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".cedar") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if err := l.LoadFromFile(path); err != nil {
			return err
		}
	}

	return nil
}

// LoadFromFile loads a single Cedar policy file.
func (l *Loader) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", path, err)
	}

	// Use filename without extension as policy ID
	base := filepath.Base(path)
	id := strings.TrimSuffix(base, ".cedar")

	// Include parent directory in ID for uniqueness
	dir := filepath.Dir(path)
	parentDir := filepath.Base(dir)
	if parentDir != "." && parentDir != "/" {
		id = fmt.Sprintf("%s/%s", parentDir, id)
	}

	if err := l.engine.AddPolicy(id, content); err != nil {
		return fmt.Errorf("loading policy from %s: %w", path, err)
	}

	return nil
}

// LoadFromBytes loads a policy from raw bytes with a given ID.
func (l *Loader) LoadFromBytes(id string, content []byte) error {
	return l.engine.AddPolicy(id, content)
}

// LoadBuiltinPolicies loads the built-in default policies.
func (l *Loader) LoadBuiltinPolicies() error {
	policies := map[string]string{
		"builtin/require-workflow": requireWorkflowPolicy,
		"builtin/require-tests":    requireTestsPolicy,
		"builtin/go-versions":      goVersionsPolicy,
	}

	for id, content := range policies {
		if err := l.engine.AddPolicy(id, []byte(content)); err != nil {
			return fmt.Errorf("loading builtin policy %s: %w", id, err)
		}
	}

	return nil
}

// Built-in policies
const (
	// requireWorkflowPolicy denies merge if no CI workflow exists
	requireWorkflowPolicy = `
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == true
};
`

	// requireTestsPolicy denies merge if tests haven't passed
	requireTestsPolicy = `
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.lastRunPassed == true
};
`

	// goVersionsPolicy permits build only with approved Go versions
	goVersionsPolicy = `
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.goVersions.containsAny(["1.24", "1.25"])
};
`
)
