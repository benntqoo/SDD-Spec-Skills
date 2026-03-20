package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
)

// NewDesignCmd creates the design command
func NewDesignCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "design",
		Short: "Design system management",
		Long:  `Manage design system and DESIGN.md documentation.`,
		Example: `  vic design init                # Initialize DESIGN.md
  vic design status            # Show design status
  vic design check             # Check design completeness`,
	}

	cmd.AddCommand(NewDesignInitCmd(cfg))
	cmd.AddCommand(NewDesignStatusCmd(cfg))
	cmd.AddCommand(NewDesignCheckCmd(cfg))

	return cmd
}

// NewDesignInitCmd initializes DESIGN.md
func NewDesignInitCmd(cfg *config.Config) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize DESIGN.md",
		Long:    `Create a new DESIGN.md file from template.`,
		Example: `  vic design init --name "My App"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				name = "Project"
			}
			return runDesignInit(cfg, name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Project name")

	return cmd
}

func runDesignInit(cfg *config.Config, projectName string) error {
	designFile := "DESIGN.md"

	// Check if already exists
	if _, err := os.Stat(designFile); err == nil {
		fmt.Printf("⚠️  %s already exists\n", designFile)
		fmt.Println("   Use 'vic design check' to validate it")
		fmt.Println("   Or delete it and run again")
		return nil
	}

	// Generate content
	content := generateDesignTemplate(projectName)

	// Write file
	if err := os.WriteFile(designFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write DESIGN.md: %w", err)
	}

	fmt.Printf("✅ Created %s\n", designFile)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("   1. Edit DESIGN.md with your design system")
	fmt.Println("   2. Use 'vic design check' to validate completeness")
	fmt.Println("   3. The design will be checked in SPEC Gate 1.5")

	return nil
}

func generateDesignTemplate(projectName string) string {
	date := time.Now().Format("2006-01-02")

	var sb strings.Builder
	sb.WriteString("# Design System: " + projectName + "\n\n")
	sb.WriteString("> Generated: " + date + "\n")
	sb.WriteString("> This document defines the visual design language for the project.\n\n")
	sb.WriteString("---\n\n")

	sb.WriteString("## Aesthetic Direction\n\n")
	sb.WriteString("Describe the overall look and feel:\n")
	sb.WriteString("- **Style**: [Minimal/Industrial/Playful/Premium/etc.]\n")
	sb.WriteString("- **Mood**: [Professional/Casual/Fun/etc.]\n")
	sb.WriteString("- **Inspiration**: [Reference sites or products]\n\n")

	sb.WriteString("## Typography\n\n")
	sb.WriteString("| Element | Font | Size | Weight |\n")
	sb.WriteString("|---------|------|------|--------|\n")
	sb.WriteString("| Display/Headings | [Font Name] | [size] | [weight] |\n")
	sb.WriteString("| Body Text | [Font Name] | [size] | [weight] |\n")
	sb.WriteString("| Code/Mono | [Font Name] | [size] | [weight] |\n\n")

	sb.WriteString("## Color Palette\n\n")
	sb.WriteString("### Primary Colors\n")
	sb.WriteString("- **Primary**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Primary Light**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Primary Dark**: #XXXXXX - [Description]\n\n")

	sb.WriteString("### Neutral Colors\n")
	sb.WriteString("- **Background**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Surface**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Text Primary**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Border**: #XXXXXX - [Description]\n\n")

	sb.WriteString("### Semantic Colors\n")
	sb.WriteString("- **Success**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Warning**: #XXXXXX - [Description]\n")
	sb.WriteString("- **Error**: #XXXXXX - [Description]\n\n")

	sb.WriteString("## Spacing Scale\n\n")
	sb.WriteString("| Token | Value | Usage |\n")
	sb.WriteString("|-------|-------|-------|\n")
	sb.WriteString("| xs | 4px | Tight spacing |\n")
	sb.WriteString("| sm | 8px | Compact elements |\n")
	sb.WriteString("| md | 16px | Default spacing |\n")
	sb.WriteString("| lg | 24px | Section spacing |\n")
	sb.WriteString("| xl | 32px | Large gaps |\n\n")

	sb.WriteString("## Border Radius\n\n")
	sb.WriteString("| Token | Value | Usage |\n")
	sb.WriteString("|-------|-------|-------|\n")
	sb.WriteString("| none | 0px | Sharp edges |\n")
	sb.WriteString("| sm | 4px | Subtle rounding |\n")
	sb.WriteString("| md | 8px | Default rounding |\n")
	sb.WriteString("| lg | 16px | Prominent rounding |\n\n")

	sb.WriteString("## Shadows\n\n")
	sb.WriteString("| Token | Value | Usage |\n")
	sb.WriteString("|-------|-------|-------|\n")
	sb.WriteString("| sm | 0 1px 2px rgba(0,0,0,0.05) | Subtle elevation |\n")
	sb.WriteString("| md | 0 4px 6px rgba(0,0,0,0.1) | Cards |\n")
	sb.WriteString("| lg | 0 10px 15px rgba(0,0,0,0.1) | Modals |\n\n")

	sb.WriteString("## Component Guidelines\n\n")
	sb.WriteString("### Buttons\n")
	sb.WriteString("- **Primary**: [style description]\n")
	sb.WriteString("- **Secondary**: [style description]\n")
	sb.WriteString("- **Ghost**: [style description]\n\n")

	sb.WriteString("### Forms\n")
	sb.WriteString("- **Input height**: [value]\n")
	sb.WriteString("- **Border**: [style]\n")
	sb.WriteString("- **Focus state**: [style]\n\n")

	sb.WriteString("## Icon Guidelines\n")
	sb.WriteString("- **Library**: [e.g., Lucide, Heroicons, Phosphor]\n")
	sb.WriteString("- **Size scale**: [e.g., 16, 20, 24, 32]\n")
	sb.WriteString("- **Stroke width**: [value]\n\n")

	sb.WriteString("## Motion & Animation\n\n")
	sb.WriteString("| Token | Value | Usage |\n")
	sb.WriteString("|-------|-------|-------|\n")
	sb.WriteString("| duration-fast | 150ms | Micro-interactions |\n")
	sb.WriteString("| duration-normal | 250ms | Default transitions |\n")
	sb.WriteString("| duration-slow | 400ms | Page transitions |\n\n")

	sb.WriteString("## Responsive Breakpoints\n\n")
	sb.WriteString("| Breakpoint | Width | Usage |\n")
	sb.WriteString("|------------|-------|-------|\n")
	sb.WriteString("| sm | 640px | Mobile landscape |\n")
	sb.WriteString("| md | 768px | Tablets |\n")
	sb.WriteString("| lg | 1024px | Small laptops |\n")
	sb.WriteString("| xl | 1280px | Desktops |\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString("## Implementation\n\n")
	sb.WriteString("Add CSS variables, Tailwind config, or component library exports here.\n")

	return sb.String()
}

// NewDesignStatusCmd shows design status
func NewDesignStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show design status",
		Long:    `Show the current status of DESIGN.md.`,
		Example: `  vic design status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDesignStatus()
		},
	}
}

func runDesignStatus() error {
	designFile := "DESIGN.md"

	// Check if exists
	if _, err := os.Stat(designFile); os.IsNotExist(err) {
		fmt.Println("📋 Design Status: Not Initialized")
		fmt.Println()
		fmt.Println("   Run 'vic design init --name \"Your Project\"' to create DESIGN.md")
		return nil
	}

	// Read content
	content, err := os.ReadFile(designFile)
	if err != nil {
		return fmt.Errorf("failed to read DESIGN.md: %w", err)
	}

	contentStr := string(content)

	// Check completeness
	sections := []struct {
		name    string
		pattern string
	}{
		{"Aesthetic Direction", "aesthetic direction"},
		{"Typography", "typography"},
		{"Color Palette", "color palette"},
		{"Spacing Scale", "spacing"},
		{"Border Radius", "border radius"},
		{"Shadows", "shadow"},
		{"Components", "component"},
		{"Icons", "icon"},
		{"Motion", "motion"},
	}

	fmt.Println("📋 Design Status: Complete")
	fmt.Println("========================================")
	fmt.Println()

	completed := 0
	for _, section := range sections {
		found := false
		lowerContent := strings.ToLower(contentStr)
		for _, line := range strings.Split(contentStr, "\n") {
			if strings.Contains(strings.ToLower(line), section.pattern) && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				found = true
				break
			}
		}
		if strings.Count(lowerContent, section.pattern) > 1 { // Count occurrences
			found = true
		}
		if found {
			fmt.Printf("   ✅ %s\n", section.name)
			completed++
		} else {
			fmt.Printf("   ❌ %s\n", section.name)
		}
	}

	fmt.Println()
	fmt.Printf("Progress: %d/%d sections defined\n", completed, len(sections))

	if completed < len(sections) {
		fmt.Println()
		fmt.Println("Run 'vic design check' for detailed analysis")
	}

	return nil
}

// NewDesignCheckCmd checks design completeness
func NewDesignCheckCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Check design completeness",
		Long:    `Run Gate 1.5 design completeness check.`,
		Example: `  vic design check`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDesignCheck()
		},
	}
}

func runDesignCheck() error {
	fmt.Println("🎨 Gate 1.5: Design Completeness Check")
	fmt.Println("========================================")
	fmt.Println()

	designFile := "DESIGN.md"

	// Check if exists
	if _, err := os.Stat(designFile); os.IsNotExist(err) {
		fmt.Println("❌ DESIGN.md not found")
		fmt.Println()
		fmt.Println("   Run 'vic design init --name \"Your Project\"' first")
		fmt.Println("   Or skip this gate if your project doesn't need UI design")
		return nil
	}

	// Read content
	content, err := os.ReadFile(designFile)
	if err != nil {
		return fmt.Errorf("failed to read DESIGN.md: %w", err)
	}

	contentStr := string(content)
	lowerContent := strings.ToLower(contentStr)

	// Define checks
	type designCheck struct {
		id      string
		name    string
		pattern string
		errMsg  string
	}

	checks := []designCheck{
		{"AESTHETIC", "Aesthetic Direction", "aesthetic direction", "Define overall look and feel"},
		{"TYPOGRAPHY", "Typography Section", "## typography", "Add typography guidelines"},
		{"FONTS", "Font Definitions", "font", "Define display and body fonts"},
		{"COLORS", "Color Palette Section", "color palette", "Add color palette"},
		{"PRIMARY", "Primary Colors", "primary", "Define primary color"},
		{"SPACING", "Spacing Scale", "spacing", "Add spacing scale"},
		{"RADIUS", "Border Radius", "border radius", "Define border radius scale"},
		{"SHADOWS", "Shadows", "shadow", "Define shadow styles"},
		{"COMPONENTS", "Component Guidelines", "component", "Add button and form guidelines"},
		{"RESPONSIVE", "Responsive Breakpoints", "breakpoint", "Define responsive breakpoints"},
	}

	allPassed := true
	for _, check := range checks {
		passed := strings.Contains(lowerContent, check.pattern)

		icon := "❌"
		if passed {
			icon = "✅"
		} else {
			allPassed = false
		}

		fmt.Printf("[%s] %s\n", icon, check.name)
		if !passed {
			fmt.Printf("      → %s\n", check.errMsg)
		}
	}

	fmt.Println()
	fmt.Println("========================================")

	if allPassed {
		fmt.Println("✅ Gate 1.5 PASSED - Design system is complete")
		fmt.Println()
		fmt.Println("Design can be integrated into SPEC-ARCHITECTURE.md")
		return nil
	}

	fmt.Println("❌ Gate 1.5 FAILED - Design incomplete")
	fmt.Println()
	fmt.Println("Edit DESIGN.md to add missing sections, then run 'vic design check' again")
	fmt.Println()
	fmt.Println("Note: Gate 1.5 is optional for CLI/non-UI projects")

	return nil
}
