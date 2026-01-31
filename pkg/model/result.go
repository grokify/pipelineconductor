package model

import "time"

// ComplianceResult is the result of a compliance scan.
type ComplianceResult struct {
	Timestamp      time.Time    `json:"timestamp"`
	Summary        ScanSummary  `json:"summary"`
	Repos          []RepoResult `json:"repos"`
	ScanDurationMs int64        `json:"scanDurationMs"`
	Config         ScanConfig   `json:"config"`
}

// ScanSummary provides aggregate statistics for a scan.
type ScanSummary struct {
	TotalRepos     int     `json:"total"`
	CompliantRepos int     `json:"compliant"`
	NonCompliant   int     `json:"nonCompliant"`
	Skipped        int     `json:"skipped"`
	Errors         int     `json:"errors"`
	ComplianceRate float64 `json:"complianceRate"`
}

// RepoResult is the compliance result for a single repository.
type RepoResult struct {
	Repo       Repo        `json:"repo"`
	Compliant  bool        `json:"compliant"`
	Violations []Violation `json:"violations,omitempty"`
	Warnings   []Warning   `json:"warnings,omitempty"`
	Skipped    bool        `json:"skipped"`
	SkipReason string      `json:"skipReason,omitempty"`
	Error      string      `json:"error,omitempty"`
	ScanTimeMs int64       `json:"scanTimeMs"`
}

// Violation represents a policy violation.
type Violation struct {
	Policy      string   `json:"policy"`
	Rule        string   `json:"rule"`
	Message     string   `json:"message"`
	Severity    Severity `json:"severity"`
	Remediation string   `json:"remediation,omitempty"`
	File        string   `json:"file,omitempty"`
	Line        int      `json:"line,omitempty"`
}

// Warning represents a non-blocking issue.
type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

// Severity indicates the severity level of a violation.
type Severity string

// Severity levels for policy violations.
const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// ScanConfig captures the configuration used for a scan.
type ScanConfig struct {
	Orgs       []string   `json:"orgs"`
	PolicyRepo string     `json:"policyRepo"`
	PolicyRef  string     `json:"policyRef"`
	Profile    string     `json:"profile"`
	Filter     RepoFilter `json:"filter"`
}

// RemediationPlan describes planned changes to fix violations.
type RemediationPlan struct {
	Repo    Repo    `json:"repo"`
	Patches []Patch `json:"patches"`
	PRTitle string  `json:"prTitle"`
	PRBody  string  `json:"prBody"`
	DryRun  bool    `json:"dryRun"`
}

// Patch represents a file change.
type Patch struct {
	Path      string `json:"path"`
	Operation string `json:"operation"` // "create", "update", "delete"
	Content   string `json:"content,omitempty"`
	Diff      string `json:"diff,omitempty"`
}

// RemediationResult is the result of applying a remediation.
type RemediationResult struct {
	Repo     Repo   `json:"repo"`
	Success  bool   `json:"success"`
	PRURL    string `json:"prUrl,omitempty"`
	PRNumber int    `json:"prNumber,omitempty"`
	Error    string `json:"error,omitempty"`
	DryRun   bool   `json:"dryRun"`
}

// IsCompliant returns true if there are no violations.
func (r *RepoResult) IsCompliant() bool {
	return len(r.Violations) == 0 && r.Error == ""
}

// ViolationCount returns the total number of violations.
func (c *ComplianceResult) ViolationCount() int {
	count := 0
	for _, r := range c.Repos {
		count += len(r.Violations)
	}
	return count
}
