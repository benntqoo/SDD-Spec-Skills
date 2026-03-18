package utils

import (
	"fmt"
	"os"

	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/types"
	"gopkg.in/yaml.v3"
)

// ============================================
// File Operations
// ============================================

// LoadYAML loads a YAML file into a struct
func LoadYAML[T any](path string) (T, error) {
	var result T
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}

// SaveYAML saves a struct to a YAML file
func SaveYAML(path string, data interface{}) error {
	// Ensure directory exists
	dir := fmt.Sprintf("%s/", path[:len(path)-len(path[len(path)-len(FileName(path))-1:])-1])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal with indentation
	d, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(path, d, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// FileName extracts filename from path
func FileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}

// LoadEvents loads events from file
func LoadEvents(cfg *config.Config) (types.EventsFile, error) {
	return LoadYAML[types.EventsFile](cfg.EventsFile)
}

// SaveEvents saves events to file
func SaveEvents(cfg *config.Config, events types.EventsFile) error {
	return SaveYAML(cfg.EventsFile, events)
}

// LoadTechRecords loads tech records from file
func LoadTechRecords(cfg *config.Config) (types.TechRecordsFile, error) {
	records, err := LoadYAML[types.TechRecordsFile](cfg.TechRecordsFile)
	// Normalize: if loaded from "records" key, copy to TechRecords
	if len(records.Records) > 0 && len(records.TechRecords) == 0 {
		records.TechRecords = records.Records
	}
	return records, err
}

// SaveTechRecords saves tech records to file
func SaveTechRecords(cfg *config.Config, records types.TechRecordsFile) error {
	return SaveYAML(cfg.TechRecordsFile, records)
}

// LoadRiskZones loads risk zones from file
func LoadRiskZones(cfg *config.Config) (types.RiskZonesFile, error) {
	return LoadYAML[types.RiskZonesFile](cfg.RiskZonesFile)
}

// SaveRiskZones saves risk zones to file
func SaveRiskZones(cfg *config.Config, risks types.RiskZonesFile) error {
	return SaveYAML(cfg.RiskZonesFile, risks)
}

// LoadState loads state from file
func LoadState(cfg *config.Config) (types.StateFile, error) {
	return LoadYAML[types.StateFile](cfg.StateFile)
}

// SaveState saves state to file
func SaveState(cfg *config.Config, state types.StateFile) error {
	return SaveYAML(cfg.StateFile, state)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureFile ensures a file exists with minimal content
func EnsureFile(path, content string) error {
	if FileExists(path) {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
