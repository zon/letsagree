# agreectl set config — Humanity Protocol

## Purpose

`agreectl set config` provisions `backend/config/humanity-protocol.yaml` for local development as part of its single-command config setup. Both when an env file is provided and when pulling credentials from the existing `humanity-protocol` secret, it writes the full YAML config back into the `humanity-protocol` Kubernetes secret in the ralph namespace so workflow pods can mount it as a config file.

## Requirements

- `set config` MUST accept the following additional flags for Humanity Protocol configuration:

  | Flag                | Default                               | Description                                                  |
  |---------------------|---------------------------------------|--------------------------------------------------------------|
  | `--ralph-namespace` | `ralph-letsagree`                     | Kubernetes namespace for the HP secret                       |
  | `--hp-secret`       | `humanity-protocol`                   | Name of the Humanity Protocol secret in Ralph's namespace    |
  | `--hp-env`          | _(none)_                              | Path to an env file containing Humanity Protocol credentials |
  | `--oidc-issuer`     | `https://api.sandbox.humanity.org/v2` | OIDC issuer base URL                                         |
  | `--oidc-redirect`   | _(required)_                          | Absolute callback URL registered with the provider           |

- The `--context` flag from the existing `set config` command MUST also be used when accessing the HP secret.

### Env file parsing

- When `--hp-env` is provided, the command MUST parse the file for the following keys:
  - `HUMANITY_CLIENT_ID` → `clientID`
  - `HUMANITY_CLIENT_SECRET` → `clientSecret`
  - `HUMANITY_PUBLIC_KEY` → `publicKey`
- Lines beginning with `#` MUST be ignored.
- If any of the three required keys are missing from the env file, the command MUST exit with an error naming the missing keys.

### Determining the source of credentials

- When `--hp-env` is provided, credentials MUST be read from the env file.
- When `--hp-env` is not provided and the `--hp-secret` secret does not contain a `humanity-protocol.yaml` key, the command MUST exit with an error instructing the user to provide `--hp-env`.
- When `--hp-env` is not provided and the `--hp-secret` secret already contains a `humanity-protocol.yaml` key, the command MUST read credential values from that stored config.

### Writing the HP secret

- After resolving credentials from either source, the command MUST upsert the `--hp-secret` secret in `--ralph-namespace` with the key `humanity-protocol.yaml` set to the full YAML content of the humanity protocol config.
- If the secret already exists it MUST be updated (upserted), not duplicated.

### Writing the local config

- The command MUST write `backend/config/humanity-protocol.yaml` relative to the repo root, creating `backend/config/` if it does not exist.
- The written YAML MUST contain the fields: `clientID`, `clientSecret`, `publicKey`, `issuerURL`, `redirectURL`.
- `clientID`, `clientSecret`, and `publicKey` MUST come from the resolved credential source (env file or HP secret).
- `issuerURL` MUST be the value of `--oidc-issuer`.
- `redirectURL` MUST be the value of `--oidc-redirect`.

## Scenarios

### Scenario: env file provided — writes local config and HP secret

Given the env file at `backend/config/env.sandbox` contains valid `HUMANITY_CLIENT_ID`, `HUMANITY_CLIENT_SECRET`, and `HUMANITY_PUBLIC_KEY`  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the `humanity-protocol` secret in `ralph-letsagree` is upserted with key `humanity-protocol.yaml` containing the full YAML config  
And `backend/config/humanity-protocol.yaml` is written with all five fields populated

### Scenario: env file provided — upserts existing HP secret

Given the `humanity-protocol` secret already exists in the `ralph-letsagree` namespace  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the existing secret is updated, not duplicated

### Scenario: no env file and no HP data in secret — exits with error

Given the `humanity-protocol` secret does not contain a `humanity-protocol.yaml` key  
And `--hp-env` is not provided  
When `agreectl set config --oidc-redirect https://example.com/auth/callback` is run  
Then the command exits with an error instructing the user to provide `--hp-env`

### Scenario: no env file and HP data present — pulls from HP secret

Given the `humanity-protocol` secret in `ralph-letsagree` contains a `humanity-protocol.yaml` key  
And `--hp-env` is not provided  
When `agreectl set config --oidc-redirect https://example.com/auth/callback` is run  
Then `backend/config/humanity-protocol.yaml` is written using the credential values from the HP secret  
And the `humanity-protocol` secret is upserted with the updated config (reflecting current OIDC options)

### Scenario: env file missing required key — exits with error

Given the env file is missing `HUMANITY_CLIENT_SECRET`  
When `agreectl set config --hp-env backend/config/env.sandbox --oidc-redirect https://example.com/auth/callback` is run  
Then the command exits with an error naming `HUMANITY_CLIENT_SECRET` as missing

### Scenario: output directory is created if missing

Given `backend/config/` does not exist  
When the command runs successfully  
Then `backend/config/` is created and `backend/config/humanity-protocol.yaml` is written

### Scenario: custom context and ralph-namespace

Given `--context k3s-prod --ralph-namespace infra` are provided  
Then the `humanity-protocol` secret is read from and written to namespace `infra` using context `k3s-prod`
