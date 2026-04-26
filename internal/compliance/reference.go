package compliance

import (
	"context"
	"fmt"
	"strings"

	"github.com/plexusone/pipelineconductor/internal/collector"
	"github.com/plexusone/pipelineconductor/pkg/model"
)

// ReferenceRepo represents a reference workflow repository.
type ReferenceRepo struct {
	Owner  string
	Name   string
	Branch string
}

// ParseReferenceRepo parses a reference repo string (e.g., "grokify/.github").
func ParseReferenceRepo(ref string) (*ReferenceRepo, error) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid reference repo format: %s (expected owner/repo)", ref)
	}
	return &ReferenceRepo{
		Owner:  parts[0],
		Name:   parts[1],
		Branch: "main",
	}, nil
}

// FullName returns the full name of the reference repo.
func (r *ReferenceRepo) FullName() string {
	return r.Owner + "/" + r.Name
}

// WorkflowRef returns the reusable workflow reference for a workflow path.
// Format: owner/repo/.github/workflows/workflow.yaml@branch
func (r *ReferenceRepo) WorkflowRef(workflowPath string) string {
	return fmt.Sprintf("%s/%s/%s@%s", r.Owner, r.Name, workflowPath, r.Branch)
}

// ReferenceWorkflows holds the fetched reference workflows.
type ReferenceWorkflows struct {
	Repo      *ReferenceRepo
	Workflows map[string]string // type -> content
}

// FetchReferenceWorkflows fetches reference workflows from the specified repository.
func FetchReferenceWorkflows(ctx context.Context, c collector.Collector, refRepo *ReferenceRepo, rules []WorkflowRule) (*ReferenceWorkflows, error) {
	ref := &ReferenceWorkflows{
		Repo:      refRepo,
		Workflows: make(map[string]string),
	}

	repo := model.Repo{
		Owner:    refRepo.Owner,
		Name:     refRepo.Name,
		FullName: refRepo.FullName(),
	}

	for _, rule := range rules {
		content, err := c.GetFileContent(ctx, repo, rule.Path)
		if err != nil {
			// Workflow may not exist in reference repo - that's OK
			continue
		}
		ref.Workflows[rule.Type] = content
	}

	return ref, nil
}

// HasWorkflow returns true if the reference has the specified workflow type.
func (r *ReferenceWorkflows) HasWorkflow(workflowType string) bool {
	_, ok := r.Workflows[workflowType]
	return ok
}

// GetExpectedRef returns the expected reusable workflow reference for a type.
func (r *ReferenceWorkflows) GetExpectedRef(rule WorkflowRule) string {
	return r.Repo.WorkflowRef(rule.Path)
}
