package main

import (
	"context"
	"log"
	"time"

	"github.com/BleSSSeDDD/reviewer-assignment/internal/infrastructure/aiprocessor"
	"github.com/BleSSSeDDD/reviewer-assignment/internal/infrastructure/httpclient"
	"github.com/BleSSSeDDD/reviewer-assignment/internal/infrastructure/kafka"
	"github.com/BleSSSeDDD/reviewer-assignment/internal/usecase"
)

func main() {
	p := httpclient.NewParser()
	a := aiprocessor.NewStub()
	k := kafka.NewPublisher(
		[]string{"internship-kafka:9092"}, // адрес Kafka из docker-compose
		"internships",                     // имя топика
	)
	defer k.Close()

	scrapper := usecase.NewScraperUsecase(p, a, k)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; ; i++ {
		if i%2 == 0 {
			err := scrapper.Run(ctx, "artemiy daun")
			if err != nil {
				log.Printf("❌ Ошибка: %v", err)
			} else {
				log.Println("✅ Успешно отправлено в Kafka")
			}
		} else {
			err := scrapper.Run(ctx, "ivan sigma")
			if err != nil {
				log.Printf("❌ Ошибка: %v", err)
			} else {
				log.Println("✅ Успешно отправлено в Kafka")
			}
		}
		time.Sleep(20 * time.Second)
	}
}
