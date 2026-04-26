# Plan: Workflow Compliance Check CLI

## Overview

Add a new `check` command to pipelineconductor that scans public repos across GitHub orgs and users, filters by language (Go, TypeScript, Crystal), and checks workflow compliance against reference workflows in `grokify/.github`.

**Key Requirements:**

- CLI uses Cobra library (already in use)
- JSON output backed by JSON Schema (Draft 2020-12)
- Schema must pass `schemalint lint --profile scale` (strict mode)

## New Files to Create

```
cmd/pipelineconductor/cmd/check.go       # Main command implementation
internal/compliance/compliance.go         # Core compliance checking logic
internal/compliance/reference.go          # Reference workflow fetching
internal/compliance/matcher.go            # Workflow matching/comparison
internal/compliance/rules.go              # Language-to-workflow mapping
pkg/model/compliance_check.go             # Go structs (source of truth)
schema/check_result.schema.json           # Generated JSON Schema
schema/schema.go                          # //go:embed for schema access
internal/report/check_markdown.go         # Markdown report for check results
```

## Files to Modify

```
internal/collector/collector.go          # Add ListUserRepos interface
internal/collector/github.go             # Implement ListUserRepos, ListReposMultiSource
```

## Command Interface

```bash
pipelineconductor check \
  --orgs agentplexus \
  --users grokify \
  --languages Go,TypeScript,Crystal \
  --ref-repo grokify/.github \
  --format markdown
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--users`, `-u` | nil | GitHub users to scan |
| `--languages`, `-l` | nil | Filter: Go, TypeScript, Crystal |
| `--ref-repo`, `-r` | `grokify/.github` | Reference workflow repository |
| `--ref-branch` | `main` | Branch in reference repo |
| `--output`, `-o` | stdout | Output file path |
| `--format`, `-f` | `json` | Output: json, markdown |
| `--strict` | false | Require exact reusable workflow usage |

## Language-to-Workflow Rules

```go
"Go": [
    {Type: "go-ci",          Path: ".github/workflows/go-ci.yaml",          Severity: High},
    {Type: "go-lint",        Path: ".github/workflows/go-lint.yaml",        Severity: Medium},
    {Type: "go-sast-codeql", Path: ".github/workflows/go-sast-codeql.yaml", Severity: Low},
]
"TypeScript": [
    {Type: "ts-ci",   Path: ".github/workflows/ts-ci.yaml",   Severity: High},
    {Type: "ts-lint", Path: ".github/workflows/ts-lint.yaml", Severity: Medium},
]
"Crystal": []  // No reference workflows yet
```

## Core Data Models (Go structs - source of truth)

Structs designed for schemalint scale profile compliance:

- No unions (anyOf/oneOf) - use string enums instead
- Explicit types on all fields
- No additionalProperties
- camelCase JSON field names

```go
// CheckResult - top-level result
type CheckResult struct {
    SchemaVersion  string            `json:"schemaVersion"`  // "1.0.0"
    Timestamp      string            `json:"timestamp"`      // RFC3339 format
    Summary        CheckSummary      `json:"summary"`
    Repos          []RepoCheckResult `json:"repos"`
    ScanDurationMs int64             `json:"scanDurationMs"`
    Config         CheckConfig       `json:"config"`
}

// CheckSummary - aggregate statistics
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

// LanguageComplianceStats - per-language breakdown (array, not map)
type LanguageComplianceStats struct {
    Language       string  `json:"language"`
    TotalRepos     int     `json:"totalRepos"`
    CompliantRepos int     `json:"compliantRepos"`
    ComplianceRate float64 `json:"complianceRate"`
}

// RepoCheckResult - per-repo result
type RepoCheckResult struct {
    Owner             string              `json:"owner"`
    Name              string              `json:"name"`
    FullName          string              `json:"fullName"`
    HTMLURL           string              `json:"htmlUrl"`
    Languages         []string            `json:"languages"`
    Compliant         bool                `json:"compliant"`
    ComplianceLevel   string              `json:"complianceLevel"`  // "full", "partial", "none"
    RequiredWorkflows []WorkflowCheck     `json:"requiredWorkflows"`
    ActualWorkflows   []WorkflowInfo      `json:"actualWorkflows"`
    Missing           []MissingWorkflow   `json:"missing"`
    Skipped           bool                `json:"skipped"`
    SkipReason        string              `json:"skipReason"`
    Error             string              `json:"error"`
    ScanTimeMs        int64               `json:"scanTimeMs"`
}

// WorkflowCheck - per-workflow check result
type WorkflowCheck struct {
    WorkflowType   string `json:"workflowType"`    // "go-ci", "go-lint", etc.
    Required       bool   `json:"required"`
    Present        bool   `json:"present"`
    UsesReusable   bool   `json:"usesReusable"`
    ReusableRef    string `json:"reusableRef"`
    ExpectedRef    string `json:"expectedRef"`
    MatchType      string `json:"matchType"`       // "exact", "equivalent", "partial", "none"
    ActualWorkflow string `json:"actualWorkflow"`
}

// WorkflowInfo - info about existing workflow
type WorkflowInfo struct {
    Name             string   `json:"name"`
    Path             string   `json:"path"`
    UsesReusable     bool     `json:"usesReusable"`
    ReusableRefs     []string `json:"reusableRefs"`
    DetectedLanguage string   `json:"detectedLanguage"`
}

// MissingWorkflow - a required workflow that's missing
type MissingWorkflow struct {
    Language     string `json:"language"`
    WorkflowType string `json:"workflowType"`
    RefPath      string `json:"refPath"`
    Severity     string `json:"severity"`          // "high", "medium", "low"
    Description  string `json:"description"`
}

// CheckConfig - scan configuration
type CheckConfig struct {
    Orgs       []string `json:"orgs"`
    Users      []string `json:"users"`
    RefRepo    string   `json:"refRepo"`
    RefBranch  string   `json:"refBranch"`
    Languages  []string `json:"languages"`
    Strict     bool     `json:"strict"`
}
```

## JSON Schema Generation

**Go-first approach**: Generate schema from Go structs using `github.com/invopop/jsonschema`.

```go
// schema/generate.go (build tool)
func main() {
    r := jsonschema.Reflector{
        DoNotReference: true,  // Inline all definitions for scale profile
    }
    schema := r.Reflect(&model.CheckResult{})
    schema.Version = "https://json-schema.org/draft/2020-12/schema"
    // Write to schema/check_result.schema.json
}
```

**Schema validation in CI:**

```bash
schemalint lint --profile scale schema/check_result.schema.json
```

**Schema embedding:**

```go
// schema/schema.go
//go:embed check_result.schema.json
var CheckResultSchema []byte
```

## Implementation Flow

1. **Parse flags** - Get orgs, users, languages, ref-repo
2. **Create collector** - GitHub API client with rate limiting
3. **Fetch reference workflows** - From grokify/.github
4. **List repos** - From orgs + users, filtered by language
5. **For each repo:**
   - Get workflows via API
   - Determine required workflows based on repo languages
   - Check each required workflow:
     - Does workflow file exist? (by type matching: go-ci, go-lint, etc.)
     - Does it use reusable workflow from reference repo? (exact match)
     - Or has equivalent inline workflow? (partial match)
   - Record compliance level
6. **Calculate summary** - Compliant, partial, non-compliant counts by language
7. **Generate report** - JSON or Markdown format

## Compliance Matching Logic

```
For each required workflow type (e.g., "go-ci"):
  1. Check if repo uses reusable workflow: grokify/.github/.github/workflows/go-ci.yaml@main
     → Match: exact (fully compliant)
  2. Check if repo has workflow named go-ci.* or ci.* with Go build steps
     → Match: equivalent (partial - has functionality but not using shared workflow)
  3. No matching workflow found
     → Match: none (non-compliant)
```

## Example Output (Markdown)

```markdown
# Workflow Compliance Report

**Reference:** grokify/.github@main | **Scanned:** 45 repos | **Compliance:** 26.7%

## Summary

| Status | Count |
|--------|-------|
| Fully Compliant | 12 |
| Partial | 18 |
| Non-Compliant | 15 |

## Non-Compliant: grokify/mogo

**Languages:** Go | **Level:** partial

| Workflow | Present | Reusable | Match |
|----------|---------|----------|-------|
| go-ci | ✅ | No | equivalent |
| go-lint | ✅ | No | equivalent |
| go-sast-codeql | ❌ | No | none |
```

## Verification

1. **Generate schema:**

   ```bash
   go run ./schema/generate
   ```

2. **Validate schema (strict):**

   ```bash
   schemalint lint --profile scale schema/check_result.schema.json
   ```

3. **Build:** `go build ./cmd/pipelineconductor`

4. **Run check (JSON):**

   ```bash
   ./pipelineconductor check --orgs agentplexus --users grokify --languages Go -f json
   ```

5. **Run check (Markdown):**

   ```bash
   ./pipelineconductor check --orgs agentplexus --users grokify --languages Go -f markdown
   ```

6. **Verify grokify/mogo shows as non-compliant** (has Go workflows but doesn't use reference)

7. **Run tests:** `go test -v ./...`

8. **Lint:** `golangci-lint run`

## Implementation Order

1. `pkg/model/compliance_check.go` - Data models (schemalint-compliant)
2. `schema/generate/main.go` - Schema generation tool
3. `schema/check_result.schema.json` - Generate and validate with schemalint
4. `schema/schema.go` - Embed schema
5. `internal/collector/collector.go` + `github.go` - Add user repo listing
6. `internal/compliance/rules.go` - Language-to-workflow mapping
7. `internal/compliance/reference.go` - Fetch reference workflows
8. `internal/compliance/matcher.go` - Workflow comparison
9. `internal/compliance/compliance.go` - Main checker
10. `internal/report/check_json.go` - JSON formatter with schema validation
11. `internal/report/check_markdown.go` - Markdown formatter
12. `cmd/pipelineconductor/cmd/check.go` - CLI command
13. Tests for new packages
