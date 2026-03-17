package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/aiprocessor"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/httpclient"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/kafka"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/usecase"
)

const (
	SITES_FILE     = "/config/sites.txt"
	MAX_CONCURRENT = 3
)

func main() {
	p := httpclient.NewParser()
	a := aiprocessor.NewAiProcessor("http://internship-ollama:11434", "qwen2.5:3b")
	k := kafka.NewPublisher(
		[]string{"internship-kafka:9092"},
		"internships",
	)
	defer k.Close()

	scrapper := usecase.NewScraperUsecase(p, a, k)

	urls, err := readSitesFromFile(SITES_FILE)
	if err != nil {
		log.Fatalf("Ошибка чтения файла с сайтами: %v", err)
	}

	if len(urls) == 0 {
		log.Fatal("Нет URL для парсинга в файле sites.txt")
	}

	log.Printf("Найдено %d URL для парсинга", len(urls))
	log.Printf("Запускаю парсинг с максимальной параллельностью %d", MAX_CONCURRENT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlChan := make(chan string, len(urls))
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	results := make(chan struct {
		url string
		err error
	}, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < MAX_CONCURRENT; i++ {
		wg.Add(1)
		go workerParser(ctx, &wg, scrapper, urlChan, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	failCount := 0

	for result := range results {
		if result.err != nil {
			log.Printf("Ошибка при парсинге %s: %v", result.url, result.err)
			failCount++
		} else {
			log.Printf("Успешно обработан: %s", result.url)
			successCount++
		}
	}

	log.Printf("Парсинг завершен. Успешно: %d, Ошибок: %d", successCount, failCount)
}

func workerParser(ctx context.Context, wg *sync.WaitGroup, scrapper *usecase.ScraperUsecase, urls <-chan string, results chan<- struct {
	url string
	err error
}) {
	defer wg.Done()

	for url := range urls {
		select {
		case <-ctx.Done():
			results <- struct {
				url string
				err error
			}{url: url, err: ctx.Err()}
			return
		default:
			// задержка между запросами к одному воркеру
			time.Sleep(2 * time.Second)

			log.Printf("Начинаю парсинг: %s", url)
			err := scrapper.Run(ctx, url)

			results <- struct {
				url string
				err error
			}{url: url, err: err}
		}
	}
}

func readSitesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && !strings.HasPrefix(url, "#") {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
