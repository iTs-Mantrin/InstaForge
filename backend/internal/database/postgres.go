package database

import (
	"fmt"
	"sync"
	"time"

	"instaforge/internal/config"
	"instaforge/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	replica *gorm.DB
	once   sync.Once
	rOnce  sync.Once
)

// Connect initializes the primary PostgreSQL connection.
func Connect(cfg config.DBConfig) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		log := logger.Get()
		log.Info().
			Str("host", cfg.Host).
			Int("port", cfg.Port).
			Str("database", cfg.Name).
			Msg("connecting to primary database")

		gormCfg := &gorm.Config{
			Logger:                 gormlogger.Default.LogMode(gormlogger.Warn),
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		}

		db, err = gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
		if err != nil {
			err = fmt.Errorf("failed to connect to primary database: %w", err)
			return
		}

		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			err = fmt.Errorf("failed to get underlying sql.DB: %w", sqlErr)
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		// Verify connection
		if pingErr := sqlDB.Ping(); pingErr != nil {
			err = fmt.Errorf("failed to ping primary database: %w", pingErr)
			return
		}

		log.Info().Msg("primary database connected successfully")
	})

	return db, err
}

// ConnectReplica initializes the read-only replica connection.
func ConnectReplica(cfg config.DBConfig) (*gorm.DB, error) {
	var err error
	rOnce.Do(func() {
		log := logger.Get()
		log.Info().
			Str("host", cfg.Host).
			Int("port", cfg.Port).
			Str("database", cfg.Name).
			Msg("connecting to replica database")

		gormCfg := &gorm.Config{
			Logger:                 gormlogger.Default.LogMode(gormlogger.Warn),
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		}

		replica, err = gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
		if err != nil {
			err = fmt.Errorf("failed to connect to replica database: %w", err)
			return
		}

		sqlDB, sqlErr := replica.DB()
		if sqlErr != nil {
			err = fmt.Errorf("failed to get replica sql.DB: %w", sqlErr)
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		log.Info().Msg("replica database connected successfully")
	})

	return replica, err
}

// Get returns the primary database connection.
func Get() *gorm.DB {
	if db == nil {
		panic("database not initialized — call database.Connect() first")
	}
	return db
}

// GetReplica returns the read-replica connection (falls back to primary if not configured).
func GetReplica() *gorm.DB {
	if replica != nil {
		return replica
	}
	return Get()
}

// WithTransaction executes a function within a database transaction.
func WithTransaction(fn func(tx *gorm.DB) error) error {
	return Get().Transaction(fn)
}

// HealthCheck verifies database connectivity.
func HealthCheck() error {
	sqlDB, err := Get().DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// RetryableConnect attempts to connect with exponential backoff.
func RetryableConnect(cfg config.DBConfig, maxRetries int) (*gorm.DB, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		conn, err := Connect(cfg)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		backoff := time.Duration(1<<uint(i)) * time.Second
		logger.Get().Warn().
			Int("attempt", i+1).
			Int("max_retries", maxRetries).
			Dur("backoff", backoff).
			Err(err).
			Msg("database connection retry")
		time.Sleep(backoff)
	}
	return nil, fmt.Errorf("database connection failed after %d retries: %w", maxRetries, lastErr)
}
