package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	kafka "github.com/segmentio/kafka-go"
)

const (
	readerMinBytes = 10_000     // 10 KB
	readerMaxBytes = 10_000_000 // 10 MB
)

type LedgerService interface {
	RecordEntry(ctx context.Context, event models.TransactionCompletedEvent) error
}

type Consumer struct {
	reader  *kafka.Reader
	service LedgerService
	log     *slog.Logger
}

func NewConsumer(cfg config.Config, service LedgerService, log *slog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: readerMinBytes,
		MaxBytes: readerMaxBytes,
	})

	return &Consumer{
		reader:  reader,
		service: service,
		log:     log,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	c.log.InfoContext(ctx, "kafka consumer started",
		slog.String("topic", c.reader.Config().Topic),
		slog.String("group", c.reader.Config().GroupID),
	)

	defer func() {
		if err := c.reader.Close(); err != nil {
			c.log.ErrorContext(ctx, "kafka reader close error", slog.Any("error", err))
		}
	}()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.InfoContext(ctx, "kafka consumer stopped")
				return nil
			}
			c.log.ErrorContext(ctx, "kafka FetchMessage error", slog.Any("error", err))
			continue
		}

		var event models.TransactionCompletedEvent
		if err = json.Unmarshal(msg.Value, &event); err != nil {
			c.log.ErrorContext(ctx, "kafka message unmarshal error",
				slog.Any("error", err),
				slog.Int64("offset", msg.Offset),
				slog.Int("partition", msg.Partition),
			)
			// Commit poison pill to avoid infinite retry
			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				c.log.ErrorContext(ctx, "kafka CommitMessages error after unmarshal failure",
					slog.Any("error", commitErr),
				)
			}
			continue
		}

		//FIXME: Handle idempotent duplicate events in the consumer.
		// In an at-least-once Kafka delivery model, a repeated TransactionCompleted event may be received
		// for an already processed transaction. ConflictError should be treated as a successful
		// idempotent re-delivery: commit the offset and move on. Only transient errors should prevent
		// the offset from being committed and trigger a retry.
		if err = c.service.RecordEntry(ctx, event); err != nil {
			c.log.ErrorContext(ctx, "RecordEntry error",
				slog.Any("error", err),
				slog.String("transaction_id", event.TransactionID),
				slog.Int64("offset", msg.Offset),
			)
			// Do not commit — retry on next consumer start
			continue
		}

		if err = c.reader.CommitMessages(ctx, msg); err != nil {
			c.log.ErrorContext(ctx, "kafka CommitMessages error",
				slog.Any("error", err),
				slog.Int64("offset", msg.Offset),
			)
		}
	}
}
