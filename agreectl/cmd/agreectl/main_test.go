package main

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetConfig_customContextNamespaceAndSecret(t *testing.T) {
	stub := &cluster.StubK8sClient{
		Secret:    cluster.AnyHPSecret(),
		RetNodeIP: "10.0.0.1",
	}

	cfg := SetConfig{
		Context:        "k3s-prod",
		RalphNamespace: "infra",
		HPSecret:       "hp-creds",
		OIDCIssuer:     "https://api.sandbox.humanity.org/v2",
		OIDCRedirect:   "https://example.com/auth/callback",
	}

	require.NoError(t, cfg.RunWith(func(_ string) (cluster.K8sClient, error) {
		return stub, nil
	}, &files.CapturingConfigWriter{}))

	assert.Equal(t, "infra", stub.Calls.Namespace)
	assert.Equal(t, "hp-creds", stub.Calls.Secret)
}

func TestSetConfig_customContextNamespaceAndHPSecret(t *testing.T) {
	stub := &cluster.StubK8sClient{
		Secret: cluster.AnyHPSecret(),
	}

	cfg := SetConfig{
		Context:        "k3s-prod",
		RalphNamespace: "infra",
		HPSecret:       "hp-creds",
		OIDCIssuer:     "https://api.sandbox.humanity.org/v2",
		OIDCRedirect:   "https://example.com/auth/callback",
	}

	require.NoError(t, cfg.RunWith(func(_ string) (cluster.K8sClient, error) {
		return stub, nil
	}, &files.CapturingConfigWriter{}))

	assert.Equal(t, "infra", stub.Calls.Namespace)
	assert.Equal(t, "hp-creds", stub.Calls.Secret)
}