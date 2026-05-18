# agreectl set config — Postgres Orchestration

## Purpose

Fetch the CNPG app secret from Kubernetes, copy it into the ralph namespace, resolve the external host when not provided, and write `backend/config/postgres.json` with in-cluster host and port replaced.

## Orchestration

**Module:** `agreectl/internal/orchestration`

```go
func (o *Orchestration) Postgres(in opts.Opts) error {
	secret, err := o.cluster.GetSecret(in.Namespace, in.DBSecret)
	if err != nil {
		return err
	}

	if err := o.cluster.UpsertSecret(in.RalphNamespace, in.DBSecret, secret.Data()); err != nil {
		return err
	}

	host := in.DBHost
	if host == "" {
		host, err = o.cluster.NodeIP()
		if err != nil {
			return err
		}
	}

	return o.files.WriteJSON(files.PostgresConfigPath, files.PostgresConfig{
		Host:     host,
		Port:     in.DBPort,
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	})
}
```

### Helpers

- **`cluster.GetSecret(namespace, name)`** — fetches the named Kubernetes secret and returns a `Secret` exposing `User()`, `Password()`, `DBName()`, and `Data()` from its data fields
- **`cluster.UpsertSecret(namespace, name, data)`** — creates or updates a Kubernetes secret by name using apply semantics
- **`cluster.NodeIP()`** — lists cluster nodes and returns the `InternalIP` of the first ready node
- **`files.WriteJSON(path, v)`** — marshals `v` as JSON and writes it to `path`, creating any missing parent directories
- **`files.PostgresConfigPath`** — output path constant: `"backend/config/postgres.json"`

## Tests

**Module:** `agreectl/internal/orchestration`

```go
func TestPostgres_autoDetectsNodeIP(t *testing.T) {
	ip := cluster.AnyNodeIP()
	svc := orchestration.WithMocks(cluster.WithNodeIP(ip))
	require.NoError(t, svc.Postgres(opts.Any()))
	assert.Equal(t, ip, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_usesProvidedHost(t *testing.T) {
	svc := orchestration.WithMocks(cluster.ThatFailsOnNodeIP())
	require.NoError(t, svc.Postgres(opts.WithDBHost("localhost")))
	assert.Equal(t, "localhost", files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Host)
}

func TestPostgres_copiesSecretFields(t *testing.T) {
	s := cluster.AnySecret()
	svc := orchestration.WithMocks(cluster.WithSecret(s))
	require.NoError(t, svc.Postgres(opts.Any()))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, s.User(), cfg.User)
	assert.Equal(t, s.Password(), cfg.Password)
	assert.Equal(t, s.DBName(), cfg.DBName)
}

func TestPostgres_usesOptsPort(t *testing.T) {
	port := opts.AnyDBPort()
	svc := orchestration.WithMocks()
	require.NoError(t, svc.Postgres(opts.WithDBPort(port)))
	assert.Equal(t, port, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Port)
}

func TestPostgres_copiesSecretToRalphNamespace(t *testing.T) {
	secret := cluster.AnySecret()
	svc := orchestration.WithMocks(cluster.WithSecret(secret))
	require.NoError(t, svc.Postgres(opts.WithRalphNamespace("ralph-letsagree")))
	assert.Equal(t, secret.Data(), cluster.UpsertedSecretData(t))
}
```

### Helpers

- **`orchestration.WithMocks(overrides ...any)`** — constructs an `Orchestration` with default stub implementations; accepts optional override values for `cluster` or `files`
- **`cluster.AnyNodeIP()`** — returns an arbitrary valid node IP string
- **`cluster.WithNodeIP(ip)`** — returns a `K8sClient` stub whose `NodeIP()` returns `ip`
- **`cluster.ThatFailsOnNodeIP()`** — returns a `K8sClient` stub that fails the test if `NodeIP()` is called
- **`cluster.ThatFailsOnUpsert()`** — returns a `K8sClient` stub that fails the test if `UpsertSecret()` is called
- **`cluster.UpsertedSecretData(t)`** — returns the data map passed to the most recent `UpsertSecret` call; fails the test if no upsert occurred
- **`cluster.AnySecret()`** — returns a `Secret` stub with arbitrary but stable field values
- **`cluster.WithSecret(s)`** — returns a `K8sClient` stub whose `GetSecret` returns `s`
- **`files.WrittenAt(t, path, out)`** — unmarshals the JSON captured by the `ConfigWriter` stub for `path` into `out` and returns it; fails the test if nothing was written to that path
- **`opts.Any()`** — returns `Opts` with defaults matching the CLI flags (no host, port 30432, namespace `letsagree`, ralph-namespace `ralph-letsagree`, db-secret `letsagree-app`); defined in `agreectl/internal/opts`
- **`opts.WithDBHost(host)`** — returns `Opts` with `DBHost` set to the given value; defined in `agreectl/internal/opts`
- **`opts.AnyDBPort()`** — returns an arbitrary port number distinct from the default; defined in `agreectl/internal/opts`
- **`opts.WithDBPort(port)`** — returns `Opts` with `DBPort` set to the given value; defined in `agreectl/internal/opts`
- **`opts.WithRalphNamespace(ns)`** — returns `Opts` with `RalphNamespace` set to the given value; defined in `agreectl/internal/opts`
