// Package policy provides Cedar policy evaluation for CI/CD compliance.
package policy

import (
	"fmt"

	"github.com/cedar-policy/cedar-go"
	"github.com/grokify/pipelineconductor/pkg/model"
)

// Action constants for CI/CD policy evaluation.
const (
	ActionBuild   = "build"
	ActionTest    = "test"
	ActionLint    = "lint"
	ActionMerge   = "merge"
	ActionDeploy  = "deploy"
	ActionRelease = "release"
)

// Engine evaluates Cedar policies against repository contexts.
type Engine struct {
	policySet *cedar.PolicySet
	entities  cedar.EntityMap
}

// NewEngine creates a new policy evaluation engine.
func NewEngine() *Engine {
	return &Engine{
		policySet: cedar.NewPolicySet(),
		entities:  make(cedar.EntityMap),
	}
}

// AddPolicy adds a policy to the engine.
func (e *Engine) AddPolicy(id string, policyText []byte) error {
	var policy cedar.Policy
	if err := policy.UnmarshalCedar(policyText); err != nil {
		return fmt.Errorf("parsing policy %s: %w", id, err)
	}
	e.policySet.Add(cedar.PolicyID(id), &policy)
	return nil
}

// AddPolicyFromFile loads and adds a policy from Cedar text.
func (e *Engine) AddPolicyFromFile(id string, content []byte) error {
	return e.AddPolicy(id, content)
}

// Evaluate evaluates policies for a repository and action.
func (e *Engine) Evaluate(ctx *model.PolicyContext, action string) *EvaluationResult {
	// Build Cedar request
	req := e.buildRequest(ctx, action)

	// Evaluate against policy set
	decision, diagnostic := cedar.Authorize(e.policySet, e.entities, req)

	result := &EvaluationResult{
		Allowed:    decision == cedar.Allow,
		Action:     action,
		RepoName:   ctx.Repo.FullName,
		Diagnostic: diagnostic,
	}

	// Extract reasons from diagnostic
	for _, reason := range diagnostic.Reasons {
		result.Reasons = append(result.Reasons, string(reason.PolicyID))
	}

	for _, err := range diagnostic.Errors {
		result.Errors = append(result.Errors, err.String())
	}

	return result
}

// EvaluateAll evaluates all standard CI/CD actions for a repository.
func (e *Engine) EvaluateAll(ctx *model.PolicyContext) []*EvaluationResult {
	actions := []string{ActionBuild, ActionTest, ActionLint, ActionMerge}
	var results []*EvaluationResult
	for _, action := range actions {
		results = append(results, e.Evaluate(ctx, action))
	}
	return results
}

// buildRequest constructs a Cedar request from a policy context.
func (e *Engine) buildRequest(ctx *model.PolicyContext, action string) cedar.Request {
	// Principal: the CI system or actor
	principal := cedar.NewEntityUID(cedar.EntityType("CISystem"), cedar.String("pipelineconductor"))

	// Action: the CI/CD action being evaluated
	actionUID := cedar.NewEntityUID(cedar.EntityType("Action"), cedar.String(action))

	// Resource: the repository
	resource := cedar.NewEntityUID(cedar.EntityType("Repository"), cedar.String(ctx.Repo.FullName))

	// Context: all the policy-relevant data
	context := cedar.NewRecord(cedar.RecordMap{
		// Repository info
		"repoName":     cedar.String(ctx.Repo.Name),
		"repoOrg":      cedar.String(ctx.Repo.Org),
		"repoFullName": cedar.String(ctx.Repo.FullName),
		"archived":     cedar.Boolean(ctx.Repo.Archived),
		"fork":         cedar.Boolean(ctx.Repo.Fork),
		"languages":    stringSliceToSet(ctx.Repo.Language),
		"topics":       stringSliceToSet(ctx.Repo.Topics),

		// CI configuration
		"hasWorkflow":          cedar.Boolean(ctx.CI.HasWorkflow),
		"usesReusableWorkflow": cedar.Boolean(ctx.CI.UsesReusableWorkflow),
		"reusableWorkflowRef":  cedar.String(ctx.CI.ReusableWorkflowRef),
		"lastRunPassed":        cedar.Boolean(ctx.CI.LastRunPassed),
		"requiredChecks":       stringSliceToSet(ctx.CI.RequiredChecks),
		"osMatrix":             stringSliceToSet(ctx.CI.OSMatrix),

		// Go-specific
		"goVersions": stringSliceToSet(ctx.Go.Versions),
		"goProfile":  cedar.String(ctx.Go.Profile),
		"hasGoMod":   cedar.Boolean(ctx.Go.HasGoMod),

		// Dependencies
		"hasRenovate":          cedar.Boolean(ctx.Dependencies.HasRenovate),
		"hasDependabot":        cedar.Boolean(ctx.Dependencies.HasDependabot),
		"oldestDependencyDays": cedar.Long(int64(ctx.Dependencies.OldestDependencyDays)),
		"hasVulnerabilities":   cedar.Boolean(ctx.Dependencies.HasVulnerabilities),
		"vulnerabilityCount":   cedar.Long(int64(ctx.Dependencies.VulnerabilityCount)),

		// Branch protection
		"branchProtectionEnabled":       cedar.Boolean(ctx.BranchProtection.Enabled),
		"requireReviews":                cedar.Boolean(ctx.BranchProtection.RequireReviews),
		"requireStatusChecks":           cedar.Boolean(ctx.BranchProtection.RequireStatusChecks),
		"branchProtectionEnforceAdmins": cedar.Boolean(ctx.BranchProtection.EnforceAdmins),
	})

	return cedar.Request{
		Principal: principal,
		Action:    actionUID,
		Resource:  resource,
		Context:   context,
	}
}

// stringSliceToSet converts a string slice to a Cedar Set.
func stringSliceToSet(slice []string) cedar.Set {
	values := make([]cedar.Value, len(slice))
	for i, s := range slice {
		values[i] = cedar.String(s)
	}
	return cedar.NewSet(values...)
}

// EvaluationResult contains the result of a policy evaluation.
type EvaluationResult struct {
	Allowed    bool
	Action     string
	RepoName   string
	Reasons    []string
	Errors     []string
	Diagnostic cedar.Diagnostic
}

// ToViolation converts a denied evaluation to a model.Violation.
func (r *EvaluationResult) ToViolation() *model.Violation {
	if r.Allowed {
		return nil
	}

	severity := model.SeverityMedium
	if r.Action == ActionMerge || r.Action == ActionDeploy {
		severity = model.SeverityHigh
	}

	message := fmt.Sprintf("Policy denied %s action", r.Action)
	if len(r.Reasons) > 0 {
		message = fmt.Sprintf("Policy denied %s action (policies: %v)", r.Action, r.Reasons)
	}

	return &model.Violation{
		Policy:   fmt.Sprintf("cedar/%s", r.Action),
		Rule:     r.Action,
		Message:  message,
		Severity: severity,
	}
}
