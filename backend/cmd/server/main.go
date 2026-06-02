package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"instaforge/internal/analytics"
	"instaforge/internal/api"
	"instaforge/internal/api/controllers"
	"instaforge/internal/cache"
	"instaforge/internal/config"
	"instaforge/internal/database"
	"instaforge/internal/health"
	"instaforge/internal/logger"
	"instaforge/internal/middleware"
	"instaforge/internal/queue"
	"instaforge/internal/security"
	"instaforge/internal/services"
	"instaforge/internal/storage"
	"instaforge/internal/telemetry"
	"instaforge/internal/workers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize structured logger
	logger.Init(cfg.App)
	log := logger.Get()
	log.Info().
		Str("name", cfg.App.Name).
		Str("env", cfg.App.Env).
		Int("port", cfg.App.Port).
		Msg("starting instaforge backend")

	// Connect to PostgreSQL (primary)
	db, err := database.RetryableConnect(cfg.DB, 5)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	log.Info().Msg("postgresql connected")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("failed to run database migrations")
	}
	log.Info().Msg("database migrations completed")

	// Connect to PostgreSQL replica (optional)
	if cfg.DBReplica.Host != "" && cfg.DBReplica.Host != cfg.DB.Host {
		if _, err := database.ConnectReplica(cfg.DBReplica); err != nil {
			log.Warn().Err(err).Msg("failed to connect to replica, using primary for reads")
		}
	}

	// Connect to Redis
	redisConfig := cfg.Redis
	if len(redisConfig.Addrs) == 1 && redisConfig.Addrs[0] == "localhost:6379" {
		_, err = cache.ConnectSingle(redisConfig)
	} else {
		_, err = cache.Connect(redisConfig)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	log.Info().Msg("redis connected")

	// Connect to Kafka
	if _, err := queue.ConnectProducer(cfg.Kafka); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to kafka")
	}
	if _, err := queue.ConnectConsumer(cfg.Kafka); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to kafka consumer")
	}
	log.Info().Msg("kafka connected")

	// Initialize OpenTelemetry tracing
	var tracerProvider *telemetry.TracerProvider
	if cfg.OTEL.Enabled {
		tp, err := telemetry.InitTracer(context.Background(), cfg.OTEL.OTLPEndpoint, cfg.App.Name)
		if err != nil {
			log.Warn().Err(err).Msg("failed to initialize otel tracer, continuing without tracing")
		} else {
			tracerProvider = tp
		}
	}

	// Initialize storage provider
	storageProvider := storage.NewProvider(&cfg.Storage)
	_ = storageProvider // available for services that need file storage

	// Initialize services
	instagramProvider := services.NewThirdPartyInstagramProvider(cfg.Instagram)
	if cfg.App.Env == "development" {
		instagramProvider = nil // fall through to mock
	}

	var instagramService *services.InstagramService
	if instagramProvider != nil && cfg.Instagram.APIKey != "" {
		instagramService = services.NewInstagramService(instagramProvider, cfg)
	} else {
		log.Warn().Msg("no instagram API key configured, using mock provider")
		mockProvider := services.NewMockInstagramProvider()
		instagramService = services.NewInstagramService(mockProvider, cfg)
	}

	downloadService := services.NewDownloadService(cfg, instagramService)
	userService := services.NewUserService(cfg, instagramService)

	// Initialize JWT manager
	jwtManager := security.NewJWTManager(cfg.JWT)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(&cfg.RateLimit)

	// Initialize health checker
	checker := health.NewChecker()
	checker.RegisterDefaults()

	// Initialize controllers
	instagramCtrl := controllers.NewInstagramController(instagramService, userService)
	downloadCtrl := controllers.NewDownloadController(downloadService)
	healthCtrl := controllers.NewHealthController(checker)

	// Initialize worker pool
	downloadHandler := workers.DownloadWorkerHandler(downloadService, &cfg.Kafka)
	workerPool := workers.NewWorkerPool(cfg.Kafka, downloadHandler, 10)
	workerPool.Start()

	// Initialize cache warmer
	if cfg.Cache.WarmEnabled {
		cacheWarmService := services.NewCacheWarmService(cfg, instagramService)
		go cacheWarmService.Warmup()
	}

	// Start analytics aggregator (runs every 5 minutes)
	aggInterval := 5 * time.Minute
	aggregator := analytics.NewAggregator(aggInterval)
	aggregator.Start()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:       cfg.App.Name,
		Prefork:       cfg.App.Env == "production",
		ReadTimeout:   15 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   120 * time.Second,
		ReadBufferSize: 8192,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(favicon.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.CORSConfig())

	// Setup routes
	api.SetupRoutes(app, api.RouterConfig{
		InstagramCtrl: instagramCtrl,
		DownloadCtrl:  downloadCtrl,
		HealthCtrl:    healthCtrl,
		JWTManager:    jwtManager,
		RateLimiter:   rateLimiter,
		Checker:       checker,
	})

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		log.Info().Str("address", addr).Msg("http server listening")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down...")

	// Stop worker pool
	workerPool.Stop()

	// Stop analytics aggregator
	aggregator.Stop()

	// Shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("server shutdown failed")
	}

	// Shutdown OpenTelemetry tracer
	if tracerProvider != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("otel tracer shutdown error")
		}
		cancel()
	}

	// Close connections
	cache.Close()
	queue.Close()

	log.Info().Msg("server exited")
}

func init() {
	// Set timezone
	os.Setenv("TZ", "UTC")
}
