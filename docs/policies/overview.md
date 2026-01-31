# Policy Overview

PipelineConductor uses **Cedar**, a policy language developed by AWS, to define CI/CD compliance rules. Cedar provides a declarative way to express what actions are permitted or forbidden based on repository context.

## Why Cedar?

| Feature | Benefit |
|---------|---------|
| **Declarative** | Express intent, not implementation |
| **Fast** | Sub-millisecond evaluation |
| **Type-safe** | Compile-time policy validation |
| **Go-native** | First-class Go support via cedar-go |
| **Readable** | Human-friendly syntax |

## Policy Structure

A Cedar policy consists of:

1. **Effect** - `permit` or `forbid`
2. **Principal** - Who is performing the action (the CI system)
3. **Action** - What is being done (`merge`, `build`, `test`, etc.)
4. **Resource** - What is being acted upon (the repository)
5. **Conditions** - When the policy applies (`when` clause)

## Basic Example

```cedar
// Require CI workflow for merge
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == false
};
```

This policy **forbids** merging when a repository has no CI workflow.

## Actions

PipelineConductor evaluates these actions:

| Action | Description |
|--------|-------------|
| `merge` | Merging code to the default branch |
| `build` | Building the project |
| `test` | Running tests |
| `lint` | Running linters |
| `deploy` | Deploying to production |
| `release` | Creating a release |

## Context Variables

Policies have access to repository context:

### Repository Info

| Variable | Type | Description |
|----------|------|-------------|
| `context.repoName` | String | Repository name |
| `context.repoOrg` | String | Organization name |
| `context.repoFullName` | String | Full name (org/repo) |
| `context.archived` | Boolean | Is archived |
| `context.fork` | Boolean | Is a fork |
| `context.languages` | Set | Programming languages |
| `context.topics` | Set | Repository topics |

### CI Configuration

| Variable | Type | Description |
|----------|------|-------------|
| `context.hasWorkflow` | Boolean | Has any workflow |
| `context.usesReusableWorkflow` | Boolean | Uses reusable workflows |
| `context.reusableWorkflowRef` | String | Reusable workflow reference |
| `context.lastRunPassed` | Boolean | Last CI run passed |
| `context.osMatrix` | Set | OS platforms in matrix |

### Go-Specific

| Variable | Type | Description |
|----------|------|-------------|
| `context.goVersions` | Set | Go versions in matrix |
| `context.goProfile` | String | Go profile name |
| `context.hasGoMod` | Boolean | Has go.mod file |

### Dependencies

| Variable | Type | Description |
|----------|------|-------------|
| `context.hasRenovate` | Boolean | Uses Renovate |
| `context.hasDependabot` | Boolean | Uses Dependabot |
| `context.oldestDependencyDays` | Long | Age of oldest dependency |
| `context.hasVulnerabilities` | Boolean | Has known vulnerabilities |
| `context.vulnerabilityCount` | Long | Number of vulnerabilities |

### Branch Protection

| Variable | Type | Description |
|----------|------|-------------|
| `context.branchProtectionEnabled` | Boolean | Protection enabled |
| `context.requireReviews` | Boolean | Requires PR reviews |
| `context.requireStatusChecks` | Boolean | Requires status checks |
| `context.branchProtectionEnforceAdmins` | Boolean | Enforced for admins |

## Policy Files

Cedar policies are stored in `.cedar` files:

```
policies/
├── go/
│   ├── merge.cedar
│   ├── versions.cedar
│   └── matrix.cedar
├── security/
│   └── branch-protection.cedar
└── dependencies/
    └── vulnerabilities.cedar
```

## Loading Policies

Policies can be loaded from:

1. **Built-in policies** - Included with PipelineConductor
2. **Local directory** - Via `--policy-dir` flag
3. **Remote repository** - Via `--policy-repo` flag (coming soon)

```bash
# Use built-in policies (default)
pipelineconductor scan --orgs myorg

# Add custom policies
pipelineconductor scan --orgs myorg --policy-dir ./policies/

# Disable built-in policies
pipelineconductor scan --orgs myorg --builtin-policies=false --policy-dir ./policies/
```

## Next Steps

- [Cedar Syntax](cedar-syntax.md) - Learn Cedar policy language
- [Built-in Policies](builtin.md) - Review default policies
- [Writing Policies](writing.md) - Create custom policies
- [Examples](examples.md) - Policy examples
