# InstaForge Backend

High-throughput Instagram media download API with rate limiting, API key authentication, analytics aggregation, OpenTelemetry tracing, and S3-compatible storage.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Runtime** | Go 1.21+ |
| **Framework** | [Fiber v2](https://gofiber.io/) |
| **Database** | PostgreSQL (via [GORM](https://gorm.io/)) |
| **Cache/Queue** | Redis |
| **Message Queue** | RabbitMQ (optional, Redis fallback) |
| **Authentication** | JWT + API Key |
| **Metrics** | Prometheus + [Grafana](deploy/dashboards/instaforge.json) |
| **Tracing** | OpenTelemetry (OTLP HTTP) |
| **Storage** | Local filesystem / S3-compatible (AWS SDK v2) |
| **Deployment** | Docker / Kubernetes |

## Architecture

```
                        ┌──────────────┐
                        │   Clients     │
                        │ (API Keys or │
                        │   JWT User)  │
                        └──────┬───────┘
                               │
                        ┌──────▼───────┐
                        │  Fiber HTTP  │
                        │   Server     │
                        │  (cmd/main)  │
                        └──────┬───────┘
                               │
                  ┌────────────┼────────────┐
                  │            │            │
           ┌──────▼─────┐ ┌───▼────┐ ┌────▼──────┐
           │ Middleware │ │Routes  │ │ InstaDownloader │
           │ - RateLimit│ │ /api/  │ │ lib (Rapid)│
           │ - Auth     │ │   v1/  │ └───────────┘
           │ - CORS     │ └───┬────┘
           │ - RequestID│     │
           │ - Security │  ┌──▼───────────┐
           └────────────┘  │  Controllers │
                           └──┬───────────┘
                              │
                    ┌─────────┼──────────┐
                    │         │          │
             ┌──────▼──┐ ┌───▼───┐ ┌────▼─────┐
             │Services │ │Queue  │ │ Repos    │
             │(IG,User)│ │Worker │ │ (GORM)   │
             └─────────┘ └───────┘ └────┬─────┘
                                        │
                                  ┌─────▼──────┐
                                  │ PostgreSQL │
                                  └────────────┘

  Background goroutines (started in main.go):
    ├── AnalyticsAggregator (flushes AnalyticsEvents → AnalyticsSummaries every 5m)
    ├── Telemetry TracerProvider (OTLP, conditional)
    └── Storage Provider (minio-based health checks)
```

## Features

- **Instagram Media Download** — single posts, carousels, reels, stories, user feeds via public API (no login)
- **API Key & JWT Authentication** — dual auth: API keys for machine-to-machine, JWT for web users
- **Rate Limiting** — per-endpoint, per-IP, configurable burst/refill; fails closed on Redis outage (503)
- **Request Tracing** — OpenTelemetry OTLP HTTP exporter for distributed tracing
- **Analytics Aggregation** — batch rollup of request events into timestamped summaries
- **Prometheus Metrics** — request count, duration histogram, queue depth, worker count, active downloads
- **Grafana Dashboard** — pre-built 10-panel dashboard
- **Health Checks** — liveness, readiness with database/cache/queue probes
- **S3/Compatible Storage** — pluggable storage provider (local filesystem or S3 via AWS SDK v2)
- **Graceful Shutdown** — signal handling, worker drain, tracer flush, DB pool close

## API Endpoints

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health/live` | Liveness probe (always 200) |
| `GET` | `/health/ready` | Readiness probe (DB, cache, queue checks) |
| `GET` | `/metrics` | Prometheus metrics (pull-based) |
| `POST` | `/api/v1/preview` | Preview an Instagram URL |
| `GET` | `/api/v1/user/:username` | Search Instagram user |
| `GET` | `/api/v1/user/:username/feed` | Get user's media feed |
| `GET` | `/api/v1/user/:username/stories` | Get user's stories |
| `POST` | `/api/v1/download` | Queue a download |
| `GET` | `/api/v1/download/:id` | Get download status |
| `GET` | `/api/v1/download/stream` | Stream downloaded media |
| `GET` | `/api/v1/ping` | Ping (rate-limited) |

### Authenticated Endpoints (JWT)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/auth/profile` | Get current user profile |

## Environment Configuration

Config is loaded from environment variables (see `internal/config/config.go`). Key groups:

### Server
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `ENVIRONMENT` | `development` | Runtime environment |

### Database
| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | — | PostgreSQL connection string |
| `DATABASE_MAX_OPEN_CONNS` | `25` | Max open connections |
| `DATABASE_MAX_IDLE_CONNS` | `10` | Max idle connections |

### Redis
| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | — | Redis password |

### Rate Limiter
| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_REQUESTS` | `100` | Requests per window |
| `RATE_LIMIT_WINDOW` | `1m` | Rate limit window duration |
| `RATE_LIMIT_BURST` | `20` | Burst capacity |

### Storage
| Variable | Default | Description |
|----------|---------|-------------|
| `STORAGE_PROVIDER` | `local` | Storage provider (`local` or `s3`) |
| `STORAGE_TEMP_DIR` | `/tmp/instaforge` | Temp directory for local provider |
| `S3_ACCESS_KEY` | — | S3 access key |
| `S3_SECRET_KEY` | — | S3 secret key |
| `S3_BUCKET` | — | S3 bucket name |
| `S3_REGION` | — | S3 region |
| `S3_ENDPOINT` | — | S3 endpoint URL (for minio/compatible) |

### Telemetry
| Variable | Default | Description |
|----------|---------|-------------|
| `OTEL_ENABLED` | `false` | Enable OpenTelemetry |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://otel-collector:4318` | OTLP HTTP endpoint |

### JWT
| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | — | JWT signing secret (min 32 chars) |
| `JWT_ACCESS_EXPIRY` | `15m` | Access token TTL |

### Queue
| Variable | Default | Description |
|----------|---------|-------------|
| `QUEUE_DRIVER` | `redis` | Queue backend (`redis` or `rabbitmq`) |
| `RABBITMQ_URL` | — | RabbitMQ connection URL |

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis
- (Optional) RabbitMQ
- (Optional) Docker & Docker Compose

### Local Development

```bash
# Clone and enter backend
cd backend

# Copy environment template
cp .env.example .env

# Install dependencies
go mod tidy

# Run database migrations
go run cmd/migrate/main.go

# Start the server
go run cmd/server/main.go
```

### Docker Compose

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: instaforge
      POSTGRES_USER: instaforge
      POSTGRES_PASSWORD: instaforge
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://instaforge:instaforge@postgres:5432/instaforge?sslmode=disable
      REDIS_URL: redis:6379
```

### Kubernetes

See [kubernetes/](kubernetes/) for manifests. Deploy:

```bash
kubectl apply -f kubernetes/
```

## Project Structure

```
backend/
├── cmd/
│   ├── migrate/          # Database migration runner
│   └── server/           # Application entrypoint
│       └── main.go
├── deploy/
│   └── dashboards/
│       └── instaforge.json   # Grafana dashboard
├── kubernetes/
│   └── *.yaml               # K8s manifests
├── internal/
│   ├── analytics/
│   │   ├── aggregator.go     # Batch event → summary aggregation
│   │   ├── models.go         # AnalyticsEvent, AnalyticsSummary
│   │   └── repository.go     # GORM persistence
│   ├── api/
│   │   ├── controllers/      # HTTP handlers
│   │   └── router.go         # Route setup
│   ├── cache/                # Redis caching layer
│   ├── config/               # Environment config loader
│   ├── database/             # PostgreSQL connection & migration
│   ├── dto/                  # Shared data transfer objects
│   ├── events/               # Event bus
│   ├── health/               # Liveness & readiness checks
│   ├── logger/               # Structured logging (zerolog)
│   ├── metrics/              # Prometheus metric definitions
│   ├── middleware/           # RequestID, auth, rate limit, logging, CORS, security
│   ├── models/               # GORM models (User, APIKey, Media, Download, etc.)
│   ├── queue/                # Asynchronous job queue
│   ├── repositories/         # Data access layer (GORM)
│   ├── security/             # JWT, API key generator, password hashing
│   ├── services/             # Business logic (Instagram, User)
│   ├── storage/              # Pluggable storage provider (local + S3)
│   ├── telemetry/            # OpenTelemetry trace provider
│   ├── validators/           # Request validation
│   └── workers/              # Background job workers
├── pkg/
│   └── response/             # Standardized API response helpers
├── tests/                    # Integration tests
├── .env.example              # Environment template
├── Dockerfile
├── go.mod
└── README.md
```

## Middleware Chain

```
Request → RequestID → RequestLogger → SecurityHeaders → CORS → RateLimiter ↴
                                                                      ↳ Auth (JWT/API Key on protected routes)
```

All middleware is applied at the `/api/v1` group level. Health and metrics endpoints live outside the group and only get `RequestID`.

## Monitoring

### Prometheus Metrics

Available at `GET /metrics`:

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total requests by method, path, status |
| `http_request_duration_ms` | Histogram | Request duration in ms |
| `worker_queue_depth` | Gauge | Current queue length |
| `worker_active_jobs` | Gauge | Currently processing jobs |

### Grafana Dashboard

Import `deploy/dashboards/instaforge.json` into Grafana for 10 pre-built panels covering:

- Request rate & latency (p50, p95, p99)
- Active workers & queue depth
- Error rate by status code
- Download throughput
- Redis & DB health
- Aggregated analytics

### OpenTelemetry Tracing

Enable with `OTEL_ENABLED=true` and an OTLP HTTP collector endpoint. The tracer provider exports:

- HTTP request spans (method, path, status)
- Database query spans (via GORM hooks, future)
- Worker job spans (future)

## Security

### Authentication

Two authentication mechanisms:

1. **JWT (`/api/v1/auth/*`)** — Bearer token in `Authorization: Bearer <token>`. Short-lived (configurable via `JWT_ACCESS_EXPIRY`).
2. **API Key (`X-API-Key` header)** — Pre-generated keys with format `if_<hex>`. Middleware hashes the raw key with SHA-256, looks up the hash in the database, checks expiry, and updates `last_used` asynchronously.

### Rate Limiting

- Token bucket per endpoint per client IP
- Configurable refill rate and burst capacity
- **Fails closed**: Returns `503` if Redis is unavailable (prevents downstream overload)

### HTTP Security

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`
- CORS configurable per origin

## Testing

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Integration tests (requires PostgreSQL + Redis)
go test ./tests/...
```

### Known Gaps & Future Improvements

- **Worker tests are in progress** — test files exist but may need integration fixtures
- **Database migrations** — auto-migrated via GORM (`AutoMigrate`); consider a versioned migration tool for production
- **CORS** — currently permissive; tighten origin restrictions in production
- **Log retention** — no log rotation or remote shipping built-in
- **API key revocation** — implemented (status field), but no admin endpoint to manage keys yet
- **Media retention** — no TTL-based cleanup of downloaded media
- **mTLS / network policies** — not configured; add in multi-tenant deployments
- **Load testing** — no formal benchmarks; recommended before production traffic
