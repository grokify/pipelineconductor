// Package collector provides interfaces and implementations for collecting
// repository and workflow data from various sources.
package collector

import (
	"context"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

// Collector defines the interface for collecting repository data.
type Collector interface {
	// ListRepos returns repositories matching the filter criteria.
	ListRepos(ctx context.Context, orgs []string, filter model.RepoFilter) ([]model.Repo, error)

	// ListUserRepos returns public repositories for the specified users.
	ListUserRepos(ctx context.Context, users []string, filter model.RepoFilter) ([]model.Repo, error)

	// ListReposMultiSource returns repositories from both orgs and users.
	ListReposMultiSource(ctx context.Context, orgs, users []string, filter model.RepoFilter) ([]model.Repo, error)

	// GetWorkflows returns CI/CD workflows for a repository.
	GetWorkflows(ctx context.Context, repo model.Repo) ([]model.Workflow, error)

	// GetBranchProtection returns branch protection settings.
	GetBranchProtection(ctx context.Context, repo model.Repo, branch string) (*model.BranchProtection, error)

	// GetLatestWorkflowRun returns the most recent workflow run.
	GetLatestWorkflowRun(ctx context.Context, repo model.Repo, workflowID int64) (*model.WorkflowRun, error)

	// GetFileContent returns the content of a file from a repository.
	GetFileContent(ctx context.Context, repo model.Repo, path string) (string, error)

	// GetLanguages returns the languages used in a repository.
	GetLanguages(ctx context.Context, repo model.Repo) ([]string, error)
}

// ListOptions configures listing behavior.
type ListOptions struct {
	PerPage int
	Page    int
}

// DefaultListOptions returns sensible defaults.
func DefaultListOptions() ListOptions {
	return ListOptions{
		PerPage: 100,
		Page:    1,
	}
}
