package orchestration

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"agreectl/internal/opts"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withMocks(overrides ...any) *Orchestration {
	defaultK8s := cluster.WithSecret(cluster.AnySecret())
	defaultCW := &files.CapturingConfigWriter{}
	var k8s K8sClient = defaultK8s
	var cw ConfigWriter = defaultCW

	for _, o := range overrides {
		switch v := o.(type) {
		case K8sClient:
			k8s = v
		case ConfigWriter:
			cw = v
		}
	}

	return New(k8s, cw)
}

func WithMocks(overrides ...any) *Orchestration {
	return withMocks(overrides...)
}

func TestPostgres_autoDetectsNodeIP(t *testing.T) {
	ip := cluster.AnyNodeIP()
	svc := withMocks(cluster.WithNodeIP(ip))
	require.NoError(t, svc.Postgres(opts.Any()))
	assert.Equal(t, ip, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_usesProvidedHost(t *testing.T) {
	svc := withMocks(cluster.ThatFailsOnNodeIP())
	require.NoError(t, svc.Postgres(opts.WithDBHost("localhost")))
	assert.Equal(t, "localhost", files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_copiesSecretFields(t *testing.T) {
	s := cluster.AnySecret()
	svc := withMocks(cluster.WithSecret(s))
	require.NoError(t, svc.Postgres(opts.Any()))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, s.User(), cfg.User)
	assert.Equal(t, s.Password(), cfg.Password)
	assert.Equal(t, s.DBName(), cfg.DBName)
}

func TestPostgres_usesOptsPort(t *testing.T) {
	port := opts.AnyDBPort()
	svc := withMocks()
	require.NoError(t, svc.Postgres(opts.WithDBPort(port)))
	assert.Equal(t, port, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Port)
}

func TestPostgres_copiesSecretToRalphNamespace(t *testing.T) {
	secret := cluster.AnySecret()
	svc := withMocks(cluster.WithSecret(secret))
	require.NoError(t, svc.Postgres(opts.WithRalphNamespace("ralph-letsagree")))
	assert.Equal(t, secret.Data(), cluster.UpsertedSecretData(t))
}

func TestHumanityProtocol_envFileUpsertClusterSecret(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := withMocks(
		cluster.ThatFailsOnGetSecret(),
		files.WithHPEnv(creds),
	)
	require.NoError(t, svc.HumanityProtocol(opts.WithHPEnvFile("any.env")))
	assert.Equal(t, creds.ToSecretData(), cluster.UpsertedSecretData(t))
}

func TestHumanityProtocol_envFileWritesLocalConfig(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := withMocks(files.WithHPEnv(creds))
	require.NoError(t, svc.HumanityProtocol(opts.WithHPEnvFile("any.env")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, creds.ClientID, cfg.ClientID)
	assert.Equal(t, creds.ClientSecret, cfg.ClientSecret)
	assert.Equal(t, creds.PublicKey, cfg.PublicKey)
}

func TestHumanityProtocol_secretPresentWritesLocalConfig(t *testing.T) {
	secret := cluster.AnyHPSecret()
	svc := withMocks(cluster.WithSecret(secret))
	require.NoError(t, svc.HumanityProtocol(opts.AnyHP()))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, secret.ClientID(), cfg.ClientID)
	assert.Equal(t, secret.ClientSecret(), cfg.ClientSecret)
	assert.Equal(t, secret.PublicKey(), cfg.PublicKey)
}

func TestHumanityProtocol_secretPresentSkipsUpsert(t *testing.T) {
	svc := withMocks(cluster.ThatFailsOnUpsert())
	require.NoError(t, svc.HumanityProtocol(opts.AnyHP()))
}

func TestHumanityProtocol_writesOIDCOptions(t *testing.T) {
	svc := withMocks()
	require.NoError(t, svc.HumanityProtocol(opts.WithOIDCOptions("https://issuer.example.com", "https://app.example.com/auth/callback")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, "https://issuer.example.com", cfg.IssuerURL)
	assert.Equal(t, "https://app.example.com/auth/callback", cfg.RedirectURL)
}