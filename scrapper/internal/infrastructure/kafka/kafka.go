package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) domain.Publisher {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		RequiredAcks: kafka.RequireOne,
	}

	return &publisher{writer: w}
}

func (p *publisher) Publish(ctx context.Context, internship *vacancy.CompanyInternship) error {
	data, err := proto.Marshal(internship)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

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
