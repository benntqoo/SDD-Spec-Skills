package chunker

import (
	"regexp"
	"strings"
	"time"
)

// TypeScriptChunker extracts TypeScript/JavaScript declarations using regex
type TypeScriptChunker struct{}

// Ensure TypeScriptChunker implements Chunker
var _ Chunker = (*TypeScriptChunker)(nil)

func (TypeScriptChunker) Name() string { return "TypeScript" }
func (TypeScriptChunker) Extensions() []string {
	return []string{".ts", ".tsx", ".mts", ".cts", ".js", ".mjs", ".cjs"}
}

func (t *TypeScriptChunker) ExtractChunks(filePath string, content string) []Chunk {
	var chunks []Chunk

	modulePath := extractModulePath(filePath)
	lines := strings.Split(content, "\n")
	isTS := strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx") ||
		strings.HasSuffix(filePath, ".mts") || strings.HasSuffix(filePath, ".cts")

	lang := "typescript"
	if !isTS {
		lang = "javascript"
	}

	// Match: function name(...), async function name(...)
	funcRE := regexp.MustCompile(`(?m)^\s*(?:export\s+)?(?:async\s+)?function\s+(\w+)\s*\(`)
	// Match: const name = async () => {}, const name = function(){}
	constFuncRE := regexp.MustCompile(`(?m)^\s*(?:export\s+)?const\s+(\w+)\s*=\s*(?:async\s+)?(?:\([^)]*\)|[^\s=])\s*(?:=>|function)`)
	// Match: class ClassName { or class ClassName extends ...
	classRE := regexp.MustCompile(`(?m)^\s*(?:export\s+)?(?:abstract\s+)?class\s+(\w+)(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{`)
	// Match: interface InterfaceName {
	interfaceRE := regexp.MustCompile(`(?m)^\s*(?:export\s+)?interface\s+(\w+)(?:\s+extends\s+[^{]+)?\s*\{`)
	// Match: type TypeName = { or type TypeName =
	typeRE := regexp.MustCompile(`(?m)^\s*(?:export\s+)?type\s+(\w+)\s*=`)
	processMatch := func(re *regexp.Regexp, chunkType string) {
		matches := re.FindAllStringSubmatchIndex(content, -1)
		for _, match := range matches {
			name := content[match[2]:match[3]]
			lineStart := strings.Count(content[:match[0]], "\n") + 1

			// Find approximate end (next top-level closing brace or next declaration)
			endPos := findBlockEnd(content, match[1])
			endLine := strings.Count(content[:endPos], "\n") + 1

			code := extractCodeLines(lines, lineStart, endLine)
			doc := extractJSDoc(content, match[0])

			chunks = append(chunks, Chunk{
				FilePath:   filePath,
				ChunkType:  chunkType,
				ChunkName:  name,
				ModulePath: modulePath,
				StartLine:  lineStart,
				EndLine:    endLine,
				Code:       code,
				Doc:        doc,
				Lang:       lang,
				UpdatedAt:  time.Now().Unix(),
			})
		}
	}

	processMatch(funcRE, "function")
	processMatch(constFuncRE, "const")
	processMatch(classRE, "class")
	processMatch(interfaceRE, "interface")
	processMatch(typeRE, "type")

	return chunks
}

func findBlockEnd(content string, start int) int {
	// Find the matching closing brace for the block starting at 'start'
	depth := 0
	inString := false
	stringChar := byte(0)
	escaped := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' && inString {
			escaped = true
			continue
		}

		if !inString && (c == '"' || c == '\'' || c == '`') {
			inString = true
			stringChar = c
			continue
		}

		if inString && c == stringChar {
			inString = false
			continue
		}

		if !inString {
			if c == '{' {
				depth++
			} else if c == '}' {
				depth--
				if depth == 0 {
					return i + 1
				}
			}
		}
	}
	return len(content)
}

func extractJSDoc(content string, declStart int) string {
	// Look backwards from declStart for /** ... */ style comments
	searchFrom := declStart
	if searchFrom > len(content) {
		return ""
	}

	// Find the start of the line containing decl
	lineStart := strings.LastIndex(content[:declStart], "\n") + 1
	beforeLine := strings.TrimSpace(content[max(0, lineStart-200):lineStart])

	// Check for JSDoc comment ending right before this line
	docEnd := len(beforeLine)
	if strings.HasSuffix(beforeLine, "*/") {
		docStart := strings.LastIndex(beforeLine[:docEnd-1], "/**")
		if docStart >= 0 {
			doc := beforeLine[docStart+3 : docEnd-2]
			// Clean up
			doc = strings.ReplaceAll(doc, "*\n", "\n")
			doc = strings.ReplaceAll(doc, " *", " ")
			return strings.TrimSpace(doc)
		}
	}

	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
