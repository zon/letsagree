# Backend Bootstrap

## Purpose

The backend is built with a version sourced from a `VERSION` file at the repo root, exposes that version via a `/version` endpoint, and is hosted locally via `air` for live-reload development.

## Requirements

- The repository MUST contain a `VERSION` file at the repo root containing a valid [semver](https://semver.org) string (e.g. `1.0.0`).
- The backend MUST expose a `GET /version` endpoint.
- The `GET /version` endpoint MUST return `200 OK` with `Content-Type: application/json`.
- The `GET /version` response body MUST be the version string as a bare JSON string (e.g. `"1.2.3"`).
- The backend MUST be hostable locally via `air`.

## Scenarios

### Scenario: GET /version returns the version

Given the backend was built from a repo with `VERSION` containing `1.2.3`  
When a client sends `GET /version`  
Then the response status is `200 OK`  
And the response body matches the contents of the `VERSION` file

### Scenario: air hosts the backend

Given `air` is installed  
When `just backend` is run  
Then `air` watches `cmd/server` and serves the backend with live reload
