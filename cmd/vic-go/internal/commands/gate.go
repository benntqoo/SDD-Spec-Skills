package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/utils"
)

// NewGateCmd creates the gate command
func NewGateCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gate",
		Short: "Manage gate checks",
		Long:  `Manage gate checks and validate phase requirements.`,
		Example: `  vic gate status              # Show all gates
  vic gate pass --gate 0       # Pass gate 0
  vic gate check --phase 0     # Check phase 0 gates`,
	}

	cmd.AddCommand(NewGateStatusCmd(cfg))
	cmd.AddCommand(NewGatePassCmd(cfg))
	cmd.AddCommand(NewGateCheckCmd(cfg))

	return cmd
}

// NewGateStatusCmd creates the gate status subcommand
func NewGateStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show all gate status",
		Long:    `Show the status of all gates.`,
		Example: `  vic gate status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGateStatus(cfg)
		},
	}
}

func runGateStatus(cfg *config.Config) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	fmt.Println("🚪 Gate Status")
	fmt.Println("========================================")

	// Show all gates
	for i := 0; i <= 3; i++ {
		phase, ok := phaseFile.Phases[i]
		if !ok {
			continue
		}

		fmt.Printf("\n📍 Phase %d: %s\n", i, phase.Name)

		for g := 0; g < 2; g++ {
			gateKey := fmt.Sprintf("gate_%d", i*2+g)
			gate, ok := phase.Gates[gateKey]
			if !ok {
				continue
			}

			icon := "⏳"
			if gate.Status == "passed" {
				icon = "✅"
			} else if gate.Status == "failed" {
				icon = "❌"
			}

			fmt.Printf("   %s Gate %d: %s", icon, i*2+g, gate.Name)
			if gate.Status == "passed" {
				fmt.Printf(" (checked: %s)\n", gate.CheckedAt)
			} else {
				fmt.Println()
			}
		}
	}

	return nil
}

// NewGatePassCmd creates the gate pass subcommand
func NewGatePassCmd(cfg *config.Config) *cobra.Command {
	var gateNum int
	var notes string

	cmd := &cobra.Command{
		Use:   "pass",
		Short: "Mark a gate as passed",
		Long:  `Mark a specific gate as passed with optional notes.`,
		Example: `  vic gate pass --gate 0
  vic gate pass --gate 1 --notes "Requirements complete"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGatePass(cfg, gateNum, notes)
		},
	}

	cmd.Flags().IntVarP(&gateNum, "gate", "g", -1, "Gate number (0-7)")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes")

	return cmd
}

func runGatePass(cfg *config.Config, gateNum int, notes string) error {
	if gateNum < 0 || gateNum > 7 {
		return fmt.Errorf("gate number must be between 0 and 7")
	}

	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	// Calculate phase and gate index
	phaseIndex := gateNum / 2
	gateKey := fmt.Sprintf("gate_%d", gateNum)

	phase, ok := phaseFile.Phases[phaseIndex]
	if !ok {
		return fmt.Errorf("phase %d not found", phaseIndex)
	}

	gate, ok := phase.Gates[gateKey]
	if !ok {
		return fmt.Errorf("gate %d not found", gateNum)
	}

	// Update gate status
	now := time.Now().Format("2006-01-02")
	gate.Status = "passed"
	gate.CheckedAt = now
	gate.CheckedBy = "sisyphus"
	if notes != "" {
		gate.Notes = notes
	}

	// Save back
	phase.Gates[gateKey] = gate
	phaseFile.Phases[phaseIndex] = phase

	if err := utils.SavePhaseStatus(cfg, phaseFile); err != nil {
		return fmt.Errorf("failed to save phase status: %w", err)
	}

	fmt.Printf("✅ Gate %d (%s) marked as PASSED\n", gateNum, gate.Name)
	if notes != "" {
		fmt.Printf("   Notes: %s\n", notes)
	}

	return nil
}

// NewGateCheckCmd creates the gate check subcommand
func NewGateCheckCmd(cfg *config.Config) *cobra.Command {
	var phaseNum int

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check and validate gates",
		Long:  `Check gates for a specific phase and validate requirements.`,
		Example: `  vic gate check --phase 0
  vic gate check            # Check current phase`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGateCheck(cfg, phaseNum)
		},
	}

	cmd.Flags().IntVarP(&phaseNum, "phase", "p", -1, "Phase number (0-3), -1 for current")

	return cmd
}

func runGateCheck(cfg *config.Config, phaseNum int) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to load phase status: %w", err)
	}

	// Determine which phase to check
	if phaseNum < 0 {
		phaseNum = phaseFile.CurrentPhase
	}

	if phaseNum < 0 || phaseNum > 3 {
		return fmt.Errorf("phase number must be between 0 and 3")
	}

	phase, ok := phaseFile.Phases[phaseNum]
	if !ok {
		return fmt.Errorf("phase %d not found", phaseNum)
	}

	fmt.Printf("🔍 Checking Phase %d: %s\n", phaseNum, phase.Name)
	fmt.Println("========================================")

	// Check each gate
	allPassed := true
	for g := 0; g < 2; g++ {
		gateKey := fmt.Sprintf("gate_%d", phaseNum*2+g)
		gate, ok := phase.Gates[gateKey]
		if !ok {
			continue
		}

		icon := "⏳"
		if gate.Status == "passed" {
			icon = "✅"
		} else if gate.Status == "failed" {
			icon = "❌"
			allPassed = false
		} else {
			allPassed = false
		}

		fmt.Printf("\n%s Gate %d: %s\n", icon, phaseNum*2+g, gate.Name)
		fmt.Printf("   Status: %s\n", gate.Status)
		if gate.CheckedAt != "" {
			fmt.Printf("   Checked: %s\n", gate.CheckedAt)
		}
		if gate.Notes != "" {
			fmt.Printf("   Notes: %s\n", gate.Notes)
		}
	}

	fmt.Println("\n========================================")
	if allPassed {
		fmt.Println("✅ All gates passed - can advance to next phase")
		fmt.Printf("   Run: vic phase advance --to %d\n", phaseNum+1)
	} else {
		fmt.Println("❌ Some gates not passed yet")
		fmt.Println("   Run: vic gate pass --gate <number> to mark as passed")
	}

	return nil
}

// ValidateGateCheck performs automatic gate validation
// Returns true if all gates for the phase are passed
func ValidateGateCheck(cfg *config.Config, phaseNum int) (bool, error) {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return false, err
	}

	phase, ok := phaseFile.Phases[phaseNum]
	if !ok {
		return false, fmt.Errorf("phase %d not found", phaseNum)
	}

	// Check all gates
	for g := 0; g < 2; g++ {
		gateKey := fmt.Sprintf("gate_%d", phaseNum*2+g)
		gate, ok := phase.Gates[gateKey]
		if !ok {
			continue
		}
		if gate.Status != "passed" {
			return false, nil
		}
	}

	return true, nil
}

// AutoValidateAndAdvance automatically validates gates and advances if all passed
func AutoValidateAndAdvance(cfg *config.Config, toPhase int) error {
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		return err
	}

	_ = phaseFile // used for validation

	// Validate all previous phases
	for i := 0; i < toPhase; i++ {
		passed, err := ValidateGateCheck(cfg, i)
		if err != nil {
			return err
		}
		if !passed {
			return fmt.Errorf("phase %d gates not all passed", i)
		}
	}

	// Advance phase
	return runPhaseAdvance(cfg, toPhase, true)
}
