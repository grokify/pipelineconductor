// Workflow Compliance Policies
// These policies enforce workflow compliance requirements based on compliance check results.

// ===== MERGE POLICIES =====

// Allow merge only for fully compliant repositories
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.compliant == true
};

// Block merge for non-compliant repositories
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "none"
}
unless {
    // Allow bypass for archived repos (no active development)
    context.archived == true
};

// ===== BUILD POLICIES =====

// Allow build for any repo with workflows
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.hasWorkflow == true
};

// Prefer repos using reusable workflows
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.usesReusableWorkflows == true &&
    context.exactMatchCount > 0
};

// ===== DEPLOY POLICIES =====

// Only allow deploy for fully compliant repos
permit(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.missingWorkflowCount == 0
};

// Block deploy if there are filename mismatches (strict mode)
forbid(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.hasFilenameMismatch == true
};

// ===== RELEASE POLICIES =====

// Allow release only for fully compliant repos with exact workflow matches
permit(
    principal,
    action == Action::"release",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.usesReusableWorkflows == true &&
    context.exactMatchCount > 0
};

// ===== LINT POLICIES =====

// Always allow lint action
permit(
    principal,
    action == Action::"lint",
    resource
)
when {
    context.hasWorkflow == true
};

// ===== TEST POLICIES =====

// Allow test for any repo with Go CI workflow
permit(
    principal,
    action == Action::"test",
    resource
)
when {
    context.hasWorkflow == true &&
    context.languages.contains("Go")
};
