package collector

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

func TestNewLocalCollector(t *testing.T) {
	coll := NewLocalCollector("/test/path")

	if coll.BasePath != "/test/path" {
		t.Errorf("BasePath = %q, want %q", coll.BasePath, "/test/path")
	}
	if coll.Workers <= 0 {
		t.Errorf("Workers = %d, want > 0", coll.Workers)
	}
}

func TestNewLocalCollectorWithConfig(t *testing.T) {
	cfg := LocalCollectorConfig{
		BasePath: "/test/path",
		Workers:  4,
		Verbose:  true,
	}
	coll := NewLocalCollectorWithConfig(cfg)

	if coll.BasePath != "/test/path" {
		t.Errorf("BasePath = %q, want %q", coll.BasePath, "/test/path")
	}
	if coll.Workers != 4 {
		t.Errorf("Workers = %d, want 4", coll.Workers)
	}
	if !coll.Verbose {
		t.Error("Verbose = false, want true")
	}
}

func TestNewLocalCollectorWithConfig_DefaultWorkers(t *testing.T) {
	cfg := LocalCollectorConfig{
		BasePath: "/test/path",
		Workers:  0, // Should use default
	}
	coll := NewLocalCollectorWithConfig(cfg)

	if coll.Workers <= 0 {
		t.Errorf("Workers = %d, want > 0", coll.Workers)
	}
}

func TestLocalCollector_ListRepos(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	orgDir := filepath.Join(tempDir, "testorg")
	if err := os.MkdirAll(orgDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a Go repo
	goRepo := filepath.Join(orgDir, "go-repo")
	if err := os.MkdirAll(goRepo, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goRepo, "go.mod"), []byte("module example.com/go-repo"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a TypeScript repo (needs .github dir to be detected as valid repo)
	tsRepo := filepath.Join(orgDir, "ts-repo")
	tsGitHubDir := filepath.Join(tsRepo, ".github")
	if err := os.MkdirAll(tsGitHubDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tsRepo, "package.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tsRepo, "tsconfig.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	coll := NewLocalCollector(tempDir)
	ctx := context.Background()

	t.Run("no filter", func(t *testing.T) {
		repos, err := coll.ListRepos(ctx, []string{"testorg"}, model.RepoFilter{})
		if err != nil {
			t.Fatalf("ListRepos error: %v", err)
		}

		if len(repos) != 2 {
			t.Errorf("len(repos) = %d, want 2", len(repos))
		}
	})

	t.Run("filter by Go", func(t *testing.T) {
		repos, err := coll.ListRepos(ctx, []string{"testorg"}, model.RepoFilter{
			IncludeLanguages: []string{"Go"},
		})
		if err != nil {
			t.Fatalf("ListRepos error: %v", err)
		}

		if len(repos) != 1 {
			t.Errorf("len(repos) = %d, want 1", len(repos))
		}
		if len(repos) > 0 && repos[0].Name != "go-repo" {
			t.Errorf("repo name = %q, want %q", repos[0].Name, "go-repo")
		}
	})

	t.Run("filter by TypeScript", func(t *testing.T) {
		repos, err := coll.ListRepos(ctx, []string{"testorg"}, model.RepoFilter{
			IncludeLanguages: []string{"TypeScript"},
		})
		if err != nil {
			t.Fatalf("ListRepos error: %v", err)
		}

		if len(repos) != 1 {
			t.Errorf("len(repos) = %d, want 1", len(repos))
		}
		if len(repos) > 0 && repos[0].Name != "ts-repo" {
			t.Errorf("repo name = %q, want %q", repos[0].Name, "ts-repo")
		}
	})
}

func TestLocalCollector_GetWorkflows(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "testorg", "test-repo")
	workflowsDir := filepath.Join(repoDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create go.mod
	if err := os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte("module example.com/test-repo"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a workflow file that uses a reusable workflow
	ciWorkflow := `name: Go CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  ci:
    uses: testorg/.github/.github/workflows/go-ci.yaml@main
    secrets: inherit
`
	if err := os.WriteFile(filepath.Join(workflowsDir, "go-ci.yaml"), []byte(ciWorkflow), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a workflow file with steps
	lintWorkflow := `name: Lint
on: [push]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4
`
	if err := os.WriteFile(filepath.Join(workflowsDir, "lint.yaml"), []byte(lintWorkflow), 0600); err != nil {
		t.Fatal(err)
	}

	coll := NewLocalCollector(tempDir)
	ctx := context.Background()

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	workflows, err := coll.GetWorkflows(ctx, repo)
	if err != nil {
		t.Fatalf("GetWorkflows error: %v", err)
	}

	if len(workflows) != 2 {
		t.Fatalf("len(workflows) = %d, want 2", len(workflows))
	}

	// Find the reusable workflow
	var ciWf *model.Workflow
	for i := range workflows {
		if workflows[i].Name == "Go CI" {
			ciWf = &workflows[i]
			break
		}
	}

	if ciWf == nil {
		t.Fatal("Go CI workflow not found")
	}

	if !ciWf.UsesReusableWorkflow {
		t.Error("UsesReusableWorkflow = false, want true")
	}

	if len(ciWf.ReusableWorkflowRefs) != 1 {
		t.Fatalf("len(ReusableWorkflowRefs) = %d, want 1", len(ciWf.ReusableWorkflowRefs))
	}

	ref := ciWf.ReusableWorkflowRefs[0]
	if ref.Owner != "testorg" {
		t.Errorf("ref.Owner = %q, want %q", ref.Owner, "testorg")
	}
	if ref.Repo != ".github" {
		t.Errorf("ref.Repo = %q, want %q", ref.Repo, ".github")
	}
}

func TestLocalCollector_GetWorkflows_NoWorkflowsDir(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "testorg", "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	coll := NewLocalCollector(tempDir)
	ctx := context.Background()

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	workflows, err := coll.GetWorkflows(ctx, repo)
	if err != nil {
		t.Fatalf("GetWorkflows error: %v", err)
	}

	if len(workflows) != 0 {
		t.Errorf("len(workflows) = %d, want 0", len(workflows))
	}
}

func TestLocalCollector_DetectLanguages(t *testing.T) {
	tests := []struct {
		name      string
		files     map[string]string
		wantLangs []string
	}{
		{
			name: "Go project",
			files: map[string]string{
				"go.mod": "module example.com/test",
			},
			wantLangs: []string{"Go"},
		},
		{
			name: "TypeScript project",
			files: map[string]string{
				"package.json":  "{}",
				"tsconfig.json": "{}",
			},
			wantLangs: []string{"TypeScript"},
		},
		{
			name: "JavaScript project",
			files: map[string]string{
				"package.json": "{}",
			},
			wantLangs: []string{"JavaScript"},
		},
		{
			name: "Python project",
			files: map[string]string{
				"pyproject.toml": "[project]",
			},
			wantLangs: []string{"Python"},
		},
		{
			name: "Rust project",
			files: map[string]string{
				"Cargo.toml": "[package]",
			},
			wantLangs: []string{"Rust"},
		},
		{
			name: "Crystal project",
			files: map[string]string{
				"shard.yml": "name: test",
			},
			wantLangs: []string{"Crystal"},
		},
		{
			name: "Multi-language project",
			files: map[string]string{
				"go.mod":       "module example.com/test",
				"package.json": "{}",
			},
			wantLangs: []string{"Go", "JavaScript"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for filename, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}

			coll := NewLocalCollector(tempDir)
			langs := coll.detectLanguages(tempDir)

			if len(langs) != len(tt.wantLangs) {
				t.Errorf("len(langs) = %d, want %d; got %v", len(langs), len(tt.wantLangs), langs)
				return
			}

			for _, want := range tt.wantLangs {
				if !slices.Contains(langs, want) {
					t.Errorf("language %q not found in %v", want, langs)
				}
			}
		})
	}
}

func TestLocalCollector_GetFileContent(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "testorg", "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := "test file content"
	if err := os.WriteFile(filepath.Join(repoDir, "test.txt"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	coll := NewLocalCollector(tempDir)
	ctx := context.Background()

	repo := model.Repo{
		Owner:     "testorg",
		Name:      "test-repo",
		FullName:  "testorg/test-repo",
		LocalPath: repoDir,
	}

	got, err := coll.GetFileContent(ctx, repo, "test.txt")
	if err != nil {
		t.Fatalf("GetFileContent error: %v", err)
	}

	if got != content {
		t.Errorf("content = %q, want %q", got, content)
	}
}

func TestExtractTriggers(t *testing.T) {
	tests := []struct {
		name string
		on   any
		want []string
	}{
		{
			name: "string trigger",
			on:   "push",
			want: []string{"push"},
		},
		{
			name: "array triggers",
			on:   []any{"push", "pull_request"},
			want: []string{"push", "pull_request"},
		},
		{
			name: "map triggers",
			on: map[string]any{
				"push":         map[string]any{"branches": []string{"main"}},
				"pull_request": nil,
			},
			want: []string{"push", "pull_request"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTriggers(tt.on)

			if len(got) != len(tt.want) {
				t.Errorf("len(triggers) = %d, want %d", len(got), len(tt.want))
				return
			}

			for _, want := range tt.want {
				if !slices.Contains(got, want) {
					t.Errorf("trigger %q not found in %v", want, got)
				}
			}
		})
	}
}

func TestExtractRunsOn(t *testing.T) {
	tests := []struct {
		name   string
		runsOn any
		want   []string
	}{
		{
			name:   "string",
			runsOn: "ubuntu-latest",
			want:   []string{"ubuntu-latest"},
		},
		{
			name:   "array",
			runsOn: []any{"ubuntu-latest", "macos-latest"},
			want:   []string{"ubuntu-latest", "macos-latest"},
		},
		{
			name:   "nil",
			runsOn: nil,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRunsOn(tt.runsOn)

			if len(got) != len(tt.want) {
				t.Errorf("len(runsOn) = %d, want %d", len(got), len(tt.want))
				return
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("runsOn[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseGoMod(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	content := `module github.com/example/test

go 1.23

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.0
)

require github.com/stretchr/testify v1.9.0

replace github.com/old/pkg => github.com/new/pkg v1.0.0

replace (
	example.com/a => example.com/b v1.0.0
	example.com/c => example.com/d v1.0.0
)
`
	if err := os.WriteFile(goModPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	moduleName, replaceCount, dependencies := ParseGoMod(goModPath)

	if moduleName != "github.com/example/test" {
		t.Errorf("moduleName = %q, want %q", moduleName, "github.com/example/test")
	}

	if replaceCount != 3 {
		t.Errorf("replaceCount = %d, want 3", replaceCount)
	}

	expectedDeps := []string{
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
		"github.com/stretchr/testify",
	}

	if len(dependencies) != len(expectedDeps) {
		t.Errorf("len(dependencies) = %d, want %d", len(dependencies), len(expectedDeps))
	}

	for _, want := range expectedDeps {
		if !slices.Contains(dependencies, want) {
			t.Errorf("dependency %q not found", want)
		}
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	// File doesn't exist
	if fileExists(filePath) {
		t.Error("fileExists = true for non-existent file")
	}

	// Create file
	if err := os.WriteFile(filePath, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}

	// File exists
	if !fileExists(filePath) {
		t.Error("fileExists = false for existing file")
	}

	// Directory should return false
	if fileExists(tempDir) {
		t.Error("fileExists = true for directory")
	}
}

func TestDirExists(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")

	// Directory doesn't exist
	if dirExists(subDir) {
		t.Error("dirExists = true for non-existent directory")
	}

	// Create directory
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Directory exists
	if !dirExists(subDir) {
		t.Error("dirExists = false for existing directory")
	}

	// File should return false
	filePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
	if dirExists(filePath) {
		t.Error("dirExists = true for file")
	}
}
