# check

The `check` command scans repositories for workflow compliance against a reference repository containing reusable workflows.

## Synopsis

```bash
pipelineconductor check [flags]
```

## Description

The check command analyzes repositories to determine if they are using the required reusable workflows from a reference repository. It supports:

- **GitHub API scanning**: Scan public repositories from organizations and users
- **Local filesystem scanning**: Scan locally cloned repositories without requiring a GitHub token
- **Multiple output formats**: JSON for programmatic use, Markdown for human-readable reports
- **Strict mode**: Require exact reusable workflow usage (no equivalent matching)
- **Dashboard generation**: Create Dashforge-compatible dashboard JSON

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--local` | | Scan local filesystem instead of GitHub API | - |
| `--languages` | `-l` | Filter by languages (Go, TypeScript, Crystal) | (required) |
| `--users` | `-u` | GitHub users to scan | - |
| `--ref-repo` | `-r` | Reference workflow repository (owner/repo) | `grokify/.github` |
| `--ref-branch` | | Branch in reference repo | `main` |
| `--output` | `-o` | Output file path | stdout |
| `--format` | `-f` | Output format: json, markdown, html | `json` |
| `--strict` | | Require exact reusable workflow usage | `false` |
| `--include-archived` | | Include archived repositories | `false` |
| `--include-forks` | | Include forked repositories | `false` |
| `--dashboard` | `-d` | Generate Dashforge dashboard JSON to this path | - |
| `--data-url` | | Data URL for dashboard | relative path |

## Compliance Levels

The check command categorizes repositories into three compliance levels:

| Level | Description |
|-------|-------------|
| **Full** | All required workflows use the exact reusable workflows from the reference repo |
| **Partial** | Workflows are present but use equivalent implementations (not reusable workflows) |
| **None** | Required workflows are missing |

## Match Types

| Type | Description |
|------|-------------|
| `exact` | Uses the reusable workflow from the reference repository |
| `equivalent` | Has a workflow that performs the same function but doesn't use the reusable workflow |
| `none` | No matching workflow found |

## Examples

### Local Filesystem Scanning

Scan locally cloned repositories without requiring a GitHub token:

```bash
# Scan all Go repos in a directory
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --ref-repo myorg/.github

# Scan with markdown output
pipelineconductor check \
  --local ~/projects \
  --orgs mycompany \
  --languages Go,TypeScript \
  --format markdown \
  --output compliance-report.md
```

### GitHub API Scanning

Scan repositories via the GitHub API (requires token):

```bash
# Scan organization repositories
export GITHUB_TOKEN=ghp_xxx
pipelineconductor check \
  --orgs myorg \
  --languages Go \
  --ref-repo myorg/.github

# Scan user repositories
pipelineconductor check \
  --users johndoe \
  --languages TypeScript \
  --format markdown
```

### Strict Mode

Require exact reusable workflow usage (no equivalent matching):

```bash
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --strict
```

### Generate Dashboard

Create a Dashforge-compatible dashboard alongside the JSON output:

```bash
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --output results.json \
  --dashboard dashboard.json
```

## Output Formats

### JSON Output

```json
{
  "schemaVersion": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "summary": {
    "totalRepos": 50,
    "compliantRepos": 25,
    "partialRepos": 15,
    "nonCompliant": 10,
    "complianceRate": 50.0,
    "byLanguage": [
      {
        "language": "Go",
        "totalRepos": 50,
        "compliantRepos": 25,
        "complianceRate": 50.0
      }
    ]
  },
  "repos": [
    {
      "fullName": "myorg/myrepo",
      "complianceLevel": "full",
      "requiredWorkflows": [
        {
          "workflowType": "go-ci",
          "present": true,
          "usesReusable": true,
          "matchType": "exact"
        }
      ]
    }
  ]
}
```

### Markdown Output

The markdown format produces a human-readable report with:

- Summary statistics table
- Per-language breakdown
- Non-compliant repositories with missing workflows
- Partially compliant repositories with filename mismatch warnings
- Fully compliant repositories list

### HTML Output

The HTML format generates a self-contained web page with:

- Summary cards showing total, compliant, partial, and non-compliant counts
- Visual progress bar showing compliance breakdown
- Per-language statistics table
- Interactive filter buttons to show/hide repository categories
- Detailed tables for each compliance level with workflow status
- Configuration section showing scan parameters

```bash
# Generate HTML report
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --format html \
  --output compliance-report.html
```

## Filename Mismatch Detection

When a repository has an equivalent workflow but with a non-standard filename, the check command will flag this as a filename mismatch:

```markdown
### myorg/myrepo

| Workflow | Present | Reusable | Match |
|----------|---------|----------|-------|
| go-ci | Yes | No | equivalent |

**Filename Mismatches:**

- `go-ci`: expected `go-ci.yaml`, found `ci.yaml`
```

## Workflow Rules

The check command uses predefined rules for each supported language:

### Go

| Workflow | Expected Filename | Severity |
|----------|-------------------|----------|
| go-ci | `go-ci.yaml` | high |
| go-lint | `go-lint.yaml` | medium |
| go-sast-codeql | `go-sast-codeql.yaml` | medium |

### TypeScript

| Workflow | Expected Filename | Severity |
|----------|-------------------|----------|
| ts-ci | `ts-ci.yaml` | high |
| ts-lint | `ts-lint.yaml` | medium |

## Cedar Policy Evaluation

The check command can evaluate Cedar policies against compliance results:

```bash
# Evaluate policies from a directory
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --policies ./policies/compliance/ \
  --policy-action merge

# Fail if any repo is denied
pipelineconductor check \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --policies ./policies/strict-compliance.cedar \
  --policy-action deploy \
  --fail-on-deny
```

### Policy Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--policies` | Path to Cedar policy files or directory | - |
| `--policy-action` | Action to evaluate (merge, build, deploy, release) | `merge` |
| `--fail-on-deny` | Exit with error if any policy denies the action | `false` |

### Available Policy Context

Cedar policies can access the following compliance context:

| Field | Type | Description |
|-------|------|-------------|
| `complianceLevel` | string | "full", "partial", or "none" |
| `compliant` | boolean | True if fully compliant |
| `complianceRate` | long | Percentage (0-100) |
| `missingWorkflowCount` | long | Number of missing workflows |
| `missingWorkflows` | set | Set of missing workflow types |
| `hasFilenameMismatch` | boolean | True if any filename mismatches |
| `usesReusableWorkflows` | boolean | True if using reusable workflows |
| `exactMatchCount` | long | Count of exact matches |
| `equivalentMatchCount` | long | Count of equivalent matches |

### Example Policy

```cedar
// Only allow merge for fully compliant repos
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.compliant == true
};
```

## See Also

- [remediate](remediate.md) - Generate compliant workflow files
- [apply](apply.md) - Apply remediation with git operations
- [GitHub Actions Integration](../integration/github-actions.md)
- [Cedar Policies](../policies/overview.md)
