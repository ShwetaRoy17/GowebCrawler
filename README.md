# GowebCrawler

A production-grade concurrent web crawler written in Go. Supports both a CLI mode for quick crawls and an HTTP server mode for integration with other services.

## Features

- Concurrent crawling with configurable goroutine pool
- Exponential backoff retry on network failures
- Token bucket rate limiting to avoid overloading servers
- robots.txt compliance with per-domain caching
- HTTP 429 / Retry-After handling
- Middleware chain — logging, metrics, gzip compression
- REST API for programmatic crawl control
- Postgres persistence for job history
- Docker and Kubernetes ready

## Project Structure

```
GowebCrawler/
├── cmd/
│   ├── crawler/        # CLI entrypoint
│   └── server/         # HTTP server entrypoint
├── internal/
│   ├── config/         # Viper-based configuration
│   ├── database/       # Postgres integration
│   ├── fetcher/        # HTTP client with middleware chain
│   ├── models/         # Shared data types
│   └── parser/         # HTML link extraction
├── deploy/
│   └── k8s/            # Kubernetes manifests
├── config.yaml         # Runtime configuration
├── docker-compose.yml  # Local development stack
└── Dockerfile          # Multi-stage build
```

## Prerequisites

- Go 1.25+
- Docker Desktop
- kubectl (for Kubernetes deployment)

## Quick Start

### CLI mode

Build and run a crawl directly from the command line:

```bash
make build
./bin/gowebcrawler start --seed https://golang.org --depth 2
```

Available flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--seed` | required | URL to start crawling from |
| `--depth` | 3 | How deep to follow links |
| `--concurrency` | 10 | Number of concurrent goroutines |

### Server mode

Start the HTTP server:

```bash
go run ./cmd/server/
```

The server starts on port 8080 and exposes three endpoints.

## REST API

### Start a crawl

```
POST /start
Content-Type: application/json

{
  "seed": "https://golang.org",
  "depth": 2,
  "concurrency": 5
}
```

Response:

```json
{
  "id": "job-1",
  "status": "in_progress",
  "pages": 0,
  "seed": "https://golang.org",
  "depth": 2,
  "concurrency": 5
}
```

### Check job status

```
GET /status?id=job-1
```

Response:

```json
{
  "id": "job-1",
  "status": "completed",
  "pages": 310,
  "seed": "https://golang.org",
  "depth": 2,
  "concurrency": 5
}
```

Job status values: `in_progress`, `completed`, `failed`

### List all jobs

```
GET /jobs
```

Response:

```json
[
  {
    "id": "job-1",
    "status": "completed",
    "pages": 310,
    "seed": "https://golang.org",
    "depth": 2,
    "concurrency": 5
  }
]
```

## Configuration

Configuration is loaded in priority order — flags override env vars, env vars override `config.yaml`, `config.yaml` overrides defaults.

### config.yaml

```yaml
max_depth: 2
concurrency: 6
user_agent: "goWebC/1.0 (+https://github.com/ShwetaRoy17/GowebCrawler)"
rate_limit: 5       # requests per second
burst: 10           # token bucket burst size
timeout: 30         # HTTP timeout in seconds
```

### Environment variables

All config values can be overridden with environment variables prefixed `WEBCRAWLER_`:

```bash
WEBCRAWLER_MAX_DEPTH=3
WEBCRAWLER_CONCURRENCY=10
WEBCRAWLER_RATE_LIMIT=5
```

### .env file (server mode)

Create a `.env` file in the project root for local development:

```
DATABASE_URL=postgres://postgres:secret@localhost:5432/crawler
```

## Local Development with Docker Compose

Start the full stack — crawler, Postgres, Redis, and Prometheus:

```bash
docker compose up
```

Services:

| Service | Port | Description |
|---------|------|-------------|
| postgres | 5432 | Job persistence |
| redis | 6379 | Available for caching |
| prometheus | 9090 | Metrics collection |

Stop everything:

```bash
docker compose down
```

Stop and delete all data:

```bash
docker compose down -v
```

## Docker

Build the image:

```bash
make docker-build
```

Run the crawler:

```bash
docker run go-crawler start --seed https://golang.org --depth 1
```

The image is ~22MB using a multi-stage build with a distroless base.

## Kubernetes

Apply all manifests:

```bash
kubectl apply -f deploy/k8s/
```

This creates:

- `Deployment` — 2 replicas of the server
- `Service` — ClusterIP on port 8080
- `ConfigMap` — non-sensitive configuration
- `Secret` — database URL (base64 encoded)
- `HorizontalPodAutoscaler` — scales 2 to 20 pods at 70% CPU

Check status:

```bash
kubectl get all
kubectl logs deployment/webcrawler
```

## Makefile targets

```bash
make build        # compile binary to bin/gowebcrawler
make run          # run with ARGS="start --seed <url>"
make test         # run all tests
make lint         # run golangci-lint
make docker-build # build Docker image
make clean        # remove build artifacts
```

## CI/CD

GitHub Actions runs on every push:

- **lint** — golangci-lint on all packages
- **test** — go test with race detector
- **docker** — builds and pushes to ghcr.io/shwetaroy17/gowebcrawler:latest (main branch only)

## Architecture

```
POST /start
     ↓
Server creates job → saves to Postgres → returns job ID
     ↓
Background goroutine starts crawl
     ↓
Fetcher pipeline:
  rate limiter → robots.txt check → HTTP request
  └── middleware chain: logging → metrics → compression → http.Client
     ↓
Parser extracts links → recursive crawl up to max_depth
     ↓
Job marked complete → Postgres updated
```