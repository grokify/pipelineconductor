// Strict Compliance Policies
// These policies enforce strict workflow compliance requirements.
// Use these for organizations that require exact reusable workflow usage.

// ===== STRICT MERGE POLICIES =====

// Only allow merge for repos using exact reusable workflows (no equivalent matches)
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.exactMatchCount > 0 &&
    context.equivalentMatchCount == 0 &&
    context.hasFilenameMismatch == false
};

// Block merge for repos with any equivalent matches (should use reusable)
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.equivalentMatchCount > 0 &&
    context.exactMatchCount == 0
};

// ===== STRICT DEPLOY POLICIES =====

// Only deploy repos with 100% exact matches
permit(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.complianceRate == 100 &&
    context.usesReusableWorkflows == true &&
    context.hasFilenameMismatch == false
};

// Block deploy for any non-full compliance
forbid(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.complianceLevel != "full"
};

// ===== STRICT RELEASE POLICIES =====

// Require full compliance with exact matches for release
permit(
    principal,
    action == Action::"release",
    resource
)
when {
    context.complianceLevel == "full" &&
    context.exactMatchCount > 0 &&
    context.equivalentMatchCount == 0 &&
    context.missingWorkflowCount == 0
};

// Block release if any filename mismatches exist
forbid(
    principal,
    action == Action::"release",
    resource
)
when {
    context.hasFilenameMismatch == true
};
