package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// Skill metadata
type SkillInfo struct {
	Name        string
	Description string
	Category    string
	Path        string
}

// NewSkillCmd creates the skill command
func NewSkillCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "View skill documentation",
		Long: `View skill documentation.
		
Skills are markdown files that guide AI behavior.
Use this command to understand what each skill does.`,
		Example: `  vic skill list               # List all available skills
  vic skill show requirements   # View requirements skill
  vic skill show architecture   # View architecture skill
  vic skill help requirements   # Quick help for a skill`,
	}

	cmd.AddCommand(NewSkillListCmd(cfg))
	cmd.AddCommand(NewSkillShowCmd(cfg))
	cmd.AddCommand(NewSkillHelpCmd(cfg))

	return cmd
}

// NewSkillListCmd lists all available skills
func NewSkillListCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all available skills",
		Long:    `List all available skills with their descriptions.`,
		Example: `  vic skill list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillList(cfg)
		},
	}
}

func runSkillList(cfg *config.Config) error {
	fmt.Println("📚 Available Skills")
	fmt.Println("========================================")
	fmt.Println()

	// Find skills directory
	skillsDir := filepath.Join(cfg.ProjectDir, "skills")
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		fmt.Println("❌ No skills directory found")
		fmt.Println("   Skills should be in: ./skills/")
		return nil
	}

	// Read skill metadata
	skills := getSkillList(skillsDir)

	if len(skills) == 0 {
		fmt.Println("❌ No skills found")
		return nil
	}

	// Group by category
	categories := make(map[string][]SkillInfo)
	for _, skill := range skills {
		categories[skill.Category] = append(categories[skill.Category], skill)
	}

	// Print grouped skills
	for category, skillList := range categories {
		fmt.Printf("📁 %s\n", category)
		for _, skill := range skillList {
			fmt.Printf("   • %-20s %s\n", skill.Name, skill.Description)
		}
		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Printf("Total: %d skills\n\n", len(skills))
	fmt.Println("View a skill:")
	fmt.Printf("   vic skill show <name>\n")
	fmt.Println("   Example: vic skill show requirements")

	return nil
}

// getSkillList reads all skills from the skills directory
func getSkillList(skillsDir string) []SkillInfo {
	skills := make([]SkillInfo, 0)

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return skills
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsDir, entry.Name())
		skillFile := filepath.Join(skillPath, "SKILL.md")

		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			continue
		}

		// Read skill description from first non-header lines
		description := readSkillDescription(skillFile)

		// Determine category based on name
		category := getSkillCategory(entry.Name())

		skills = append(skills, SkillInfo{
			Name:        entry.Name(),
			Description: description,
			Category:    category,
			Path:        skillFile,
		})
	}

	return skills
}

// readSkillDescription extracts the first line description from a skill file
func readSkillDescription(skillFile string) string {
	content, err := os.ReadFile(skillFile)
	if err != nil {
		return "No description"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip headers and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		// Return first non-empty, non-header line
		if len(line) > 80 {
			return line[:77] + "..."
		}
		return line
	}

	return "No description"
}

// getSkillCategory determines the category of a skill
func getSkillCategory(skillName string) string {
	categories := map[string][]string{
		"Self-Awareness": {"context-tracker"},
		"Vibe":           {"requirements", "architecture", "design-review", "debugging"},
		"QA":             {"qa"},
		"SDD":            {"sdd-orchestrator", "spec-architect", "spec-contract-diff", "spec-traceability"},
	}

	for category, skills := range categories {
		for _, skill := range skills {
			if skill == skillName {
				return category
			}
		}
	}

	return "Other"
}

// NewSkillShowCmd shows a specific skill
func NewSkillShowCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "show <skill-name>",
		Short: "Show skill content",
		Long:  `Show the full content of a skill file.`,
		Example: `  vic skill show requirements
  vic skill show architecture`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Println("❌ Please specify a skill name")
				fmt.Println("   Usage: vic skill show <skill-name>")
				fmt.Println("   Run 'vic skill list' to see available skills")
				return nil
			}
			return runSkillShow(cfg, args[0])
		},
	}
}

func runSkillShow(cfg *config.Config, skillName string) error {
	skillsDir := filepath.Join(cfg.ProjectDir, "skills")
	skillFile := filepath.Join(skillsDir, skillName, "SKILL.md")

	// Check if skill exists
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		fmt.Printf("❌ Skill '%s' not found\n", skillName)
		fmt.Println()
		fmt.Println("Run 'vic skill list' to see available skills:")

		// Suggest similar skills
		similar := findSimilarSkills(skillsDir, skillName)
		if len(similar) > 0 {
			fmt.Println("Did you mean:")
			for _, s := range similar {
				fmt.Printf("  • %s\n", s)
			}
		}
		return nil
	}

	// Read and display skill
	content, err := os.ReadFile(skillFile)
	if err != nil {
		return fmt.Errorf("failed to read skill: %w", err)
	}

	fmt.Printf("📄 Skill: %s\n", skillName)
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println(string(content))

	return nil
}

// findSimilarSkills finds skills with similar names
func findSimilarSkills(skillsDir, query string) []string {
	query = strings.ToLower(query)
	similar := make([]string, 0)

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return similar
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := strings.ToLower(entry.Name())
			// Simple similarity check
			if strings.Contains(name, query) ||
				strings.HasPrefix(name, query) ||
				queryHasCommonPrefix(name, query) {
				similar = append(similar, entry.Name())
			}
		}
	}

	return similar
}

// queryHasCommonPrefix checks if two strings share a prefix
func queryHasCommonPrefix(a, b string) bool {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	count := 0
	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			count++
		} else {
			break
		}
	}

	return count >= 3 // At least 3 common chars
}

// NewSkillHelpCmd shows quick help for a skill
func NewSkillHelpCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "help <skill-name>",
		Short: "Show quick help for a skill",
		Long:  `Show quick help summary for a skill (first section only).`,
		Example: `  vic skill help requirements
  vic skill help architecture`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Println("❌ Please specify a skill name")
				fmt.Println("   Usage: vic skill help <skill-name>")
				return nil
			}
			return runSkillHelp(cfg, args[0])
		},
	}
}

func runSkillHelp(cfg *config.Config, skillName string) error {
	skillsDir := filepath.Join(cfg.ProjectDir, "skills")
	skillFile := filepath.Join(skillsDir, skillName, "SKILL.md")

	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		fmt.Printf("❌ Skill '%s' not found\n", skillName)
		fmt.Println("   Run 'vic skill list' to see available skills")
		return nil
	}

	content, err := os.ReadFile(skillFile)
	if err != nil {
		return fmt.Errorf("failed to read skill: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Print first section only (until next ## or empty section)
	fmt.Printf("📋 Quick Help: %s\n", skillName)
	fmt.Println("========================================")

	sectionEnded := false
	printedLines := 0
	maxLines := 30 // Limit output

	for _, line := range lines {
		if printedLines >= maxLines {
			fmt.Println("...")
			break
		}

		// Stop at second ## heading
		if strings.HasPrefix(strings.TrimSpace(line), "##") {
			if sectionEnded {
				break
			}
			sectionEnded = true
		}

		fmt.Println(line)
		printedLines++
	}

	return nil
}
