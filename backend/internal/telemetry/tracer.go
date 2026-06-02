package telemetry

import (
	"context"
	"fmt"
	"time"

	"instaforge/internal/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TracerProvider wraps the OpenTelemetry SDK tracer provider for lifecycle management.
type TracerProvider struct {
	provider *sdktrace.TracerProvider
}

// InitTracer creates and configures an OTLP HTTP trace exporter.
// endpoint should be the OTLP collector address without scheme (e.g. "otel-collector:4318").
func InitTracer(ctx context.Context, endpoint, serviceName string) (*TracerProvider, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create otlp trace exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.26.0",
		attribute.String("service.name", serviceName),
		attribute.String("service.version", "1.0.0"),
		attribute.String("telemetry.sdk.language", "go"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	logger.Get().Info().Str("endpoint", endpoint).Str("service", serviceName).Msg("otel tracer provider initialized")
	return &TracerProvider{provider: tp}, nil
}

// Shutdown gracefully flushes and shuts down the tracer provider.
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp == nil || tp.provider == nil {
		return nil
	}
	return tp.provider.Shutdown(ctx)
}
