package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/metrics"
	"instaforge/internal/queue"

	"github.com/IBM/sarama"
)

// WorkerPool manages a pool of workers for processing download requests.
type WorkerPool struct {
	cfg         *config.KafkaConfig
	consumer    sarama.ConsumerGroup
	handler     *queue.ConsumerGroupHandler
	numWorkers  int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	activeJobs  int32
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(cfg config.KafkaConfig, processFn func(msg *queue.Message) error, numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		cfg:        &cfg,
		consumer:   queue.GetConsumer(),
		handler:    &queue.ConsumerGroupHandler{ProcessFn: processFn},
		numWorkers: numWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins consuming messages from Kafka.
func (wp *WorkerPool) Start() {
	log := logger.Get()

	log.Info().
		Str("topic", wp.cfg.DownloadTopic).
		Str("group", wp.cfg.ConsumerGroup).
		Int("workers", wp.numWorkers).
		Msg("starting worker pool")

	topics := []string{wp.cfg.DownloadTopic, wp.cfg.RetryTopic}

	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.consumeLoop(topics, i)
	}

	// Start DLQ cleaner
	go wp.dlqMonitor()

	log.Info().Msg("worker pool started")
}

// Stop gracefully shuts down the worker pool.
func (wp *WorkerPool) Stop() {
	logger.Get().Info().Msg("stopping worker pool...")
	wp.cancel()
	wp.wg.Wait()
	logger.Get().Info().Msg("worker pool stopped")
}

func (wp *WorkerPool) consumeLoop(topics []string, workerID int) {
	defer wp.wg.Done()

	log := logger.Get().With().Int("worker_id", workerID).Logger()

	for {
		select {
		case <-wp.ctx.Done():
			log.Info().Msg("worker shutting down")
			return
		default:
			err := wp.consumer.Consume(wp.ctx, topics, wp.handler)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Error().Err(err).Msg("consumer error, retrying in 5s")
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// dlqMonitor periodically checks DLQ and alerts on dead letters.
func (wp *WorkerPool) dlqMonitor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-ticker.C:
			logger.Get().Warn().Str("dlq_topic", wp.cfg.DLQTopic).Msg("DLQ monitor tick — check for dead letters")
		}
	}
}

// UpdateActiveJobs updates the active jobs metric.
func (wp *WorkerPool) UpdateActiveJobs(count int) {
	metrics.SetActiveWorkers(count)
}

// QueueDepth returns approximate queue depth (Kafka lag).
func (wp *WorkerPool) QueueDepth() (int64, error) {
	// In production, this would query Kafka consumer group lag
	// For now, return a placeholder
	return 0, fmt.Errorf("queue depth monitoring requires Kafka admin client")
}
