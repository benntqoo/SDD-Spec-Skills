package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/utils"
)

// ChangeLogEntry represents a single SPEC change entry
type ChangeLogEntry struct {
	ID       string   `json:"id"`
	Date     string   `json:"date"`
	File     string   `json:"file"`
	Type     string   `json:"type"` // "tech_stack", "api", "module", "security", "other"
	Changes  []string `json:"changes"`
	Impact   string   `json:"impact"` // "high", "medium", "low"
	Reviewed bool     `json:"reviewed"`
	Reviewer string   `json:"reviewer,omitempty"`
	Note     string   `json:"note,omitempty"`
}

// ChangeLog tracks all SPEC changes
type ChangeLog struct {
	Version     string           `json:"version"`
	LastUpdated string           `json:"last_updated"`
	Entries     []ChangeLogEntry `json:"entries"`
}

// specHash stores the last known hash of SPEC files
type specHash struct {
	SpecRequirements string `json:"spec_requirements"`
	SpecArchitecture string `json:"spec_architecture"`
	Design           string `json:"design"`
}

// RunSpecWatch monitors SPEC files for changes and logs them
func RunSpecWatch(cfg *config.Config) error {
	fmt.Println("👁️  SPEC Watch Mode")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Monitoring SPEC files for changes...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Load or create hash file
	hashFile := filepath.Join(cfg.VICDir, "status", "spec-hash.json")
	lastHash := loadSpecHash(hashFile)

	// Get current hashes
	currentHash := getCurrentSpecHash(cfg)

	// Check for changes
	if hasChanges(lastHash, currentHash) {
		fmt.Println("📝 Changes detected!")

		// Detect what changed
		changes := detectChanges(cfg, lastHash, currentHash)

		// Log changes
		if err := logChanges(cfg, changes); err != nil {
			fmt.Printf("⚠️  Failed to log changes: %v\n", err)
		}

		// Run drift detection
		fmt.Println()
		fmt.Println("🔍 Running tech drift detection...")
		if err := RunGate2(cfg, "plain"); err != nil {
			fmt.Printf("⚠️  Gate 2 check failed: %v\n", err)
		}

		// Save new hashes
		saveSpecHash(hashFile, currentHash)
	} else {
		fmt.Println("✅ No changes detected")
	}

	return nil
}

// RunSpecChanges shows change history
func RunSpecChanges(cfg *config.Config) error {
	fmt.Println("📋 SPEC Change History")
	fmt.Println("========================================")
	fmt.Println()

	changeLog := loadChangeLog(cfg)

	if len(changeLog.Entries) == 0 {
		fmt.Println("No SPEC changes recorded yet")
		fmt.Println()
		fmt.Println("Run 'vic spec watch' or 'vic spec diff' to detect changes")
		return nil
	}

	// Print entries in reverse order (newest first)
	for i := len(changeLog.Entries) - 1; i >= 0; i-- {
		entry := changeLog.Entries[i]

		statusIcon := "📝"
		if entry.Reviewed {
			statusIcon = "✅"
		}

		impactIcon := "🔴"
		if entry.Impact == "medium" {
			impactIcon = "🟡"
		} else if entry.Impact == "low" {
			impactIcon = "🟢"
		}

		fmt.Printf("[%s] %s %s - %s\n", statusIcon, impactIcon, entry.Date, entry.File)
		fmt.Printf("    Type: %s\n", entry.Type)
		for _, change := range entry.Changes {
			fmt.Printf("    • %s\n", change)
		}
		if entry.Note != "" {
			fmt.Printf("    Note: %s\n", entry.Note)
		}
		fmt.Println()
	}

	return nil
}

// RunSpecDiff detects and displays SPEC changes
func RunSpecDiff(cfg *config.Config) error {
	fmt.Println("🔄 SPEC Diff Detection")
	fmt.Println("========================================")
	fmt.Println()

	// Load hash file
	hashFile := filepath.Join(cfg.VICDir, "status", "spec-hash.json")
	lastHash := loadSpecHash(hashFile)
	currentHash := getCurrentSpecHash(cfg)

	if !hasChanges(lastHash, currentHash) {
		fmt.Println("✅ No changes since last check")
		fmt.Println("   Run 'vic spec watch' first to start monitoring")
		return nil
	}

	// Detect what changed
	changes := detectChanges(cfg, lastHash, currentHash)

	fmt.Printf("📝 %d change(s) detected:\n\n", len(changes))

	for _, change := range changes {
		fmt.Printf("  [%s] %s\n", change.Type, change.File)
		for _, c := range change.Changes {
			fmt.Printf("      • %s\n", c)
		}
		fmt.Printf("      Impact: %s\n\n", change.Impact)
	}

	// Prompt to log
	fmt.Println("Run 'vic spec log --accept' to accept these changes")
	fmt.Println("Or review manually and run 'vic spec log --add \"change description\"'")

	return nil
}

// NewSpecWatchCmd creates the spec watch command
func NewSpecWatchCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Short:   "Watch SPEC files for changes",
		Long:    `Monitor SPEC files and detect changes, auto-run tech drift check.`,
		Example: `  vic spec watch`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSpecWatch(cfg)
		},
	}
}

// NewSpecChangesCmd creates the spec changes command
func NewSpecChangesCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "changes",
		Short:   "Show SPEC change history",
		Long:    `Display all recorded SPEC changes.`,
		Example: `  vic spec changes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSpecChanges(cfg)
		},
	}
}

// NewSpecDiffCmd creates the spec diff command
func NewSpecDiffCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "diff",
		Short:   "Detect SPEC changes since last check",
		Long:    `Compare current SPEC files with last known state.`,
		Example: `  vic spec diff`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSpecDiff(cfg)
		},
	}
}

// loadSpecHash loads the last known SPEC hashes
func loadSpecHash(hashFile string) specHash {
	var hash specHash

	data, err := os.ReadFile(hashFile)
	if err != nil {
		return specHash{}
	}

	json.Unmarshal(data, &hash)
	return hash
}

// saveSpecHash saves the current SPEC hashes
func saveSpecHash(hashFile string, hash specHash) {
	os.MkdirAll(filepath.Dir(hashFile), 0755)

	data, _ := json.MarshalIndent(hash, "", "  ")
	os.WriteFile(hashFile, data, 0644)
}

// getCurrentSpecHash computes hashes of current SPEC files
func getCurrentSpecHash(cfg *config.Config) specHash {
	return specHash{
		SpecRequirements: hashFile(cfg.SpecRequirements),
		SpecArchitecture: hashFile(cfg.SpecArchitecture),
		Design:           hashFile("DESIGN.md"),
	}
}

// hashFile computes SHA256 hash of a file
func hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// hasChanges compares two hash sets
func hasChanges(old, new specHash) bool {
	return old.SpecRequirements != new.SpecRequirements ||
		old.SpecArchitecture != new.SpecArchitecture ||
		old.Design != new.Design
}

// detectChanges identifies what specifically changed
func detectChanges(cfg *config.Config, old, new specHash) []ChangeLogEntry {
	var changes []ChangeLogEntry

	// Check SPEC-REQUIREMENTS.md
	if old.SpecRequirements != new.SpecRequirements {
		if utils.FileExists(cfg.SpecRequirements) {
			content, _ := os.ReadFile(cfg.SpecRequirements)
			contentStr := string(content)

			entry := ChangeLogEntry{
				ID:       fmt.Sprintf("REQ-%d", time.Now().Unix()),
				Date:     time.Now().Format("2006-01-02 15:04"),
				File:     "SPEC-REQUIREMENTS.md",
				Type:     detectChangeType(contentStr),
				Changes:  extractChanges(contentStr),
				Impact:   assessImpact(contentStr),
				Reviewed: false,
			}
			changes = append(changes, entry)
		}
	}

	// Check SPEC-ARCHITECTURE.md
	if old.SpecArchitecture != new.SpecArchitecture {
		if utils.FileExists(cfg.SpecArchitecture) {
			content, _ := os.ReadFile(cfg.SpecArchitecture)
			contentStr := string(content)

			entry := ChangeLogEntry{
				ID:       fmt.Sprintf("ARCH-%d", time.Now().Unix()),
				Date:     time.Now().Format("2006-01-02 15:04"),
				File:     "SPEC-ARCHITECTURE.md",
				Type:     detectArchChangeType(contentStr),
				Changes:  extractArchChanges(contentStr),
				Impact:   assessArchImpact(contentStr),
				Reviewed: false,
			}
			changes = append(changes, entry)
		}
	}

	// Check DESIGN.md
	if old.Design != new.Design {
		if _, err := os.Stat("DESIGN.md"); err == nil {
			entry := ChangeLogEntry{
				ID:       fmt.Sprintf("DES-%d", time.Now().Unix()),
				Date:     time.Now().Format("2006-01-02 15:04"),
				File:     "DESIGN.md",
				Type:     "design",
				Changes:  []string{"Design system updated"},
				Impact:   "medium",
				Reviewed: false,
			}
			changes = append(changes, entry)
		}
	}

	return changes
}

// detectChangeType identifies the type of requirements change
func detectChangeType(content string) string {
	lower := strings.ToLower(content)

	if strings.Contains(lower, "user story") || strings.Contains(lower, "feature") {
		return "requirements"
	}
	if strings.Contains(lower, "acceptance criteria") || strings.Contains(lower, "验收标准") {
		return "criteria"
	}
	if strings.Contains(lower, "non-functional") || strings.Contains(lower, "performance") {
		return "quality"
	}

	return "other"
}

// detectArchChangeType identifies the type of architecture change
func detectArchChangeType(content string) string {
	lower := strings.ToLower(content)

	if strings.Contains(lower, "technology") || strings.Contains(lower, "tech stack") || strings.Contains(lower, "技术栈") {
		return "tech_stack"
	}
	if strings.Contains(lower, "api") || strings.Contains(lower, "endpoint") {
		return "api"
	}
	if strings.Contains(lower, "module") || strings.Contains(lower, "component") || strings.Contains(lower, "组件") {
		return "module"
	}
	if strings.Contains(lower, "security") || strings.Contains(lower, "auth") {
		return "security"
	}
	if strings.Contains(lower, "database") || strings.Contains(lower, "data model") {
		return "data"
	}

	return "other"
}

// extractChanges extracts key changes from requirements
func extractChanges(content string) []string {
	var changes []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]") {
			task := strings.TrimPrefix(trimmed, "- [ ]")
			task = strings.TrimPrefix(task, "- [x]")
			task = strings.TrimSpace(task)
			if task != "" && !strings.Contains(task, "...") {
				changes = append(changes, fmt.Sprintf("Task: %s", task))
			}
		}
	}

	if len(changes) == 0 {
		changes = append(changes, "Content updated")
	}

	return changes
}

// extractArchChanges extracts key changes from architecture
func extractArchChanges(content string) []string {
	var changes []string

	lower := strings.ToLower(content)

	// Check for tech stack changes
	techPatterns := []string{"postgresql", "mysql", "mongodb", "redis", "react", "vue", "angular", "go", "python", "rust", "docker", "kubernetes"}
	for _, tech := range techPatterns {
		if strings.Contains(lower, tech) {
			changes = append(changes, fmt.Sprintf("Technology: %s", strings.ToUpper(tech)))
		}
	}

	// Check for API changes
	if strings.Contains(lower, "api") || strings.Contains(lower, "endpoint") {
		changes = append(changes, "API endpoints defined")
	}

	// Check for security
	if strings.Contains(lower, "security") || strings.Contains(lower, "auth") {
		changes = append(changes, "Security requirements defined")
	}

	if len(changes) == 0 {
		changes = append(changes, "Architecture updated")
	}

	return changes
}

// assessImpact determines the impact level of a change
func assessImpact(content string) string {
	lower := strings.ToLower(content)

	highImpact := []string{"performance", "security", "scalability", "breaking"}
	mediumImpact := []string{"feature", "user story", "acceptance"}

	for _, keyword := range highImpact {
		if strings.Contains(lower, keyword) {
			return "high"
		}
	}

	for _, keyword := range mediumImpact {
		if strings.Contains(lower, keyword) {
			return "medium"
		}
	}

	return "low"
}

// assessArchImpact determines impact of architecture change
func assessArchImpact(content string) string {
	lower := strings.ToLower(content)

	highImpact := []string{"security", "auth", "database", "breaking"}
	mediumImpact := []string{"technology", "tech", "api", "endpoint", "module"}

	for _, keyword := range highImpact {
		if strings.Contains(lower, keyword) {
			return "high"
		}
	}

	for _, keyword := range mediumImpact {
		if strings.Contains(lower, keyword) {
			return "medium"
		}
	}

	return "low"
}

// loadChangeLog loads the change log
func loadChangeLog(cfg *config.Config) ChangeLog {
	logFile := filepath.Join(cfg.VICDir, "status", "change-log.yaml")

	data, err := os.ReadFile(logFile)
	if err != nil {
		return ChangeLog{
			Version:     "1.0",
			LastUpdated: time.Now().Format("2006-01-02"),
			Entries:     []ChangeLogEntry{},
		}
	}

	// Try JSON first
	var log ChangeLog
	if err := json.Unmarshal(data, &log); err != nil {
		// Return empty log on parse error
		return ChangeLog{
			Version:     "1.0",
			LastUpdated: time.Now().Format("2006-01-02"),
			Entries:     []ChangeLogEntry{},
		}
	}

	return log
}

// logChanges saves detected changes to the log
func logChanges(cfg *config.Config, changes []ChangeLogEntry) error {
	if len(changes) == 0 {
		return nil
	}

	logFile := filepath.Join(cfg.VICDir, "status", "change-log.yaml")
	os.MkdirAll(filepath.Dir(logFile), 0755)

	// Load existing log
	log := loadChangeLog(cfg)

	// Append new entries
	log.Entries = append(log.Entries, changes...)
	log.LastUpdated = time.Now().Format("2006-01-02")

	// Save
	data, _ := json.MarshalIndent(log, "", "  ")
	return os.WriteFile(logFile, data, 0644)
}
