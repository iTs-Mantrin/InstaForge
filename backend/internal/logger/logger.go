package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"instaforge/internal/config"

	"github.com/rs/zerolog"
)

var (
	log  zerolog.Logger
	once sync.Once
)

// Get returns the global logger instance.
func Get() *zerolog.Logger {
	if log.GetLevel() == zerolog.NoLevel {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
		return &l
	}
	return &log
}

// Init initializes the structured JSON logger.
func Init(cfg config.AppConfig) {
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			level = zerolog.InfoLevel
		}

		zerolog.SetGlobalLevel(level)

		// Default: JSON output to stdout
		output := io.Writer(os.Stdout)

		// Pretty console in development
		if cfg.Env == "development" {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
		}

		log = zerolog.New(output).
			With().
			Timestamp().
			Str("service", cfg.Name).
			Str("environment", cfg.Env).
			Str("version", cfg.Version).
			Logger()
	})
}

// SetLevel updates the log level at runtime.
func SetLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return
	}
	zerolog.SetGlobalLevel(lvl)
	log = log.Level(lvl)
}

// WithRequest returns a logger enriched with request context.
func WithRequest(method, path, requestID, ip string) zerolog.Logger {
	return Get().With().
		Str("method", method).
		Str("path", path).
		Str("request_id", requestID).
		Str("ip", ip).
		Logger()
}


// Fatal logs a fatal error and exits.
func Fatal(err error, msg string) {
	Get().Fatal().Err(err).Msg(msg)
}
