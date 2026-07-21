package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/infrastructure/aiprocessor"
	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/infrastructure/httpclient"
	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/infrastructure/kafka"
	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/usecase"
)

type config struct {
	ollamaURL      string
	ollamaModel    string
	aiTimeout      time.Duration
	kafkaBrokers   []string
	kafkaTopic     string
	sitesFile      string
	concurrency    int
	scrapeInterval time.Duration
}

func loadConfig() config {
	return config{
		ollamaURL:      envOr("OLLAMA_URL", "http://ollama:11434"),
		ollamaModel:    envOr("OLLAMA_MODEL", "qwen2.5:3b"),
		aiTimeout:      envDurationOr("AI_TIMEOUT", 10*time.Minute),
		kafkaBrokers:   strings.Split(envOr("KAFKA_BROKERS", "kafka:9092"), ","),
		kafkaTopic:     envOr("KAFKA_TOPIC", "internships"),
		sitesFile:      envOr("SITES_FILE", "/config/sites.txt"),
		concurrency:    envIntOr("SCRAPE_CONCURRENCY", 3),
		scrapeInterval: envDurationOr("SCRAPE_INTERVAL", 6*time.Hour),
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg := loadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	urls, err := readSitesFile(cfg.sitesFile)
	if err != nil {
		slog.Error("failed to read sites file", "file", cfg.sitesFile, "error", err)
		os.Exit(1)
	}
	if len(urls) == 0 {
		slog.Error("sites file contains no URLs", "file", cfg.sitesFile)
		os.Exit(1)
	}

	// Модель тянется отдельно (task ollama-pull), при холодном старте ждём её.
	if err := aiprocessor.WaitForModel(ctx, cfg.ollamaURL, cfg.ollamaModel); err != nil {
		slog.Error("ollama model is not available", "model", cfg.ollamaModel, "error", err)
		os.Exit(1)
	}

	publisher := kafka.NewPublisher(cfg.kafkaBrokers, cfg.kafkaTopic)
	defer publisher.Close()

	scraper := usecase.NewScraperUsecase(
		httpclient.NewParser(),
		aiprocessor.NewAiProcessor(cfg.ollamaURL, cfg.ollamaModel, cfg.aiTimeout),
		publisher,
	)

	slog.Info("scraper started",
		"urls", len(urls),
		"concurrency", cfg.concurrency,
		"interval", cfg.scrapeInterval,
	)

	runCycle(ctx, scraper, urls, cfg.concurrency)

	ticker := time.NewTicker(cfg.scrapeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down")
			return
		case <-ticker.C:
			runCycle(ctx, scraper, urls, cfg.concurrency)
		}
	}
}

// runCycle обходит все URL пулом воркеров и ждёт завершения обхода.
func runCycle(ctx context.Context, scraper *usecase.ScraperUsecase, urls []string, concurrency int) {
	start := time.Now()

	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	succeeded, failed := 0, 0

	for range concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				err := scraper.Run(ctx, url)

				mu.Lock()
				if err != nil {
					failed++
				} else {
					succeeded++
				}
				mu.Unlock()

				if err != nil && ctx.Err() == nil {
					slog.Error("failed to scrape", "url", url, "error", err)
				}
			}
		}()
	}

feed:
	for _, url := range urls {
		select {
		case jobs <- url:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()

	slog.Info("scrape cycle finished",
		"succeeded", succeeded,
		"failed", failed,
		"took", time.Since(start).Round(time.Second),
	)
}

// readSitesFile читает URL построчно, пропуская пустые строки и #-комментарии.
func readSitesFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		slog.Warn("invalid value, using default", "key", key, "value", v, "default", fallback)
		return fallback
	}
	return n
}

func envDurationOr(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		slog.Warn("invalid value, using default", "key", key, "value", v, "default", fallback)
		return fallback
	}
	return d
}
