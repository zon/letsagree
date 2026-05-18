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

func TestPostgres_localFile_usesNodeIPAndNodePort(t *testing.T) {
	ip := cluster.AnyNodeIP()
	port := opts.AnyDBPort()
	svc := withMocks(cluster.WithNodeIP(ip))
	require.NoError(t, svc.Postgres(opts.WithDBPort(port)))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, ip, cfg.Host)
	assert.Equal(t, port, cfg.Port)
}

func TestPostgres_clusterSecret_usesQualifiedHostAndSecretPort(t *testing.T) {
	s := cluster.AnySecret()
	o := opts.Any()
	svc := withMocks(cluster.WithSecret(s))
	require.NoError(t, svc.Postgres(o))
	expected := files.PostgresConfig{
		Host:     s.QualifiedHost(o.Namespace),
		Port:     s.Port(),
		User:     s.User(),
		Password: s.Password(),
		DBName:   s.DBName(),
	}
	assert.Equal(t, expected.ToSecretData(), cluster.UpsertedSecretData(t))
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

func TestHumanityProtocol_envFileUpsertClusterSecret(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := withMocks(files.WithHPEnv(creds))
	o := opts.WithHPEnvFile("any.env")
	require.NoError(t, svc.HumanityProtocol(o))
	expected := files.HumanityProtocolConfig{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		PublicKey:    creds.PublicKey,
		IssuerURL:    o.OIDCIssuer,
		RedirectURL:  o.OIDCRedirect,
	}
	assert.Equal(t, expected.ToSecretData(), cluster.UpsertedSecretData(t))
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

func TestHumanityProtocol_existingSecretWritesLocalConfig(t *testing.T) {
	creds := files.AnyHPCredentials()
	hpConfig := files.HumanityProtocolConfig{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		PublicKey:    creds.PublicKey,
	}
	svc := withMocks(cluster.WithSecret(cluster.SecretFromStringData(hpConfig.ToSecretData())))
	require.NoError(t, svc.HumanityProtocol(opts.Any()))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, creds.ClientID, cfg.ClientID)
	assert.Equal(t, creds.ClientSecret, cfg.ClientSecret)
	assert.Equal(t, creds.PublicKey, cfg.PublicKey)
}

func TestHumanityProtocol_noHPDataInBackendSecret_returnsError(t *testing.T) {
	svc := withMocks()
	err := svc.HumanityProtocol(opts.Any())
	require.Error(t, err)
}

func TestHumanityProtocol_writesOIDCOptions(t *testing.T) {
	hpConfig := files.HumanityProtocolConfig{ClientID: "any", ClientSecret: "any", PublicKey: "any"}
	svc := withMocks(cluster.WithSecret(cluster.SecretFromStringData(hpConfig.ToSecretData())))
	require.NoError(t, svc.HumanityProtocol(opts.WithOIDCOptions("https://issuer.example.com", "https://app.example.com/auth/callback")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, "https://issuer.example.com", cfg.IssuerURL)
	assert.Equal(t, "https://app.example.com/auth/callback", cfg.RedirectURL)
}