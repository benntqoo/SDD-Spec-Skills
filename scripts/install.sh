#!/usr/bin/env bash
#
# install.sh - Installation script for vic CLI tool
#
# Usage:
#   ./install.sh              # Install vic
#   ./install.sh --uninstall  # Uninstall vic
#   ./install.sh --help       # Show help
#
# Supports: Linux (Intel/ARM), macOS (Intel/Apple Silicon)
#

set -euo pipefail

# ============================================================================
# CONFIGURATION
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}" && git rev-parse --show-toplevel 2>/dev/null || echo "${SCRIPT_DIR}")"
VIC_SOURCE_DIR="${REPO_ROOT}/cmd/vic-go"
BINARY_NAME="vic"
INSTALL_DIR_USER="${HOME}/.local/bin"
INSTALL_DIR_SYSTEM="/usr/local/bin"
PRECOMMIT_CONFIG="${REPO_ROOT}/.pre-commit-config.yaml"

# Colors for output (detect terminal support)
if [[ -t 1 ]] && command -v tput &>/dev/null && [[ $(tput colors 2>/dev/null || echo 0) -ge 8 ]]; then
    RED=$'\033[0;31m'
    GREEN=$'\033[0;32m'
    YELLOW=$'\033[1;33m'
    BLUE=$'\033[0;34m'
    BOLD=$'\033[1m'
    NC=$'\033[0m'
else
    RED=""
    GREEN=""
    YELLOW=""
    BLUE=""
    BOLD=""
    NC=""
fi

# ============================================================================
# UTILITY FUNCTIONS
# ============================================================================

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
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

# Detect OS: linux, darwin, or windows (git bash/wsl)
detect_os() {
    local os
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          os="unknown" ;;
    esac
    echo "${os}"
}

# Detect architecture: amd64, arm64, or armv7l
detect_arch() {
    local arch
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        armv7l)     arch="arm" ;;
        *)          arch="unknown" ;;
    esac
    echo "${arch}"
}

# Check if command exists
command_exists() {
    command -v "$1" &>/dev/null
}

# Check if running with sudo privileges
has_sudo() {
    if [[ "${OS}" == "windows" ]]; then
        return 1
    fi
    sudo -n true 2>/dev/null
}

# Prompt user for yes/no
prompt_yes_no() {
    local prompt="$1"
    local default="${2:-no}"
    local response

    while true; do
        if [[ "${default}" == "yes" ]]; then
            read -r -p "${prompt} [Y/n] " response
            response="${response:-Y}"
        else
            read -r -p "${prompt} [y/N] " response
            response="${response:-N}"
        fi

        case "${response}" in
            [Yy]|[Yy][Ee][Ss]) return 0 ;;
            [Nn]|[Nn][Oo]) return 1 ;;
            *) echo "Please answer yes or no." ;;
        esac
    done
}

# ============================================================================
# PREREQUISITE CHECKS
# ============================================================================

check_go() {
    log_info "Checking for Go installation..."

    if ! command_exists go; then
        log_error "Go is not installed."
        echo ""
        echo "Please install Go first:"
        echo "  macOS:  brew install go"
        echo "  Linux:  sudo apt install golang-go  # Debian/Ubuntu"
        echo "          sudo dnf install golang     # Fedora"
        echo ""
        echo "Or download from: https://go.dev/dl/"
        exit 1
    fi

    local go_version
    go_version="$(go version | grep -oE 'go[0-9]+\.[0-9]+' | head -1)"
    log_info "Found Go: $(go version)"
}

check_git() {
    if ! command_exists git; then
        log_warning "Git is not installed. Some features may be limited."
    fi
}

check_precommit() {
    if command_exists pre-commit; then
        log_info "Found pre-commit: $(pre-commit --version)"
        return 0
    fi
    log_warning "pre-commit not found. Install with: pip install pre-commit"
    return 1
}

# ============================================================================
# SOURCE VERIFICATION
# ============================================================================

verify_source() {
    log_info "Verifying vic-go source..."

    if [[ ! -d "${VIC_SOURCE_DIR}" ]]; then
        log_error "vic-go source directory not found: ${VIC_SOURCE_DIR}"
        exit 1
    fi

    if [[ ! -f "${VIC_SOURCE_DIR}/main.go" ]]; then
        log_error "main.go not found in ${VIC_SOURCE_DIR}"
        exit 1
    fi

    if [[ ! -f "${VIC_SOURCE_DIR}/go.mod" ]]; then
        log_error "go.mod not found in ${VIC_SOURCE_DIR}"
        exit 1
    fi

    log_success "Source verified: ${VIC_SOURCE_DIR}"
}

# ============================================================================
# BUILD
# ============================================================================

build_vic() {
    log_info "Building vic binary..."

    local build_dir="${VIC_SOURCE_DIR}/dist"
    mkdir -p "${build_dir}"

    # Build with optimization flags for smaller binary
    if ! go build -ldflags="-s -w" -o "${build_dir}/${BINARY_NAME}" .; then
        log_error "Build failed!"
        exit 1
    fi

    # Make executable
    chmod +x "${build_dir}/${BINARY_NAME}"

    # Verify build
    if [[ ! -f "${build_dir}/${BINARY_NAME}" ]]; then
        log_error "Build output not found!"
        exit 1
    fi

    log_success "Built: ${build_dir}/${BINARY_NAME}"
}

# ============================================================================
# INSTALLATION
# ============================================================================

install_to_user_dir() {
    log_info "Installing to user directory: ${INSTALL_DIR_USER}"

    # Create directory if needed
    mkdir -p "${INSTALL_DIR_USER}"

    # Remove old installation if exists
    if [[ -f "${INSTALL_DIR_USER}/${BINARY_NAME}" ]]; then
        rm -f "${INSTALL_DIR_USER}/${BINARY_NAME}"
    fi

    # Copy binary
    cp "${VIC_SOURCE_DIR}/dist/${BINARY_NAME}" "${INSTALL_DIR_USER}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR_USER}/${BINARY_NAME}"

    log_success "Installed to: ${INSTALL_DIR_USER}/${BINARY_NAME}"
}

install_to_system_dir() {
    log_info "Installing to system directory: ${INSTALL_DIR_SYSTEM}"

    if ! has_sudo; then
        log_info "Sudo required. Prompting..."
        sudo mkdir -p "${INSTALL_DIR_SYSTEM}" 2>/dev/null || true
    else
        sudo mkdir -p "${INSTALL_DIR_SYSTEM}" 2>/dev/null || true
    fi

    # Remove old installation if exists
    if [[ -f "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}" ]]; then
        if has_sudo; then
            sudo rm -f "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
        else
            rm -f "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
        fi
    fi

    # Copy binary
    if has_sudo; then
        sudo cp "${VIC_SOURCE_DIR}/dist/${BINARY_NAME}" "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
    else
        log_error "Cannot write to ${INSTALL_DIR_SYSTEM} without sudo."
        exit 1
    fi

    log_success "Installed to: ${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
}

update_path() {
    log_info "Checking PATH configuration..."

    local shell_rc=""
    local shell_config=""

    # Detect shell
    case "${SHELL:-}" in
        */zsh)  shell_rc="${HOME}/.zshrc" ;;
        */bash) shell_rc="${HOME}/.bashrc" ;;
        */fish) shell_rc="${HOME}/.config/fish/config.fish" ;;
        *)      shell_rc="${HOME}/.profile" ;;
    esac

    # Check if already in PATH
    if [[ ":${PATH}:" == *":${INSTALL_DIR_USER}:"* ]]; then
        log_info "PATH already contains: ${INSTALL_DIR_USER}"
        return 0
    fi

    log_info "Adding ${INSTALL_DIR_USER} to PATH in ${shell_rc}"

    # Backup
    if [[ -f "${shell_rc}" ]]; then
        cp "${shell_rc}" "${shell_rc}.bak"
        log_info "Backed up: ${shell_rc} -> ${shell_rc}.bak"
    fi

    # Add to PATH
    case "${SHELL:-}" in
        */fish)
            echo "" >> "${shell_rc}"
            echo "# Added by vic install script" >> "${shell_rc}"
            echo "set -gx PATH ${INSTALL_DIR_USER} \$PATH" >> "${shell_rc}"
            ;;
        *)
            echo "" >> "${shell_rc}"
            echo "# Added by vic install script" >> "${shell_rc}"
            echo 'export PATH="'${INSTALL_DIR_USER}':$PATH"' >> "${shell_rc}"
            ;;
    esac

    log_success "Updated: ${shell_rc}"
    log_warning "Please restart your shell or run: source ${shell_rc}"
}

update_precommit_config() {
    log_info "Checking pre-commit configuration..."

    if [[ ! -f "${PRECOMMIT_CONFIG}" ]]; then
        log_warning "No .pre-commit-config.yaml found. Skipping pre-commit update."
        return 0
    fi

    # Check if local vic entry already exists
    if grep -q "local" "${PRECOMMIT_CONFIG}" && grep -q "vic-gate-check" "${PRECOMMIT_CONFIG}"; then
        # Check if using local path
        if grep -qE "^[[:space:]]+entry:.*vic.*gate check" "${PRECOMMIT_CONFIG}"; then
            log_info "pre-commit already configured to use local vic"
            return 0
        fi
    fi

    log_info "pre-commit config found. Local vic-go will be used if pre-commit is installed."

    # The .pre-commit-config.yaml already uses relative paths (./cmd/vic-go/vic)
    # So no modification needed if running from repo root
    if [[ -f "${VIC_SOURCE_DIR}/dist/vic" ]]; then
        log_success "Local vic binary available for pre-commit hooks"
    fi
}

# ============================================================================
# UNINSTALL
# ============================================================================

uninstall() {
    log_info "Uninstalling vic..."

    local uninstalled=false

    # Remove from user directory
    if [[ -f "${INSTALL_DIR_USER}/${BINARY_NAME}" ]]; then
        rm -f "${INSTALL_DIR_USER}/${BINARY_NAME}"
        log_success "Removed: ${INSTALL_DIR_USER}/${BINARY_NAME}"
        uninstalled=true
    fi

    # Remove from system directory
    if [[ -f "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}" ]]; then
        if has_sudo; then
            sudo rm -f "${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
            log_success "Removed: ${INSTALL_DIR_SYSTEM}/${BINARY_NAME}"
            uninstalled=true
        else
            log_warning "Cannot remove ${INSTALL_DIR_SYSTEM}/${BINARY_NAME} without sudo"
        fi
    fi

    # Remove from build directory
    if [[ -f "${VIC_SOURCE_DIR}/dist/${BINARY_NAME}" ]]; then
        rm -f "${VIC_SOURCE_DIR}/dist/${BINARY_NAME}"
        log_info "Removed build artifact: ${VIC_SOURCE_DIR}/dist/${BINARY_NAME}"
    fi

    if [[ -d "${VIC_SOURCE_DIR}/dist" ]] && [[ -z "$(ls -A "${VIC_SOURCE_DIR}/dist" 2>/dev/null)" ]]; then
        rmdir "${VIC_SOURCE_DIR}/dist" 2>/dev/null || true
    fi

    if [[ "${uninstalled}" == "true" ]]; then
        log_success "vic has been uninstalled."
        echo ""
        echo "Note: PATH modifications in ~/.bashrc, ~/.zshrc, etc. were not removed."
        echo "      You may manually remove the added lines if desired."
    else
        log_warning "vic was not found in standard locations."
    fi
}

# ============================================================================
# MAIN
# ============================================================================

show_help() {
    cat << 'EOF'
vic Install Script
==================

Usage: ./install.sh [OPTIONS]

Options:
  -h, --help     Show this help message
  -u, --uninstall  Uninstall vic
  -l, --local    Install to ~/.local/bin (default if no sudo)
  -s, --system   Install to /usr/local/bin (requires sudo)
  -n, --no-precommit  Skip pre-commit configuration

Examples:
  ./install.sh           # Interactive install
  ./install.sh --local   # Install to ~/.local/bin
  ./install.sh --system  # Install to /usr/local/bin
  ./install.sh --uninstall  # Remove vic

Requires: Go 1.21+
EOF
}

main() {
    local install_mode="auto"
    local skip_precommit=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -u|--uninstall)
                uninstall
                exit 0
                ;;
            -l|--local)
                install_mode="local"
                ;;
            -s|--system)
                install_mode="system"
                ;;
            -n|--no-precommit)
                skip_precommit=true
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
        shift
    done

    # Detect environment
    OS="$(detect_os)"
    ARCH="$(detect_arch)"
    log_info "Detected: ${OS} (${ARCH})"

    if [[ "${OS}" == "windows" ]]; then
        log_warning "Windows detected. For best experience, use WSL or Git Bash."
        log_info "The script will attempt to proceed anyway..."
    fi

    # Run checks
    check_go
    check_git

    # Verify source
    verify_source

    # Build
    build_vic

    # Determine install location
    if [[ "${install_mode}" == "auto" ]]; then
        # Try system first if we have sudo, otherwise local
        if has_sudo && prompt_yes_no "Install to system directory ${INSTALL_DIR_SYSTEM}?" "no"; then
            install_mode="system"
        else
            install_mode="local"
        fi
    fi

    # Install
    case "${install_mode}" in
        system)
            install_to_system_dir
            ;;
        local|*)
            install_to_user_dir
            update_path
            ;;
    esac

    # Pre-commit config
    if [[ "${skip_precommit}" == "false" ]]; then
        check_precommit || true
        update_precommit_config
    fi

    # Show success message
    echo ""
    echo "============================================================"
    log_success "Installation complete!"
    echo "============================================================"
    echo ""
    echo "Next steps:"
    echo ""
    echo "  1. Restart your shell or source your config:"
    if [[ "${install_mode}" == "local" ]]; then
        echo "     source ~/.bashrc   # or ~/.zshrc"
    fi
    echo ""
    echo "  2. Verify installation:"
    echo "     which vic"
    echo "     vic --version"
    echo ""
    echo "  3. Initialize a project:"
    echo "     vic init --name \"My Project\" --tech \"Go,PostgreSQL\""
    echo ""
    echo "  4. (Optional) Install pre-commit hooks:"
    if ! command_exists pre-commit; then
        echo "     pip install pre-commit"
    fi
    echo "     pre-commit install"
    echo "     pre-commit install --hook-type commit-msg"
    echo ""
    echo "============================================================"
    echo ""
}

# Run
main "$@"
