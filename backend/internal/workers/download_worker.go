package workers

import (
	"encoding/json"
	"fmt"
	"time"

	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/queue"
	"instaforge/internal/services"
)

// DownloadWorkerHandler creates the message handler for download workers.
func DownloadWorkerHandler(downloadService *services.DownloadService, cfg *config.KafkaConfig) func(msg *queue.Message) error {
	return func(msg *queue.Message) error {
		start := time.Now()

		log := logger.Get().With().
			Str("message_id", msg.ID).
			Str("type", msg.Type).
			Logger()

		log.Info().Msg("worker processing download request")

		// Extract request ID from payload
		requestID, ok := msg.Payload["request_id"].(string)
		if !ok {
			return fmt.Errorf("missing request_id in message payload")
		}

		// Process the download
		err := downloadService.ProcessDownload(requestID)
		if err != nil {
			log.Error().Err(err).Str("request_id", requestID).Msg("download processing failed")

			// Check retry count
			if msg.RetryCount < 3 {
				// Send to retry topic
				if retryErr := queue.PublishRetry(*cfg, msg); retryErr != nil {
					log.Error().Err(retryErr).Msg("failed to send to retry topic")
				}
				log.Info().Int("retry_count", msg.RetryCount).Msg("sent to retry topic")
			} else {
				// Send to dead letter queue
				if dlqErr := queue.PublishDLQ(*cfg, msg, err.Error()); dlqErr != nil {
					log.Error().Err(dlqErr).Msg("failed to send to DLQ")
				}
				log.Warn().Msg("max retries reached, sent to DLQ")
			}

			return err
		}

		elapsed := time.Since(start)
		log.Info().
			Dur("elapsed", elapsed).
			Str("request_id", requestID).
			Msg("download processed successfully")

		return nil
	}
}

// RetryWorkerHandler processes retry topic messages.
func RetryWorkerHandler(downloadService *services.DownloadService, cfg *config.KafkaConfig) func(msg *queue.Message) error {
	return func(msg *queue.Message) error {
		log := logger.Get().With().
			Str("message_id", msg.ID).
			Int("retry_count", msg.RetryCount).
			Logger()

		log.Info().Msg("retry worker processing message")

		requestID, ok := msg.Payload["request_id"].(string)
		if !ok {
			return fmt.Errorf("missing request_id in retry message")
		}

		err := downloadService.ProcessDownload(requestID)
		if err != nil {
			log.Error().Err(err).Msg("retry processing failed")

			if msg.RetryCount >= 3 {
				if dlqErr := queue.PublishDLQ(*cfg, msg, err.Error()); dlqErr != nil {
					log.Error().Err(dlqErr).Msg("failed to send retry to DLQ")
				}
			}
			return err
		}

		log.Info().Msg("retry processed successfully")
		return nil
	}
}

// DLQAnalyzer logs and analyzes dead letter queue messages.
func DLQAnalyzer() func(msg *queue.Message) error {
	return func(msg *queue.Message) error {
		data, err := json.Marshal(msg)
		payloadStr := ""
		if err != nil {
			payloadStr = fmt.Sprintf("{\"error\": \"failed to marshal: %v\", \"message_id\": %q}", err, msg.ID)
		} else {
			payloadStr = string(data)
		}
		logger.Get().Error().
			Str("message_id", msg.ID).
			Str("type", msg.Type).
			Int("retry_count", msg.RetryCount).
			Str("payload", payloadStr).
			Msg("DLQ: dead letter message")
		return nil
	}
}
