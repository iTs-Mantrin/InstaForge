package tests

import (
	"os"
	"testing"

	"instaforge/internal/config"
)

func init() {
	// Reset config before any test runs so first Load() picks up test env vars
	config.ResetForTesting()
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg := config.Load()

	if cfg.App.Name != "instaforge" {
		t.Errorf("expected App.Name = instaforge, got %s", cfg.App.Name)
	}
	if cfg.App.Port != 8080 {
		t.Errorf("expected App.Port = 8080, got %d", cfg.App.Port)
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("expected DB.Host = localhost, got %s", cfg.DB.Host)
	}
	if cfg.DB.MaxOpenConns != 50 {
		t.Errorf("expected DB.MaxOpenConns = 50, got %d", cfg.DB.MaxOpenConns)
	}
}

func TestLoadConfig_EnvironmentOverrides(t *testing.T) {
	config.ResetForTesting()
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DB_HOST", "prod-db.example.com")
	os.Setenv("CACHE_TTL_SECONDS", "600")
	os.Setenv("FEATURE_STORIES_ENABLED", "false")

	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("CACHE_TTL_SECONDS")
		os.Unsetenv("FEATURE_STORIES_ENABLED")
	}()

	// Force reload
	cfg := config.Load()

	if cfg.App.Port != 9090 {
		t.Errorf("expected App.Port = 9090, got %d", cfg.App.Port)
	}
	if cfg.DB.Host != "prod-db.example.com" {
		t.Errorf("expected DB.Host = prod-db.example.com, got %s", cfg.DB.Host)
	}
	if cfg.Cache.TTLSeconds != 600 {
		t.Errorf("expected Cache.TTLSeconds = 600, got %d", cfg.Cache.TTLSeconds)
	}
	if cfg.Features.StoriesEnabled != false {
		t.Errorf("expected Features.StoriesEnabled = false, got %v", cfg.Features.StoriesEnabled)
	}
}

func TestDBConfig_DSN(t *testing.T) {
	config.ResetForTesting()
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "disable")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSL_MODE")
	}()

	cfg := config.Load()
	dsn := cfg.DB.DSN()

	expected := "host=testhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	if dsn != expected {
		t.Errorf("expected DSN = %q, got %q", expected, dsn)
	}
}
