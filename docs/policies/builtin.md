# Built-in Policies

PipelineConductor includes built-in policies that enforce common CI/CD best practices. These policies are enabled by default.

## Included Policies

### ci/workflow-required

**Severity:** High

Requires repositories to have at least one GitHub Actions workflow.

```cedar
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.hasWorkflow == false
};
```

**Remediation:** Create a `.github/workflows/ci.yml` file.

### security/branch-protection

**Severity:** Medium

Requires branch protection to be enabled on the default branch.

```cedar
forbid(
    principal,
    action == Action::"merge",
    resource
)
when {
    context.branchProtectionEnabled == false
};
```

**Remediation:** Enable branch protection in repository settings.

## Disabling Built-in Policies

To use only custom policies:

```bash
pipelineconductor scan --orgs myorg --builtin-policies=false --policy-dir ./policies/
```

## Extending Built-in Policies

You can add custom policies alongside built-in ones:

```bash
pipelineconductor scan --orgs myorg --policy-dir ./policies/
```

Your custom policies will be evaluated together with built-in policies.

## Policy Precedence

When multiple policies apply:

1. **Forbid always wins** - If any `forbid` matches, action is denied
2. **Permit allows** - If a `permit` matches (and no `forbid`), action is allowed
3. **Default deny** - If no policies match, action is denied

## Viewing Built-in Policies

Validate built-in policies to confirm they load correctly:

```bash
pipelineconductor validate --builtin
```

Output:

```
Validating built-in policies...
✓ Built-in policies are valid
```

## See Also

- [Writing Policies](writing.md) - Create custom policies
- [Examples](examples.md) - Policy examples
