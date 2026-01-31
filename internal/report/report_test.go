package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/grokify/pipelineconductor/pkg/model"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"markdown", FormatMarkdown, false},
		{"md", FormatMarkdown, false},
		{"sarif", FormatSARIF, false},
		{"csv", FormatCSV, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}
	if len(builder.formatters) != 4 {
		t.Errorf("NewBuilder() has %d formatters, want 4", len(builder.formatters))
	}
}

func sampleResult() *model.ComplianceResult {
	return &model.ComplianceResult{
		Timestamp:      time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
		ScanDurationMs: 1234,
		Config: model.ScanConfig{
			Orgs:    []string{"testorg"},
			Profile: "default",
		},
		Summary: model.ScanSummary{
			TotalRepos:     2,
			CompliantRepos: 1,
			NonCompliant:   1,
			ComplianceRate: 50.0,
		},
		Repos: []model.RepoResult{
			{
				Repo: model.Repo{
					FullName: "testorg/repo1",
					Name:     "repo1",
					Owner:    "testorg",
				},
				Compliant:  true,
				ScanTimeMs: 100,
			},
			{
				Repo: model.Repo{
					FullName: "testorg/repo2",
					Name:     "repo2",
					Owner:    "testorg",
				},
				Compliant: false,
				Violations: []model.Violation{
					{
						Policy:      "ci/workflow-required",
						Rule:        "has-workflow",
						Message:     "No CI/CD workflow found",
						Severity:    model.SeverityHigh,
						Remediation: "Create a .github/workflows/ci.yml file",
					},
				},
				ScanTimeMs: 150,
			},
		},
	}
}

func TestBuilderGenerateJSON(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	output, err := builder.Generate(result, FormatJSON)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(string(output), `"fullName": "testorg/repo1"`) {
		t.Error("JSON output missing expected repo1")
	}
	if !strings.Contains(string(output), `"compliant": false`) {
		t.Error("JSON output missing expected non-compliant status")
	}
}

func TestBuilderGenerateMarkdown(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	output, err := builder.Generate(result, FormatMarkdown)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	md := string(output)
	if !strings.Contains(md, "# Compliance Report") {
		t.Error("Markdown output missing header")
	}
	if !strings.Contains(md, "testorg/repo1") {
		t.Error("Markdown output missing repo1")
	}
	if !strings.Contains(md, "No CI/CD workflow found") {
		t.Error("Markdown output missing violation message")
	}
}

func TestBuilderGenerateSARIF(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	output, err := builder.Generate(result, FormatSARIF)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	sarif := string(output)
	if !strings.Contains(sarif, `"version": "2.1.0"`) {
		t.Error("SARIF output missing version")
	}
	if !strings.Contains(sarif, `"name": "PipelineConductor"`) {
		t.Error("SARIF output missing tool name")
	}
	if !strings.Contains(sarif, "ci/workflow-required") {
		t.Error("SARIF output missing rule ID")
	}
}

func TestBuilderGenerateCSV(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	output, err := builder.Generate(result, FormatCSV)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	csv := string(output)
	// Check header
	if !strings.Contains(csv, "repo,org,compliant") {
		t.Error("CSV output missing header")
	}
	// Check data
	if !strings.Contains(csv, "testorg/repo1") {
		t.Error("CSV output missing repo1")
	}
	if !strings.Contains(csv, "testorg/repo2") {
		t.Error("CSV output missing repo2")
	}
}

func TestBuilderWrite(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	var buf bytes.Buffer
	err := builder.Write(result, FormatJSON, &buf)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Write() produced empty output")
	}
}

func TestBuilderGenerateUnsupportedFormat(t *testing.T) {
	builder := NewBuilder()
	result := sampleResult()

	_, err := builder.Generate(result, Format("unsupported"))
	if err == nil {
		t.Error("Generate() with unsupported format should return error")
	}
}
