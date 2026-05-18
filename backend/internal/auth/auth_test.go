package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	state1 := NewState()
	state2 := NewState()
	assert.NotEmpty(t, state1)
	assert.NotEmpty(t, state2)
	assert.NotEqual(t, state1, state2)
	assert.LessOrEqual(t, len(state1), 44)
}

func TestSetStateCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetStateCookie(c, "test-state")
	header := w.Header().Get("Set-Cookie")
	assert.Contains(t, header, "state=test-state")
	assert.Contains(t, header, "HttpOnly")
	assert.Contains(t, header, "Path=/")
}

func TestSetSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetSessionCookie(c, "session-token")
	header := w.Header().Get("Set-Cookie")
	assert.Contains(t, header, "session=session-token")
	assert.Contains(t, header, "HttpOnly")
	assert.Contains(t, header, "Path=/")
}

func TestStateCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "state", Value: "my-state"})
	assert.Equal(t, "my-state", StateCookie(c))
}

func TestStateCookie_empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	assert.Equal(t, "", StateCookie(c))
}

func TestSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "session", Value: "my-session"})
	assert.Equal(t, "my-session", SessionCookie(c))
}

func TestClearSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ClearSessionCookie(c)
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "session", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.True(t, cookies[0].HttpOnly)
}

func TestSetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetCurrentUser(c, 42)
	assert.Equal(t, uint(42), CurrentUser(c))
}

func TestCurrentUser_notSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	assert.Equal(t, uint(0), CurrentUser(c))
}

func TestWithMocks(t *testing.T) {
	o := WithMocks()
	assert.NotNil(t, o.provider)
	assert.NotNil(t, o.sessions)
	assert.NotNil(t, o.users)
}

func TestNewTestContext(t *testing.T) {
	w, c := NewTestContext(http.MethodGet, "/test")
	assert.NotNil(t, w)
	assert.NotNil(t, c)
	assert.Equal(t, http.MethodGet, c.Request.Method)
	assert.Equal(t, "/test", c.Request.URL.Path)
}

func TestWithRequestCookie(t *testing.T) {
	_, c := NewTestContext(http.MethodGet, "/test")
	WithRequestCookie(c, "state", "abc")
	cookies := c.Request.Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "state", cookies[0].Name)
	assert.Equal(t, "abc", cookies[0].Value)
}

func TestStateCookieValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetStateCookie(c, "state-val")
	val := StateCookieValue(w)
	assert.Contains(t, val, "state=state-val")
}

func TestSessionCookieValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetSessionCookie(c, "session-val")
	val := SessionCookieValue(w)
	assert.Contains(t, val, "session=session-val")
}