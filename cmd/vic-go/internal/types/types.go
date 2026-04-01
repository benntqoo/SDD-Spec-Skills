package types

import (
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Custom time formats to handle various YAML timestamp formats
var timeFormats = []string{
	time.RFC3339,                // 2006-01-02T15:04:05Z07:00
	"2006-01-02 15:04:05+00:00", // 2026-03-17 10:00:00+00:00
	"2006-01-02 15:04:05Z",      // 2026-03-17 10:00:00Z
	"2006-01-02 15:04:05",       // 2026-03-17 10:00:00
	"2006-01-02",                // 2026-01-02
}

// ParseTime parses time from various formats
func ParseTime(s string) time.Time {
	for _, fmt := range timeFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t
		}
	}
	// Try to find and replace space with T
	s = strings.Replace(s, " ", "T", 1)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

// ============================================
// Event Types
// ============================================

// Event represents an event in the system
type Event struct {
	ID        string     `yaml:"id"`
	Type      string     `yaml:"type"` // decision_made, risk_identified, dependency_recorded
	Timestamp CustomTime `yaml:"timestamp"`
	Agent     string     `yaml:"agent"` // AI or human
	Data      EventData  `yaml:"data"`
}

// CustomTime handles multiple timestamp formats
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	ct.Time = ParseTime(s)
	return nil
}

// EventData contains the actual event data
type EventData map[string]interface{}

// TechRecord represents a technical decision
type TechRecord struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Decision  string   `yaml:"decision"`
	Category  string   `yaml:"category"`
	Reason    string   `yaml:"reason"`
	Impact    string   `yaml:"impact"` // low, medium, high
	Status    string   `yaml:"status"` // planned, in_progress, completed, deprecated
	Files     []string `yaml:"files"`
	CreatedAt string   `yaml:"created_at"`
	UpdatedAt string   `yaml:"updated_at"`
}

// RiskRecord represents a risk
type RiskRecord struct {
	ID          string `yaml:"id"`
	Area        string `yaml:"area"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Impact      string `yaml:"impact"` // low, medium, high, critical
	Status      string `yaml:"status"` // identified, mitigating, resolved, accepted
	CreatedAt   string `yaml:"created_at"`
	UpdatedAt   string `yaml:"updated_at"`
}

// DependencyRecord represents a module dependency
type DependencyRecord struct {
	Module    string   `yaml:"module"`
	DependsOn []string `yaml:"depends_on"`
	CreatedAt string   `yaml:"created_at"`
}

// ============================================
// State Types
// ============================================

// State represents the current state
type State struct {
	LastFolded      string             `yaml:"last_folded"`
	ActiveDecisions int                `yaml:"active_decisions"`
	ActiveRisks     int                `yaml:"active_risks"`
	TechRecords     []TechRecord       `yaml:"tech_records"`
	Risks           []RiskRecord       `yaml:"risks"`
	Dependencies    []DependencyRecord `yaml:"dependencies"`
}

// TechRecordsFile represents the tech-records.yaml structure
type TechRecordsFile struct {
	Version     string       `yaml:"version,omitempty"`
	Records     []TechRecord `yaml:"records"`      // matches Python vic format
	TechRecords []TechRecord `yaml:"tech_records"` // alternative format
	Metadata    interface{}  `yaml:"metadata,omitempty"`
}

// RiskZonesFile represents the risk-zones.yaml structure
type RiskZonesFile struct {
	Risks []RiskRecord `yaml:"risks"`
}

// EventsFile represents the events.yaml structure
type EventsFile struct {
	Events []Event `yaml:"events"`
}

// StateFile represents the state.yaml structure
type StateFile struct {
	State State `yaml:"state"`
}

// ============================================
// Phase Types (SDD Flow)
// ============================================

// Phase represents a development phase
type Phase struct {
	Name            string          `yaml:"name"`
	Status          string          `yaml:"status"` // pending, in_progress, completed
	StartedAt       string          `yaml:"started_at"`
	CompletedAt     string          `yaml:"completed_at"`
	Completion      int             `yaml:"completion"` // 0-100
	OutputsRequired []string        `yaml:"outputs_required"`
	Gates           map[string]Gate `yaml:"gates"`
}

// Gate represents a gate check
type Gate struct {
	Name        string `yaml:"name"`
	Status      string `yaml:"status"` // pending, passed, failed
	CheckedAt   string `yaml:"checked_at"`
	CheckedBy   string `yaml:"checked_by"`
	Notes       string `yaml:"notes"`
	Description string `yaml:"description"`
	Phase       int    `yaml:"phase"`
}

// PhaseFile represents the phase status file
type PhaseFile struct {
	CycleID      string        `yaml:"cycle_id"`
	CycleName    string        `yaml:"cycle_name"`
	CurrentPhase int           `yaml:"current_phase"`
	CurrentGate  int           `yaml:"current_gate"`
	StartedAt    string        `yaml:"started_at"`
	LastUpdated  string        `yaml:"last_updated"`
	Phases       map[int]Phase `yaml:"phases"`
}

// ============================================
// Product Types
// ============================================

// ProductRecord represents a product redesign record
type ProductRecord struct {
	ID        string `yaml:"id"`
	Original  string `yaml:"original"`
	Real      string `yaml:"real"`
	Mode      string `yaml:"mode"`
	Timestamp string `yaml:"timestamp"`
}
