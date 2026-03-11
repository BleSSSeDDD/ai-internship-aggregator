package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
)

const MAX_RETRIES = 3

type vacancyParser struct {
	client *http.Client
	delay  <-chan time.Time
}

func NewParser() domain.Parser {
	client := &http.Client{Timeout: 10 * time.Second}
	ticker := time.NewTicker(1 * time.Second)

	return &vacancyParser{
		client: client,
		delay:  ticker.C,
	}
}

// func (p *vacancyParser) GetRawContent(ctx context.Context, url string) (string, error) {
// 	var lastErr error

// 	for attempt := 0; attempt < MAX_RETRIES; attempt++ {
// 		if attempt > 0 {
// 			time.Sleep(time.Duration(1<<attempt) * time.Second)
// 		}

// 		content, status, err := p.getWithDelay(ctx, url)
// 		if err == nil && status == http.StatusOK {
// 			return content, nil
// 		}

// 		if !shouldRetry(err, status) {
// 			return "", fmt.Errorf("неретраябельная ошибка: статус %d, err %v", status, err)
// 		}

// 		lastErr = err
// 		log.Printf("Попытка %d не удалась: %v", attempt+1, err)
// 	}

// 	return "", fmt.Errorf("все %d попыток провалились: %w", MAX_RETRIES, lastErr)
// }

// func (p *vacancyParser) getWithDelay(ctx context.Context, url string) (string, int, error) {
// 	select {
// 	case <-p.delay:
// 	case <-ctx.Done():
// 		return "", 0, ctx.Err()
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
// 	if err != nil {
// 		return "", 0, err
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
// 	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
// 	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
// 	req.Header.Set("Connection", "keep-alive")
// 	req.Header.Set("Upgrade-Insecure-Requests", "1")

// 	resp, err := p.client.Do(req)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	return string(body), resp.StatusCode, err
// }

// func shouldRetry(err error, statusCode int) bool {
// 	if err != nil {
// 		return true
// 	}

// 	if statusCode >= 500 && statusCode < 600 {
// 		return true
// 	}

// 	if statusCode == 429 {
// 		return true
// 	}

// 	return false
// }

// для мока вот такое есть:
func (p *vacancyParser) GetRawContent(ctx context.Context, url string) (string, error) {
	if url == "ivan sigma" {
		return string("qwe"), nil
	}
	return string("da eto tak)"), nil
}
