package policy

import (
	"github.com/grokify/pipelineconductor/pkg/model"
)

// ContextBuilder builds a PolicyContext from repository data.
type ContextBuilder struct {
	profile *model.Profile
}

// NewContextBuilder creates a new context builder with an optional profile.
func NewContextBuilder(profile *model.Profile) *ContextBuilder {
	return &ContextBuilder{profile: profile}
}

// Build creates a PolicyContext from a repository and its associated data.
func (b *ContextBuilder) Build(repo model.Repo, workflows []model.Workflow, bp *model.BranchProtection) *model.PolicyContext {
	ctx := &model.PolicyContext{}

	// Repo context
	ctx.Repo = model.RepoContext{
		Name:     repo.Name,
		Org:      repo.Owner,
		FullName: repo.FullName,
		Language: repo.Languages,
		Topics:   repo.Topics,
		Archived: repo.Archived,
		Fork:     repo.Fork,
	}

	// CI context
	ctx.CI = b.buildCIContext(workflows)

	// Go context (if applicable)
	ctx.Go = b.buildGoContext(repo, workflows)

	// Branch protection context
	if bp != nil {
		ctx.BranchProtection = model.BranchProtectionContext{
			Enabled:             bp.Enabled,
			RequireReviews:      bp.RequireReviews,
			RequireStatusChecks: bp.RequireStatusChecks,
			EnforceAdmins:       bp.EnforceAdmins,
		}
	}

	return ctx
}

func (b *ContextBuilder) buildCIContext(workflows []model.Workflow) model.CIContext {
	ctx := model.CIContext{
		HasWorkflow: len(workflows) > 0,
	}

	var osMatrix []string
	osMatrixSet := make(map[string]bool)

	for _, wf := range workflows {
		// Check for reusable workflow usage
		if len(wf.ReusableWorkflowRefs) > 0 {
			ctx.UsesReusableWorkflow = true
			if len(wf.ReusableWorkflowRefs) > 0 {
				ctx.ReusableWorkflowRef = wf.ReusableWorkflowRefs[0].FullRef
			}
		}

		// Collect OS matrix from jobs
		for _, job := range wf.Jobs {
			for _, os := range job.RunsOn {
				if !osMatrixSet[os] {
					osMatrixSet[os] = true
					osMatrix = append(osMatrix, os)
				}
			}

			// Check job matrix
			if job.Matrix != nil {
				for _, os := range job.Matrix.OS {
					if !osMatrixSet[os] {
						osMatrixSet[os] = true
						osMatrix = append(osMatrix, os)
					}
				}
			}
		}
	}

	ctx.OSMatrix = osMatrix
	return ctx
}

func (b *ContextBuilder) buildGoContext(repo model.Repo, workflows []model.Workflow) model.GoContext {
	ctx := model.GoContext{}

	// Check if this is a Go project
	isGo := false
	for _, lang := range repo.Languages {
		if lang == "Go" {
			isGo = true
			break
		}
	}

	if !isGo {
		return ctx
	}

	// Set profile if provided
	if b.profile != nil {
		ctx.Profile = b.profile.Name
		ctx.Versions = b.profile.Go.Versions
	}

	// Extract Go versions from workflow matrix
	var versions []string
	versionSet := make(map[string]bool)

	for _, wf := range workflows {
		for _, job := range wf.Jobs {
			if job.Matrix != nil {
				for _, v := range job.Matrix.GoVersion {
					if !versionSet[v] {
						versionSet[v] = true
						versions = append(versions, v)
					}
				}
			}
		}
	}

	if len(versions) > 0 {
		ctx.Versions = versions
	}

	return ctx
}

// BuildFromRepoResult creates a PolicyContext from a compliance scan result.
func (b *ContextBuilder) BuildFromRepoResult(result *model.RepoResult, workflows []model.Workflow, bp *model.BranchProtection) *model.PolicyContext {
	return b.Build(result.Repo, workflows, bp)
}

// WithProfile sets the profile for context building.
func (b *ContextBuilder) WithProfile(profile *model.Profile) *ContextBuilder {
	b.profile = profile
	return b
}
