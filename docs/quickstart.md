# Quick Start

Get up and running with PipelineConductor in minutes.

## Prerequisites

- PipelineConductor [installed](installation.md)
- GitHub personal access token with `repo` scope

## Step 1: Set Your GitHub Token

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

## Step 2: Run Your First Scan

Scan a single organization:

```bash
pipelineconductor scan --orgs myorg
```

This outputs a JSON compliance report to stdout.

## Step 3: Generate a Readable Report

Use Markdown format for human-readable output:

```bash
pipelineconductor scan --orgs myorg --format markdown
```

Example output:

```markdown
# Compliance Report

**Generated:** 2025-01-15T10:30:00Z
**Duration:** 1234ms
**Profile:** default

## Summary

| Metric | Value |
|--------|-------|
| Total Repos | 42 |
| Compliant | 38 |
| Non-Compliant | 4 |
| Compliance Rate | 90.5% |

## Repositories

### ✅ PASS myorg/api-server

*Scanned in 100ms*

### ❌ FAIL myorg/legacy-tool

**Violations:**

- 🟠 **[high]** ci/workflow-required: No CI/CD workflow found
  - 💡 Remediation: Create a .github/workflows/ci.yml file
```

## Step 4: Save the Report

Write the report to a file:

```bash
pipelineconductor scan --orgs myorg --format markdown --output report.md
```

## Step 5: Filter Repositories

Scan only Go repositories:

```bash
pipelineconductor scan --orgs myorg --languages Go
```

Scan multiple organizations:

```bash
pipelineconductor scan --orgs org1,org2,org3
```

Exclude archived and forked repos (default behavior):

```bash
pipelineconductor scan --orgs myorg
```

Include them:

```bash
pipelineconductor scan --orgs myorg --include-archived --include-forks
```

## Step 6: Use a Profile

Profiles define expected CI/CD configurations:

```bash
# Use the modern profile (latest Go, fewer platforms)
pipelineconductor scan --orgs myorg --profile modern

# Use the legacy profile (older Go, Linux only)
pipelineconductor scan --orgs myorg --profile legacy
```

## Step 7: Validate Custom Policies

If you have custom Cedar policies:

```bash
# Validate policy syntax
pipelineconductor validate ./policies/

# Use custom policies in scan
pipelineconductor scan --orgs myorg --policy-dir ./policies/
```

## Common Workflows

### Daily Compliance Check

```bash
#!/bin/bash
DATE=$(date +%Y-%m-%d)
pipelineconductor scan \
  --orgs myorg \
  --format markdown \
  --output "reports/compliance-${DATE}.md"
```

### SARIF for GitHub Security

```bash
pipelineconductor scan \
  --orgs myorg \
  --format sarif \
  --output results.sarif
```

### CSV for Spreadsheet Analysis

```bash
pipelineconductor scan \
  --orgs myorg \
  --format csv \
  --output compliance.csv
```

## Next Steps

- [CLI Reference](cli/overview.md) - Full command documentation
- [Policies](policies/overview.md) - Learn about policy-as-code
- [Profiles](profiles/overview.md) - Configure profiles for your needs
- [GitHub Actions Integration](integration/github-actions.md) - Automate scans in CI
