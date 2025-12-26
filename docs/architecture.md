# Architecture

Sequence Insights is a small, production-shaped service that demonstrates:

- Clean separation: handlers → services → storage
- PostgreSQL persistence
- Structured logging
- Basic auth via static bearer token (optional)
- Background worker for async post-processing

## Layout

- `server/cmd/api` – entrypoint
- `server/internal/handlers` – HTTP routes/handlers
- `server/internal/services` – domain logic (metrics)
- `server/internal/storage` – persistence implementations
- `server/internal/jobs` – in-process worker queue
- `server/migrations` – SQL schema used by docker-compose

## Local dev

1. `cp .env.example .env`
2. `make docker-up`
3. `curl http://localhost:8080/health`
