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
	var format string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync embedding index for vic ask",
		Long: `Build or update the embedding index used by vic ask.

This command extracts code chunks from your project, generates embeddings
using Ollama, and stores them in a local SQLite database.

By default, it performs an incremental sync (only changes since last build).
Use --full to force a complete rebuild.

Examples:
  vic sync              # Incremental sync (recommended)
  vic sync --full       # Full rebuild
  vic sync --format json # JSON output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepsSync(cfg, fullSync, format)
		},
	}

	cmd.Flags().BoolVar(&fullSync, "full", false, "Force full rebuild of the embedding index")
	cmd.Flags().StringVarP(&format, "format", "f", "plain", "Output format (plain, json)")

	return cmd
}

func runDepsSync(cfg *config.Config, fullSync bool, format string) error {
	// Check if Ollama is available
	embedder := embedding.NewEmbedder()
	if !embedder.IsAvailable() {
		if format == "json" {
			fmt.Printf(`{"success":false,"message":"Ollama is not available","error":"Install Ollama from https://ollama.com and run: ollama pull all-minilm-l6-v2"}`)
			fmt.Println()
		} else {
			fmt.Println("❌ Ollama is not available")
			fmt.Println("")
			fmt.Println("Please install Ollama and pull the embedding model:")
			fmt.Println("  1. Install Ollama: https://ollama.com")
			fmt.Println("  2. Pull the model: ollama pull all-minilm-l6-v2")
			fmt.Println("  3. Ollama will start automatically")
		}
		return nil
	}

	if format != "json" {
		fmt.Println("🔄 Syncing embedding index...")
	}

	sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)

	start := time.Now()
	var err error
	var added, updated, removed int

	if fullSync {
		if format != "json" {
			fmt.Println("   Mode: full rebuild")
		}
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
		if format != "json" {
			fmt.Println("   Mode: incremental")
		}
		added, updated, removed, err = sync.IncrementalSync()
	}

	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	// JSON output
	if format == "json" {
		mode := "incremental"
		if fullSync {
			mode = "full"
		}
		fmt.Printf(`{"success":true,"message":"Sync completed","data":{"mode":%q,"added":%d,"updated":%d,"removed":%d,"elapsed_ms":%d}}`, mode, added, updated, removed, elapsed.Milliseconds())
		fmt.Println()
		return nil
	}

	// Plain output
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
