package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"server/internal/auth"
	"server/internal/oidc"
	"server/internal/store"
)

func TestWiring_authRoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	o := auth.WithMocks()
	RegisterRoutes(r, o)

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{http.MethodGet, "/auth/login", http.StatusFound},
		{http.MethodPost, "/auth/logout", http.StatusNoContent},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, tc.code, w.Code, "expected %d for %s %s", tc.code, tc.method, tc.path)
	}
}

func TestWiring_callbackWithMismatchedState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	o := auth.WithMocks()
	RegisterRoutes(r, o)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/callback?code=xyz&state=wrong", nil)
	req.AddCookie(&http.Cookie{Name: "state", Value: "correct"})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWiring_callbackWithValidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	o := auth.WithMocks(oidc.NewStubProvider(oidc.AnyIDToken()))
	RegisterRoutes(r, o)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/callback?code=xyz&state=abc", nil)
	req.AddCookie(&http.Cookie{Name: "state", Value: "abc"})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestWiring_protectedRouteWithoutSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	o := auth.WithMocks(store.NoSessions())
	RegisterRoutes(r, o)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWiring_protectedRouteWithSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	session := store.AnySession()
	o := auth.WithMocks(store.WithSession(session))
	RegisterRoutes(r, o)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: session.Token})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}