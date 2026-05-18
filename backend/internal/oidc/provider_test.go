package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tmp := t.TempDir()
	fixture := tmp + "/config.yaml"
	writeFile(t, fixture, `
clientId: my-client-id
clientSecret: my-client-secret
issuerUrl: https://issuer.example.com
redirectUrl: https://app.example.com/callback
`)

	cfg, err := LoadConfig(fixture)
	assert.NoError(t, err)
	assert.Equal(t, "my-client-id", cfg.ClientID)
	assert.Equal(t, "my-client-secret", cfg.ClientSecret)
	assert.Equal(t, "https://issuer.example.com", cfg.IssuerURL)
	assert.Equal(t, "https://app.example.com/callback", cfg.RedirectURL)
}

func TestLoadConfig_notFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_invalidYAML(t *testing.T) {
	tmp := t.TempDir()
	fixture := tmp + "/config.yaml"
	writeFile(t, fixture, "invalid: yaml: content:")

	_, err := LoadConfig(fixture)
	assert.Error(t, err)
}

func TestStubProvider_AuthURL(t *testing.T) {
	p := NewStubProvider(AnyIDToken())
	url := p.AuthURL("my-state")
	assert.Contains(t, url, "state=my-state")
}

func TestStubProvider_Exchange(t *testing.T) {
	token := IDTokenWithIsHuman(true)
	p := NewStubProvider(token)
	result, err := p.Exchange(context.Background(), "code123")
	assert.NoError(t, err)
	assert.Equal(t, token, result)
}

func TestAnyIDToken(t *testing.T) {
	token := AnyIDToken()
	assert.Equal(t, "hp_test_sub", token.Sub())
	assert.True(t, token.IsHuman())
}

func TestIDTokenWithIsHuman(t *testing.T) {
	token := IDTokenWithIsHuman(false)
	assert.False(t, token.IsHuman())
}

func TestIDTokenWithSub(t *testing.T) {
	token := IDTokenWithSub("hp-custom")
	assert.Equal(t, "hp-custom", token.Sub())
}

func writeFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write file: %v", err)
	}
}

type mockOIDC struct {
	server *httptest.Server
	tokens map[string]struct{ Sub, IsHuman bool }
}

func (m *mockOIDC) handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"authorization_endpoint": "` + m.server.URL + `/authorize",
			"token_endpoint": "` + m.server.URL + `/token",
			"jwks_uri": "` + m.server.URL + `/jwks"
		}`))
	case "/jwks":
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keys":[]}`))
	case "/authorize":
		q := r.URL.Query()
		code := "code-" + q.Get("state")
		if q.Get("response_type") == "code" {
			loc := q.Get("redirect_uri") + "?code=" + code + "&state=" + q.Get("state")
			http.Redirect(w, r, loc, http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	case "/token":
		r.ParseForm()
		code := r.FormValue("code")
		sub, isHuman := m.subAndHumanFor(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"at","token_type":"Bearer","id_token":"` +
			m.signedIDToken(sub, isHuman) + `"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *mockOIDC) authCodeFor(sub string, isHuman bool) string {
	return "code-" + sub
}

func (m *mockOIDC) subAndHumanFor(code string) (string, bool) {
	return "hp_test_sub", true
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (m *mockOIDC) signedIDToken(sub string, isHuman bool) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	claims := map[string]interface{}{
		"sub":      sub,
		"is_human": isHuman,
		"aud":      []string{"test-client"},
		"iss":      m.server.URL,
		"exp":      9999999999,
		"iat":      1000000000,
	}
	claimsJSON, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
	return header + "." + payload + "."
}

func mockOIDCServer() *mockOIDC {
	m := &mockOIDC{tokens: make(map[string]struct{ Sub, IsHuman bool })}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	return m
}