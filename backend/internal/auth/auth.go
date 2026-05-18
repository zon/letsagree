package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"server/internal/oidc"
	"server/internal/store"
)

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

func WithMocks(overrides ...any) *Orchestration {
	o := &Orchestration{
		provider: oidc.NewStubProvider(oidc.AnyIDToken()),
		sessions: store.StubSessions(),
		users:    store.StubUsers(),
	}
	for _, override := range overrides {
		switch v := override.(type) {
		case *oidc.StubProvider:
			o.provider = v
		case *store.MockSessionStore:
			o.sessions = v
		case *store.MockUserStore:
			o.users = v
		}
	}
	return o
}

func NewState() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func SetStateCookie(c *gin.Context, value string) {
	c.SetCookie("state", value, 600, "/", "", false, true)
}

func StateCookie(c *gin.Context) string {
	if cookie, err := c.Cookie("state"); err == nil {
		return cookie
	}
	return ""
}

func SetSessionCookie(c *gin.Context, value string) {
	c.SetCookie("session", value, 86400, "/", "", false, true)
}

func SessionCookie(c *gin.Context) string {
	if cookie, err := c.Cookie("session"); err == nil {
		return cookie
	}
	return ""
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", false, true)
}

func SetCurrentUser(c *gin.Context, userID uint) {
	c.Set("current_user", userID)
}

func CurrentUser(c *gin.Context) uint {
	if id, exists := c.Get("current_user"); exists {
		return id.(uint)
	}
	return 0
}

func NewTestContext(method, path string) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, path, nil)
	return w, c
}

func WithRequestCookie(c *gin.Context, name, value string) {
	c.Request.AddCookie(&http.Cookie{Name: name, Value: value})
}

func StateCookieValue(w *httptest.ResponseRecorder) string {
	return w.Header().Get("Set-Cookie")
}

func SessionCookieValue(w *httptest.ResponseRecorder) string {
	return w.Header().Get("Set-Cookie")
}