package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"instaforge/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// s3Provider stores files in an S3-compatible bucket.
type s3Provider struct {
	client *s3.Client
	bucket string
}

func newS3Provider(cfg *config.StorageConfig) (*s3Provider, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey, cfg.S3SecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	clientOpts := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = true // required for MinIO and most S3-compatible stores
		},
	}

	// Custom endpoint (for MinIO, DigitalOcean Spaces, etc.)
	if cfg.S3Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)

	return &s3Provider{
		client: client,
		bucket: cfg.S3Bucket,
	}, nil
}

func (p *s3Provider) Upload(ctx context.Context, key string, reader io.Reader) error {
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("s3 upload %s: %w", key, err)
	}
	return nil
}

func (p *s3Provider) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 download %s: %w", key, err)
	}
	return out.Body, nil
}

func (p *s3Provider) Delete(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete %s: %w", key, err)
	}
	return nil
}

func (p *s3Provider) Exists(ctx context.Context, key string) (bool, error) {
	_, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
			var nf *types.NotFound
		if errors.As(err, &nf) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *s3Provider) URL(ctx context.Context, key string) (string, error) {
	// Presigned URL valid for 1 hour
	presignClient := s3.NewPresignClient(p.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("s3 presign %s: %w", key, err)
	}
	return req.URL, nil
}
