package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"instaforge/internal/config"
	"instaforge/internal/logger"

	"github.com/IBM/sarama"
)

var (
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	onceP    sync.Once
	onceC    sync.Once
	ctx      = context.Background()
)

// Message wraps a queue message with metadata.
type Message struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	RetryCount int                   `json:"retry_count,omitempty"`
}

// ConnectProducer initializes the Kafka sync producer.
func ConnectProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	var err error
	onceP.Do(func() {
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Retry.Max = 5
		saramaConfig.Producer.Return.Successes = true
		saramaConfig.Producer.Return.Errors = true
		saramaConfig.Producer.Compression = sarama.CompressionSnappy
		saramaConfig.Producer.Flush.Frequency = 500 * time.Millisecond
		saramaConfig.ClientID = cfg.ClientID
		saramaConfig.Net.DialTimeout = 10 * time.Second
		saramaConfig.Net.ReadTimeout = 30 * time.Second
		saramaConfig.Net.WriteTimeout = 30 * time.Second

		producer, err = sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
		if err != nil {
			err = fmt.Errorf("failed to create kafka producer: %w", err)
			return
		}
		logger.Get().Info().Strs("brokers", cfg.Brokers).Msg("kafka producer connected")
	})
	return producer, err
}

// ConnectConsumer initializes the Kafka consumer group.
func ConnectConsumer(cfg config.KafkaConfig) (sarama.ConsumerGroup, error) {
	var err error
	onceC.Do(func() {
		saramaConfig := sarama.NewConfig()
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
		saramaConfig.Consumer.Return.Errors = true
		saramaConfig.Consumer.MaxProcessingTime = 30 * time.Second
		saramaConfig.ClientID = cfg.ClientID
		saramaConfig.Net.DialTimeout = 10 * time.Second

		consumer, err = sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, saramaConfig)
		if err != nil {
			err = fmt.Errorf("failed to create kafka consumer: %w", err)
			return
		}
		logger.Get().Info().Strs("brokers", cfg.Brokers).Str("group", cfg.ConsumerGroup).Msg("kafka consumer connected")
	})
	return consumer, err
}

// GetProducer returns the Kafka producer singleton.
func GetProducer() sarama.SyncProducer {
	if producer == nil {
		panic("kafka producer not initialized")
	}
	return producer
}

// GetConsumer returns the Kafka consumer singleton.
func GetConsumer() sarama.ConsumerGroup {
	if consumer == nil {
		panic("kafka consumer not initialized")
	}
	return consumer
}

// Publish sends a message to a Kafka topic.
func Publish(topic string, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("queue marshal error: %w", err)
	}

	partition, offset, err := GetProducer().SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{Key: []byte("message_type"), Value: []byte(msg.Type)},
			{Key: []byte("timestamp"), Value: []byte(msg.Timestamp.Format(time.RFC3339))},
		},
	})
	if err != nil {
		return fmt.Errorf("kafka publish error (topic=%s): %w", topic, err)
	}

	logger.Get().Debug().
		Str("topic", topic).
		Str("message_id", msg.ID).
		Int32("partition", partition).
		Int64("offset", offset).
		Msg("message published to kafka")

	return nil
}

// PublishRetry sends a message to the retry topic with incremented retry count.
func PublishRetry(cfg config.KafkaConfig, msg *Message) error {
	msg.RetryCount++
	return Publish(cfg.RetryTopic, msg)
}

// PublishDLQ sends a message to the dead letter queue after max retries.
func PublishDLQ(cfg config.KafkaConfig, msg *Message, reason string) error {
	// Copy payload to avoid mutating the original message
	dlqPayload := make(map[string]interface{}, len(msg.Payload)+1)
	for k, v := range msg.Payload {
		dlqPayload[k] = v
	}
	dlqPayload["dlq_reason"] = reason

	dlqMsg := &Message{
		ID:          msg.ID,
		Type:        msg.Type,
		Payload:     dlqPayload,
		Timestamp:   msg.Timestamp,
		RetryCount:  msg.RetryCount,
	}
	return Publish(cfg.DLQTopic, dlqMsg)
}

// HealthCheck verifies Kafka connectivity.
func HealthCheck() error {
	_, _, err := GetProducer().SendMessage(&sarama.ProducerMessage{
		Topic: "__health_check",
		Value: sarama.ByteEncoder([]byte("ping")),
	})
	return err
}

// Close gracefully shuts down Kafka connections.
func Close() error {
	if producer != nil {
		if err := producer.Close(); err != nil {
			return err
		}
	}
	if consumer != nil {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler.
type ConsumerGroupHandler struct {
	ProcessFn func(msg *Message) error
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var queueMsg Message
		if err := json.Unmarshal(msg.Value, &queueMsg); err != nil {
			logger.Get().Error().Err(err).Msg("failed to deserialize kafka message")
			session.MarkMessage(msg, "")
			continue
		}

		if err := h.ProcessFn(&queueMsg); err != nil {
			logger.Get().Error().
				Err(err).
				Str("message_id", queueMsg.ID).
				Str("topic", msg.Topic).
				Msg("consumer processing failed")
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
