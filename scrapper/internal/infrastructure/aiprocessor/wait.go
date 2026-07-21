package aiprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const modelPollInterval = 10 * time.Second

type ollamaModel struct {
	Name string `json:"name"`
}

type ollamaTagsResponse struct {
	Models []ollamaModel `json:"models"`
}

// WaitForModel блокируется, пока модель не станет доступна в Ollama
// (при холодном старте её ещё тянет ollama pull), либо пока не отменят контекст.
func WaitForModel(ctx context.Context, ollamaURL, model string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	ticker := time.NewTicker(modelPollInterval)
	defer ticker.Stop()

	for {
		available, err := modelAvailable(ctx, client, ollamaURL, model)
		if available {
			slog.Info("ollama model is available", "model", model)
			return nil
		}
		if err != nil {
			slog.Warn("ollama is not ready yet", "error", err)
		} else {
			slog.Info("waiting for model to be pulled", "model", model)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func modelAvailable(ctx context.Context, client *http.Client, ollamaURL, model string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ollamaURL+"/api/tags", nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ollama is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var tags ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false, fmt.Errorf("decode ollama tags: %w", err)
	}

	for _, m := range tags.Models {
		if m.Name == model {
			return true, nil
		}
	}
	return false, nil
}
