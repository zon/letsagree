default:
  @dev

dev:
  just backend

backend:
  cd backend && air

frontend:
  cd frontend && bun watch

install:
  cd agreectl && go install ./cmd/agreectl

build-backend:
  cd backend && go build -o bin/server ./cmd/server

version := `cat ./VERSION`

push-backend:
  podman build --build-arg VERSION={{version}} -t ghcr.io/zon/letsagree-backend:{{version}} -t ghcr.io/zon/letsagree-backend:latest backend/
  podman push ghcr.io/zon/letsagree-backend:{{version}}
  podman push ghcr.io/zon/letsagree-backend:latest

push-frontend:
  podman build --build-arg VERSION={{version}} -t ghcr.io/zon/letsagree-frontend:{{version}} -t ghcr.io/zon/letsagree-frontend:latest frontend/
  podman push ghcr.io/zon/letsagree-frontend:{{version}}
  podman push ghcr.io/zon/letsagree-frontend:latest