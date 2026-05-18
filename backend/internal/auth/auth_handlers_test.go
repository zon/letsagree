package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"server/internal/oidc"
	"server/internal/store"
)

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
	svc := WithMocks(oidc.NewStubProvider(oidc.IDTokenWithIsHuman(false)))
	w, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "biometric verification required")
}

func TestCallback_authenticatesHuman(t *testing.T) {
	svc := WithMocks(oidc.NewStubProvider(oidc.AnyIDToken()))
	w, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "/")
	assert.NotEmpty(t, SessionCookieValue(w))
}

func TestCallback_upsertsUserOnRelogin(t *testing.T) {
	token := oidc.IDTokenWithSub("hp_abc")
	users := store.StubUsers()
	svc := WithMocks(oidc.NewStubProvider(token), users)
	_, c := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c, "state", "abc")
	svc.Callback(c)
	assert.Equal(t, "hp_abc", users.UpsertedSub(t))
	_, c2 := NewTestContext(http.MethodGet, "/auth/callback?code=xyz&state=abc")
	WithRequestCookie(c2, "state", "abc")
	svc.Callback(c2)
	assert.Equal(t, 1, len(users.Users()))
}

func TestCallback_upsertsUser(t *testing.T) {
	token := oidc.AnyIDToken()
	users := store.StubUsers()
	svc := WithMocks(oidc.NewStubProvider(token), users)
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

func TestRequireAuth_noSession_returns401(t *testing.T) {
	svc := WithMocks(store.NoSessions())
	w, c := NewTestContext(http.MethodGet, "/protected")
	svc.RequireAuth(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_validSession_proceedsAndSetsCurrentUser(t *testing.T) {
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
	header := SessionCookieValue(w)
	assert.Contains(t, header, "session=")
	assert.Contains(t, header, "Max-Age=0")
}