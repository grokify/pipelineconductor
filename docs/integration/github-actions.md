# GitHub Actions Integration

Automate compliance scanning with GitHub Actions.

## Basic Workflow

Create `.github/workflows/compliance.yml`:

```yaml
name: Compliance Scan

on:
  schedule:
    - cron: '0 6 * * *'  # Daily at 6 AM UTC
  workflow_dispatch:      # Manual trigger

jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
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
            --format markdown \
            --output report.md

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: compliance-report
          path: report.md
```

## Token Setup

### Create a Personal Access Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Create a new token (classic) with scopes:
   - `repo` - Access repositories
   - `read:org` - Read organization membership
3. Copy the token

### Add as Repository Secret

1. Go to repository **Settings** > **Secrets and variables** > **Actions**
2. Click **New repository secret**
3. Name: `SCAN_TOKEN`
4. Value: Your personal access token

## Workflow Examples

### SARIF Upload to Security Tab

```yaml
name: Compliance Scan with SARIF

on:
  schedule:
    - cron: '0 6 * * *'
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

      - name: Run scan
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs ${{ github.repository_owner }} \
            --format sarif \
            --output results.sarif

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
          category: compliance
```

### Multiple Output Formats

```yaml
name: Compliance Reports

on:
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Monday
  workflow_dispatch:

jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Run scans
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          # JSON for archival
          pipelineconductor scan --orgs myorg --format json --output report.json

          # Markdown for review
          pipelineconductor scan --orgs myorg --format markdown --output report.md

          # CSV for spreadsheets
          pipelineconductor scan --orgs myorg --format csv --output report.csv

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: compliance-reports
          path: |
            report.json
            report.md
            report.csv
```

### Slack Notification

```yaml
name: Compliance Scan with Slack

on:
  schedule:
    - cron: '0 6 * * *'

jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Run scan
        id: scan
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs myorg \
            --format json \
            --output report.json

          # Extract summary
          TOTAL=$(jq '.summary.total' report.json)
          COMPLIANT=$(jq '.summary.compliant' report.json)
          RATE=$(jq '.summary.complianceRate' report.json)

          echo "total=$TOTAL" >> $GITHUB_OUTPUT
          echo "compliant=$COMPLIANT" >> $GITHUB_OUTPUT
          echo "rate=$RATE" >> $GITHUB_OUTPUT

      - name: Send Slack notification
        uses: slackapi/slack-github-action@v1
        with:
          payload: |
            {
              "text": "Compliance Report",
              "blocks": [
                {
                  "type": "section",
                  "text": {
                    "type": "mrkdwn",
                    "text": "*Daily Compliance Report*\n• Total: ${{ steps.scan.outputs.total }}\n• Compliant: ${{ steps.scan.outputs.compliant }}\n• Rate: ${{ steps.scan.outputs.rate }}%"
                  }
                }
              ]
            }
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### Fail on Non-Compliance

```yaml
name: Compliance Gate

on:
  pull_request:
    branches: [main]

jobs:
  check:
    runs-on: ubuntu-latest

    steps:
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Run scan
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs ${{ github.repository_owner }} \
            --format json \
            --output report.json

      - name: Check compliance
        run: |
          COMPLIANT=$(jq '.summary.compliant' report.json)
          TOTAL=$(jq '.summary.total' report.json)

          if [ "$COMPLIANT" -ne "$TOTAL" ]; then
            echo "::error::Not all repositories are compliant"
            jq '.repos[] | select(.compliant == false) | .repo.fullName' report.json
            exit 1
          fi

          echo "All repositories are compliant"
```

## Caching

Speed up workflows by caching Go modules:

```yaml
steps:
  - name: Set up Go
    uses: actions/setup-go@v5
    with:
      go-version: '1.25'
      cache: true

  - name: Install PipelineConductor
    run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest
```

## Scheduled Scans

Common cron schedules:

| Schedule | Cron | Description |
|----------|------|-------------|
| Daily | `0 6 * * *` | 6 AM UTC daily |
| Weekly | `0 6 * * 1` | Monday at 6 AM |
| Monthly | `0 6 1 * *` | 1st of month at 6 AM |

## Troubleshooting

### Token Permission Errors

Ensure your token has:

- `repo` scope for private repos
- `read:org` for organization access

### Rate Limiting

PipelineConductor includes built-in rate limit handling with:

- Automatic retry on 429 (Too Many Requests) and 403 (rate limit exceeded)
- Exponential backoff with jitter
- Respect for GitHub's `X-RateLimit-*` and `Retry-After` headers
- Up to 5 retry attempts by default

For very large scans across many organizations, you can still add delays between orgs if needed:

```yaml
- name: Scan with delay
  run: |
    for org in org1 org2 org3; do
      pipelineconductor scan --orgs $org --output "${org}.json"
      sleep 60  # Wait between orgs
    done
```

### Verbose Output

Add `-v` for debugging:

```yaml
- name: Run scan (verbose)
  run: pipelineconductor scan --orgs myorg -v
```

## See Also

- [SARIF Integration](../reports/sarif.md) - Security tab setup
- [Configuration](../cli/config.md) - Config file options
