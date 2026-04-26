package remediator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

func TestNewGenerator(t *testing.T) {
	cfg := GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
	}

	gen := NewGenerator(cfg)

	if gen.Config.RefRepo != "testorg/.github" {
		t.Errorf("RefRepo = %q, want %q", gen.Config.RefRepo, "testorg/.github")
	}
	if gen.Config.RefBranch != "main" {
		t.Errorf("RefBranch = %q, want %q", gen.Config.RefBranch, "main")
	}
}

func TestNewGenerator_DefaultBranch(t *testing.T) {
	cfg := GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "", // Should default to "main"
	}

	gen := NewGenerator(cfg)

	if gen.Config.RefBranch != "main" {
		t.Errorf("RefBranch = %q, want %q (default)", gen.Config.RefBranch, "main")
	}
}

func TestGenerator_GetTemplate(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{RefRepo: "testorg/.github"})

	tests := []struct {
		workflowType string
		wantOK       bool
		wantFilename string
	}{
		{"go-ci", true, "go-ci.yaml"},
		{"go-lint", true, "go-lint.yaml"},
		{"go-sast-codeql", true, "go-sast-codeql.yaml"},
		{"unknown", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.workflowType, func(t *testing.T) {
			tmpl, ok := gen.GetTemplate(tt.workflowType)

			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}

			if ok && tmpl.Filename != tt.wantFilename {
				t.Errorf("Filename = %q, want %q", tmpl.Filename, tt.wantFilename)
			}
		})
	}
}

func TestGenerator_ListTemplates(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{RefRepo: "testorg/.github"})

	templates := gen.ListTemplates()

	// Should have at least the 3 Go templates
	if len(templates) < 3 {
		t.Errorf("len(templates) = %d, want >= 3", len(templates))
	}

	// Verify all templates have required fields
	for _, tmpl := range templates {
		if tmpl.Name == "" {
			t.Error("template has empty Name")
		}
		if tmpl.Filename == "" {
			t.Error("template has empty Filename")
		}
		if tmpl.Type == "" {
			t.Error("template has empty Type")
		}
		if tmpl.Template == "" {
			t.Error("template has empty Template")
		}
	}
}

func TestGenerator_RenderTemplate(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "myorg/.github",
		RefBranch: "develop",
	})

	tmpl, ok := gen.GetTemplate("go-ci")
	if !ok {
		t.Fatal("go-ci template not found")
	}

	data := TemplateData{
		RefRepo:   "myorg/.github",
		RefBranch: "develop",
		RepoName:  "my-service",
		RepoOwner: "myorg",
	}

	content, err := gen.renderTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("renderTemplate error: %v", err)
	}

	// Check that template variables were replaced
	if !strings.Contains(content, "myorg/.github/.github/workflows/go-ci.yaml@develop") {
		t.Error("content does not contain expected reusable workflow reference")
	}
	if strings.Contains(content, "{{") {
		t.Error("content contains unrendered template variables")
	}
}

func TestGenerator_GenerateForRepo_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
		DryRun:    true,
	})

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	missing := []model.MissingWorkflow{
		{
			WorkflowType: "go-ci",
			Language:     "Go",
		},
		{
			WorkflowType: "go-lint",
			Language:     "Go",
		},
	}

	generated, err := gen.GenerateForRepo(repo, missing)
	if err != nil {
		t.Fatalf("GenerateForRepo error: %v", err)
	}

	if len(generated) != 2 {
		t.Errorf("len(generated) = %d, want 2", len(generated))
	}

	// In dry-run mode, files should NOT be created
	workflowsDir := filepath.Join(repoDir, ".github", "workflows")
	if _, err := os.Stat(workflowsDir); !os.IsNotExist(err) {
		t.Error("workflows directory was created in dry-run mode")
	}

	// Check generated file metadata
	for _, gf := range generated {
		if gf.Content == "" {
			t.Errorf("generated file %s has empty content", gf.WorkflowType)
		}
		if !gf.IsNew {
			t.Errorf("generated file %s should be marked as new", gf.WorkflowType)
		}
	}
}

func TestGenerator_GenerateForRepo_WriteFiles(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
		DryRun:    false,
	})

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	missing := []model.MissingWorkflow{
		{
			WorkflowType: "go-ci",
			Language:     "Go",
		},
	}

	generated, err := gen.GenerateForRepo(repo, missing)
	if err != nil {
		t.Fatalf("GenerateForRepo error: %v", err)
	}

	if len(generated) != 1 {
		t.Fatalf("len(generated) = %d, want 1", len(generated))
	}

	// File should be created
	expectedPath := filepath.Join(repoDir, ".github", "workflows", "go-ci.yaml")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}

	if len(content) == 0 {
		t.Error("generated file is empty")
	}

	// Verify content has correct reusable workflow reference
	if !strings.Contains(string(content), "testorg/.github/.github/workflows/go-ci.yaml@main") {
		t.Error("generated file does not contain expected workflow reference")
	}
}

func TestGenerator_GenerateForRepo_ExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "test-repo")
	workflowsDir := filepath.Join(repoDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create existing file
	existingContent := "existing content"
	existingPath := filepath.Join(workflowsDir, "go-ci.yaml")
	if err := os.WriteFile(existingPath, []byte(existingContent), 0600); err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
		DryRun:    true, // Don't actually overwrite
	})

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	missing := []model.MissingWorkflow{
		{
			WorkflowType: "go-ci",
			Language:     "Go",
		},
	}

	generated, err := gen.GenerateForRepo(repo, missing)
	if err != nil {
		t.Fatalf("GenerateForRepo error: %v", err)
	}

	if len(generated) != 1 {
		t.Fatalf("len(generated) = %d, want 1", len(generated))
	}

	gf := generated[0]
	if gf.IsNew {
		t.Error("IsNew = true, want false (file exists)")
	}
	if !gf.WouldOverwrite {
		t.Error("WouldOverwrite = false, want true")
	}
}

func TestGenerator_GenerateForRepo_NoLocalPath(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
	})

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: "", // Missing local path
	}

	missing := []model.MissingWorkflow{
		{WorkflowType: "go-ci"},
	}

	_, err := gen.GenerateForRepo(repo, missing)
	if err == nil {
		t.Error("expected error for repo without local path")
	}
}

func TestGenerator_GenerateForRepo_UnknownWorkflowType(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
		DryRun:    true,
	})

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	missing := []model.MissingWorkflow{
		{
			WorkflowType: "unknown-workflow",
			Language:     "Unknown",
		},
	}

	generated, err := gen.GenerateForRepo(repo, missing)
	if err != nil {
		t.Fatalf("GenerateForRepo error: %v", err)
	}

	// Unknown workflow type should be skipped
	if len(generated) != 0 {
		t.Errorf("len(generated) = %d, want 0 (unknown workflow should be skipped)", len(generated))
	}
}

func TestWorkflowTemplate_Fields(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{RefRepo: "testorg/.github"})

	expectedTemplates := []struct {
		typeName string
		filename string
		language string
		name     string
	}{
		{"go-ci", "go-ci.yaml", "Go", "Go CI"},
		{"go-lint", "go-lint.yaml", "Go", "Go Lint"},
		{"go-sast-codeql", "go-sast-codeql.yaml", "Go", "CodeQL"},
	}

	for _, et := range expectedTemplates {
		t.Run(et.typeName, func(t *testing.T) {
			tmpl, ok := gen.GetTemplate(et.typeName)
			if !ok {
				t.Fatalf("template %s not found", et.typeName)
			}

			if tmpl.Filename != et.filename {
				t.Errorf("Filename = %q, want %q", tmpl.Filename, et.filename)
			}
			if tmpl.Language != et.language {
				t.Errorf("Language = %q, want %q", tmpl.Language, et.language)
			}
			if tmpl.Name != et.name {
				t.Errorf("Name = %q, want %q", tmpl.Name, et.name)
			}
			if tmpl.Type != et.typeName {
				t.Errorf("Type = %q, want %q", tmpl.Type, et.typeName)
			}
		})
	}
}

func TestGoTemplates_ValidYAML(t *testing.T) {
	gen := NewGenerator(GeneratorConfig{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
	})

	data := TemplateData{
		RefRepo:   "testorg/.github",
		RefBranch: "main",
		RepoName:  "test-repo",
		RepoOwner: "testorg",
	}

	templates := gen.ListTemplates()
	for _, tmpl := range templates {
		t.Run(tmpl.Type, func(t *testing.T) {
			content, err := gen.renderTemplate(tmpl, data)
			if err != nil {
				t.Fatalf("renderTemplate error: %v", err)
			}

			// Basic YAML validation - check for required fields
			if !strings.Contains(content, "name:") {
				t.Error("generated workflow missing 'name:' field")
			}
			if !strings.Contains(content, "on:") {
				t.Error("generated workflow missing 'on:' field")
			}
			if !strings.Contains(content, "jobs:") {
				t.Error("generated workflow missing 'jobs:' field")
			}
			if !strings.Contains(content, "uses:") {
				t.Error("generated workflow missing 'uses:' field for reusable workflow")
			}
		})
	}
}
