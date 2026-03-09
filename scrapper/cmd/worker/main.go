package main

import (
	"context"
	"log"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/aiprocessor"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/httpclient"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/infrastructure/kafka"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/usecase"
)

func main() {
	p := httpclient.NewParser()
	a := aiprocessor.NewAiProcessor("http://internship-ollama:11434", "qwen2.5:3b")
	k := kafka.NewPublisher(
		[]string{"internship-kafka:9092"}, // адрес Kafka из docker-compose
		"internships",                     // имя топика
	)
	defer k.Close()

	scrapper := usecase.NewScraperUsecase(p, a, k)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Я ЖИВОЙ Я ЖИВОЙ ЩА ПОДОЖДИ ЧУТКА")

	err := scrapper.Run(ctx, "https://yandex.ru/yaintern/backend")
	if err != nil {
		log.Printf("Ошибка: %v", err)
	} else {
		log.Println("Успешно отправлено в Kafka")
	}
}
