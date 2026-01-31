# PipelineConductor - Product Requirements Document

## Executive Summary

PipelineConductor is a multi-repo CI/CD pipeline orchestration and compliance system designed to manage 600+ repositories across multiple GitHub organizations. It provides centralized policy enforcement, automated compliance scanning, and remediation capabilities through automated pull requests.

## Problem Statement

Organizations with hundreds of repositories face significant challenges:

- **Inconsistent CI/CD**: Each repo has snowflake pipelines with duplicated logic
- **Policy drift**: Security and compliance requirements are inconsistently applied
- **Manual overhead**: Keeping pipelines synchronized requires constant manual effort
- **Scaling pain**: Adding new languages, tools, or requirements means touching every repo
- **Audit complexity**: No centralized view of compliance status across repos

## Solution

PipelineConductor implements a three-layer architecture:

1. **Policy & Orchestration (org-wide)**: Centralized policy definitions and compliance rules
2. **Reusable Workflows (language/stack specific)**: Shared CI/CD building blocks
3. **Thin Repo-Level Config**: Minimal YAML (~15 lines) with no logic

### Core Principle

> "Repos describe what they are. Org CI decides how they're built. Security lives centrally."

If repos contain CI logic, the architecture has failed.

## Target Users

| User | Needs |
|------|-------|
| Platform Engineers | Manage CI/CD at scale, enforce consistency |
| Security Teams | Ensure compliance, audit policy adherence |
| Developers | Simple onboarding, minimal CI maintenance |
| DevOps/SRE | Reliable builds, predictable pipelines |

## Functional Requirements

### FR-1: Multi-Repo Compliance Scanning

- Scan repositories across multiple GitHub organizations
- Support scheduled scans (weekly/monthly) and on-demand execution
- API-first discovery with selective git inspection
- Generate compliance reports (JSON, SARIF, Markdown, CSV)

### FR-2: Policy-as-Code with Cedar

- Define CI/CD policies using Cedar policy language
- Support policy profiles (default, modern, legacy)
- Evaluate policies against collected repository context
- Provide clear allow/deny decisions with explanations

### FR-3: Automated Remediation

- Generate pull requests to fix policy violations
- Support dry-run mode for preview
- Batch remediation across multiple repos
- Track remediation status and outcomes

### FR-4: Go CI/CD Support (Initial Focus)

- Multi-platform builds: macOS, Ubuntu, Windows
- Multi-version Go matrix: 1.24, 1.25 (configurable)
- Legacy Go support (1.12+) via profiles
- Integration with golangci-lint, go test, go build

### FR-5: Dependency Management Integration

- Support Mend Renovate for consolidated PRs
- Policy-based auto-merge (tests pass, dependency age threshold)
- Vulnerability scanning integration
- Dependency freshness tracking

### FR-6: Reporting and Visibility

- Compliance dashboard data (JSON output)
- SARIF output for GitHub Security tab integration
- Historical tracking of compliance status
- Exception/waiver tracking

## Non-Functional Requirements

### NFR-1: Scalability

- Handle 600+ repositories
- Support multiple GitHub organizations
- Efficient API usage with rate limiting awareness

### NFR-2: Extensibility

- Language-agnostic core architecture
- Pluggable collectors (GitHub, GitLab)
- Extensible policy framework

### NFR-3: Security

- Read-only operations by default
- Explicit opt-in for remediation PRs
- No secrets in policy definitions
- Audit trail for all actions

### NFR-4: Usability

- Simple CLI interface
- GitHub Actions integration
- Minimal configuration for repos
- Clear error messages and remediation guidance

## Success Metrics

| Metric | Target |
|--------|--------|
| Repos with compliant CI | >95% |
| Average policy evaluation time | <30s per repo |
| Automated remediation success rate | >90% |
| Time to onboard new repo | <5 minutes |
| Policy change rollout time | <1 hour org-wide |

## Constraints

- Must work with GitHub Actions (GitLab support later)
- Cedar policy language (not OPA/Rego for CI/CD)
- Go as primary implementation language
- Open source platform, policies can be public or private

## Out of Scope (Initial Release)

- Kubernetes/infrastructure policy (use OPA for that)
- Real-time pipeline blocking (advisory mode first)
- Non-GitHub platforms (GitLab in future)
- Languages other than Go (Phase 2+)

## Dependencies

- [cedar-go](https://github.com/cedar-policy/cedar-go) - Cedar policy evaluation
- GitHub API v4 (GraphQL) and v3 (REST)
- GitHub Actions for execution environment

## Glossary

| Term | Definition |
|------|------------|
| Profile | Named configuration for CI matrix (Go versions, OS platforms) |
| Policy | Cedar rule defining allowed/denied CI actions |
| Collector | Component that gathers repo metadata from APIs |
| Remediator | Component that generates fix PRs |
| Context | Structured data about a repo passed to Cedar for evaluation |
