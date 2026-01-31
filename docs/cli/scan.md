# scan Command

The `scan` command scans repositories across one or more GitHub organizations, evaluates them against policies, and generates a compliance report.

## Synopsis

```bash
pipelineconductor scan [flags]
```

## Description

The scan command:

1. Lists repositories from specified GitHub organizations
2. Filters repositories based on criteria (language, topics, archived status)
3. Fetches workflow files and branch protection settings
4. Evaluates Cedar policies against each repository
5. Generates a compliance report in the specified format

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | stdout |
| `--format` | `-f` | Output format: `json`, `markdown`, `sarif`, `csv` | `json` |
| `--include-archived` | | Include archived repositories | `false` |
| `--include-forks` | | Include forked repositories | `false` |
| `--languages` | | Filter by languages (comma-separated) | all |
| `--topics` | | Filter by topics (comma-separated) | all |
| `--policy-dir` | | Directory containing Cedar policy files | - |
| `--builtin-policies` | | Use built-in policies | `true` |
| `--evaluate-policies` | | Evaluate Cedar policies | `true` |

## Examples

### Basic Scan

```bash
pipelineconductor scan --orgs myorg
```

### Multiple Organizations

```bash
pipelineconductor scan --orgs org1,org2,org3
```

### Filter by Language

```bash
# Only Go repositories
pipelineconductor scan --orgs myorg --languages Go

# Go and Python
pipelineconductor scan --orgs myorg --languages Go,Python
```

### Filter by Topic

```bash
pipelineconductor scan --orgs myorg --topics api,production
```

### Include Archived Repositories

```bash
pipelineconductor scan --orgs myorg --include-archived
```

### Output Formats

```bash
# JSON (default)
pipelineconductor scan --orgs myorg --format json

# Markdown (human-readable)
pipelineconductor scan --orgs myorg --format markdown

# SARIF (GitHub Security integration)
pipelineconductor scan --orgs myorg --format sarif

# CSV (spreadsheet)
pipelineconductor scan --orgs myorg --format csv
```

### Save to File

```bash
pipelineconductor scan --orgs myorg --format markdown --output report.md
```

### Use Custom Policies

```bash
pipelineconductor scan --orgs myorg --policy-dir ./policies/
```

### Disable Built-in Policies

```bash
pipelineconductor scan --orgs myorg --builtin-policies=false --policy-dir ./policies/
```

### Verbose Output

```bash
pipelineconductor scan --orgs myorg -v
```

Output includes:

```
Loaded built-in policies
Using profile: default
Scanning organizations: [myorg]
Found 42 repositories
Scanning: myorg/repo1
Scanning: myorg/repo2
...
Report written to: report.json
```

## Output

### JSON Format

```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "summary": {
    "total": 42,
    "compliant": 38,
    "nonCompliant": 4,
    "complianceRate": 90.5
  },
  "repos": [
    {
      "repo": {
        "fullName": "myorg/api-server",
        "name": "api-server"
      },
      "compliant": true,
      "violations": [],
      "scanTimeMs": 100
    }
  ],
  "scanDurationMs": 1234
}
```

### Markdown Format

See [Report Formats](../reports/formats.md) for details.

### SARIF Format

See [SARIF Integration](../reports/sarif.md) for details.

## Built-in Checks

When `--evaluate-policies` is enabled (default), the scan performs these checks:

| Check | Severity | Description |
|-------|----------|-------------|
| Workflow exists | High | Repository has at least one GitHub Actions workflow |
| Branch protection | Medium | Default branch has protection enabled |

## Rate Limiting

PipelineConductor handles GitHub API rate limits automatically:

- **Automatic retry** on 429 (Too Many Requests) and 403 (rate limit exceeded)
- **Exponential backoff** starting at 1 second, up to 60 seconds maximum
- **Jitter** to prevent thundering herd problems
- **Header awareness** respects `X-RateLimit-Remaining`, `X-RateLimit-Reset`, and `Retry-After`
- **Up to 5 retries** by default before failing

This allows scanning large organizations without manual intervention for rate limits.

## See Also

- [validate](validate.md) - Validate policy files
- [Policies](../policies/overview.md) - Policy-as-code documentation
- [Report Formats](../reports/formats.md) - Output format details
