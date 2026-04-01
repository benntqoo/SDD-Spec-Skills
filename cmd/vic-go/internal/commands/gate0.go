package commands

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// Gate 0 check items for SPEC-REQUIREMENTS.md
var gate0Checks = []struct {
	id      string
	name    string
	pattern string
	errMsg  string
}{
	{
		id:      "USER_STORIES",
		name:    "User Stories Section",
		pattern: `(?i)(##\s*User Stories|As a [a-zA-Z]+.*can|user story)`,
		errMsg:  "Missing User Stories section - define who users are and what they can do",
	},
	{
		id:      "KEY_FEATURES",
		name:    "Key Features Section",
		pattern: `(?i)(##\s*Key Features|##\s*Features|feature list)`,
		errMsg:  "Missing Key Features section - list main capabilities",
	},
	{
		id:      "ACCEPTANCE",
		name:    "Acceptance Criteria",
		pattern: `(?i)(acceptance criteria|given.*when.*then|验收标准)`,
		errMsg:  "Missing Acceptance Criteria - define what 'done' means",
	},
	{
		id:      "NON_FUNC",
		name:    "Non-Functional Requirements",
		pattern: `(?i)(non-functional|performance|security|scalability|响应时间|并发)`,
		errMsg:  "Missing Non-Functional Requirements section",
	},
	{
		id:      "OUT_OF_SCOPE",
		name:    "Out of Scope",
		pattern: `(?i)(out of scope|not included|不在范围内)`,
		errMsg:  "Missing Out of Scope section - clarify boundaries",
	},
}

type gate0Result struct {
	checkID    string
	checkName  string
	passed     bool
	message    string
	lineNumber int
}

// RunGate0 validates SPEC-REQUIREMENTS.md structure
func RunGate0(cfg *config.Config, format string) error {
	// Use GateReport for JSON output
	report := NewGateReport(0)

	// Plain output header
	if format != "json" {
		fmt.Println("🔍 Gate 0: Requirements Completeness Check")
		fmt.Println("========================================")
		fmt.Println()
	}

	// Check file exists
	if !fileExists(cfg.SpecRequirements) {
		fmt.Println("❌ SPEC-REQUIREMENTS.md not found")
		fmt.Println()
		fmt.Println("   Run 'vic spec init' to create it")
		fmt.Println("   Then edit SPEC-REQUIREMENTS.md with your requirements")
		return nil
	}

	fmt.Printf("📄 Analyzing: %s\n\n", cfg.SpecRequirements)

	// Read file content
	content, err := os.ReadFile(cfg.SpecRequirements)
	if err != nil {
		return fmt.Errorf("failed to read SPEC-REQUIREMENTS.md: %w", err)
	}
	contentStr := string(content)

	// Run checks
	results := make([]gate0Result, 0)
	allPassed := true

	for _, check := range gate0Checks {
		result := gate0Result{
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

	// Check for TODOs/TBDs
	todoResults := checkForTODOs(contentStr)
	results = append(results, todoResults...)

	// Check if features have acceptance criteria
	featureCheck := checkFeaturesHaveCriteria(contentStr)
	results = append(results, featureCheck)

	// Collect results into report
	for _, r := range results {
		report.AddCheck(r.checkID, r.checkName, r.passed, r.message)
	}

	// Finalize and output report
	todoCount := 0
	for _, r := range todoResults {
		if !r.passed {
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
		statusText := "FAIL"
		if r.passed {
			statusIcon = "✅"
			statusText = "PASS"
		}
		fmt.Printf("[%s] %s: %s\n", statusIcon, r.checkName, statusText)
		if !r.passed {
			fmt.Printf("      → %s\n", r.message)
		}
	}

	fmt.Println()
	fmt.Println("========================================")

	if allPassed && todoCount == 0 {
		fmt.Println("✅ Gate 0 PASSED - Requirements document is complete")
		fmt.Println()
		fmt.Println("Next: Run 'vic spec gate 1' to check Architecture")
		return nil
	}

	fmt.Println("❌ Gate 0 FAILED - Requirements incomplete")
	fmt.Println()
	if !allPassed {
		fmt.Println("Required sections missing - edit SPEC-REQUIREMENTS.md to add:")
		for _, r := range results {
			if !r.passed && r.checkID != "TODO" {
				fmt.Printf("   • %s\n", r.checkName)
			}
		}
	}
	if todoCount > 0 {
		fmt.Println()
		fmt.Println("⚠️  Found TBD/TODO/FIXME in document - resolve these before proceeding")
	}
	fmt.Println()
	fmt.Println("Edit SPEC-REQUIREMENTS.md to fix these issues, then run 'vic spec gate 0' again")

	return nil
}

// checkFeaturesHaveCriteria verifies features have acceptance criteria
func checkFeaturesHaveCriteria(content string) gate0Result {
	// Find feature list section
	featureRe := regexp.MustCompile(`(?mi)^[-*]\s*(.+?)(?:\s*\n|$)`)
	features := featureRe.FindAllStringSubmatch(content, -1)

	if len(features) == 0 {
		return gate0Result{
			checkID:   "FEATURES",
			checkName: "Features Have Criteria",
			passed:    true,
			message:   "No features to check",
		}
	}

	// Count features with criteria (look for "given/when/then" near feature)
	criteriaRe := regexp.MustCompile(`(?i)(given|when|then|acceptance|验收)`)
	featuresWithCriteria := 0

	for _, match := range features {
		if len(match) > 1 {
			feature := match[1]
			// Look in a window around the feature
			idx := strings.Index(content, feature)
			if idx >= 0 {
				windowStart := idx
				windowEnd := idx + 500
				if windowEnd > len(content) {
					windowEnd = len(content)
				}
				window := content[windowStart:windowEnd]
				if criteriaRe.MatchString(window) {
					featuresWithCriteria++
				}
			}
		}
	}

	passRatio := float64(featuresWithCriteria) / float64(len(features))
	if passRatio >= 0.5 {
		return gate0Result{
			checkID:   "FEATURES",
			checkName: "Features Have Criteria",
			passed:    true,
			message:   fmt.Sprintf("%d/%d features have acceptance criteria", featuresWithCriteria, len(features)),
		}
	}

	return gate0Result{
		checkID:   "FEATURES",
		checkName: "Features Have Criteria",
		passed:    false,
		message:   fmt.Sprintf("Only %d/%d features have acceptance criteria - add criteria for each feature", featuresWithCriteria, len(features)),
	}
}

// NewGate0Cmd creates the gate 0 command
func NewGate0Cmd(cfg *config.Config) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "gate0",
		Short: "Validate SPEC-REQUIREMENTS.md structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGate0(cfg, outputFormat)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}
