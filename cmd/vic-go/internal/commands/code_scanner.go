package commands

import (
	"bufio"
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
	Type     string `json:"type"`           // TODO, FIXME, XXX, HACK
	Priority string `json:"priority,omitempty"` // low, medium, high
	Text      string `json:"text"`
	Context   string `json:"context,omitempty"`
}

// CodeViolation represents a code quality violation
type CodeViolation struct {
	Rule      string `json:"rule"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Severity  string `json:"severity"` // error, warning, info
	Message   string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// CodeScanner scans code files for quality issues
type CodeScanner struct {
	projectDir string
	skipDirs   []string
	sourceExts map[string]bool
}

// NewCodeScanner creates a new code scanner
func NewCodeScanner(projectDir string) *CodeScanner {
	return &CodeScanner{
		projectDir: projectDir,
		skipDirs: []string{
			"node_modules", "vendor", ".git", "dist", "build",
			".venv", "venv", "__pycache__", "bin", "obj", ".vic-sdd",
		},
		sourceExts: map[string]bool{
			".go": true, ".py": true, ".js": true,
			".ts": true, ".tsx": true, ".jsx": true,
			".java": true, ".rs": true, ".cpp": true,
			".c": true, ".cc": true, ".h": true,
		},
	}
}

// FindTODOs scans code for TODO/FIXME/XXX/HACK comments
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

			todoPatterns := map[string]string{
				`(?i)\bTODO\b`:     "TODO",
				`(?i)\bFIXME\b`:    "FIXME",
				`(?i)\bXXX\b`:      "XXX",
				`(?i)\bHACK\b`:     "HACK",
			}

			for pattern, todoType := range todoPatterns {
				re := regexp.MustCompile(pattern)
				if matches := re.FindStringSubmatchIndex(line); matches != nil && len(matches) > 0 {
					todo := TODOComment{
						File:    filepath.Base(path),
						Line:    lineNum,
						Type:     todoType,
						Text:     strings.TrimSpace(line),
					}
					if len(matches) > 0 {
						todo.Column = matches[0]
					}

					switch todoType {
					case "FIXME", "XXX":
						todo.Priority = "high"
					case "HACK":
						todo.Priority = "medium"
					default:
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

			for _, rule := range rules {
				if s.checkRule(rule, line) {
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
func (s *CodeScanner) checkRule(rule string, line string) bool {
	switch rule {
	case "NO-TODO-IN-CODE":
		return regexp.MustCompile(`(?i)\b(TODO|FIXME|XXX|HACK)\b`).MatchString(line)
	case "NO-CONSOLE-IN-PROD":
		return regexp.MustCompile(`(?i)(console\.(log|warn|error|debug|info)\b)`).MatchString(line)
	case "NO-HARD-CODED-SECRETS":
		return regexp.MustCompile(`(?i)(password|api[_-]?key|secret[_-]?token|private[_-]?key)\s*=\s*["']`).MatchString(line)
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
		"NO-TODO-IN-CODE":      "Unresolved TODO/FIXME/XXX/HACK comment found",
		"NO-CONSOLE-IN-PROD":  "Console statement should not be in production code",
		"NO-HARD-CODED-SECRETS": "Hardcoded secrets detected",
	}
	return messages[rule]
}
