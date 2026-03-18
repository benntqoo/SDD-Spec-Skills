package commands

import (
	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// NewRootCmd creates the root command
func NewRootCmd(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vic",
		Short: "vic - VIBE-SDD CLI for AI project memory",
		Long: `vic - VIBE-SDD CLI

Simplified command-line tool for AI and humans to operate .vic-sdd/ files.

Quick Start:
    vic init                    # Initialize project
    vic record tech            # Record technical decision
    vic status                 # Show project status
    vic check                  # Check code alignment

For more information, visit: https://github.com/vic-sdd/vic
`,
		Version: "1.0.0",
	}

	// Add subcommands
	rootCmd.AddCommand(NewInitCmd(cfg))
	rootCmd.AddCommand(NewRecordCmd(cfg))
	rootCmd.AddCommand(NewCheckCmd(cfg))
	rootCmd.AddCommand(NewValidateCmd(cfg))
	rootCmd.AddCommand(NewFoldCmd(cfg))
	rootCmd.AddCommand(NewStatusCmd(cfg))
	rootCmd.AddCommand(NewSearchCmd(cfg))
	rootCmd.AddCommand(NewHistoryCmd(cfg))
	rootCmd.AddCommand(NewExportCmd(cfg))
	rootCmd.AddCommand(NewImportCmd(cfg))
	rootCmd.AddCommand(NewSpecCmd(cfg))
	rootCmd.AddCommand(NewPhaseCmd(cfg))
	rootCmd.AddCommand(NewGateCmd(cfg))

	// Add flags
	rootCmd.PersistentFlags().StringVarP(&cfg.OutputFormat, "output", "o", cfg.OutputFormat, "Output format (json, yaml, plain)")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", cfg.Verbose, "Verbose output")

	return rootCmd
}
