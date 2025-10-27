# URL Shortener

![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14-blue)
![Redis](https://img.shields.io/badge/Redis-7-red)

Production-ready URL shortening service built with Go, featuring PostgreSQL, Redis caching, Prometheus monitoring, and comprehensive testing.

---

## Features

- **URL Shortening** - Generate short URLs with customizable TTL expiration
- **JWT Authentication** - Secure user authentication with token-based access
- **Redis Caching** - High-performance caching layer for improved response times
- **Structured Logging** - JSON logs with rotation using lumberjack
- **Prometheus Metrics** - Built-in metrics collection and monitoring
- **Grafana Dashboards** - Pre-configured dashboards for visualization
- **Comprehensive Testing** - Unit, integration, and benchmark tests
- **Load Testing** - K6 load testing scenarios included
- **Clean Architecture** - SOLID principles and dependency injection
- **Docker Support** - Full Docker Compose setup with all services

---

## Quick Start

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- PostgreSQL 14 (or use Docker)
- Redis 7 (or use Docker)

### Installation

```bash
# Clone the repository
git clone https://github.com/MaxGot69/url-shortener.git
cd url-shortener

# Start all services with Docker
docker-compose up -d

# Run migrations (automatically via AutoMigrate)
# Build and run the application
go run cmd/server/main.go
```

---

## Configuration

Set environment variables (or use defaults):

```bash
export PORT=8081
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=maxim
export DB_PASSWORD=secret
export DB_NAME=urlshortener
export REDIS_ADDR=localhost:6379
export JWT_SECRET=your-secret-key
export LOG_LEVEL=info
```

---

## API Endpoints

### URL Shortening

**Create Short URL**
```http
POST /api/v1/shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

**Response:**
```json
{
  "short_url": "http://localhost:8081/abc123"
}
```

### URL Redirection

**Redirect to Original URL**
```http
GET /{shortCode}
```

Returns HTTP 302 redirect to original URL.

### Authentication

**Register User**
```http
POST /api/v1/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Login**
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Statistics

**Get URL Statistics** (requires authentication)
```http
GET /api/v1/stats/{shortCode}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com",
  "clicks": 42
}
```

### Monitoring

**Health Check**
```http
GET /health
```

**Prometheus Metrics**
```http
GET /metrics
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      HTTP Layer                         │
│  Handlers (shorten, redirect, auth, stats)              │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────────┐
│                   Middleware Layer                      │
│  JWT Auth, Request ID, CORS, Logging                    │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────────┐
│                    Service Layer                        │
│  URL Service (validation, generation, business logic)   │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────────┐
│                 Repository Layer                        │
│  PostgreSQL (persistent storage)                        │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────────┐
│                    Cache Layer                          │
│  Redis (hot data caching)                               │
└─────────────────────────────────────────────────────────┘
```

---

## Testing

### Unit Tests

```bash
go test ./internal/handler/... -v
```

### Integration Tests

```bash
go test ./tests/integration/... -v
```

### Benchmark Tests

```bash
go test -bench=. ./benchmark/ -benchmem
```

### Load Testing with K6

```bash
# Install K6
brew install k6  # macOS
# or
# https://k6.io/docs/getting-started/installation/

# Run load tests
k6 run load-test/k6-test.js
k6 run load-test/redirect-test.js
```

---

## Monitoring

### Prometheus

Prometheus is available at `http://localhost:9090`

Metrics collected:
- `url_shortens_total` - Total shortened URLs
- `url_redirects_total` - Total redirects
- `active_urls_count` - Active URLs count
- `http_requests_total` - HTTP request counter
- `http_request_duration_seconds` - Request duration histogram

### Grafana

Grafana is available at `http://localhost:3000`

Default credentials:
- Username: `admin`
- Password: `admin`

---

## Project Structure

```
url-shortener/
├── cmd/
│   └── server/           # Application entry point
│       └── main.go
├── internal/
│   ├── config/           # Configuration management
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # JWT auth, request ID
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── server/           # Router setup
│   └── service/          # Business logic
├── pkg/
│   ├── cache/            # Redis client
│   ├── database/         # Database connection
│   └── logger/           # Structured logging
├── tests/
│   └── integration/      # Integration tests
├── benchmark/            # Benchmark tests
├── load-test/            # K6 load tests
├── migrations/           # Database migrations
├── docker-compose.yml    # Docker services
└── prometheus.yml        # Prometheus config
```

---

## Development

### Requirements

- Go 1.25+
- Make (optional)

### Build

```bash
go build -o url-shortener cmd/server/main.go
```

### Run

```bash
./url-shortener
```

### Database Migrations

Migrations run automatically via GORM AutoMigrate. Manual migrations available in `migrations/` directory.

---

## Technologies

- **Go 1.25+** - Programming language
- **Chi** - HTTP router
- **GORM** - ORM library
- **PostgreSQL** - Database
- **Redis** - Caching
- **JWT** - Authentication
- **Prometheus** - Metrics
- **Grafana** - Visualization
- **K6** - Load testing
- **Docker** - Containerization

---

## License

MIT License - see LICENSE file for details

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## Author

MaxGot69