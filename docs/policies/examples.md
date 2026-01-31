# Policy Examples

Complete Cedar policy examples for common CI/CD scenarios.

## Go Version Enforcement

Ensure repositories use supported Go versions:

```cedar
// policies/go/versions.cedar

// Block outdated Go versions
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.goVersions.containsAny(["1.18", "1.19", "1.20"])
};

// Prefer latest Go versions
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.goVersions.containsAny(["1.24", "1.25"])
};
```

## Multi-Platform Matrix

Require testing across multiple operating systems:

```cedar
// policies/go/matrix.cedar

// Allow build when testing on Linux
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.contains("ubuntu-latest")
};

// Prefer builds that test on all platforms
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.containsAll([
        "ubuntu-latest",
        "macos-latest",
        "windows-latest"
    ])
};
```

## Dependency Management

Control merging based on dependency health:

```cedar
// policies/dependencies.cedar

// Block repos with known vulnerabilities
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasVulnerabilities == true
};

// Block repos with stale dependencies (over 90 days)
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.oldestDependencyDays > 90
};

// Allow when using dependency automation
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    (context.hasRenovate == true || context.hasDependabot == true) &&
    context.hasVulnerabilities == false
};

// Allow repos with fresh dependencies
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.oldestDependencyDays <= 30 &&
    context.hasVulnerabilities == false
};
```

## Reusable Workflows

Encourage use of centralized CI/CD workflows:

```cedar
// policies/go/reusable-workflow.cedar

// Prefer repos using reusable workflows for merge
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.usesReusableWorkflow == true
};

// Allow build for repos with any workflow
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.hasWorkflow == true
};

// Strongly prefer reusable workflows for consistency
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.usesReusableWorkflow == true &&
    context.reusableWorkflowRef != ""
};
```

## Branch Protection

Enforce branch protection rules:

```cedar
// policies/security/branch-protection.cedar

// Require branch protection for merge
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.branchProtectionEnabled == false
};

// Require PR reviews for production repos
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.topics.contains("production") &&
    context.requireReviews == false
};

// Require status checks
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.requireStatusChecks == false
};
```

## Release Gating

Control when releases can be created:

```cedar
// policies/release.cedar

// Block releases with vulnerabilities
forbid(
    principal,
    action == Action::"release",
    resource
)
when {
    context.hasVulnerabilities == true
};

// Block releases if last CI run failed
forbid(
    principal,
    action == Action::"release",
    resource
)
when {
    context.lastRunPassed == false
};

// Allow release when all checks pass
permit(
    principal,
    action == Action::"release",
    resource
)
when {
    context.hasWorkflow == true &&
    context.lastRunPassed == true &&
    context.hasVulnerabilities == false &&
    context.branchProtectionEnabled == true
};
```

## Organization-Specific Rules

Different rules for different organizations:

```cedar
// policies/org-specific.cedar

// Stricter rules for public org
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.repoOrg == "public-org" &&
    context.branchProtectionEnabled == false
};

// Relaxed rules for sandbox org
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.repoOrg == "sandbox-org"
};

// Require all checks for production org
forbid(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.repoOrg == "production-org" &&
    (context.hasWorkflow == false ||
     context.branchProtectionEnabled == false ||
     context.hasVulnerabilities == true)
};
```

## Archive and Fork Handling

Special rules for archived repos and forks:

```cedar
// policies/repo-status.cedar

// Skip archived repos
permit(
    principal,
    action,
    resource
)
when {
    context.archived == true
};

// Relaxed rules for forks
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.fork == true &&
    context.hasWorkflow == true
};
```

## Complete Policy Set

A comprehensive policy set for a Go organization:

```cedar
// policies/complete.cedar

// ===================
// Core Requirements
// ===================

// Require CI workflow
forbid(principal, action == Action::"merge", resource)
when { context.hasWorkflow == false };

// Require branch protection
forbid(principal, action == Action::"merge", resource)
when { context.branchProtectionEnabled == false };

// ===================
// Security
// ===================

// Block vulnerable repos from merge and release
forbid(principal, action == Action::"merge", resource)
when { context.hasVulnerabilities == true };

forbid(principal, action == Action::"release", resource)
when { context.hasVulnerabilities == true };

// ===================
// Go-Specific
// ===================

// Require go.mod
forbid(principal, action == Action::"build", resource)
when {
    context.languages.contains("Go") &&
    context.hasGoMod == false
};

// Require modern Go versions
forbid(principal, action == Action::"merge", resource)
when {
    context.languages.contains("Go") &&
    context.goVersions.containsAny(["1.18", "1.19", "1.20"])
};

// ===================
// Positive Permits
// ===================

// Allow merge when fully compliant
permit(principal, action == Action::"merge", resource)
when {
    context.hasWorkflow == true &&
    context.branchProtectionEnabled == true &&
    context.hasVulnerabilities == false
};

// Allow build with workflow
permit(principal, action == Action::"build", resource)
when { context.hasWorkflow == true };

// Allow test with workflow
permit(principal, action == Action::"test", resource)
when { context.hasWorkflow == true };
```

## See Also

- [Cedar Syntax](cedar-syntax.md) - Language reference
- [Writing Policies](writing.md) - Policy authoring guide
- [validate Command](../cli/validate.md) - Policy validation
