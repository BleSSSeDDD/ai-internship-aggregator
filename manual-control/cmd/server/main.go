package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/handlers"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	kafkaProducer, err := kafka.NewProducer([]string{"internship-kafka:9092"})
	if err != nil {
		log.Fatalf("Ошибка создания Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	log.Println("Kafka producer инициализирован")

	h := handlers.NewHandlers(kafkaProducer)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.GET("/", h.Index)
	router.POST("/submit", h.Submit)
	router.GET("/health", h.Health)

	srv := &http.Server{
		Addr:         ":2228",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Сервер запущен на порту 2228")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	log.Println("Сервер остановлен")
}
