# Backend Publish

## Purpose

The backend can be packaged as a container image and published to the GitHub Container Registry (ghcr.io) with a version tag sourced from the repo's `VERSION` file and a rolling `latest` tag.

## Requirements

- The repository MUST contain a `Containerfile` at `backend/Containerfile` that produces a runnable backend image.
- The `just push-backend` recipe MUST build the container image using the `Containerfile`.
- The `just push-backend` recipe MUST tag the image with the version string read from `./VERSION` (e.g. `ghcr.io/<owner>/<repo>-backend:0.3.0`).
- The `just push-backend` recipe MUST also tag the image as `latest` (e.g. `ghcr.io/<owner>/<repo>-backend:latest`).
- The `just push-backend` recipe MUST push both tags to ghcr.io.

## Scenarios

### Scenario: push-backend tags with VERSION

Given `./VERSION` contains `0.3.0`  
When `just push-backend` is run  
Then the image is pushed to `ghcr.io/<owner>/<repo>-backend:0.3.0`  
And the image is pushed to `ghcr.io/<owner>/<repo>-backend:latest`

### Scenario: container image runs the backend

Given the image is built from `backend/Containerfile`  
When a container is started from that image  
Then the backend HTTP server starts and responds to `GET /version`
