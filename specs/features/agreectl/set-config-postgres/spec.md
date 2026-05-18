# agreectl set config — Postgres

## Purpose

`agreectl set config` writes postgres config to two destinations, each adjusted for its environment: a local file for dev machine access and a Kubernetes secret for in-cluster workflow pods.

## Requirements

- `agreectl` MUST be a Go CLI tool using [Kong](https://github.com/alecthomas/kong).
- `agreectl` MUST expose a `set config` subcommand.
- `set config` MUST accept the following flags:

  | Flag                | Default           | Description                                              |
  |---------------------|-------------------|----------------------------------------------------------|
  | `--context`         | `microk8s`        | kubectl context to use                                   |
  | `--namespace`       | `letsagree`       | Kubernetes namespace containing the CNPG secret          |
  | `--ralph-namespace` | `ralph-letsagree` | Kubernetes namespace where the ralph workflow runs       |
  | `--db-secret`       | `letsagree-app`   | Name of the CNPG app secret                              |
  | `--db-port`         | `30432`           | NodePort for local postgres access                       |
  | `--postgres-secret` | `postgres`        | Name of the postgres secret in Ralph's namespace         |

- `set config` MUST read the following fields from the CNPG app secret: `host`, `port`, `user`, `password`, `dbname`.
- `set config` MUST write `backend/config/postgres.json` with values adjusted for local access:
  - `host`: internal IP of the first ready cluster node (auto-detected)
  - `port`: the `--db-port` NodePort value
  - `user`, `password`, `dbname`: copied from the secret unchanged
- `set config` MUST upsert the `--postgres-secret` secret in `--ralph-namespace` with values adjusted for in-cluster access:
  - `host`: the secret's `host` qualified with the source namespace (e.g. `letsagree-rw.letsagree`)
  - `port`: the secret's `port` value (ClusterIP port)
  - `user`, `password`, `dbname`: copied from the secret unchanged

## Scenarios

### Scenario: local file uses NodeIP and NodePort

Given the secret `letsagree-app` exists in namespace `letsagree`  
And the first cluster node has InternalIP `192.168.1.10`  
When `agreectl set config` is run with no flags  
Then `backend/config/postgres.json` is written with `host=192.168.1.10` and `port=30432`

### Scenario: cluster secret uses qualified host and ClusterIP port

Given the secret `letsagree-app` has `host=letsagree-rw` and `port=5432`  
When `agreectl set config` is run  
Then the `postgres` secret in `ralph-letsagree` is upserted with `host=letsagree-rw.letsagree` and `port=5432`

### Scenario: custom context, namespace, and secret

Given `agreectl set config` is run with `--context k3s-prod --namespace myns --db-secret myns-app`  
Then the secret `myns-app` is fetched from namespace `myns` using context `k3s-prod`

### Scenario: custom ralph-namespace

Given `agreectl set config` is run with `--ralph-namespace infra`  
Then the postgres secret is upserted into namespace `infra`

### Scenario: output directory is created if missing

Given `backend/config/` does not exist  
When `agreectl set config` is run  
Then `backend/config/` is created  
And `backend/config/postgres.json` is written successfully

### Scenario: --db-port overrides the default NodePort

Given `agreectl set config` is run with `--db-port 31000`  
Then `backend/config/postgres.json` has `port=31000`
