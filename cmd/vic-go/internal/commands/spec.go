package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/utils"
)

const specRequirementsTemplate = `# SPEC-REQUIREMENTS.md

## Project Overview

> Brief description of the project

## User Stories

- [ ] As a user, I can...
- [ ] As an admin, I can...

## Key Features

1. Feature 1
2. Feature 2
3. Feature 3

## Non-Functional Requirements

- Performance:
- Security:
- Scalability:

## Out of Scope

- Feature X
- Feature Y
`

const specArchitectureTemplate = `# SPEC-ARCHITECTURE.md

## Architecture Overview

> High-level architecture description

## System Design

### Components

- Component 1
- Component 2

### Data Flow

> How data flows through the system

## Technology Stack

| Layer | Technology |
|-------|------------|
| Frontend | |
| Backend | |
| Database | |
| Infrastructure | |

## API Design

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/resource | List resources |
| POST | /api/resource | Create resource |

## Security

- Authentication:
- Authorization:
- Data Protection:

## Open Questions

- Question 1
- Question 2
`

// NewSpecCmd creates the spec command
func NewSpecCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "SPEC management commands",
		Long:  `Manage SPEC documents (REQUIREMENTS, ARCHITECTURE).`,
	}

	cmd.AddCommand(NewSpecInitCmd(cfg))
	cmd.AddCommand(NewSpecStatusCmd(cfg))
	cmd.AddCommand(NewSpecGateCmd(cfg))
	cmd.AddCommand(NewSpecMergeCmd(cfg))

	return cmd
}

// NewSpecInitCmd initializes SPEC documents
func NewSpecInitCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Short:   "Initialize SPEC documents",
		Long:    `Initialize SPEC-REQUIREMENTS.md and SPEC-ARCHITECTURE.md files.`,
		Example: `  vic spec init`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpecInit(cfg)
		},
	}
}

func runSpecInit(cfg *config.Config) error {
	fmt.Println("📄 Initializing SPEC documents...")

	// Ensure VIC directory
	if err := cfg.EnsureVICDir(); err != nil {
		return fmt.Errorf("failed to create .vic-sdd/: %w", err)
	}

	// Create SPEC-REQUIREMENTS.md if not exists
	if !utils.FileExists(cfg.SpecRequirements) {
		content := strings.Replace(specRequirementsTemplate, "Project Overview", fmt.Sprintf("Project Overview\n\n> Generated: %s", time.Now().Format("2006-01-02")), 1)
		if err := os.WriteFile(cfg.SpecRequirements, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write SPEC-REQUIREMENTS.md: %w", err)
		}
		fmt.Printf("   ✅ Created SPEC-REQUIREMENTS.md\n")
	} else {
		fmt.Printf("   ℹ️  SPEC-REQUIREMENTS.md already exists\n")
	}

	// Create SPEC-ARCHITECTURE.md if not exists
	if !utils.FileExists(cfg.SpecArchitecture) {
		content := strings.Replace(specArchitectureTemplate, "Architecture Overview", fmt.Sprintf("Architecture Overview\n\n> Generated: %s", time.Now().Format("2006-01-02")), 1)
		if err := os.WriteFile(cfg.SpecArchitecture, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write SPEC-ARCHITECTURE.md: %w", err)
		}
		fmt.Printf("   ✅ Created SPEC-ARCHITECTURE.md\n")
	} else {
		fmt.Printf("   ℹ️  SPEC-ARCHITECTURE.md already exists\n")
	}

	fmt.Println("\n✅ SPEC documents initialized")
	fmt.Println("   Edit these files to define your project requirements and architecture")

	return nil
}

// NewSpecStatusCmd shows SPEC status
func NewSpecStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show SPEC status",
		Long:    `Show status of SPEC documents.`,
		Example: `  vic spec status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpecStatus(cfg)
		},
	}
}

func runSpecStatus(cfg *config.Config) error {
	fmt.Println("📋 SPEC Status")
	fmt.Println("========================================")

	reqExists := utils.FileExists(cfg.SpecRequirements)
	archExists := utils.FileExists(cfg.SpecArchitecture)

	if reqExists {
		fmt.Printf("   ✅ SPEC-REQUIREMENTS.md: exists\n")
	} else {
		fmt.Printf("   ❌ SPEC-REQUIREMENTS.md: not found\n")
	}

	if archExists {
		fmt.Printf("   ✅ SPEC-ARCHITECTURE.md: exists\n")
	} else {
		fmt.Printf("   ❌ SPEC-ARCHITECTURE.md: not found\n")
	}

	if !reqExists || !archExists {
		fmt.Println("\nRun 'vic spec init' to create SPEC documents")
	}

	return nil
}

// NewSpecGateCmd runs SPEC gate checks
func NewSpecGateCmd(cfg *config.Config) *cobra.Command {
	var gate int

	cmd := &cobra.Command{
		Use:   "gate [0-3]",
		Short: "Run SPEC gate check",
		Long: `Run SPEC gate validation:
  Gate 0: Requirements Completeness
  Gate 1: Architecture Completeness  
  Gate 2: Code Alignment
  Gate 3: Test Coverage`,
		Example: `  vic spec gate 0
  vic spec gate 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				fmt.Sscanf(args[0], "%d", &gate)
			}
			return runSpecGate(cfg, gate)
		},
	}

	cmd.Flags().IntVarP(&gate, "gate", "g", 0, "Gate number (0-3)")

	return cmd
}

func runSpecGate(cfg *config.Config, gate int) error {
	gateNames := []string{
		"Requirements Completeness",
		"Architecture Completeness",
		"Code Alignment",
		"Test Coverage",
	}

	if gate < 0 || gate > 3 {
		return fmt.Errorf("gate must be between 0 and 3")
	}

	fmt.Printf("🚪 Running SPEC Gate %d: %s\n", gate, gateNames[gate])
	fmt.Println("----------------------------------------")

	switch gate {
	case 0:
		return runGate0(cfg)
	case 1:
		return runGate1(cfg)
	case 2:
		return runGate2(cfg)
	case 3:
		return runGate3(cfg)
	}

	return nil
}

func runGate0(cfg *config.Config) error {
	reqExists := utils.FileExists(cfg.SpecRequirements)

	if !reqExists {
		fmt.Println("❌ SPEC-REQUIREMENTS.md not found")
		fmt.Println("\nRun 'vic spec init' to create it")
		return nil
	}

	fmt.Println("✅ SPEC-REQUIREMENTS.md exists")
	fmt.Println("\n⚠️  Full requirements validation not implemented yet")
	fmt.Println("   (Will analyze document content in future version)")

	return nil
}

func runGate1(cfg *config.Config) error {
	archExists := utils.FileExists(cfg.SpecArchitecture)

	if !archExists {
		fmt.Println("❌ SPEC-ARCHITECTURE.md not found")
		fmt.Println("\nRun 'vic spec init' to create it")
		return nil
	}

	fmt.Println("✅ SPEC-ARCHITECTURE.md exists")
	fmt.Println("\n⚠️  Full architecture validation not implemented yet")
	fmt.Println("   (Will analyze document content in future version)")

	return nil
}

func runGate2(cfg *config.Config) error {
	// This is the same as vic check
	fmt.Println("Running code alignment check...")
	fmt.Println("✅ Use 'vic check' for full code alignment")
	return nil
}

func runGate3(cfg *config.Config) error {
	fmt.Println("⚠️  Test coverage gate not implemented yet")
	fmt.Println("   (Will integrate with test framework in future version)")
	return nil
}

// NewSpecMergeCmd merges SPEC to final docs
func NewSpecMergeCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "merge",
		Short:   "Merge to final documents",
		Long:    `Merge SPEC documents to final PROJECT.md.`,
		Example: `  vic spec merge`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpecMerge(cfg)
		},
	}
}

func runSpecMerge(cfg *config.Config) error {
	fmt.Println("🔀 Merging SPEC to final documents...")

	// Load SPEC files
	reqContent := ""
	if utils.FileExists(cfg.SpecRequirements) {
		data, _ := os.ReadFile(cfg.SpecRequirements)
		reqContent = string(data)
	}

	archContent := ""
	if utils.FileExists(cfg.SpecArchitecture) {
		data, _ := os.ReadFile(cfg.SpecArchitecture)
		archContent = string(data)
	}

	// Create merged content
	merged := fmt.Sprintf(`# Project Documentation
> Generated: %s

## Requirements

%s

## Architecture

%s

---
*This document was generated from SPEC-REQUIREMENTS.md and SPEC-ARCHITECTURE.md*
`, time.Now().Format("2006-01-02"), reqContent, archContent)

	// Write to PROJECT.md
	if err := os.WriteFile(cfg.ProjectState, []byte(merged), 0644); err != nil {
		return fmt.Errorf("failed to write PROJECT.md: %w", err)
	}

	fmt.Println("✅ Merged to PROJECT.md")

	return nil
}
