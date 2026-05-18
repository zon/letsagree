# Frontend Auth

## Purpose

The frontend initiates and completes the Humanity Protocol OIDC login flow by directing the user to the backend auth endpoints, handling post-auth redirects, and providing logout. Access to protected pages requires a valid session cookie issued by the backend.

## Requirements

### Login

- The frontend MUST expose a `GET /login` route that renders a login page.
- The login page MUST display a button that initiates the Humanity Protocol login flow.
- Clicking the login button MUST redirect the browser to the backend `GET /auth/login` endpoint.
- If the user is already authenticated, `GET /login` MUST redirect to the home page.

### Post-Auth Redirect

- After a successful OIDC callback, the backend redirects the browser to `GET /` on the frontend.
- The frontend MUST expose a `GET /` route that serves the authenticated home page.
- If the session cookie is absent or invalid when accessing `GET /`, the frontend MUST redirect to `GET /login`.

### Biometric Verification Error

- If the backend returns `403 Forbidden` during the auth flow, the frontend MUST display an error page indicating that biometric verification via Humanity Protocol is required.

### Logout

- The authenticated home page MUST display a logout button.
- Clicking the logout button MUST send `POST /auth/logout` to the backend.
- After a successful logout response, the frontend MUST redirect to `GET /login`.

## Scenarios

### Scenario: Unauthenticated user visits home page

Given a user has no `session` cookie  
When the user navigates to `GET /`  
Then the frontend redirects to `GET /login`

### Scenario: Login page renders for unauthenticated user

Given a user has no `session` cookie  
When the user navigates to `GET /login`  
Then the login page is displayed with a "Login with Humanity Protocol" button

### Scenario: Login button initiates OIDC flow

Given the user is on the login page  
When the user clicks the login button  
Then the browser is redirected to the backend `GET /auth/login` endpoint

### Scenario: Authenticated user is redirected away from login

Given a user has a valid `session` cookie  
When the user navigates to `GET /login`  
Then the frontend redirects to `GET /`

### Scenario: Successful login lands on home page

Given the OIDC flow completed and the backend set a `session` cookie  
When the backend redirects the browser to `GET /`  
Then the home page is displayed  
And a logout button is visible

### Scenario: Non-human user sees biometric error

Given the backend returned `403 Forbidden` during the OIDC callback  
When the browser follows the redirect  
Then an error page is displayed indicating that biometric verification via Humanity Protocol is required

### Scenario: Logout clears session and returns to login

Given the user is authenticated and on the home page  
When the user clicks the logout button  
Then `POST /auth/logout` is sent to the backend  
And the frontend redirects to `GET /login`
