# remediate

The `remediate` command generates compliant workflow files for repositories that are missing required workflows.

## Synopsis

```bash
pipelineconductor remediate --local <path> --orgs <org> --languages <langs> [flags]
```

## Description

The remediate command scans local repositories and generates workflow files that call the reference reusable workflows from the organization's `.github` repository. It:

- Identifies repositories with missing workflows
- Generates workflow files using templates that call reusable workflows
- Supports dry-run mode to preview changes
- Can target specific repositories or process all non-compliant repos

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--local` | | Base path for local filesystem scanning | (required) |
| `--languages` | `-l` | Filter by languages (Go, TypeScript, Crystal) | (required) |
| `--ref-repo` | `-r` | Reference workflow repository (owner/repo) | `plexusone/.github` |
| `--ref-branch` | | Branch in reference repo | `main` |
| `--repo` | | Target specific repository name | - |
| `--dry-run` | | Show what would be generated without writing files | `false` |
| `--overwrite` | | Overwrite existing workflow files | `false` |
| `--output` | `-o` | Output remediation report to file | stdout |
| `--format` | `-f` | Output format: text, json | `text` |

## Examples

### Dry Run

Preview what workflows would be generated:

```bash
pipelineconductor remediate \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --dry-run
```

Output:

```
Remediation Report
==================

Reference: myorg/.github@main
Dry Run: true
Total Repos: 50
Repos Remediated: 10
Files Would generate: 30

## myorg/repo1
   Path: /Users/dev/go/src/github.com/myorg/repo1
   - .github/workflows/go-ci.yaml (would be created)
   - .github/workflows/go-lint.yaml (would be created)
   - .github/workflows/go-sast-codeql.yaml (would be created)

## myorg/repo2
   Path: /Users/dev/go/src/github.com/myorg/repo2
   - .github/workflows/go-ci.yaml (would be created)
```

### Generate Workflows

Generate workflow files for all non-compliant repos:

```bash
pipelineconductor remediate \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go
```

### Target Specific Repository

Generate workflows for a single repository:

```bash
pipelineconductor remediate \
  --local ~/go/src/github.com \
  --orgs myorg \
  --repo my-service \
  --languages Go
```

### Overwrite Existing Files

Replace existing workflow files with updated versions:

```bash
pipelineconductor remediate \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --overwrite
```

### JSON Output

Get machine-readable output for scripting:

```bash
pipelineconductor remediate \
  --local ~/go/src/github.com \
  --orgs myorg \
  --languages Go \
  --format json \
  --output remediation.json
```

## Generated Workflow Templates

### Go CI Workflow

Generated file: `.github/workflows/go-ci.yaml`

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  ci:
    uses: myorg/.github/.github/workflows/go-ci.yaml@main
    secrets: inherit
```

### Go Lint Workflow

Generated file: `.github/workflows/go-lint.yaml`

```yaml
name: Go Lint

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  lint:
    uses: myorg/.github/.github/workflows/go-lint.yaml@main
    secrets: inherit
```

### Go SAST CodeQL Workflow

Generated file: `.github/workflows/go-sast-codeql.yaml`

```yaml
name: Go SAST CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 6 * * 1'

permissions:
  contents: read
  security-events: write

jobs:
  codeql:
    uses: myorg/.github/.github/workflows/go-sast-codeql.yaml@main
    secrets: inherit
```

## Output Formats

### Text Format (default)

Human-readable report showing:

- Summary statistics
- Per-repository file list
- Status indicators (created/overwritten)

### JSON Format

Machine-readable output for CI/CD integration:

```json
{
  "dryRun": false,
  "refRepo": "myorg/.github",
  "refBranch": "main",
  "totalRepos": 50,
  "remediatedRepos": 10,
  "filesGenerated": 30,
  "repos": [
    {
      "owner": "myorg",
      "name": "repo1",
      "fullName": "myorg/repo1",
      "localPath": "/path/to/repo1",
      "files": [
        {
          "workflowType": "go-ci",
          "filename": "go-ci.yaml",
          "relativePath": ".github/workflows/go-ci.yaml",
          "absolutePath": "/path/to/repo1/.github/workflows/go-ci.yaml",
          "wouldOverwrite": false
        }
      ]
    }
  ]
}
```

## See Also

- [check](check.md) - Check workflow compliance
- [apply](apply.md) - Apply remediation with git operations
