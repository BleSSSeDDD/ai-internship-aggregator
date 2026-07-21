package httpclient

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/domain"
)

const (
	maxRetries     = 3
	clientTimeout  = 15 * time.Second
	requestDelay   = 1 * time.Second // общий rate limit на все воркеры
	retryBaseDelay = 2 * time.Second
	maxBodySize    = 10 << 20 // 10 MiB, защита от бесконечных ответов
)

type vacancyParser struct {
	client     *http.Client
	limiter    <-chan time.Time
	retryDelay time.Duration
}

func NewParser() domain.Parser {
	return &vacancyParser{
		client:     &http.Client{Timeout: clientTimeout},
		limiter:    time.NewTicker(requestDelay).C,
		retryDelay: retryBaseDelay,
	}
}

// GetRawContent скачивает страницу с ретраями и экспоненциальным бэкоффом.
func (p *vacancyParser) GetRawContent(ctx context.Context, url string) (string, error) {
	slog.Info("fetching page", "url", url)
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(p.retryDelay << (attempt - 1)):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		content, status, err := p.get(ctx, url)
		if err == nil && status == http.StatusOK {
			return content, nil
		}

		if !shouldRetry(err, status) {
			return "", fmt.Errorf("non-retryable response: status %d, err %v", status, err)
		}

		lastErr = fmt.Errorf("status %d: %w", status, err)
		slog.Warn("attempt failed", "url", url, "attempt", attempt+1, "status", status, "error", err)
	}

	return "", fmt.Errorf("all %d attempts failed: %w", maxRetries, lastErr)
}

func (p *vacancyParser) get(ctx context.Context, url string) (string, int, error) {
	select {
	case <-p.limiter:
	case <-ctx.Done():
		return "", 0, ctx.Err()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	return string(body), resp.StatusCode, err
}

func shouldRetry(err error, statusCode int) bool {
	if err != nil {
		return true
	}
	return statusCode >= 500 || statusCode == http.StatusTooManyRequests
}
