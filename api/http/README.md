# GoKit HTTP API

The HTTP adapter exposes the unified interpretation pipeline as JSON.

## Start

```bash
go run ./cmd/gokit --serve --addr :8080
```

## Endpoints

### Health

`GET /health`

```json
{"status":"ok"}
```

### Interpret

`POST /v1/interpret`

Request:

```json
{"input":"what is 20% of 500"}
```

Response includes the original input, selected source and kind, typed/executed value, confidence, assumptions, and execution plan.

The adapter intentionally keeps transport concerns outside the interpretation packages, so applications can use the Go APIs directly without running an HTTP server.
