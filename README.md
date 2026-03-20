# TIMGCPSMOPERATOR (timgcpsm-operator)

[![GitHub Release](https://img.shields.io/github/v/release/renatoruis/TimGCPSMOperator)](https://github.com/renatoruis/TimGCPSMOperator/releases)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![CI](https://github.com/renatoruis/TimGCPSMOperator/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/renatoruis/TimGCPSMOperator/actions/workflows/ci.yaml)

Kubernetes operator **only for [Google Cloud Secret Manager](https://cloud.google.com/secret-manager)**. It syncs secrets into Kubernetes `Secret` objects on a configurable interval, detects changes (including updates made in the **GCP console** or **Secret Manager API**), and can **restart a `Deployment`** when the synced data changes.

This is intentionally **not** External Secrets Operator: one backend, one workflow, and rollout-on-change as a first-class feature.

## Features

- **GCP Secret Manager only** — reads secret versions via Application Default Credentials (Workload Identity on GKE).
- **Centralized project config** — optional `TimGcpSmSecretConfig` with `projectId` for many `TimGcpSmSecret` resources.
- **Polling + hash** — periodic sync (`syncInterval`, default `5m`) so changes outside the cluster are picked up.
- **Deployment restart** — optional `deploymentName`; rollout when secret payload changes.
- **Payload modes** — `decodeFormat: text` (single key, default key `value`) or `json` (JSON object → multiple keys).
- **Status** — conditions, `secretHash`, retries, last error.

## Requirements

- GKE (or GCP) com identidade GCP para o pod — em GKE o habitual é [**Workload Identity**](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity). **O cluster não ganha acesso ao Secret Manager só por existir:** é preciso IAM na GCP (ver [docs/gcp-permissoes.md](docs/gcp-permissoes.md)).
- **GCP IAM (só leitura):** o código usa apenas [`AccessSecretVersion`](https://cloud.google.com/secret-manager/docs/reference/libraries) — não cria, altera nem apaga segredos no GSM. Concede ao Service Account GCP usado pelo pod **`roles/secretmanager.secretAccessor`** nos segredos (ou projeto) necessários.

## Installation

O operador instala-se num **namespace dedicado** — `timgcpsm-operator-system` — definido em [`config/manager/namespace.yaml`](config/manager/namespace.yaml). O `Deployment`, `ServiceAccount` e o pod do controller ficam **só** aí; os teus `TimGcpSmSecret` e `Secret` de aplicação podem continuar noutros namespaces (via `spec.namespace` no CR).

### Release (when published)

```bash
kubectl apply -f https://github.com/renatoruis/TimGCPSMOperator/releases/latest/download/install.yaml
```

### From this repository

```bash
kubectl apply -f config/crd/timgcpsmsecret-crd.yaml
kubectl apply -f config/crd/timgcpsmsecretconfig-crd.yaml
kubectl apply -f config/manager/namespace.yaml
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/service_account.yaml
kubectl apply -f config/rbac/role_binding.yaml
kubectl apply -f config/manager/deployment.yaml
```

### Build

```bash
make build
make docker-build IMG=ghcr.io/renatoruis/timgcpsm-operator:v1.0.0
```

## Quick start

1. **Optional:** create `TimGcpSmSecretConfig` with `spec.projectId` (see `examples/timgcpsmsecretconfig-example.yaml`).
2. Create `TimGcpSmSecret` with `secretId`, `secretName`, and either `projectId` or `gcpSmConfig` (see `examples/`).

```bash
kubectl apply -f examples/timgcpsmsecretconfig-example.yaml
kubectl apply -f examples/timgcpsmsecret-with-config.yaml
```

## CRDs

| Resource | Short name | Purpose |
|----------|------------|---------|
| `TimGcpSmSecret` | `tgs` | Sync one GSM secret into a Kubernetes `Secret` |
| `TimGcpSmSecretConfig` | `tgsc` | Default GCP `projectId` |

**Scope:** ambos os CRDs são **sempre namespaced** (`ClusterRole` só no RBAC do operador). **Não existe** recurso cluster-wide tipo “config global único”.

Para **um único `projectId` partilhado por todo o cluster**, cria **um** `TimGcpSmSecretConfig` num namespace central (ex. `timgcpsm-operator-system`) e, em cada `TimGcpSmSecret` noutros namespaces, usa `gcpSmConfig: <nome>` e `gcpSmConfigNamespace: timgcpsm-operator-system`. Alternativa: define `projectId` diretamente em cada `TimGcpSmSecret` sem config partilhado.

O **operador** corre num namespace fixo; os **CRs** `TimGcpSmSecret` / `TimGcpSmSecretConfig` podem existir em **qualquer** namespace (conforme RBAC dos utilizadores).

### `TimGcpSmSecret` spec (summary)

| Field | Description |
|-------|-------------|
| `projectId` | GCP project (omit if using `gcpSmConfig`) |
| `gcpSmConfig` / `gcpSmConfigNamespace` | Reference to `TimGcpSmSecretConfig` |
| `secretId` | Secret Manager secret id (short name) |
| `secretVersion` | Version id or `latest` (default) |
| `secretName` | Target Kubernetes `Secret` name |
| `deploymentName` | Optional deployment to roll on change |
| `namespace` | Namespace for `Secret` / `Deployment` (defaults to the CR’s namespace) |
| `syncInterval` | e.g. `30s`, `5m` (min 30s, max 1h; default 5m) |
| `decodeFormat` | `text` or `json` |
| `secretKey` | Key for `text` mode (default `value`) |

## Operations

```bash
kubectl get tgs -A
kubectl describe tgs example-app-secrets
kubectl logs -n timgcpsm-operator-system deployment/timgcpsm-operator-controller -f
```

## Capacity

See [CAPACITY.md](CAPACITY.md) for tuning notes (e.g. `MaxConcurrentReconciles: 10`).

## Community

- [Contributing](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security policy](SECURITY.md) — report vulnerabilities **privately** (not via public issues)
- [Checklist de configuração do repositório (GitHub)](docs/github-configuracao.md) *(mantenedores)*

## License

Licensed under the [Apache License 2.0](LICENSE). See [NOTICE](NOTICE) for attribution.
