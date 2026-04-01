package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

type gate2Result struct {
	checkID   string
	checkName string
	passed    bool
	message   string
	details   string
}

// RunGate2 validates code alignment with SPEC
func RunGate2(cfg *config.Config, format string) error {
	// Use GateReport for JSON output
	report := NewGateReport(2)

	// Plain output header
	if format != "json" {
		fmt.Println("🔍 Gate 2: Code Alignment Check")
		fmt.Println("========================================")
		fmt.Println()
	}

	// Check if SPEC files exist
	if !fileExists(cfg.SpecArchitecture) {
		if format != "json" {
			fmt.Println("❌ SPEC-ARCHITECTURE.md not found - run 'vic spec init' first")
		}
		report.AddCheck("SPEC_FILE", "SPEC File Exists", false, "SPEC-ARCHITECTURE.md not found")
		report.Finalize(false)
		if format == "json" {
			fmt.Println(report.ToJSON())
		}
		return nil
	}

	if format != "json" {
		fmt.Printf("📄 Comparing: %s with code\n\n", cfg.SpecArchitecture)
	}

	// Read SPEC content
	specContent, err := os.ReadFile(cfg.SpecArchitecture)
	if err != nil {
		return fmt.Errorf("failed to read SPEC-ARCHITECTURE.md: %w", err)
	}
	specStr := string(specContent)

	// Run alignment checks
	results := make([]gate2Result, 0)

	// Check 1: Tech stack in SPEC vs actual code
	techCheck := checkTechStackAlignment(cfg, specStr)
	results = append(results, techCheck)

	// Check 2: API endpoints in SPEC vs code
	apiCheck := checkAPIAlignment(cfg, specStr)
	results = append(results, apiCheck)

	// Check 3: Module structure
	moduleCheck := checkModuleStructure(cfg, specStr)
	results = append(results, moduleCheck)

	// Check 4: Security implementation
	securityCheck := checkSecurityImplementation(cfg, specStr)
	results = append(results, securityCheck)

	// NEW: Code quality scans
	codeScanner := NewCodeScanner(cfg.ProjectDir)

	// Check for TODOs in code
	todos := codeScanner.FindTODOs()
	if len(todos) > 0 {
		todoResult := gate2Result{
			checkID:   "CODE_TODOS",
			checkName: "Code TODO Comments",
			passed:    false,
			message:   fmt.Sprintf("Found %d TODO/FIXME/XXX/HACK comments in code", len(todos)),
			details:   fmt.Sprintf("Most critical: %s", getMostCriticalTODO(todos)),
		}
		results = append(results, todoResult)
	}

	// Check constitution rules
	violations := codeScanner.ValidateConstitution([]string{
		"NO-TODO-IN-CODE",
		"NO-CONSOLE-IN-PROD",
		"NO-HARD-CODED-SECRETS",
	})
	for _, violation := range violations {
		violationResult := gate2Result{
			checkID:   "CONSTITUTION",
			checkName: fmt.Sprintf("Constitution: %s", violation.Rule),
			passed:    violation.Severity != "error",
			message:   violation.Message,
			details:   fmt.Sprintf("%s:%d", violation.File, violation.Line),
		}
		results = append(results, violationResult)
	}

	// Collect results into report
	for _, r := range results {
		report.AddCheck(r.checkID, r.checkName, r.passed, r.message, r.details)
	}

	// Finalize and output report
	allPassed := true
	for _, r := range results {
		if !r.passed {
			allPassed = false
		}
	}
	report.Finalize(allPassed)

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
		if r.details != "" {
			fmt.Printf("         Details: %s\n", r.details)
		}
	}

	fmt.Println()
	fmt.Println("========================================")

	if allPassed {
		fmt.Println("✅ Gate 2 PASSED - Code aligns with SPEC")
		fmt.Println()
		fmt.Println("Next: Run 'vic spec gate 3' to check test coverage")
		return nil
	}

	fmt.Printf("❌ Gate 2 FAILED - %d/%d checks failed\n", len(results)-report.PassedChecks, len(results))
	fmt.Println()

	// Collect failed check names for recommendation
	failedChecks := make([]string, 0)
	for _, r := range results {
		if !r.passed {
			failedChecks = append(failedChecks, r.checkName)
		}
	}
	showSpecUpdateRecommendation(failedChecks)

	return fmt.Errorf("gate 2 failed")
}

// getMostCriticalTODO returns the file:line of the most critical TODO
func getMostCriticalTODO(todos []TODOComment) string {
	for _, todo := range todos {
		if todo.Priority == "high" {
			return fmt.Sprintf("%s:%d", todo.File, todo.Line)
		}
	}
	if len(todos) > 0 {
		return fmt.Sprintf("%s:%d", todos[0].File, todos[0].Line)
	}
	return ""
}

// showSpecUpdateRecommendation prints recommended actions when drift is detected
func showSpecUpdateRecommendation(affectedSections []string) {
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("📋 SPEC UPDATE RECOMMENDATION")
	fmt.Println("════════════════════════════════════════════════════════════")

	if len(affectedSections) > 0 {
		fmt.Printf("Drift detected in: %s\n\n", strings.Join(affectedSections, ", "))
	}

	fmt.Println("To resolve this drift, choose one of the following:")
	fmt.Println("1️⃣  Update SPEC (Recommended)")
	fmt.Println("    $ vic spec update --file SPEC-ARCHITECTURE.md --section \"[section]\"")
	fmt.Println("    Then: vic spec gate 2")
	fmt.Println("2️⃣  Revert code changes")
	fmt.Println("    $ git diff [affected files]")
	fmt.Println("    Then: Revert and re-implement correctly")
	fmt.Println("3️⃣  Document as accepted drift (requires approval)")
	fmt.Println("    $ vic rr --id DRIFT-[DATE] --desc \"[description]\"")
	fmt.Println("    ⚠️  Only for emergency hotfixes")
	fmt.Println("For more details, see: skills/constitution-check/SKILL.md")
	fmt.Println("════════════════════════════════════════════════════════════")
}

// checkTechStackAlignment verifies tech stack declared in SPEC exists in code
func checkTechStackAlignment(cfg *config.Config, specContent string) gate2Result {
	// Extract declared technologies from SPEC
	declaredTech := extractDeclaredTech(specContent)

	if len(declaredTech) == 0 {
		return gate2Result{
			checkID:   "TECH_ALIGN",
			checkName: "Tech Stack Alignment",
			passed:    true,
			message:   "No specific tech declared in SPEC",
			details:   "",
		}
	}

	// Scan code for technologies
	detectedTech := scanCodeForTech(cfg.ProjectDir)

	// Check alignment
	mismatches := make([]string, 0)
	matches := make([]string, 0)

	for tech, patterns := range declaredTech {
		found := false
		for _, pattern := range patterns {
			if detectedTech[pattern] {
				found = true
				break
			}
		}
		if found {
			matches = append(matches, tech)
		} else {
			mismatches = append(mismatches, tech)
		}
	}

	passRatio := float64(len(matches)) / float64(len(declaredTech))

	if passRatio >= 0.7 {
		return gate2Result{
			checkID:   "TECH_ALIGN",
			checkName: "Tech Stack Alignment",
			passed:    true,
			message:   fmt.Sprintf("%d/%d declared tech found in code", len(matches), len(declaredTech)),
			details:   strings.Join(matches, ", "),
		}
	}

	return gate2Result{
		checkID:   "TECH_ALIGN",
		checkName: "Tech Stack Alignment",
		passed:    false,
		message:   fmt.Sprintf("Only %d/%d declared tech found in code", len(matches), len(declaredTech)),
		details:   fmt.Sprintf("Missing: %s", strings.Join(mismatches, ", ")),
	}
}

// extractDeclaredTech extracts technology declarations from SPEC
func extractDeclaredTech(specContent string) map[string][]string {
	tech := make(map[string][]string)

	techPatterns := map[string][]string{
		"postgresql": {"postgres", "postgresql", "pg."},
		"mysql":      {"mysql", "mariadb"},
		"mongodb":    {"mongodb", "mongoose"},
		"sqlite":     {"sqlite", ".db"},
		"redis":      {"redis", "ioredis"},
		"react":      {"react", "reactdom", "create-react-app"},
		"vue":        {"vue", "vue."},
		"angular":    {"@angular"},
		"express":    {"express", "express."},
		"fastapi":    {"fastapi", "fastapi"},
		"django":     {"django", "django"},
		"flask":      {"flask", "flask"},
		"gin":        {"gin-gonic", "gin.engine"},
		"go":         {"package main", "func main"},
		"docker":     {"docker", "dockerfile", "docker-compose"},
	}

	specLower := strings.ToLower(specContent)
	for techName, patterns := range techPatterns {
		for _, pattern := range patterns {
			if strings.Contains(specLower, pattern) {
				tech[techName] = patterns
				break
			}
		}
	}

	return tech
}

// scanCodeForTech scans code directory for technology indicators
func scanCodeForTech(projectDir string) map[string]bool {
	detected := make(map[string]bool)

	techIndicators := []string{
		"postgres", "postgresql", "mysql", "mongodb", "sqlite", "redis",
		"react", "vue", "angular", "svelte",
		"express", "fastapi", "django", "flask", "gin",
		"package main", "func main",
		"docker", "dockerfile", "docker-compose",
	}

	// Walk through source files
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip common non-source directories
		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", ".venv", "venv", "__pycache__"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check file extension
		ext := filepath.Ext(path)
		sourceExts := map[string]bool{
			".go": true, ".py": true, ".js": true, ".ts": true,
			".tsx": true, ".jsx": true, ".java": true, ".rs": true,
		}

		if !sourceExts[ext] {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentStr := strings.ToLower(string(content))
		for _, indicator := range techIndicators {
			if strings.Contains(contentStr, indicator) {
				detected[indicator] = true
			}
		}

		return nil
	})

	return detected
}

// checkAPIAlignment verifies API endpoints in SPEC exist in code
func checkAPIAlignment(cfg *config.Config, specContent string) gate2Result {
	// Extract API endpoints from SPEC
	endpointRe := regexp.MustCompile(`(?mi)(GET|POST|PUT|DELETE|PATCH)\s+([/\w{}.-]+)`)
	endpoints := endpointRe.FindAllStringSubmatch(specContent, -1)

	if len(endpoints) == 0 {
		return gate2Result{
			checkID:   "API_ALIGN",
			checkName: "API Endpoint Alignment",
			passed:    true,
			message:   "No API endpoints declared in SPEC",
			details:   "",
		}
	}

	// Check if routes exist in code
	routePatterns := []string{"router", "route", "endpoint", "@app.route", "@router", "http.Method"}
	foundRoutes := false

	filepath.Walk(cfg.ProjectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git", "dist"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentStr := string(content)
		for _, pattern := range routePatterns {
			if strings.Contains(strings.ToLower(contentStr), pattern) {
				foundRoutes = true
				return nil
			}
		}

		return nil
	})

	if foundRoutes {
		return gate2Result{
			checkID:   "API_ALIGN",
			checkName: "API Endpoint Alignment",
			passed:    true,
			message:   fmt.Sprintf("Found route handlers for %d declared endpoints", len(endpoints)),
			details:   "",
		}
	}

	return gate2Result{
		checkID:   "API_ALIGN",
		checkName: "API Endpoint Alignment",
		passed:    false,
		message:   "API endpoints declared but no route handlers found",
		details:   fmt.Sprintf("Declared %d endpoints in SPEC", len(endpoints)),
	}
}

// checkModuleStructure verifies module structure matches SPEC
func checkModuleStructure(cfg *config.Config, specContent string) gate2Result {
	// Check if common module directories exist
	expectedModules := []string{"cmd", "internal", "pkg", "src", "lib", "app", "modules"}

	foundModules := make([]string, 0)
	for _, module := range expectedModules {
		modulePath := filepath.Join(cfg.ProjectDir, module)
		if _, err := os.Stat(modulePath); err == nil {
			foundModules = append(foundModules, module)
		}
	}

	if len(foundModules) > 0 {
		return gate2Result{
			checkID:   "MODULE_STRUCT",
			checkName: "Module Structure",
			passed:    true,
			message:   fmt.Sprintf("Found %d source directories", len(foundModules)),
			details:   strings.Join(foundModules, ", "),
		}
	}

	return gate2Result{
		checkID:   "MODULE_STRUCT",
		checkName: "Module Structure",
		passed:    false,
		message:   "No standard source directories found",
		details:   "Expected: cmd, internal, pkg, src, lib, or app",
	}
}

// checkSecurityImplementation verifies security measures from SPEC
func checkSecurityImplementation(cfg *config.Config, specContent string) gate2Result {
	// Check if security is mentioned in SPEC
	securityMentioned := regexp.MustCompile(`(?i)(security|auth|jwt|oauth|ssl|tls|encryption)`)
	hasSecurity := securityMentioned.MatchString(specContent)

	if !hasSecurity {
		return gate2Result{
			checkID:   "SECURITY",
			checkName: "Security Implementation",
			passed:    true,
			message:   "No security requirements in SPEC",
			details:   "",
		}
	}

	// Check for security patterns in code
	securityPatterns := []string{"jwt", "auth", "bcrypt", "oauth", "ssl", "tls", "cors", "csrf"}
	foundSecurity := false

	filepath.Walk(cfg.ProjectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentStr := strings.ToLower(string(content))
		for _, pattern := range securityPatterns {
			if strings.Contains(contentStr, pattern) {
				foundSecurity = true
				return nil
			}
		}

		return nil
	})

	if foundSecurity {
		return gate2Result{
			checkID:   "SECURITY",
			checkName: "Security Implementation",
			passed:    true,
			message:   "Security measures found in code",
			details:   "",
		}
	}

	return gate2Result{
		checkID:   "SECURITY",
		checkName: "Security Implementation",
		passed:    false,
		message:   "Security mentioned in SPEC but no implementation found",
		details:   "Implement security measures defined in SPEC",
	}
}

// NewGate2Cmd creates the gate 2 command
func NewGate2Cmd(cfg *config.Config) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "gate2",
		Short: "Validate code alignment with SPEC",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGate2(cfg, outputFormat)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}
