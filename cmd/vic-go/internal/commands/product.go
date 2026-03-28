package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
	"gopkg.in/yaml.v3"
)

// NewProductCmd creates the product command
func NewProductCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Product redesign",
		Long:  `Track product redesign decisions and the four modes.`,
		Example: `  vic product record --original "Photo upload" --real "Smart listings"
  vic product list          # List product decisions
  vic product modes        # Show the four modes`,
	}

	cmd.AddCommand(NewProductRecordCmd(cfg))
	cmd.AddCommand(NewProductListCmd(cfg))
	cmd.AddCommand(NewProductModesCmd(cfg))

	return cmd
}

// ============================================
// Record Command
// ============================================

// NewProductRecordCmd creates the product record subcommand
func NewProductRecordCmd(cfg *config.Config) *cobra.Command {
	var original, real, mode string

	cmd := &cobra.Command{
		Use:     "record",
		Short:   "Record product redesign",
		Long:    `Record a product redesign decision from original request to real product.`,
		Example: `  vic product record --original "Photo upload" --real "Smart listings" --mode expansion`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProductRecord(cfg, original, real, mode)
		},
	}

	cmd.Flags().StringVar(&original, "original", "", "Original request (required)")
	cmd.Flags().StringVar(&real, "real", "", "Real product purpose (required)")
	cmd.Flags().StringVar(&mode, "mode", "selective", "Mode (expansion/selective/hold/reduction)")

	return cmd
}

func runProductRecord(cfg *config.Config, original, real, mode string) error {
	if original == "" || real == "" {
		fmt.Println("📋 Record Product Redesign")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  vic product record \\")
		fmt.Println("    --original \"User said: photo upload\" \\")
		fmt.Println("    --real \"Help sellers create sellable listings\" \\")
		fmt.Println("    --mode expansion")
		fmt.Println("")
		fmt.Println("Modes:")
		fmt.Println("  expansion   - Dream big, propose ambitious version")
		fmt.Println("  selective   - Neutral presentation, let user choose")
		fmt.Println("  hold        - Maximum rigor, no expansions")
		fmt.Println("  reduction   - Find minimum viable version")
		return nil
	}

	// Load existing records
	records, _ := loadProductRecords(cfg)

	// Create new record
	recordID := fmt.Sprintf("PROD-%03d", len(records.Records)+1)
	record := types.ProductRedesign{
		ID:           recordID,
		Type:         "product-redesign",
		Timestamp:    time.Now(),
		OriginalReq:  original,
		RealProduct:  real,
		Mode:         mode,
		Decisions:    []types.ProductDecision{},
		UserApproved: false,
	}

	records.Records = append(records.Records, record)

	// Save
	if err := saveProductRecords(cfg, records); err != nil {
		return fmt.Errorf("failed to save product records: %w", err)
	}

	fmt.Println("✅ Product redesign recorded")
	fmt.Printf("   ID: %s\n", recordID)
	fmt.Println("")
	fmt.Printf("   Original: %s\n", original)
	fmt.Printf("   Real Product: %s\n", real)
	fmt.Printf("   Mode: %s\n", mode)
	fmt.Println("")
	fmt.Println("   Next: 'vic product list' to view all decisions")

	return nil
}

// ============================================
// List Command
// ============================================

// NewProductListCmd creates the product list subcommand
func NewProductListCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List product decisions",
		Long:    `Show all product redesign decisions.`,
		Example: `  vic product list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProductList(cfg)
		},
	}
}

func runProductList(cfg *config.Config) error {
	records, err := loadProductRecords(cfg)
	if err != nil || len(records.Records) == 0 {
		fmt.Println("📋 No product redesign decisions yet")
		fmt.Println("")
		fmt.Println("Use 'vic product record' to record a redesign decision")
		return nil
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📋 Product Redesign Decisions")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	for i := len(records.Records) - 1; i >= 0; i-- {
		record := records.Records[i]

		modeIcon := "🎯"
		switch record.Mode {
		case "expansion":
			modeIcon = "🚀"
		case "selective":
			modeIcon = "⚖️"
		case "hold":
			modeIcon = "🔒"
		case "reduction":
			modeIcon = "✂️"
		}

		fmt.Printf("   %s %s [%s]\n", modeIcon, record.ID, record.Mode)
		fmt.Printf("      Original: %s\n", truncate(record.OriginalReq, 50))
		fmt.Printf("      Real:    %s\n", truncate(record.RealProduct, 50))
		fmt.Printf("      Time: %s\n", record.Timestamp.Format("2006-01-02"))
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("   Total: %d decisions\n", len(records.Records))

	return nil
}

// ============================================
// Modes Command
// ============================================

// NewProductModesCmd creates the product modes subcommand
func NewProductModesCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "modes",
		Short:   "Show the four modes",
		Long:    `Display the four product redesign modes with examples.`,
		Example: `  vic product modes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProductModes()
		},
	}
}

func runProductModes() error {
	fmt.Println(`
═══════════════════════════════════════════════════════════
  🚀 Product Redesign - The Four Modes
═══════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────┐
│ 1. EXPANSION - Scope Expansion                          │
├─────────────────────────────────────────────────────────┤
│ Dream big. Propose the ambitious version.                │
│                                                         │
│ Example:                                                │
│ "User wants photo upload"                               │
│ → "What if we auto-detect product from photo?"          │
│ → "What if we draft descriptions automatically?"         │
│ → "What if we suggest the best hero image?"            │
│                                                         │
│ Icon: 🚀 | When: Explore possibilities                    │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ 2. SELECTIVE - Selective Expansion                     │
├─────────────────────────────────────────────────────────┤
│ Present opportunities neutrally. Let user choose.       │
│                                                         │
│ Example:                                                │
│ ┌────────────────────────────────────┐                 │
│ │ Opportunity A                      │                 │
│ │ Effort: Medium | Impact: High      │                 │
│ │ Why: Saves users 10 mins           │                 │
│ └────────────────────────────────────┘                 │
│                                                         │
│ Icon: ⚖️ | When: Evaluate trade-offs                    │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ 3. HOLD - Hold Scope                                   │
├─────────────────────────────────────────────────────────┤
│ Maximum rigor. No expansions. Stay focused.              │
│                                                         │
│ "Let's lock down the current scope and ship."           │
│                                                         │
│ Icon: 🔒 | When: Stay focused, time constrained          │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ 4. REDUCTION - Scope Reduction                         │
├─────────────────────────────────────────────────────────┤
│ Find the minimum viable version. Cut everything else.  │
│                                                         │
│ "What's the ONE thing this must do?"                   │
│ "What can we cut and add later?"                        │
│                                                         │
│ Icon: ✂️ | When: Ship fast                             │
└─────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════
  Usage: vic product record --original "..." --real "..." --mode [mode]
═══════════════════════════════════════════════════════════
`)

	return nil
}

// ============================================
// Helper Functions
// ============================================

func loadProductRecords(cfg *config.Config) (*types.ProductRecordsFile, error) {
	productFile := filepath.Join(cfg.GetVICDir(), "status", "product-records.yaml")

	if !utils.FileExists(productFile) {
		return &types.ProductRecordsFile{
			Version: "1.0",
			Records: []types.ProductRedesign{},
		}, nil
	}

	data, err := os.ReadFile(productFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read product records: %w", err)
	}

	var records types.ProductRecordsFile
	if err := yaml.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse product records: %w", err)
	}

	return &records, nil
}

func saveProductRecords(cfg *config.Config, records *types.ProductRecordsFile) error {
	productFile := filepath.Join(cfg.GetVICDir(), "status", "product-records.yaml")

	// Ensure directory exists
	statusDir := filepath.Join(cfg.GetVICDir(), "status")
	if !utils.FileExists(statusDir) {
		if err := os.MkdirAll(statusDir, 0755); err != nil {
			return fmt.Errorf("failed to create status directory: %w", err)
        }
	}

	data, err := yaml.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal product records: %w", err)
	}

	header := []byte(`# Product Redesign Records - VIBE-SDD
# Auto-generated by vic-go
# Do not edit manually - use vic product commands

`)
	if err := os.WriteFile(productFile, append(header, data...), 0644); err != nil {
		return fmt.Errorf("failed to write product records: %w", err)
	}

	return nil
}
