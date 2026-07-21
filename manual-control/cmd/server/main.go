package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/internal/handlers"
	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/internal/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	brokers := strings.Split(envOr("KAFKA_BROKERS", "kafka:9092"), ",")
	topic := envOr("KAFKA_TOPIC", "internships")
	port := envOr("PORT", "2228")

	publisher, err := kafka.NewPublisher(brokers, topic)
	if err != nil {
		slog.Error("failed to create kafka producer", "brokers", brokers, "error", err)
		os.Exit(1)
	}
	defer publisher.Close()

	slog.Info("kafka producer ready", "brokers", brokers, "topic", topic)

	h := handlers.NewHandlers(publisher)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", h.Index)
	router.POST("/submit", h.Submit)
	router.GET("/health", h.Health)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("admin panel listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
