default:
  @dev

dev:
  just backend

backend:
  cd backend && air

frontend:
  cd frontend && bun watch

build-backend:
  cd backend && go build -o bin/server ./cmd/server