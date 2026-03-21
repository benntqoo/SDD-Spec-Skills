package chunker

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Multiplexer combines all language-specific chunkers
type Multiplexer struct {
	chunkers []Chunker
}

// NewMultiplexer creates a new multiplexer with all supported chunkers
func NewMultiplexer() *Multiplexer {
	return &Multiplexer{
		chunkers: []Chunker{
			&GoChunker{},
			&PythonChunker{},
			&TypeScriptChunker{},
		},
	}
}

// ExtractChunks extracts chunks from a single file, delegating to the appropriate language chunker
func (m *Multiplexer) ExtractChunks(filePath string, content string) []Chunk {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, c := range m.chunkers {
		for _, e := range c.Extensions() {
			if e == ext {
				return c.ExtractChunks(filePath, content)
			}
		}
	}
	return nil
}

// WalkAndExtract walks a directory and extracts all chunks from supported source files
func (m *Multiplexer) WalkAndExtract(projectDir string) ([]Chunk, error) {
	var allChunks []Chunk

	skipDirs := map[string]bool{
		".git": true, "vendor": true, "node_modules": true, ".venv": true,
		"venv": true, "testdata": true, "_test": true, ".idea": true,
		".vscode": true, "dist": true, "build": true, "target": true,
		"__pycache__": true, ".next": true, ".nuxt": true,
	}

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(projectDir, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)

		// Skip .vic-sdd and embeddings dirs
		if strings.HasPrefix(rel, ".vic-sdd") || strings.HasPrefix(rel, "embeddings") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			// Skip certain directories
			dirName := filepath.Base(path)
			if skipDirs[dirName] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip non-source files
		ext := strings.ToLower(filepath.Ext(path))
		supported := false
		for _, c := range m.chunkers {
			for _, e := range c.Extensions() {
				if e == ext {
					supported = true
					break
				}
			}
		}
		if !supported {
			return nil
		}

		// Read and parse file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		chunks := m.ExtractChunks(path, string(content))
		for i := range chunks {
			chunks[i].UpdatedAt = time.Now().Unix()
		}
		allChunks = append(allChunks, chunks...)
		return nil
	})

	return allChunks, err
}
