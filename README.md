# PipelineConductor

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/pipelineconductor/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/pipelineconductor
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/pipelineconductor
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/pipelineconductor
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/pipelineconductor
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/pipelineconductor
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fpipelineconductor
 [loc-svg]: https://tokei.rs/b1/github/plexusone/pipelineconductor
 [repo-url]: https://github.com/plexusone/pipelineconductor
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/pipelineconductor/blob/main/LICENSE

**Orchestrate and harmonize multi-repo CI/CD pipelines with policy-driven automation.**

PipelineConductor is a tool for managing CI/CD pipeline consistency across hundreds of repositories. It scans repositories, evaluates them against Cedar policies, generates compliance reports, and can automatically remediate violations via pull requests.

## Features

- 🏢 **Multi-org scanning**: Scan repositories across multiple GitHub organizations
- 📜 **Policy-as-code**: Define CI/CD policies using [Cedar](https://www.cedarpolicy.com/)
- ⚙️ **Profile system**: Named configurations for different project types (default, modern, legacy)
- 📊 **Compliance reports**: Generate JSON, SARIF, Markdown, and CSV reports
- 🔧 **Automated remediation**: Create PRs to fix policy violations
- ⚡ **API-first**: Efficient GitHub API usage with selective git inspection

## Installation

```bash
go install github.com/plexusone/pipelineconductor/cmd/pipelineconductor@latest
```

Or build from source:

```bash
git clone https://github.com/plexusone/pipelineconductor.git
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
┌────────────────────────────────────────────────────────────────┐
│                      PipelineConductor CLI                     │
├────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Collectors  │  │    Policy    │  │     Remediator       │  │
│  │ - GitHub API │  │    Engine    │  │ - PR Generator       │  │
│  │ - GitLab API │  │ - Cedar      │  │ - Patch Builder      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                            │                                   │
│                    ┌───────┴────────┐                          │
│                    │   pkg/model    │                          │
│                    └────────────────┘                          │
└────────────────────────────────────────────────────────────────┘
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
