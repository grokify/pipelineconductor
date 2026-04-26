package model

// CheckResult is the top-level result of a workflow compliance check.
type CheckResult struct {
	SchemaVersion  string            `json:"schemaVersion"`
	Timestamp      string            `json:"timestamp"`
	Summary        CheckSummary      `json:"summary"`
	Repos          []RepoCheckResult `json:"repos"`
	ScanDurationMs int64             `json:"scanDurationMs"`
	Config         CheckConfig       `json:"config"`
}

// CheckSummary provides aggregate statistics for a compliance check.
type CheckSummary struct {
	TotalRepos     int                       `json:"totalRepos"`
	CompliantRepos int                       `json:"compliantRepos"`
	PartialRepos   int                       `json:"partialRepos"`
	NonCompliant   int                       `json:"nonCompliant"`
	Skipped        int                       `json:"skipped"`
	Errors         int                       `json:"errors"`
	ComplianceRate float64                   `json:"complianceRate"`
	ByLanguage     []LanguageComplianceStats `json:"byLanguage"`
}

// LanguageComplianceStats provides per-language compliance breakdown.
type LanguageComplianceStats struct {
	Language       string  `json:"language"`
	TotalRepos     int     `json:"totalRepos"`
	CompliantRepos int     `json:"compliantRepos"`
	ComplianceRate float64 `json:"complianceRate"`
}

// RepoCheckResult is the compliance result for a single repository.
type RepoCheckResult struct {
	Owner             string            `json:"owner"`
	Name              string            `json:"name"`
	FullName          string            `json:"fullName"`
	HTMLURL           string            `json:"htmlUrl"`
	Languages         []string          `json:"languages"`
	Compliant         bool              `json:"compliant"`
	ComplianceLevel   string            `json:"complianceLevel"`
	RequiredWorkflows []WorkflowCheck   `json:"requiredWorkflows"`
	ActualWorkflows   []WorkflowInfo    `json:"actualWorkflows"`
	Missing           []MissingWorkflow `json:"missing"`
	Skipped           bool              `json:"skipped"`
	SkipReason        string            `json:"skipReason"`
	Error             string            `json:"error"`
	ScanTimeMs        int64             `json:"scanTimeMs"`
}

// WorkflowCheck is the result of checking a single required workflow.
type WorkflowCheck struct {
	WorkflowType     string `json:"workflowType"`
	Required         bool   `json:"required"`
	Present          bool   `json:"present"`
	UsesReusable     bool   `json:"usesReusable"`
	ReusableRef      string `json:"reusableRef"`
	ExpectedRef      string `json:"expectedRef"`
	MatchType        string `json:"matchType"`
	ActualWorkflow   string `json:"actualWorkflow"`
	FilenameMismatch bool   `json:"filenameMismatch,omitempty"`
	ExpectedFilename string `json:"expectedFilename,omitempty"`
	ActualFilename   string `json:"actualFilename,omitempty"`
}

// WorkflowInfo provides information about an existing workflow in a repository.
type WorkflowInfo struct {
	Name             string   `json:"name"`
	Path             string   `json:"path"`
	UsesReusable     bool     `json:"usesReusable"`
	ReusableRefs     []string `json:"reusableRefs"`
	DetectedLanguage string   `json:"detectedLanguage"`
}

// MissingWorkflow represents a required workflow that is not present.
type MissingWorkflow struct {
	Language     string `json:"language"`
	WorkflowType string `json:"workflowType"`
	RefPath      string `json:"refPath"`
	Severity     string `json:"severity"`
	Description  string `json:"description"`
}

// CheckConfig captures the configuration used for a compliance check.
type CheckConfig struct {
	Orgs      []string `json:"orgs"`
	Users     []string `json:"users"`
	RefRepo   string   `json:"refRepo"`
	RefBranch string   `json:"refBranch"`
	Languages []string `json:"languages"`
	Strict    bool     `json:"strict"`
}

// ComplianceLevel constants.
const (
	ComplianceLevelFull    = "full"
	ComplianceLevelPartial = "partial"
	ComplianceLevelNone    = "none"
)

// MatchType constants.
const (
	MatchTypeExact      = "exact"
	MatchTypeEquivalent = "equivalent"
	MatchTypePartial    = "partial"
	MatchTypeNone       = "none"
)

// SeverityLevel constants for missing workflows.
const (
	SeverityLevelHigh   = "high"
	SeverityLevelMedium = "medium"
	SeverityLevelLow    = "low"
)
