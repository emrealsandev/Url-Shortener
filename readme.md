# 🔗 Short URL Service

A simple and performant URL shortening service written in Go, using MongoDB as the backend store.  
Includes Base62 encoding, unique sequence number generation, and supports Docker-based development.

---

## 📦 Features

- Shorten long URLs using Base62 encoding.
- Persistent sequence-based short link generation.
- MongoDB backend with indexes for performance.
- `.env`-based config loading (via `envconfig`).
- Docker-based development environment with `mongo` and `mongo-express`.
- Hot-reloading with [air](https://github.com/cosmtrek/air) for development.
- Lightweight and extendable architecture.

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/emrealsandev/short-url.git
cd short-url
```

---

### 2. Setup Environment Variables

Copy the provided `.env.example` file and customize as needed:

```bash
cp .env.example .env
```

---

### 3. Run MongoDB with Docker

A `docker-compose.yml` is provided for spinning up MongoDB and Mongo Express:

```bash
docker-compose up -d
```

> 📌 If you want a visual interface for Mongo, open `http://localhost:8081` (mongo-express).

---

### 4. Run DB Migration

This will initialize the MongoDB with required indexes or setup:

```bash
go run cmd/migration
```

---

### 5. Start the Server (with Air)

To enable hot-reloading for development, use `air` with the provided config:

```bash
air -c .air.toml
```

Or build and run manually:

```bash
go build -o main ./cmd/api
./main
```

---

## 🧪 Sample Request

```bash
POST /shorten
Content-Type: application/json

{
  "url": "https://your-long-url.com"
}
```

Response:

```json
{
  "short_url": "http://localhost:3000/abc123"
}
```



---

## ⚙️ Configuration

You can adjust server behavior using `.env` variables. Example:

```env
APP_ENVIRONMENT=local-prod
PORT=8080
BASE_URL=http://localhost:8080

MONGO_URI=mongodb://localhost:27017
MONGO_DB=shortener

SEQUENCE_SALT=EXAMPLE_SALT (hex format with underscore seperators)
```

Also see `.air.toml` for customizing Air's hot reload behavior.

---

## 🛠 Requirements

- Go 1.20+
- Docker + Docker Compose
- (Optional) [Air](https://github.com/cosmtrek/air) for hot reload

---

## 📌 TODO / Roadmap

- [ ] Rate limiting
- [ ] Expiry for short URLs
- [ ] Not found or expire page displaying
- [ ] QR code support

---

## 🤝 Contributing

PRs and issues are welcome! Feel free to fork and customize as needed.

---

## 📄 License

MIT License. See `LICENSE` file.

---

## ✨ Author

Made with ❤️ by [Emre Alsan](https://github.com/emrealsandev)
