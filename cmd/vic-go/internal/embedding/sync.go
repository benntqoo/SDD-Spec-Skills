package embedding

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"encoding/json"

	"github.com/vic-sdd/vic/internal/embedding/chunker"
)

// Sync handles incremental embedding index synchronization
type Sync struct {
	store        *Store
	embedder     *Embedder
	chunker      *chunker.Multiplexer
	projectDir   string
	indexFile    string
	manifestFile string
}

// NewSync creates a new Sync instance
func NewSync(projectDir, embeddingDir, indexFile string) *Sync {
	return &Sync{
		store:        nil,
		embedder:     NewEmbedder(),
		chunker:      chunker.NewMultiplexer(),
		projectDir:   projectDir,
		indexFile:    indexFile,
		manifestFile: filepath.Join(embeddingDir, "manifest.json"),
	}
}

// FullSync performs a complete rebuild of the embedding index
func (s *Sync) FullSync() error {
	if !s.embedder.IsAvailable() {
		return fmt.Errorf("Ollama is not available. Please install Ollama and pull the model:\n  ollama pull all-minilm-l6-v2")
	}

	store, err := OpenStore(s.indexFile)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer store.Close()

	// Clear existing data
	if err := store.Clear(); err != nil {
		return fmt.Errorf("failed to clear store: %w", err)
	}

	// Extract all chunks
	chunks, err := s.chunker.WalkAndExtract(s.projectDir)
	if err != nil {
		return fmt.Errorf("failed to extract chunks: %w", err)
	}

	if len(chunks) == 0 {
		// Update manifest with empty index
		m := &Manifest{
			Version:    "1.0",
			Dimension:  384,
			ChunkCount: 0,
			LastBuild:  time.Now().Unix(),
			LastSync:   time.Now().Unix(),
			ProjectDir: s.projectDir,
			Model:      s.embedder.model,
		}
		return store.SetManifest(s.manifestFile, m)
	}

	// Build embedding texts (doc + code for better context)
	texts := make([]string, len(chunks))
	for i, c := range chunks {
		text := c.ChunkName
		if c.Doc != "" {
			text += ". " + c.Doc
		}
		text += ". " + c.Code
		texts[i] = text
	}

	// Batch embed
	var allVectors [][]float64
	batchSize := 10
	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}
		vecs, err := s.embedder.Embed(texts[i:end])
		if err != nil {
			return fmt.Errorf("failed to embed batch %d-%d: %w", i, end, err)
		}
		allVectors = append(allVectors, vecs...)
	}

	// Store
	if err := store.InsertChunks(chunks, allVectors); err != nil {
		return fmt.Errorf("failed to insert chunks: %w", err)
	}

	// Update manifest
	m := &Manifest{
		Version:    "1.0",
		Dimension:  384,
		ChunkCount: len(chunks),
		LastBuild:  time.Now().Unix(),
		LastSync:   time.Now().Unix(),
		ProjectDir: s.projectDir,
		Model:      s.embedder.model,
	}
	return store.SetManifest(s.manifestFile, m)
}

// IncrementalSync checks file modification times and only re-processes changed files
func (s *Sync) IncrementalSync() (added, updated, removed int, err error) {
	if !s.embedder.IsAvailable() {
		return 0, 0, 0, nil // Silently skip if Ollama not available
	}

	store, err := OpenStore(s.indexFile)
	if err != nil {
		return 0, 0, 0, err
	}
	defer store.Close()

	manifest, err := store.GetManifest(s.manifestFile)
	if err != nil {
		// If manifest doesn't exist, do a full sync
		return 0, 0, 0, s.FullSync()
	}

	// Track which files we've seen
	seenFiles := make(map[string]bool)

	// Walk project files
	skipDirs := map[string]bool{
		".git": true, "vendor": true, "node_modules": true, ".venv": true,
		"venv": true, "testdata": true, "_test": true, ".idea": true,
		".vscode": true, "dist": true, "build": true, "target": true,
		"__pycache__": true, ".next": true, ".nuxt": true,
	}

	supportedExts := map[string]bool{
		".go": true, ".py": true, ".pyi": true,
		".ts": true, ".tsx": true, ".mts": true, ".cts": true,
		".js": true, ".mjs": true, ".cjs": true,
	}

	changedFiles := []string{}

	err = filepath.Walk(s.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(s.projectDir, path)
		rel = filepath.ToSlash(rel)

		if strings.HasPrefix(rel, ".vic-sdd") || strings.HasPrefix(rel, "embeddings") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExts[ext] {
			return nil
		}

		seenFiles[path] = true
		mtime := info.ModTime().Unix()

		if mtime > manifest.LastBuild {
			changedFiles = append(changedFiles, path)
		}
		return nil
	})

	if err != nil {
		return 0, 0, 0, err
	}

	// Process changed and new files
	for _, filePath := range changedFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			// File may have been deleted or is inaccessible, skip silently
			continue
		}

		chunks := s.chunker.ExtractChunks(filePath, string(content))
		if len(chunks) == 0 {
			// No extractable chunks in this file, skip silently
			continue
		}

		// Delete existing chunks for this file
		deleted, err := store.DeleteChunksByFile(filePath)
		if err != nil {
			// Failed to delete existing chunks, skip this file
			continue
		}
		if deleted > 0 {
			updated += int(deleted)
		}

		// Re-embed the chunks
		texts := make([]string, len(chunks))
		for i, c := range chunks {
			text := c.ChunkName
			if c.Doc != "" {
				text += ". " + c.Doc
			}
			text += ". " + c.Code
			texts[i] = text
		}

		vecs, err := s.embedder.Embed(texts)
		if err != nil {
			// Failed to generate embeddings, skip this file
			continue
		}

		if err := store.InsertChunks(chunks, vecs); err != nil {
			// Failed to insert chunks, skip this file
			continue
		}
		if deleted == 0 {
			added += len(chunks)
		}
	}

	// Handle removed files (deleted from project)
	// Get all files currently in the index and check if they still exist on disk
	indexedFiles, err := store.GetAllIndexedFiles()
	if err == nil {
		for _, filePath := range indexedFiles {
			// If the file was not seen during the walk, it has been deleted
			if !seenFiles[filePath] {
				deleted, delErr := store.DeleteChunksByFile(filePath)
				if delErr == nil {
					removed += int(deleted)
				}
			}
		}
	}

	// Update manifest
	manifest.LastSync = time.Now().Unix()
	count, _ := store.ChunkCount()
	manifest.ChunkCount = count
	if err := store.SetManifest(s.manifestFile, manifest); err != nil {
		return added, updated, removed, err
	}

	return added, updated, removed, nil
}

// Stats returns summary statistics about the embedding index
func (s *Sync) Stats() (*Manifest, error) {
	manifest, err := ReadManifest(s.manifestFile)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

// ReadManifest reads the manifest file
func ReadManifest(manifestFile string) (*Manifest, error) {
	data, err := os.ReadFile(manifestFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{Version: "1.0", Dimension: 384}, nil
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// DetectChangedFiles returns files modified since last sync
func (s *Sync) DetectChangedFiles() ([]string, error) {
	manifest, err := ReadManifest(s.manifestFile)
	if err != nil || manifest.LastBuild == 0 {
		// No manifest, return empty (will trigger full scan)
		return nil, nil
	}

	var changed []string
	supportedExts := map[string]bool{
		".go": true, ".py": true, ".pyi": true,
		".ts": true, ".tsx": true, ".mts": true, ".cts": true,
		".js": true, ".mjs": true, ".cjs": true,
	}
	skipDirs := map[string]bool{
		".git": true, "vendor": true, "node_modules": true, ".venv": true,
		"venv": true, "testdata": true, "_test": true, ".idea": true,
		".vscode": true, "dist": true, "build": true, "target": true,
		"__pycache__": true, ".next": true, ".nuxt": true,
	}

	err = filepath.Walk(s.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(s.projectDir, path)
		rel = filepath.ToSlash(rel)

		if strings.HasPrefix(rel, ".vic-sdd") || strings.HasPrefix(rel, "embeddings") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExts[ext] {
			return nil
		}

		if info.ModTime().Unix() > manifest.LastBuild {
			changed = append(changed, path)
		}
		return nil
	})

	sort.Strings(changed)
	return changed, nil
}
