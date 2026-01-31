# PipelineConductor - Technical Requirements Document

## Architecture Overview

PipelineConductor follows a modular architecture with clear separation between data collection, policy evaluation, and remediation.

```
┌─────────────────────────────────────────────────────────────────┐
│                      PipelineConductor CLI                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Collectors  │  │    Policy    │  │     Remediator       │  │
│  │              │  │    Engine    │  │                      │  │
│  │ - GitHub API │  │              │  │ - PR Generator       │  │
│  │ - GitLab API │  │ - Cedar      │  │ - Patch Builder      │  │
│  │ - Git        │  │ - Profiles   │  │ - GitHub API         │  │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────┘  │
│         │                 │                      │              │
│         ▼                 ▼                      ▼              │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    pkg/model                             │   │
│  │  - Repo, CIWorkflow, PolicyContext, ComplianceResult    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│                        Report Generator                          │
│                   JSON | SARIF | Markdown | CSV                  │
└─────────────────────────────────────────────────────────────────┘
```

## Repository Structure

```
pipelineconductor/
├── .github/
│   └── workflows/
│       ├── ci.yml                 # CI for pipelineconductor itself
│       ├── compliance.yml         # Scheduled multi-repo compliance
│       └── release.yml            # Release automation
├── cmd/
│   └── pipelineconductor/
│       └── main.go                # CLI entry point
├── internal/
│   ├── collector/
│   │   ├── collector.go           # Collector interface
│   │   ├── github.go              # GitHub API collector
│   │   ├── gitlab.go              # GitLab API collector (future)
│   │   └── git.go                 # Git inspection utilities
│   ├── policy/
│   │   ├── engine.go              # Cedar evaluation engine
│   │   ├── loader.go              # Policy file loader
│   │   ├── context.go             # Context builder for Cedar
│   │   └── profiles.go            # Profile management
│   ├── remediator/
│   │   ├── remediator.go          # Remediator interface
│   │   ├── pr.go                  # PR creation logic
│   │   └── patch.go               # Workflow patch generation
│   └── report/
│       ├── report.go              # Report generation
│       ├── json.go                # JSON output
│       ├── sarif.go               # SARIF output
│       └── markdown.go            # Markdown output
├── pkg/
│   └── model/
│       ├── repo.go                # Repository model
│       ├── workflow.go            # CI workflow model
│       ├── policy.go              # Policy context model
│       └── result.go              # Compliance result model
├── policies/
│   └── examples/
│       └── go/
│           ├── merge.cedar        # Example merge policy
│           ├── versions.cedar     # Example version policy
│           └── dependencies.cedar # Example dependency policy
├── configs/
│   └── profiles/
│       ├── default.yaml           # Default Go profile
│       ├── modern.yaml            # Modern Go profile
│       └── legacy.yaml            # Legacy Go profile
├── go.mod
├── go.sum
├── PRD.md
├── TRD.md
├── ROADMAP.md
└── README.md
```

## Component Specifications

### 1. CLI (`cmd/pipelineconductor`)

#### Commands

```bash
# Run compliance scan
pipelineconductor scan \
  --orgs grokify,otherorg \
  --policy-repo github.com/grokify/pipelineconductor-policy@v1 \
  --profile default \
  --output report.json \
  --format json,sarif,markdown

# Dry-run remediation
pipelineconductor remediate \
  --orgs grokify \
  --dry-run \
  --policy-repo github.com/grokify/pipelineconductor-policy@v1

# Validate policies
pipelineconductor validate \
  --policy-repo github.com/grokify/pipelineconductor-policy@v1

# List repos and their compliance status
pipelineconductor list \
  --orgs grokify \
  --filter non-compliant
```

#### Configuration File

```yaml
# pipelineconductor.yaml
orgs:
  - grokify
  - otherorg

policy:
  repo: github.com/grokify/pipelineconductor-policy
  ref: v1

defaults:
  profile: default

scan:
  schedule: weekly
  include:
    - language: go
  exclude:
    - archived: true
    - fork: true

remediation:
  enabled: true
  dry_run: false
  branch_prefix: pipelineconductor/

output:
  formats:
    - json
    - sarif
  path: ./reports
```

### 2. Collectors (`internal/collector`)

#### Interface

```go
type Collector interface {
    // ListRepos returns all repos matching criteria
    ListRepos(ctx context.Context, opts ListOptions) ([]model.Repo, error)

    // GetRepoDetails fetches detailed info for a repo
    GetRepoDetails(ctx context.Context, repo model.Repo) (*model.RepoDetails, error)

    // GetWorkflows fetches CI workflow configurations
    GetWorkflows(ctx context.Context, repo model.Repo) ([]model.Workflow, error)

    // GetBranchProtection fetches branch protection rules
    GetBranchProtection(ctx context.Context, repo model.Repo, branch string) (*model.BranchProtection, error)
}
```

#### GitHub Collector Data Sources

| Data | API | Endpoint |
|------|-----|----------|
| Repo list | GraphQL | `organization.repositories` |
| Repo metadata | REST | `GET /repos/{owner}/{repo}` |
| Workflows | REST | `GET /repos/{owner}/{repo}/actions/workflows` |
| Branch protection | REST | `GET /repos/{owner}/{repo}/branches/{branch}/protection` |
| Required checks | REST | `GET /repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks` |
| Languages | REST | `GET /repos/{owner}/{repo}/languages` |

### 3. Policy Engine (`internal/policy`)

#### Cedar Integration

```go
type Engine struct {
    policySet *cedar.PolicySet
    schema    *cedar.Schema
}

func (e *Engine) Evaluate(ctx context.Context, input PolicyInput) (*PolicyResult, error) {
    // Build Cedar request
    req := cedar.Request{
        Principal: input.Principal,
        Action:    input.Action,
        Resource:  input.Resource,
        Context:   input.Context,
    }

    // Evaluate against policy set
    decision, diagnostics := e.policySet.IsAuthorized(entities, req)

    return &PolicyResult{
        Decision:    decision,
        Diagnostics: diagnostics,
    }, nil
}
```

#### Policy Context Schema

```go
type PolicyContext struct {
    // Repository information
    Repo struct {
        Name     string   `json:"name"`
        Org      string   `json:"org"`
        Language []string `json:"language"`
        Archived bool     `json:"archived"`
        Fork     bool     `json:"fork"`
    } `json:"repo"`

    // CI configuration
    CI struct {
        HasWorkflow           bool     `json:"hasWorkflow"`
        UsesReusableWorkflow  bool     `json:"usesReusableWorkflow"`
        ReusableWorkflowRef   string   `json:"reusableWorkflowRef"`
        RequiredChecks        []string `json:"requiredChecks"`
        LastRunPassed         bool     `json:"lastRunPassed"`
        OSMatrix              []string `json:"osMatrix"`
    } `json:"ci"`

    // Go-specific (when applicable)
    Go struct {
        Versions []string `json:"versions"`
        Profile  string   `json:"profile"`
    } `json:"go"`

    // Dependencies
    Dependencies struct {
        HasRenovate        bool `json:"hasRenovate"`
        HasDependabot      bool `json:"hasDependabot"`
        OldestDependencyDays int `json:"oldestDependencyDays"`
        HasVulnerabilities bool `json:"hasVulnerabilities"`
        VulnerabilityCount int  `json:"vulnerabilityCount"`
    } `json:"dependencies"`

    // Branch protection
    BranchProtection struct {
        Enabled              bool `json:"enabled"`
        RequireReviews       bool `json:"requireReviews"`
        RequireStatusChecks  bool `json:"requireStatusChecks"`
        EnforceAdmins        bool `json:"enforceAdmins"`
    } `json:"branchProtection"`
}
```

### 4. Profiles (`configs/profiles`)

#### Default Profile

```yaml
# configs/profiles/default.yaml
name: default
description: Standard Go CI configuration for active projects

go:
  versions:
    - "1.24"
    - "1.25"

os:
  - ubuntu-latest
  - macos-latest
  - windows-latest

checks:
  required:
    - test
    - lint
    - build

lint:
  enabled: true
  tool: golangci-lint

test:
  coverage: true
  race: true
```

#### Legacy Profile

```yaml
# configs/profiles/legacy.yaml
name: legacy
description: Legacy Go CI for older projects requiring Go 1.12-1.18

go:
  versions:
    - "1.12"

os:
  - ubuntu-latest

checks:
  required:
    - test
    - build

lint:
  enabled: false  # Modern linters don't support old Go

test:
  coverage: false
  race: false  # Race detector unreliable on old Go
```

### 5. Remediator (`internal/remediator`)

#### Remediation Actions

| Violation | Remediation |
|-----------|-------------|
| Missing CI workflow | Create `.github/workflows/ci.yml` with reusable workflow call |
| Wrong Go versions | Update workflow matrix to match profile |
| Missing branch protection | Log warning (requires admin, cannot auto-fix) |
| Old reusable workflow ref | Update `@v1` to `@v2` |
| Missing required checks | Update workflow to include required jobs |

#### PR Generation

```go
type PRGenerator struct {
    client *github.Client
}

func (g *PRGenerator) CreateRemediationPR(ctx context.Context, opts PROptions) (*github.PullRequest, error) {
    // 1. Create branch
    branchName := fmt.Sprintf("pipelineconductor/remediate-%s", time.Now().Format("20060102"))

    // 2. Apply patches
    for _, patch := range opts.Patches {
        // Create or update file
    }

    // 3. Create PR
    pr, _, err := g.client.PullRequests.Create(ctx, opts.Owner, opts.Repo, &github.NewPullRequest{
        Title: github.String("[PipelineConductor] Fix CI compliance violations"),
        Body:  github.String(opts.Description),
        Head:  github.String(branchName),
        Base:  github.String(opts.BaseBranch),
    })

    return pr, err
}
```

### 6. Report Generator (`internal/report`)

#### Output Formats

**JSON**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "summary": {
    "total": 150,
    "compliant": 142,
    "non_compliant": 8,
    "compliance_rate": 94.67
  },
  "repos": [
    {
      "name": "grokify/example",
      "compliant": false,
      "violations": [
        {
          "policy": "go/versions",
          "message": "Go 1.21 not in allowed versions [1.24, 1.25]",
          "severity": "warning",
          "remediation": "Update go-version in workflow"
        }
      ]
    }
  ]
}
```

**SARIF** (for GitHub Security tab)
```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "PipelineConductor",
          "version": "0.1.0"
        }
      },
      "results": []
    }
  ]
}
```

## Data Flow

```
1. Scan Initiation
   CLI/GitHub Action triggers scan

2. Discovery Phase (API-first)
   ├── List all repos in specified orgs
   ├── Filter by criteria (language, archived, fork)
   └── Collect metadata via GitHub API

3. Inspection Phase (selective git)
   ├── For repos needing deep inspection
   ├── Fetch workflow files
   └── Parse and extract configuration

4. Evaluation Phase (Cedar)
   ├── Build PolicyContext for each repo
   ├── Load policies from policy repo
   └── Evaluate each repo against policies

5. Reporting Phase
   ├── Aggregate results
   ├── Generate reports in requested formats
   └── Upload artifacts (if in GitHub Actions)

6. Remediation Phase (optional)
   ├── For non-compliant repos with auto-fix enabled
   ├── Generate patches
   └── Create pull requests
```

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | Go 1.24+ | Native GitHub Actions support, strong typing |
| Policy Engine | Cedar (cedar-go) | Static analysis, compile-time testing |
| GitHub API | google/go-github | Well-maintained, full API coverage |
| CLI Framework | cobra | Industry standard for Go CLIs |
| Config | viper | Flexible configuration management |
| Logging | slog | Go 1.21+ structured logging |
| Testing | testify | Assertions and mocking |

## Security Considerations

1. **Token Scoping**: Use minimal required GitHub token permissions
2. **Read-Only Default**: Remediation requires explicit opt-in
3. **No Secret Storage**: Policies must not contain secrets
4. **Audit Logging**: All actions logged with context
5. **Rate Limiting**: Respect GitHub API rate limits

## Performance Targets

| Metric | Target |
|--------|--------|
| Repos scanned per minute | 100+ |
| Policy evaluation per repo | <100ms |
| Memory usage | <500MB for 1000 repos |
| API calls per repo | <10 average |

## Testing Strategy

1. **Unit Tests**: All packages have >80% coverage
2. **Integration Tests**: Mock GitHub API responses
3. **Policy Tests**: Cedar policies tested with golden inputs
4. **E2E Tests**: Test against real repos in test org

## Shared Libraries Strategy

PipelineConductor leverages shared libraries for reusable functionality and will contribute back where appropriate.

### github.com/grokify/gogithub

Extended GitHub API functionality. **Check gogithub first** before implementing GitHub-specific code.

| Package | Use For |
|---------|---------|
| `auth` | Token management, client creation |
| `repo` | Repository listing, metadata |
| `pr` | Pull request creation/management |
| `release` | Release operations |
| `search` | Search API queries |
| `graphql` | GraphQL operations |

**Decision tree:**

1. **Exists in gogithub?** → Import and use
2. **Reusable GitHub utility?** → Add to gogithub, then import
3. **CI/CD-specific logic?** → Implement in PipelineConductor

### github.com/grokify/omniscm (Future)

Multi-SCM abstraction to be extracted after implementing 3+ providers:

```
omniscm/
├── provider/
│   ├── interface.go      # Provider interface
│   └── registry.go       # Provider registration
├── github/
│   └── github.go         # GitHub implementation
├── gitlab/
│   └── gitlab.go         # GitLab implementation
├── bitbucket/
│   └── bitbucket.go      # Bitbucket implementation
└── model/
    ├── repo.go           # Provider-agnostic repo
    ├── pr.go             # Provider-agnostic PR
    └── workflow.go       # Provider-agnostic workflow
```

**Extraction criteria:**

- 3+ working provider implementations
- Stable Collector interface (no breaking changes for 2+ releases)
- Clear separation of provider-agnostic vs provider-specific code
