package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"instaforge/internal/config"
	"instaforge/internal/logger"
)

// Common errors.
var (
	ErrNotFound    = errors.New("object not found")
	ErrInvalidPath = errors.New("invalid object path")
)

// Provider defines the interface for media storage backends.
type Provider interface {
	// Upload stores a file and returns its storage key.
	Upload(ctx context.Context, key string, reader io.Reader) error

	// Download retrieves a stored file by key.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a stored file by key.
	Delete(ctx context.Context, key string) error

	// Exists checks whether a key exists in storage.
	Exists(ctx context.Context, key string) (bool, error)

	// URL returns a public (or presigned) URL for the key, or empty if not available.
	URL(ctx context.Context, key string) (string, error)
}

// NewProvider creates the appropriate storage provider based on config.
func NewProvider(cfg *config.StorageConfig) Provider {
	switch cfg.Provider {
	case "s3":
		p, err := newS3Provider(cfg)
		if err != nil {
			logger.Get().Warn().Err(err).Msg("failed to init S3 provider, falling back to local")
			return newLocalProvider(cfg)
		}
		logger.Get().Info().Str("bucket", cfg.S3Bucket).Str("region", cfg.S3Region).Msg("S3 storage provider initialized")
		return p
	default:
		logger.Get().Info().Str("dir", cfg.TempDir).Msg("local storage provider initialized")
		return newLocalProvider(cfg)
	}
}

// localProvider stores files on the local filesystem.
type localProvider struct {
	baseDir string
}

func newLocalProvider(cfg *config.StorageConfig) *localProvider {
	return &localProvider{baseDir: cfg.TempDir}
}

func (p *localProvider) Upload(_ context.Context, key string, reader io.Reader) error {
	fullPath := filepath.Join(p.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (p *localProvider) Download(_ context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(p.baseDir, key)
	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return f, nil
}

func (p *localProvider) Delete(_ context.Context, key string) error {
	fullPath := filepath.Join(p.baseDir, key)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (p *localProvider) Exists(_ context.Context, key string) (bool, error) {
	_, err := os.Stat(filepath.Join(p.baseDir, key))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *localProvider) URL(_ context.Context, key string) (string, error) {
	return "", nil // local provider has no public URL
}
