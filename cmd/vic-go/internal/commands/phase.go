package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
)

// Phase names
var phaseNames = map[int]string{
	0: "需求凝固",
	1: "架构设计",
	2: "代码实现",
	3: "验证发布",
}

// NewPhaseCmd creates the phase command
func NewPhaseCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase",
		Short: "Manage development phases",
		Long:  `Manage development phases and view current phase status.`,
		Example: `  vic phase status           # Show current phase
  vic phase advance --to 1    # Advance to phase 1
  vic phase check             # Check phase requirements`,
	}

	cmd.AddCommand(NewPhaseStatusCmd(cfg))
	cmd.AddCommand(NewPhaseAdvanceCmd(cfg))
	cmd.AddCommand(NewPhaseCheckCmd(cfg))

	return cmd
}

// NewPhaseStatusCmd creates the phase status subcommand
func NewPhaseStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show current phase status",
		Long:    `Show the current development phase and its status.`,
		Example: `  vic phase status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPhaseStatus(cfg)
		},
	}
}

func runPhaseStatus(cfg *config.Config) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	fmt.Println("📍 Current Phase Status")
	fmt.Println("========================================")
	fmt.Printf("Cycle: %s\n", phaseFile.CycleID)
	fmt.Printf("Current Phase: %d - %s\n", phaseFile.CurrentPhase, phaseNames[phaseFile.CurrentPhase])
	fmt.Printf("Current Gate: Gate %d\n", phaseFile.CurrentGate)
	fmt.Printf("Started: %s\n", phaseFile.StartedAt)

	// Show each phase
	fmt.Println("\n📊 Phase Details:")
	for i := 0; i <= 3; i++ {
		phase, ok := phaseFile.Phases[i]
		if !ok {
			continue
		}
		icon := "⏳"
		if phase.Status == "completed" {
			icon = "✅"
		} else if phase.Status == "in_progress" {
			icon = "🔄"
		}

		fmt.Printf("\n%s Phase %d: %s\n", icon, i, phase.Name)
		fmt.Printf("   Status: %s\n", phase.Status)
		fmt.Printf("   Completion: %d%%\n", phase.Completion)

		if phase.StartedAt != "" {
			fmt.Printf("   Started: %s\n", phase.StartedAt)
		}
		if phase.CompletedAt != "" {
			fmt.Printf("   Completed: %s\n", phase.CompletedAt)
		}

		// Show gates
		fmt.Printf("   Gates: ")
		for g := 0; g < 2; g++ {
			gateKey := fmt.Sprintf("gate_%d", i*2+g)
			if gate, ok := phase.Gates[gateKey]; ok {
				if gate.Status == "passed" {
					fmt.Printf("✅ Gate%d ", i*2+g)
				} else if gate.Status == "failed" {
					fmt.Printf("❌ Gate%d ", i*2+g)
				} else {
					fmt.Printf("⏳ Gate%d ", i*2+g)
				}
			}
		}
		fmt.Println()
	}

	return nil
}

// NewPhaseAdvanceCmd creates the phase advance subcommand
func NewPhaseAdvanceCmd(cfg *config.Config) *cobra.Command {
	var toPhase int
	var force bool

	cmd := &cobra.Command{
		Use:   "advance",
		Short: "Advance to next phase",
		Long:  `Advance to a specific phase. Validates requirements before advancing.`,
		Example: `  vic phase advance --to 1    # Advance to phase 1
  vic phase advance -f          # Force advance without validation`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPhaseAdvance(cfg, toPhase, force)
		},
	}

	cmd.Flags().IntVarP(&toPhase, "to", "t", -1, "Phase to advance to (0-3)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force advance without validation")

	return cmd
}

func runPhaseAdvance(cfg *config.Config, toPhase int, force bool) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	currentPhase := phaseFile.CurrentPhase

	// Determine target phase
	if toPhase < 0 {
		toPhase = currentPhase + 1
	}

	if toPhase > 3 {
		return fmt.Errorf("already at final phase (Phase 3)")
	}

	if toPhase <= currentPhase {
		return fmt.Errorf("already at or past phase %d", toPhase)
	}

	// Validate gates before advancing
	if !force {
		fmt.Println("🔒 Validating gates before advancing...")
		currentPhaseData, ok := phaseFile.Phases[currentPhase]
		if ok {
			for g := 0; g < 2; g++ {
				gateKey := fmt.Sprintf("gate_%d", currentPhase*2+g)
				if gate, ok := currentPhaseData.Gates[gateKey]; ok {
					if gate.Status != "passed" {
						return fmt.Errorf("Gate %d (%s) must pass before advancing", currentPhase*2+g, gate.Name)
					}
				}
			}
		}
		fmt.Println("✅ All required gates passed")
	}

	// Advance phase
	now := time.Now().Format("2006-01-02")

	// Mark current phase as completed
	if currentPhaseData, ok := phaseFile.Phases[currentPhase]; ok {
		currentPhaseData.Status = "completed"
		currentPhaseData.CompletedAt = now
		currentPhaseData.Completion = 100
		phaseFile.Phases[currentPhase] = currentPhaseData
	}

	// Mark next phase as in_progress
	if nextPhaseData, ok := phaseFile.Phases[toPhase]; ok {
		nextPhaseData.Status = "in_progress"
		nextPhaseData.StartedAt = now
		nextPhaseData.Completion = 0
		phaseFile.Phases[toPhase] = nextPhaseData
	}

	// Update current phase
	phaseFile.CurrentPhase = toPhase
	phaseFile.CurrentGate = toPhase * 2
	phaseFile.LastUpdated = now

	// Save
	if err := utils.SavePhaseStatus(cfg, phaseFile); err != nil {
		return fmt.Errorf("failed to save phase status: %w", err)
	}

	fmt.Printf("\n✅ Advanced from Phase %d to Phase %d (%s)\n", currentPhase, toPhase, phaseNames[toPhase])
	fmt.Println("   Use 'vic phase status' to view details")

	return nil
}

// NewPhaseCheckCmd creates the phase check subcommand
func NewPhaseCheckCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Check phase requirements",
		Long:    `Check if current phase requirements are met.`,
		Example: `  vic phase check`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPhaseCheck(cfg)
		},
	}
}

func runPhaseCheck(cfg *config.Config) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	currentPhase := phaseFile.CurrentPhase

	fmt.Printf("🔍 Checking Phase %d: %s\n", currentPhase, phaseNames[currentPhase])
	fmt.Println("========================================")

	// Check current phase data
	phase, ok := phaseFile.Phases[currentPhase]
	if !ok {
		return fmt.Errorf("phase %d not found", currentPhase)
	}

	// Check outputs
	fmt.Println("\n📦 Required Outputs:")
	for _, output := range phase.OutputsRequired {
		fmt.Printf("   - %s\n", output)
	}

	// Check gates
	fmt.Println("\n🚪 Gate Status:")
	for g := 0; g < 2; g++ {
		gateKey := fmt.Sprintf("gate_%d", currentPhase*2+g)
		if gate, ok := phase.Gates[gateKey]; ok {
			icon := "⏳"
			if gate.Status == "passed" {
				icon = "✅"
			} else if gate.Status == "failed" {
				icon = "❌"
			}
			fmt.Printf("   %s Gate %d: %s (%s)\n", icon, currentPhase*2+g, gate.Name, gate.Status)
		}
	}

	return nil
}

// InitializePhaseFile creates initial phase file
func InitializePhaseFile(cfg *config.Config, cycleName string) error {
	now := time.Now().Format("2006-01-02")
	cycleID := fmt.Sprintf("cycle-%s", now)

	// Create default phases
	phases := make(map[int]types.Phase)
	for i := 0; i <= 3; i++ {
		phases[i] = types.Phase{
			Name:            phaseNames[i],
			Status:          "pending",
			StartedAt:       "",
			CompletedAt:     "",
			Completion:      0,
			OutputsRequired: []string{},
			Gates:           make(map[string]types.Gate),
		}
	}

	// Add gates
	for i := 0; i <= 3; i++ {
		for g := 0; g < 2; g++ {
			gateNum := i*2 + g
			gateName := ""
			switch gateNum {
			case 0:
				gateName = "需求完整性"
			case 1:
				gateName = "需求可测试"
			case 2:
				gateName = "架构完整性"
			case 3:
				gateName = "技术选型合理"
			case 4:
				gateName = "代码可编译"
			case 5:
				gateName = "代码对齐SPEC"
			case 6:
				gateName = "功能测试通过"
			case 7:
				gateName = "发布就绪"
			}
			phases[i].Gates[fmt.Sprintf("gate_%d", gateNum)] = types.Gate{
				Name:   gateName,
				Status: "pending",
				Phase:  i,
			}
		}
	}

	// Set first phase to in_progress
	phases[0] = types.Phase{
		Name:            phaseNames[0],
		Status:          "in_progress",
		StartedAt:       now,
		CompletedAt:     "",
		Completion:      0,
		OutputsRequired: []string{},
		Gates:           phases[0].Gates,
	}

	phaseFile := types.PhaseFile{
		CycleID:      cycleID,
		CycleName:    cycleName,
		CurrentPhase: 0,
		CurrentGate:  0,
		StartedAt:    now,
		LastUpdated:  now,
		Phases:       phases,
	}

	return utils.SavePhaseStatus(cfg, &phaseFile)
}
