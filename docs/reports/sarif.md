# SARIF Integration

SARIF (Static Analysis Results Interchange Format) is a standard format for static analysis results. PipelineConductor can generate SARIF reports that integrate with GitHub's Security tab.

## What is SARIF?

SARIF is an OASIS standard (version 2.1.0) for expressing static analysis results. GitHub supports uploading SARIF files to display results in the Security tab.

## Generating SARIF

```bash
pipelineconductor scan --orgs myorg --format sarif --output results.sarif
```

## SARIF Structure

PipelineConductor generates SARIF with:

- **Tool information** - PipelineConductor version and metadata
- **Rules** - Policy definitions with IDs and descriptions
- **Results** - Violation instances with locations and messages
- **Fixes** - Remediation suggestions

### Example Output

```json
{
  "version": "2.1.0",
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "PipelineConductor",
          "version": "0.1.0",
          "informationUri": "https://github.com/grokify/pipelineconductor",
          "rules": [
            {
              "id": "ci/workflow-required/has-workflow",
              "name": "ci/workflow-required",
              "shortDescription": {
                "text": "No CI/CD workflow found"
              },
              "defaultConfiguration": {
                "level": "error"
              }
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "ci/workflow-required/has-workflow",
          "ruleIndex": 0,
          "level": "error",
          "message": {
            "text": "[myorg/legacy-tool] No CI/CD workflow found"
          },
          "locations": [
            {
              "logicalLocations": [
                {
                  "name": "legacy-tool",
                  "fullyQualifiedName": "myorg/legacy-tool",
                  "kind": "repository"
                }
              ]
            }
          ],
          "fixes": [
            {
              "description": {
                "text": "Create a .github/workflows/ci.yml file"
              }
            }
          ]
        }
      ]
    }
  ]
}
```

## Severity Mapping

PipelineConductor severities map to SARIF levels:

| PipelineConductor | SARIF Level |
|-------------------|-------------|
| critical | error |
| high | error |
| medium | warning |
| low | note |
| info | note |

## GitHub Actions Integration

Upload SARIF results to GitHub Security:

```yaml
# .github/workflows/compliance.yml
name: Compliance Scan

on:
  schedule:
    - cron: '0 6 * * *'  # Daily at 6 AM
  workflow_dispatch:

jobs:
  scan:
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      contents: read

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Run compliance scan
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs ${{ github.repository_owner }} \
            --format sarif \
            --output results.sarif

      - name: Upload SARIF results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
          category: pipelineconductor
```

### Required Permissions

The workflow needs:

- `security-events: write` - To upload SARIF results
- `contents: read` - To access repository

The `SCAN_TOKEN` secret needs:

- `repo` scope - To scan repositories
- `read:org` - To list organization repos

## Viewing Results

After uploading SARIF:

1. Go to your repository on GitHub
2. Click **Security** tab
3. Click **Code scanning alerts**
4. Filter by tool: "PipelineConductor"

## Result Details

Each alert shows:

- **Rule** - Policy that was violated
- **Repository** - Affected repository
- **Message** - Violation description
- **Fix** - Remediation suggestion

## Dismissing Alerts

You can dismiss alerts in GitHub:

- **False positive** - Alert doesn't apply
- **Won't fix** - Acknowledged but not fixing
- **Used in tests** - Expected in test context

## Multiple Organizations

For multi-org scans, results are aggregated:

```bash
pipelineconductor scan \
  --orgs org1,org2,org3 \
  --format sarif \
  --output results.sarif
```

All violations appear in the same SARIF file with repository-qualified locations.

## Limitations

- SARIF is best for code-level findings
- Repository-level compliance is shown as logical locations
- Some GitHub Security features may not apply (e.g., code navigation)

## See Also

- [Output Formats](formats.md) - All report formats
- [GitHub Actions Integration](../integration/github-actions.md) - CI/CD setup
- [SARIF Specification](https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html) - Official spec
