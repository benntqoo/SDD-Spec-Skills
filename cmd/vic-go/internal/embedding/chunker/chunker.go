package chunker

// Chunk represents a single code chunk extracted from source files
type Chunk struct {
	ID         int64
	FilePath   string
	ChunkType  string // func, class, def, struct, module
	ChunkName  string // function/class name
	ModulePath string // e.g. internal/commands
	StartLine  int
	EndLine    int
	Code       string
	Doc        string
	Lang       string
	UpdatedAt  int64
}

// Chunker extracts code chunks from source files
type Chunker interface {
	Name() string
	Extensions() []string
	ExtractChunks(filePath string, content string) []Chunk
}
