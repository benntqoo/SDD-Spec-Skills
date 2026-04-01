package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
	"gopkg.in/yaml.v3"
)

// NewReplanCmd creates the replan command
func NewReplanCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replan",
		Short: "Adaptive replanning",
		Long:  `Trigger adaptive replanning when research reveals discrepancies or when reassessment is needed.`,
		Example: `  vic replan                    # Start replanning
  vic replan --trigger research  # Trigger by research
  vic replan --list             # List past replans`,
	}

	cmd.AddCommand(NewReplanTriggerCmd(cfg))
	cmd.AddCommand(NewReplanListCmd(cfg))
	cmd.AddCommand(NewReplanShowCmd(cfg))

	return cmd
}

// ============================================
// Trigger Command
// ============================================

// NewReplanTriggerCmd creates the replan trigger subcommand
func NewReplanTriggerCmd(cfg *config.Config) *cobra.Command {
	var trigger string
	var finding string

	cmd := &cobra.Command{
		Use:   "trigger",
		Short: "Trigger a replan",
		Long:  `Trigger adaptive replanning with a finding or discrepancy.`,
		Example: `  vic replan trigger --finding "Found library that does this in 10 lines"
  vic replan trigger --trigger research --finding "..."`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplanTrigger(cfg, trigger, finding)
		},
	}

	cmd.Flags().StringVar(&trigger, "trigger", "user_request", "Trigger type (research/slice/tech-surprise/env-change)")
	cmd.Flags().StringVar(&finding, "finding", "", "What was discovered")

	return cmd
}

func runReplanTrigger(cfg *config.Config, trigger, finding string) error {
	if finding == "" {
		fmt.Println("📋 Adaptive Replanning")
		fmt.Println("")
		fmt.Println("This command helps reassess the plan when new information is discovered.")
		fmt.Println("")
		fmt.Println("Use --finding to specify what triggered this replan:")
		fmt.Println("  vic replan trigger --finding \"Found a better library\"")
		fmt.Println("")
		fmt.Println("Trigger types:")
		fmt.Println("  research    - Research revealed discrepancies")
		fmt.Println("  slice       - Slice/milestone completed")
		fmt.Println("  tech-surprise - Technical discovery")
		fmt.Println("  env-change  - Environment changed")
		return nil
	}

	// Load existing replan log
	replanLog, _ := loadReplanLog(cfg)

	// Create new replan entry
	replanID := fmt.Sprintf("REPLAN-%03d", len(replanLog.ReplanHistory)+1)
	decision := types.ReplanDecision{
		ID:           replanID,
		Trigger:      types.ReplanTrigger(trigger),
		Timestamp:    time.Now(),
		Finding:      finding,
		OriginalPlan: "See SPEC-ARCHITECTURE.md",
		NewPlan:      "Assessment pending",
		Reason:       "Awaiting user guidance",
		UserApproved: false,
	}

	replanLog.ReplanHistory = append(replanLog.ReplanHistory, decision)
	replanLog.LastReplan = time.Now()

	// Save
	if err := saveReplanLog(cfg, replanLog); err != nil {
		return fmt.Errorf("failed to save replan log: %w", err)
	}

	fmt.Println("✅ Replan triggered")
	fmt.Printf("   ID: %s\n", replanID)
	fmt.Printf("   Trigger: %s\n", trigger)
	fmt.Printf("   Finding: %s\n", finding)
	fmt.Println("")
	fmt.Println("   Next steps:")
	fmt.Println("   1. Review the finding against current plan")
	fmt.Println("   2. Assess impact using 'vic replan show'")
	fmt.Println("   3. Update SPEC if changes are approved")

	return nil
}

// ============================================
// List Command
// ============================================

// NewReplanListCmd creates the replan list subcommand
func NewReplanListCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List past replans",
		Long:    `Show history of replan decisions.`,
		Example: `  vic replan list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplanList(cfg)
		},
	}
}

func runReplanList(cfg *config.Config) error {
	replanLog, err := loadReplanLog(cfg)
	if err != nil || len(replanLog.ReplanHistory) == 0 {
		fmt.Println("📋 No replans recorded yet")
		fmt.Println("")
		fmt.Println("Use 'vic replan trigger --finding \"...\"' to trigger a replan")
		return nil
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📋 Replan History")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("   Last replan: %s\n", replanLog.LastReplan.Format("2006-01-02 15:04"))
	fmt.Println("")

	for i := len(replanLog.ReplanHistory) - 1; i >= 0; i-- {
		entry := replanLog.ReplanHistory[i]
		approved := "⏳"
		if entry.UserApproved {
			approved = "✅"
		}

		triggerIcon := "📝"
		switch entry.Trigger {
		case types.ReplanTriggerResearch:
			triggerIcon = "🔬"
		case types.ReplanTriggerSliceComplete:
			triggerIcon = "📦"
		case types.ReplanTriggerTechSurprise:
			triggerIcon = "⚡"
		case types.ReplanTriggerEnvChange:
			triggerIcon = "🌍"
		}

		fmt.Printf("   %s %s %s\n", approved, triggerIcon, entry.ID)
		fmt.Printf("      Trigger: %s\n", entry.Trigger)
		fmt.Printf("      Finding: %s\n", truncate(entry.Finding, 60))
		fmt.Printf("      Time: %s\n", entry.Timestamp.Format("2006-01-02 15:04"))
		if entry.Impact.EffortLevel != "" {
			fmt.Printf("      Impact: %s effort\n", entry.Impact.EffortLevel)
		}
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Total replans: %d\n", len(replanLog.ReplanHistory))

	return nil
}

// ============================================
// Show Command
// ============================================

// NewReplanShowCmd creates the replan show subcommand
func NewReplanShowCmd(cfg *config.Config) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "show [id]",
		Short: "Show replan details",
		Long:  `Show detailed information about a replan decision.`,
		Example: `  vic replan show REPLAN-001
  vic replan show --latest`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				id = args[0]
			}
			return runReplanShow(cfg, id)
		},
	}

	cmd.Flags().Bool("latest", false, "Show latest replan")

	return cmd
}

func runReplanShow(cfg *config.Config, id string) error {
	replanLog, err := loadReplanLog(cfg)
	if err != nil || len(replanLog.ReplanHistory) == 0 {
		fmt.Println("📋 No replans recorded yet")
		return nil
	}

	var entry types.ReplanDecision
	if id == "" {
		entry = replanLog.ReplanHistory[len(replanLog.ReplanHistory)-1]
	} else {
		for _, e := range replanLog.ReplanHistory {
			if e.ID == id {
				entry = e
				break
			}
		}
		if entry.ID == "" {
			return fmt.Errorf("replan %s not found", id)
		}
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📋 Replan Detail: " + entry.ID)
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("\n   Trigger: %s\n", entry.Trigger)
	fmt.Printf("   Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Approved: %v\n", entry.UserApproved)

	fmt.Println("")
	fmt.Println("   📝 Finding:")
	fmt.Printf("   %s\n", wrapText(entry.Finding, 70))

	fmt.Println("")
	fmt.Println("   📄 Original Plan:")
	fmt.Printf("   %s\n", wrapText(entry.OriginalPlan, 70))

	fmt.Println("")
	fmt.Println("   🔄 New Plan:")
	fmt.Printf("   %s\n", wrapText(entry.NewPlan, 70))

	fmt.Println("")
	fmt.Println("   💡 Reason:")
	fmt.Printf("   %s\n", wrapText(entry.Reason, 70))

	if len(entry.Impact.ScopeChanges) > 0 {
		fmt.Println("")
		fmt.Println("   📊 Impact:")
		if entry.Impact.ScopeChanges != nil {
			fmt.Println("   Scope Changes:")
			for _, change := range entry.Impact.ScopeChanges {
				fmt.Printf("      • %s\n", change)
			}
		}
		if entry.Impact.TimelineDelta != "" {
			fmt.Printf("   Timeline: %s\n", entry.Impact.TimelineDelta)
		}
		if entry.Impact.EffortLevel != "" {
			fmt.Printf("   Effort: %s\n", entry.Impact.EffortLevel)
		}
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// ============================================
// Helper Functions
// ============================================

func loadReplanLog(cfg *config.Config) (*types.ReplanLogFile, error) {
	replanFile := filepath.Join(cfg.GetVICDir(), "status", "replan-log.yaml")

	if !utils.FileExists(replanFile) {
		return &types.ReplanLogFile{
			Version:       "1.0",
			ReplanHistory: []types.ReplanDecision{},
		}, nil
	}

	data, err := os.ReadFile(replanFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read replan log: %w", err)
	}

	var replanLog types.ReplanLogFile
	if err := yaml.Unmarshal(data, &replanLog); err != nil {
		return nil, fmt.Errorf("failed to parse replan log: %w", err)
	}

	return &replanLog, nil
}

func saveReplanLog(cfg *config.Config, replanLog *types.ReplanLogFile) error {
	replanFile := filepath.Join(cfg.GetVICDir(), "status", "replan-log.yaml")

	// Ensure directory exists
	if !utils.FileExists(cfg.ProjectDir + "/status") {
		if err := os.MkdirAll(cfg.ProjectDir+"/status", 0755); err != nil {
			return fmt.Errorf("failed to create status directory: %w", err)
		}
	}

	data, err := yaml.Marshal(replanLog)
	if err != nil {
		return fmt.Errorf("failed to marshal replan log: %w", err)
	}

	header := []byte(`# Replan Log - VIBE-SDD
# Auto-generated by vic-go
# Do not edit manually - use vic replan commands

`)
	if err := os.WriteFile(replanFile, append(header, data...), 0644); err != nil {
		return fmt.Errorf("failed to write replan log: %w", err)
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wrapText(s string, width int) string {
	// Simple word wrap
	words := ""
	for i, r := range s {
		words += string(r)
		if (i+1)%width == 0 {
			words += "\n      "
		}
	}
	return words
}
