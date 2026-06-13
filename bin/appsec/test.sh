#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
ENVIRONMENT_FILE="bin/shared/environment.sh"
source "$ENVIRONMENT_FILE"

# Console Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
RESET='\033[0;0m'

log_info() { echo -e "[INFO] $*"; }
log_error() { echo -e "[ERROR] $*" >&2; }

# Helper function to print results
print_result() {
    local name=$1
    local status=$2
    if [[ $status -eq 0 ]]; then
        echo -e "  - $name: ${GREEN}[OK]${RESET}"
        return 0
    else
        echo -e "  - $name: ${RED}[KO]${RESET}"
        return 1
    fi
}

setup_environment
show_config "full"

echo "--------------------------------------------------"
log_info "Starting security analysis..."
echo "--------------------------------------------------"

# Snyk Authentication
snyk auth "$SNYK_TOKEN" > /dev/null 2>&1

# Track global exit status
GLOBAL_EXIT=0

# Create a temporary directory to store logs if tools fail
LOG_DIR=$(mktemp -d)
trap 'rm -rf "$LOG_DIR"' EXIT

# --- SNYK CODE ---
if snyk code test --severity-threshold=medium --include-ignores > "$LOG_DIR/snyk_code.log" 2>&1; then
    print_result "Snyk Code (SAST)" 0
else
    print_result "Snyk Code (SAST)" 1
    GLOBAL_EXIT=1
    log_error "Snyk Code found Medium/High severity vulnerabilities:"
    cat "$LOG_DIR/snyk_code.log" >&2
fi

# --- SNYK SCA ---
if snyk test --all-projects --severity-threshold=medium --include-ignores > "$LOG_DIR/snyk_sca.log" 2>&1; then
    print_result "Snyk SCA (Dependencies)" 0
else
    print_result "Snyk SCA (Dependencies)" 1
    GLOBAL_EXIT=1
    log_error "Snyk SCA found vulnerabilities in your libraries:"
    cat "$LOG_DIR/snyk_sca.log" >&2
fi

# --- SNYK IAC ---
if snyk iac test --severity-threshold=high > "$LOG_DIR/snyk_iac.log" 2>&1; then
    print_result "Snyk IaC" 0
else
    if grep -qE "Could not find any valid IaC files|SNYK-CLI-0012|monthly limit" "$LOG_DIR/snyk_iac.log"; then
        echo -e "  - Snyk IaC: [SKIPPED] (No files found or limit reached)"
    else
        print_result "Snyk IaC" 1
        GLOBAL_EXIT=1
        log_error "Snyk IaC found configuration issues:"
        cat "$LOG_DIR/snyk_iac.log" >&2
    fi
fi

# --- GITLEAKS ---
if gitleaks detect > "$LOG_DIR/gitleaks.log" 2>&1; then
    print_result "Gitleaks (Secrets)" 0
else
    print_result "Gitleaks (Secrets)" 1
    GLOBAL_EXIT=1
    log_error "Gitleaks detected potential exposed credentials/tokens:"
    cat "$LOG_DIR/gitleaks.log" >&2
fi

# --- GOLANGCI-LINT ---
if golangci-lint run > "$LOG_DIR/golangci_lint.log" 2>&1; then
    print_result "GolangCI-Lint" 0
else
    print_result "GolangCI-Lint" 1
    GLOBAL_EXIT=1
    log_error "Linter errors detected in Go:"
    cat "$LOG_DIR/golangci_lint.log" >&2
fi

echo "--------------------------------------------------"
if [[ $GLOBAL_EXIT -eq 0 ]]; then
    log_info "Analysis finished successfully."
    exit 0
else
    log_error "Security analysis failed. Please review the [KO] steps above."
    exit 1
fi