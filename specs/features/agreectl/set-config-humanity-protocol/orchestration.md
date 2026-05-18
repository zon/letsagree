# agreectl set config — Humanity Protocol Orchestration

## Purpose

Parse Humanity Protocol credentials from an env file or extract them from the existing `humanity-protocol` secret, write the full YAML config back into the `humanity-protocol` secret, and write `backend/config/humanity-protocol.yaml`.

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
		creds = parsed
	} else {
		secret, err := o.cluster.GetSecret(in.RalphNamespace, in.HPSecret)
		if err != nil {
			return errors.New("humanity protocol config not found in secret; provide --hp-env")
		}
		yamlData := secret.Data()["humanity-protocol.yaml"]
		if yamlData == "" {
			return errors.New("humanity protocol config not found in secret; provide --hp-env")
		}
		var existing files.HumanityProtocolConfig
		if err := yaml.Unmarshal([]byte(yamlData), &existing); err != nil {
			return err
		}
		creds = files.HPCredentials{
			ClientID:     existing.ClientID,
			ClientSecret: existing.ClientSecret,
			PublicKey:    existing.PublicKey,
		}
	}

	config := files.HumanityProtocolConfig{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		PublicKey:    creds.PublicKey,
		IssuerURL:    in.OIDCIssuer,
		RedirectURL:  in.OIDCRedirect,
	}

	if err := o.cluster.UpsertSecret(in.RalphNamespace, in.HPSecret, config.ToSecretData()); err != nil {
		return err
	}

	return o.files.WriteYAML(files.HumanityProtocolConfigPath, config)
}
```

### Helpers

- **`files.ParseHPEnv(path)`** — parses an env-format file and returns `HPCredentials` populated from `HUMANITY_CLIENT_ID`, `HUMANITY_CLIENT_SECRET`, and `HUMANITY_PUBLIC_KEY`; returns an error naming any missing keys
- **`cluster.UpsertSecret(namespace, name, data)`** — creates or updates the named Kubernetes secret with the given key/value data
- **`cluster.GetSecret(namespace, name)`** — fetches the named Kubernetes secret; returns an error when the secret is absent
- **`files.HumanityProtocolConfig.ToSecretData()`** — marshals the config as YAML and returns `map[string]string{"humanity-protocol.yaml": <yaml>}`
- **`files.WriteYAML(path, v)`** — marshals `v` as YAML and writes it to `path`, creating any missing parent directories
- **`files.HumanityProtocolConfigPath`** — output path constant: `"backend/config/humanity-protocol.yaml"`
- **`files.HumanityProtocolConfig`** — output struct with fields `ClientID`, `ClientSecret`, `PublicKey`, `IssuerURL`, `RedirectURL` and corresponding `yaml:` tags
- **`opts.Opts`** — fields used: `RalphNamespace`, `HPSecret`, `HPEnvFile`, `OIDCIssuer`, `OIDCRedirect`

## Tests

**Module:** `agreectl/internal/orchestration`

```go
func TestHumanityProtocol_envFileUpsertClusterSecret(t *testing.T) {
	creds := files.AnyHPCredentials()
	svc := orchestration.WithMocks(files.WithHPEnv(creds))
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
	svc := orchestration.WithMocks(files.WithHPEnv(creds))
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
	svc := orchestration.WithMocks(cluster.WithSecret(cluster.SecretFromStringData(hpConfig.ToSecretData())))
	require.NoError(t, svc.HumanityProtocol(opts.Any()))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, creds.ClientID, cfg.ClientID)
	assert.Equal(t, creds.ClientSecret, cfg.ClientSecret)
	assert.Equal(t, creds.PublicKey, cfg.PublicKey)
}

func TestHumanityProtocol_noHPDataInSecret_returnsError(t *testing.T) {
	svc := orchestration.WithMocks()
	err := svc.HumanityProtocol(opts.Any())
	require.Error(t, err)
}

func TestHumanityProtocol_writesOIDCOptions(t *testing.T) {
	hpConfig := files.HumanityProtocolConfig{ClientID: "any", ClientSecret: "any", PublicKey: "any"}
	svc := orchestration.WithMocks(cluster.WithSecret(cluster.SecretFromStringData(hpConfig.ToSecretData())))
	require.NoError(t, svc.HumanityProtocol(opts.WithOIDCOptions("https://issuer.example.com", "https://app.example.com/auth/callback")))
	cfg := files.WrittenYAMLAt(t, files.HumanityProtocolConfigPath, &files.HumanityProtocolConfig{})
	assert.Equal(t, "https://issuer.example.com", cfg.IssuerURL)
	assert.Equal(t, "https://app.example.com/auth/callback", cfg.RedirectURL)
}
```

### Helpers

- **`files.AnyHPCredentials()`** — returns an `HPCredentials` with arbitrary but stable values
- **`files.WithHPEnv(creds)`** — returns a `ConfigWriter` override whose `ParseHPEnv` returns `creds` for any path
- **`files.WrittenYAMLAt(t, path, out)`** — unmarshals the YAML captured by the `ConfigWriter` stub for `path` into `out` and returns it; fails the test if nothing was written to that path
- **`cluster.SecretFromStringData(data)`** — constructs a `*Secret` from a `map[string]string`; used to build an HP secret containing the YAML config for test stubs
- **`cluster.WithSecret(s)`** — returns a `K8sClient` stub whose `GetSecret` returns `s` for any call
- **`cluster.UpsertedSecretData(t)`** — returns the `map[string]string` passed to the most recent `UpsertSecret` call; fails the test if no upsert occurred
- **`opts.Any()`** — returns `Opts` with defaults: `RalphNamespace: "ralph-letsagree"`, `HPSecret: "humanity-protocol"`, `OIDCIssuer: "https://api.sandbox.humanity.org/v2"`
- **`opts.WithHPEnvFile(path)`** — returns `Opts` with `HPEnvFile` set to `path`
- **`opts.WithOIDCOptions(issuer, redirect)`** — returns `Opts` with `OIDCIssuer` and `OIDCRedirect` set
