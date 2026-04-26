package compliance

import (
	"testing"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

func TestNewWorkflowMatcher(t *testing.T) {
	refRepo := &ReferenceRepo{
		Owner:  "testorg",
		Name:   ".github",
		Branch: "main",
	}

	matcher := NewWorkflowMatcher(refRepo, true)

	if matcher.RefRepo != refRepo {
		t.Error("RefRepo not set correctly")
	}
	if !matcher.Strict {
		t.Error("Strict = false, want true")
	}
}

func TestWorkflowMatcher_MatchWorkflow_ExactMatch(t *testing.T) {
	refRepo := &ReferenceRepo{
		Owner:  "testorg",
		Name:   ".github",
		Branch: "main",
	}
	matcher := NewWorkflowMatcher(refRepo, false)

	rule := WorkflowRule{
		Type:        "go-ci",
		Path:        ".github/workflows/go-ci.yaml",
		Description: "Go CI workflow",
		Severity:    "high",
	}

	// Workflow that uses the reusable workflow
	workflow := model.Workflow{
		Name: "Go CI",
		Path: ".github/workflows/go-ci.yaml",
		Content: `name: Go CI
jobs:
  ci:
    uses: testorg/.github/.github/workflows/go-ci.yaml@main
`,
	}

	result := matcher.MatchWorkflow([]model.Workflow{workflow}, rule)

	if result.MatchType != model.MatchTypeExact {
		t.Errorf("MatchType = %q, want %q", result.MatchType, model.MatchTypeExact)
	}
	if !result.UsesReusable {
		t.Error("UsesReusable = false, want true")
	}
	if result.FilenameMismatch {
		t.Error("FilenameMismatch = true, want false")
	}
}

func TestWorkflowMatcher_MatchWorkflow_EquivalentMatch(t *testing.T) {
	refRepo := &ReferenceRepo{
		Owner:  "testorg",
		Name:   ".github",
		Branch: "main",
	}
	matcher := NewWorkflowMatcher(refRepo, false)

	rule := WorkflowRule{
		Type:        "go-ci",
		Path:        ".github/workflows/go-ci.yaml",
		Description: "Go CI workflow",
		Severity:    "high",
	}

	// Workflow with Go steps but not using reusable workflow
	workflow := model.Workflow{
		Name: "CI",
		Path: ".github/workflows/ci.yaml",
		Content: `name: CI
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go build ./...
      - run: go test ./...
`,
	}

	result := matcher.MatchWorkflow([]model.Workflow{workflow}, rule)

	if result.MatchType != model.MatchTypeEquivalent {
		t.Errorf("MatchType = %q, want %q", result.MatchType, model.MatchTypeEquivalent)
	}
	if result.UsesReusable {
		t.Error("UsesReusable = true, want false")
	}
	if !result.FilenameMismatch {
		t.Error("FilenameMismatch = false, want true (ci.yaml vs go-ci.yaml)")
	}
	if result.ExpectedFilename != "go-ci.yaml" {
		t.Errorf("ExpectedFilename = %q, want %q", result.ExpectedFilename, "go-ci.yaml")
	}
	if result.ActualFilename != "ci.yaml" {
		t.Errorf("ActualFilename = %q, want %q", result.ActualFilename, "ci.yaml")
	}
}

func TestWorkflowMatcher_MatchWorkflow_NoMatch(t *testing.T) {
	refRepo := &ReferenceRepo{
		Owner:  "testorg",
		Name:   ".github",
		Branch: "main",
	}
	matcher := NewWorkflowMatcher(refRepo, false)

	rule := WorkflowRule{
		Type:        "go-ci",
		Path:        ".github/workflows/go-ci.yaml",
		Description: "Go CI workflow",
		Severity:    "high",
	}

	// Workflow that's unrelated to Go
	workflow := model.Workflow{
		Name: "Python CI",
		Path: ".github/workflows/python-ci.yaml",
		Content: `name: Python CI
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
      - run: pip install -r requirements.txt
`,
	}

	result := matcher.MatchWorkflow([]model.Workflow{workflow}, rule)

	if result.MatchType != model.MatchTypeNone {
		t.Errorf("MatchType = %q, want %q", result.MatchType, model.MatchTypeNone)
	}
}

func TestWorkflowMatcher_MatchWorkflow_StrictMode(t *testing.T) {
	refRepo := &ReferenceRepo{
		Owner:  "testorg",
		Name:   ".github",
		Branch: "main",
	}
	matcher := NewWorkflowMatcher(refRepo, true) // strict mode

	rule := WorkflowRule{
		Type:        "go-ci",
		Path:        ".github/workflows/go-ci.yaml",
		Description: "Go CI workflow",
		Severity:    "high",
	}

	// Equivalent workflow (should not match in strict mode)
	workflow := model.Workflow{
		Name: "CI",
		Path: ".github/workflows/ci.yaml",
		Content: `name: CI
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
      - run: go test ./...
`,
	}

	result := matcher.MatchWorkflow([]model.Workflow{workflow}, rule)

	if result.MatchType != model.MatchTypeNone {
		t.Errorf("MatchType = %q, want %q (strict mode should not allow equivalent)", result.MatchType, model.MatchTypeNone)
	}
}

func TestGetExpectedFilename(t *testing.T) {
	tests := []struct {
		workflowType string
		want         string
	}{
		{"go-ci", "go-ci.yaml"},
		{"go-lint", "go-lint.yaml"},
		{"go-sast-codeql", "go-sast-codeql.yaml"},
		{"ts-ci", "ts-ci.yaml"},
		{"ts-lint", "ts-lint.yaml"},
		{"unknown", "unknown.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.workflowType, func(t *testing.T) {
			got := getExpectedFilename(tt.workflowType)
			if got != tt.want {
				t.Errorf("getExpectedFilename(%q) = %q, want %q", tt.workflowType, got, tt.want)
			}
		})
	}
}

func TestGetFilenameFromPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{".github/workflows/go-ci.yaml", "go-ci.yaml"},
		{"workflows/ci.yml", "ci.yml"},
		{"ci.yaml", "ci.yaml"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := getFilenameFromPath(tt.path)
			if got != tt.want {
				t.Errorf("getFilenameFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_IsExpectedFilename(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		path     string
		ruleType string
		want     bool
	}{
		{".github/workflows/go-ci.yaml", "go-ci", true},
		{".github/workflows/go-ci.yml", "go-ci", true},
		{".github/workflows/GO-CI.YAML", "go-ci", true},
		{".github/workflows/ci.yaml", "go-ci", false},
		{".github/workflows/go-lint.yaml", "go-lint", true},
		{".github/workflows/lint.yaml", "go-lint", false},
		{".github/workflows/go-sast-codeql.yaml", "go-sast-codeql", true},
		{".github/workflows/codeql.yaml", "go-sast-codeql", false},
	}

	for _, tt := range tests {
		t.Run(tt.path+"_"+tt.ruleType, func(t *testing.T) {
			got := matcher.isExpectedFilename(tt.path, tt.ruleType)
			if got != tt.want {
				t.Errorf("isExpectedFilename(%q, %q) = %v, want %v", tt.path, tt.ruleType, got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_ContainsGoSteps(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "go build",
			content: "run: go build ./...",
			want:    true,
		},
		{
			name:    "go test",
			content: "run: go test ./...",
			want:    true,
		},
		{
			name:    "go mod",
			content: "run: go mod tidy",
			want:    true,
		},
		{
			name:    "setup-go action",
			content: "uses: actions/setup-go@v5",
			want:    true,
		},
		{
			name:    "no go",
			content: "run: npm test",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.containsGoSteps(tt.content)
			if got != tt.want {
				t.Errorf("containsGoSteps = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_ContainsGolangciLint(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "golangci-lint action",
			content: "uses: golangci/golangci-lint-action@v4",
			want:    true,
		},
		{
			name:    "golangci-lint run",
			content: "run: golangci-lint run",
			want:    true,
		},
		{
			name:    "no lint",
			content: "run: go test ./...",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.containsGolangciLint(tt.content)
			if got != tt.want {
				t.Errorf("containsGolangciLint = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_ContainsCodeQL(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "codeql action",
			content: "uses: github/codeql-action/init@v3",
			want:    true,
		},
		{
			name:    "codeql in name",
			content: "name: CodeQL Analysis",
			want:    true,
		},
		{
			name:    "no codeql",
			content: "run: go test ./...",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.containsCodeQL(tt.content)
			if got != tt.want {
				t.Errorf("containsCodeQL = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_ContainsNodeSteps(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "npm run",
			content: "run: npm run build",
			want:    true,
		},
		{
			name:    "npm test",
			content: "run: npm test",
			want:    true,
		},
		{
			name:    "yarn test",
			content: "run: yarn test",
			want:    true,
		},
		{
			name:    "setup-node action",
			content: "uses: actions/setup-node@v4",
			want:    true,
		},
		{
			name:    "no node",
			content: "run: go test ./...",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.containsNodeSteps(tt.content)
			if got != tt.want {
				t.Errorf("containsNodeSteps = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_ContainsESLint(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "npm run lint",
			content: "run: npm run lint",
			want:    true,
		},
		{
			name:    "yarn lint",
			content: "run: yarn lint",
			want:    true,
		},
		{
			name:    "eslint directly",
			content: "run: eslint .",
			want:    true,
		},
		{
			name:    "no lint",
			content: "run: npm test",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.containsESLint(tt.content)
			if got != tt.want {
				t.Errorf("containsESLint = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectWorkflowLanguage(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Go workflow",
			content: "uses: actions/setup-go@v5",
			want:    "Go",
		},
		{
			name:    "Go with lint",
			content: "uses: golangci/golangci-lint-action@v4",
			want:    "Go",
		},
		{
			name:    "TypeScript workflow",
			content: "uses: actions/setup-node@v4",
			want:    "TypeScript",
		},
		{
			name:    "Empty content",
			content: "",
			want:    "",
		},
		{
			name:    "Unknown workflow",
			content: "run: python test.py",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := model.Workflow{Content: tt.content}
			got := DetectWorkflowLanguage(wf)
			if got != tt.want {
				t.Errorf("DetectWorkflowLanguage = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWorkflowMatcher_NormalizeRef(t *testing.T) {
	matcher := &WorkflowMatcher{}

	tests := []struct {
		ref  string
		want string
	}{
		{"owner/repo/.github/workflows/ci.yaml@main", "owner/repo/.github/workflows/ci.yaml@main"},
		{"  owner/repo/.github/workflows/ci.yaml@main  ", "owner/repo/.github/workflows/ci.yaml@main"},
		{"owner/repo/.github/workflows/ci.yaml", "owner/repo/.github/workflows/ci.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			got := matcher.normalizeRef(tt.ref)
			if got != tt.want {
				t.Errorf("normalizeRef(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}
