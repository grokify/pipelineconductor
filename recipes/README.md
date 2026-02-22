# Workflow Recipes

This directory contains workflow recipes for different language stacks. Each recipe is a thin wrapper that calls reusable workflows from `grokify/.github`.

## Architecture

```
grokify/.github/                    # Org-level repo
├── .github/workflows/              # Reusable workflows (logic)
│   ├── go-ci.yaml                  # Go CI with inputs
│   ├── go-lint.yaml                # Go lint with inputs
│   ├── ts-ci.yaml                  # TypeScript CI with inputs
│   └── ts-lint.yaml                # TypeScript lint with inputs
└── workflow-templates/             # GitHub starter templates
    ├── go-ci.yaml
    └── ts-ci.yaml

pipelineconductor/recipes/          # This directory
├── go/                             # Go recipes
│   ├── go-ci.yaml
│   └── go-lint.yaml
├── ts/                             # TypeScript recipes
│   ├── ts-ci.yaml
│   └── ts-lint.yaml
└── multi-language/                 # Multi-language examples
    └── go-ts.yaml
```

## Naming Convention

```
<language>-<action>.yaml
```

| Pattern | Examples |
|---------|----------|
| `go-*.yaml` | `go-ci.yaml`, `go-lint.yaml`, `go-codeql.yaml` |
| `ts-*.yaml` | `ts-ci.yaml`, `ts-lint.yaml` |
| `swift-*.yaml` | `swift-ci.yaml`, `swift-lint.yaml` |
| `python-*.yaml` | `python-ci.yaml`, `python-lint.yaml` |

## Usage

### Single-Language Repo (Go)

Copy to `.github/workflows/`:

```bash
cp recipes/go/go-ci.yaml .github/workflows/
cp recipes/go/go-lint.yaml .github/workflows/
```

### Single-Language Repo (TypeScript)

Copy to `.github/workflows/`:

```bash
cp recipes/ts/ts-ci.yaml .github/workflows/
cp recipes/ts/ts-lint.yaml .github/workflows/
```

### Multi-Language Repo (Go + TypeScript)

Copy all relevant recipes:

```bash
cp recipes/go/go-ci.yaml .github/workflows/
cp recipes/go/go-lint.yaml .github/workflows/
cp recipes/ts/ts-ci.yaml .github/workflows/
cp recipes/ts/ts-lint.yaml .github/workflows/
```

If TypeScript is in a subdirectory, add `working-directory`:

```yaml
jobs:
  ci:
    uses: grokify/.github/.github/workflows/ts-ci.yaml@main
    with:
      working-directory: 'ts'
```

## Pipelineconductor Detection

Pipelineconductor can identify repo languages by workflow files:

| Query | Glob Pattern |
|-------|--------------|
| All Go repos | `go-*.yaml` |
| All TypeScript repos | `ts-*.yaml` |
| Multi-language repos | Has both `go-*.yaml` AND `ts-*.yaml` |
| Repos missing lint | Has `*-ci.yaml` but no `*-lint.yaml` |

## Customization

Recipes use sensible defaults. Override via `with:` inputs:

```yaml
jobs:
  ci:
    uses: grokify/.github/.github/workflows/go-ci.yaml@main
    with:
      go-versions: '["1.25.x"]'           # Test single version
      platforms: '["ubuntu-latest"]'       # Linux only
      test-flags: '-v -race -covermode=atomic'
```

See reusable workflow files for all available inputs.

## Path Filtering

All recipes include path filters so CI only runs when relevant files change:

- Go workflows trigger on `**.go`, `go.mod`, `go.sum`
- TypeScript workflows trigger on `**.ts`, `package.json`, `package-lock.json`

This saves CI minutes in multi-language repos.
