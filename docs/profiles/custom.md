# Custom Profiles

Create custom profiles to match your organization's specific CI/CD requirements.

## Profile Schema

Profiles are defined in YAML:

```yaml
name: my-profile
description: Custom profile for my organization

go:
  versions:
    - "1.24"
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

## Creating a Profile

### 1. Create Profile Directory

```bash
mkdir -p configs/profiles
```

### 2. Create Profile File

Create `configs/profiles/enterprise.yaml`:

```yaml
name: enterprise
description: Enterprise Go CI with strict requirements

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
    - security-scan

lint:
  enabled: true
  tool: golangci-lint

test:
  coverage: true
  race: true
```

### 3. Use the Profile

```bash
pipelineconductor scan --orgs myorg --profile enterprise
```

## Profile Fields

### name (required)

Unique identifier for the profile.

```yaml
name: my-custom-profile
```

### description

Human-readable description.

```yaml
description: CI profile for microservices
```

### go.versions

List of allowed Go versions.

```yaml
go:
  versions:
    - "1.24"
    - "1.25"
```

### os

List of required operating systems in the build matrix.

```yaml
os:
  - ubuntu-latest
  - macos-latest
  - windows-latest
```

### checks.required

List of required CI checks.

```yaml
checks:
  required:
    - test
    - lint
    - build
```

### lint

Linting configuration.

```yaml
lint:
  enabled: true
  tool: golangci-lint
```

### test

Test configuration.

```yaml
test:
  coverage: true
  race: true
```

## Example Profiles

### Microservices Profile

```yaml
name: microservices
description: CI for containerized microservices

go:
  versions:
    - "1.25"

os:
  - ubuntu-latest

checks:
  required:
    - test
    - lint
    - build
    - docker-build

lint:
  enabled: true
  tool: golangci-lint

test:
  coverage: true
  race: true
```

### Library Profile

```yaml
name: library
description: CI for reusable Go libraries

go:
  versions:
    - "1.23"
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

### Minimal Profile

```yaml
name: minimal
description: Minimal CI for internal tools

go:
  versions:
    - "1.25"

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

## Loading Custom Profiles

Custom profiles are loaded automatically when placed in `configs/profiles/`.

The profile manager searches for:

1. Built-in profiles (`default`, `modern`, `legacy`)
2. Profiles in `configs/profiles/*.yaml`

## Profile Inheritance (Future)

!!! note "Coming Soon"
    Profile inheritance will allow extending base profiles.

```yaml
# Future syntax
name: strict
extends: default
description: Stricter version of default

go:
  versions:
    - "1.25"  # Override: only latest
```

## Validation

Profiles are validated when loaded. Invalid profiles will cause an error:

```bash
pipelineconductor scan --orgs myorg --profile invalid-profile
```

Output:

```
Error: profile not found: invalid-profile
```

## Best Practices

1. **Start from built-in** - Copy a built-in profile and modify
2. **Be specific** - Name profiles by purpose (`microservices`, `library`)
3. **Document requirements** - Use the description field
4. **Version control** - Store profiles in your repo
5. **Test changes** - Validate profiles before deploying

## See Also

- [Built-in Profiles](builtin.md) - Default profiles
- [Configuration](../cli/config.md) - Set default profile
