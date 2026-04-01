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
	var format string

	cmd := &cobra.Command{
		Use:   "ask",
		Short: "Semantic code search using natural language",
		Long: `Search your codebase using natural language queries.
This command uses embedding-based semantic search to find code that matches your intent,
not just keyword matches.

Examples:
  vic ask "database connection pooling"
  vic ask "how does auth middleware work"
  vic ask --format json "find error handling patterns"

Requirements:
  - Ollama must be running (http://localhost:11434)
  - Run 'vic deps sync' to build the embedding index first`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAsk(cfg, args[0], format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}

func runAsk(cfg *config.Config, query string, format string) error {
	// Check if Ollama is available
	embedder := embedding.NewEmbedder()
	if !embedder.IsAvailable() {
		if format == "json" {
			fmt.Println(`{"success":false,"errors":[{"code":"OLLAMA_UNAVAILABLE","message":"Ollama is not available"}]}`)
		} else {
			fmt.Println("❌ Ollama is not available")
			fmt.Println("")
			fmt.Println("Please install Ollama and pull the embedding model:")
			fmt.Println("  1. Install Ollama: https://ollama.com")
			fmt.Println("  2. Pull the model: ollama pull all-minilm-l6-v2")
			fmt.Println("  3. Ollama will start automatically")
			fmt.Println("")
			fmt.Println("Then run 'vic deps sync' to build the embedding index.")
		}
		return nil
	}

	// Open the store
	indexFile := cfg.EmbeddingIndexFile
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		if format == "json" {
			fmt.Println(`{"success":false,"errors":[{"code":"INDEX_NOT_FOUND","message":"No embedding index found","hint":"Run 'vic deps sync' to build the index"}]}`)
		} else {
			fmt.Println("⚠️  No embedding index found.")
			fmt.Println("   Run 'vic deps sync' to build the index first.")
		}
		return nil
	}

	store, err := embedding.OpenStore(indexFile)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer store.Close()

	// Run incremental sync to pick up changes automatically
	sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
	var syncAdded, syncUpdated, syncRemoved int
	if added, updated, removed, syncErr := sync.IncrementalSync(); syncErr == nil {
		syncAdded, syncUpdated, syncRemoved = added, updated, removed
		if added+updated+removed > 0 && format != "json" {
			fmt.Printf("🔄 Auto-synced index: +%d ~%d -%d chunks\n", added, updated, removed)
		}
	} else if format != "json" {
		fmt.Printf("⚠️  Index sync failed: %v (results may be stale)\n", syncErr)
	}

	// Embed the query
	vector, err := embedder.EmbedQuery(query)
	if err != nil {
		return fmt.Errorf("failed to embed query: %w", err)
	}

	// Search
	results, err := store.Search(vector, 5)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// JSON output
	if format == "json" {
		fmt.Printf(`{"success":true,"message":"Search completed","data":{"query":%q,"sync":{"added":%d,"updated":%d,"removed":%d},"results":[`, query, syncAdded, syncUpdated, syncRemoved)
		for i, chunk := range results {
			relPath, _ := filepath.Rel(cfg.ProjectDir, chunk.FilePath)
			if strings.HasPrefix(relPath, ".") {
				relPath = chunk.FilePath
			}
			if i > 0 {
				fmt.Printf(",")
			}
			fmt.Printf(`{"file":%q,"start_line":%d,"end_line":%d,"lang":%q,"chunk_type":%q,"chunk_name":%q,"doc":%q}`,
				relPath, chunk.StartLine, chunk.EndLine, chunk.Lang, chunk.ChunkType, chunk.ChunkName, chunk.Doc)
		}
		fmt.Printf(`]}}`)
		fmt.Println()
		return nil
	}

	// Plain output
	fmt.Printf("🔍 Query: %s\n", query)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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
