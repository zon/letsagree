# agreectl set config — Humanity Protocol

## Purpose

`agreectl set config` provisions `backend/config/humanity-protocol.yaml` for local development as part of its single-command config setup. When credentials are first introduced via an env file, it also writes them as a Kubernetes secret so subsequent runs can pull the values from the cluster instead of requiring the env file again.

## Requirements

- `set config` MUST accept the following additional flags for Humanity Protocol configuration:

  | Flag                | Default                               | Description                                                  |
  |---------------------|---------------------------------------|--------------------------------------------------------------|
  | `--ralph-namespace` | `ralph-letsagree`                     | Kubernetes namespace for the Humanity Protocol secret        |
  | `--hp-secret`       | `humanity-protocol`                   | Name of the Kubernetes secret                                |
  | `--hp-env`          | _(none)_                              | Path to an env file containing Humanity Protocol credentials |
  | `--oidc-issuer`      | `https://api.sandbox.humanity.org/v2` | OIDC issuer base URL                                         |
  | `--oidc-redirect`    | _(required)_                          | Absolute callback URL registered with the provider           |

- The `--context` flag from the existing `set config` command MUST also be used when accessing the Humanity Protocol secret.

### Env file parsing

- When `--hp-env` is provided, the command MUST parse the file for the following keys:
  - `HUMANITY_CLIENT_ID` → `clientId`
  - `HUMANITY_CLIENT_SECRET` → `clientSecret`
  - `HUMANITY_PUBLIC_KEY` → `publicKey`
- Lines beginning with `#` MUST be ignored.
- If any of the three required keys are missing from the env file, the command MUST exit with an error naming the missing keys.

### Writing the Kubernetes secret

- When `--hp-env` is provided, the command MUST write a Kubernetes secret to the specified namespace containing the parsed credential values (`clientId`, `clientSecret`, `publicKey`).
- If the secret already exists it MUST be updated (upserted), not duplicated.

### Determining the source of credentials

- When `--hp-env` is not provided and the Kubernetes secret is absent from the cluster, the command MUST exit with an error instructing the user to provide `--hp-env`.
- When `--hp-env` is not provided and the Kubernetes secret is present, the command MUST read credential values from the secret.

### Writing the local config

- The command MUST write `backend/config/humanity-protocol.yaml` relative to the repo root, creating `backend/config/` if it does not exist.
- The written YAML MUST contain the fields: `clientId`, `clientSecret`, `publicKey`, `issuerUrl`, `redirectUrl`.
- `clientId`, `clientSecret`, and `publicKey` MUST come from the resolved credential source (env file or cluster secret).
- `issuerUrl` MUST be the value of `--oidc-issuer`.
- `redirectUrl` MUST be the value of `--oidc-redirect`.

## Scenarios

### Scenario: env file provided — writes local config and cluster secret

Given the env file at `backend/config/env.sandbox` contains valid `HUMANITY_CLIENT_ID`, `HUMANITY_CLIENT_SECRET`, and `HUMANITY_PUBLIC_KEY`  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the Kubernetes secret `humanity-protocol` is written to the `ralph-letsagree` namespace with the parsed credential values  
And `backend/config/humanity-protocol.yaml` is written with all five fields populated

### Scenario: env file provided — upserts existing cluster secret

Given the Kubernetes secret `humanity-protocol` already exists in the `ralph-letsagree` namespace  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the existing secret is updated, not duplicated

### Scenario: no env file and secret absent — exits with error

Given the Kubernetes secret `humanity-protocol` does not exist in the `ralph-letsagree` namespace  
And `--hp-env` is not provided  
When `agreectl set config --oidc-redirect https://example.com/auth/callback` is run  
Then the command exits with an error instructing the user to provide `--hp-env`

### Scenario: no env file and secret present — pulls from cluster

Given the Kubernetes secret `humanity-protocol` exists in the `ralph-letsagree` namespace  
And `--hp-env` is not provided  
When `agreectl set config --oidc-redirect https://example.com/auth/callback` is run  
Then `backend/config/humanity-protocol.yaml` is written using the credential values from the cluster secret  
And no Kubernetes secret write is performed

### Scenario: env file missing required key — exits with error

Given the env file is missing `HUMANITY_CLIENT_SECRET`  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the command exits with an error naming `HUMANITY_CLIENT_SECRET` as missing

### Scenario: output directory is created if missing

Given `backend/config/` does not exist  
When the command runs successfully  
Then `backend/config/` is created and `backend/config/humanity-protocol.yaml` is written

### Scenario: custom context, namespace, and secret name

Given `--context k3s-prod --ralph-namespace infra --hp-secret hp-creds` are provided  
Then the secret `hp-creds` is read from or written to namespace `infra` using context `k3s-prod`
