// Package compliance provides workflow compliance checking functionality.
package compliance

import "github.com/plexusone/pipelineconductor/pkg/model"

// WorkflowRule defines a required workflow for a language.
type WorkflowRule struct {
	Type        string
	Path        string
	Severity    string
	Description string
}

// LanguageRules maps languages to their required workflows.
var LanguageRules = map[string][]WorkflowRule{
	"Go": {
		{
			Type:        "go-ci",
			Path:        ".github/workflows/go-ci.yaml",
			Severity:    model.SeverityLevelHigh,
			Description: "Go CI pipeline with build, test, and coverage",
		},
		{
			Type:        "go-lint",
			Path:        ".github/workflows/go-lint.yaml",
			Severity:    model.SeverityLevelMedium,
			Description: "Go linting with golangci-lint",
		},
		{
			Type:        "go-sast-codeql",
			Path:        ".github/workflows/go-sast-codeql.yaml",
			Severity:    model.SeverityLevelLow,
			Description: "Go static analysis with CodeQL",
		},
	},
	"TypeScript": {
		{
			Type:        "ts-ci",
			Path:        ".github/workflows/ts-ci.yaml",
			Severity:    model.SeverityLevelHigh,
			Description: "TypeScript CI pipeline with build and test",
		},
		{
			Type:        "ts-lint",
			Path:        ".github/workflows/ts-lint.yaml",
			Severity:    model.SeverityLevelMedium,
			Description: "TypeScript/JavaScript linting with ESLint",
		},
	},
	"Crystal": {
		// No reference workflows yet
	},
}

// GetRequiredWorkflows returns the required workflows for the given languages.
func GetRequiredWorkflows(languages []string) []WorkflowRule {
	seen := make(map[string]bool)
	var rules []WorkflowRule

	for _, lang := range languages {
		if langRules, ok := LanguageRules[lang]; ok {
			for _, rule := range langRules {
				if !seen[rule.Type] {
					seen[rule.Type] = true
					rules = append(rules, rule)
				}
			}
		}
	}

	return rules
}

// SupportedLanguages returns the list of languages with defined workflow rules.
func SupportedLanguages() []string {
	return []string{"Go", "TypeScript", "Crystal"}
}

// IsLanguageSupported returns true if the language has workflow rules defined.
func IsLanguageSupported(language string) bool {
	_, ok := LanguageRules[language]
	return ok
}
