package hooks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// PreCommitHookScript is the bash script template for pre-commit hook
const PreCommitHookScript = `#!/bin/bash
# VIBE-SDD Pre-Commit Hook
# This hook ensures all VIBE-SDD gates are passed before allowing commits

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "🚪 VIBE-SDD Pre-Commit Check"
echo "========================================"
echo ""

# Run vic gate check
vic gate check --blocking

# Check exit code
if [ $? -ne 0 ]; then
    echo ""
    echo "${RED}❌ Gate check failed - commit blocked${NC}"
    echo ""
    echo "To bypass (NOT recommended):"
    echo "  git commit --no-verify -m 'message'"
    echo ""
    exit 1
fi

echo "${GREEN}✅ Gate check passed - commit allowed${NC}"
echo ""
exit 0
`

// NewHooksInstallCmd creates the hooks install command
func NewHooksInstallCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Git hooks",
		Long:  `Install Git pre-commit hook to enforce VIBE-SDD gates before commits.`,
		Example: `  vic hooks install      # Install to local repo`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHooksInstall(cfg)
		},
	}

	return cmd
}

func runHooksInstall(cfg *config.Config) error {
	fmt.Println("🔧 Installing Git hooks...")
	fmt.Println()

	// Get git hooks directory
	gitDir := filepath.Join(cfg.ProjectDir, ".git")
	hooksDir := filepath.Join(gitDir, "hooks")

	// Check if this is a git repository
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (no .git directory)")
	}

	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Write pre-commit hook
	preCommitPath := filepath.Join(hooksDir, "pre-commit")
	if err := os.WriteFile(preCommitPath, []byte(PreCommitHookScript), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	fmt.Printf("✅ Pre-commit hook installed: %s\n", preCommitPath)
	fmt.Println()
	fmt.Println("📋 Hook will run: vic gate check --blocking")
	fmt.Println()
	fmt.Println("⚠️  To bypass the hook (NOT recommended):")
	fmt.Println("   git commit --no-verify -m 'message'")
	fmt.Println()

	return nil
}

// NewHooksUninstallCmd creates the hooks uninstall command
func NewHooksUninstallCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Git hooks",
		Long:  `Remove Git pre-commit hook.`,
		Example: `  vic hooks uninstall    # Remove from local repo`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHooksUninstall(cfg)
		},
	}

	return cmd
}

func runHooksUninstall(cfg *config.Config) error {
	fmt.Println("🔧 Uninstalling Git hooks...")
	fmt.Println()

	// Get hooks directory
	hooksDir := filepath.Join(cfg.ProjectDir, ".git", "hooks")
	preCommitPath := filepath.Join(hooksDir, "pre-commit")

	// Check if pre-commit hook exists
	if _, err := os.Stat(preCommitPath); os.IsNotExist(err) {
		fmt.Println("ℹ️  No pre-commit hook found")
		fmt.Println()
		return nil
	}

	// Remove pre-commit hook
	if err := os.Remove(preCommitPath); err != nil {
		return fmt.Errorf("failed to remove pre-commit hook: %w", err)
	}

	fmt.Printf("✅ Pre-commit hook removed: %s\n", preCommitPath)
	fmt.Println()

	return nil
}

// NewHooksCmd creates the hooks command
func NewHooksCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage Git hooks",
		Long:  `Install and manage Git hooks for VIBE-SDD gate enforcement.`,
		Example: `  vic hooks install      # Install pre-commit hook
  vic hooks uninstall    # Remove pre-commit hook`,
	}

	cmd.AddCommand(NewHooksInstallCmd(cfg))
	cmd.AddCommand(NewHooksUninstallCmd(cfg))

	return cmd
}
