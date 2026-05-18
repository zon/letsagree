# Backend Auth Orchestration

## Purpose

Four Gin handlers — Login, Callback, RequireAuth middleware, and Logout — coordinate an OIDC provider, a session store, and a user store to authenticate users via Humanity Protocol and gate protected routes on a verified session.

## Orchestration

**Module:** `backend/internal/auth`

```go
type OIDCProvider interface {
	AuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oidc.IDToken, error)
}

type SessionStore interface {
	Create(userID uint) (string, error)
	Get(token string) (*store.Session, error)
	Delete(token string) error
}

type UserStore interface {
	Upsert(sub string) (*store.User, error)
}

type Orchestration struct {
	provider OIDCProvider
	sessions SessionStore
	users    UserStore
}

func (o *Orchestration) Login(c *gin.Context) {
	state := NewState()
	SetStateCookie(c, state)
	c.Redirect(http.StatusFound, o.provider.AuthURL(state))
}

func (o *Orchestration) Callback(c *gin.Context) {
	if c.Query("state") != StateCookie(c) {
		c.Status(http.StatusBadRequest)
		return
	}

	token, err := o.provider.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !token.IsHuman() {
		c.JSON(http.StatusForbidden, gin.H{"error": "biometric verification required"})
		return
	}

	user, err := o.users.Upsert(token.Sub())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	sessionToken, err := o.sessions.Create(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	SetSessionCookie(c, sessionToken)
	c.Redirect(http.StatusFound, "/")
}

func (o *Orchestration) RequireAuth(c *gin.Context) {
	session, err := o.sessions.Get(SessionCookie(c))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	SetCurrentUser(c, session.UserID)
	c.Next()
}

func (o *Orchestration) Logout(c *gin.Context) {
	o.sessions.Delete(SessionCookie(c))
	ClearSessionCookie(c)
	c.Status(http.StatusNoContent)
}
```

### Helpers

- **`NewState()`** — generates a cryptographically random, URL-safe state string
- **`SetStateCookie(c, state)`** — writes the state value as an `HttpOnly`, `SameSite=Lax` cookie named `state` on the response
- **`StateCookie(c)`** — reads the `state` cookie value from the request; returns `""` if absent
- **`SetSessionCookie(c, token)`** — writes the session token as an `HttpOnly`, `SameSite=Lax` cookie named `session` on the response
- **`SessionCookie(c)`** — reads the `session` cookie value from the request; returns `""` if absent
- **`ClearSessionCookie(c)`** — overwrites the `session` cookie with an empty value and a past expiry
- **`SetCurrentUser(c, userID)`** — stores the user ID in the Gin context under a package-scoped key
- **`CurrentUser(c)`** — retrieves the user ID stored by `SetCurrentUser`
- **`oidc.IDToken`** — domain type wrapping the verified ID token claims; exposes `Sub() string` and `IsHuman() bool`
- **`store.Session`** — GORM model with `Token string` and `UserID uint` fields
- **`store.User`** — GORM model with `ID uint` and `Sub string` fields

## Tests

**Module:** `backend/internal/auth`

```go
func TestLogin_redirectsWithStateCookie(t *testing.T) {
	svc := WithMocks()
	w, c := NewTestContext(http.MethodGet, "/auth/login")
	svc.Login(c)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.NotEmpty(t, StateCookieValue(w))
}

func TestCallback_rejectsMismatchedState(t *testing.T) {
	svc := WithMocks()
	w, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=wrong")
	WithRequestCookie(c, "state", "correct")
	svc.Callback(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCallback_rejectsNonHuman(t *testing.T) {
	svc := WithMocks(oidc.StubProvider(oidc.IDTokenWithIsHuman(false)))
	w, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCallback_authenticatesHuman(t *testing.T) {
	svc := WithMocks(oidc.StubProvider(oidc.AnyIDToken()))
	w, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.NotEmpty(t, SessionCookieValue(w))
}

func TestCallback_upsertsUser(t *testing.T) {
	token := oidc.AnyIDToken()
	users := store.StubUsers()
	svc := WithMocks(oidc.StubProvider(token), users)
	_, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, token.Sub(), users.UpsertedSub(t))
}

func TestRequireAuth_rejectsNoSession(t *testing.T) {
	svc := WithMocks(store.NoSessions())
	w, c := NewTestContext(http.MethodGet, "/protected")
	svc.RequireAuth(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_setsCurrentUser(t *testing.T) {
	session := store.AnySession()
	svc := WithMocks(store.WithSession(session))
	_, c := NewTestContext(http.MethodGet, "/protected")
	WithRequestCookie(c, "session", session.Token)
	svc.RequireAuth(c)
	assert.Equal(t, session.UserID, CurrentUser(c))
}

func TestLogout_deletesSession(t *testing.T) {
	sessions := store.StubSessions()
	svc := WithMocks(sessions)
	_, c := NewTestContext(http.MethodPost, "/auth/logout")
	WithRequestCookie(c, "session", "tok123")
	svc.Logout(c)
	assert.False(t, sessions.Has("tok123"))
}

func TestLogout_clearsSessionCookie(t *testing.T) {
	svc := WithMocks()
	w, c := NewTestContext(http.MethodPost, "/auth/logout")
	svc.Logout(c)
	assert.Empty(t, SessionCookieValue(w))
}
```

### Helpers

- **`WithMocks(overrides ...any)`** — constructs an `Orchestration` with stub `OIDCProvider`, `SessionStore`, and `UserStore`; accepts optional overrides for any dependency
- **`NewTestContext(method, path)`** — returns an `*httptest.ResponseRecorder` and a `*gin.Context` wired to it
- **`WithRequestCookie(c, name, value)`** — adds a cookie to the test context's request
- **`StateCookieValue(w)`** — reads the `state` cookie value set on the response recorder
- **`SessionCookieValue(w)`** — reads the `session` cookie value set on the response recorder
- **`oidc.AnyIDToken()`** — returns an `*IDToken` with an arbitrary `sub` and `IsHuman: true`; defined in `backend/internal/oidc`
- **`oidc.IDTokenWithIsHuman(v)`** — returns an `*IDToken` with `IsHuman` set to `v`; defined in `backend/internal/oidc`
- **`oidc.StubProvider(token)`** — returns an `OIDCProvider` stub whose `Exchange` returns `token`; defined in `backend/internal/oidc`
- **`store.AnySession()`** — returns a `*Session` with arbitrary but stable `Token` and `UserID`; defined in `backend/internal/store`
- **`store.WithSession(s)`** — returns a `SessionStore` stub whose `Get` returns `s` for any token
- **`store.NoSessions()`** — returns a `SessionStore` stub whose `Get` always returns an error
- **`store.StubSessions()`** — returns a recording `SessionStore` stub; supports `Has(token)` to check if a token is still present
- **`store.StubUsers()`** — returns a recording `UserStore` stub; supports `UpsertedSub(t)` to retrieve the `sub` passed to `Upsert`
