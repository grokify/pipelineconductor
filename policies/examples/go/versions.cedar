// Go version enforcement policy
// Ensures builds use approved Go versions from the profile

// Allow build with current supported Go versions (1.24, 1.25)
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.goVersions.containsAny(["1.24", "1.25"])
};

// Allow test with current supported Go versions
permit(
    principal,
    action == Action::"test",
    resource
)
when {
    context.goVersions.containsAny(["1.24", "1.25"])
};

// Allow lint with current supported Go versions
permit(
    principal,
    action == Action::"lint",
    resource
)
when {
    context.goVersions.containsAny(["1.24", "1.25"])
};
