# Built-in Profiles

PipelineConductor includes three built-in profiles for common use cases.

## Default Profile

The standard profile for active Go projects.

```yaml
name: default
description: Standard Go CI configuration for active projects

go:
  versions:
    - "1.24"
    - "1.25"

os:
  - ubuntu-latest
  - macos-latest
  - windows-latest

checks:
  required:
    - test
    - lint
    - build

lint:
  enabled: true
  tool: golangci-lint

test:
  coverage: true
  race: true
```

### When to Use

- Active production projects
- Projects that need cross-platform support
- Standard compliance requirements

### Requirements

| Requirement | Value |
|-------------|-------|
| Go versions | 1.24 or 1.25 |
| Platforms | Linux, macOS, Windows |
| Required checks | test, lint, build |
| Linting | Enabled (golangci-lint) |
| Coverage | Required |
| Race detection | Required |

## Modern Profile

For projects using the latest Go features.

```yaml
name: modern
description: Modern Go CI for projects using latest Go features

go:
  versions:
    - "1.25"

os:
  - ubuntu-latest
  - macos-latest

checks:
  required:
    - test
    - lint
    - build

lint:
  enabled: true
  tool: golangci-lint

test:
  coverage: true
  race: true
```

### When to Use

- New projects starting fresh
- Projects using Go 1.25+ features
- Teams that don't need Windows support

### Requirements

| Requirement | Value |
|-------------|-------|
| Go versions | 1.25 only |
| Platforms | Linux, macOS |
| Required checks | test, lint, build |
| Linting | Enabled (golangci-lint) |
| Coverage | Required |
| Race detection | Required |

## Legacy Profile

For older projects with compatibility constraints.

```yaml
name: legacy
description: Legacy Go CI for older projects requiring Go 1.12-1.18

go:
  versions:
    - "1.12"

os:
  - ubuntu-latest

checks:
  required:
    - test
    - build

lint:
  enabled: false

test:
  coverage: false
  race: false
```

### When to Use

- Older projects that can't upgrade Go
- Projects with specific compatibility requirements
- Minimal CI overhead

### Requirements

| Requirement | Value |
|-------------|-------|
| Go versions | 1.12 |
| Platforms | Linux only |
| Required checks | test, build |
| Linting | Disabled |
| Coverage | Not required |
| Race detection | Not required |

## Profile Comparison

| Feature | Default | Modern | Legacy |
|---------|---------|--------|--------|
| Go versions | 1.24, 1.25 | 1.25 | 1.12 |
| Linux | ✅ | ✅ | ✅ |
| macOS | ✅ | ✅ | ❌ |
| Windows | ✅ | ❌ | ❌ |
| Linting | ✅ | ✅ | ❌ |
| Coverage | ✅ | ✅ | ❌ |
| Race detection | ✅ | ✅ | ❌ |

## Viewing Profile Details

Check which profile is in use:

```bash
pipelineconductor scan --orgs myorg --profile modern -v
```

Output includes:

```
Using profile: modern
```

## See Also

- [Custom Profiles](custom.md) - Create your own profiles
- [Configuration](../cli/config.md) - Set default profile
