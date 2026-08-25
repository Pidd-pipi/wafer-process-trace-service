# Wafer Process Trace Service

面向半导体工艺批次的追踪服务，默认监听 `8080`，可由 `PORT` 环境变量调整。

## Layout

```text
.
├── backend/                 # Go module, application, tests, static assets, Dockerfile
├── database/                # persistence extension point
├── output/                  # verification record
├── prompt.txt               # task prompt
├── runtime_smoke.json       # startup contract
├── .env.example
└── .gitignore
```

## API and Health

- `GET /healthz` - health check
- `GET /api/lots` - list process lots
- `POST /api/lots/status` - update a lot status
- `GET /` and `GET /app.js` - embedded trace page assets

Supported statuses are `queued`, `running`, `hold`, and `completed`.

## Run

```bash
cd backend
go test ./...
go build ./...
go run .
```
