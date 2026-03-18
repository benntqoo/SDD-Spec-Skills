package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/utils"
)

var initName string
var initTech string

// NewInitCmd creates the init command
func NewInitCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize .vic-sdd/ directory",
		Long: `Initialize .vic-sdd/ directory for a project.

This creates the basic directory structure and files needed for VIBE-SDD.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cfg)
		},
	}

	cmd.Flags().StringVarP(&initName, "name", "n", "", "Project name")
	cmd.Flags().StringVarP(&initTech, "tech", "t", "", "Tech stack (comma-separated)")

	return cmd
}

func runInit(cfg *config.Config) error {
	// Ensure VIC directory exists
	if err := cfg.EnsureVICDir(); err != nil {
		return fmt.Errorf("failed to create .vic-sdd/: %w", err)
	}

	// Create subdirectories
	subDirs := []string{"status", "tech"}
	for _, dir := range subDirs {
		if err := cfg.EnsureSubDir(dir); err != nil {
			return fmt.Errorf("failed to create %s/: %w", dir, err)
		}
	}

	// Create template files
	files := map[string]string{
		cfg.EventsFile:      "events: []",
		cfg.StateFile:       "state: {}",
		cfg.TechRecordsFile: "tech_records: []",
		cfg.RiskZonesFile:   "risks: []",
	}

	for path, content := range files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", path, err)
			}
		}
	}

	// Generate project.yaml
	utils.GenerateProjectYAML(cfg, initName, initTech)

	// Generate dependency-graph.yaml
	utils.GenerateDependencyGraph(cfg)

	// Initialize SPEC documents
	runSpecInit(cfg)

	// Initialize phase status
	InitializePhaseFile(cfg, "default")

	fmt.Printf("✅ Initialized .vic-sdd/ directory\n")
	fmt.Printf("   Project: %s\n", initName)
	fmt.Printf("   Tech: %s\n", initTech)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("   vic record tech --id DB-001 --title \"Use PostgreSQL\" --decision \"Primary DB\"\n")
	fmt.Printf("   vic status\n")

	return nil
}
