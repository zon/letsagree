# agreectl set config — Postgres Orchestration

## Purpose

Fetch the CNPG app secret from Kubernetes, resolve the external host when not provided, and write `backend/config/postgres.json` with in-cluster host and port replaced.

## Orchestration

**Module:** `agreectl/internal/orchestration`

```go
type Orchestration struct {
	cluster K8sClient
	files   ConfigWriter
}

func (o *Orchestration) Postgres(in opts.Opts) error {
	secret, err := o.cluster.GetSecret(in.Namespace, in.DBSecret)
	if err != nil {
		return err
	}

	host := in.Host
	if host == "" {
		host, err = o.cluster.NodeIP()
		if err != nil {
			return err
		}
	}

	return o.files.WriteJSON(files.PostgresConfigPath, files.PostgresConfig{
		Host:     host,
		Port:     in.Port,
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	})
}
```

### Helpers

- **`cluster.GetSecret(namespace, name)`** — fetches the named Kubernetes secret and returns a `Secret` exposing `User()`, `Password()`, and `DBName()` from its data fields
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
	require.NoError(t, svc.Postgres(opts.WithHost("localhost")))
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
	port := opts.AnyPort()
	svc := orchestration.WithMocks()
	require.NoError(t, svc.Postgres(opts.WithPort(port)))
	assert.Equal(t, port, files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{}).Port)
}
```

### Helpers

- **`orchestration.WithMocks(overrides ...any)`** — constructs an `Orchestration` with default stub implementations; accepts optional override values for `cluster` or `files`
- **`cluster.AnyNodeIP()`** — returns an arbitrary valid node IP string
- **`cluster.WithNodeIP(ip)`** — returns a `K8sClient` stub whose `NodeIP()` returns `ip`
- **`cluster.ThatFailsOnNodeIP()`** — returns a `K8sClient` stub that fails the test if `NodeIP()` is called
- **`cluster.AnySecret()`** — returns a `Secret` stub with arbitrary but stable field values
- **`cluster.WithSecret(s)`** — returns a `K8sClient` stub whose `GetSecret` returns `s`
- **`files.WrittenAt(t, path, out)`** — unmarshals the JSON captured by the `ConfigWriter` stub for `path` into `out` and returns it; fails the test if nothing was written to that path
- **`opts.Any()`** — returns `Opts` with defaults matching the CLI flags (no host, port 30432, namespace `letsagree`, db-secret `letsagree-app`); defined in `agreectl/internal/opts`
- **`opts.WithHost(host)`** — returns `Opts` with `Host` set to the given value; defined in `agreectl/internal/opts`
- **`opts.AnyPort()`** — returns an arbitrary port number distinct from the default; defined in `agreectl/internal/opts`
- **`opts.WithPort(port)`** — returns `Opts` with `Port` set to the given value; defined in `agreectl/internal/opts`
