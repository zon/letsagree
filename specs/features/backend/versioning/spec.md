# Backend Versioning

## Purpose

The repository maintains a single source of truth for the application version in a `VERSION` file at the repo root. The backend binary is built referencing this file and exposes it via a `/version` HTTP endpoint.

## Requirements

- The repository MUST contain a `VERSION` file at the repo root containing a valid [semver](https://semver.org) string (e.g. `1.0.0`).
- The repository MUST contain a `just build` command that reads the `VERSION` file and passes the version to the backend compiler at build time.
- The backend build MUST embed the version from the `VERSION` file at compile time.
- The backend MUST expose a `GET /version` endpoint.
- The `GET /version` endpoint MUST return `200 OK` with `Content-Type: application/json`.
- The `GET /version` response body MUST be the embedded version string as a bare JSON string (e.g. `"1.2.3"`).
- If the `VERSION` file is absent or contains a non-semver string, the build MUST fail with a descriptive error.

## Scenarios

### Scenario: just build embeds the version

Given the `VERSION` file at the repo root contains `1.2.3`  
When `just build` is run  
Then the backend binary is produced with version `1.2.3` embedded

### Scenario: GET /version returns the embedded version

Given the backend was built with version `1.2.3`  
When a client sends `GET /version`  
Then the response status is `200 OK`  
And the response body is `"1.2.3"`

### Scenario: VERSION file is absent

Given the `VERSION` file does not exist at the repo root  
When the build is attempted  
Then the build fails with an error indicating the missing `VERSION` file

### Scenario: VERSION file contains a non-semver string

Given the `VERSION` file contains `latest` or another non-semver value  
When the build is attempted  
Then the build fails with an error indicating the invalid version format
