package model

import (
	"testing"
)

func TestRepoResultIsCompliant(t *testing.T) {
	tests := []struct {
		name   string
		result RepoResult
		want   bool
	}{
		{
			name: "compliant - no violations or errors",
			result: RepoResult{
				Repo:       Repo{FullName: "org/repo"},
				Violations: []Violation{},
				Error:      "",
			},
			want: true,
		},
		{
			name: "non-compliant - has violations",
			result: RepoResult{
				Repo: Repo{FullName: "org/repo"},
				Violations: []Violation{
					{
						Policy:   "ci/workflow-required",
						Rule:     "has-workflow",
						Message:  "No CI/CD workflow found",
						Severity: SeverityHigh,
					},
				},
			},
			want: false,
		},
		{
			name: "non-compliant - has error",
			result: RepoResult{
				Repo:  Repo{FullName: "org/repo"},
				Error: "failed to fetch repository",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.IsCompliant()
			if got != tt.want {
				t.Errorf("RepoResult.IsCompliant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplianceResultViolationCount(t *testing.T) {
	tests := []struct {
		name   string
		result ComplianceResult
		want   int
	}{
		{
			name: "no repos",
			result: ComplianceResult{
				Repos: []RepoResult{},
			},
			want: 0,
		},
		{
			name: "repos with no violations",
			result: ComplianceResult{
				Repos: []RepoResult{
					{Repo: Repo{FullName: "org/repo1"}},
					{Repo: Repo{FullName: "org/repo2"}},
				},
			},
			want: 0,
		},
		{
			name: "repos with violations",
			result: ComplianceResult{
				Repos: []RepoResult{
					{
						Repo: Repo{FullName: "org/repo1"},
						Violations: []Violation{
							{Policy: "policy1", Message: "violation1"},
							{Policy: "policy2", Message: "violation2"},
						},
					},
					{
						Repo: Repo{FullName: "org/repo2"},
						Violations: []Violation{
							{Policy: "policy3", Message: "violation3"},
						},
					},
				},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.ViolationCount()
			if got != tt.want {
				t.Errorf("ComplianceResult.ViolationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeverityValues(t *testing.T) {
	tests := []struct {
		severity Severity
		want     string
	}{
		{SeverityCritical, "critical"},
		{SeverityHigh, "high"},
		{SeverityMedium, "medium"},
		{SeverityLow, "low"},
		{SeverityInfo, "info"},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if string(tt.severity) != tt.want {
				t.Errorf("Severity = %v, want %v", tt.severity, tt.want)
			}
		})
	}
}
