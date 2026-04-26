// Package remediator provides workflow remediation and generation functionality.
package remediator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

// WorkflowTemplate defines a workflow template for generation.
type WorkflowTemplate struct {
	Name        string
	Filename    string
	Language    string
	Type        string
	Template    string
	Description string
}

// GeneratorConfig configures the workflow generator.
type GeneratorConfig struct {
	RefRepo     string // e.g., "plexusone/.github"
	RefBranch   string // e.g., "main"
	DryRun      bool
	Verbose     bool
	OutputDir   string   // Optional: override output directory
	PathFilters []string // Optional: path filters for workflows
}

// Generator generates compliant workflow files.
type Generator struct {
	Config    GeneratorConfig
	Templates map[string]*WorkflowTemplate
}

// NewGenerator creates a new workflow generator.
func NewGenerator(cfg GeneratorConfig) *Generator {
	if cfg.RefBranch == "" {
		cfg.RefBranch = "main"
	}

	g := &Generator{
		Config:    cfg,
		Templates: make(map[string]*WorkflowTemplate),
	}

	// Register built-in templates
	g.registerGoTemplates()

	return g
}

// registerGoTemplates registers Go workflow templates.
func (g *Generator) registerGoTemplates() {
	g.Templates["go-ci"] = &WorkflowTemplate{
		Name:        "Go CI",
		Filename:    "go-ci.yaml",
		Language:    "Go",
		Type:        "go-ci",
		Description: "Go CI pipeline with build, test, and coverage",
		Template:    goCI,
	}

	g.Templates["go-lint"] = &WorkflowTemplate{
		Name:        "Go Lint",
		Filename:    "go-lint.yaml",
		Language:    "Go",
		Type:        "go-lint",
		Description: "Go linting with golangci-lint",
		Template:    goLint,
	}

	g.Templates["go-sast-codeql"] = &WorkflowTemplate{
		Name:        "CodeQL",
		Filename:    "go-sast-codeql.yaml",
		Language:    "Go",
		Type:        "go-sast-codeql",
		Description: "Go static analysis with CodeQL",
		Template:    goSASTCodeQL,
	}
}

// TemplateData contains data for template rendering.
type TemplateData struct {
	RefRepo     string
	RefBranch   string
	RepoName    string
	RepoOwner   string
	PathFilters []string
}

// GenerateForRepo generates missing workflow files for a repository.
func (g *Generator) GenerateForRepo(repo model.Repo, missing []model.MissingWorkflow) ([]GeneratedFile, error) {
	var generated []GeneratedFile

	repoPath := repo.LocalPath
	if repoPath == "" {
		return nil, fmt.Errorf("repo %s has no local path", repo.FullName)
	}

	workflowsDir := filepath.Join(repoPath, ".github", "workflows")

	// Create workflows directory if it doesn't exist
	if !g.Config.DryRun {
		if err := os.MkdirAll(workflowsDir, 0755); err != nil {
			return nil, fmt.Errorf("creating workflows directory: %w", err)
		}
	}

	data := TemplateData{
		RefRepo:     g.Config.RefRepo,
		RefBranch:   g.Config.RefBranch,
		RepoName:    repo.Name,
		RepoOwner:   repo.Owner,
		PathFilters: g.Config.PathFilters,
	}

	for _, m := range missing {
		tmpl, ok := g.Templates[m.WorkflowType]
		if !ok {
			if g.Config.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: no template for workflow type %s\n", m.WorkflowType)
			}
			continue
		}

		content, err := g.renderTemplate(tmpl, data)
		if err != nil {
			return nil, fmt.Errorf("rendering template %s: %w", tmpl.Name, err)
		}

		outputPath := filepath.Join(workflowsDir, tmpl.Filename)

		gf := GeneratedFile{
			Path:         outputPath,
			RelativePath: filepath.Join(".github", "workflows", tmpl.Filename),
			Content:      content,
			WorkflowType: m.WorkflowType,
			IsNew:        true,
		}

		// Check if file already exists
		if _, err := os.Stat(outputPath); err == nil {
			gf.IsNew = false
			gf.WouldOverwrite = true
		}

		if !g.Config.DryRun {
			// Workflow files should be world-readable (0644)
			if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil { //nolint:gosec // G306: workflow files need to be readable
				return nil, fmt.Errorf("writing %s: %w", outputPath, err)
			}
		}

		generated = append(generated, gf)
	}

	return generated, nil
}

// GeneratedFile represents a generated workflow file.
type GeneratedFile struct {
	Path           string
	RelativePath   string
	Content        string
	WorkflowType   string
	IsNew          bool
	WouldOverwrite bool
}

// renderTemplate renders a workflow template with the given data.
func (g *Generator) renderTemplate(wt *WorkflowTemplate, data TemplateData) (string, error) {
	tmpl, err := template.New(wt.Name).Parse(wt.Template)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// GetTemplate returns a template by workflow type.
func (g *Generator) GetTemplate(workflowType string) (*WorkflowTemplate, bool) {
	tmpl, ok := g.Templates[workflowType]
	return tmpl, ok
}

// ListTemplates returns all available templates.
func (g *Generator) ListTemplates() []*WorkflowTemplate {
	var templates []*WorkflowTemplate
	for _, t := range g.Templates {
		templates = append(templates, t)
	}
	return templates
}

// Go workflow templates

const goCI = `name: Go CI

on:
  push:
    branches: [main]
    paths:
      - "**.go"
      - "go.mod"
      - "go.sum"
      - ".github/workflows/go-ci.yaml"
  pull_request:
    branches: [main]
    paths:
      - "**.go"
      - "go.mod"
      - "go.sum"
      - ".github/workflows/go-ci.yaml"

jobs:
  ci:
    uses: {{.RefRepo}}/.github/workflows/go-ci.yaml@{{.RefBranch}}
`

const goLint = `name: Go Lint

on:
  push:
    branches: [main]
    paths:
      - "**.go"
      - "go.mod"
      - "go.sum"
      - ".golangci.yml"
      - ".golangci.yaml"
      - ".github/workflows/go-lint.yaml"
  pull_request:
    branches: [main]
    paths:
      - "**.go"
      - "go.mod"
      - "go.sum"
      - ".golangci.yml"
      - ".golangci.yaml"
      - ".github/workflows/go-lint.yaml"

jobs:
  lint:
    uses: {{.RefRepo}}/.github/workflows/go-lint.yaml@{{.RefBranch}}
`

const goSASTCodeQL = `name: CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: "0 6 * * 1"

permissions:
  actions: read
  contents: read
  security-events: write

jobs:
  analyze:
    uses: {{.RefRepo}}/.github/workflows/go-sast-codeql.yaml@{{.RefBranch}}
`
