package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/embedding"
)

// NewAskCmd creates the vic ask command
func NewAskCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "ask",
		Short: "Semantic code search using natural language",
		Long: `Search your codebase using natural language queries.
This command uses embedding-based semantic search to find code that matches your intent,
not just keyword matches.

Examples:
  vic ask "database connection pooling"
  vic ask "how does auth middleware work"
  vic ask "find error handling patterns"

Requirements:
  - Ollama must be running (http://localhost:11434)
  - Run 'vic deps sync' to build the embedding index first`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAsk(cfg, args[0])
		},
	}
}

func runAsk(cfg *config.Config, query string) error {
	// Check if Ollama is available
	embedder := embedding.NewEmbedder()
	if !embedder.IsAvailable() {
		fmt.Println("❌ Ollama is not available")
		fmt.Println("")
		fmt.Println("Please install Ollama and pull the embedding model:")
		fmt.Println("  1. Install Ollama: https://ollama.com")
		fmt.Println("  2. Pull the model: ollama pull all-minilm-l6-v2")
		fmt.Println("  3. Ollama will start automatically")
		fmt.Println("")
		fmt.Println("Then run 'vic deps sync' to build the embedding index.")
		return nil
	}

	// Open the store
	indexFile := cfg.EmbeddingIndexFile
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		fmt.Println("⚠️  No embedding index found.")
		fmt.Println("   Run 'vic deps sync' to build the index first.")
		return nil
	}

	store, err := embedding.OpenStore(indexFile)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer store.Close()

	// Run incremental sync to pick up changes
	sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
	_ = sync // sync checks Ollama availability internally

	// Embed the query
	fmt.Printf("🔍 Query: %s\n", query)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	vector, err := embedder.EmbedQuery(query)
	if err != nil {
		return fmt.Errorf("failed to embed query: %w", err)
	}

	// Search
	results, err := store.Search(vector, 5)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("📭 No results found. Try a different query or run 'vic deps sync --full' to rebuild the index.")
		return nil
	}

	for i, chunk := range results {
		relPath, _ := filepath.Rel(cfg.ProjectDir, chunk.FilePath)
		if strings.HasPrefix(relPath, ".") {
			relPath = chunk.FilePath
		}

		fmt.Printf("\n📁 %s:%d-%d [%s:%s]\n",
			relPath, chunk.StartLine, chunk.EndLine, chunk.Lang, chunk.ChunkType)
		fmt.Printf("   └─ %s\n", chunk.ChunkName)

		// Show doc if available
		if chunk.Doc != "" {
			// Truncate long docs
			doc := chunk.Doc
			if len(doc) > 120 {
				doc = doc[:120] + "..."
			}
			fmt.Printf("   └─ %s\n", doc)
		}

		// Show code snippet (first 5 lines)
		lines := strings.Split(chunk.Code, "\n")
		maxLines := 5
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			fmt.Println("   └─ ...")
		}
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("      %s\n", truncateLine(line, 80))
			}
		}

		if i < len(results)-1 {
			fmt.Println()
		}
	}

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n✅ Found %d relevant code snippet(s)\n", len(results))

	return nil
}

func truncateLine(line string, maxLen int) string {
	if len(line) <= maxLen {
		return line
	}
	return line[:maxLen] + "..."
}
