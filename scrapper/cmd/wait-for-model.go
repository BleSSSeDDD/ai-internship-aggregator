package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type ollamaModel struct {
	Name string `json:"name"`
}

type ollamaListResponse struct {
	Models []ollamaModel `json:"models"`
}

func main() {
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://internship-ollama:11434"
	}

	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "qwen2.5:3b"
	}

	log.Printf("Waiting for Ollama model %s to be available...", modelName)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		available, err := checkModel(ollamaURL, modelName)
		if err != nil {
			log.Printf("Error checking model: %v", err)
		}

		if available {
			log.Printf("Model %s is available! Starting parser...", modelName)
			break
		}

		log.Printf("Model %s not yet available, waiting...", modelName)
		<-ticker.C
	}

	cmd := exec.Command("./ai-parser")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Fatalf("Parser exited with error: %v", err)
	}
}

func checkModel(ollamaURL, modelName string) (bool, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(ollamaURL + "/api/tags")
	if err != nil {
		return false, fmt.Errorf("Ollama not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	var listResp ollamaListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, model := range listResp.Models {
		if model.Name == modelName {
			return true, nil
		}
	}

	return false, nil
}
