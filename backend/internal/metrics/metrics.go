package metrics

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently in flight",
		},
	)

	// Business metrics
	downloadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "downloads_total",
			Help: "Total number of download requests",
		},
		[]string{"type", "status"},
	)

	searchesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "searches_total",
			Help: "Total number of username searches",
		},
		[]string{"status"},
	)

	// Cache metrics
	cacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Worker metrics
	workerQueueDepth = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_queue_depth",
			Help: "Current depth of the worker queue",
		},
	)

	workerActiveJobs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_active_jobs",
			Help: "Number of currently active worker jobs",
		},
	)

	// Error metrics
	errorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"type", "source"},
	)

	// Rate limit metrics
	rateLimitHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ratelimit_hits_total",
			Help: "Total number of rate limit hits",
		},
	)
)

// HTTPMetricsMiddleware records HTTP metrics for each request.
func HTTPMetricsMiddleware(c *fiber.Ctx) error {
	httpRequestInFlight.Inc()
	c.Next()

	status := strconv.Itoa(c.Response().StatusCode())
	path := c.Route().Path

	httpRequestsTotal.WithLabelValues(c.Method(), path, status).Inc()
	tookMs, _ := c.Locals("took_ms").(int64)
	httpRequestDuration.WithLabelValues(c.Method(), path, status).Observe(
		float64(tookMs) / 1000.0,
	)
	httpRequestInFlight.Dec()

	return nil
}

// RecordDownload records a download metric.
func RecordDownload(downloadType string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	downloadsTotal.WithLabelValues(downloadType, status).Inc()
}

// RecordSearch records a search metric.
func RecordSearch(success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	searchesTotal.WithLabelValues(status).Inc()
}

// RecordCacheHit increments the cache hit counter.
func RecordCacheHit() {
	cacheHitsTotal.Inc()
}

// RecordCacheMiss increments the cache miss counter.
func RecordCacheMiss() {
	cacheMissesTotal.Inc()
}

// RecordError records an error metric.
func RecordError(errorType, source string) {
	errorTotal.WithLabelValues(errorType, source).Inc()
}

// RecordRateLimitHit increments the rate limit counter.
func RecordRateLimitHit() {
	rateLimitHitsTotal.Inc()
}

// SetQueueDepth sets the current queue depth.
func SetQueueDepth(depth int64) {
	workerQueueDepth.Set(float64(depth))
}

// SetActiveWorkers sets the number of active workers.
func SetActiveWorkers(count int) {
	workerActiveJobs.Set(float64(count))
}

// Handler returns the Prometheus metrics HTTP handler.
func Handler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}
