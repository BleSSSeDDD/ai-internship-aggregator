package kafka

import (
	"fmt"
	"log/slog"

	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/gen/go/vacancy"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type Publisher interface {
	SendInternship(internship *vacancy.CompanyInternship) (partition int32, offset int64, err error)
	Close() error
}

type publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewPublisher(brokers []string, topic string) (Publisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &publisher{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *publisher) SendInternship(internship *vacancy.CompanyInternship) (int32, int64, error) {
	data, err := proto.Marshal(internship)
	if err != nil {
		return 0, 0, fmt.Errorf("marshal protobuf: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(internship.CompanyName),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("send to kafka: %w", err)
	}

	slog.Info("message sent to kafka",
		"topic", p.topic,
		"partition", partition,
		"offset", offset,
	)
	return partition, offset, nil
}

func (p *publisher) Close() error {
	return p.producer.Close()
}
