# Writing Policies

This guide walks through creating custom Cedar policies for PipelineConductor.

## Getting Started

### 1. Create a Policy Directory

```bash
mkdir -p policies/go
```

### 2. Create a Policy File

Create `policies/go/versions.cedar`:

```cedar
// Require supported Go versions
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.goVersions.containsAny(["1.18", "1.19", "1.20"])
};
```

### 3. Validate the Policy

```bash
pipelineconductor validate policies/
```

### 4. Test with a Scan

```bash
pipelineconductor scan --orgs myorg --policy-dir policies/
```

## Policy Patterns

### Require Something

Use `forbid` when a required condition is missing:

```cedar
// Require CI workflow
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == false
};
```

### Block Something

Use `forbid` when a condition should prevent action:

```cedar
// Block repos with vulnerabilities
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasVulnerabilities == true
};
```

### Allow with Conditions

Use `permit` to explicitly allow when conditions are met:

```cedar
// Allow merge when using dependency automation
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    (context.hasRenovate == true || context.hasDependabot == true) &&
    context.hasVulnerabilities == false
};
```

### Tiered Requirements

Create multiple policies for graduated enforcement:

```cedar
// Strict: Require all platforms
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.containsAll(["ubuntu-latest", "macos-latest", "windows-latest"])
};

// Relaxed: Allow Linux-only for legacy
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.contains("ubuntu-latest") &&
    context.topics.contains("legacy")
};
```

## Organization-Specific Policies

### By Organization

```cedar
// Stricter rules for production org
forbid(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.repoOrg == "production-org" &&
    context.branchProtectionEnabled == false
};
```

### By Topic

```cedar
// Require security scanning for public repos
forbid(
    principal,
    action == Action::"release",
    resource
)
when {
    context.topics.contains("public") &&
    context.hasVulnerabilities == true
};
```

### By Language

```cedar
// Go-specific requirements
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.languages.contains("Go") &&
    context.hasGoMod == false
};
```

## Policy File Organization

Recommended structure:

```
policies/
├── common/
│   └── security.cedar      # Cross-language security
├── go/
│   ├── merge.cedar         # Go merge requirements
│   ├── versions.cedar      # Go version enforcement
│   └── matrix.cedar        # OS matrix requirements
├── python/
│   └── versions.cedar      # Python version enforcement
└── dependencies/
    ├── age.cedar           # Dependency freshness
    └── vulnerabilities.cedar  # Security scanning
```

## Testing Policies

### Validate Syntax

```bash
pipelineconductor validate policies/ --verbose
```

### Dry Run

Test against your repos without taking action:

```bash
pipelineconductor scan --orgs myorg --policy-dir policies/ --format markdown
```

### Incremental Testing

1. Start with permissive policies
2. Review violations in reports
3. Gradually tighten requirements

## Common Mistakes

### Forgetting the Semicolon

```cedar
// Wrong - missing semicolon
forbid(principal, action, resource)
when { context.hasWorkflow == false }

// Correct
forbid(principal, action, resource)
when { context.hasWorkflow == false };
```

### Wrong Comparison Operator

```cedar
// Wrong - assignment instead of comparison
when { context.hasWorkflow = false };

// Correct
when { context.hasWorkflow == false };
```

### Invalid Context Variable

```cedar
// Wrong - typo in variable name
when { context.hasWorkFlow == false };

// Correct
when { context.hasWorkflow == false };
```

## Debugging Policies

### Enable Verbose Mode

```bash
pipelineconductor scan --orgs myorg --policy-dir policies/ -v
```

### Check Policy Loading

```bash
pipelineconductor validate policies/ --verbose
```

### Review Violations

The report shows which policies triggered violations:

```json
{
  "violations": [
    {
      "policy": "go/versions",
      "rule": "supported-versions",
      "message": "Go version 1.18 not in allowed versions",
      "severity": "medium"
    }
  ]
}
```

## See Also

- [Cedar Syntax](cedar-syntax.md) - Language reference
- [Examples](examples.md) - Complete policy examples
- [Built-in Policies](builtin.md) - Default policies
