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

push-backend:
  VERSION=$$(cat ./VERSION) && \
  docker build -t ghcr.io/zon/letsagree-backend:$$VERSION -t ghcr.io/zon/letsagree-backend:latest backend/ && \
  docker push ghcr.io/zon/letsagree-backend:$$VERSION && \
  docker push ghcr.io/zon/letsagree-backend:latest