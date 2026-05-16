default:
  @dev

dev:
  just backend

backend:
  air

build-backend:
  go build -o bin/server ./cmd/server