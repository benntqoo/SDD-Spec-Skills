package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vic-sdd/vic/internal/config"
	"github.com/vic-sdd/vic/internal/embedding"
)

// NewDepsSyncCmd creates the vic deps sync command
func NewDepsSyncCmd(cfg *config.Config) *cobra.Command {
	var fullSync bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync embedding index for vic ask",
		Long: `Build or update the embedding index used by vic ask.

This command extracts code chunks from your project, generates embeddings
using Ollama, and stores them in a local SQLite database.

By default, it performs an incremental sync (only changes since last build).
Use --full to force a complete rebuild.

Examples:
  vic deps sync         # Incremental sync (recommended)
  vic deps sync --full  # Full rebuild`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepsSync(cfg, fullSync)
		},
	}

	cmd.Flags().BoolVar(&fullSync, "full", false, "Force full rebuild of the embedding index")

	return cmd
}

func runDepsSync(cfg *config.Config, fullSync bool) error {
	// Check if Ollama is available
	embedder := embedding.NewEmbedder()
	if !embedder.IsAvailable() {
		fmt.Println("❌ Ollama is not available")
		fmt.Println("")
		fmt.Println("Please install Ollama and pull the embedding model:")
		fmt.Println("  1. Install Ollama: https://ollama.com")
		fmt.Println("  2. Pull the model: ollama pull all-minilm-l6-v2")
		fmt.Println("  3. Ollama will start automatically")
		return nil
	}

	fmt.Println("🔄 Syncing embedding index...")

	sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)

	start := time.Now()
	var err error
	var added, updated, removed int

	if fullSync {
		fmt.Println("   Mode: full rebuild")
		err = sync.FullSync()
		if err == nil {
			// Get stats after full sync
			store, storeErr := embedding.OpenStore(cfg.EmbeddingIndexFile)
			if storeErr == nil {
				count, _ := store.ChunkCount()
				added = count
				store.Close()
			}
		}
	} else {
		fmt.Println("   Mode: incremental")
		added, updated, removed, err = sync.IncrementalSync()
	}

	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Sync Summary:")

	if fullSync {
		fmt.Printf("   ✅ %d chunks indexed\n", added)
	} else {
		if added > 0 {
			fmt.Printf("   ✅ %d new chunks added\n", added)
		}
		if updated > 0 {
			fmt.Printf("   🔄 %d chunks updated\n", updated)
		}
		if removed > 0 {
			fmt.Printf("   🗑  %d chunks removed\n", removed)
		}
		if added == 0 && updated == 0 && removed == 0 {
			fmt.Println("   ✅ No changes detected")
		}
	}

	fmt.Printf("   ⏱  Completed in %v\n", elapsed.Round(time.Millisecond))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Embedding index is ready. Use 'vic ask' to search.")

	return nil
}
