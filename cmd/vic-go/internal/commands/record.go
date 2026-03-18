package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"github.com/vic-sdd/vic/internal/utils"
)

// NewRecordCmd creates the record command
func NewRecordCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record decisions, risks, or dependencies",
		Long:  `Record technical decisions, risks, or dependencies to the project memory.`,
	}

	cmd.AddCommand(NewRecordTechCmd(cfg))
	cmd.AddCommand(NewRecordRiskCmd(cfg))
	cmd.AddCommand(NewRecordDepCmd(cfg))

	return cmd
}

// NewRecordTechCmd records a technical decision
func NewRecordTechCmd(cfg *config.Config) *cobra.Command {
	var id, title, decision, category, reason, impact, status, files string

	cmd := &cobra.Command{
		Use:   "tech",
		Short: "Record a technical decision",
		Example: `  vic record tech --id DB-001 --title "Use PostgreSQL" --decision "Primary DB" --reason "ACID compliance"
  vic rt --id AUTH-001 --title "JWT Auth" --decision "Use JWT" --status in_progress`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordTech(cfg, id, title, decision, category, reason, impact, status, files)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "", "", "Decision ID (required)")
	cmd.Flags().StringVarP(&title, "title", "", "", "Decision title (required)")
	cmd.Flags().StringVarP(&decision, "decision", "", "", "The decision (required)")
	cmd.Flags().StringVarP(&category, "category", "", "general", "Category")
	cmd.Flags().StringVarP(&reason, "reason", "", "", "Reason for decision")
	cmd.Flags().StringVarP(&impact, "impact", "", "medium", "Impact level (low/medium/high)")
	cmd.Flags().StringVarP(&status, "status", "", "planned", "Status (planned/in_progress/completed/deprecated)")
	cmd.Flags().StringVarP(&files, "files", "", "", "Related files (comma-separated)")

	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("decision")

	return cmd
}

func runRecordTech(cfg *config.Config, id, title, decision, category, reason, impact, status, files string) error {
	now := time.Now().Format("2006-01-02")

	// Parse files
	var fileList []string
	if files != "" {
		fileList = strings.Split(files, ",")
		for i := range fileList {
			fileList[i] = strings.TrimSpace(fileList[i])
		}
	}

	// Create tech record
	record := types.TechRecord{
		ID:        id,
		Title:     title,
		Decision:  decision,
		Category:  category,
		Reason:    reason,
		Impact:    impact,
		Status:    status,
		Files:     fileList,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Load existing records
	records, err := utils.LoadTechRecords(cfg)
	if err != nil {
		records = types.TechRecordsFile{}
	}

	// Check for duplicate ID
	updated := false
	for i, r := range records.TechRecords {
		if r.ID == id {
			records.TechRecords[i] = record
			records.TechRecords[i].UpdatedAt = now
			updated = true
			break
		}
	}

	if !updated {
		records.TechRecords = append(records.TechRecords, record)
	}

	// Save
	if err := utils.SaveTechRecords(cfg, records); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// Also add event
	event := types.Event{
		ID:        uuid.New().String(),
		Type:      "decision_made",
		Timestamp: types.CustomTime{Time: time.Now()},
		Agent:     "cli",
		Data:      types.EventData{"id": id, "title": title},
	}

	events, _ := utils.LoadEvents(cfg)
	events.Events = append(events.Events, event)
	utils.SaveEvents(cfg, events)

	action := "Recorded"
	if updated {
		action = "Updated"
	}
	fmt.Printf("✅ %s technical decision: %s\n", action, id)
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("   Decision: %s\n", decision)

	return nil
}

// NewRecordRiskCmd records a risk
func NewRecordRiskCmd(cfg *config.Config) *cobra.Command {
	var id, area, desc, category, impact, status string

	cmd := &cobra.Command{
		Use:     "risk",
		Short:   "Record a risk",
		Example: `  vic record risk --id RISK-001 --area auth --desc "JWT not validated" --impact critical`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordRisk(cfg, id, area, desc, category, impact, status)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "", "", "Risk ID (required)")
	cmd.Flags().StringVarP(&area, "area", "", "", "Risk area (required)")
	cmd.Flags().StringVarP(&desc, "desc", "", "", "Risk description (required)")
	cmd.Flags().StringVarP(&category, "category", "", "", "Category")
	cmd.Flags().StringVarP(&impact, "impact", "", "medium", "Impact (low/medium/high/critical)")
	cmd.Flags().StringVarP(&status, "status", "", "identified", "Status")

	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("area")
	cmd.MarkFlagRequired("desc")

	return cmd
}

func runRecordRisk(cfg *config.Config, id, area, desc, category, impact, status string) error {
	now := time.Now().Format("2006-01-02")

	risk := types.RiskRecord{
		ID:          id,
		Area:        area,
		Description: desc,
		Category:    category,
		Impact:      impact,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Load existing risks
	risks, err := utils.LoadRiskZones(cfg)
	if err != nil {
		risks = types.RiskZonesFile{}
	}

	// Check for duplicate
	updated := false
	for i, r := range risks.Risks {
		if r.ID == id {
			risks.Risks[i] = risk
			risks.Risks[i].UpdatedAt = now
			updated = true
			break
		}
	}

	if !updated {
		risks.Risks = append(risks.Risks, risk)
	}

	if err := utils.SaveRiskZones(cfg, risks); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// Add event
	event := types.Event{
		ID:        uuid.New().String(),
		Type:      "risk_identified",
		Timestamp: types.CustomTime{Time: time.Now()},
		Agent:     "cli",
		Data:      types.EventData{"id": id, "area": area},
	}

	events, _ := utils.LoadEvents(cfg)
	events.Events = append(events.Events, event)
	utils.SaveEvents(cfg, events)

	action := "Recorded"
	if updated {
		action = "Updated"
	}
	fmt.Printf("✅ %s risk: %s\n", action, id)
	fmt.Printf("   Area: %s\n", area)
	fmt.Printf("   Description: %s\n", desc)

	return nil
}

// NewRecordDepCmd records a dependency
func NewRecordDepCmd(cfg *config.Config) *cobra.Command {
	var module, deps string

	cmd := &cobra.Command{
		Use:     "dep",
		Short:   "Record module dependencies",
		Example: `  vic record dep --module auth-service --deps user-service,jwt-service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordDep(cfg, module, deps)
		},
	}

	cmd.Flags().StringVarP(&module, "module", "", "", "Module name (required)")
	cmd.Flags().StringVarP(&deps, "deps", "", "", "Dependencies (comma-separated)")

	cmd.MarkFlagRequired("module")
	cmd.MarkFlagRequired("deps")

	return cmd
}

func runRecordDep(cfg *config.Config, module, deps string) error {
	now := time.Now().Format("2006-01-02")

	depList := strings.Split(deps, ",")
	for i := range depList {
		depList[i] = strings.TrimSpace(depList[i])
	}

	dep := types.DependencyRecord{
		Module:    module,
		DependsOn: depList,
		CreatedAt: now,
	}

	// Load state
	stateFile, err := utils.LoadState(cfg)
	if err != nil {
		stateFile = types.StateFile{State: types.State{}}
	}

	// Check for existing module
	updated := false
	for i, d := range stateFile.State.Dependencies {
		if d.Module == module {
			stateFile.State.Dependencies[i] = dep
			updated = true
			break
		}
	}

	if !updated {
		stateFile.State.Dependencies = append(stateFile.State.Dependencies, dep)
	}

	if err := utils.SaveState(cfg, stateFile); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// Add event
	event := types.Event{
		ID:        uuid.New().String(),
		Type:      "dependency_recorded",
		Timestamp: types.CustomTime{Time: time.Now()},
		Agent:     "cli",
		Data:      types.EventData{"module": module},
	}

	events, _ := utils.LoadEvents(cfg)
	events.Events = append(events.Events, event)
	utils.SaveEvents(cfg, events)

	action := "Recorded"
	if updated {
		action = "Updated"
	}
	fmt.Printf("✅ %s dependency: %s\n", action, module)
	fmt.Printf("   Depends on: %s\n", deps)

	return nil
}
