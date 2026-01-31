# PipelineConductor - Roadmap

## Overview

This roadmap outlines the phased implementation of PipelineConductor, starting with Go-only support and expanding to multi-language, multi-platform capabilities.

## Phase 1: Foundation (Weeks 1-4)

**Goal**: Establish core architecture and basic Go compliance scanning.

### Milestone 1.1: Project Setup
- [x] Create repository structure
- [x] Initialize Go module
- [x] Create PRD.md, TRD.md, ROADMAP.md
- [x] Set up CI workflow for pipelineconductor itself
- [x] Add golangci-lint configuration
- [x] Create initial README.md

### Milestone 1.2: Core Models
- [x] Define `pkg/model/repo.go` - Repository model
- [x] Define `pkg/model/workflow.go` - CI workflow model
- [x] Define `pkg/model/policy.go` - Policy context model
- [x] Define `pkg/model/result.go` - Compliance result model
- [x] Add JSON serialization support
- [x] Write unit tests for models

### Milestone 1.3: GitHub Collector
- [x] Implement `internal/collector/collector.go` - Interface definition
- [x] Implement `internal/collector/github.go` - GitHub API integration
- [x] List repositories across organizations
- [x] Fetch repository metadata (languages, topics, archived status)
- [x] Fetch workflow files from `.github/workflows/`
- [x] Fetch branch protection rules
- [ ] Handle API rate limiting
- [ ] Write integration tests with mocked responses

### Milestone 1.4: Basic CLI
- [x] Set up cobra CLI framework in `cmd/pipelineconductor/`
- [x] Implement `scan` command (list repos, fetch metadata)
- [x] Implement `--orgs` flag for organization filtering
- [x] Implement `--output` flag for report path
- [x] Add configuration file support (viper)
- [ ] Write CLI tests

## Phase 2: Policy Engine (Weeks 5-8)

**Goal**: Integrate Cedar policy evaluation and define initial Go policies.

### Milestone 2.1: Cedar Integration
- [x] Add cedar-go dependency
- [x] Implement `internal/policy/engine.go` - Policy evaluation
- [x] Implement `internal/policy/loader.go` - Load policies from files/repo
- [x] Implement `internal/policy/context.go` - Build Cedar context from repo data
- [x] Write unit tests with sample policies

### Milestone 2.2: Profile System
- [x] Define profile schema (YAML)
- [x] Create `configs/profiles/default.yaml`
- [x] Create `configs/profiles/modern.yaml`
- [x] Create `configs/profiles/legacy.yaml`
- [x] Implement `internal/policy/profiles.go` - Profile loading and management
- [ ] Support profile inheritance/composition

### Milestone 2.3: Initial Go Policies
- [x] Create `policies/examples/go/merge.cedar` - Merge gating
- [x] Create `policies/examples/go/versions.cedar` - Go version enforcement
- [x] Create `policies/examples/go/matrix.cedar` - OS matrix enforcement
- [x] Create `policies/examples/go/dependencies.cedar` - Dependency age/security
- [x] Create `policies/examples/go/reusable-workflow.cedar` - Reusable workflow policy
- [ ] Write policy test cases
- [ ] Document policy syntax and examples

### Milestone 2.4: CLI Policy Commands
- [x] Implement `validate` command - Validate policy syntax
- [ ] Implement `--policy-repo` flag - Load policies from GitHub
- [x] Implement `--profile` flag - Select evaluation profile
- [x] Add policy evaluation to `scan` command
- [x] Output compliance results per repo

## Phase 3: Reporting (Weeks 9-10)

**Goal**: Generate comprehensive compliance reports in multiple formats.

### Milestone 3.1: Report Generation
- [x] Implement `internal/report/report.go` - Report builder
- [x] Implement `internal/report/json.go` - JSON output
- [x] Implement `internal/report/markdown.go` - Markdown output
- [x] Implement `internal/report/sarif.go` - SARIF output for GitHub Security
- [x] Implement `internal/report/csv.go` - CSV output for spreadsheets

### Milestone 3.2: Report Features
- [x] Summary statistics (total, compliant, non-compliant, rate)
- [x] Per-repo violation details
- [ ] Trend tracking (compare with previous scan)
- [ ] Filterable output (by org, language, compliance status)
- [x] Implement `--format` flag for output format selection

### Milestone 3.3: GitHub Actions Integration
- [ ] Create `.github/workflows/compliance.yml` for scheduled scans
- [ ] Upload reports as artifacts
- [ ] Post summary to PR comments (for PR-triggered scans)
- [ ] Integrate SARIF with GitHub Security tab

## Phase 4: Remediation (Weeks 11-14)

**Goal**: Automated PR creation to fix compliance violations.

### Milestone 4.1: Patch Generation
- [ ] Implement `internal/remediator/patch.go` - Generate file patches
- [ ] Create workflow file from template
- [ ] Update existing workflow file (Go versions, matrix)
- [ ] Update reusable workflow references

### Milestone 4.2: PR Creation
- [ ] Implement `internal/remediator/pr.go` - GitHub PR creation
- [ ] Create feature branch
- [ ] Commit patches
- [ ] Create PR with descriptive body
- [ ] Add labels and assignees

### Milestone 4.3: Remediation CLI
- [ ] Implement `remediate` command
- [ ] Implement `--dry-run` flag (preview changes without creating PRs)
- [ ] Implement `--batch` flag (batch PRs across repos)
- [ ] Track remediation status
- [ ] Write integration tests

### Milestone 4.4: Safety Features
- [ ] Require explicit opt-in for remediation
- [ ] Support allowlist/blocklist for repos
- [ ] Limit concurrent PRs per org
- [ ] Implement rollback capability

## Phase 5: Production Hardening (Weeks 15-16)

**Goal**: Production-ready release with documentation and testing.

### Milestone 5.1: Testing & Quality
- [ ] Achieve >80% code coverage
- [ ] Run golangci-lint with strict configuration
- [ ] Add E2E tests against test organization
- [ ] Performance testing with 100+ repos
- [ ] Security review of token handling

### Milestone 5.2: Documentation
- [ ] Comprehensive README.md with quick start
- [ ] CLI usage documentation
- [ ] Policy authoring guide
- [ ] Profile configuration guide
- [ ] GitHub Actions integration guide
- [ ] Troubleshooting guide

### Milestone 5.3: Release
- [ ] Set up GoReleaser for binary distribution
- [ ] Create GitHub Action for pipelineconductor
- [ ] Publish to GitHub Marketplace (Action)
- [ ] Tag v0.1.0 release
- [ ] Announce on relevant channels

## Phase 6: Multi-Language Support (Future)

**Goal**: Extend beyond Go to support additional languages.

### Milestone 6.1: Language Abstraction
- [ ] Abstract language-specific logic into plugins
- [ ] Define language plugin interface
- [ ] Refactor Go support as first plugin

### Milestone 6.2: Additional Languages
- [ ] Swift support (macOS-focused profiles)
- [ ] Python support (pyproject.toml, multiple Python versions)
- [ ] TypeScript/Node support (LTS versions, npm/yarn/pnpm)
- [ ] Polyglot repo support (Go + TS, Go + Python, etc.)

### Milestone 6.3: Language-Specific Policies
- [ ] Swift version policies
- [ ] Python version and dependency policies
- [ ] Node.js LTS policies
- [ ] Shared security policies across languages

## Phase 7: Multi-SCM Library & Advanced Features (Future)

**Goal**: Extract multi-SCM abstraction library and add enterprise features.

### Milestone 7.1: Multi-SCM Provider Support
- [ ] Implement `internal/collector/gitlab.go` - GitLab API integration
- [ ] Implement `internal/collector/bitbucket.go` - Bitbucket API integration
- [ ] Support GitLab CI/CD configuration parsing
- [ ] Support Bitbucket Pipelines configuration parsing
- [ ] Support GitLab merge request creation
- [ ] Support Bitbucket pull request creation
- [ ] Unified error handling across providers

### Milestone 7.2: OmniSCM Library Extraction
- [ ] Stabilize Collector interface with 3+ providers (GitHub, GitLab, Bitbucket)
- [ ] Create `github.com/grokify/omniscm` repository
- [ ] Extract and refactor collector code to new library
- [ ] Define provider-agnostic models (Repo, PR, Branch, Workflow)
- [ ] Update PipelineConductor to use external dependency
- [ ] Document library API and usage patterns

### Milestone 7.3: Dependency Automation
- [ ] Deep Mend Renovate integration
- [ ] Policy-based auto-merge decisions
- [ ] Dependency age tracking and alerting
- [ ] Vulnerability remediation automation

### Milestone 7.4: Enterprise Features
- [ ] Multi-tenant policy management
- [ ] Role-based access control for policies
- [ ] Compliance audit exports
- [ ] Integration with ticketing systems (Jira, Linear)

### Milestone 7.5: Observability
- [ ] Metrics export (Prometheus format)
- [ ] Distributed tracing support
- [ ] Health check endpoints
- [ ] Alerting integration

## Release Schedule

| Version | Target | Scope |
|---------|--------|-------|
| v0.1.0 | Week 10 | Go scanning + reporting |
| v0.2.0 | Week 14 | Automated remediation |
| v0.3.0 | Week 16 | Production hardening |
| v0.4.0 | TBD | Multi-language support |
| v1.0.0 | TBD | Stable API, enterprise features |

## Success Criteria for v0.1.0

- [ ] Scan 100+ Go repos across 2+ orgs
- [ ] Evaluate 3+ Cedar policies per repo
- [ ] Generate JSON and Markdown reports
- [ ] <30s average scan time per repo
- [ ] Documentation for basic usage
- [ ] CI passing with >80% coverage

## Dependencies and Blockers

| Dependency | Status | Notes |
|------------|--------|-------|
| cedar-go library | Available | github.com/cedar-policy/cedar-go |
| GitHub API access | Required | Need org-level token for scanning |
| Test organization | Required | Need repos for E2E testing |
| Policy repo | To create | Separate repo for org policies |

## Shared Libraries

PipelineConductor leverages and contributes to shared libraries in the grokify ecosystem:

### github.com/grokify/gogithub

Extended GitHub functionality. Use existing functions when available; add new functionality there for reuse across projects.

**Available packages:**

| Package | Purpose |
|---------|---------|
| `auth` | GitHub authentication helpers |
| `config` | Configuration management |
| `repo` | Repository operations (listing, metadata) |
| `pr` | Pull request operations |
| `release` | Release management |
| `search` | Search API (issues, PRs, users) |
| `graphql` | GraphQL client and queries |
| `errors` | GitHub-specific error handling |

**Strategy:**

1. Check gogithub first for existing functionality
2. Implement internally in PipelineConductor if specific to CI/CD scanning
3. Migrate reusable GitHub utilities to gogithub
4. Keep PipelineConductor focused on CI/CD policy orchestration

### github.com/grokify/omniscm (Future)

Multi-SCM abstraction library to be extracted from PipelineConductor after implementing 3+ providers.

**Planned packages:**

| Package | Purpose |
|---------|---------|
| `provider` | Provider interface and registry |
| `github` | GitHub implementation |
| `gitlab` | GitLab implementation |
| `bitbucket` | Bitbucket implementation |
| `model` | Provider-agnostic models |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:
- Code style and linting
- Testing requirements
- PR process
- Policy contribution guidelines
