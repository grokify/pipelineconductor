#!/bin/bash
# batch-update-workflows.sh - Batch update workflows across repositories
#
# This script wraps pipelineconductor commands for convenient batch workflow updates.
#
# Usage:
#   ./batch-update-workflows.sh --check                    # Check compliance
#   ./batch-update-workflows.sh --remediate --dry-run      # Preview remediation
#   ./batch-update-workflows.sh --remediate                # Apply remediation
#   ./batch-update-workflows.sh --apply --push             # Apply and push
#   ./batch-update-workflows.sh --apply --push --create-pr # Apply, push, and create PRs

set -euo pipefail

# Default configuration
LOCAL_PATH="${LOCAL_PATH:-$HOME/go/src/github.com}"
ORGS="${ORGS:-plexusone}"
LANGUAGES="${LANGUAGES:-Go}"
REF_REPO="${REF_REPO:-plexusone/.github}"
BRANCH="${BRANCH:-ci/update-workflows}"
PIPELINECONDUCTOR="${PIPELINECONDUCTOR:-pipelineconductor}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
    cat <<EOF
Usage: $0 [OPTIONS] COMMAND

Commands:
  --check           Check workflow compliance across repositories
  --remediate       Generate compliant workflow files
  --apply           Apply workflows and commit changes

Options:
  --local PATH      Base path for local repos (default: \$HOME/go/src/github.com)
  --orgs ORGS       Comma-separated list of organizations (default: plexusone)
  --languages LANG  Comma-separated list of languages (default: Go)
  --ref-repo REPO   Reference workflow repository (default: plexusone/.github)
  --repo NAME       Target specific repository
  --dry-run         Preview changes without applying
  --push            Push commits to remote (apply only)
  --create-pr       Create pull requests (requires --push)
  --branch NAME     Branch name for commits (default: ci/update-workflows)
  --verbose         Enable verbose output
  -h, --help        Show this help message

Environment Variables:
  LOCAL_PATH        Base path for local repos
  ORGS              Organizations to scan
  LANGUAGES         Languages to filter
  REF_REPO          Reference workflow repository
  PIPELINECONDUCTOR Path to pipelineconductor binary

Examples:
  # Check compliance for plexusone Go repos
  $0 --check

  # Preview what would be generated
  $0 --remediate --dry-run

  # Generate and commit workflows locally
  $0 --remediate

  # Full workflow: generate, commit, push, create PRs
  $0 --apply --push --create-pr

  # Update specific repo
  $0 --apply --repo vibium-wcag --push
EOF
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Parse arguments
COMMAND=""
DRY_RUN=""
PUSH=""
CREATE_PR=""
VERBOSE=""
REPO=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --check)
            COMMAND="check"
            shift
            ;;
        --remediate)
            COMMAND="remediate"
            shift
            ;;
        --apply)
            COMMAND="apply"
            shift
            ;;
        --local)
            LOCAL_PATH="$2"
            shift 2
            ;;
        --orgs)
            ORGS="$2"
            shift 2
            ;;
        --languages)
            LANGUAGES="$2"
            shift 2
            ;;
        --ref-repo)
            REF_REPO="$2"
            shift 2
            ;;
        --repo)
            REPO="$2"
            shift 2
            ;;
        --branch)
            BRANCH="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN="--dry-run"
            shift
            ;;
        --push)
            PUSH="--push"
            shift
            ;;
        --create-pr)
            CREATE_PR="--create-pr"
            shift
            ;;
        --verbose|-v)
            VERBOSE="--verbose"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

if [[ -z "$COMMAND" ]]; then
    log_error "No command specified"
    usage
    exit 1
fi

# Check if pipelineconductor exists
if ! command -v "$PIPELINECONDUCTOR" &> /dev/null; then
    # Try to find it in common locations
    if [[ -x "./pipelineconductor" ]]; then
        PIPELINECONDUCTOR="./pipelineconductor"
    elif [[ -x "$HOME/go/bin/pipelineconductor" ]]; then
        PIPELINECONDUCTOR="$HOME/go/bin/pipelineconductor"
    else
        log_error "pipelineconductor not found. Please install it or set PIPELINECONDUCTOR env var."
        exit 1
    fi
fi

log_info "Using pipelineconductor: $PIPELINECONDUCTOR"
log_info "Local path: $LOCAL_PATH"
log_info "Organizations: $ORGS"
log_info "Languages: $LANGUAGES"
log_info "Reference repo: $REF_REPO"

# Build common arguments
COMMON_ARGS=(
    --local "$LOCAL_PATH"
    --orgs "$ORGS"
    --languages "$LANGUAGES"
    --ref-repo "$REF_REPO"
)

if [[ -n "$REPO" ]]; then
    COMMON_ARGS+=(--repo "$REPO")
fi

if [[ -n "$VERBOSE" ]]; then
    COMMON_ARGS+=($VERBOSE)
fi

case $COMMAND in
    check)
        log_info "Checking workflow compliance..."
        "$PIPELINECONDUCTOR" check "${COMMON_ARGS[@]}" --format markdown
        ;;
    remediate)
        log_info "Running remediation..."
        REMEDIATE_ARGS=("${COMMON_ARGS[@]}")
        if [[ -n "$DRY_RUN" ]]; then
            REMEDIATE_ARGS+=($DRY_RUN)
        fi
        "$PIPELINECONDUCTOR" remediate "${REMEDIATE_ARGS[@]}"
        ;;
    apply)
        log_info "Applying workflows..."
        APPLY_ARGS=("${COMMON_ARGS[@]}")
        APPLY_ARGS+=(--branch "$BRANCH")
        if [[ -n "$DRY_RUN" ]]; then
            APPLY_ARGS+=($DRY_RUN)
        fi
        if [[ -n "$PUSH" ]]; then
            APPLY_ARGS+=($PUSH)
        fi
        if [[ -n "$CREATE_PR" ]]; then
            APPLY_ARGS+=($CREATE_PR)
        fi
        "$PIPELINECONDUCTOR" apply "${APPLY_ARGS[@]}"
        ;;
esac

log_success "Done!"
