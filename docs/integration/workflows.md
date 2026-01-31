# CI/CD Workflows

Best practices for integrating PipelineConductor into your CI/CD workflows.

## Workflow Patterns

### Centralized Policy Repository

Store policies in a dedicated repository:

```
myorg/ci-policies/
├── policies/
│   ├── go/
│   │   ├── merge.cedar
│   │   └── versions.cedar
│   └── security/
│       └── branch-protection.cedar
├── profiles/
│   ├── default.yaml
│   └── modern.yaml
└── .github/
    └── workflows/
        └── validate.yml
```

Validate policies on push:

```yaml
# .github/workflows/validate.yml
name: Validate Policies

on:
  push:
    paths:
      - 'policies/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Validate policies
        run: pipelineconductor validate policies/ --verbose
```

### Scheduled Compliance Scans

Run daily scans across your organization:

```yaml
# .github/workflows/daily-scan.yml
name: Daily Compliance Scan

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

      - name: Clone policy repo
        run: git clone https://github.com/myorg/ci-policies.git

      - name: Run scan
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs myorg \
            --policy-dir ci-policies/policies/ \
            --format markdown \
            --output report.md

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: compliance-report-${{ github.run_id }}
          path: report.md
          retention-days: 30
```

### PR-Based Compliance Checks

Check compliance on PRs to the policy repo:

```yaml
# .github/workflows/pr-check.yml
name: PR Compliance Check

on:
  pull_request:
    branches: [main]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Validate policies
        run: pipelineconductor validate policies/ --verbose

      - name: Test scan with new policies
        env:
          GITHUB_TOKEN: ${{ secrets.SCAN_TOKEN }}
        run: |
          pipelineconductor scan \
            --orgs myorg \
            --policy-dir policies/ \
            --format markdown \
            --output report.md

      - name: Comment on PR
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('report.md', 'utf8');

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '## Compliance Scan Preview\n\n' + report
            });
```

## Multi-Environment Strategy

### Development vs Production

Use different profiles for different environments:

```yaml
jobs:
  scan-dev:
    runs-on: ubuntu-latest
    steps:
      - name: Scan dev repos
        run: |
          pipelineconductor scan \
            --orgs myorg-dev \
            --profile legacy \
            --output dev-report.json

  scan-prod:
    runs-on: ubuntu-latest
    steps:
      - name: Scan production repos
        run: |
          pipelineconductor scan \
            --orgs myorg-prod \
            --profile default \
            --output prod-report.json
```

### Branch-Based Policies

Different policies for different branch protection levels:

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - org: myorg-strict
            profile: enterprise
          - org: myorg-standard
            profile: default
          - org: myorg-legacy
            profile: legacy

    steps:
      - name: Scan ${{ matrix.org }}
        run: |
          pipelineconductor scan \
            --orgs ${{ matrix.org }} \
            --profile ${{ matrix.profile }} \
            --output ${{ matrix.org }}.json
```

## Report Aggregation

Combine reports from multiple scans:

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        org: [org1, org2, org3]

    steps:
      - name: Scan ${{ matrix.org }}
        run: |
          pipelineconductor scan \
            --orgs ${{ matrix.org }} \
            --format json \
            --output ${{ matrix.org }}.json

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: report-${{ matrix.org }}
          path: ${{ matrix.org }}.json

  aggregate:
    needs: scan
    runs-on: ubuntu-latest
    steps:
      - name: Download all reports
        uses: actions/download-artifact@v4
        with:
          pattern: report-*
          merge-multiple: true

      - name: Aggregate reports
        run: |
          jq -s '
            {
              timestamp: now | todate,
              summary: {
                total: map(.summary.total) | add,
                compliant: map(.summary.compliant) | add,
                nonCompliant: map(.summary.nonCompliant) | add
              },
              orgs: map({org: .config.orgs[0], summary: .summary})
            }
          ' *.json > aggregate.json

      - name: Upload aggregate
        uses: actions/upload-artifact@v4
        with:
          name: aggregate-report
          path: aggregate.json
```

## Notification Patterns

### Email Summary

```yaml
- name: Send email
  uses: dawidd6/action-send-mail@v3
  with:
    server_address: smtp.gmail.com
    server_port: 587
    username: ${{ secrets.EMAIL_USER }}
    password: ${{ secrets.EMAIL_PASS }}
    subject: "Compliance Report - ${{ github.run_id }}"
    body: file://report.md
    to: team@example.com
    from: compliance-bot@example.com
```

### GitHub Issue

```yaml
- name: Create issue for failures
  if: steps.scan.outputs.failures > 0
  uses: actions/github-script@v7
  with:
    script: |
      github.rest.issues.create({
        owner: context.repo.owner,
        repo: context.repo.repo,
        title: 'Compliance violations detected',
        body: 'See workflow run: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}',
        labels: ['compliance', 'automated']
      });
```

## Best Practices

1. **Use dedicated tokens** - Create a token specifically for scanning
2. **Cache installations** - Use Go module caching
3. **Store reports** - Keep historical reports for trending
4. **Notify on changes** - Alert when compliance degrades
5. **Version policies** - Tag policy releases for reproducibility
6. **Test policy changes** - Validate before merging

## See Also

- [GitHub Actions Integration](github-actions.md) - Basic setup
- [SARIF Integration](../reports/sarif.md) - Security tab
- [Profiles](../profiles/overview.md) - Profile configuration
