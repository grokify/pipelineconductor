# Configuration

PipelineConductor can be configured via command-line flags, environment variables, or a configuration file.

## Configuration File

Create a configuration file at `~/.pipelineconductor.yaml` or `.pipelineconductor.yaml` in your project:

```yaml
# GitHub authentication
github_token: ${GITHUB_TOKEN}

# Organizations to scan
orgs:
  - myorg
  - otherorg

# Profile for evaluation
profile: default

# Policy repository (optional)
policy_repo: myorg/policies@main

# Output settings
output: report.json
format: json

# Verbose logging
verbose: false
```

## Configuration Precedence

Configuration values are loaded in this order (later overrides earlier):

1. Default values
2. Configuration file
3. Environment variables
4. Command-line flags

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | GitHub personal access token |
| `PIPELINECONDUCTOR_CONFIG` | Path to config file |

## Configuration Options

### github_token

GitHub personal access token for API access.

```yaml
github_token: ghp_xxxxxxxxxxxx
```

!!! warning "Security"
    Use `${GITHUB_TOKEN}` to reference an environment variable instead of hardcoding the token.

### orgs

List of GitHub organizations to scan.

```yaml
orgs:
  - myorg
  - otherorg
  - thirdorg
```

### profile

Profile name for policy evaluation.

```yaml
profile: default  # or: modern, legacy
```

### policy_repo

Remote repository containing Cedar policies.

```yaml
policy_repo: myorg/policies@main
```

Format: `owner/repo@ref` where `ref` can be a branch, tag, or commit SHA.

### output

Default output file path.

```yaml
output: reports/compliance.json
```

### format

Default output format.

```yaml
format: json  # or: markdown, sarif, csv
```

### verbose

Enable verbose logging.

```yaml
verbose: true
```

## Multiple Configuration Files

You can use different configuration files for different environments:

```bash
# Development
pipelineconductor scan --config ./config/dev.yaml

# Production
pipelineconductor scan --config ./config/prod.yaml
```

## Example Configurations

### Minimal Configuration

```yaml
github_token: ${GITHUB_TOKEN}
orgs:
  - myorg
```

### Full Configuration

```yaml
# Authentication
github_token: ${GITHUB_TOKEN}

# Target organizations
orgs:
  - myorg
  - shared-libs
  - internal-tools

# Policy configuration
profile: modern
policy_repo: myorg/ci-policies@v1.0.0

# Output settings
output: reports/compliance.json
format: json

# Logging
verbose: false
```

### CI/CD Configuration

For use in GitHub Actions or other CI systems:

```yaml
# .pipelineconductor.ci.yaml
github_token: ${GITHUB_TOKEN}
orgs:
  - ${GITHUB_REPOSITORY_OWNER}
profile: default
format: sarif
output: results.sarif
verbose: true
```

## See Also

- [scan Command](scan.md) - Scan command reference
- [Profiles](../profiles/overview.md) - Profile configuration
