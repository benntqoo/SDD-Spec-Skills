package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Embedder calls Ollama API to generate embeddings
type Embedder struct {
	model   string
	baseURL string
	client  *http.Client
}

// EmbeddingResponse is the Ollama API response
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// NewEmbedder creates a new Ollama embedder
func NewEmbedder() *Embedder {
	return &Embedder{
		model:   "all-minilm-l6-v2",
		baseURL: "http://localhost:11434",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsAvailable checks if Ollama is running and the model is available
func (e *Embedder) IsAvailable() bool {
	// First check if Ollama is running
	req, err := http.NewRequest("GET", e.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	// Check if the model is available
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(body, &tags); err != nil {
		return false
	}

	for _, m := range tags.Models {
		if strings.Contains(strings.ToLower(m.Name), "minilm") ||
			strings.Contains(strings.ToLower(m.Name), "nomic") ||
			strings.Contains(strings.ToLower(m.Name), "embed") {
			return true
		}
	}
	return false
}

// EmbedQuery embeds a single query text and returns the vector
func (e *Embedder) EmbedQuery(query string) ([]float64, error) {
	payload := map[string]interface{}{
		"model":   e.model,
		"prompt":  query,
		"options": map[string]interface{}{"embedding_only": true},
	}

	body, err := e.doRequest(payload)
	if err != nil {
		return nil, err
	}

	var resp EmbeddingResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}

	return resp.Embedding, nil
}

// Embed texts and return multiple vectors
func (e *Embedder) Embed(texts []string) ([][]float64, error) {
	var vectors [][]float64
	for _, text := range texts {
		vec, err := e.EmbedQuery(text)
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, vec)
	}
	return vectors, nil
}

func (e *Embedder) doRequest(payload map[string]interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", e.baseURL+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama at %s: %w", e.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error: %s - %s", resp.Status, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, nil
}
