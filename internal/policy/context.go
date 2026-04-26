package policy

import (
	"slices"

	"github.com/plexusone/pipelineconductor/pkg/model"
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
	if !slices.Contains(repo.Languages, "Go") {
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

// BuildFromComplianceResult creates a PolicyContext from a compliance check result.
func (b *ContextBuilder) BuildFromComplianceResult(result model.RepoCheckResult, workflows []model.Workflow, refRepo string) *model.PolicyContext {
	// Build base context from repo
	repo := model.Repo{
		Owner:     result.Owner,
		Name:      result.Name,
		FullName:  result.FullName,
		Languages: result.Languages,
		HTMLURL:   result.HTMLURL,
	}

	ctx := b.Build(repo, workflows, nil)

	// Add compliance context
	ctx.Compliance = b.buildComplianceContext(result, refRepo)

	return ctx
}

// buildComplianceContext creates a ComplianceContext from a RepoCheckResult.
func (b *ContextBuilder) buildComplianceContext(result model.RepoCheckResult, refRepo string) model.ComplianceContext {
	compliance := model.ComplianceContext{
		Level:     result.ComplianceLevel,
		Compliant: result.Compliant,
		RefRepo:   refRepo,
	}

	// Count missing workflows
	compliance.MissingWorkflowCount = len(result.Missing)
	for _, m := range result.Missing {
		compliance.MissingWorkflows = append(compliance.MissingWorkflows, m.WorkflowType)
	}

	// Count match types and check for filename mismatches
	totalRequired := len(result.RequiredWorkflows)
	for _, wf := range result.RequiredWorkflows {
		switch wf.MatchType {
		case model.MatchTypeExact:
			compliance.ExactMatchCount++
			if wf.UsesReusable {
				compliance.UsesReusableWorkflows = true
			}
		case model.MatchTypeEquivalent:
			compliance.EquivalentMatchCount++
		}
		if wf.FilenameMismatch {
			compliance.HasFilenameMismatch = true
		}
	}

	// Calculate compliance rate
	if totalRequired > 0 {
		matchedCount := compliance.ExactMatchCount + compliance.EquivalentMatchCount
		compliance.ComplianceRate = float64(matchedCount) / float64(totalRequired) * 100
	}

	return compliance
}
