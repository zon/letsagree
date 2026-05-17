package orchestration

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"agreectl/internal/opts"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgres_autoDetectsNodeIP(t *testing.T) {
	ip := cluster.AnyNodeIP()
	svc := WithMocks(cluster.WithNodeIP(ip))
	require.NoError(t, svc.Postgres(opts.Any()))
	assert.Equal(t, ip, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_usesProvidedHost(t *testing.T) {
	svc := WithMocks(cluster.ThatFailsOnNodeIP())
	require.NoError(t, svc.Postgres(opts.WithHost("localhost")))
	assert.Equal(t, "localhost", files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_copiesSecretFields(t *testing.T) {
	s := cluster.AnySecret()
	svc := WithMocks(cluster.WithSecret(s))
	require.NoError(t, svc.Postgres(opts.Any()))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, s.User(), cfg.User)
	assert.Equal(t, s.Password(), cfg.Password)
	assert.Equal(t, s.DBName(), cfg.DBName)
}

func TestPostgres_usesOptsPort(t *testing.T) {
	port := opts.AnyPort()
	svc := WithMocks()
	require.NoError(t, svc.Postgres(opts.WithPort(port)))
	assert.Equal(t, port, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Port)
}