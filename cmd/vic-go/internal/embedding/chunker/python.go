package chunker

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// PythonChunker extracts Python classes and functions using regex
type PythonChunker struct{}

// Ensure PythonChunker implements Chunker
var _ Chunker = (*PythonChunker)(nil)

func (PythonChunker) Name() string         { return "Python" }
func (PythonChunker) Extensions() []string { return []string{".py", ".pyi"} }

func (p *PythonChunker) ExtractChunks(filePath string, content string) []Chunk {
	var chunks []Chunk

	modulePath := extractPythonModulePath(filePath)
	lines := strings.Split(content, "\n")

	// Match class definitions: class ClassName(bases):
	classRE := regexp.MustCompile(`(?m)^(\s*)class\s+(\w+)\s*(?:\([^)]*\))?\s*:`)

	// Match function definitions: def func_name(args):
	funcRE := regexp.MustCompile(`(?m)^(\s*)def\s+(\w+)\s*\(([^)]*)\)\s*(?:->\s*[^:]+)?\s*:`)

	// Extract classes
	classMatches := classRE.FindAllStringSubmatchIndex(content, -1)
	for _, match := range classMatches {
		chunk := p.extractPythonChunk(filePath, content, lines, match, "class", modulePath)
		if chunk.ChunkName != "" {
			chunks = append(chunks, chunk)
		}
	}

	// Extract functions (not inside classes)
	funcMatches := funcRE.FindAllStringSubmatchIndex(content, -1)
	for _, match := range funcMatches {
		indent := extractIndent(string(content[match[2]:match[3]]))
		// Only process top-level functions (not indented, i.e., inside a class)
		if indent == "" || indent == "\t" {
			// Check if this function is inside a class
			funcStart := match[0]
			isInClass := false
			for _, classMatch := range classMatches {
				classEnd := classMatch[1]
				if classEnd < funcStart {
					// Class ends before this function
					classIndent := extractIndent(string(content[classMatch[2]:classMatch[3]]))
					// Get the next line after class definition to find its end
					classBodyStart := getLineEnd(content, classMatch[1]) + 1
					// Find the next line at same or lower indentation
					rest := content[classBodyStart:]
					classEnd = findClassEnd(rest, len(classIndent))

					if classEnd >= 0 && funcStart < classBodyStart+classEnd {
						isInClass = true
						break
					}
				}
			}

			if !isInClass {
				chunk := p.extractPythonChunk(filePath, content, lines, match, "func", modulePath)
				if chunk.ChunkName != "" {
					chunks = append(chunks, chunk)
				}
			}
		}
	}

	return chunks
}

func (p *PythonChunker) extractPythonChunk(filePath, content string, lines []string, match []int, chunkType, modulePath string) Chunk {
	// Extract name (group 2 for class, group 3 for func)
	nameGroup := 2
	if chunkType == "func" {
		nameGroup = 3
	}
	name := content[match[nameGroup*2]:match[nameGroup*2+1]]

	// Find line number
	lineStart := getLineStart(content, match[0])
	lineNum := strings.Count(content[:lineStart], "\n") + 1

	// Find the end of the definition line
	lineEnd := getLineEnd(content, match[1])

	// Find the matching indentation level to determine the end of the block
	indent := extractIndent(string(content[match[2]:match[3]]))
	indentLen := len(indent)

	// Find the end of this block by finding the next line at same or lower indentation
	blockEnd := lineEnd
	searchFrom := match[1]

	for {
		nextLineStart := strings.Index(content[searchFrom:], "\n")
		if nextLineStart < 0 {
			break
		}
		searchFrom += nextLineStart + 1

		if searchFrom >= len(content) {
			break
		}

		// Get the line content (skip empty lines)
		lineContent := ""
		if searchFrom < len(content) {
			lineContent = content[searchFrom:]
			if eol := strings.Index(lineContent, "\n"); eol >= 0 {
				lineContent = lineContent[:eol]
			}
		}

		// Skip empty lines and comment-only lines
		trimmed := strings.TrimLeft(lineContent, " \t")
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Check indentation
		lineIndent := len(lineContent) - len(strings.TrimLeft(lineContent, " \t"))

		if lineIndent <= indentLen && !isOnlyWhitespace(lineContent) {
			// Found end of block
			blockEnd = searchFrom - 1
			break
		}

		blockEnd = searchFrom + len(lineContent)
	}

	chunkContent := content[match[0] : blockEnd+1]
	endLine := strings.Count(content[:blockEnd+1], "\n") + 1

	// Extract docstring
	doc := extractPythonDocstring(chunkContent)

	return Chunk{
		FilePath:   filePath,
		ChunkType:  chunkType,
		ChunkName:  name,
		ModulePath: modulePath,
		StartLine:  lineNum,
		EndLine:    endLine,
		Code:       strings.Join(lines[lineNum-1:endLine], "\n"),
		Doc:        doc,
		Lang:       "python",
		UpdatedAt:  time.Now().Unix(),
	}
}

func extractPythonModulePath(filePath string) string {
	return filepath.ToSlash(filepath.Dir(filePath))
}

func extractPythonFilePath(content string, lineStart int) string {
	// Return a placeholder since we need the original filePath
	// This will be overwritten by the caller
	return ""
}

func extractIndent(s string) string {
	return s
}

func getLineStart(content string, pos int) int {
	for pos > 0 {
		if content[pos-1] == '\n' {
			return pos
		}
		pos--
	}
	return 0
}

func getLineEnd(content string, pos int) int {
	for pos < len(content) && content[pos] != '\n' {
		pos++
	}
	return pos
}

func isOnlyWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

func findClassEnd(rest string, baseIndent int) int {
	lines := strings.SplitAfter(rest, "\n")
	for i, line := range lines {
		if isOnlyWhitespace(line) {
			continue
		}
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if indent < baseIndent {
			// Found line at same or lower indent
			if i == 0 {
				return 0
			}
			// Return position before this line
			pos := 0
			for j := 0; j < i; j++ {
				pos += len(lines[j])
			}
			return pos
		}
	}
	// End of content
	pos := 0
	for _, line := range lines {
		pos += len(line)
	}
	return pos
}

func extractPythonDocstring(content string) string {
	// Find triple-quoted string at the beginning
	content = strings.TrimLeft(content, " \t\n")

	// Check for triple double quote docstring
	if strings.HasPrefix(content, `"""`) {
		return extractTripleQuoteDocstring(content, `"""`)
	}

	// Check for triple single quote docstring
	if strings.HasPrefix(content, `'''`) {
		return extractTripleQuoteDocstring(content, `'''`)
	}

	return ""
}

func extractTripleQuoteDocstring(content, quote string) string {
	firstQuote := strings.Index(content, quote)
	if firstQuote < 0 {
		return ""
	}

	// Find the closing quote (not immediately after opening)
	start := firstQuote + len(quote)
	if start >= len(content) {
		return ""
	}

	// Check for immediate empty docstring
	if strings.HasPrefix(content[start:], quote) {
		return ""
	}

	// Find closing quote
	endQuote := strings.Index(content[start:], quote)
	if endQuote < 0 {
		return ""
	}

	doc := content[start : start+endQuote]
	// Clean up the docstring
	lines := strings.Split(doc, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return strings.Join(cleaned, " ")
}
