// Reusable workflow policy
// Encourages use of centralized CI/CD workflows

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
