package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTransactionCompleted(ctx context.Context, event models.TransactionEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.TransactionID),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("writer.WriteMessages: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
