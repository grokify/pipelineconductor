# Cedar Syntax

Cedar is a policy language that uses a declarative syntax to express authorization rules. This guide covers the Cedar syntax as used in PipelineConductor.

## Basic Structure

```cedar
effect(
    principal,
    action,
    resource
)
when {
    conditions
};
```

## Effects

### Permit

Allow an action when conditions are met:

```cedar
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == true
};
```

### Forbid

Deny an action when conditions are met:

```cedar
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasVulnerabilities == true
};
```

## Principal, Action, Resource

In PipelineConductor:

- **Principal** - Always the CI system (`CISystem::"pipelineconductor"`)
- **Action** - The CI/CD action being evaluated
- **Resource** - The repository being scanned

Use wildcards to match any:

```cedar
// Match any principal, action, or resource
permit(principal, action, resource);

// Match specific action
permit(principal, action == Action::"build", resource);
```

## Conditions

### Boolean Comparisons

```cedar
when {
    context.hasWorkflow == true
};

when {
    context.archived == false
};
```

### Numeric Comparisons

```cedar
when {
    context.oldestDependencyDays > 90
};

when {
    context.vulnerabilityCount >= 1
};
```

### String Comparisons

```cedar
when {
    context.repoOrg == "myorg"
};

when {
    context.goProfile == "modern"
};
```

### Set Operations

Check if a set contains a value:

```cedar
when {
    context.osMatrix.contains("ubuntu-latest")
};
```

Check if a set contains any of several values:

```cedar
when {
    context.osMatrix.containsAny(["ubuntu-latest", "macos-latest"])
};
```

Check if a set contains all values:

```cedar
when {
    context.osMatrix.containsAll(["ubuntu-latest", "macos-latest", "windows-latest"])
};
```

### Logical Operators

**AND** (implicit when on separate lines, explicit with `&&`):

```cedar
when {
    context.hasWorkflow == true &&
    context.branchProtectionEnabled == true
};
```

**OR**:

```cedar
when {
    context.hasRenovate == true || context.hasDependabot == true
};
```

**NOT**:

```cedar
when {
    !(context.archived == true)
};
```

### Parentheses for Grouping

```cedar
when {
    (context.hasRenovate == true || context.hasDependabot == true) &&
    context.hasVulnerabilities == false
};
```

## Comments

```cedar
// Single-line comment

/*
 * Multi-line
 * comment
 */

permit(
    principal,
    action == Action::"merge",  // Merge action
    resource
);
```

## Multiple Policies

Multiple policies in a file are evaluated together:

```cedar
// Policy 1: Require workflow
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == false
};

// Policy 2: Block vulnerable repos
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasVulnerabilities == true
};

// Policy 3: Allow when compliant
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == true &&
    context.hasVulnerabilities == false
};
```

## Policy Evaluation

Cedar policies are evaluated using a **default-deny** model:

1. If any `forbid` policy matches → **Denied**
2. If any `permit` policy matches → **Allowed**
3. If no policies match → **Denied** (default)

!!! note "Forbid Takes Precedence"
    If both `permit` and `forbid` policies match, `forbid` wins.

## Type Reference

| Type | Examples |
|------|----------|
| Boolean | `true`, `false` |
| Long | `0`, `90`, `365` |
| String | `"ubuntu-latest"`, `"myorg"` |
| Set | `["a", "b"]`, `context.languages` |
| EntityUID | `Action::"merge"`, `Repository::"org/repo"` |

## Best Practices

1. **Use comments** - Explain the intent of each policy
2. **Be specific** - Target specific actions rather than using wildcards
3. **Test policies** - Use `validate` command before deploying
4. **Group related policies** - Organize by domain (security, CI, dependencies)
5. **Prefer forbid for security** - Explicitly deny risky conditions

## See Also

- [Cedar Documentation](https://docs.cedarpolicy.com/)
- [Writing Policies](writing.md)
- [Examples](examples.md)
