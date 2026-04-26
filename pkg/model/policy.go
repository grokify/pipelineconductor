package model

// PolicyContext is the canonical context passed to Cedar for policy evaluation.
type PolicyContext struct {
	Repo             RepoContext             `json:"repo"`
	CI               CIContext               `json:"ci"`
	Go               GoContext               `json:"go"`
	Dependencies     DependenciesContext     `json:"dependencies"`
	BranchProtection BranchProtectionContext `json:"branchProtection"`
	Compliance       ComplianceContext       `json:"compliance"`
}

// RepoContext contains repository information for policy evaluation.
type RepoContext struct {
	Name     string   `json:"name"`
	Org      string   `json:"org"`
	FullName string   `json:"fullName"`
	Language []string `json:"language"`
	Topics   []string `json:"topics"`
	Archived bool     `json:"archived"`
	Fork     bool     `json:"fork"`
}

// CIContext contains CI/CD workflow information for policy evaluation.
type CIContext struct {
	HasWorkflow          bool     `json:"hasWorkflow"`
	UsesReusableWorkflow bool     `json:"usesReusableWorkflow"`
	ReusableWorkflowRef  string   `json:"reusableWorkflowRef"`
	RequiredChecks       []string `json:"requiredChecks"`
	LastRunPassed        bool     `json:"lastRunPassed"`
	OSMatrix             []string `json:"osMatrix"`
}

// GoContext contains Go-specific information for policy evaluation.
type GoContext struct {
	Versions  []string `json:"versions"`
	Profile   string   `json:"profile"`
	HasGoMod  bool     `json:"hasGoMod"`
	GoModTidy bool     `json:"goModTidy"`
}

// DependenciesContext contains dependency information for policy evaluation.
type DependenciesContext struct {
	HasRenovate          bool `json:"hasRenovate"`
	HasDependabot        bool `json:"hasDependabot"`
	OldestDependencyDays int  `json:"oldestDependencyDays"`
	HasVulnerabilities   bool `json:"hasVulnerabilities"`
	VulnerabilityCount   int  `json:"vulnerabilityCount"`
}

// BranchProtectionContext contains branch protection info for policy evaluation.
type BranchProtectionContext struct {
	Enabled             bool `json:"enabled"`
	RequireReviews      bool `json:"requireReviews"`
	RequireStatusChecks bool `json:"requireStatusChecks"`
	EnforceAdmins       bool `json:"enforceAdmins"`
}

// ComplianceContext contains workflow compliance information for policy evaluation.
type ComplianceContext struct {
	// Level is the compliance level: "full", "partial", or "none"
	Level string `json:"level"`
	// Compliant indicates if the repo is fully compliant
	Compliant bool `json:"compliant"`
	// ComplianceRate is the percentage of required workflows that are present (0-100)
	ComplianceRate float64 `json:"complianceRate"`
	// MissingWorkflowCount is the number of required workflows that are missing
	MissingWorkflowCount int `json:"missingWorkflowCount"`
	// MissingWorkflows lists the types of missing workflows
	MissingWorkflows []string `json:"missingWorkflows"`
	// HasFilenameMismatch indicates if any workflows have non-standard filenames
	HasFilenameMismatch bool `json:"hasFilenameMismatch"`
	// UsesReusableWorkflows indicates if the repo uses any reusable workflows
	UsesReusableWorkflows bool `json:"usesReusableWorkflows"`
	// ExactMatchCount is the number of exact matches (using reusable workflows)
	ExactMatchCount int `json:"exactMatchCount"`
	// EquivalentMatchCount is the number of equivalent matches (same function, not reusable)
	EquivalentMatchCount int `json:"equivalentMatchCount"`
	// RefRepo is the reference repository for compliance checking
	RefRepo string `json:"refRepo"`
}

// Profile defines a named CI/CD configuration profile.
type Profile struct {
	Name        string        `json:"name" yaml:"name"`
	Description string        `json:"description" yaml:"description"`
	Go          ProfileGo     `json:"go" yaml:"go"`
	OS          []string      `json:"os" yaml:"os"`
	Checks      ProfileChecks `json:"checks" yaml:"checks"`
	Lint        ProfileLint   `json:"lint" yaml:"lint"`
	Test        ProfileTest   `json:"test" yaml:"test"`
}

// ProfileGo contains Go-specific profile settings.
type ProfileGo struct {
	Versions []string `json:"versions" yaml:"versions"`
}

// ProfileChecks contains required check settings.
type ProfileChecks struct {
	Required []string `json:"required" yaml:"required"`
}

// ProfileLint contains linting settings.
type ProfileLint struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Tool    string `json:"tool" yaml:"tool"`
}

// ProfileTest contains test settings.
type ProfileTest struct {
	Coverage bool `json:"coverage" yaml:"coverage"`
	Race     bool `json:"race" yaml:"race"`
}
