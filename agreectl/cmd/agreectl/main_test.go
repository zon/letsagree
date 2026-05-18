package main

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func secretWithHP(t *testing.T) *cluster.Secret {
	t.Helper()
	hpConfig := files.HumanityProtocolConfig{
		ClientID: "hp-client-id", ClientSecret: "hp-client-secret", PublicKey: "hp-public-key",
	}
	return cluster.SecretFromStringData(hpConfig.ToSecretData())
}

func TestSetConfig_customRalphNamespace(t *testing.T) {
	stub := &cluster.StubK8sClient{
		Secret: secretWithHP(t),
	}

	cfg := SetConfig{
		Context:        "k3s-prod",
		RalphNamespace: "infra",
		HPSecret:       "humanity-protocol",
		OIDCIssuer:     "https://api.sandbox.humanity.org/v2",
	}

	require.NoError(t, cfg.RunWith(func(_ string) (cluster.K8sClient, error) {
		return stub, nil
	}, &files.CapturingConfigWriter{}))

	assert.True(t, slices.Contains(stub.Calls, cluster.GetSecretCall{Namespace: "infra", Secret: "humanity-protocol"}))
}
