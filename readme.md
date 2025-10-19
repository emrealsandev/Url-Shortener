### 🔗 URL Shortener (Go + Fiber, MongoDB, Redis)
A simple, fast, and production-ready URL shortening service written in Go. It uses MongoDB for persistence, Redis for aggressive caching, and Fiber for HTTP. It supports custom aliases, automatic Base62 code generation, per-request settings, and sensible rate limits.

Demo: https://shortener.emrealsan.com/

---

### 📦 Features
- Base62 short code generation from a monotonic sequence (XOR with `SEQUENCE_SALT`)
- Custom alias support on creation
- MongoDB persistence with sequence and URL storage
- Redis caching for hot paths (code→URL and URL→code) with configurable TTL
- Per-request dynamic settings loaded via provider (cached in Redis hash)
- Expirable short URLs (TTL in hours)
- Rate limiting
    - API: 20 requests/minute per IP
    - Redirects: 5 requests/second per IP
- Health and readiness endpoints
- Docker-based local setup (MongoDB, Redis); optional Air for hot reload

---

### 🚀 Getting Started
#### 1) Clone the repository
```bash
git clone https://github.com/emrealsandev/short-url.git
cd short-url
```

#### 2) Configure environment
Copy `.env.example` to `.env` and adjust values:
```bash
cp .env.example .env
```

Copy `.Dockerfile_local` to `.Dockerfile` for local setup:
```bash
cp .Dockerfile_local .Dockerfile
```

Minimal `.env` example:
```env
APP_ENVIRONMENT=local
PORT=8080
BASE_URL=http://localhost:8080

MONGO_URI=mongodb://host.docker.internal:27017
MONGO_DB=shortener

# Salt used when generating Base62 codes (use your own!)
SEQUENCE_SALT=0xDE_AD_BE_EF

# Redis configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3) Start the stack (Docker Compose)
- Start all services (app, MongoDB, Redis):
  ```bash
  docker compose up -d
  ```
- Tail application logs (includes Air hot-reload output):
  ```bash
  docker compose logs -f app
  ```
- Hot reload: Air starts automatically inside the `app` container (`command: ["air", "-c", "/app/.air.toml"]`). No manual Air command needed.
- App URL: `http://localhost:8080`
- Default ports: MongoDB `27017`, Redis `6379`
- Rebuild after changes to `Dockerfile` or dependencies:
  ```bash
  docker compose up -d --build
  ```

#### 4) Run migrations (if present)
Some setups include a small migration/initialization step (indexes, seed). If your project contains `cmd/migration`, run:
```bash
go run cmd/migration
```

#### 5) Run the API
- With Go directly:
```bash
go build -o main ./cmd/api
./main
```

- With Air (hot reload) if configured:
```bash
# Air typically reads .air.toml or air.toml
# Example check (in containerized dev):
docker compose exec app sh -lc 'ls -l /app/.air.toml || ls -l /app/air.toml'
```

---

### 🧪 API Reference
Base path: `http://localhost:8080`

- Health checks
    - `GET /healthz` → `200 ok`
    - `GET /readyz` → `200 ready`

- Create short URL
    - `POST /v1/shorten`
    - Request body:
      ```json
      {
        "url": "https://your-long-url.com/with/path?utm=x",
        "custom_alias": "optional-custom"  // optional
      }
      ```
    - Responses:
        - `200`:
          ```json
          {
            "code": "abc123",
            "short_url": "http://localhost:8080/abc123"
          }
          ```
        - `400` with body `invalid_url`
        - `409` with body `conflict` (custom alias taken or duplicate insert)
        - `500` with body `internal`

- Redirect
    - `GET /:code` → `302 Found` to original URL
    - Errors:
        - `404` when not found/disabled/expired

---

### ⚙️ Settings and Behavior
Settings are fetched per request and cached in Redis for 5 minutes by default. Two key settings influence behavior:
- `TtlTime` (hours): If greater than 0, newly created short URLs get an `ExpiresAt` of now + `TtlTime` hours.
- `RedisTtlTime` (minutes): Cache TTL for both mappings
    - `c:<code>` → URL
    - `u:<normalized_url>` → code
      Default is 5 minutes if not set.

Note: Settings are stored/retrieved via the repository and cached as a Redis hash. The settings provider key defaults to the code path; ensure consistent keying in your environment.

---

### 🔐 Rate Limiting
- API (`/v1/...`): 20 requests per minute per IP
- Redirects (`/:code`): 5 requests per second per IP
  These are enforced via Fiber’s `limiter` middleware.

---

### 🧰 Architecture Overview
- `cmd/api`: Server bootstrap (Fiber)
- `internal/short`: Core shortening logic
- `internal/repo`: Persistence models and repository (MongoDB)
- `internal/cache`: Redis client and helpers
- `internal/server`:
    - `handlers`: HTTP handlers for shorten/redirect
    - `middleware`: settings injection and rate limiters
    - `routes.go`: endpoint registration
- `internal/config`: env config loader and settings provider
- `pkg/base62`: Base62 encoder

---

### 📄 Configuration Reference (.env)
- `APP_ENVIRONMENT` (default: `dev`)
- `PORT` (default: `8080`)
- `BASE_URL` (required): e.g., `https://sho.rt`
- `MONGO_URI` (required)
- `MONGO_DB` (default: `shortener`)
- `SEQUENCE_SALT` (default: `_`): underscore-separated hex allowed, ex: `0xDE_AD_BE_EF`
- `REDIS_ADDR` (default: `localhost:6379`)
- `REDIS_PASSWORD` (default: empty)
- `REDIS_DB` (default: 0)

Security tips:
- Use a strong, secret `SEQUENCE_SALT`
- Put the service behind HTTPS (TLS termination via reverse proxy)
- Consider enabling auth/rate limiting at the edge in production

---

### 📦 Requirements
- Go 1.20+
- MongoDB 6/7+
- Redis 6/7+
- Docker (optional, for local setup)
- Air (optional, for hot reload)

---

### 🗺️ Roadmap Ideas
- Admin UI for managing links (enable/disable, manual expiry)
- Click analytics and unique visitor metrics
- API authentication (tokens)
- Batch shortening
- QR code generation
---

### 🤝 Contributing
PRs and issues are welcome. Please open an issue for discussion before large changes.

---

### 📄 License
MIT License. See `LICENSE`.

---

### ✨ Author
Made with ❤️ by [Emre Alsan](https://github.com/emrealsandev)
