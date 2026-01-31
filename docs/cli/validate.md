# validate Command

The `validate` command validates Cedar policy syntax and loads policies from a directory or file.

## Synopsis

```bash
pipelineconductor validate [path] [flags]
```

## Description

The validate command checks that Cedar policy files are syntactically correct without running a full scan. This is useful for:

- Testing policies during development
- CI/CD pipeline validation
- Verifying policy syntax before deployment

## Arguments

| Argument | Description |
|----------|-------------|
| `path` | Path to a policy file or directory (optional if `--builtin` is used) |

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--builtin` | Validate built-in policies | `false` |
| `--verbose` | Show policy details | `false` |

## Examples

### Validate Built-in Policies

```bash
pipelineconductor validate --builtin
```

Output:

```
Validating built-in policies...
✓ Built-in policies are valid
```

### Validate a Directory

```bash
pipelineconductor validate ./policies/
```

Output:

```
Validating policies in directory: ./policies/

✓ 5 policy file(s) validated successfully
```

### Validate with Verbose Output

```bash
pipelineconductor validate ./policies/ --verbose
```

Output:

```
Validating policies in directory: ./policies/
  ✓ ./policies/go/merge.cedar
  ✓ ./policies/go/versions.cedar
  ✓ ./policies/go/matrix.cedar
  ✓ ./policies/go/dependencies.cedar
  ✓ ./policies/go/reusable-workflow.cedar

✓ 5 policy file(s) validated successfully
```

### Validate a Single File

```bash
pipelineconductor validate ./policies/go/merge.cedar
```

Output:

```
Validating policy file: ./policies/go/merge.cedar
✓ 1 policy file(s) validated successfully
```

### Handling Errors

When a policy has syntax errors:

```bash
pipelineconductor validate ./policies/
```

Output:

```
Validating policies in directory: ./policies/

Errors found:
  ✗ ./policies/broken.cedar: parsing error at line 5

1 policy file(s) have errors
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All policies are valid |
| 1 | One or more policies have errors |

## Use in CI/CD

Add policy validation to your CI pipeline:

```yaml
# .github/workflows/validate-policies.yml
name: Validate Policies

on:
  push:
    paths:
      - 'policies/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install PipelineConductor
        run: go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest

      - name: Validate policies
        run: pipelineconductor validate ./policies/ --verbose
```

## See Also

- [scan](scan.md) - Run compliance scans
- [Cedar Syntax](../policies/cedar-syntax.md) - Cedar policy language reference
- [Writing Policies](../policies/writing.md) - Policy authoring guide
