# agreectl set config — Postgres

## Purpose

`agreectl set config` writes local development config files by pulling secrets from a running Kubernetes cluster. This spec covers the postgres config: reading the CNPG app secret and writing `backend/config/postgres.json` with host and port replaced to be reachable from outside the cluster. It also writes the full postgres config into the shared `backend` Kubernetes secret in the ralph namespace so workflow pods can mount it as a config file.

## Requirements

- `agreectl` MUST be a Go CLI tool using [Kong](https://github.com/alecthomas/kong).
- `agreectl` MUST expose a `set config` subcommand.
- `set config` MUST accept the following flags:

  | Flag                | Default           | Description                                              |
  |---------------------|-------------------|----------------------------------------------------------|
  | `--context`         | `microk8s`        | kubectl context to use                                   |
  | `--namespace`       | `letsagree`       | Kubernetes namespace containing the secret               |
  | `--ralph-namespace` | `ralph-letsagree` | Kubernetes namespace for the postgres secret             |
  | `--db-secret`       | `letsagree-app`   | Name of the CNPG app secret                              |
  | `--postgres-secret` | `postgres`        | Name of the postgres secret in Ralph's namespace         |
  | `--db-host`         | _(auto)_          | Override the external postgres host                      |
  | `--db-port`         | `30432`           | External NodePort for postgres                           |

- `set config` MUST read the following fields from the named Kubernetes secret: `user`, `password`, `dbname`.
- When `--db-host` is not provided, `set config` MUST auto-detect the external host by querying the node list from the specified context and using the `InternalIP` of the first ready node.
- `set config` MUST write `backend/config/postgres.json` relative to the repo root, creating `backend/config/` if it does not exist.
- The written JSON MUST contain the fields `host`, `port`, `user`, `password`, and `dbname`.
- The `host` field in the output MUST be the external host (auto-detected or overridden), not the in-cluster service name from the secret.
- The `port` field in the output MUST be the `--db-port` value, not the in-cluster port from the secret.
- The `user`, `password`, and `dbname` fields MUST be copied from the secret unchanged.
- `set config` MUST upsert the `--postgres-secret` secret in `--ralph-namespace` with the key `postgres.json` set to the full JSON content of the postgres config.

## Scenarios

### Scenario: writes postgres config with auto-detected host

Given a running cluster accessible via context `microk8s`  
And the secret `letsagree-app` exists in namespace `letsagree`  
And the first cluster node has InternalIP `192.168.1.10`  
When `agreectl set config` is run with no flags  
Then `backend/config/postgres.json` is written with:
```json
{
  "host": "192.168.1.10",
  "port": 30432,
  "user": "app",
  "password": "<secret password>",
  "dbname": "app"
}
```

### Scenario: writes postgres config to postgres secret

Given the secret `letsagree-app` exists in namespace `letsagree`  
When `agreectl set config` is run  
Then the `postgres` secret in `ralph-letsagree` is upserted with key `postgres.json` containing the full JSON config

### Scenario: --host overrides auto-detection

Given `agreectl set config` is run with `--db-host localhost`  
Then the `host` field in `backend/config/postgres.json` is `"localhost"`  
And the node list is not queried

### Scenario: custom context, namespace, and secret

Given `agreectl set config` is run with `--context k3s-prod --namespace myns --db-secret myns-app`  
Then the secret `myns-app` is fetched from namespace `myns` using context `k3s-prod`

### Scenario: custom ralph-namespace

Given `agreectl set config` is run with `--ralph-namespace infra`  
Then the backend secret is upserted into namespace `infra`

### Scenario: output directory is created if missing

Given `backend/config/` does not exist  
When `agreectl set config` is run  
Then `backend/config/` is created  
And `backend/config/postgres.json` is written successfully

### Scenario: --port overrides the default NodePort

Given `agreectl set config` is run with `--db-port 31000`  
Then the `port` field in `backend/config/postgres.json` is `31000`
