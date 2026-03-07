package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	vacancy "github.com/BleSSSeDDD/reviewer-assignment/generated"
	"github.com/BleSSSeDDD/reviewer-assignment/internal/domain"
)

type publisher struct {
	writer *kafka.Writer
}

// NewPublisher создает настоящего publisher'а
func NewPublisher(brokers []string, topic string) domain.Publisher {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &publisher{writer: w}
}

func (p *publisher) Publish(ctx context.Context, internship *vacancy.CompanyInternship) error {
	// 1. Сериализуем Protobuf
	data, err := proto.Marshal(internship)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	// 2. Отправляем в Kafka
	msg := kafka.Message{
		Key:   []byte(internship.CompanyName),
		Value: data,
		Headers: []kafka.Header{
			{Key: "source", Value: []byte(internship.SourceSite)},
		},
		Time: time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write to kafka: %w", err)
	}

	return nil
}

func (p *publisher) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
