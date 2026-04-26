package compliance

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

// WorkflowMatcher checks workflows for compliance.
type WorkflowMatcher struct {
	RefRepo *ReferenceRepo
	Strict  bool
}

// NewWorkflowMatcher creates a new workflow matcher.
func NewWorkflowMatcher(refRepo *ReferenceRepo, strict bool) *WorkflowMatcher {
	return &WorkflowMatcher{
		RefRepo: refRepo,
		Strict:  strict,
	}
}

// MatchResult represents the result of matching a workflow.
type MatchResult struct {
	MatchType        string
	UsesReusable     bool
	ReusableRef      string
	ActualWorkflow   string
	FilenameMismatch bool
	ExpectedFilename string
	ActualFilename   string
}

// MatchWorkflow checks if a repository has a matching workflow for the given rule.
func (m *WorkflowMatcher) MatchWorkflow(workflows []model.Workflow, rule WorkflowRule) *MatchResult {
	expectedRef := m.RefRepo.WorkflowRef(rule.Path)
	expectedFilename := getExpectedFilename(rule.Type)

	// First, check for exact reusable workflow usage
	for _, wf := range workflows {
		if m.usesReusableWorkflow(wf, expectedRef) {
			result := &MatchResult{
				MatchType:        model.MatchTypeExact,
				UsesReusable:     true,
				ReusableRef:      expectedRef,
				ActualWorkflow:   wf.Path,
				ExpectedFilename: expectedFilename,
				ActualFilename:   getFilenameFromPath(wf.Path),
			}
			// Check for filename mismatch even on exact matches
			if !m.isExpectedFilename(wf.Path, rule.Type) {
				result.FilenameMismatch = true
			}
			return result
		}
	}

	// If strict mode, no equivalent matching allowed
	if m.Strict {
		return &MatchResult{
			MatchType:        model.MatchTypeNone,
			ExpectedFilename: expectedFilename,
		}
	}

	// Check for equivalent workflow by type
	for _, wf := range workflows {
		if m.isEquivalentWorkflow(wf, rule.Type) {
			result := &MatchResult{
				MatchType:        model.MatchTypeEquivalent,
				UsesReusable:     false,
				ActualWorkflow:   wf.Path,
				ExpectedFilename: expectedFilename,
				ActualFilename:   getFilenameFromPath(wf.Path),
			}
			// Check for filename mismatch
			if !m.isExpectedFilename(wf.Path, rule.Type) {
				result.FilenameMismatch = true
			}
			return result
		}
	}

	return &MatchResult{
		MatchType:        model.MatchTypeNone,
		ExpectedFilename: expectedFilename,
	}
}

// getExpectedFilename returns the expected filename for a workflow type.
func getExpectedFilename(workflowType string) string {
	switch workflowType {
	case "go-ci":
		return "go-ci.yaml"
	case "go-lint":
		return "go-lint.yaml"
	case "go-sast-codeql":
		return "go-sast-codeql.yaml"
	case "ts-ci":
		return "ts-ci.yaml"
	case "ts-lint":
		return "ts-lint.yaml"
	default:
		return workflowType + ".yaml"
	}
}

// getFilenameFromPath extracts the filename from a workflow path.
func getFilenameFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// isExpectedFilename checks if a workflow path has the expected filename.
func (m *WorkflowMatcher) isExpectedFilename(path string, ruleType string) bool {
	filename := strings.ToLower(getFilenameFromPath(path))

	switch ruleType {
	case "go-ci":
		return filename == "go-ci.yaml" || filename == "go-ci.yml"
	case "go-lint":
		return filename == "go-lint.yaml" || filename == "go-lint.yml"
	case "go-sast-codeql":
		return filename == "go-sast-codeql.yaml" || filename == "go-sast-codeql.yml"
	case "ts-ci":
		return filename == "ts-ci.yaml" || filename == "ts-ci.yml"
	case "ts-lint":
		return filename == "ts-lint.yaml" || filename == "ts-lint.yml"
	default:
		return false
	}
}

// usesReusableWorkflow checks if a workflow uses the specified reusable workflow.
func (m *WorkflowMatcher) usesReusableWorkflow(wf model.Workflow, expectedRef string) bool {
	if wf.Content == "" {
		return false
	}

	// Parse workflow content
	var workflow map[string]any
	if err := yaml.Unmarshal([]byte(wf.Content), &workflow); err != nil {
		return false
	}

	// Check jobs for reusable workflow usage
	jobs, ok := workflow["jobs"].(map[string]any)
	if !ok {
		return false
	}

	for _, job := range jobs {
		jobMap, ok := job.(map[string]any)
		if !ok {
			continue
		}
		if uses, ok := jobMap["uses"].(string); ok {
			// Normalize the reference for comparison
			if m.normalizeRef(uses) == m.normalizeRef(expectedRef) {
				return true
			}
		}
	}

	return false
}

// normalizeRef normalizes a workflow reference for comparison.
func (m *WorkflowMatcher) normalizeRef(ref string) string {
	// Remove leading/trailing whitespace
	ref = strings.TrimSpace(ref)
	// Normalize branch reference (main, master, or tag)
	parts := strings.SplitN(ref, "@", 2)
	if len(parts) == 2 {
		return parts[0] + "@" + parts[1]
	}
	return ref
}

// isEquivalentWorkflow checks if a workflow is functionally equivalent to the rule type.
func (m *WorkflowMatcher) isEquivalentWorkflow(wf model.Workflow, ruleType string) bool {
	// Match by workflow name or file name
	lowerName := strings.ToLower(wf.Name)
	lowerPath := strings.ToLower(wf.Path)

	switch ruleType {
	case "go-ci":
		return m.matchesGoCI(wf, lowerName, lowerPath)
	case "go-lint":
		return m.matchesGoLint(wf, lowerName, lowerPath)
	case "go-sast-codeql":
		return m.matchesGoCodeQL(wf, lowerName, lowerPath)
	case "ts-ci":
		return m.matchesTSCI(wf, lowerName, lowerPath)
	case "ts-lint":
		return m.matchesTSLint(wf, lowerName, lowerPath)
	}

	return false
}

// matchesGoCI checks if a workflow matches Go CI patterns.
func (m *WorkflowMatcher) matchesGoCI(wf model.Workflow, name, path string) bool {
	// Check name/path patterns
	patterns := []string{"go-ci", "go_ci", "goci", "ci", "build", "test"}
	for _, p := range patterns {
		if strings.Contains(name, p) || strings.Contains(path, p) {
			// Verify it actually contains Go-related steps
			if m.containsGoSteps(wf.Content) {
				return true
			}
		}
	}
	return false
}

// matchesGoLint checks if a workflow matches Go lint patterns.
func (m *WorkflowMatcher) matchesGoLint(wf model.Workflow, name, path string) bool {
	patterns := []string{"go-lint", "go_lint", "golint", "lint"}
	for _, p := range patterns {
		if strings.Contains(name, p) || strings.Contains(path, p) {
			if m.containsGolangciLint(wf.Content) {
				return true
			}
		}
	}
	return false
}

// matchesGoCodeQL checks if a workflow matches Go CodeQL patterns.
func (m *WorkflowMatcher) matchesGoCodeQL(wf model.Workflow, name, path string) bool {
	patterns := []string{"codeql", "code-ql", "sast", "security"}
	for _, p := range patterns {
		if strings.Contains(name, p) || strings.Contains(path, p) {
			if m.containsCodeQL(wf.Content) {
				return true
			}
		}
	}
	return false
}

// matchesTSCI checks if a workflow matches TypeScript CI patterns.
func (m *WorkflowMatcher) matchesTSCI(wf model.Workflow, name, path string) bool {
	patterns := []string{"ts-ci", "ts_ci", "typescript-ci", "node-ci", "ci", "build", "test"}
	for _, p := range patterns {
		if strings.Contains(name, p) || strings.Contains(path, p) {
			if m.containsNodeSteps(wf.Content) {
				return true
			}
		}
	}
	return false
}

// matchesTSLint checks if a workflow matches TypeScript lint patterns.
func (m *WorkflowMatcher) matchesTSLint(wf model.Workflow, name, path string) bool {
	patterns := []string{"ts-lint", "ts_lint", "eslint", "lint"}
	for _, p := range patterns {
		if strings.Contains(name, p) || strings.Contains(path, p) {
			if m.containsESLint(wf.Content) {
				return true
			}
		}
	}
	return false
}

var (
	goStepsRe      = regexp.MustCompile(`(?i)(go\s+(build|test|mod)|actions/setup-go)`)
	golangciLintRe = regexp.MustCompile(`(?i)(golangci-lint|golangci/golangci-lint-action)`)
	codeqlRe       = regexp.MustCompile(`(?i)(codeql|github/codeql-action)`)
	nodeStepsRe    = regexp.MustCompile(`(?i)(npm\s+(run|test|build)|yarn\s+(test|build)|actions/setup-node)`)
	eslintRe       = regexp.MustCompile(`(?i)(eslint|npm\s+run\s+lint|yarn\s+lint)`)
)

func (m *WorkflowMatcher) containsGoSteps(content string) bool {
	return goStepsRe.MatchString(content)
}

func (m *WorkflowMatcher) containsGolangciLint(content string) bool {
	return golangciLintRe.MatchString(content)
}

func (m *WorkflowMatcher) containsCodeQL(content string) bool {
	return codeqlRe.MatchString(content)
}

func (m *WorkflowMatcher) containsNodeSteps(content string) bool {
	return nodeStepsRe.MatchString(content)
}

func (m *WorkflowMatcher) containsESLint(content string) bool {
	return eslintRe.MatchString(content)
}

// DetectWorkflowLanguage attempts to detect the primary language a workflow targets.
func DetectWorkflowLanguage(wf model.Workflow) string {
	content := wf.Content
	if content == "" {
		return ""
	}

	// Check for Go indicators
	if goStepsRe.MatchString(content) || golangciLintRe.MatchString(content) {
		return "Go"
	}

	// Check for TypeScript/Node indicators
	if nodeStepsRe.MatchString(content) || eslintRe.MatchString(content) {
		return "TypeScript"
	}

	return ""
}
