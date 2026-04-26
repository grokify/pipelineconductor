package compliance

import (
	"context"
	"fmt"
	"time"

	"github.com/plexusone/pipelineconductor/internal/collector"
	"github.com/plexusone/pipelineconductor/pkg/model"
)

// Checker performs workflow compliance checks.
type Checker struct {
	Collector collector.Collector
	RefRepo   *ReferenceRepo
	Strict    bool
	Verbose   bool
}

// CheckerConfig configures the compliance checker.
type CheckerConfig struct {
	RefRepo   string
	RefBranch string
	Strict    bool
	Verbose   bool
}

// NewChecker creates a new compliance checker.
func NewChecker(c collector.Collector, cfg CheckerConfig) (*Checker, error) {
	refRepo, err := ParseReferenceRepo(cfg.RefRepo)
	if err != nil {
		return nil, err
	}
	if cfg.RefBranch != "" {
		refRepo.Branch = cfg.RefBranch
	}

	return &Checker{
		Collector: c,
		RefRepo:   refRepo,
		Strict:    cfg.Strict,
		Verbose:   cfg.Verbose,
	}, nil
}

// CheckRepos checks compliance for a list of repositories.
func (c *Checker) CheckRepos(ctx context.Context, repos []model.Repo, languages []string) (*model.CheckResult, error) {
	startTime := time.Now()

	// Get required workflow rules for the specified languages
	rules := GetRequiredWorkflows(languages)

	// Fetch reference workflows
	_, err := FetchReferenceWorkflows(ctx, c.Collector, c.RefRepo, rules)
	if err != nil {
		return nil, fmt.Errorf("fetching reference workflows: %w", err)
	}

	// Create matcher
	matcher := NewWorkflowMatcher(c.RefRepo, c.Strict)

	result := &model.CheckResult{
		SchemaVersion: "1.0.0",
		Timestamp:     time.Now().Format(time.RFC3339),
		Repos:         make([]model.RepoCheckResult, 0, len(repos)),
		Config: model.CheckConfig{
			RefRepo:   c.RefRepo.FullName(),
			RefBranch: c.RefRepo.Branch,
			Languages: languages,
			Strict:    c.Strict,
		},
	}

	// Check each repository
	for _, repo := range repos {
		repoResult := c.checkRepo(ctx, repo, rules, matcher)
		result.Repos = append(result.Repos, repoResult)
	}

	// Calculate summary
	result.Summary = c.calculateSummary(result.Repos, languages)
	result.ScanDurationMs = time.Since(startTime).Milliseconds()

	return result, nil
}

// checkRepo checks a single repository for compliance.
func (c *Checker) checkRepo(ctx context.Context, repo model.Repo, rules []WorkflowRule, matcher *WorkflowMatcher) model.RepoCheckResult {
	startTime := time.Now()

	result := model.RepoCheckResult{
		Owner:             repo.Owner,
		Name:              repo.Name,
		FullName:          repo.FullName,
		HTMLURL:           repo.HTMLURL,
		Languages:         repo.Languages,
		RequiredWorkflows: make([]model.WorkflowCheck, 0),
		ActualWorkflows:   make([]model.WorkflowInfo, 0),
		Missing:           make([]model.MissingWorkflow, 0),
	}

	// Get actual workflows
	workflows, err := c.Collector.GetWorkflows(ctx, repo)
	if err != nil {
		result.Error = fmt.Sprintf("failed to get workflows: %v", err)
		result.ComplianceLevel = model.ComplianceLevelNone
		result.ScanTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	// Convert workflows to WorkflowInfo
	for _, wf := range workflows {
		info := model.WorkflowInfo{
			Name:             wf.Name,
			Path:             wf.Path,
			UsesReusable:     wf.UsesReusableWorkflow,
			ReusableRefs:     extractReusableRefs(wf),
			DetectedLanguage: DetectWorkflowLanguage(wf),
		}
		result.ActualWorkflows = append(result.ActualWorkflows, info)
	}

	// Get relevant rules for this repo's languages
	repoRules := c.filterRulesForRepo(rules, repo.Languages)

	if len(repoRules) == 0 {
		// No rules apply to this repo
		result.Skipped = true
		result.SkipReason = "no applicable workflow rules for repo languages"
		result.ComplianceLevel = model.ComplianceLevelFull
		result.Compliant = true
		result.ScanTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	// Check each required workflow
	exactMatches := 0
	equivalentMatches := 0
	missingCount := 0

	for _, rule := range repoRules {
		matchResult := matcher.MatchWorkflow(workflows, rule)

		check := model.WorkflowCheck{
			WorkflowType:     rule.Type,
			Required:         true,
			Present:          matchResult.MatchType != model.MatchTypeNone,
			UsesReusable:     matchResult.UsesReusable,
			ReusableRef:      matchResult.ReusableRef,
			ExpectedRef:      c.RefRepo.WorkflowRef(rule.Path),
			MatchType:        matchResult.MatchType,
			ActualWorkflow:   matchResult.ActualWorkflow,
			FilenameMismatch: matchResult.FilenameMismatch,
			ExpectedFilename: matchResult.ExpectedFilename,
			ActualFilename:   matchResult.ActualFilename,
		}
		result.RequiredWorkflows = append(result.RequiredWorkflows, check)

		switch matchResult.MatchType {
		case model.MatchTypeExact:
			exactMatches++
		case model.MatchTypeEquivalent:
			equivalentMatches++
		case model.MatchTypeNone:
			missingCount++
			result.Missing = append(result.Missing, model.MissingWorkflow{
				Language:     getLanguageForRule(rule),
				WorkflowType: rule.Type,
				RefPath:      rule.Path,
				Severity:     rule.Severity,
				Description:  rule.Description,
			})
		}
	}

	// Determine compliance level
	totalRules := len(repoRules)
	if missingCount == 0 && equivalentMatches == 0 {
		result.ComplianceLevel = model.ComplianceLevelFull
		result.Compliant = true
	} else if missingCount == 0 {
		result.ComplianceLevel = model.ComplianceLevelPartial
		result.Compliant = false
	} else if exactMatches+equivalentMatches > 0 {
		result.ComplianceLevel = model.ComplianceLevelPartial
		result.Compliant = false
	} else if missingCount == totalRules {
		result.ComplianceLevel = model.ComplianceLevelNone
		result.Compliant = false
	} else {
		result.ComplianceLevel = model.ComplianceLevelNone
		result.Compliant = false
	}

	result.ScanTimeMs = time.Since(startTime).Milliseconds()
	return result
}

// filterRulesForRepo returns rules that apply to the repo's languages.
func (c *Checker) filterRulesForRepo(allRules []WorkflowRule, repoLanguages []string) []WorkflowRule {
	// Create a set of repo languages for quick lookup
	langSet := make(map[string]bool)
	for _, lang := range repoLanguages {
		langSet[lang] = true
	}

	// Get rules that apply to repo languages
	var applicableRules []WorkflowRule
	for _, rule := range allRules {
		ruleLang := getLanguageForRule(rule)
		if langSet[ruleLang] {
			applicableRules = append(applicableRules, rule)
		}
	}

	return applicableRules
}

// getLanguageForRule determines which language a rule belongs to.
func getLanguageForRule(rule WorkflowRule) string {
	for lang, rules := range LanguageRules {
		for _, r := range rules {
			if r.Type == rule.Type {
				return lang
			}
		}
	}
	return ""
}

// extractReusableRefs extracts reusable workflow references from a workflow.
func extractReusableRefs(wf model.Workflow) []string {
	var refs []string
	for _, ref := range wf.ReusableWorkflowRefs {
		refs = append(refs, ref.FullRef)
	}
	return refs
}

// calculateSummary calculates the compliance summary from results.
func (c *Checker) calculateSummary(repos []model.RepoCheckResult, languages []string) model.CheckSummary {
	summary := model.CheckSummary{
		TotalRepos: len(repos),
		ByLanguage: make([]model.LanguageComplianceStats, 0),
	}

	// Track per-language stats
	langStats := make(map[string]*model.LanguageComplianceStats)
	for _, lang := range languages {
		langStats[lang] = &model.LanguageComplianceStats{
			Language: lang,
		}
	}

	for _, repo := range repos {
		if repo.Skipped {
			summary.Skipped++
			continue
		}
		if repo.Error != "" {
			summary.Errors++
			continue
		}

		switch repo.ComplianceLevel {
		case model.ComplianceLevelFull:
			summary.CompliantRepos++
		case model.ComplianceLevelPartial:
			summary.PartialRepos++
		case model.ComplianceLevelNone:
			summary.NonCompliant++
		}

		// Update per-language stats
		for _, lang := range repo.Languages {
			if stats, ok := langStats[lang]; ok {
				stats.TotalRepos++
				if repo.ComplianceLevel == model.ComplianceLevelFull {
					stats.CompliantRepos++
				}
			}
		}
	}

	// Calculate compliance rates
	countedRepos := summary.TotalRepos - summary.Skipped - summary.Errors
	if countedRepos > 0 {
		summary.ComplianceRate = float64(summary.CompliantRepos) / float64(countedRepos) * 100
	}

	// Build language stats array
	for _, lang := range languages {
		stats := langStats[lang]
		if stats.TotalRepos > 0 {
			stats.ComplianceRate = float64(stats.CompliantRepos) / float64(stats.TotalRepos) * 100
		}
		summary.ByLanguage = append(summary.ByLanguage, *stats)
	}

	return summary
}
