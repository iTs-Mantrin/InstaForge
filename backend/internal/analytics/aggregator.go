package analytics

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"instaforge/internal/database"
	"instaforge/internal/logger"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// Aggregator periodically rolls raw AnalyticsEvents into AnalyticsSummaries.
type Aggregator struct {
	db       *gorm.DB
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
	lastRun  time.Time
	mu       sync.Mutex
}

// NewAggregator creates an aggregator that runs every `interval`.
func NewAggregator(interval time.Duration) *Aggregator {
	return &Aggregator{
		db:       database.Get(),
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start launches the background aggregation loop.
func (a *Aggregator) Start() {
	a.wg.Add(1)
	go a.loop()
	logger.Get().Info().Dur("interval", a.interval).Msg("analytics aggregator started")
}

// Stop signals the aggregator to shut down and waits for completion.
func (a *Aggregator) Stop() {
	close(a.stopCh)
	a.wg.Wait()
	logger.Get().Info().Msg("analytics aggregator stopped")
}

func (a *Aggregator) loop() {
	defer a.wg.Done()

	// Run once immediately on startup
	a.runOnce()

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.runOnce()
		case <-a.stopCh:
			return
		}
	}
}

func (a *Aggregator) runOnce() {
	a.mu.Lock()
	cutoff := a.lastRun
	now := time.Now().UTC()
	a.lastRun = now
	a.mu.Unlock()

	if cutoff.IsZero() {
		// First run — aggregate last hour's data
		cutoff = now.Add(-1 * time.Hour)
	}

	log := logger.Get()
	log.Debug().Time("since", cutoff).Msg("analytics aggregation starting")

	summary, err := a.computeSummary(cutoff, now)
	if err != nil {
		log.Error().Err(err).Msg("analytics aggregation failed")
		return
	}

	if err := a.persistSummary(summary); err != nil {
		log.Error().Err(err).Msg("analytics aggregation persist failed")
		return
	}

	log.Info().
		Int("endpoints", len(summary)).
		Time("from", cutoff).
		Time("to", now).
		Msg("analytics aggregation complete")
}

type endpointStats struct {
	endpoint string
	total    int64
	success  int64
	errors   int64
	durations []float64
}

func (a *Aggregator) computeSummary(since, until time.Time) ([]models.AnalyticsSummary, error) {
	var events []models.AnalyticsEvent
	err := a.db.Where("created_at > ? AND created_at <= ?", since, until).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("query analytics events: %w", err)
	}

	if len(events) == 0 {
		return nil, nil
	}

	// Group by endpoint
	grouped := make(map[string]*endpointStats)
	for _, ev := range events {
		dateKey := ev.CreatedAt.Format("2006-01-02") + "|" + ev.Endpoint
		s, ok := grouped[dateKey]
		if !ok {
			s = &endpointStats{endpoint: ev.Endpoint}
			grouped[dateKey] = s
		}
		s.total++
		s.durations = append(s.durations, float64(ev.DurationMs))
		if ev.StatusCode >= 200 && ev.StatusCode < 400 {
			s.success++
		} else {
			s.errors++
		}
	}

	summaries := make([]models.AnalyticsSummary, 0, len(grouped))
	for dateKey, s := range grouped {
		// Parse date
		dateStr := dateKey[:10]
		date, _ := time.Parse("2006-01-02", dateStr)

		// Sort durations for percentile calculation
		sort.Float64s(s.durations)

		avgDuration := 0.0
		p99Duration := 0.0
		if len(s.durations) > 0 {
			sum := 0.0
			for _, d := range s.durations {
				sum += d
			}
			avgDuration = sum / float64(len(s.durations))

			p99Idx := int(math.Ceil(float64(len(s.durations))*0.99) - 1)
			if p99Idx < 0 {
				p99Idx = 0
			}
			p99Duration = s.durations[p99Idx]
		}

		summaries = append(summaries, models.AnalyticsSummary{
			Date:          date,
			Endpoint:      s.endpoint,
			TotalReq:      s.total,
			SuccessReq:    s.success,
			ErrorReq:      s.errors,
			AvgDurationMs: math.Round(avgDuration*100) / 100,
			P99DurationMs: math.Round(p99Duration*100) / 100,
		})
	}

	return summaries, nil
}

func (a *Aggregator) persistSummary(summaries []models.AnalyticsSummary) error {
	if len(summaries) == 0 {
		return nil
	}

	for i := range summaries {
		s := &summaries[i]
		err := a.db.Where("date = ? AND endpoint = ?", s.Date, s.Endpoint).
			Assign(map[string]interface{}{
				"total_req":       s.TotalReq,
				"success_req":     s.SuccessReq,
				"error_req":       s.ErrorReq,
				"avg_duration_ms": s.AvgDurationMs,
				"p99_duration_ms": s.P99DurationMs,
			}).
			FirstOrCreate(s).Error
		if err != nil {
			return fmt.Errorf("upsert summary for %s/%s: %w", s.Date, s.Endpoint, err)
		}
	}

	return nil
}
