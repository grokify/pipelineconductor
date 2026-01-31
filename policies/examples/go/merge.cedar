// Merge gating policy for Go repositories
// Permits merge only when CI tests have passed and branch protection is enabled

// Allow merge when tests passed and workflow exists
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == true &&
    context.lastRunPassed == true
};

// Require branch protection for merge
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.branchProtectionEnabled == true &&
    context.requireStatusChecks == true
};
