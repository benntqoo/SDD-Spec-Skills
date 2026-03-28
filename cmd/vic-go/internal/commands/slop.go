package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
)

// NewSlopCmd creates the slop command
func NewSlopCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slop",
		Short: "AI Slop detection",
		Long:  `Detect AI-generated patterns in code and design.`,
		Example: `  vic slop scan           # Scan for AI slop patterns
  vic slop report         # Show last report
  vic slop list          # List configured patterns`,
	}

	cmd.AddCommand(NewSlopScanCmd(cfg))
	cmd.AddCommand(NewSlopReportCmd(cfg))
	cmd.AddCommand(NewSlopListCmd(cfg))
	cmd.AddCommand(NewSlopFixCmd(cfg))

	return cmd
}

// ============================================
// Scan Command
// ============================================

// NewSlopScanCmd creates the slop scan subcommand
func NewSlopScanCmd(cfg *config.Config) *cobra.Command {
	var patternType string

	cmd := &cobra.Command{
		Use:   "scan [directory]",
		Short: "Scan for AI slop patterns",
		Long:  `Scan source code and design files for AI-generated patterns.`,
		Example: `  vic slop scan
  vic slop scan ./src/components
  vic slop scan --type design`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runSlopScan(cfg, dir, patternType)
		},
	}

	cmd.Flags().StringVar(&patternType, "type", "", "Filter by pattern type (design/code/text)")

	return cmd
}

func runSlopScan(cfg *config.Config, dir, patternType string) error {
	fmt.Println("🔍 Scanning for AI Slop Patterns...")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	patterns := types.DefaultAISlopPatterns()

	// Filter patterns if type specified
	if patternType != "" {
		var filtered []types.AISlopPattern
		for _, p := range patterns {
			if p.Type == patternType {
				filtered = append(filtered, p)
			}
		}
		patterns = filtered
		fmt.Printf("   Scanning for %s patterns only\n\n", patternType)
	}

	var findings []types.AISlopFinding
	highCount, mediumCount, lowCount := 0, 0, 0

	// File extensions to scan
	extensions := []string{".tsx", ".jsx", ".ts", ".js", ".css", ".scss", ".md", ".html"}

	// Walk directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Skip common non-source directories
		if info.IsDir() {
			dirName := filepath.Base(path)
			if dirName == "node_modules" || dirName == "dist" || dirName == "build" ||
				dirName == ".git" || dirName == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check extension
		ext := strings.ToLower(filepath.Ext(path))
		shouldScan := false
		for _, e := range extensions {
			if ext == e {
				shouldScan = true
				break
			}
		}
		if !shouldScan {
			return nil
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Check each pattern
		for _, pattern := range patterns {
			if patternType != "" && pattern.Type != patternType {
				continue
			}

			re, err := regexp.Compile("(?i)" + pattern.Pattern)
			if err != nil {
				continue
			}

			matches := re.FindAllString(string(content), -1)
			if len(matches) > 0 {
				// Find line numbers
				lines := strings.Split(string(content), "\n")
				for lineNum, line := range lines {
					if re.MatchString(line) {
						finding := types.AISlopFinding{
							ID:        fmt.Sprintf("SLOP-%03d", len(findings)+1),
							PatternID: pattern.ID,
							Type:      pattern.Type,
							File:      path,
							Line:      lineNum + 1,
							Severity:  pattern.Severity,
							Match:     strings.TrimSpace(line),
							Fix:       pattern.Alternative,
						}
						findings = append(findings, finding)

						switch pattern.Severity {
						case "high":
							highCount++
						case "medium":
							mediumCount++
						case "low":
							lowCount++
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("⚠️  Scan error: %v\n", err)
	}

	// Calculate score
	total := len(findings)
	report := types.AISlopReport{
		TotalPatterns:  total,
		HighSeverity:   highCount,
		MediumSeverity: mediumCount,
		LowSeverity:    lowCount,
		Findings:       findings,
	}
	report.Score = report.CalculateScore()

	// Print results
	fmt.Printf("   Total findings: %d\n", total)
	fmt.Printf("   HIGH: %d | MEDIUM: %d | LOW: %d\n", highCount, mediumCount, lowCount)
	fmt.Println("")

	// Score
	scoreIcon := "✅"
	switch report.Score {
	case "A":
		scoreIcon = "✅"
	case "B":
		scoreIcon = "👍"
	case "C":
		scoreIcon = "⚠️"
	case "D":
		scoreIcon = "❌"
	}

	fmt.Printf("   AI Slop Score: %s %s\n", scoreIcon, report.Score)
	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// Print findings by severity
	if len(findings) > 0 {
		fmt.Println("")
		fmt.Println("📋 Findings:")

		// High first
		for _, f := range findings {
			if f.Severity == "high" {
				printAISlopFinding(f)
			}
		}
		for _, f := range findings {
			if f.Severity == "medium" {
				printAISlopFinding(f)
			}
		}
		for _, f := range findings {
			if f.Severity == "low" {
				printAISlopFinding(f)
			}
		}
	} else {
		fmt.Println("")
		fmt.Println("✅ No AI Slop patterns detected!")
	}

	// Summary by type
	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("📊 Summary by Type:")

	typeCounts := make(map[string]int)
	for _, f := range findings {
		typeCounts[f.Type]++
	}
	for t, c := range typeCounts {
		fmt.Printf("   %s: %d\n", strings.ToUpper(t), c)
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// Save report
	if err := saveAISlopReport(cfg, &report); err != nil {
		fmt.Printf("⚠️  Could not save report: %v\n", err)
	}

	return nil
}

func printAISlopFinding(f types.AISlopFinding) {
	severityIcon := "⚠️"
	switch f.Severity {
	case "high":
		severityIcon = "🔴"
	case "medium":
		severityIcon = "🟡"
	case "low":
		severityIcon = "🔵"
	}

	truncatedMatch := f.Match
	if len(truncatedMatch) > 60 {
		truncatedMatch = truncatedMatch[:57] + "..."
	}
	truncatedFix := f.Fix
	if len(truncatedFix) > 50 {
		truncatedFix = truncatedFix[:47] + "..."
	}

	fmt.Printf("\n%s [%s] %s - %s:%d\n", severityIcon, f.PatternID, strings.ToUpper(f.Type), f.File, f.Line)
	fmt.Printf("   Match: %s\n", truncatedMatch)
	fmt.Printf("   Fix: %s\n", truncatedFix)
}

func slopTruncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ============================================
// Report Command
// ============================================

// NewSlopReportCmd creates the slop report subcommand
func NewSlopReportCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "report",
		Short:   "Show last scan report",
		Long:    `Display the last AI Slop scan report.`,
		Example: `  vic slop report`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSlopReport(cfg)
		},
	}
}

func runSlopReport(cfg *config.Config) error {
	fmt.Println("📋 Use 'vic slop scan' to run a new scan")
	return nil
}

// ============================================
// List Command
// ============================================

// NewSlopListCmd creates the slop list subcommand
func NewSlopListCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List AI Slop patterns",
		Long:  `Show all configured AI Slop detection patterns.`,
		Example: `  vic slop list
  vic slop list --type design`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSlopList()
		},
	}
}

func runSlopList() error {
	patterns := types.DefaultAISlopPatterns()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  🚫 AI Slop Patterns")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	byType := make(map[string][]types.AISlopPattern)
	for _, p := range patterns {
		byType[p.Type] = append(byType[p.Type], p)
	}

	for typeName, typePatterns := range byType {
		fmt.Printf("  📁 %s\n", strings.ToUpper(typeName))
		fmt.Println("  ────────────────────────────────────────────────────")
		for _, p := range typePatterns {
			severityIcon := "⚠️"
			if p.Severity == "high" {
				severityIcon = "🔴"
			} else if p.Severity == "low" {
				severityIcon = "🔵"
			}
			fmt.Printf("    %s [%s] %s\n", severityIcon, p.ID, slopTruncate(p.Pattern, 40))
			fmt.Printf("       Fix: %s\n", slopTruncate(p.Alternative, 50))
		}
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  Total patterns: %d\n", len(patterns))

	return nil
}

// ============================================
// Helper Functions
// ============================================

func saveAISlopReport(cfg *config.Config, report *types.AISlopReport) error {
	slopFile := filepath.Join(cfg.GetVICDir(), "status", "slop-report.yaml")

	// Ensure directory exists
	if !utils.FileExists(cfg.ProjectDir + "/status") {
		if err := os.MkdirAll(cfg.ProjectDir+"/status", 0755); err != nil {
			return err
		}
	}

	// Simple report format
	content := fmt.Sprintf("# AI Slop Report\n")
	content += fmt.Sprintf("Score: %s | Total: %d | HIGH: %d | MED: %d | LOW: %d\n",
		report.Score, report.TotalPatterns, report.HighSeverity, report.MediumSeverity, report.LowSeverity)

	return os.WriteFile(slopFile, []byte(content), 0644)
}

// ============================================
// Fix Command
// ============================================

// NewSlopFixCmd creates the slop fix subcommand
func NewSlopFixCmd(cfg *config.Config) *cobra.Command {
	var dryRun bool
	var patternType string

	cmd := &cobra.Command{
		Use:   "fix [directory]",
		Short: "Auto-fix AI slop patterns",
		Long:  `Scan and auto-fix detected AI slop patterns. Use --dry-run to preview changes.`,
		Example: `  vic slop fix              # Scan and fix in current directory
  vic slop fix --dry-run     # Preview changes without applying
  vic slop fix ./src          # Fix in specific directory
  vic slop fix --type code    # Fix code patterns only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runSlopFix(cfg, dir, dryRun, patternType)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Preview changes without applying (default: true)")
	cmd.Flags().StringVar(&patternType, "type", "", "Filter by pattern type (design/code/text)")

	return cmd
}

func runSlopFix(cfg *config.Config, dir string, dryRun bool, filterType string) error {
	mode := "📝 PREVIEW"
	if !dryRun {
		mode = "🔧 FIXING"
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  %s AI Slop Patterns\n", mode)
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Directory: %s\n", dir)
	fmt.Printf("   Mode: %s\n", map[bool]string{true: "DRY RUN (no changes)", false: "LIVE (changes applied)"}[dryRun])
	if filterType != "" {
		fmt.Printf("   Filter: %s patterns only\n", filterType)
	}
	fmt.Println("")

	// Scan first
	patterns := types.DefaultAISlopPatterns()
	var findings []types.AISlopFinding
	fixableCount := 0
	warnOnlyCount := 0

	// File extensions to scan
	extensions := []string{".tsx", ".jsx", ".ts", ".js", ".css", ".scss", ".html"}

	// Walk directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Skip common non-source directories
		if info.IsDir() {
			dirName := filepath.Base(path)
			if dirName == "node_modules" || dirName == "dist" || dirName == "build" ||
				dirName == ".git" || dirName == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check extension
		ext := strings.ToLower(filepath.Ext(path))
		shouldScan := false
		for _, e := range extensions {
			if ext == e {
				shouldScan = true
				break
			}
		}
		if !shouldScan {
			return nil
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		modified := false

		// Check each pattern
		for _, pattern := range patterns {
			if filterType != "" && pattern.Type != filterType {
				continue
			}

			// Skip warn-only patterns
			if pattern.Severity == "low" {
				warnOnlyCount++
				continue
			}

			re, err := regexp.Compile("(?i)" + pattern.Pattern)
			if err != nil {
				continue
			}

			// Find matches and process
			for lineNum, line := range lines {
				if re.MatchString(line) {
					finding := types.AISlopFinding{
						ID:        fmt.Sprintf("FIX-%03d", len(findings)+1),
						PatternID: pattern.ID,
						Type:      pattern.Type,
						File:      path,
						Line:      lineNum + 1,
						Severity:  pattern.Severity,
						Match:     strings.TrimSpace(line),
						Fix:       pattern.Alternative,
					}
					findings = append(findings, finding)

					// Auto-fixable patterns
					fixable := false
					switch pattern.ID {
					case "C002": // console.log
						fixable = true
						if !dryRun {
							lines[lineNum] = "// [FIXED] " + lines[lineNum]
							modified = true
						}
					case "T001": // Hello World text
						fixable = true
						if !dryRun {
							lines[lineNum] = "// [FIXED] Placeholder text - replace with actual content"
							modified = true
						}
					case "C005": // setTimeout
						fixable = true
						if !dryRun {
							lines[lineNum] = "// [WARN] " + lines[lineNum] + " // Review for proper async handling"
							modified = true
						}
					}

					if fixable {
						fixableCount++
					}
				}
			}
		}

		// Write modified content
		if modified && !dryRun {
			if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644); err != nil {
				fmt.Printf("   ⚠️  Failed to write: %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("⚠️  Scan error: %v\n", err)
	}

	// Print findings summary
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Printf("   Total issues found: %d\n", len(findings))
	fmt.Printf("   Auto-fixable: %d\n", fixableCount)
	fmt.Printf("   Warn-only: %d\n", warnOnlyCount)
	fmt.Println("")

	// Print findings
	if len(findings) > 0 {
		fmt.Println("📋 Issues:")
		fmt.Println("───────────────────────────────────────────────────────────")

		for _, f := range findings {
			severityIcon := "⚠️"
			if f.Severity == "high" {
				severityIcon = "🔴"
			} else if f.Severity == "medium" {
				severityIcon = "🟡"
			}

			autoFixable := "manual"
			switch f.PatternID {
			case "C002", "T001", "C005":
				autoFixable = "auto"
			}

			fmt.Printf("\n%s [%s] %s - %s:%d\n", severityIcon, f.PatternID, strings.ToUpper(f.Type), f.File, f.Line)
			fmt.Printf("   BEFORE: %s\n", slopTruncate(f.Match, 60))
			fmt.Printf("   FIX:    %s\n", slopTruncate(f.Fix, 60))
			fmt.Printf("   STATUS: %s\n", map[string]string{"auto": "✅ Will fix", "manual": "⚠️ Manual review"}[autoFixable])
		}
	}

	// Summary
	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")

	if dryRun {
		fmt.Println("")
		fmt.Println("💡 To apply fixes, run: vic slop fix --dry-run=false")
	} else {
		fmt.Println("")
		fmt.Printf("✅ Applied %d fixes\n", fixableCount)
		fmt.Println("⚠️  Manual review required for remaining issues")
	}

	return nil
}
