# Roadmap

This document outlines the development roadmap for PipelineConductor.

## Current Status

**Version:** 0.1.0-dev

### Completed Features

- [x] Multi-org repository scanning
- [x] Cedar policy evaluation engine
- [x] Profile system (default, modern, legacy)
- [x] Multiple report formats (JSON, Markdown, SARIF, CSV)
- [x] Policy validation command
- [x] GitHub API integration
- [x] Built-in policies
- [x] Example Cedar policies

### Recently Completed

- [x] API rate limiting handling (with exponential backoff and GitHub X-RateLimit header support)
- [x] GitHub Actions CI/CD workflows (ci.yaml, lint.yaml, sast_codeql.yaml)

### Backlog

- [ ] Policy repository support (`--policy-repo`) - load policies from remote git repo
- [ ] Profile inheritance - extend base profiles with overrides
- [ ] README.md update - link to documentation site
- [ ] Integration tests - test against real GitHub API

## Upcoming Releases

### v0.1.0 - Initial Release

**Target:** Q1 2025

- Core scanning functionality
- Cedar policy evaluation
- Report generation (JSON, Markdown, SARIF, CSV)
- Profile system
- Documentation

### v0.2.0 - Remediation

**Target:** Q2 2025

- Automated PR creation for violations
- Workflow file generation
- Dry-run mode
- Batch remediation

### v0.3.0 - Enterprise Features

**Target:** Q3 2025

- GitLab support
- Bitbucket support
- Trend tracking
- Dashboard integration
- Policy repository sync

## Feature Requests

Features under consideration:

### Multi-SCM Support

- GitLab CI/CD configuration parsing
- Bitbucket Pipelines support
- Azure DevOps integration

### Advanced Policies

- Policy templates
- Policy versioning
- Policy inheritance
- Custom severity levels

### Reporting

- HTML reports
- Trend charts
- Email notifications
- Webhook integrations

### Performance

- Parallel scanning
- Caching
- Incremental scans

## Contributing

See [Contributing](contributing.md) for how to help with development.

## Feedback

We welcome feedback on the roadmap:

- File an issue on GitHub
- Join discussions
- Submit feature requests

## Version History

| Version | Date | Highlights |
|---------|------|------------|
| 0.1.0-dev | Current | Initial development |

See [CHANGELOG](https://github.com/grokify/pipelineconductor/blob/main/CHANGELOG.md) for detailed release notes.
