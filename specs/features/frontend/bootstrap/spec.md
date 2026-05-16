# Frontend Bootstrap

## Purpose

The frontend is built with a version sourced from the `VERSION` file at the repo root, serves an about page displaying both frontend and backend versions, and is hosted locally via `bun watch` for live-reload development.

## Requirements

- The frontend MUST be hostable locally via `bun watch`.
- A `just frontend` command MUST start the frontend development server.
- The frontend MUST accept command-line arguments via [parseArgs](https://nodejs.org/api/util.html#utilparseargsconfig).
- The frontend MUST accept an `--addr` flag controlling the HTTP listen address (default: `:3000`).
- The frontend MUST accept a `--backend` flag specifying the backend base URL (default: `http://localhost:8080`).
- The frontend MUST expose a `GET /about` route that returns an HTML page.
- The about page MUST display the frontend version read from the `VERSION` file at the repo root.
- The about page MUST display the backend version fetched from the backend `GET /version` endpoint.
- If the backend is unreachable, the about page MUST still render and MUST indicate the backend version is unavailable.

## Scenarios

### Scenario: just frontend starts the dev server

Given `bun` is installed  
When `just frontend` is run  
Then `bun watch` starts the frontend server with live reload

### Scenario: GET /about shows both versions

Given the frontend was started from a repo with `VERSION` containing `1.2.3`  
And the backend is running and returns `"1.2.3"` from `GET /version`  
When a client sends `GET /about`  
Then the response status is `200 OK`  
And the page displays the frontend version `1.2.3`  
And the page displays the backend version `1.2.3`

### Scenario: GET /about when backend is unreachable

Given the backend is not running  
When a client sends `GET /about`  
Then the response status is `200 OK`  
And the page displays the frontend version  
And the page indicates the backend version is unavailable

### Scenario: --addr flag overrides listen address

Given the frontend is started with `--addr :4000`  
When a client sends `GET /about` to port `4000`  
Then the response status is `200 OK`

### Scenario: --addr defaults to :3000

Given the frontend is started with no `--addr` flag  
Then the server listens on `:3000`

### Scenario: --backend flag overrides backend URL

Given the frontend is started with `--backend http://localhost:9090`  
When the about page is requested  
Then the frontend fetches the backend version from `http://localhost:9090/version`
