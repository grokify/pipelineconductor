# Dashboard Integration Example

This example shows how to generate a Dashforge dashboard from pipelineconductor compliance check results.

## Quick Start

### 1. Run a compliance check with dashboard output

```bash
# Generate both JSON data and dashboard definition
pipelineconductor check \
  --users grokify \
  --languages Go,TypeScript \
  -o data.json \
  --dashboard dashboard.json
```

This creates two files:
- `data.json` - The compliance check results
- `dashboard.json` - A Dashforge dashboard that visualizes the results

### 2. View the dashboard

**Option A: Local HTTP server**

```bash
# Copy the Dashforge viewer
cp -r ~/go/src/github.com/grokify/dashforge/viewer .

# Start a local server
python3 -m http.server 8080

# Open in browser
open "http://localhost:8080/viewer/?dashboard=../dashboard.json"
```

**Option B: GitHub Pages**

1. Push `data.json` and `dashboard.json` to a GitHub Pages-enabled repo
2. Access via: `https://yourusername.github.io/repo/viewer/?dashboard=../dashboard.json`

## Dashboard Features

The generated dashboard includes:

### Key Metrics (Row 1)
- **Total Repositories** - Count of scanned repos
- **Fully Compliant** - Repos using reusable workflows exactly
- **Partial Compliance** - Repos with equivalent but not reusable workflows
- **Non-Compliant** - Repos missing required workflows

### Charts (Row 2)
- **Compliance Rate** - Overall percentage metric with color thresholds
- **Compliance by Language** - Bar chart showing rates per language
- **Repositories by Language** - Pie chart of language distribution

### Tables (Rows 3-4)
- **Repositories Needing Attention** - Non-compliant and partial repos
- **All Repositories** - Complete list with pagination

## Customizing the Dashboard

The generated `dashboard.json` follows the Dashforge IR format and can be customized:

```json
{
  "widgets": [
    {
      "id": "custom-widget",
      "type": "chart",
      "position": { "x": 0, "y": 0, "w": 6, "h": 4 },
      "dataSourceId": "repos",
      "config": {
        "marks": [{ "geometry": "bar", "encode": { "x": "fullName", "y": "scanTimeMs" }}]
      }
    }
  ]
}
```

## Data URL Options

By default, the dashboard references data via a relative path. You can customize this:

```bash
# Use absolute URL (e.g., for GitHub Pages)
pipelineconductor check \
  --users grokify \
  --languages Go \
  -o data.json \
  --dashboard dashboard.json \
  --data-url "https://grokify.github.io/reports/data.json"
```

## Directory Structure for GitHub Pages

```
your-repo/
├── viewer/
│   └── index.html          # Dashforge viewer
├── dashboards/
│   └── compliance.json     # Dashboard definition
└── data/
    └── compliance-data.json # Compliance check results
```

Access via: `https://you.github.io/your-repo/viewer/?dashboard=../dashboards/compliance.json`

## Automated Updates

Use GitHub Actions to automatically update compliance data:

```yaml
name: Update Compliance Dashboard
on:
  schedule:
    - cron: '0 0 * * *'  # Daily
  workflow_dispatch:

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run compliance check
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          pipelineconductor check \
            --users grokify \
            --languages Go,TypeScript \
            -o data/compliance.json \
            --dashboard dashboards/compliance.json \
            --data-url "./data/compliance.json"

      - name: Commit changes
        run: |
          git config user.name "github-actions"
          git config user.email "github-actions@github.com"
          git add data/ dashboards/
          git diff --staged --quiet || git commit -m "chore: update compliance dashboard"
          git push
```
