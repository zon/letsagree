# Frontend Auth Orchestration

## Purpose

Four route handlers — `homePage`, `loginPage`, `notHumanPage`, and `logout` — coordinate a backend session client and the render module to gate protected pages on a valid session and drive the Humanity Protocol login flow.

## Orchestration

**Module:** `frontend/src/orchestration.ts`

```typescript
type PageResponse =
    | { type: "html"; content: string }
    | { type: "redirect"; to: string }

async function homePage(
    session: string | null,
    fetchUser: () => Promise<User | null>,
): Promise<PageResponse> {
    if (!session) return { type: "redirect", to: "/login" }
    const user = await fetchUser()
    if (!user) return { type: "redirect", to: "/login" }
    return { type: "html", content: renderHome(user) }
}

function loginPage(session: string | null): PageResponse {
    if (session) return { type: "redirect", to: "/" }
    return { type: "html", content: renderLogin() }
}

function notHumanPage(): PageResponse {
    return { type: "html", content: renderNotHuman() }
}

async function logout(doLogout: () => Promise<void>): Promise<PageResponse> {
    await doLogout()
    return { type: "redirect", to: "/login" }
}
```

### Helpers

- **`renderHome(user)`** — returns an HTML string for the authenticated home page including a logout form; defined in `frontend/src/render.ts`
- **`renderLogin()`** — returns an HTML string for the login page with a Humanity Protocol login button; defined in `frontend/src/render.ts`
- **`renderNotHuman()`** — returns an HTML string for the biometric verification error page; defined in `frontend/src/render.ts`
- **`backend.fetchUser(backendUrl, session)`** — retrieves the current user for the given session token from the backend; returns `null` if the session is absent, invalid, or expired; defined in `frontend/src/backend.ts`
- **`backend.logout(backendUrl, session)`** — sends `POST /auth/logout` to the backend to delete the session; defined in `frontend/src/backend.ts`
- **`User`** — domain type representing an authenticated user; defined in `frontend/src/backend.ts`

## Tests

**Module:** `frontend/src/orchestration.test.ts`

```typescript
test("home page redirects to login when session cookie is absent", async () => {
    const response = await homePage(sessions.absent(), backend.userNotFound())
    assertRedirectsTo(response, "/login")
})

test("home page redirects to login when session is invalid", async () => {
    const response = await homePage(sessions.any(), backend.userNotFound())
    assertRedirectsTo(response, "/login")
})

test("home page renders for authenticated user", async () => {
    const user = users.any()
    const response = await homePage(sessions.any(), backend.findingUser(user))
    assertRendersHome(response, user)
})

test("login page renders when no session cookie", () => {
    const response = loginPage(sessions.absent())
    assertRendersLogin(response)
})

test("login page redirects authenticated user to home", () => {
    const response = loginPage(sessions.any())
    assertRedirectsTo(response, "/")
})

test("not human page renders biometric error", () => {
    const response = notHumanPage()
    assertRendersNotHuman(response)
})

test("logout redirects to login", async () => {
    const response = await logout(backend.thatLogsOut())
    assertRedirectsTo(response, "/login")
})
```

### Helpers

- **`sessions.any()`** — returns a non-empty session cookie string; defined in `frontend/src/backend.ts`
- **`sessions.absent()`** — returns `null`, representing a missing session cookie; defined in `frontend/src/backend.ts`
- **`users.any()`** — returns a `User` with arbitrary but valid fields; defined in `frontend/src/backend.ts`
- **`backend.findingUser(user)`** — returns a `fetchUser` stub that resolves to `user`; defined in `frontend/src/backend.ts`
- **`backend.userNotFound()`** — returns a `fetchUser` stub that resolves to `null`; defined in `frontend/src/backend.ts`
- **`backend.thatLogsOut()`** — returns a `doLogout` stub that resolves without error; defined in `frontend/src/backend.ts`
- **`assertRedirectsTo(response, path)`** — asserts `response.type === "redirect"` and `response.to === path`
- **`assertRendersHome(response, user)`** — asserts `response.type === "html"` and the content includes home page markers for `user`
- **`assertRendersLogin(response)`** — asserts `response.type === "html"` and the content includes the login page marker
- **`assertRendersNotHuman(response)`** — asserts `response.type === "html"` and the content includes the biometric error marker
