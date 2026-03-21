package kafka

import (
	"fmt"
	"log"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    "internships",
	}, nil
}

func (p *Producer) SendInternship(internship *vacancy.CompanyInternship) (int32, int64, error) {
	// Сериализуем в protobuf
	data, err := proto.Marshal(internship)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка сериализации protobuf: %w", err)
	}

	// Отправляем в Kafka
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
		Key:   sarama.StringEncoder(internship.CompanyName),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка отправки в Kafka: %w", err)
	}

	log.Printf("📤 Отправлено в Kafka [%s] partition=%d offset=%d", p.topic, partition, offset)
	return partition, offset, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
