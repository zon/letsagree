# agreectl set config — Humanity Protocol Orchestration

## Purpose

Parse Humanity Protocol credentials from an env file or pull them from a Kubernetes secret, write the cluster secret when an env file is provided, and write `backend/config/humanity-protocol.yaml` with all five config fields.

## Orchestration

**Module:** `agreectl/internal/orchestration`

```go
type K8sClient interface {
	GetSecret(namespace, name string) (*cluster.Secret, error)
	NodeIP() (string, error)
	UpsertSecret(namespace, name string, data map[string]string) error
}

type ConfigWriter interface {
	WriteJSON(path string, v any) error
	WriteYAML(path string, v any) error
	ParseHPEnv(path string) (files.HPCredentials, error)
}

func (o *Orchestration) HumanityProtocol(in opts.Opts) error {
	var creds files.HPCredentials

	if in.HPEnvFile != "" {
		parsed, err := o.files.ParseHPEnv(in.HPEnvFile)
		if err != nil {
			return err
		}
		if err := o.cluster.UpsertSecret(in.RalphNamespace, in.HPSecret, parsed.ToSecretData()); err != nil {
			return err
		}
		creds = parsed
	} else {
		secret, err := o.cluster.GetSecret(in.RalphNamespace, in.HPSecret)
		if err != nil {
			return err
		}
		creds = files.HPCredentials{
			ClientID:     secret.ClientID(),
			ClientSecret: secret.ClientSecret(),
			PublicKey:    secret.PublicKey(),
		}
	}

	return o.files.WriteYAML(files.HumanityProtocolConfigPath, files.HumanityProtocolConfig{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		PublicKey:    creds.PublicKey,
		IssuerURL:    in.OIDCIssuer,
		RedirectURL:  in.OIDCRedirect,
	})
}
```

### Helpers

- **`files.ParseHPEnv(path)`** — parses an env-format file and returns `HPCredentials` populated from `HUMANITY_CLIENT_ID`, `HUMANITY_CLIENT_SECRET`, and `HUMANITY_PUBLIC_KEY`; returns an error naming any missing keys
- **`files.HPCredentials.ToSecretData()`** — returns a `map[string]string` with keys `clientId`, `clientSecret`, and `publicKey` for writing to a Kubernetes secret
- **`cluster.UpsertSecret(namespace, name, data)`** — creates or updates the named Kubernetes secret with the given key/value data
- **`cluster.GetSecret(namespace, name)`** — fetches the named Kubernetes secret; returns an error when the secret is absent
- **`cluster.Secret.ClientID()`** — returns the `clientId` field from the secret data
- **`cluster.Secret.ClientSecret()`** — returns the `clientSecret` field from the secret data
- **`cluster.Secret.PublicKey()`** — returns the `publicKey` field from the secret data
- **`files.WriteYAML(path, v)`** — marshals `v` as YAML and writes it to `path`, creating any missing parent directories
- **`files.HumanityProtocolConfigPath`** — output path constant: `"backend/config/humanity-protocol.yaml"`
- **`files.HumanityProtocolConfig`** — output struct with fields `ClientID`, `ClientSecret`, `PublicKey`, `IssuerURL`, `RedirectURL` and corresponding `yaml:` tags
- **`opts.Opts`** — gains fields `RalphNamespace`, `HPSecret`, `HPEnvFile`, `OIDCIssuer`, `OIDCRedirect`

## Tests

**Module:** `agreectl/internal/orchestration`

```go
func TestHumanityProtocol_envFileUpsertClusterSecret(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := orchestration.WithMocks(
		cluster.ThatFailsOnGetSecret(),
		files.WithHPEnv(creds),
	)
	require.NoError(t, svc.HumanityProtocol(opts.WithHPEnvFile("any.env")))
	assert.Equal(t, creds.ToSecretData(), cluster.UpsertedSecretData(t))
}

func TestHumanityProtocol_envFileWritesLocalConfig(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := orchestration.WithMocks(files.WithHPEnv(creds))
	require.NoError(t, svc.HumanityProtocol(opts.WithHPEnvFile("any.env")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, creds.ClientID, cfg.ClientID)
	assert.Equal(t, creds.ClientSecret, cfg.ClientSecret)
	assert.Equal(t, creds.PublicKey, cfg.PublicKey)
}

func TestHumanityProtocol_secretPresentWritesLocalConfig(t *testing.T) {
	secret := cluster.AnyHPSecret()
	svc := orchestration.WithMocks(cluster.WithSecret(secret))
	require.NoError(t, svc.HumanityProtocol(opts.AnyHP()))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, secret.ClientID(), cfg.ClientID)
	assert.Equal(t, secret.ClientSecret(), cfg.ClientSecret)
	assert.Equal(t, secret.PublicKey(), cfg.PublicKey)
}

func TestHumanityProtocol_secretPresentSkipsUpsert(t *testing.T) {
	svc := orchestration.WithMocks(cluster.ThatFailsOnUpsert())
	require.NoError(t, svc.HumanityProtocol(opts.AnyHP()))
}

func TestHumanityProtocol_writesOIDCOptions(t *testing.T) {
	svc := orchestration.WithMocks()
	require.NoError(t, svc.HumanityProtocol(opts.WithOIDCOptions("https://issuer.example.com", "https://app.example.com/auth/callback")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, "https://issuer.example.com", cfg.IssuerURL)
	assert.Equal(t, "https://app.example.com/auth/callback", cfg.RedirectURL)
}
```

### Helpers

- **`files.AnyHPCredentials()`** — returns an `HPCredentials` with arbitrary but stable values; defined in `agreectl/internal/files`
- **`files.WithHPEnv(creds)`** — returns a `ConfigWriter` override whose `ParseHPEnv` returns `creds` for any path
- **`files.WrittenYAMLAt(t, path, out)`** — unmarshals the YAML captured by the `ConfigWriter` stub for `path` into `out` and returns it; fails the test if nothing was written to that path; defined in `agreectl/internal/files`
- **`cluster.AnyHPSecret()`** — returns a `*Secret` stub with arbitrary but stable `ClientID`, `ClientSecret`, and `PublicKey` values; defined in `agreectl/internal/cluster`
- **`cluster.ThatFailsOnGetSecret()`** — returns a `K8sClient` stub that fails the test if `GetSecret` is called
- **`cluster.ThatFailsOnUpsert()`** — returns a `K8sClient` stub that fails the test if `UpsertSecret` is called
- **`cluster.UpsertedSecretData(t)`** — returns the `map[string]string` passed to the most recent `UpsertSecret` call on the active stub; fails the test if no upsert occurred
- **`opts.AnyHP()`** — returns `Opts` with HP defaults: `RalphNamespace: "ralph-letsagree"`, `HPSecret: "humanity-protocol"`, `OIDCIssuer: "https://api.sandbox.humanity.org/v2"`; defined in `agreectl/internal/opts`
- **`opts.WithHPEnvFile(path)`** — returns `Opts` with `HPEnvFile` set to `path`
- **`opts.WithOIDCOptions(issuer, redirect)`** — returns `Opts` with `OIDCIssuer` and `OIDCRedirect` set
