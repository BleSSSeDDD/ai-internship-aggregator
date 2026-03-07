package main

import (
	"context"
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
		[]string{"localhost:9094"}, // адрес Kafka из docker-compose
		"internships",              // имя топика
	)
	defer k.Close()

	scrapper := usecase.NewScraperUsecase(p, a, k)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scrapper.Run(ctx, "artemiy daun")

	time.Sleep(time.Second)
}
