package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"gopkg.in/yaml.v3"
)

// ============================================
// Debug Session Types
// ============================================

// DebugSession represents a debugging session
type DebugSession struct {
	Active         bool      `yaml:"active"`
	StartedAt      time.Time `yaml:"started_at"`
	Problem        string    `yaml:"problem"`
	CurrentPhase   string    `yaml:"current_phase"` // survey, pattern, hypothesis, implement
	Attempts       int       `yaml:"attempts"`
	RootCauseFound bool      `yaml:"root_cause_found"`
	RootCause      string    `yaml:"root_cause"`
	Hypothesis     string    `yaml:"hypothesis"`
}

// DebugFinding represents a finding during debugging
type DebugFinding struct {
	ID          string `yaml:"id"`
	Phase       string `yaml:"phase"`
	Description string `yaml:"description"`
	Evidence    string `yaml:"evidence"`
	Timestamp   string `yaml:"timestamp"`
}

// DebugLog represents the debug log
type DebugLog struct {
	Sessions  []DebugSession `yaml:"sessions"`
	Findings  []DebugFinding `yaml:"findings"`
	SessionID string         `yaml:"session_id"`
}

// ============================================
// Debug Command
// ============================================

// NewDebugCmd creates the debug command
func NewDebugCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Systematic debugging mode",
		Long:  `Systematic debugging with 4-phase root cause analysis.`,
		Example: `  vic debug start --problem "Login fails in production"
  vic debug survey   # Gather evidence
  vic debug pattern  # Find similar issues
  vic debug hypothesis --explain "Token expired"
  vic debug implement --fix "Added token refresh"
  vic debug status
  vic debug report`,
	}

	cmd.AddCommand(NewDebugStartCmd(cfg))
	cmd.AddCommand(NewDebugSurveyCmd(cfg))
	cmd.AddCommand(NewDebugPatternCmd(cfg))
	cmd.AddCommand(NewDebugHypothesisCmd(cfg))
	cmd.AddCommand(NewDebugImplementCmd(cfg))
	cmd.AddCommand(NewDebugStatusCmd(cfg))
	cmd.AddCommand(NewDebugReportCmd(cfg))
	cmd.AddCommand(NewDebugHistoryCmd(cfg))

	return cmd
}

// NewDebugStartCmd starts a new debug session
func NewDebugStartCmd(cfg *config.Config) *cobra.Command {
	var problem string

	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start debug session",
		Long:    `Initialize a new debugging session.`,
		Example: `  vic debug start --problem "Login fails in production"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugStart(cfg, problem)
		},
	}

	cmd.Flags().StringVar(&problem, "problem", "", "Problem description (required)")
	cmd.MarkFlagRequired("problem")

	return cmd
}

func runDebugStart(cfg *config.Config, problem string) error {
	session := DebugSession{
		Active:         true,
		StartedAt:      time.Now(),
		Problem:        problem,
		CurrentPhase:   "survey",
		Attempts:       0,
		RootCauseFound: false,
		RootCause:      "",
		Hypothesis:     "",
	}

	saveDebugSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🔍 Debug Session Started")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Problem: %s\n", problem)
	fmt.Printf("   Started: %s\n", session.StartedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("")
	fmt.Println("   4-Phase Debug Cycle:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  1️⃣  SURVEY     → Gather evidence                        │")
	fmt.Println("   │  2️⃣  PATTERN    → Find similar issues                    │")
	fmt.Println("   │  3️⃣  HYPOTHESIS → Form and test hypothesis              │")
	fmt.Println("   │  4️⃣  IMPLEMENT  → Fix root cause                        │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("⚠️  IRON RULE: Never fix without finding root cause first")
	fmt.Println("")
	fmt.Println("💡 Next: vic debug survey")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugSurveyCmd enters SURVEY phase
func NewDebugSurveyCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "survey",
		Short:   "SURVEY phase - gather evidence",
		Long:    `Phase 1: Gather evidence about the problem.`,
		Example: `  vic debug survey`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugSurvey(cfg)
		},
	}
}

func runDebugSurvey(cfg *config.Config) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	session.CurrentPhase = "survey"
	saveDebugSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  1️⃣  SURVEY Phase - Gather Evidence")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Problem: %s\n", session.Problem)
	fmt.Println("")
	fmt.Println("   📋 Survey Checklist:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  [ ] Reproduce the error                                │")
	fmt.Println("   │  [ ] Collect: logs, stack traces, environment            │")
	fmt.Println("   │  [ ] Question assumptions                                │")
	fmt.Println("   │  [ ] Identify what is NOT causing the issue             │")
	fmt.Println("   │  [ ] Note: error type, when it occurs, frequency         │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("   Questions to answer:")
	fmt.Println("   • What is the exact error message?")
	fmt.Println("   • When does it occur? (first time? intermittently?)")
	fmt.Println("   • In what environment? (prod/staging/local?)")
	fmt.Println("   • What changed recently?")
	fmt.Println("")
	fmt.Println("   ✅ When done: vic debug pattern")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugPatternCmd enters PATTERN phase
func NewDebugPatternCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "pattern",
		Short:   "PATTERN phase - find similar issues",
		Long:    `Phase 2: Search for similar patterns in codebase.`,
		Example: `  vic debug pattern`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugPattern(cfg)
		},
	}
}

func runDebugPattern(cfg *config.Config) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	session.CurrentPhase = "pattern"
	saveDebugSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  2️⃣  PATTERN Phase - Find Similar Issues")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("   🔎 Pattern Analysis:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  [ ] Search for similar issues in codebase             │")
	fmt.Println("   │  [ ] Check: error handling patterns, edge cases        │")
	fmt.Println("   │  [ ] Compare: working vs non-working code              │")
	fmt.Println("   │  [ ] Review: recent commits, dependency changes        │")
	fmt.Println("   │  [ ] Research: known issues with this library/pattern │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("   Search patterns to consider:")
	fmt.Println("   • grep -r \"error\" src/")
	fmt.Println("   • git log --oneline -10")
	fmt.Println("   • Check if library has known issues")
	fmt.Println("   • Compare with similar working features")
	fmt.Println("")
	fmt.Println("   ✅ When patterns found: vic debug hypothesis --explain \"...\"")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugHypothesisCmd enters HYPOTHESIS phase
func NewDebugHypothesisCmd(cfg *config.Config) *cobra.Command {
	var explanation string

	cmd := &cobra.Command{
		Use:     "hypothesis",
		Short:   "HYPOTHESIS phase - form and test hypothesis",
		Long:    `Phase 3: Form a testable hypothesis.`,
		Example: `  vic debug hypothesis --explain "Token expired because server time is off"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugHypothesis(cfg, explanation)
		},
	}

	cmd.Flags().StringVar(&explanation, "explain", "", "Your hypothesis (required)")
	cmd.MarkFlagRequired("explain")

	return cmd
}

func runDebugHypothesis(cfg *config.Config, explanation string) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	session.CurrentPhase = "hypothesis"
	session.Hypothesis = explanation
	session.Attempts++
	saveDebugSession(cfg, session)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  3️⃣  HYPOTHESIS Phase - Form & Test")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Hypothesis: %s\n", explanation)
	fmt.Printf("   Attempt: #%d\n", session.Attempts)
	fmt.Println("")
	fmt.Println("   🎯 Hypothesis Testing:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  1. Make ONE change                                      │")
	fmt.Println("   │  2. Test minimally                                       │")
	fmt.Println("   │  3. Verify hypothesis                                    │")
	fmt.Println("   │  4. If wrong: form new hypothesis                        │")
	fmt.Println("   │  5. If correct: proceed to implement                     │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("   ⚠️  STOP after 3 failed attempts!")
	fmt.Println("   If 3+ attempts fail → question the architecture itself")
	fmt.Println("")
	fmt.Println("   ✅ When hypothesis confirmed: vic debug implement --fix \"...\"")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugImplementCmd enters IMPLEMENT phase
func NewDebugImplementCmd(cfg *config.Config) *cobra.Command {
	var fix string
	var rootCause string

	cmd := &cobra.Command{
		Use:     "implement",
		Short:   "IMPLEMENT phase - fix root cause",
		Long:    `Phase 4: Implement the fix for root cause.`,
		Example: `  vic debug implement --fix "Added token refresh logic" --root-cause "Token expires without refresh"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugImplement(cfg, fix, rootCause)
		},
	}

	cmd.Flags().StringVar(&fix, "fix", "", "Fix description (required)")
	cmd.Flags().StringVar(&rootCause, "root-cause", "", "Root cause (optional)")
	cmd.MarkFlagRequired("fix")

	return cmd
}

func runDebugImplement(cfg *config.Config, fix, rootCause string) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	session.CurrentPhase = "implement"
	session.RootCauseFound = true
	if rootCause != "" {
		session.RootCause = rootCause
	}
	saveDebugSession(cfg, session)

	// Save finding
	finding := DebugFinding{
		ID:          fmt.Sprintf("FIND-%03d", session.Attempts),
		Phase:       "implement",
		Description: fmt.Sprintf("Fix: %s", fix),
		Evidence:    rootCause,
		Timestamp:   time.Now().Format("2006-01-02 15:04"),
	}
	saveDebugFinding(cfg, finding)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  4️⃣  IMPLEMENT Phase - Fix Root Cause")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Fix: %s\n", fix)
	if rootCause != "" {
		fmt.Printf("   Root Cause: %s\n", rootCause)
	}
	fmt.Println("")
	fmt.Println("   ✅ Implementation Complete!")
	fmt.Println("")
	fmt.Println("   📋 Final Steps:")
	fmt.Println("   ┌─────────────────────────────────────────────────────────┐")
	fmt.Println("   │  [ ] Add regression test                                 │")
	fmt.Println("   │  [ ] Document the fix                                    │")
	fmt.Println("   │  [ ] Verify fix works                                    │")
	fmt.Println("   │  [ ] Record in vic rr                                    │")
	fmt.Println("   └─────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("   Run 'vic debug report' for full session summary")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugStatusCmd shows debug status
func NewDebugStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show debug status",
		Long:    `Display current debug session status.`,
		Example: `  vic debug status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugStatus(cfg)
		},
	}
}

func runDebugStatus(cfg *config.Config) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	phaseNum := "1️⃣ "
	phaseName := "SURVEY"
	switch session.CurrentPhase {
	case "pattern":
		phaseNum = "2️⃣ "
		phaseName = "PATTERN"
	case "hypothesis":
		phaseNum = "3️⃣ "
		phaseName = "HYPOTHESIS"
	case "implement":
		phaseNum = "4️⃣ "
		phaseName = "IMPLEMENT"
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🔍 Debug Status")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Problem:  %s\n", session.Problem)
	fmt.Printf("   Started:  %s\n", session.StartedAt.Format("2006-01-02 15:04"))
	fmt.Printf("   Phase:    %s %s\n", phaseNum, phaseName)
	fmt.Printf("   Attempts: %d\n", session.Attempts)
	if session.RootCauseFound {
		fmt.Println("   Status:   ✅ Root cause found")
		if session.RootCause != "" {
			fmt.Printf("            → %s\n", session.RootCause)
		}
	} else if session.Attempts >= 3 {
		fmt.Println("   ⚠️  WARNING: 3+ attempts - consider questioning architecture")
	}
	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugReportCmd generates debug report
func NewDebugReportCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "report",
		Short:   "Generate debug report",
		Long:    `Generate a full debug session report.`,
		Example: `  vic debug report`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugReport(cfg)
		},
	}
}

func runDebugReport(cfg *config.Config) error {
	session, err := loadDebugSession(cfg)
	if err != nil {
		fmt.Println("❌ No active debug session. Run 'vic debug start' first.")
		return nil
	}

	findings := loadDebugFindings(cfg)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📋 Debug Report")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Problem:     %s\n", session.Problem)
	fmt.Printf("   Started:     %s\n", session.StartedAt.Format("2006-01-02 15:04"))
	fmt.Printf("   Duration:    %s\n", time.Since(session.StartedAt).Round(time.Second))
	fmt.Printf("   Phase:       %s\n", session.CurrentPhase)
	fmt.Printf("   Attempts:    %d\n", session.Attempts)
	fmt.Println("")

	if session.RootCauseFound {
		fmt.Println("   ✅ ROOT CAUSE FOUND")
		if session.RootCause != "" {
			fmt.Printf("   → %s\n", session.RootCause)
		}
	} else {
		fmt.Println("   ⚠️  Root cause NOT yet determined")
		fmt.Println("   Continue: vic debug pattern → vic debug hypothesis")
	}

	if len(findings) > 0 {
		fmt.Println("")
		fmt.Println("   📝 Findings:")
		for _, f := range findings {
			fmt.Printf("   • [%s] %s\n", f.Phase, f.Description)
		}
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// NewDebugHistoryCmd shows debug history
func NewDebugHistoryCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "history",
		Short:   "Show debug history",
		Long:    `Display past debug sessions.`,
		Example: `  vic debug history`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugHistory(cfg)
		},
	}
}

func runDebugHistory(cfg *config.Config) error {
	log, err := loadDebugLog(cfg)
	if err != nil || len(log.Sessions) == 0 {
		fmt.Println("❌ No debug history found.")
		return nil
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📜 Debug History")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	for _, s := range log.Sessions {
		statusIcon := "🔍"
		if s.RootCauseFound {
			statusIcon = "✅"
		}
		fmt.Printf("   %s %s | %s\n", statusIcon, s.StartedAt.Format("2006-01-02 15:04"), truncate(s.Problem, 50))
		fmt.Printf("       Phase: %s | Attempts: %d\n", s.CurrentPhase, s.Attempts)
		if s.RootCause != "" {
			fmt.Printf("       Root Cause: %s\n", truncate(s.RootCause, 60))
		}
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("   Total sessions: %d\n", len(log.Sessions))
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// ============================================
// Helper Functions
// ============================================

func loadDebugSession(cfg *config.Config) (DebugSession, error) {
	sessionFile := cfg.ProjectDir + "/status/debug-session.yaml"
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return DebugSession{}, err
	}

	var session DebugSession
	if err := yaml.Unmarshal(data, &session); err != nil {
		return DebugSession{}, err
	}
	return session, nil
}

func saveDebugSession(cfg *config.Config, session DebugSession) {
	sessionFile := cfg.ProjectDir + "/status/debug-session.yaml"
	data, _ := yaml.Marshal(session)
	os.WriteFile(sessionFile, data, 0644)
}

func loadDebugFindings(cfg *config.Config) []DebugFinding {
	logFile := filepath.Join(cfg.GetVICDir(), "status", "debug-log.yaml")
	data, err := os.ReadFile(logFile)
	if err != nil {
		return []DebugFinding{}
	}

	var log DebugLog
	if err := yaml.Unmarshal(data, &log); err != nil {
		return []DebugFinding{}
	}
	return log.Findings
}

func saveDebugFinding(cfg *config.Config, finding DebugFinding) {
	logFile := filepath.Join(cfg.GetVICDir(), "status", "debug-log.yaml")

	var log DebugLog
	data, _ := os.ReadFile(logFile)
	if len(data) > 0 {
		yaml.Unmarshal(data, &log)
	}

	log.Findings = append(log.Findings, finding)
	log.SessionID = fmt.Sprintf("DEBUG-%d", len(log.Sessions)+1)

	output, _ := yaml.Marshal(log)
	os.WriteFile(logFile, output, 0644)
}

func loadDebugLog(cfg *config.Config) (DebugLog, error) {
	logFile := filepath.Join(cfg.GetVICDir(), "status", "debug-log.yaml")
	data, err := os.ReadFile(logFile)
	if err != nil {
		return DebugLog{}, err
	}

	var log DebugLog
	if err := yaml.Unmarshal(data, &log); err != nil {
		return DebugLog{}, err
	}
	return log, nil
}
