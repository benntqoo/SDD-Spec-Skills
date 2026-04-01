# Phase 3: Gate 检查增强 - Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 增强 Gate 检查功能，实现智能检测、详细报告和更好的代码规范验证

**Architecture:** 基于现有 gate 检查框架，添加新的检查类型、JSON 输出格式和规则引擎支持

**Tech Stack:** Go (cobra CLI), SQLite (gate 结果存储), Regex (模式匹配), YAML (规则配置)

---

## File Structure

| 文件 | 职责 | 新建/修改 |
|------|--------|----------|
| `gate_report.go` | Gate 报告生成器和 JSON 输出 | 新建 |
| `code_scanner.go` | 代码扫描器（TODO、规范检查） | 新建 |
| `gate_checker.go` | Gate 检查引擎集成 | 新建 |
| `gate0_test.go` | Gate 0 单元测试 | 新建 |
| `gate1_test.go` | Gate 1 单元测试 | 新建 |
| `gate2_test.go` | Gate 2 单元测试 | 新建 |
| `gate3_test.go` | Gate 3 单元测试 | 新建 |
| `gate_utils.go` | 通用工具函数 | 修改 |
| `gate0.go` | Gate 0 实现 | 修改 |
| `gate1.go` | Gate 1 实现 | 修改 |
| `gate2.go` | Gate 2 实现 | 修改 |
| `gate3.go` | Gate 3 实现 | 修改 |

---

### Task 1: 创建 Gate 报告结构

**Files:**
- Create: `cmd/vic-go/internal/commands/gate_report.go`

- [ ] **Step 1: Write failing test**

```go
// Test creating gate report
func TestGateReportCreation(t *testing.T) {
    report := NewGateReport(0)
    report.AddCheck("TEST", "Test Check", true, "Test passed")

    output := report.ToJSON()
    if !strings.Contains(output, `"success":true`) {
        t.Error("JSON output format incorrect")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGateReportCreation -v`
Expected: FAIL with "NewGateReport not defined"

- [ ] **Step 3: Write minimal implementation**

```go
// gate_report.go
package commands

import (
	"encoding/json"
	"fmt"
	"time"
)

// GateCheck represents a single gate check result
type GateCheck struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Passed    bool      `json:"passed"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Severity  string    `json:"severity,omitempty"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
}

// GateReport represents the full gate check report
type GateReport struct {
	GateNumber   int          `json:"gate_number"`
	GateName     string       `json:"gate_name"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     string       `json:"duration"`
	TotalChecks  int          `json:"total_checks"`
	PassedChecks int          `json:"passed_checks"`
	FailedChecks int          `json:"failed_checks"`
	Checks       []GateCheck  `json:"checks"`
	Summary     string       `json:"summary"`
	Success      bool         `json:"success"`
	Recommendations []string    `json:"recommendations,omitempty"`
}

// NewGateReport creates a new gate report
func NewGateReport(gateNum int) *GateReport {
	return &GateReport{
		GateNumber:   gateNum,
		GateName:     getGateName(gateNum),
		StartTime:    time.Now(),
		Checks:       []GateCheck{},
	}
}

// AddCheck adds a check result to the report
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGateReportCreation -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate_report.go
git commit -m "feat: add gate report structure with JSON output"
```

---

### Task 2: 创建代码扫描器

**Files:**
- Create: `cmd/vic-go/internal/commands/code_scanner.go`
- Test: `cmd/vic-go/internal/commands/code_scanner_test.go`

- [ ] **Step 1: Write failing test**

```go
func TestCodeScannerFindTODOs(t *testing.T) {
	scanner := NewCodeScanner("testdata/sample.go")
	todos := scanner.FindTODOs()

	if len(todos) == 0 {
		t.Error("Expected to find TODOs in sample file")
	}
}

func TestCodeScannerValidateConstitution(t *testing.T) {
	scanner := NewCodeScanner("testdata/sample.go")
	violations := scanner.ValidateConstitution([]string{"NO-TODO-IN-CODE"})

	if len(violations) > 0 {
		t.Error("Should find constitution violations")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestCodeScannerFindTODOs -v`
Expected: FAIL with "testdata not found"

- [ ] **Step 3: Write minimal implementation**

```go
// code_scanner.go
package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TODOComment represents a TODO comment found in code
type TODOComment struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Type     string `json:"type"`    // TODO, FIXME, XXX, etc.
	Priority  string `json:"priority,omitempty"` // low, medium, high
	Text      string `json:"text"`
	Context   string `json:"context,omitempty"`
}

// CodeViolation represents a code quality violation
type CodeViolation struct {
	Rule     string `json:"rule"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity  string `json:"severity"` // error, warning, info
	Message   string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// CodeScanner scans code files for quality issues
type CodeScanner struct {
	projectDir   string
	skipDirs     []string
	sourceExts   map[string]bool
}

// NewCodeScanner creates a new code scanner
func NewCodeScanner(projectDir string) *CodeScanner {
	return &CodeScanner{
		projectDir: projectDir,
		skipDirs: []string{"node_modules", "vendor", ".git", "dist", "build", ".venv", "venv", "__pycache__", "bin", "obj"},
		sourceExts: map[string]bool{
			".go":   true, ".py":  true, ".js":   true,
			".ts":   true, ".tsx":  true, ".jsx":  true,
			".java": true, ".rs":   true, ".cpp":  true,
		},
	}
}

// FindTODOs scans code for TODO/FIXME/XXX comments
func (s *CodeScanner) FindTODOs() []TODOComment {
	todos := []TODOComment{}

	filepath.Walk(s.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			for _, skip := range s.skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		ext := filepath.Ext(path)
		if !s.sourceExts[ext] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Match TODO patterns
			patterns := map[string]string{
				`(?i)\bTODO\b`:     "TODO",
				`(?i)\bFIXME\b`:    "FIXME",
				`(?i)\bXXX\b`:      "XXX",
				`(?i)\bHACK\b`:     "HACK",
			}

			for pattern, todoType := range patterns {
				re := regexp.MustCompile(pattern)
				if matches := re.FindStringSubmatchIndex(line, -1); matches != nil {
					todo := TODOComment{
						File:    filepath.Base(path),
						Line:    lineNum,
						Column:   matches[1],
						Type:     todoType,
						Text:     strings.TrimSpace(line[matches[1]:]),
					}

					// Determine priority
					if strings.Contains(strings.ToUpper(line), "FIXME") || strings.Contains(strings.ToUpper(line), "XXX") {
						todo.Priority = "high"
					} else if strings.Contains(strings.ToUpper(line), "HACK") {
						todo.Priority = "medium"
					} else {
						todo.Priority = "low"
					}

					todos = append(todos, todo)
				}
			}
		}

		return nil
	})

	return todos
}

// ValidateConstitution validates code against constitution rules
func (s *CodeScanner) ValidateConstitution(rules []string) []CodeViolation {
	violations := []CodeViolation{}

	filepath.Walk(s.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			for _, skip := range s.skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		ext := filepath.Ext(path)
		if !s.sourceExts[ext] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Check each rule
			for _, rule := range rules {
				if s.checkRule(rule, line, path, lineNum) {
					violation := CodeViolation{
						Rule:     rule,
						File:     filepath.Base(path),
						Line:     lineNum,
						Severity:  s.determineSeverity(rule),
						Message:   s.getRuleMessage(rule),
					}
					violations = append(violations, violation)
				}
			}
		}

		return nil
	})

	return violations
}

// checkRule checks if a line violates a specific rule
func (s *CodeScanner) checkRule(rule string, line, path string, lineNum int) bool {
	switch rule {
	case "NO-TODO-IN-CODE":
		return regexp.MustCompile(`(?i)\b(TODO|FIXME|XXX|HACK)\b`).MatchString(line)
	case "NO-CONSOLE-IN-PROD":
		return regexp.MustCompile(`(?i)\b(console\.(log|warn|error|debug)\b)`).MatchString(line)
	case "NO-HARD-CODED-SECRETS":
		return regexp.MustCompile(`(?i)(password|api[_-]?key|secret[_-]?token)\s*=\s*["']`).MatchString(line)
	}
	return false
}

// determineSeverity returns severity level for a rule
func (s *CodeScanner) determineSeverity(rule string) string {
	switch rule {
	case "NO-TODO-IN-CODE":
		return "warning"
	case "NO-CONSOLE-IN-PROD":
		return "error"
	case "NO-HARD-CODED-SECRETS":
		return "error"
	}
	return "warning"
}

// getRuleMessage returns user-friendly message for a rule
func (s *CodeScanner) getRuleMessage(rule string) string {
	messages := map[string]string{
		"NO-TODO-IN-CODE":      "Unresolved TODO/FIXME comment found",
		"NO-CONSOLE-IN-PROD":  "Console statement should not be in production code",
		"NO-HARD-CODED-SECRETS": "Hardcoded secrets detected",
	}
	return messages[rule]
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/commands/code_scanner_test.go -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/code_scanner.go cmd/vic-go/internal/commands/code_scanner_test.go
git commit -m "feat: add code scanner for TODO and constitution validation"
```

---

### Task 3: 增强 Gate 0 - 添加 JSON 输出

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate0.go`

- [ ] **Step 1: Write failing test**

```go
func TestGate0JSONOutput(t *testing.T) {
    // Test that gate0 outputs JSON format
    // Implementation will add --json flag
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGate0JSONOutput -v`
Expected: FAIL with "no JSON output"

- [ ] **Step 3: Write implementation**

```go
// Add to gate0.go
// Add format flag
var outputFormat string

// Update NewGate0Cmd
func NewGate0Cmd(cfg *config.Config) *cobra.Command {
    var outputFormat string

    cmd := &cobra.Command{
        Use:   "gate0",
        Short: "Validate SPEC-REQUIREMENTS.md structure",
        RunE: func(cmd *cobra.Command, args []string) error {
            return RunGate0(cfg, outputFormat)
        },
    }

    cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

    return cmd
}

// Update RunGate0 signature
func RunGate0(cfg *config.Config, format string) error {
    // ... existing code ...

    // After collecting results, use new report format
    if format == "json" {
        report := NewGateReport(0)
        for _, r := range results {
            report.AddCheck(r.checkID, r.checkName, r.passed, r.message)
        }
        report.Finalize(allPassed && todoCount == 0)
        fmt.Println(report.ToJSON())
    } else {
        // existing plain output
        // ...
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGate0JSONOutput -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate0.go
git commit -m "feat: add JSON output to gate0 command"
```

---

### Task 4: 增强 Gate 1 - 添加 JSON 输出

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate1.go`

- [ ] **Step 1: Write failing test**

```go
func TestGate1JSONOutput(t *testing.T) {
    // Test that gate1 outputs JSON format
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGate1JSONOutput -v`
Expected: FAIL with "no JSON output"

- [ ] **Step 3: Write implementation**

```go
// Similar to Task 3, add format flag and JSON output
var outputFormat string

// Update NewGate1Cmd
func NewGate1Cmd(cfg *config.Config) *cobra.Command {
    var outputFormat string

    cmd := &cobra.Command{
        Use:   "gate1",
        Short: "Validate SPEC-ARCHITECTURE.md structure",
        RunE: func(cmd *cobra.Command, args []string) error {
            return RunGate1(cfg, outputFormat)
        },
    }

    cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

    return cmd
}

// Update RunGate1 signature
func RunGate1(cfg *config.Config, format string) error {
    // ... existing code ...

    if format == "json" {
        report := NewGateReport(1)
        for _, r := range results {
            report.AddCheck(r.checkID, r.checkName, r.passed, r.message)
        }
        report.Finalize(allPassed && todoCount == 0)
        fmt.Println(report.ToJSON())
    } else {
        // existing plain output
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGate1JSONOutput -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate1.go
git commit -m "feat: add JSON output to gate1 command"
```

---

### Task 5: 增强 Gate 2 - 添加代码扫描集成

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate2.go`

- [ ] **Step 1: Write failing test**

```go
func TestGate2CodeScanning(t *testing.T) {
    // Test that gate2 uses code scanner
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGate2CodeScanning -v`
Expected: FAIL with "no code scanning"

- [ ] **Step 3: Write implementation**

```go
// Update RunGate2 to integrate code scanner
func RunGate2(cfg *config.Config, format string) error {
    // ... existing checks ...

    // Add new check: Code quality violations
    scanner := NewCodeScanner(cfg.ProjectDir)

    // Check for TODOs in code
    todos := scanner.FindTODOs()
    if len(todos) > 0 {
        todoResult := gate2Result{
            checkID:   "CODE_TODOS",
            checkName: "Code TODO Comments",
            passed:    false,
            message:   fmt.Sprintf("Found %d TODO/FIXME/XXX comments in code", len(todos)),
            details:   fmt.Sprintf("Most critical: %s", getMostCriticalTODO(todos)),
        }
        results = append(results, todoResult)
        allPassed = false
    }

    // Check constitution rules
    violations := scanner.ValidateConstitution([]string{
        "NO-TODO-IN-CODE",
        "NO-CONSOLE-IN-PROD",
        "NO-HARD-CODED-SECRETS",
    })

    for _, violation := range violations {
        if violation.Severity == "error" {
            allPassed = false
        }
        results = append(results, gate2Result{
            checkID:   "CONSTITUTION",
            checkName: fmt.Sprintf("Constitution: %s", violation.Rule),
            passed:    violation.Severity != "error",
            message:   violation.Message,
            details:   fmt.Sprintf("%s:%d", violation.File, violation.Line),
        })
    }

    // ... continue with existing output logic ...
}

func getMostCriticalTODO(todos []TODOComment) string {
    // Return file:line of most critical TODO
    for _, todo := range todos {
        if todo.Priority == "high" {
            return fmt.Sprintf("%s:%d", todo.File, todo.Line)
        }
    }
    if len(todos) > 0 {
        return fmt.Sprintf("%s:%d", todos[0].File, todos[0].Line)
    }
    return ""
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGate2CodeScanning -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate2.go
git commit -m "feat: add code scanning to gate2 command"
```

---

### Task 6: 增强 Gate 3 - 添加 JSON 输出

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate3.go`

- [ ] **Step 1: Write failing test**

```go
func TestGate3JSONOutput(t *testing.T) {
    // Test that gate3 outputs JSON format
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGate3JSONOutput -v`
Expected: FAIL with "no JSON output"

- [ ] **Step 3: Write implementation**

```go
// Similar to Tasks 3 & 4, add format flag and JSON output
var outputFormat string

// Update NewGate3Cmd
func NewGate3Cmd(cfg *config.Config) *cobra.Command {
    var outputFormat string

    cmd := &cobra.Command{
        Use:   "gate3",
        Short: "Validate test coverage",
        RunE: func(cmd *cobra.Command, args []string) error {
            return RunGate3(cfg, outputFormat)
        },
    }

    cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

    return cmd
}

// Update RunGate3 signature
func RunGate3(cfg *config.Config, format string) error {
    // ... existing code ...

    if format == "json" {
        report := NewGateReport(3)
        for _, r := range results {
            report.AddCheck(r.checkID, r.checkName, r.passed, r.message, r.details)
        }
        report.Finalize(passedCount >= 3)
        fmt.Println(report.ToJSON())
    } else {
        // existing plain output
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGate3JSONOutput -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate3.go
git commit -m "feat: add JSON output to gate3 command"
```

---

### Task 7: 更新 spec.go 命令添加 --format 标志

**Files:**
- Modify: `cmd/vic-go/internal/commands/spec.go`

- [ ] **Step 1: Write failing test**

```go
func TestSpecGateJSONOutput(t *testing.T) {
    // Test that spec gate accepts format flag
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestSpecGateJSONOutput -v`
Expected: FAIL with "format parameter not passed"

- [ ] **Step 3: Write implementation**

```go
// Update NewSpecGateCmd to add format parameter
func NewSpecGateCmd(cfg *config.Config) *cobra.Command {
    var gate float64
    var outputFormat string

    cmd := &cobra.Command{
        Use:   "gate [0-3|1.5]",
        Short: "Run SPEC gate check",
        Long: `Run SPEC gate validation:
  Gate 0:  Requirements Completeness
  Gate 1:  Architecture Completeness
  Gate 1.5: Design Completeness (optional, for UI projects)
  Gate 2:  Code Alignment
  Gate 3:  Test Coverage`,
        Example: `  vic spec gate 0
  vic spec gate 1
  vic spec gate 1.5
  vic spec gate 2`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runSpecGate(cfg, gate, outputFormat)
        },
    }

    cmd.Flags().Float64VarP(&gate, "gate", "g", 0, "Gate number (0-3, or 1.5)")
    cmd.Flags().StringVarP(&outputFormat, "format", "f", "plain", "Output format (plain, json)")

    return cmd
}

// Update runSpecGate signature
func runSpecGate(cfg *config.Config, gate float64, format string) error {
    // ... existing validation logic ...

    // Pass format to individual gate functions
    switch gate {
    case 0:
        return RunGate0(cfg, format)
    case 1:
        return RunGate1(cfg, format)
    case 1.5:
        return RunDesignGate(cfg)
    case 2:
        return RunGate2(cfg, format)
    case 3:
        return RunGate3(cfg, format)
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestSpecGateJSONOutput -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/spec.go
git commit -m "feat: add --format flag to vic spec gate command"
```

---

### Task 8: 添加完整测试覆盖

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate_utils.go`
- Create: `cmd/vic-go/internal/commands/gate0_test.go`
- Create: `cmd/vic-go/internal/commands/gate1_test.go`
- Create: `cmd/vic-go/internal/commands/gate2_test.go`
- Create: `cmd/vic-go/internal/commands/gate3_test.go`

- [ ] **Step 1: Write failing test**

```go
// gate0_test.go
package commands

import (
	"testing"
)

func TestRunGate0AllPass(t *testing.T) {
    // Test gate0 with all checks passing
}

func TestRunGate0MissingFile(t *testing.T) {
    // Test gate0 with missing SPEC-REQUIREMENTS.md
}

// gate1_test.go
package commands

import (
	"testing"
)

func TestRunGate1AllPass(t *testing.T) {
    // Test gate1 with all checks passing
}

func TestRunGate1MissingFile(t *testing.T) {
    // Test gate1 with missing SPEC-ARCHITECTURE.md
}

// gate2_test.go
package commands

import (
	"testing"
)

func TestRunGate2TechMismatch(t *testing.T) {
    // Test gate2 with tech stack mismatch
}

func TestRunGate2MissingAPI(t *testing.T) {
    // Test gate2 with missing API implementation
}

// gate3_test.go
package commands

import (
	"testing"
)

func TestRunGate3NoTests(t *testing.T) {
    // Test gate3 with no test files
}

func TestRunGate3PartialCoverage(t *testing.T) {
    // Test gate3 with partial test coverage
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestRunGate0AllPass -v`
Expected: FAIL with "no test file"

- [ ] **Step 3: Write minimal implementation**

```go
// Implement test files with proper setup and assertions
// Each test should test a specific scenario
```

- [ ] **Step 4: Run test to verify they pass**

Run: `go test ./internal/commands -run TestGate0AllPass -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate0_test.go cmd/vic-go/internal/commands/gate1_test.go cmd/vic-go/internal/commands/gate2_test.go cmd/vic-go/internal/commands/gate3_test.go
git commit -m "test: add comprehensive gate tests"
```

---

### Task 9: 集成到 gate.go 命令

**Files:**
- Modify: `cmd/vic-go/internal/commands/gate.go`

- [ ] **Step 1: Write failing test**

```go
func TestGateSmartCheckJSON(t *testing.T) {
    // Test gate smart check with JSON output
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/commands -run TestGateSmartCheckJSON -v`
Expected: FAIL with "no JSON output"

- [ ] **Step 3: Write implementation**

```go
// Update NewGateCheckCmd to add format flag
func NewGateCheckCmd(cfg *config.Config) *cobra.Command {
    var phaseNum int
    var blocking bool
    var format string

    cmd := &cobra.Command{
        Use:   "check",
        Short: "Check gate status for pre-commit",
        Long: `Check gate status and optionally block if gates are not passed.

This command is designed to be used in pre-commit hooks to ensure
all VIBE-SDD gates are passed before allowing a commit.

Examples:
  vic gate check                    # Check current phase gates
  vic gate check --phase 1          # Check phase 1 gates
  vic gate check --blocking         # Exit with error if gates not passed
  vic gate check --format json      # JSON output`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runGateStatusCheck(cfg, phaseNum, format)
        },
    }

    cmd.Flags().IntVarP(&phaseNum, "phase", "p", -1, "Phase number (0-3), -1 for current")
    cmd.Flags().BoolVarP(&blocking, "blocking", "b", false, "Exit with error if gates not passed (for pre-commit)")
    cmd.Flags().StringVarP(&format, "format", "f", "plain", "Output format (plain, json)")

    return cmd
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/commands -run TestGateSmartCheckJSON -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/vic-go/internal/commands/gate.go
git commit -m "feat: add --format flag to gate check command"
```

---

### Task 10: 构建验证和文档更新

**Files:**
- Modify: `D:\Code\aaa\docs\TODO.md`
- Modify: `D:\Code\aaa\docs\FIX-LOG.md`

- [ ] **Step 1: Write failing test**

```go
func TestBuildSuccess(t *testing.T) {
    // Test that build succeeds
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go build -o vic.exe . 2>&1`
Expected: SUCCESS (no errors)

- [ ] **Step 3: Run all gate tests**

Run: `go test ./internal/commands -v`
Expected: All tests pass

- [ ] **Step 4: Update TODO.md**

```markdown
### 3.2 添加智能检测 ✅ 完成
- [x] 检测 TODO/FIXME 注释 - 代码扫描器实现 ✅
- [x] 检测代码与规范不一致 - Constitution 验证实现 ✅
- [x] 检测缺失测试 - Gate 3 检查增强 ✅
- [x] 生成详细报告 - JSON 输出格式实现 ✅
```

- [ ] **Step 5: Update FIX-LOG.md**

```markdown
## 2026-03-31 - Phase 3: Gate 检查增强

### 已完成

| 文件 | 修复内容 | 状态 |
|------|----------|------|
| `gate_report.go` | 创建 Gate 报告生成器 | ✅ 完成 |
| `code_scanner.go` | 创建代码扫描器 | ✅ 完成 |
| `gate0.go` | 添加 JSON 输出 | ✅ 完成 |
| `gate1.go` | 添加 JSON 输出 | ✅ 完成 |
| `gate2.go` | 添加代码扫描集成 | ✅ 完成 |
| `gate3.go` | 添加 JSON 输出 | ✅ 完成 |
| `gate.go` | 添加 --format 标志 | ✅ 完成 |
| `*_test.go` | 添加完整测试覆盖 | ✅ 完成 |

### 新增功能

1. **JSON 输出格式**
   - 所有 gate 命令支持 `--format json`
   - 结构化的报告数据
   - 包含检查详情和建议

2. **代码扫描器**
   - 检测代码中的 TODO/FIXME/XXX 注释
   - Constitution 规则验证（NO-TODO-IN-CODE, NO-CONSOLE-IN-PROD 等）
   - 按严重级别分类违规

3. **测试覆盖**
   - 完整的单元测试
   - 测试各种场景

### 验证

```
✅ go build -o vic.exe .              # 构建成功
✅ go test ./internal/commands -v    # 所有测试通过
✅ vic spec gate 0 --format json   # JSON 输出正常
✅ vic spec gate 1 --format json   # JSON 输出正常
✅ vic spec gate 2 --format json   # JSON 输出正常
✅ vic spec gate 3 --format json   # JSON 输出正常
```
```

- [ ] **Step 6: Final commit**

```bash
git add docs/TODO.md docs/FIX-LOG.md
git commit -m "docs: update Phase 3 completion status"
```

---

## Summary

Phase 3 将实现以下增强：

1. **Gate 报告结构** - 新增统一的报告格式支持 JSON 输出
2. **代码扫描器** - 实现 TODO/FIXME 检测和 Constitution 规则验证
3. **Gate 0/1/2/3 增强** - 添加 JSON 输出支持
4. **Gate 2 智能检测** - 集成代码扫描器检测代码质量问题
5. **测试覆盖** - 为所有 gate 命令添加完整单元测试
6. **命令集成** - 更新 gate check 和 spec gate 命令支持 --format 标志

## Notes

- 每个 gate 检查的 JSON 输出应包含：gate_number, gate_name, total_checks, passed_checks, failed_checks, checks[], summary, success
- 代码扫描器应支持多种编程语言（Go, Python, JS/TS）
- Constitution 规则应可配置（从 .vic-sdd/constitution.yaml 读取）
- 所有新功能都应有对应的单元测试
