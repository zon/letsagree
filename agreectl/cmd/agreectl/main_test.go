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
		Secret:    cluster.AnySecret(),
		RetNodeIP: "10.0.0.1",
	}

	cfg := SetConfig{
		Context:   "k3s-prod",
		Namespace: "myns",
		DBSecret:  "myns-app",
		Port:      30432,
	}

	require.NoError(t, cfg.RunWith(func(_ string) (cluster.K8sClient, error) {
		return stub, nil
	}, &files.CapturingConfigWriter{}))

	assert.Equal(t, "myns", stub.Calls.Namespace)
	assert.Equal(t, "myns-app", stub.Calls.Secret)
}