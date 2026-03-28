package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"gopkg.in/yaml.v3"
)

// ============================================
// TDD Session Types
// ============================================

// TDDSession represents a TDD session state
type TDDSession struct {
	Active       bool   `yaml:"active"`
	CurrentPhase string `yaml:"current_phase"` // red, green, refactor
	Cycles       int    `yaml:"cycles"`
	TestsAdded   int    `yaml:"tests_added"`
	TestsPassed  int    `yaml:"tests_passed"`
	CurrentTest  string `yaml:"current_test"`
	Feature      string `yaml:"feature"`
}

// TDDCheckpoint represents a TDD checkpoint
type TDDCheckpoint struct {
	ID       string `yaml:"id"`
	Phase    string `yaml:"phase"`
	TestName string `yaml:"test_name"`
	Passed   bool   `yaml:"passed"`
	Notes    string `yaml:"notes"`
}

// TDDHistory represents the history file
type TDDHistory struct {
	Checkpoints []TDDCheckpoint `yaml:"checkpoints"`
}

// ============================================
// TDD Command
// ============================================

// NewTddCmd creates the tdd command
func NewTddCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tdd",
		Short: "TDD (Test-Driven Development) mode",
		Long:  `Enforce TDD methodology with Red-Green-Refactor cycle.`,
		Example: `  vic tdd start --feature "user login"
  vic tdd red    # Write failing test
  vic tdd green  # Make it pass
  vic tdd refactor
  vic tdd status
  vic tdd checkpoint --note "Login test passing"
  vic tdd history`,
	}

	cmd.AddCommand(NewTddStartCmd(cfg))
	cmd.AddCommand(NewTddRedCmd(cfg))
	cmd.AddCommand(NewTddGreenCmd(cfg))
	cmd.AddCommand(NewTddRefactorCmd(cfg))
	cmd.AddCommand(NewTddStatusCmd(cfg))
	cmd.AddCommand(NewTddCheckpointCmd(cfg))
	cmd.AddCommand(NewTddHistoryCmd(cfg))

	return cmd
}

// NewTddStartCmd starts a new TDD session
func NewTddStartCmd(cfg *config.Config) *cobra.Command {
	var feature string

	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start TDD session",
		Long:    `Initialize a new TDD session for a feature.`,
		Example: `  vic tdd start --feature "user login"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddStart(cfg, feature)
		},
	}

	cmd.Flags().StringVar(&feature, "feature", "", "Feature name for TDD session (required)")
	cmd.MarkFlagRequired("feature")

	return cmd
}

func runTddStart(cfg *config.Config, feature string) error {
	session := TDDSession{
		Active:       true,
		CurrentPhase: "red",
		Cycles:       0,
		TestsAdded:   0,
		TestsPassed:  0,
		Feature:      feature,
	}

	// Save session
	if err := saveTDDSession(cfg, session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🚀 TDD Session Started")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Feature: %s\n", feature)
	fmt.Println("")
	fmt.Println("   Current Phase: 🔴 RED")
	fmt.Println("")
	fmt.Println("   📋 TDD Cycle:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  🔴 RED    → Write failing test first                   │")
	fmt.Println("   │  🟢 GREEN  → Write minimal code to pass                  │")
	fmt.Println("   │  🔵 REFACTOR → Improve code structure                   │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("💡 Next: vic tdd red --test \"should validate email\"")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddRedCmd transitions to RED phase
func NewTddRedCmd(cfg *config.Config) *cobra.Command {
	var testName string

	cmd := &cobra.Command{
		Use:     "red",
		Short:   "RED phase - write failing test",
		Long:    `Enter RED phase: write a failing test before implementation.`,
		Example: `  vic tdd red --test "should validate email"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddRed(cfg, testName)
		},
	}

	cmd.Flags().StringVar(&testName, "test", "", "Test name/description")
	cmd.MarkFlagRequired("test")

	return cmd
}

func runTddRed(cfg *config.Config, testName string) error {
	session, err := loadTDDSession(cfg)
	if err != nil {
		fmt.Println("❌ No active TDD session. Run 'vic tdd start' first.")
		return nil
	}

	session.CurrentPhase = "red"
	session.CurrentTest = testName
	saveTDDSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🔴 RED Phase")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Test: %s\n", testName)
	fmt.Println("")
	fmt.Println("   ⚠️  IRON RULE: Write the test BEFORE implementation")
	fmt.Println("")
	fmt.Println("   Checkpoints:")
	fmt.Println("   1. Write test that describes desired behavior")
	fmt.Println("   2. Run test → it MUST fail")
	fmt.Println("   3. Verify failure message is clear")
	fmt.Println("")
	fmt.Printf("   ✅ When test fails: vic tdd green --test \"%s\" --passed\n", testName)
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddGreenCmd transitions to GREEN phase
func NewTddGreenCmd(cfg *config.Config) *cobra.Command {
	var testName string
	var passed bool

	cmd := &cobra.Command{
		Use:     "green",
		Short:   "GREEN phase - make test pass",
		Long:    `Enter GREEN phase: write minimal code to pass the test.`,
		Example: `  vic tdd green --test "should validate email" --passed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddGreen(cfg, testName, passed)
		},
	}

	cmd.Flags().StringVar(&testName, "test", "", "Test name/description")
	cmd.Flags().BoolVar(&passed, "passed", false, "Test passed successfully")
	cmd.MarkFlagRequired("test")

	return cmd
}

func runTddGreen(cfg *config.Config, testName string, passed bool) error {
	session, err := loadTDDSession(cfg)
	if err != nil {
		fmt.Println("❌ No active TDD session. Run 'vic tdd start' first.")
		return nil
	}

	session.CurrentPhase = "green"
	session.Cycles++
	if passed {
		session.TestsPassed++
		session.TestsAdded++
	}
	saveTDDSession(cfg, session)

	// Create checkpoint
	checkpoint := TDDCheckpoint{
		ID:       fmt.Sprintf("CP-%03d", session.Cycles),
		Phase:    "green",
		TestName: testName,
		Passed:   passed,
	}
	saveTDDCheckpoint(cfg, checkpoint)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🟢 GREEN Phase")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Test: %s\n", testName)
	fmt.Printf("   Cycles completed: %d\n", session.Cycles)
	if passed {
		fmt.Println("")
		fmt.Println("   ✅ Test passed!")
		fmt.Println("")
		fmt.Println("   Next: vic tdd refactor (or vic tdd red for next test)")
	} else {
		fmt.Println("")
		fmt.Println("   ⚠️  Still in progress - write implementation")
	}
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddRefactorCmd transitions to REFACTOR phase
func NewTddRefactorCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "refactor",
		Short:   "REFACTOR phase - improve code",
		Long:    `Enter REFACTOR phase: improve code structure without changing behavior.`,
		Example: `  vic tdd refactor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddRefactor(cfg)
		},
	}
}

func runTddRefactor(cfg *config.Config) error {
	session, err := loadTDDSession(cfg)
	if err != nil {
		fmt.Println("❌ No active TDD session. Run 'vic tdd start' first.")
		return nil
	}

	session.CurrentPhase = "refactor"
	saveTDDSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🔵 REFACTOR Phase")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("   Focus Areas:")
	fmt.Println("   1. Remove duplication")
	fmt.Println("   2. Improve naming")
	fmt.Println("   3. Extract methods if needed")
	fmt.Println("   4. Keep it simple")
	fmt.Println("")
	fmt.Println("   ⚠️  Rule: Don't change behavior - only improve structure")
	fmt.Println("")
	fmt.Println("   ✅ When done: vic tdd red --test \"next test name\"")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddStatusCmd shows TDD status
func NewTddStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show TDD session status",
		Long:    `Display current TDD session status and progress.`,
		Example: `  vic tdd status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddStatus(cfg)
		},
	}
}

func runTddStatus(cfg *config.Config) error {
	session, err := loadTDDSession(cfg)
	if err != nil {
		fmt.Println("❌ No active TDD session. Run 'vic tdd start' first.")
		return nil
	}

	phaseIcon := "🔴"
	phaseDesc := "RED"
	switch session.CurrentPhase {
	case "green":
		phaseIcon = "🟢"
		phaseDesc = "GREEN"
	case "refactor":
		phaseIcon = "🔵"
		phaseDesc = "REFACTOR"
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📊 TDD Status")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Feature:    %s\n", session.Feature)
	fmt.Printf("   Phase:       %s %s\n", phaseIcon, phaseDesc)
	fmt.Printf("   Cycles:     %d\n", session.Cycles)
	fmt.Printf("   Tests Added: %d\n", session.TestsAdded)
	fmt.Printf("   Tests Passed: %d\n", session.TestsPassed)
	if session.CurrentTest != "" {
		fmt.Printf("   Current:    %s\n", session.CurrentTest)
	}
	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddCheckpointCmd saves a checkpoint
func NewTddCheckpointCmd(cfg *config.Config) *cobra.Command {
	var note string

	cmd := &cobra.Command{
		Use:     "checkpoint",
		Short:   "Save TDD checkpoint",
		Long:    `Save a checkpoint in the current TDD session.`,
		Example: `  vic tdd checkpoint --note "Login test passing"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddCheckpoint(cfg, note)
		},
	}

	cmd.Flags().StringVar(&note, "note", "", "Checkpoint note")

	return cmd
}

func runTddCheckpoint(cfg *config.Config, note string) error {
	session, err := loadTDDSession(cfg)
	if err != nil {
		fmt.Println("❌ No active TDD session. Run 'vic tdd start' first.")
		return nil
	}

	checkpoint := TDDCheckpoint{
		ID:       fmt.Sprintf("CP-%03d", session.Cycles+1),
		Phase:    session.CurrentPhase,
		TestName: session.CurrentTest,
		Passed:   session.TestsPassed > 0,
		Notes:    note,
	}
	saveTDDCheckpoint(cfg, checkpoint)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  💾 Checkpoint %s saved\n", checkpoint.ID)
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Phase:  %s\n", checkpoint.Phase)
	fmt.Printf("   Test:   %s\n", checkpoint.TestName)
	if note != "" {
		fmt.Printf("   Note:   %s\n", note)
	}
	fmt.Println("")
	fmt.Println("   Run 'vic tdd history' to see all checkpoints")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewTddHistoryCmd shows TDD history
func NewTddHistoryCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "history",
		Short:   "Show TDD history",
		Long:    `Display TDD session history and checkpoints.`,
		Example: `  vic tdd history`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTddHistory(cfg)
		},
	}
}

func runTddHistory(cfg *config.Config) error {
	history, err := loadTDDHistory(cfg)
	if err != nil || len(history.Checkpoints) == 0 {
		fmt.Println("❌ No TDD history found. Run 'vic tdd start' first.")
		return nil
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📜 TDD History")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	for _, cp := range history.Checkpoints {
		passedIcon := "❌"
		if cp.Passed {
			passedIcon = "✅"
		}
		phaseIcon := "🔴"
		if cp.Phase == "green" {
			phaseIcon = "🟢"
		} else if cp.Phase == "refactor" {
			phaseIcon = "🔵"
		}

		fmt.Printf("   %s | %s %s\n", cp.ID, phaseIcon, passedIcon)
		fmt.Printf("       Test: %s\n", cp.TestName)
		if cp.Notes != "" {
			fmt.Printf("       Note: %s\n", cp.Notes)
		}
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("   Total checkpoints: %d\n", len(history.Checkpoints))
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// ============================================
// Helper Functions
// ============================================

func loadTDDSession(cfg *config.Config) (TDDSession, error) {
	sessionFile := filepath.Join(cfg.GetVICDir(), "status", "tdd-session.yaml")
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return TDDSession{}, err
	}

	var session TDDSession
	if err := yaml.Unmarshal(data, &session); err != nil {
		return TDDSession{}, err
	}
	return session, nil
}

func saveTDDSession(cfg *config.Config, session TDDSession) error {
	sessionFile := filepath.Join(cfg.GetVICDir(), "status", "tdd-session.yaml")
	data, err := yaml.Marshal(session)
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFile, data, 0644)
}

func loadTDDHistory(cfg *config.Config) (TDDHistory, error) {
historyFile := filepath.Join(cfg.GetVICDir(), "status", "tdd-history.yaml")
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return TDDHistory{}, err
	}

	var history TDDHistory
	if err := yaml.Unmarshal(data, &history); err != nil {
		return TDDHistory{}, err
	}
	return history, nil
}

func saveTDDCheckpoint(cfg *config.Config, checkpoint TDDCheckpoint) {
historyFile := filepath.Join(cfg.GetVICDir(), "status", "tdd-history.yaml")
	history, _ := loadTDDHistory(cfg)
	history.Checkpoints = append(history.Checkpoints, checkpoint)
	data, _ := yaml.Marshal(history)
	os.WriteFile(historyFile, data, 0644)
}
