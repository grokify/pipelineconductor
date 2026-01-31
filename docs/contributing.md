# Contributing

Thank you for your interest in contributing to PipelineConductor!

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- GitHub account

### Fork and Clone

```bash
# Fork the repository on GitHub, then clone
git clone https://github.com/YOUR_USERNAME/pipelineconductor.git
cd pipelineconductor
```

### Build and Test

```bash
# Build
go build -v ./...

# Run tests
go test -v ./...

# Run linter
golangci-lint run
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

- Follow Go conventions
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
go test -v ./...

# Run linter
golangci-lint run

# Test CLI
go run ./cmd/pipelineconductor --help
```

### 4. Commit Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat: add new policy evaluation feature"
git commit -m "fix: resolve race condition in scanner"
git commit -m "docs: update CLI reference"
```

### 5. Push and Create PR

```bash
git push origin feature/my-feature
```

Then create a Pull Request on GitHub.

## Code Style

### Go Conventions

- Use `gofmt` for formatting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `golangci-lint run` before committing

### Naming

- Use descriptive names
- Exported functions should have doc comments
- Package names should be lowercase, single words

### Error Handling

- Always handle errors
- Use `fmt.Errorf("context: %w", err)` for wrapping
- Return errors rather than logging them

### Testing

- Write table-driven tests
- Use meaningful test names
- Test edge cases

## Project Structure

```
pipelineconductor/
├── cmd/
│   └── pipelineconductor/  # CLI application
├── internal/
│   ├── collector/          # GitHub API integration
│   ├── policy/             # Cedar policy engine
│   └── report/             # Report generation
├── pkg/
│   └── model/              # Shared data models
├── policies/
│   └── examples/           # Example Cedar policies
├── configs/
│   └── profiles/           # Profile configurations
└── docs/                   # Documentation
```

## Types of Contributions

### Bug Reports

File an issue with:

- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Version information

### Feature Requests

File an issue with:

- Use case description
- Proposed solution
- Alternatives considered

### Code Contributions

Areas we welcome contributions:

- **New collectors** - GitLab, Bitbucket support
- **Report formats** - New output formats
- **Policy features** - Cedar policy enhancements
- **Documentation** - Guides, examples, tutorials
- **Tests** - Improved coverage

### Policy Contributions

Share useful Cedar policies:

1. Add to `policies/examples/`
2. Include comments explaining the policy
3. Add documentation in `docs/policies/examples.md`

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation updated
- [ ] Commit messages follow conventions

### PR Description

Include:

- What changes are made
- Why they're needed
- How to test
- Related issues

### Review Process

1. Maintainer reviews code
2. CI checks pass
3. Changes requested (if any)
4. Approval
5. Merge

## Release Process

Releases are managed by maintainers:

1. Update version in code
2. Update CHANGELOG
3. Create git tag
4. GitHub Actions builds releases

## Getting Help

- File an issue for bugs or questions
- Check existing issues first
- Join discussions

## Code of Conduct

Be respectful and inclusive. We welcome contributors of all backgrounds and experience levels.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
