# PipelineConductor Tasks

## Ready to Deploy

- [ ] **Run remediation on plexusone repos** - Fix the 30 non-compliant/partial repos identified by compliance check
- [ ] **Release v0.2.0** - Tag a release with new compliance features (check, remediate, apply commands, local scanning, GitHub Action)

## Quality & CI

- [ ] **CI/CD for pipelineconductor** - Set up GitHub Actions using the reusable workflows we're enforcing
- [ ] **Integration tests** - End-to-end tests for CLI commands (check, remediate, apply)
- [ ] **Test coverage badge** - Add coverage reporting with gocoverbadge

## Features

- [ ] **Diff output** - Show what remediation would change (`--diff` flag)
- [ ] **TypeScript verification** - Test TypeScript compliance rules with real repos

## Operational

- [ ] **Slack/Teams notifications** - Alert on compliance changes via webhooks
- [ ] **Scheduled compliance reports** - Cron job examples for weekly compliance reports

## Completed

- [x] Local filesystem scanning (LocalCollector)
- [x] Compliance checking with filename mismatch detection
- [x] Remediation generator for workflow files
- [x] Batch apply with git commit/push/PR support
- [x] gitscan integration
- [x] GitHub Action (action.yaml)
- [x] MkDocs documentation for new commands
- [x] Unit tests for LocalCollector, WorkflowMatcher, Generator
- [x] Cedar policy integration (--policies, --policy-action, --fail-on-deny flags)
- [x] HTML report format (--format html with interactive filtering)
