package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

type gate3Result struct {
	checkID   string
	checkName string
	passed    bool
	message   string
	details   string
}

// RunGate3 validates test coverage
func RunGate3(cfg *config.Config, format string) error {
	// Use GateReport for JSON output
	report := NewGateReport(3)

	// Plain output header
	if format != "json" {
		fmt.Println("🔍 Gate 3: Test Coverage Check")
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
		fmt.Printf("📄 Checking test coverage in: %s\n\n", cfg.ProjectDir)
	}

	// Run test checks
	results := make([]gate3Result, 0)

	// Check 1: Test files exist
	testFilesCheck := checkTestFilesExist(cfg.ProjectDir)
	results = append(results, testFilesCheck)

	// Check 2: Test framework used
	frameworkCheck := checkTestFramework(cfg.ProjectDir)
	results = append(results, frameworkCheck)

	// Check 3: Key modules have tests
	moduleCoverageCheck := checkModuleTestCoverage(cfg.ProjectDir)
	results = append(results, moduleCoverageCheck)

	// Check 4: Critical paths covered
	criticalPathCheck := checkCriticalPathCoverage(cfg.ProjectDir)
	results = append(results, criticalPathCheck)

	// Collect results into report
	for _, r := range results {
		report.AddCheck(r.checkID, r.checkName, r.passed, r.message, r.details)
	}

	// Finalize and output report
	passedCount := 0
	for _, r := range results {
		if r.passed {
			passedCount++
		}
	}
	overallSuccess := passedCount >= 3 // At least 3 of 4 checks should pass
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
		if r.details != "" {
			fmt.Printf("         Details: %s\n", r.details)
		}
	}

	fmt.Println()
	fmt.Println("========================================")

	if passedCount >= 3 { // At least 3 of 4 checks should pass
		fmt.Println("✅ Gate 3 PASSED - Adequate test coverage")
		fmt.Println()
		fmt.Println("🎉 All Gates passed! Your implementation is complete.")
		fmt.Println()
		fmt.Println("Run 'vic spec merge' to finalize documentation")
		return nil
	}

	fmt.Printf("❌ Gate 3 FAILED - %d/%d checks passed\n", passedCount, len(results))
	fmt.Println()
	fmt.Println("Improve test coverage, then run 'vic spec gate 3' again")

	return nil
}

// checkTestFilesExist verifies test files exist
func checkTestFilesExist(projectDir string) gate3Result {
	testPatterns := []string{
		"*_test.go",
		"*_test.py",
		"*.test.js",
		"*.test.ts",
		"*.spec.js",
		"*.spec.ts",
		"*.test.jsx",
		"*.test.tsx",
		"test_*.py",
		"tests/*.py",
		"*_test.rs",
		"*_test.go",
	}

	testFiles := make([]string, 0)

	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", ".venv", "venv"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		for _, pattern := range testPatterns {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				testFiles = append(testFiles, path)
				return nil
			}
		}

		return nil
	})

	if len(testFiles) > 0 {
		return gate3Result{
			checkID:   "TEST_FILES",
			checkName: "Test Files Exist",
			passed:    true,
			message:   fmt.Sprintf("Found %d test file(s)", len(testFiles)),
			details:   strings.Join(testFiles[:min(5, len(testFiles))], ", "),
		}
	}

	return gate3Result{
		checkID:   "TEST_FILES",
		checkName: "Test Files Exist",
		passed:    false,
		message:   "No test files found",
		details:   "Expected: *_test.go, *_test.py, *.test.js, etc.",
	}
}

// checkTestFramework verifies a test framework is configured
func checkTestFramework(projectDir string) gate3Result {
	// Common test framework indicators
	frameworkPatterns := map[string][]string{
		"Go":         {"testing", "_test.go", "go test", "testify"},
		"Python":     {"pytest", "unittest", "test_", "tests/"},
		"JavaScript": {"jest", "mocha", "vitest", "testing-library"},
		"TypeScript": {"jest", "mocha", "vitest", "@testing-library"},
		"Rust":       {"#\\[test\\]", "#[tokio::test]"},
		"Java":       {"@Test", "junit", "assertj"},
	}

	foundFrameworks := make([]string, 0)

	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build"}
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

		for framework, indicators := range frameworkPatterns {
			for _, indicator := range indicators {
				if strings.Contains(contentStr, strings.ToLower(indicator)) {
					found := false
					for _, f := range foundFrameworks {
						if f == framework {
							found = true
							break
						}
					}
					if !found {
						foundFrameworks = append(foundFrameworks, framework)
					}
					break
				}
			}
		}

		return nil
	})

	// Also check config files
	configFiles := []string{
		"package.json", "go.mod", "Cargo.toml", "requirements.txt",
		"pytest.ini", "jest.config.js", "jest.config.ts", "vitest.config.ts",
		"pyproject.toml", "setup.py",
	}

	for _, configFile := range configFiles {
		configPath := filepath.Join(projectDir, configFile)
		if content, err := os.ReadFile(configPath); err == nil {
			contentStr := strings.ToLower(string(content))
			for framework, indicators := range frameworkPatterns {
				for _, indicator := range indicators {
					if strings.Contains(contentStr, strings.ToLower(indicator)) {
						found := false
						for _, f := range foundFrameworks {
							if f == framework {
								found = true
								break
							}
						}
						if !found {
							foundFrameworks = append(foundFrameworks, framework)
						}
						break
					}
				}
			}
		}
	}

	if len(foundFrameworks) > 0 {
		return gate3Result{
			checkID:   "TEST_FRAMEWORK",
			checkName: "Test Framework Configured",
			passed:    true,
			message:   fmt.Sprintf("Found: %s", strings.Join(foundFrameworks, ", ")),
			details:   "",
		}
	}

	return gate3Result{
		checkID:   "TEST_FRAMEWORK",
		checkName: "Test Framework Configured",
		passed:    false,
		message:   "No test framework detected",
		details:   "Configure a test framework (jest, pytest, go test, etc.)",
	}
}

// checkModuleTestCoverage verifies key modules have tests
func checkModuleTestCoverage(projectDir string) gate3Result {
	// Find source directories
	sourceDirs := []string{"cmd", "internal", "pkg", "src", "lib", "app", "modules"}
	testDirs := []string{"test", "tests", "__tests__"}

	hasSource := false
	hasTests := false

	for _, dir := range sourceDirs {
		dirPath := filepath.Join(projectDir, dir)
		if _, err := os.Stat(dirPath); err == nil {
			hasSource = true
			break
		}
	}

	for _, dir := range testDirs {
		dirPath := filepath.Join(projectDir, dir)
		if _, err := os.Stat(dirPath); err == nil {
			hasTests = true
			break
		}
	}

	if !hasSource {
		return gate3Result{
			checkID:   "MODULE_COVERAGE",
			checkName: "Module Test Coverage",
			passed:    true,
			message:   "No source modules to test",
			details:   "",
		}
	}

	if hasTests {
		return gate3Result{
			checkID:   "MODULE_COVERAGE",
			checkName: "Module Test Coverage",
			passed:    true,
			message:   "Test directory found alongside source",
			details:   "",
		}
	}

	// Check if test files are co-located with source
	testFileCount := 0
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
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

		baseName := filepath.Base(path)
		if strings.HasSuffix(baseName, "_test.go") ||
			strings.HasSuffix(baseName, "_test.py") ||
			strings.HasSuffix(baseName, ".test.js") ||
			strings.HasSuffix(baseName, ".test.ts") {
			testFileCount++
		}

		return nil
	})

	if testFileCount > 0 {
		return gate3Result{
			checkID:   "MODULE_COVERAGE",
			checkName: "Module Test Coverage",
			passed:    true,
			message:   fmt.Sprintf("Found %d co-located test file(s)", testFileCount),
			details:   "",
		}
	}

	return gate3Result{
		checkID:   "MODULE_COVERAGE",
		checkName: "Module Test Coverage",
		passed:    false,
		message:   "Test files not found near source modules",
		details:   "Add test files (*_test.go, *_test.py, *.test.js)",
	}
}

// checkCriticalPathCoverage verifies critical paths are tested
func checkCriticalPathCoverage(projectDir string) gate3Result {
	// Common critical path patterns
	criticalPatterns := []string{
		"main.go", "main.py", "index.js", "app.js",
		"handler", "controller", "service", "business",
		"auth", "login", "signup", "user",
		"api", "route", "endpoint",
	}

	criticalFiles := make([]string, 0)
	testedFiles := make(map[string]bool)

	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", ".venv", "venv"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		baseName := strings.ToLower(filepath.Base(path))

		// Check if it's a critical file
		for _, pattern := range criticalPatterns {
			if strings.Contains(baseName, pattern) {
				criticalFiles = append(criticalFiles, path)

				// Check if there's a corresponding test file
				dir := filepath.Dir(path)
				name := filepath.Base(path)
				ext := filepath.Ext(name)
				base := name[:len(name)-len(ext)]

				testVariants := []string{
					filepath.Join(dir, base+"_test"+ext),
					filepath.Join(dir, base+".test"+ext),
					filepath.Join(dir, "test_"+name),
					filepath.Join(dir, "tests", base+"_test"+ext),
				}

				for _, testPath := range testVariants {
					if _, err := os.Stat(testPath); err == nil {
						testedFiles[path] = true
						break
					}
				}
				break
			}
		}

		return nil
	})

	if len(criticalFiles) == 0 {
		return gate3Result{
			checkID:   "CRITICAL_PATH",
			checkName: "Critical Path Coverage",
			passed:    true,
			message:   "No obvious critical paths detected",
			details:   "",
		}
	}

	coverageRatio := float64(len(testedFiles)) / float64(len(criticalFiles))

	if coverageRatio >= 0.3 {
		return gate3Result{
			checkID:   "CRITICAL_PATH",
			checkName: "Critical Path Coverage",
			passed:    true,
			message:   fmt.Sprintf("%d/%d critical files have tests", len(testedFiles), len(criticalFiles)),
			details:   "",
		}
	}

	return gate3Result{
		checkID:   "CRITICAL_PATH",
		checkName: "Critical Path Coverage",
		passed:    false,
		message:   fmt.Sprintf("Only %d/%d critical files have tests", len(testedFiles), len(criticalFiles)),
		details:   "Add tests for main files, handlers, and auth modules",
	}
}

// NewGate3Cmd creates the gate 3 command
func NewGate3Cmd(cfg *config.Config) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "gate3",
		Short: "Validate test coverage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGate3(cfg, outputFormat)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
