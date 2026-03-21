package chunker

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GoChunker extracts exported Go declarations using AST
type GoChunker struct{}

// Ensure GoChunker implements Chunker
var _ Chunker = (*GoChunker)(nil)

func (GoChunker) Name() string         { return "Go" }
func (GoChunker) Extensions() []string { return []string{".go"} }

func (g *GoChunker) ExtractChunks(filePath string, content string) []Chunk {
	var chunks []Chunk

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil
	}

	modulePath := extractModulePath(filePath)
	lines := strings.Split(content, "\n")

	for _, decl := range file.Decls {
		var chunk Chunk
		chunk.Lang = "go"
		chunk.FilePath = filePath
		chunk.ModulePath = modulePath
		chunk.UpdatedAt = time.Now().Unix()

		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Only export functions (starts with uppercase)
			if !d.Name.IsExported() {
				continue
			}
			chunk.ChunkType = "func"
			chunk.ChunkName = d.Name.Name
			chunk.StartLine = fset.Position(d.Pos()).Line
			chunk.EndLine = fset.Position(d.End()).Line
			chunk.Code = extractCodeLines(lines, chunk.StartLine, chunk.EndLine)
			chunk.Doc = extractDoc(d.Doc)
			chunk.Doc = cleanGoDoc(chunk.Doc)

		case *ast.GenDecl:
			// Handle type, const, var declarations
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if !typeSpec.Name.IsExported() {
							continue
						}
						chunk.ChunkType = "type"
						chunk.ChunkName = typeSpec.Name.Name
						chunk.StartLine = fset.Position(typeSpec.Pos()).Line
						chunk.EndLine = fset.Position(typeSpec.End()).Line
						chunk.Code = extractCodeLines(lines, chunk.StartLine, chunk.EndLine)
						chunk.Doc = extractDoc(d.Doc)
						chunk.Doc = cleanGoDoc(chunk.Doc)
						chunks = append(chunks, chunk)
					}
				}
				continue

			} else if d.Tok == token.CONST || d.Tok == token.VAR {
				for _, spec := range d.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range valueSpec.Names {
							if !name.IsExported() {
								continue
							}
							chunk = Chunk{
								Lang:       "go",
								FilePath:   filePath,
								ModulePath: modulePath,
								UpdatedAt:  time.Now().Unix(),
							}

							if d.Tok == token.CONST {
								chunk.ChunkType = "const"
							} else {
								chunk.ChunkType = "var"
							}
							chunk.ChunkName = name.Name

							if i == 0 && valueSpec.Values != nil && len(valueSpec.Values) > 0 {
								chunk.StartLine = fset.Position(valueSpec.Pos()).Line
								chunk.EndLine = fset.Position(valueSpec.End()).Line
							} else {
								chunk.StartLine = fset.Position(name.Pos()).Line
								chunk.EndLine = chunk.StartLine
							}

							chunk.Code = extractCodeLines(lines, chunk.StartLine, chunk.EndLine)

							if i == 0 {
								chunk.Doc = extractDoc(d.Doc)
								chunk.Doc = cleanGoDoc(chunk.Doc)
							}

							chunks = append(chunks, chunk)
						}
					}
				}
				continue
			}
		}

		if chunk.ChunkName != "" {
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}

func extractModulePath(filePath string) string {
	// Get relative path from project root
	rel, err := filepath.Rel(findProjectRoot(filePath), filePath)
	if err != nil {
		return filepath.Dir(filePath)
	}
	// Get directory containing the file
	dir := filepath.Dir(rel)
	if dir == "." {
		return "root"
	}
	return filepath.ToSlash(dir)
}

func findProjectRoot(filePath string) string {
	dir := filepath.Dir(filePath)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Dir(filePath)
}

func extractCodeLines(lines []string, startLine, endLine int) string {
	if startLine < 1 {
		startLine = 1
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine > endLine {
		return ""
	}

	startLine-- // Convert to 0-indexed
	codeLines := lines[startLine:endLine]
	return strings.Join(codeLines, "\n")
}

func extractDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	return doc.Text()
}

func cleanGoDoc(doc string) string {
	// Remove leading comment markers and normalize whitespace
	lines := strings.Split(doc, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, " ")
}
