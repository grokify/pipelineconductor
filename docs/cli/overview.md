# CLI Overview

PipelineConductor provides a command-line interface for scanning repositories and managing CI/CD compliance.

## Commands

| Command | Description |
|---------|-------------|
| `scan` | Scan repositories for compliance |
| `validate` | Validate Cedar policy files |
| `version` | Print version information |

## Global Flags

These flags are available for all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `$HOME/.pipelineconductor.yaml` |
| `--github-token` | GitHub personal access token | `$GITHUB_TOKEN` |
| `--orgs` | GitHub organizations to scan | (required) |
| `--policy-repo` | Policy repository (e.g., `owner/repo@ref`) | - |
| `--profile` | Profile to use for evaluation | `default` |
| `-v, --verbose` | Enable verbose output | `false` |

## Usage Pattern

```bash
pipelineconductor [command] [flags]
```

## Examples

```bash
# Scan with verbose output
pipelineconductor scan --orgs myorg -v

# Use a specific config file
pipelineconductor scan --config ./myconfig.yaml

# Validate policies
pipelineconductor validate ./policies/
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid flags, API error, etc.) |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | GitHub personal access token |
| `PIPELINECONDUCTOR_CONFIG` | Path to config file |

## Getting Help

```bash
# General help
pipelineconductor --help

# Command-specific help
pipelineconductor scan --help
pipelineconductor validate --help
```
