package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GateCheck represents a single gate check result
type GateCheck struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message"`
	Details  string `json:"details,omitempty"`
	Severity string `json:"severity,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// GateReport represents a full gate check report
type GateReport struct {
	GateNumber      int         `json:"gate_number"`
	GateName        string      `json:"gate_name"`
	StartTime       time.Time   `json:"start_time"`
	EndTime         time.Time   `json:"end_time"`
	Duration        string      `json:"duration"`
	TotalChecks     int         `json:"total_checks"`
	PassedChecks    int         `json:"passed_checks"`
	FailedChecks    int         `json:"failed_checks"`
	Checks          []GateCheck `json:"checks"`
	Summary         string      `json:"summary"`
	Success         bool        `json:"success"`
	Recommendations []string    `json:"recommendations,omitempty"`
}

// NewGateReport creates a new gate report
func NewGateReport(gateNum int) *GateReport {
	return &GateReport{
		GateNumber: gateNum,
		GateName:   getGateName(gateNum),
		StartTime:  time.Now(),
		Checks:     []GateCheck{},
	}
}

// AddCheck adds a check result to report
func (r *GateReport) AddCheck(id, name string, passed bool, message string, details ...string) *GateReport {
	check := GateCheck{
		ID:      id,
		Name:    name,
		Passed:  passed,
		Message: message,
	}
	if len(details) > 0 {
		check.Details = details[0]
	}
	r.Checks = append(r.Checks, check)
	return r
}

// Finalize completes the report and calculates summary
func (r *GateReport) Finalize(success bool) {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime).String()
	r.TotalChecks = len(r.Checks)

	for _, check := range r.Checks {
		if check.Passed {
			r.PassedChecks++
		} else {
			r.FailedChecks++
		}
	}

	r.Success = success
	if success {
		r.Summary = "✅ Gate PASSED"
	} else {
		r.Summary = "❌ Gate FAILED"
	}
}

// ToJSON converts report to JSON string
func (r *GateReport) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// ToPlain converts report to plain text output
func (r *GateReport) ToPlain() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🔍 Gate %d: %s\n", r.GateNumber, r.GateName))
	sb.WriteString("========================================\n\n")

	for _, check := range r.Checks {
		icon := "❌"
		if check.Passed {
			icon = "✅"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", icon, check.Name))
		if check.Passed {
			sb.WriteString(fmt.Sprintf("      %s\n", check.Message))
		} else {
			sb.WriteString(fmt.Sprintf("      → %s\n", check.Message))
		}
		if check.Details != "" {
			sb.WriteString(fmt.Sprintf("         Details: %s\n", check.Details))
		}
	}

	sb.WriteString("\n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("📊 Summary: %s\n", r.Summary))
	sb.WriteString(fmt.Sprintf("   Passed: %d/%d\n", r.PassedChecks, r.TotalChecks))
	sb.WriteString(fmt.Sprintf("   Duration: %s\n", r.Duration))

	return sb.String()
}

func getGateName(gateNum int) string {
	names := map[int]string{
		0: "Requirements Completeness",
		1: "Architecture Completeness",
		2: "Code Alignment",
		3: "Test Coverage",
	}
	if name, ok := names[gateNum]; ok {
		return name
	}
	return fmt.Sprintf("Gate %d", gateNum)
}
