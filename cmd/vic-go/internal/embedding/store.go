package embedding

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vic-sdd/vic/internal/embedding/chunker"
	_ "modernc.org/sqlite"
)

// Store is the SQLite vector store for code chunks
type Store struct {
	db *sql.DB
}

// Manifest stores metadata about the embedding index
type Manifest struct {
	Version    string `json:"version"`
	Dimension  int    `json:"dimension"`
	ChunkCount int    `json:"chunk_count"`
	LastBuild  int64  `json:"last_build"`
	LastSync   int64  `json:"last_sync"`
	ProjectDir string `json:"project_dir"`
	Model      string `json:"model"`
}

// OpenStore opens or creates a SQLite store at the given path
func OpenStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{db: db}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS chunks (
		id          INTEGER PRIMARY KEY,
		file_path   TEXT NOT NULL,
		chunk_type  TEXT NOT NULL,
		chunk_name  TEXT NOT NULL,
		module_path TEXT NOT NULL,
		start_line  INTEGER,
		end_line    INTEGER,
		code        TEXT NOT NULL,
		doc         TEXT,
		lang        TEXT,
		updated_at  INTEGER
	);
	CREATE TABLE IF NOT EXISTS vectors (
		chunk_id INTEGER PRIMARY KEY,
		vector   BLOB
	);
	CREATE INDEX IF NOT EXISTS idx_file ON chunks(file_path);
	CREATE INDEX IF NOT EXISTS idx_module ON chunks(module_path);
	CREATE INDEX IF NOT EXISTS idx_updated ON chunks(updated_at);
	`
	_, err := s.db.Exec(schema)
	return err
}

// InsertChunks inserts chunks and their vectors into the store
func (s *Store) InsertChunks(chunks []chunker.Chunk, vectors [][]float64) error {
	if len(chunks) != len(vectors) {
		return fmt.Errorf("chunks and vectors count mismatch: %d vs %d", len(chunks), len(vectors))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	chunkStmt, err := tx.Prepare(`
		INSERT INTO chunks (file_path, chunk_type, chunk_name, module_path, start_line, end_line, code, doc, lang, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer chunkStmt.Close()

	vecStmt, err := tx.Prepare(`INSERT INTO vectors (chunk_id, vector) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer vecStmt.Close()

	for i, c := range chunks {
		result, err := chunkStmt.Exec(
			c.FilePath, c.ChunkType, c.ChunkName,
			c.ModulePath, c.StartLine, c.EndLine,
			c.Code, c.Doc, c.Lang, c.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert chunk %s.%s: %w", c.FilePath, c.ChunkName, err)
		}

		chunkID, _ := result.LastInsertId()
		vecBytes := float64SliceToBytes(vectors[i])
		_, err = vecStmt.Exec(chunkID, vecBytes)
		if err != nil {
			return fmt.Errorf("failed to insert vector: %w", err)
		}
	}

	return tx.Commit()
}

// Search finds the top-k most similar chunks to the query vector using cosine similarity
func (s *Store) Search(query []float64, topK int) ([]chunker.Chunk, error) {
	queryNorm := normalize(query)

	rows, err := s.db.Query(`
		SELECT c.id, c.file_path, c.chunk_type, c.chunk_name, c.module_path,
		       c.start_line, c.end_line, c.code, c.doc, c.lang, c.updated_at,
		       v.vector
		FROM chunks c JOIN vectors v ON c.id = v.chunk_id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query chunks: %w", err)
	}
	defer rows.Close()

	type scoredChunk struct {
		chunk chunker.Chunk
		score float64
	}
	var results []scoredChunk

	for rows.Next() {
		var c chunker.Chunk
		var vecBytes []byte
		err := rows.Scan(&c.ID, &c.FilePath, &c.ChunkType, &c.ChunkName,
			&c.ModulePath, &c.StartLine, &c.EndLine, &c.Code, &c.Doc,
			&c.Lang, &c.UpdatedAt, &vecBytes)
		if err != nil {
			return nil, err
		}

		vector := bytesToFloat64Slice(vecBytes)
		score := cosineSim(queryNorm, vector)
		results = append(results, scoredChunk{chunk: c, score: score})
	}

	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	var chunks []chunker.Chunk
	for i := 0; i < topK && i < len(results); i++ {
		chunks = append(chunks, results[i].chunk)
	}

	return chunks, nil
}

// GetManifest reads the manifest file
func (s *Store) GetManifest(manifestPath string) (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
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

// SetManifest writes the manifest file
func (s *Store) SetManifest(manifestPath string, m *Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(manifestPath, data, 0644)
}

// DeleteChunksByFile removes all chunks for a given file
func (s *Store) DeleteChunksByFile(filePath string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id FROM chunks WHERE file_path = ?", filePath)
	if err != nil {
		return 0, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rows.Close()

	if len(ids) == 0 {
		return 0, nil
	}

	for _, id := range ids {
		_, err := tx.Exec("DELETE FROM vectors WHERE chunk_id = ?", id)
		if err != nil {
			return 0, err
		}
	}

	result, err := tx.Exec("DELETE FROM chunks WHERE file_path = ?", filePath)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// ChunkCount returns the total number of chunks in the store
func (s *Store) ChunkCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM chunks").Scan(&count)
	return count, err
}

// Clear removes all chunks and vectors
func (s *Store) Clear() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("DELETE FROM vectors")
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM chunks")
	if err != nil {
		return err
	}
	return tx.Commit()
}

// ============================================================================
// Vector math helpers
// ============================================================================

func float64SliceToBytes(v []float64) []byte {
	b := make([]byte, len(v)*8)
	for i, f := range v {
		bits := float64Bits(f)
		b[i*8] = byte(bits)
		b[i*8+1] = byte(bits >> 8)
		b[i*8+2] = byte(bits >> 16)
		b[i*8+3] = byte(bits >> 24)
		b[i*8+4] = byte(bits >> 32)
		b[i*8+5] = byte(bits >> 40)
		b[i*8+6] = byte(bits >> 48)
		b[i*8+7] = byte(bits >> 56)
	}
	return b
}

func bytesToFloat64Slice(b []byte) []float64 {
	v := make([]float64, len(b)/8)
	for i := 0; i < len(v); i++ {
		bits := uint64(b[i*8]) |
			uint64(b[i*8+1])<<8 |
			uint64(b[i*8+2])<<16 |
			uint64(b[i*8+3])<<24 |
			uint64(b[i*8+4])<<32 |
			uint64(b[i*8+5])<<40 |
			uint64(b[i*8+6])<<48 |
			uint64(b[i*8+7])<<56
		v[i] = float64FromBits(bits)
	}
	return v
}

func float64Bits(f float64) uint64        { return uint64(f) }
func float64FromBits(bits uint64) float64 { return float64(bits) }

func normalize(v []float64) []float64 {
	var mag float64
	for _, f := range v {
		mag += f * f
	}
	mag = sqrt(mag)
	if mag == 0 {
		return v
	}
	result := make([]float64, len(v))
	for i, f := range v {
		result[i] = f / mag
	}
	return result
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 20; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func dot(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func cosineSim(a, b []float64) float64 {
	return dot(a, b)
}
