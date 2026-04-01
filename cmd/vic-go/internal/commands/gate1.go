package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// Gate 1 check items for SPEC-ARCHITECTURE.md
var gate1Checks = []struct {
	id          string
	name        string
	pattern     string
	errMsg      string
	suggestions []string
}{
	{
		id:          "TECH_STACK",
		name:        "Technology Stack Section",
		pattern:     `(?i)(##\s*Technology Stack|技术栈|技术选型)`,
		errMsg:      "Missing Technology Stack section",
		suggestions: []string{"Add a table listing: Frontend, Backend, Database, Infrastructure"},
	},
	{
		id:          "SYSTEM_DESIGN",
		name:        "System Design Section",
		pattern:     `(?i)(##\s*System Design|##\s*Architecture|系统设计|系统架构)`,
		errMsg:      "Missing System Design section",
		suggestions: []string{"Add architecture diagram and component descriptions"},
	},
	{
		id:          "API_DESIGN",
		name:        "API Design Section",
		pattern:     `(?i)(##\s*API|##\s*Endpoints|API设计|接口设计)`,
		errMsg:      "Missing API Design section",
		suggestions: []string{"Define main API endpoints with methods and descriptions"},
	},
	{
		id:          "DATA_MODEL",
		name:        "Data Model Section",
		pattern:     `(?i)(##\s*Data Model|##\s*Schema|##\s*Database|data model|数据模型)`,
		errMsg:      "Missing Data Model section",
		suggestions: []string{"Define main entities and their relationships"},
	},
	{
		id:          "SECURITY",
		name:        "Security Section",
		pattern:     `(?i)(##\s*Security|##\s*Auth|安全|认证|授权)`,
		errMsg:      "Missing Security section",
		suggestions: []string{"Document authentication, authorization, and data protection"},
	},
	{
		id:          "DECISION_RATIONALE",
		name:        "Decision Rationale",
		pattern:     `(?i)(选择理由|rationale|为什么选择|decision reason)`,
		errMsg:      "Missing rationale for tech decisions",
		suggestions: []string{"Explain WHY each technology was chosen, not just WHAT"},
	},
}

type gate1Result struct {
	checkID   string
	checkName string
	passed    bool
	message   string
	details   string
}

// RunGate1 validates SPEC-ARCHITECTURE.md structure
func RunGate1(cfg *config.Config, format string) error {
	// Use GateReport for JSON output
	report := NewGateReport(1)

	// Plain output header
	if format != "json" {
		fmt.Println("🔍 Gate 1: Architecture Completeness Check")
		fmt.Println("========================================")
		fmt.Println()
	}

	// Check file exists
	if !fileExists(cfg.SpecArchitecture) {
		fmt.Println("❌ SPEC-ARCHITECTURE.md not found")
		fmt.Println()
		fmt.Println("   Run 'vic spec init' to create it")
		fmt.Println("   Then edit SPEC-ARCHITECTURE.md with your architecture")
		return nil
	}

	fmt.Printf("📄 Analyzing: %s\n\n", cfg.SpecArchitecture)

	// Read file content
	content, err := os.ReadFile(cfg.SpecArchitecture)
	if err != nil {
		return fmt.Errorf("failed to read SPEC-ARCHITECTURE.md: %w", err)
	}
	contentStr := string(content)

	// Run section checks
	results := make([]gate1Result, 0)
	allPassed := true

	for _, check := range gate1Checks {
		result := gate1Result{
			checkID:   check.id,
			checkName: check.name,
			passed:    false,
			message:   check.errMsg,
		}

		re := regexp.MustCompile(check.pattern)
		if re.MatchString(contentStr) {
			result.passed = true
			result.message = "✅ Found"
		} else {
			allPassed = false
		}

		results = append(results, result)
	}

	// Check tech stack completeness
	techCheck := checkTechStackCompleteness(contentStr)
	results = append(results, techCheck)

	// Check for TODOs
	todoResults := checkForTODOs(contentStr)
	for _, r := range todoResults {
		if r.checkID == "TODO" {
			results = append(results, gate1Result{
				checkID:   r.checkID,
				checkName: r.checkName,
				passed:    r.passed,
				message:   r.message,
			})
		}
	}

	// Collect results into report
	for _, r := range results {
		report.AddCheck(r.checkID, r.checkName, r.passed, r.message)
	}

	// Finalize and output report
	todoCount := 0
	for _, r := range results {
		if r.checkID == "TODO" && !r.passed {
			todoCount++
		}
	}
	overallSuccess := allPassed && todoCount == 0
	report.Finalize(overallSuccess)

	if format == "json" {
		fmt.Println(report.ToJSON())
		return nil
	}

	// Plain output
	for _, r := range results {
		statusIcon := "❌"
		if r.passed {
			statusIcon = "✅"
		}
		fmt.Printf("[%s] %s\n", statusIcon, r.checkName)
		if r.passed {
			fmt.Printf("      %s\n", r.message)
		} else {
			fmt.Printf("      → %s\n", r.message)
		}
	}

	fmt.Println()
	fmt.Println("========================================")

	if allPassed && todoCount == 0 {
		fmt.Println("✅ Gate 1 PASSED - Architecture document is complete")
		fmt.Println()
		fmt.Println("Next: Implement features, then run 'vic spec gate 2' to check alignment")
		return nil
	}

	fmt.Println("❌ Gate 1 FAILED - Architecture incomplete")
	fmt.Println()

	// Show suggestions
	fmt.Println("Required sections missing - edit SPEC-ARCHITECTURE.md:")
	for _, r := range results {
		if !r.passed && r.checkID != "TODO" {
			fmt.Printf("\n   📝 %s\n", r.checkName)
			for _, check := range gate1Checks {
				if check.id == r.checkID {
					for _, s := range check.suggestions {
						fmt.Printf("      💡 %s\n", s)
					}
					break
				}
			}
		}
	}

	if todoCount > 0 {
		fmt.Println()
		fmt.Println("⚠️  Found TBD/TODO/FIXME in document - resolve these before proceeding")
	}

	fmt.Println()
	fmt.Println("Run 'vic spec gate 1' again after fixing issues")

	return nil
}

// checkTechStackCompleteness verifies tech stack has necessary layers
func checkTechStackCompleteness(content string) gate1Result {
	// Required tech layers
	requiredLayers := map[string]string{
		"frontend": `(?i)(frontend|前端|react|vue|angular|svelte)`,
		"backend":  `(?i)(backend|后端|server|api|express|fastapi|gin|nestjs)`,
		"database": `(?i)(database|db|数据库|postgres|mysql|mongodb|redis|sqlite)`,
	}

	foundLayers := 0
	layerNames := make([]string, 0)

	for layerName, pattern := range requiredLayers {
		re := regexp.MustCompile(pattern)
		if re.MatchString(content) {
			foundLayers++
			layerNames = append(layerNames, layerName)
		}
	}

	if foundLayers >= 2 {
		return gate1Result{
			checkID:   "TECH_LAYERS",
			checkName: "Technology Layers Defined",
			passed:    true,
			message:   fmt.Sprintf("Found %d/3 layers: %v", foundLayers, layerNames),
		}
	}

	return gate1Result{
		checkID:   "TECH_LAYERS",
		checkName: "Technology Layers Defined",
		passed:    false,
		message:   fmt.Sprintf("Only %d/3 layers defined - add more technology layers", foundLayers),
	}
}

// NewGate1Cmd creates the gate 1 command
func NewGate1Cmd(cfg *config.Config) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "gate1",
		Short: "Validate SPEC-ARCHITECTURE.md structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGate1(cfg, outputFormat)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}
