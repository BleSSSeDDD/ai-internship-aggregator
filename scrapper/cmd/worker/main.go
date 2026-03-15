package main

import (
	"context"
	"log"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/aiprocessor"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/httpclient"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/kafka"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/usecase"
)

func main() {
	p := httpclient.NewParser()
	a := aiprocessor.NewStub()
	k := kafka.NewPublisher(
        []string{"localhost:9094"},
        "internships",
    )
	defer k.Close()

	scrapper := usecase.NewScraperUsecase(p, a, k)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; ; i++ {
		if i%2 == 0 {
			err := scrapper.Run(ctx, "крутая ссылка")
			if err != nil {
				log.Printf("Ошибка: %v", err)
			} else {
				log.Println("Успешно отправлено в Kafka")
			}
		} else {
			err := scrapper.Run(ctx, "ivan sigma")
			if err != nil {
				log.Printf("Ошибка: %v", err)
			} else {
				log.Println("Успешно отправлено в Kafka")
			}
		}
		time.Sleep(20 * time.Second)
	}
}
