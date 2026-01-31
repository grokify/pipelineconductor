// OS matrix enforcement policy
// Ensures CI runs on required platforms

// Allow build when testing on Linux
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.contains("ubuntu-latest")
};

// Prefer builds that also test on macOS and Windows
permit(
    principal,
    action == Action::"build",
    resource
)
when {
    context.osMatrix.contains("ubuntu-latest") &&
    context.osMatrix.contains("macos-latest") &&
    context.osMatrix.contains("windows-latest")
};

// Allow test on any supported platform
permit(
    principal,
    action == Action::"test",
    resource
)
when {
    context.osMatrix.containsAny(["ubuntu-latest", "macos-latest", "windows-latest"])
};
