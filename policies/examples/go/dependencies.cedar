// Dependency management policy
// Controls merge based on dependency age and vulnerability status

// Block merge if vulnerabilities exist
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasVulnerabilities == true
};

// Block merge if dependencies are too old (over 90 days)
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.oldestDependencyDays > 90
};

// Allow merge when using dependency management tools
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    (context.hasRenovate == true || context.hasDependabot == true) &&
    context.hasVulnerabilities == false
};

// Allow merge for repos with recent dependencies
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.oldestDependencyDays <= 30 &&
    context.hasVulnerabilities == false
};
