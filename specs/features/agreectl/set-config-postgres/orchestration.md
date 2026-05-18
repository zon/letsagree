# agreectl set config — Postgres Orchestration

## Purpose

Fetch the CNPG app secret from Kubernetes and write postgres config to two destinations, each adjusted for its environment: the local file gets the NodeIP + NodePort so the dev machine can reach postgres directly; the ralph namespace secret gets the qualified in-cluster service name + ClusterIP port so workflow pods can reach postgres across namespaces.

## Orchestration

**Module:** `agreectl/internal/orchestration`

```go
func (o *Orchestration) Postgres(in opts.Opts) error {
	secret, err := o.cluster.GetSecret(in.Namespace, in.DBSecret)
	if err != nil {
		return err
	}

	nodeIP, err := o.cluster.NodeIP()
	if err != nil {
		return err
	}

	localConfig := files.PostgresConfig{
		Host:     nodeIP,
		Port:     in.DBPort,
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	}

	clusterConfig := files.PostgresConfig{
		Host:     secret.QualifiedHost(in.Namespace),
		Port:     secret.Port(),
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	}

	if err := o.cluster.UpsertSecret(in.RalphNamespace, in.PostgresSecret, clusterConfig.ToSecretData()); err != nil {
		return err
	}

	return o.files.WriteJSON(files.PostgresConfigPath, localConfig)
}
```

### Helpers

- **`cluster.GetSecret(namespace, name)`** — fetches the named Kubernetes secret and returns a `Secret` exposing `Host()`, `Port()`, `QualifiedHost(namespace)`, `User()`, `Password()`, `DBName()`, and `Data()`
- **`cluster.NodeIP()`** — lists cluster nodes and returns the `InternalIP` of the first ready node
- **`cluster.Secret.QualifiedHost(namespace)`** — returns `host + "." + namespace`, making the in-cluster service name routable from other namespaces
- **`cluster.Secret.Port()`** — parses and returns the `port` field as an int
- **`cluster.UpsertSecret(namespace, name, data)`** — creates or updates a Kubernetes secret by name using apply semantics
- **`files.PostgresConfig`** — output struct with fields `Host`, `Port`, `User`, `Password`, `DBName` and corresponding `json:` tags
- **`files.PostgresConfig.ToSecretData()`** — marshals the config as JSON and returns `map[string]string{"postgres.json": <json>}`
- **`files.WriteJSON(path, v)`** — marshals `v` as JSON and writes it to `path`, creating any missing parent directories
- **`files.PostgresConfigPath`** — output path constant: `"backend/config/postgres.json"`
- **`opts.Opts`** — fields used: `Namespace`, `DBSecret`, `DBPort`, `RalphNamespace`, `PostgresSecret`

## Tests

**Module:** `agreectl/internal/orchestration`

```go
func TestPostgres_localFile_usesNodeIPAndNodePort(t *testing.T) {
	ip := cluster.AnyNodeIP()
	port := opts.AnyDBPort()
	svc := orchestration.WithMocks(cluster.WithNodeIP(ip))
	require.NoError(t, svc.Postgres(opts.WithDBPort(port)))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, ip, cfg.Host)
	assert.Equal(t, port, cfg.Port)
}

func TestPostgres_clusterSecret_usesQualifiedHostAndSecretPort(t *testing.T) {
	s := cluster.AnySecret()
	o := opts.Any()
	svc := orchestration.WithMocks(cluster.WithSecret(s))
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
	svc := orchestration.WithMocks(cluster.WithSecret(s))
	require.NoError(t, svc.Postgres(opts.Any()))
	cfg := files.WrittenAt(t, files.PostgresConfigPath, &files.PostgresConfig{})
	assert.Equal(t, s.User(), cfg.User)
	assert.Equal(t, s.Password(), cfg.Password)
	assert.Equal(t, s.DBName(), cfg.DBName)
}
```

### Helpers

- **`orchestration.WithMocks(overrides ...any)`** — constructs an `Orchestration` with default stub implementations; accepts optional override values for `cluster` or `files`
- **`cluster.AnyNodeIP()`** — returns an arbitrary valid node IP string
- **`cluster.WithNodeIP(ip)`** — returns a `K8sClient` stub whose `NodeIP()` returns `ip`
- **`cluster.UpsertedSecretData(t)`** — returns the data map passed to the most recent `UpsertSecret` call; fails the test if no upsert occurred
- **`cluster.AnySecret()`** — returns a `Secret` stub with arbitrary but stable field values including `host` and `port`
- **`cluster.WithSecret(s)`** — returns a `K8sClient` stub whose `GetSecret` returns `s` for any call
- **`files.WrittenAt(t, path, out)`** — unmarshals the JSON captured by the `ConfigWriter` stub for `path` into `out` and returns it; fails the test if nothing was written to that path
- **`opts.Any()`** — returns `Opts` with defaults matching the CLI flags
- **`opts.AnyDBPort()`** — returns the default NodePort value
- **`opts.WithDBPort(port)`** — returns `Opts` with `DBPort` set to the given value
