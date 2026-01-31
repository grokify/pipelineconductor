# PipelineConductor

[![CI](https://github.com/grokify/pipelineconductor/workflows/CI/badge.svg)](https://github.com/grokify/pipelineconductor/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/pipelineconductor)](https://goreportcard.com/report/github.com/grokify/pipelineconductor)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Orchestrate and harmonize multi-repo CI/CD pipelines with policy-driven automation.**

PipelineConductor is a tool for managing CI/CD pipeline consistency across hundreds of repositories. It scans repositories, evaluates them against Cedar policies, generates compliance reports, and can automatically remediate violations via pull requests.

## Features

- **Multi-org scanning**: Scan repositories across multiple GitHub organizations
- **Policy-as-code**: Define CI/CD policies using [Cedar](https://www.cedarpolicy.com/)
- **Profile system**: Named configurations for different project types (default, modern, legacy)
- **Compliance reports**: Generate JSON, SARIF, Markdown, and CSV reports
- **Automated remediation**: Create PRs to fix policy violations
- **API-first**: Efficient GitHub API usage with selective git inspection

## Installation

```bash
go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/pipelineconductor.git
cd pipelineconductor
go build -o pipelineconductor ./cmd/pipelineconductor
```

## Quick Start

1. Set your GitHub token:

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

2. Scan your organization:

```bash
pipelineconductor scan --orgs myorg --output report.json
```

3. View the compliance report:

```bash
pipelineconductor scan --orgs myorg --format markdown
```

## Usage

### Scan Command

Scan repositories for compliance:

```bash
# Basic scan
pipelineconductor scan --orgs myorg

# Multiple organizations
pipelineconductor scan --orgs org1,org2,org3

# Filter by language
pipelineconductor scan --orgs myorg --languages Go,Python

# Include archived repos
pipelineconductor scan --orgs myorg --include-archived

# Output to file
pipelineconductor scan --orgs myorg --output report.json --format json
```

### Configuration File

Create `~/.pipelineconductor.yaml` or `.pipelineconductor.yaml`:

```yaml
github_token: ${GITHUB_TOKEN}
orgs:
  - myorg
  - otherorg
profile: default
verbose: true
```

## Profiles

PipelineConductor uses profiles to define expected CI/CD configurations:

| Profile | Go Versions | Platforms | Use Case |
|---------|-------------|-----------|----------|
| `default` | 1.24, 1.25 | Linux, macOS, Windows | Standard projects |
| `modern` | 1.25 | Linux, macOS | Latest features |
| `legacy` | 1.12 | Linux | Older projects |

## Documentation

- [PRD.md](PRD.md) - Product Requirements Document
- [TRD.md](TRD.md) - Technical Requirements Document
- [MRD.md](MRD.md) - Market Requirements Document
- [ROADMAP.md](ROADMAP.md) - Implementation Roadmap

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      PipelineConductor CLI                       │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Collectors  │  │    Policy    │  │     Remediator       │  │
│  │ - GitHub API │  │    Engine    │  │ - PR Generator       │  │
│  │ - GitLab API │  │ - Cedar      │  │ - Patch Builder      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                            │                                     │
│                    ┌───────┴────────┐                           │
│                    │   pkg/model    │                           │
│                    └────────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
