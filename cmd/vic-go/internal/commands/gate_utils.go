package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// checkForTBDs checks for unresolved TBD/TODO/FIXME markers
func checkForTBDs(content string) []gate0Result {
	results := make([]gate0Result, 0)
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	todosFound := 0

	todoPatterns := []string{
		`(?i)\bTBD\b`,
		`(?i)\bTODO\b`,
		`(?i)\bFIXME\b`,
		`(?i)\bXXX\b`,
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, pattern := range todoPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(line) {
				todosFound++
				results = append(results, gate0Result{
					checkID:    "TODO",
					checkName:  "No Unresolved TODOs",
					passed:     false,
					message:    fmt.Sprintf("Line %d: %s", lineNum, strings.TrimSpace(line)),
					lineNumber: lineNum,
				})
				break
			}
		}
	}

	if todosFound == 0 {
		results = append(results, gate0Result{
			checkID:   "TODO",
			checkName: "No Unresolved TODOs",
			passed:    true,
			message:   "No TBD/TODO/FIXME markers found",
		})
	}

	return results
}

// fileExists is a simple check for file existence
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
