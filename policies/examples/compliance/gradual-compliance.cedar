// Gradual Compliance Policies
// These policies allow for gradual adoption of workflow compliance.
// Use these for organizations transitioning to reusable workflows.

// ===== GRADUAL MERGE POLICIES =====

// Allow merge for full or partial compliance
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "full" ||
    context.complianceLevel == "partial"
};

// Allow merge for non-compliant repos if they have any workflow
permit(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.complianceLevel == "none" &&
    context.hasWorkflow == true
};

// Block merge only for repos with no workflows at all
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == false
}
unless {
    context.archived == true
};

// ===== GRADUAL BUILD POLICIES =====

// Allow build for any repo (no compliance requirement)
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.hasWorkflow == true
};

// ===== GRADUAL DEPLOY POLICIES =====

// Allow deploy for partial or full compliance (not none)
permit(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.complianceLevel != "none"
};

// Allow deploy for non-compliant if compliance rate is above 50%
permit(
    principal,
    action == Action::"deploy",
    resource
)
when {
    context.complianceRate >= 50
};

// ===== GRADUAL RELEASE POLICIES =====

// Allow release for partial compliance with at least one exact match
permit(
    principal,
    action == Action::"release",
    resource
)
when {
    context.complianceLevel == "partial" &&
    context.exactMatchCount >= 1
};

// Allow release for full compliance (any match type)
permit(
    principal,
    action == Action::"release",
    resource
)
when {
    context.complianceLevel == "full"
};
