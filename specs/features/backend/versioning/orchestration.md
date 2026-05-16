# Version Endpoint

## Purpose

Expose the build-time embedded version as a bare JSON string via `GET /version`.

## Orchestration

**Module:** `internal/version`

```go
var Version string

func HandleGetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, Version)
}
```

- `Version` — package-level variable injected at compile time via `-ldflags`

## Tests

**Module:** `internal/version_test`

- **returns embedded version** — set `Version` to `"1.2.3"`; `GET /version` responds `200` with body `"1.2.3"`
- **uses JSON content type** — `GET /version` response has `Content-Type: application/json`
