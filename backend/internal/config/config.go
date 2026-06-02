package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	App     AppConfig
	DB      DBConfig
	DBReplica DBConfig
	Redis   RedisConfig
	Kafka   KafkaConfig
	Instagram InstagramConfig
	JWT     JWTConfig
	RateLimit RateLimitConfig
	Cache   CacheConfig
	Cloudflare CloudflareConfig
	OTEL    OTELConfig
	Storage StorageConfig
	Features FeatureConfig
}

type AppConfig struct {
	Name      string
	Env       string
	Port      int
	LogLevel  string
	Version   string
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

type RedisConfig struct {
	Addrs         []string
	Password      string
	DB            int
	PoolSize      int
	MinIdleConns  int
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

type KafkaConfig struct {
	Brokers         []string
	ClientID        string
	ConsumerGroup   string
	DownloadTopic   string
	RetryTopic      string
	DLQTopic        string
}

type InstagramConfig struct {
	APIKey      string
	APIBaseURL  string
	RapidAPIKey string
	RapidAPIHost string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type CacheConfig struct {
	TTLSeconds     int
	UserTTLSeconds int
	WarmEnabled    bool
}

type CloudflareConfig struct {
	APIToken  string
	ZoneID    string
	AccountID string
}

type OTELConfig struct {
	Enabled bool
	OTLPEndpoint string
}

type StorageConfig struct {
	Provider  string
	TempDir   string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
	S3Endpoint  string
}

type FeatureConfig struct {
	CarouselEnabled     bool
	StoriesEnabled      bool
	UsernameSearchEnabled bool
}

var (
	cfg  *Config
	once sync.Once
)

// ResetForTesting clears the cached config for testing.
// Should only be used in test files.
func ResetForTesting() {
	cfg = nil
	once = sync.Once{}
}

// Load reads configuration from environment variables.
func Load() *Config {
	once.Do(func() {
		cfg = &Config{
			App: AppConfig{
				Name:     getEnv("APP_NAME", "instaforge"),
				Env:      getEnv("APP_ENV", "development"),
				Port:     getEnvInt("APP_PORT", 8080),
				LogLevel: getEnv("APP_LOG_LEVEL", "info"),
				Version:  "1.0.0",
			},
			DB: DBConfig{
				Host:            getEnv("DB_HOST", "localhost"),
				Port:            getEnvInt("DB_PORT", 5432),
				User:            getEnv("DB_USER", "instaforge"),
				Password:        getEnv("DB_PASSWORD", "changeme"),
				Name:            getEnv("DB_NAME", "instaforge"),
				SSLMode:         getEnv("DB_SSL_MODE", "disable"),
				MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 50),
				MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
				ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
			},
			DBReplica: DBConfig{
				Host:     getEnv("DB_REPLICA_HOST", "localhost"),
				Port:     getEnvInt("DB_REPLICA_PORT", 5432),
				User:     getEnv("DB_REPLICA_USER", "instaforge"),
				Password: getEnv("DB_REPLICA_PASSWORD", "changeme"),
				Name:     getEnv("DB_REPLICA_NAME", "instaforge"),
				SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			},
			Redis: RedisConfig{
				Addrs:         strings.Split(getEnv("REDIS_ADDRS", "localhost:6379"), ","),
				Password:      getEnv("REDIS_PASSWORD", ""),
				DB:            getEnvInt("REDIS_DB", 0),
				PoolSize:      getEnvInt("REDIS_POOL_SIZE", 50),
				MinIdleConns:  getEnvInt("REDIS_MIN_IDLE_CONNS", 10),
				DialTimeout:   time.Duration(getEnvInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
				ReadTimeout:   time.Duration(getEnvInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
				WriteTimeout:  time.Duration(getEnvInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
			},
			Kafka: KafkaConfig{
				Brokers:       strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
				ClientID:      getEnv("KAFKA_CLIENT_ID", "instaforge"),
				ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "instaforge-workers"),
				DownloadTopic: getEnv("KAFKA_DOWNLOAD_TOPIC", "download_requests"),
				RetryTopic:    getEnv("KAFKA_RETRY_TOPIC", "download_retry"),
				DLQTopic:      getEnv("KAFKA_DLQ_TOPIC", "download_dlq"),
			},
			Instagram: InstagramConfig{
				APIKey:      getEnv("INSTAGRAM_API_KEY", ""),
				APIBaseURL:  getEnv("INSTAGRAM_API_BASE_URL", ""),
				RapidAPIKey: getEnv("INSTAGRAM_RAPIDAPI_KEY", ""),
				RapidAPIHost: getEnv("INSTAGRAM_RAPIDAPI_HOST", ""),
			},
			JWT: JWTConfig{
				Secret:        getEnv("JWT_SECRET", "changeme"),
				AccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
				RefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "720h")),
			},
			RateLimit: RateLimitConfig{
				Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
				Window:   time.Duration(getEnvInt("RATE_LIMIT_WINDOW", 60)) * time.Second,
			},
			Cache: CacheConfig{
				TTLSeconds:     getEnvInt("CACHE_TTL_SECONDS", 300),
				UserTTLSeconds: getEnvInt("CACHE_USER_TTL_SECONDS", 600),
				WarmEnabled:    getEnvBool("CACHE_WARM_ENABLED", false),
			},
			Cloudflare: CloudflareConfig{
				APIToken:  getEnv("CF_API_TOKEN", ""),
				ZoneID:    getEnv("CF_ZONE_ID", ""),
				AccountID: getEnv("CF_ACCOUNT_ID", ""),
			},
			OTEL: OTELConfig{
				Enabled:      getEnvBool("OTEL_ENABLED", false),
				OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4318"),
			},
			Storage: StorageConfig{
				Provider:    getEnv("STORAGE_PROVIDER", "local"),
				TempDir:     getEnv("STORAGE_TEMP_DIR", "/tmp/instaforge"),
				S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
				S3SecretKey: getEnv("S3_SECRET_KEY", ""),
				S3Bucket:    getEnv("S3_BUCKET", ""),
				S3Region:    getEnv("S3_REGION", ""),
				S3Endpoint:  getEnv("S3_ENDPOINT", ""),
			},
			Features: FeatureConfig{
				CarouselEnabled:      getEnvBool("FEATURE_CAROUSEL_ENABLED", true),
				StoriesEnabled:       getEnvBool("FEATURE_STORIES_ENABLED", true),
				UsernameSearchEnabled: getEnvBool("FEATURE_USERNAME_SEARCH_ENABLED", true),
			},
		}

		// Validation
		if cfg.App.Env == "production" && cfg.JWT.Secret == "changeme" {
			panic("JWT_SECRET must be changed in production")
		}
	})

	return cfg
}

// Get returns the singleton config instance.
func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
