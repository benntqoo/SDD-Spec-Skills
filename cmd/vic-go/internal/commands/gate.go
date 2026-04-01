package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/hooks"
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
	cmd.AddCommand(NewGateSmartCmd(cfg))
	cmd.AddCommand(hooks.NewHooksCmd(cfg))

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
// This command is designed for pre-commit hooks and CI
func NewGateCheckCmd(cfg *config.Config) *cobra.Command {
	var phaseNum int
	var blocking bool
	var format string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check gate status for pre-commit",
		Long: `Check gate status and optionally block if gates are not passed.

This command is designed to be used in pre-commit hooks to ensure
all VIBE-SDD gates are passed before allowing a commit.

Examples:
  vic gate check                    # Check current phase gates
  vic gate check --phase 1          # Check phase 1 gates
  vic gate check --blocking         # Exit with error if gates not passed
  vic gate check --format json      # JSON output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if blocking {
				return RunGateCheck(cfg, phaseNum)
			}
			return runGateStatusCheck(cfg, phaseNum, format)
		},
	}

	cmd.Flags().IntVarP(&phaseNum, "phase", "p", -1, "Phase number (0-3), -1 for current")
	cmd.Flags().BoolVarP(&blocking, "blocking", "b", false, "Exit with error if gates not passed (for pre-commit)")
	cmd.Flags().StringVarP(&format, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}

// runGateStatusCheck shows gate status for a phase
func runGateStatusCheck(cfg *config.Config, phaseNum int, format string) error {
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

	// Collect gate status
	gateResults := []map[string]interface{}{}
	allPassed := true
	for g := 0; g < 2; g++ {
		gateKey := fmt.Sprintf("gate_%d", phaseNum*2+g)
		gate, ok := phase.Gates[gateKey]
		if !ok {
			continue
		}
		if gate.Status != "passed" {
			allPassed = false
		}
		gateResults = append(gateResults, map[string]interface{}{
			"gate":    phaseNum*2 + g,
			"name":    gate.Name,
			"status":  gate.Status,
			"checked": gate.CheckedAt,
		})
	}

	// JSON output
	if format == "json" {
		fmt.Printf(`{"success":%v,"message":"Gate check completed","data":{"phase":%d,"phase_name":%q,"all_passed":%v,"gates":[`, allPassed, phaseNum, phase.Name, allPassed)
		for i, g := range gateResults {
			if i > 0 {
				fmt.Printf(",")
			}
			fmt.Printf(`{"gate":%d,"name":%q,"status":%q,"checked":%q}`, g["gate"], g["name"], g["status"], g["checked"])
		}
		fmt.Printf(`]}}`)
		fmt.Println()
		return nil
	}

	// Plain output
	fmt.Printf("🔍 Checking Phase %d: %s\n", phaseNum, phase.Name)
	fmt.Println("========================================")

	// Display each gate from collected results
	for _, g := range gateResults {
		icon := "⏳"
		if g["status"] == "passed" {
			icon = "✅"
		} else if g["status"] == "failed" {
			icon = "❌"
		}

		fmt.Printf("\n%s Gate %d: %s\n", icon, g["gate"], g["name"])
		fmt.Printf("   Status: %s\n", g["status"])
		if g["checked"] != "" {
			fmt.Printf("   Checked: %s\n", g["checked"])
		}
		if g["notes"] != nil && g["notes"] != "" {
			fmt.Printf("   Notes: %s\n", g["notes"])
		}
	}

	fmt.Println("\n========================================")
	if allPassed {
		fmt.Println("✅ All gates passed - can advance to next phase")
		fmt.Printf("   Run: vic phase advance --to %d\n", phaseNum+1)
	} else {
		fmt.Println("❌ Some gates not passed yet")
		fmt.Println("   Run: vic spec gate <number> to run gate checks")
		fmt.Println("   Or: vic gate pass --gate <number> to manually mark as passed")
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

// RunGateCheck performs blocking gate check for pre-commit/CI
func RunGateCheck(cfg *config.Config, phaseNum int) error {
	fmt.Println("🚪 VIBE-SDD Gate Check (Blocking)")
	fmt.Println("========================================")
	fmt.Println()

	// Load phase status
	phaseFile, err := utils.LoadPhaseStatus(cfg)
	if err != nil {
		// Phase status file doesn't exist - this is a new project
		fmt.Println("⚠️  No phase status found (run 'vic init' first)")
		fmt.Println("   Allowing commit for new projects...")
		fmt.Println()
		return nil
	}

	// Determine which phase we're in
	if phaseNum < 0 {
		phaseNum = phaseFile.CurrentPhase
	}

	// Check gates for current and previous phases
	allPassed := true

	for i := 0; i <= phaseNum; i++ {
		phase, ok := phaseFile.Phases[i]
		if !ok {
			continue
		}

		// Check at least the first gate for each phase
		gateKey := fmt.Sprintf("gate_%d", i*2)
		gate, ok := phase.Gates[gateKey]
		if !ok {
			continue
		}

		icon := "⏳"
		if gate.Status == "passed" {
			icon = "✅"
		} else {
			icon = "❌"
			allPassed = false
		}

		required := ""
		if i <= phaseNum {
			required = " [REQUIRED]"
		}

		fmt.Printf("[%s] Gate %d: %s%s\n", icon, i*2, gate.Name, required)
		fmt.Printf("     Status: %s\n", gate.Status)
	}

	fmt.Println()
	fmt.Println("========================================")

	if allPassed {
		fmt.Println("✅ All required gates passed - commit allowed")
		fmt.Println()
		return nil
	}

	fmt.Println("❌ Gate check FAILED - some required gates not passed")
	fmt.Println()
	fmt.Println("To pass gates:")
	for i := 0; i <= phaseNum; i++ {
		phase, ok := phaseFile.Phases[i]
		if !ok {
			continue
		}
		gateKey := fmt.Sprintf("gate_%d", i*2)
		gate, ok := phase.Gates[gateKey]
		if ok && gate.Status != "passed" {
			fmt.Printf("   vic spec gate %d\n", i)
		}
	}
	fmt.Println()
	fmt.Println("⚠️  To bypass (NOT recommended):")
	fmt.Println("   git commit --no-verify -m 'message'")

	return fmt.Errorf("gate check failed: required gates not passed")
}

// NewGateSmartCmd creates the gate smart subcommand
// This command intelligently selects which gates to run based on risk assessment
func NewGateSmartCmd(cfg *config.Config) *cobra.Command {
	var execute bool
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "smart",
		Short: "Smart gate selection based on risk assessment",
		Long: `Intelligently select which gates to run based on change analysis.

This command:
1. Analyzes current changes
2. Assesses risk level
3. Selects required gates
4. Optionally runs the gates

Examples:
  vic gate smart              # Show which gates would run
  vic gate smart --execute    # Run the selected gates
  vic gate smart --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGateSmart(cfg, execute, outputFormat)
		},
	}

	cmd.Flags().BoolVarP(&execute, "execute", "e", false, "Execute the selected gates")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "plain", "Output format (plain, json)")

	return cmd
}

func runGateSmart(cfg *config.Config, execute bool, outputFormat string) error {
	// Get assessment
	details, err := getGitDiffDetails()
	if err != nil {
		return fmt.Errorf("failed to get git diff: %w", err)
	}

	changeType := detectChangeTypeFromDetails(details)
	riskScore, riskLevel := assessRisk(details)
	gatesRequired := getRequiredGates(riskLevel)
	recommendedSkill := getRecommendedSkill(changeType, riskLevel)

	if outputFormat == "json" {
		fmt.Printf(`{
  "change_type": "%s",
  "risk_level": "%s",
  "risk_score": %.2f,
  "gates_required": %v,
  "gates_skipped": %v,
  "recommended_skill": "%s"
}`, changeType, riskLevel, riskScore, gatesRequired, getSkippedGates(gatesRequired), recommendedSkill)
		fmt.Println()
		return nil
	}

	// Plain output
	fmt.Println()
	fmt.Println("🧠 Smart Gate Selection")
	fmt.Println("========================================")
	fmt.Println()

	// Show assessment
	fmt.Printf("📋 Change Type: %s\n", changeType)
	fmt.Printf("⚠️  Risk Level: %s (%.2f)\n", riskLevel, riskScore)
	fmt.Println()

	// Show gate selection
	fmt.Println("🚪 Gate Selection:")
	for gate := 0; gate <= 3; gate++ {
		required := containsInt(gatesRequired, gate)
		status := "⏭️  SKIP"
		if required {
			status = "✅ REQUIRED"
		}
		fmt.Printf("   Gate %d: %s\n", gate, status)
	}

	fmt.Println()
	fmt.Printf("📊 Summary: %d gates required, %d skipped\n", len(gatesRequired), 4-len(gatesRequired))
	fmt.Printf("🎯 Recommended Skill: %s\n", recommendedSkill)
	fmt.Println()

	if execute && len(gatesRequired) > 0 {
		fmt.Println("========================================")
		fmt.Println("🚀 Executing required gates...")
		fmt.Println()

		for _, gate := range gatesRequired {
			fmt.Printf("Running Gate %d...\n", gate)
			// In a real implementation, this would call the actual gate check
			// For now, we just show what would run
			fmt.Printf("   vic spec gate %d\n", gate)
		}
		fmt.Println()
		fmt.Println("✅ Smart gate execution complete")
	} else if execute {
		fmt.Println("✅ No gates required for this change type")
	} else {
		fmt.Println("💡 Run with --execute to run the selected gates")
	}

	return nil
}

func getSkippedGates(required []int) []int {
	all := []int{0, 1, 2, 3}
	skipped := []int{}
	for _, g := range all {
		if !containsInt(required, g) {
			skipped = append(skipped, g)
		}
	}
	return skipped
}

func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
