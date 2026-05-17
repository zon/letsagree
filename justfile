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

version := `cat ./VERSION`

push-backend:
  podman build -t ghcr.io/zon/letsagree-backend:{{version}} -t ghcr.io/zon/letsagree-backend:latest backend/
  podman push ghcr.io/zon/letsagree-backend:{{version}}
  podman push ghcr.io/zon/letsagree-backend:latest

push-frontend:
  podman build -t ghcr.io/zon/letsagree-frontend:{{version}} -t ghcr.io/zon/letsagree-frontend:latest frontend/
  podman push ghcr.io/zon/letsagree-frontend:{{version}}
  podman push ghcr.io/zon/letsagree-frontend:latest