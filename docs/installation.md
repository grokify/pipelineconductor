# Installation

PipelineConductor can be installed via Go install, from source, or using pre-built binaries.

## Requirements

- Go 1.24 or later (for building from source)
- GitHub personal access token with `repo` scope

## Go Install

The simplest way to install PipelineConductor:

```bash
go install github.com/grokify/pipelineconductor/cmd/pipelineconductor@latest
```

This installs the `pipelineconductor` binary to your `$GOPATH/bin` directory.

## Build from Source

Clone and build the repository:

```bash
git clone https://github.com/grokify/pipelineconductor.git
cd pipelineconductor
go build -o pipelineconductor ./cmd/pipelineconductor
```

Optionally, install to your PATH:

```bash
go install ./cmd/pipelineconductor
```

## Pre-built Binaries

Download pre-built binaries from the [GitHub Releases](https://github.com/grokify/pipelineconductor/releases) page.

=== "Linux (amd64)"

    ```bash
    curl -LO https://github.com/grokify/pipelineconductor/releases/latest/download/pipelineconductor_linux_amd64.tar.gz
    tar xzf pipelineconductor_linux_amd64.tar.gz
    sudo mv pipelineconductor /usr/local/bin/
    ```

=== "macOS (arm64)"

    ```bash
    curl -LO https://github.com/grokify/pipelineconductor/releases/latest/download/pipelineconductor_darwin_arm64.tar.gz
    tar xzf pipelineconductor_darwin_arm64.tar.gz
    sudo mv pipelineconductor /usr/local/bin/
    ```

=== "Windows"

    Download `pipelineconductor_windows_amd64.zip` and add to your PATH.

## Verify Installation

Check that PipelineConductor is installed correctly:

```bash
pipelineconductor version
```

Expected output:

```
pipelineconductor dev
  commit: none
  built:  unknown
```

## GitHub Token Setup

PipelineConductor requires a GitHub personal access token to access the GitHub API.

### Create a Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select scopes:
    - `repo` - Full control of private repositories (or `public_repo` for public only)
    - `read:org` - Read organization membership
4. Generate and copy the token

### Configure the Token

Set the token as an environment variable:

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

Or pass it directly to commands:

```bash
pipelineconductor scan --github-token ghp_your_token_here --orgs myorg
```

!!! warning "Token Security"
    Never commit your GitHub token to version control. Use environment variables or a secrets manager.

## Next Steps

- [Quick Start](quickstart.md) - Run your first scan
- [Configuration](cli/config.md) - Set up a configuration file
