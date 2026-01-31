# Profiles Overview

Profiles define expected CI/CD configurations for different project types. They specify Go versions, target platforms, required checks, and other CI settings.

## What is a Profile?

A profile is a named configuration that defines:

- **Go versions** to test against
- **Operating systems** for the build matrix
- **Required checks** (test, lint, build)
- **Lint settings** (tool, enabled)
- **Test settings** (coverage, race detection)

## Built-in Profiles

PipelineConductor includes three built-in profiles:

| Profile | Go Versions | Platforms | Use Case |
|---------|-------------|-----------|----------|
| `default` | 1.24, 1.25 | Linux, macOS, Windows | Standard active projects |
| `modern` | 1.25 | Linux, macOS | Latest features, new projects |
| `legacy` | 1.12 | Linux | Older projects, compatibility |

## Using Profiles

### Specify a Profile

```bash
pipelineconductor scan --orgs myorg --profile modern
```

### Default Profile

If no profile is specified, `default` is used:

```bash
# Uses "default" profile
pipelineconductor scan --orgs myorg
```

## Profile Evaluation

When a scan runs with a profile:

1. Each repository's CI configuration is extracted
2. Go versions are compared against the profile's allowed versions
3. OS matrix is compared against the profile's required platforms
4. Violations are generated for mismatches

### Example Violations

```json
{
  "violations": [
    {
      "policy": "profile/go-version",
      "rule": "allowed-versions",
      "message": "Go version 1.20 not in profile default allowed versions [1.24, 1.25]",
      "severity": "medium",
      "remediation": "Update go-version to one of: [1.24, 1.25]"
    },
    {
      "policy": "profile/os-matrix",
      "rule": "required-os",
      "message": "Profile default requires OS windows-latest but it's not in the matrix",
      "severity": "info",
      "remediation": "Add windows-latest to your OS matrix"
    }
  ]
}
```

## Profile Selection Strategy

Choose a profile based on your project needs:

### Default Profile

Use for:

- Active production projects
- Projects that need broad platform support
- Standard compliance requirements

### Modern Profile

Use for:

- New projects using latest Go features
- Projects that don't need Windows support
- Faster CI runs (fewer matrix combinations)

### Legacy Profile

Use for:

- Older projects that can't upgrade Go
- Projects with specific compatibility requirements
- Minimal CI overhead

## Configuration in YAML

Set the default profile in your config file:

```yaml
# ~/.pipelineconductor.yaml
profile: modern
```

## Next Steps

- [Built-in Profiles](builtin.md) - Detailed profile specifications
- [Custom Profiles](custom.md) - Create your own profiles
