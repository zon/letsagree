# Backend Auth

## Purpose

The backend authenticates users via the [Humanity Protocol](https://www.humanity.org) OIDC provider. Each user is verified as a unique human being before being permitted to vote. Identity is established by the `sub` claim from the ID token; the `is_human` claim is the biometric guarantee that gates all privileged actions.

## Requirements

### Configuration

- The backend MUST read Humanity Protocol configuration from `config/humanity-protocol.yaml`.
- The config file MUST contain the following properties:
  - `clientId` — the OAuth client identifier
  - `clientSecret` — the OAuth client secret
  - `publicKey` — the Humanity Protocol public key
  - `issuerUrl` — the OIDC issuer base URL (e.g. `https://api.sandbox.humanity.org/v2`)
  - `redirectUrl` — the absolute callback URL registered with the provider
- The backend MUST discover OIDC endpoints via `<issuerUrl>/.well-known/openid-configuration` at startup.

### Login

- The backend MUST expose a `GET /auth/login` endpoint.
- `GET /auth/login` MUST redirect the user to the Humanity Protocol authorization URL.
- The authorization request MUST request the scopes `openid` and `identity:read`.
- `GET /auth/login` MUST set a `state` cookie containing a randomly generated, unguessable value.
- The `state` cookie MUST be `HttpOnly`, `SameSite=Lax`, and `Secure` in production.

### Callback

- The backend MUST expose a `GET /auth/callback` endpoint.
- `GET /auth/callback` MUST reject requests whose `state` query parameter does not match the `state` cookie, returning `400 Bad Request`.
- `GET /auth/callback` MUST exchange the `code` query parameter for an ID token and access token.
- `GET /auth/callback` MUST verify the ID token signature using the provider's JWKS.
- `GET /auth/callback` MUST extract the `sub` and `is_human` claims from the ID token.
- If `is_human` is `false` or absent, `GET /auth/callback` MUST return `403 Forbidden` with a message indicating that biometric verification is required — it MUST NOT treat this as a generic auth failure.
- On a valid, human-verified callback, the backend MUST upsert a user record keyed on `sub`.
- The backend MUST create a session for the authenticated user and set a `session` cookie.
- The `session` cookie MUST be `HttpOnly`, `SameSite=Lax`, and `Secure` in production.
- Sessions MUST be stored in Postgres.

### Session Middleware

- The backend MUST provide auth middleware that validates the `session` cookie on protected routes.
- If the session cookie is missing or invalid, the middleware MUST return `401 Unauthorized`.
- If the session is valid, the middleware MUST make the authenticated user available to downstream handlers.

### Logout

- The backend MUST expose a `POST /auth/logout` endpoint.
- `POST /auth/logout` MUST delete the user's session from Postgres.
- `POST /auth/logout` MUST clear the `session` cookie.
- `POST /auth/logout` MUST return `204 No Content`.

## Scenarios

### Scenario: Login redirects to Humanity Protocol

Given the backend is configured with valid `HUMANITY_CLIENT_ID` and `HUMANITY_ISSUER_URL`  
When a client sends `GET /auth/login`  
Then the response is a redirect to the Humanity Protocol authorization URL  
And the response sets a `state` cookie  
And the authorization URL includes the `openid` and `identity:read` scopes

### Scenario: Callback rejects mismatched state

Given a client has a `state` cookie value of `abc123`  
When the client sends `GET /auth/callback?code=xyz&state=different`  
Then the response status is `400 Bad Request`

### Scenario: Callback rejects non-human user

Given the Humanity Protocol returns an ID token with `is_human: false`  
When the client sends `GET /auth/callback` with a valid code and matching state  
Then the response status is `403 Forbidden`  
And the response body indicates that biometric verification is required

### Scenario: Callback authenticates a verified human

Given the Humanity Protocol returns an ID token with `sub: "hp_abc"` and `is_human: true`  
When the client sends `GET /auth/callback` with a valid code and matching state  
Then the backend upserts a user record with `sub = "hp_abc"`  
And the response sets a `session` cookie  
And the response redirects to the application

### Scenario: Callback upserts existing user on re-login

Given a user with `sub: "hp_abc"` already exists in Postgres  
When the user completes the OIDC callback again  
Then no duplicate user record is created  
And a new session is issued

### Scenario: Protected route requires session

Given a client sends no `session` cookie  
When the client accesses a protected route  
Then the response status is `401 Unauthorized`

### Scenario: Protected route accepts valid session

Given a client has a valid `session` cookie  
When the client accesses a protected route  
Then the handler proceeds and the authenticated user is available

### Scenario: Logout clears session

Given a client has a valid `session` cookie  
When the client sends `POST /auth/logout`  
Then the session is deleted from Postgres  
And the `session` cookie is cleared  
And the response status is `204 No Content`
