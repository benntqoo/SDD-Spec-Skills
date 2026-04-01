package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/checker"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
)

// NewCheckCmd creates the check command
func NewCheckCmd(cfg *config.Config) *cobra.Command {
	var jsonOutput bool
	var category string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check code alignment with decisions",
		Long: `Check if code aligns with recorded technical decisions.

This command analyzes your source code to verify that the implemented 
technologies match your recorded technical decisions.`,
		Example: `  vic check
  vic check --json
  vic check --category database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(cfg, jsonOutput, category)
		},
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "", false, "Output as JSON")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter by category (database, auth, frontend, backend)")

	return cmd
}

func runCheck(cfg *config.Config, jsonOutput bool, categoryFilter string) error {
	if jsonOutput {
		return runCheckJSON(cfg, categoryFilter)
	}

	fmt.Println("🔍 Checking code alignment...")
	fmt.Println("----------------------------------------")

	// Load tech records
	records, err := utils.LoadTechRecords(cfg)
	if err != nil {
		return fmt.Errorf("failed to load tech records: %w", err)
	}

	if len(records.TechRecords) == 0 {
		fmt.Println("No technical decisions recorded yet.")
		return nil
	}

	// Create analyzer and scan code
	analyzer := checker.NewCodeAnalyzer()
	projectDir := cfg.ProjectDir
	if err := analyzer.ScanDirectory(projectDir); err != nil {
		fmt.Printf("⚠️  Warning: Could not scan code: %v\n", err)
	}

	// Get detected technologies
	detected := analyzer.GetDetectedTech()
	if len(detected) > 0 {
		fmt.Println("\n📦 Detected Technologies:")
		for cat, techs := range detected {
			fmt.Printf("   %s: %s\n", cat, strings.Join(techs, ", "))
		}
	}

	// Check each record
	fmt.Println("\n📋 Checking Decisions:")
	fmt.Println("----------------------------------------")

	var results []checker.CheckResult
	for _, r := range records.TechRecords {
		// Filter by category if specified
		if categoryFilter != "" && strings.ToLower(r.Category) != strings.ToLower(categoryFilter) {
			continue
		}

		result := analyzer.CheckDecision(r.ID, r.Category, r.Decision)
		results = append(results, result)

		// Print result
		icon := map[checker.CheckStatus]string{
			checker.StatusPass:    "✅",
			checker.StatusFail:    "❌",
			checker.StatusSkip:    "⏭️",
			checker.StatusUnknown: "❓",
		}[result.Status]

		fmt.Printf("\n%s [%s] %s\n", icon, r.ID, r.Title)
		fmt.Printf("   Category: %s\n", r.Category)
		fmt.Printf("   Decision: %s\n", r.Decision)
		fmt.Printf("   Status: %s\n", result.Status)
		fmt.Printf("   Message: %s\n", result.Message)
	}

	// Summary
	fmt.Println("\n" + "========================================")
	fmt.Println("Summary")
	fmt.Println("========================================")

	passCount := 0
	failCount := 0
	skipCount := 0
	unknownCount := 0

	for _, r := range results {
		switch r.Status {
		case checker.StatusPass:
			passCount++
		case checker.StatusFail:
			failCount++
		case checker.StatusSkip:
			skipCount++
		case checker.StatusUnknown:
			unknownCount++
		}
	}

	fmt.Printf("✅ Pass: %d\n", passCount)
	fmt.Printf("❌ Fail: %d\n", failCount)
	fmt.Printf("⏭️ Skip: %d\n", skipCount)
	fmt.Printf("❓ Unknown: %d\n", unknownCount)

	if failCount > 0 {
		fmt.Println("\n⚠️  Some decisions are not reflected in code!")
		fmt.Println("   Review failed checks and either update code or records.")
		showCheckUpdateRecommendation(failCount)
	}

	if passCount == 0 && unknownCount > 0 {
		fmt.Println("\n💡 Tip: Unknown status means no matching code patterns were found.")
		fmt.Println("   This is normal for planned/in-progress decisions.")
	}

	return nil
}

func runCheckJSON(cfg *config.Config, categoryFilter string) error {
	records, err := utils.LoadTechRecords(cfg)
	if err != nil {
		return fmt.Errorf("failed to load tech records: %w", err)
	}

	analyzer := checker.NewCodeAnalyzer()
	projectDir := cfg.ProjectDir
	analyzer.ScanDirectory(projectDir)

	var results []checker.CheckResult
	for _, r := range records.TechRecords {
		if categoryFilter != "" && strings.ToLower(r.Category) != strings.ToLower(categoryFilter) {
			continue
		}
		result := analyzer.CheckDecision(r.ID, r.Category, r.Decision)
		results = append(results, result)
	}

	output := map[string]interface{}{
		"results":  results,
		"detected": analyzer.GetDetectedTech(),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// NewValidateCmd creates the validate command
func NewValidateCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "validate",
		Short:   "Full validation (check + fold)",
		Long:    `Run full validation: code alignment check + event folding.`,
		Example: `  vic validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔍 Step 1: Code Alignment Check")
			fmt.Println("----------------------------------------")

			if err := runCheck(cfg, false, ""); err != nil {
				fmt.Printf("⚠️  Warning: %v\n", err)
			}

			fmt.Println("\n📦 Step 2: Fold Events")
			fmt.Println("----------------------------------------")

			if err := runFold(cfg); err != nil {
				return err
			}

			fmt.Println("\n========================================")
			fmt.Println("✅ All validations passed!")
			return nil
		},
	}
}

// NewFoldCmd creates the fold command
func NewFoldCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "fold",
		Short:   "Fold events to state",
		Long:    `Fold event history into current state snapshot.`,
		Example: `  vic fold`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFold(cfg)
		},
	}
}

func runFold(cfg *config.Config) error {
	// Load events
	events, err := utils.LoadEvents(cfg)
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}

	if len(events.Events) == 0 {
		fmt.Println("No events to fold.")
		return nil
	}

	// Load tech records
	records, _ := utils.LoadTechRecords(cfg)

	// Load risks
	risks, _ := utils.LoadRiskZones(cfg)

	// Count active
	activeDecisions := 0
	for _, r := range records.TechRecords {
		if r.Status != "completed" && r.Status != "deprecated" {
			activeDecisions++
		}
	}

	activeRisks := 0
	for _, r := range risks.Risks {
		if r.Status != "resolved" && r.Status != "accepted" {
			activeRisks++
		}
	}

	// Create state
	state := types.StateFile{
		State: types.State{
			LastFolded:      fmt.Sprintf("%d events processed", len(events.Events)),
			ActiveDecisions: activeDecisions,
			ActiveRisks:     activeRisks,
			TechRecords:     records.TechRecords,
			Risks:           risks.Risks,
		},
	}

	// Save state
	if err := utils.SaveState(cfg, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("✅ Folded %d events to state\n", len(events.Events))
	fmt.Printf("   Active decisions: %d\n", activeDecisions)
	fmt.Printf("   Active risks: %d\n", activeRisks)

	return nil
}

// NewStatusCmd creates the status command
func NewStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show project status",
		Long:    `Show current project status.`,
		Example: `  vic status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(cfg)
		},
	}
}

func runStatus(cfg *config.Config) error {
	fmt.Println("📊 Project Status")
	fmt.Println("========================================")

	// Check if VIC directory exists
	if !utils.FileExists(cfg.GetVICDir()) {
		fmt.Println("⚠️  Project not initialized")
		fmt.Println("\nRun 'vic init' to initialize")
		return nil
	}

	// Load state
	state, err := utils.LoadState(cfg)
	if err != nil {
		state = types.StateFile{}
	}

	// Load events
	events, _ := utils.LoadEvents(cfg)

	// Load tech records
	records, _ := utils.LoadTechRecords(cfg)

	// Load risks
	risks, _ := utils.LoadRiskZones(cfg)

	fmt.Printf("Last folded: %s\n", state.State.LastFolded)
	fmt.Printf("Total events: %d\n", len(events.Events))
	fmt.Printf("Active decisions: %d\n", countActiveDecisions(records.TechRecords))
	fmt.Printf("Active risks: %d\n", countActiveRisks(risks.Risks))

	// Show recent records
	if len(records.TechRecords) > 0 {
		fmt.Println("\n📋 Recent Tech Records:")
		show := records.TechRecords
		if len(show) > 5 {
			show = show[len(show)-5:]
		}
		for _, r := range show {
			fmt.Printf("   [%s] %s\n", r.ID, r.Title)
		}
	}

	return nil
}

func countActiveDecisions(records []types.TechRecord) int {
	count := 0
	for _, r := range records {
		if r.Status != "completed" && r.Status != "deprecated" {
			count++
		}
	}
	return count
}

func countActiveRisks(risks []types.RiskRecord) int {
	count := 0
	for _, r := range risks {
		if r.Status != "resolved" && r.Status != "accepted" {
			count++
		}
	}
	return count
}

// NewSearchCmd creates the search command
func NewSearchCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search records",
		Long:  `Search tech records and risk records.`,
		Example: `  vic search postgres
  vic search authentication`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a search query")
			}
			return runSearch(cfg, args[0])
		},
	}
}

func runSearch(cfg *config.Config, query string) error {
	query = strings.ToLower(query)
	fmt.Printf("🔍 Searching for: %s\n", query)
	fmt.Println("========================================")

	found := false

	// Search tech records
	records, err := utils.LoadTechRecords(cfg)
	if err == nil {
		for _, r := range records.TechRecords {
			if contains(query, r.ID, r.Title, r.Decision, r.Category, r.Reason) {
				fmt.Printf("\n📝 Tech Record: [%s] %s\n", r.ID, r.Title)
				fmt.Printf("   Decision: %s\n", r.Decision)
				fmt.Printf("   Status: %s\n", r.Status)
				found = true
			}
		}
	}

	// Search risks
	risks, err := utils.LoadRiskZones(cfg)
	if err == nil {
		for _, r := range risks.Risks {
			if contains(query, r.ID, r.Area, r.Description, r.Category) {
				fmt.Printf("\n⚠️  Risk: [%s] %s\n", r.ID, r.Area)
				fmt.Printf("   Description: %s\n", r.Description)
				fmt.Printf("   Impact: %s\n", r.Impact)
				found = true
			}
		}
	}

	if !found {
		fmt.Println("No matching records found.")
	}

	return nil
}

func contains(query string, fields ...string) bool {
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), query) {
			return true
		}
	}
	return false
}

// NewHistoryCmd creates the history command
func NewHistoryCmd(cfg *config.Config) *cobra.Command {
	var limit int
	var eventType string

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show event history",
		Long:  `Show event history.`,
		Example: `  vic history --limit 5
  vic history --type decision_made`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(cfg, limit, eventType)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of events to show")
	cmd.Flags().StringVarP(&eventType, "type", "t", "", "Filter by event type")

	return cmd
}

func runHistory(cfg *config.Config, limit int, eventType string) error {
	fmt.Println("📜 Event History")
	fmt.Println("========================================")

	events, err := utils.LoadEvents(cfg)
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}

	if len(events.Events) == 0 {
		fmt.Println("No events recorded yet.")
		return nil
	}

	// Filter by type
	var filtered []types.Event
	for _, e := range events.Events {
		if eventType == "" || e.Type == eventType {
			filtered = append(filtered, e)
		}
	}

	// Reverse to show newest first
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	// Limit
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	for _, e := range filtered {
		fmt.Printf("\n[%s] %s\n", e.Timestamp.Format("2006-01-02 15:04"), e.Type)
		if id, ok := e.Data["id"].(string); ok {
			fmt.Printf("   ID: %s\n", id)
		}
		if title, ok := e.Data["title"].(string); ok {
			fmt.Printf("   Title: %s\n", title)
		}
	}

	fmt.Printf("\nShowing %d of %d events\n", len(filtered), len(events.Events))

	return nil
}

// NewExportCmd creates the export command
func NewExportCmd(cfg *config.Config) *cobra.Command {
	var output string
	var exportType string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export data",
		Long:  `Export .vic-sdd/ data to JSON file.`,
		Example: `  vic export --output backup.json
  vic export --type tech -o tech-decisions.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if output == "" {
				output = "vibe-integrity-export.json"
			}
			return runExport(cfg, output, exportType)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	cmd.Flags().StringVarP(&exportType, "type", "t", "", "Export type (tech/risks/events/all)")

	return cmd
}

func runExport(cfg *config.Config, output, exportType string) error {
	fmt.Printf("📤 Exporting to: %s\n", output)

	data := make(map[string]interface{})

	// Export based on type
	if exportType == "" || exportType == "all" || exportType == "tech" {
		if records, err := utils.LoadTechRecords(cfg); err == nil {
			data["tech_records"] = records.TechRecords
		}
	}

	if exportType == "" || exportType == "all" || exportType == "risks" {
		if risks, err := utils.LoadRiskZones(cfg); err == nil {
			data["risks"] = risks.Risks
		}
	}

	if exportType == "" || exportType == "all" || exportType == "events" {
		if events, err := utils.LoadEvents(cfg); err == nil {
			data["events"] = events.Events
		}
	}

	// Write to file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(output, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("✅ Export complete")
	return nil
}

// NewImportCmd creates the import command
func NewImportCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import <file>",
		Short:   "Import data",
		Long:    `Import data from JSON file.`,
		Example: `  vic import backup.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires an input file")
			}
			return runImport(cfg, args[0])
		},
	}

	return cmd
}

func runImport(cfg *config.Config, input string) error {
	fmt.Printf("📥 Importing from: %s\n", input)

	// Read file
	data, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var importData map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Import tech records
	if techRecords, ok := importData["tech_records"]; ok {
		records, _ := utils.LoadTechRecords(cfg)
		newRecords := techRecords.([]interface{})
		for _, nr := range newRecords {
			if nrMap, ok := nr.(map[string]interface{}); ok {
				record := types.TechRecord{
					ID:       toString(nrMap["id"]),
					Title:    toString(nrMap["title"]),
					Decision: toString(nrMap["decision"]),
					Category: toString(nrMap["category"]),
					Reason:   toString(nrMap["reason"]),
					Impact:   toString(nrMap["impact"]),
					Status:   toString(nrMap["status"]),
				}
				records.TechRecords = append(records.TechRecords, record)
			}
		}
		utils.SaveTechRecords(cfg, records)
		fmt.Printf("   Imported %d tech records\n", len(newRecords))
	}

	// Import risks
	if riskRecords, ok := importData["risks"]; ok {
		risks, _ := utils.LoadRiskZones(cfg)
		newRisks := riskRecords.([]interface{})
		for _, nr := range newRisks {
			if nrMap, ok := nr.(map[string]interface{}); ok {
				risk := types.RiskRecord{
					ID:          toString(nrMap["id"]),
					Area:        toString(nrMap["area"]),
					Description: toString(nrMap["description"]),
					Category:    toString(nrMap["category"]),
					Impact:      toString(nrMap["impact"]),
					Status:      toString(nrMap["status"]),
				}
				risks.Risks = append(risks.Risks, risk)
			}
		}
		utils.SaveRiskZones(cfg, risks)
		fmt.Printf("   Imported %d risks\n", len(newRisks))
	}

	fmt.Println("✅ Import complete")
	return nil
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// showCheckUpdateRecommendation prints recommended actions when check fails
func showCheckUpdateRecommendation(failCount int) {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("📋 SPEC UPDATE RECOMMENDATION")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Printf("Found %d decision(s) not reflected in code.\n", failCount)
	fmt.Println("To resolve this, choose one of the following:")
	fmt.Println("1️⃣  Update SPEC (Recommended)")
	fmt.Println("    $ vic spec update --file SPEC-ARCHITECTURE.md --section \"[section]\"")
	fmt.Println("    Then: vic check")
	fmt.Println("2️⃣  Update code to match decisions")
	fmt.Println("    $ git diff [affected files]")
	fmt.Println("    Then: Re-implement correctly")
	fmt.Println("3️⃣  Document as accepted drift (requires approval)")
	fmt.Println("    $ vic rr --id DRIFT-[DATE] --desc \"[description]\"")
	fmt.Println("    ⚠️  Only for emergency hotfixes")
	fmt.Println("For more details, see: skills/constitution-check/SKILL.md")
	fmt.Println("════════════════════════════════════════════════════════════")
}
