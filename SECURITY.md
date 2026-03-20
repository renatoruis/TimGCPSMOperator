# Security policy

## Supported versions

We release security fixes for the **latest minor release** on the default branch (`main`) and the **latest published tag** when applicable. Use the newest `v*` release for production.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |

## Reporting a vulnerability

**Please do not** open a public issue for security vulnerabilities.

1. Use **[GitHub Private vulnerability reporting](https://github.com/renatoruis/TimGCPSMOperator/security/advisories/new)** (recommended), if enabled on the repository; or  
2. Contact the maintainers through a **private channel** (Security Advisory or direct message to maintainers).

Include:

- Description of the issue and impact
- Steps to reproduce (if safe to share)
- Affected versions (tag or commit), if known

We aim to acknowledge reports within a few business days and coordinate disclosure after a fix is available.

## Operator security notes

This operator syncs data from **Google Cloud Secret Manager** into **Kubernetes Secrets**. It only calls **read** APIs (`AccessSecretVersion`); it does **not** create, update, or delete secrets in GSM.

Hardening:

- Use **Workload Identity** (or equivalent) with least privilege: **`roles/secretmanager.secretAccessor`** scoped to the specific secrets (or projects) you reference in `TimGcpSmSecret` — never Secret Admin.
- Step-by-step (GKE, `gcloud`, KSA annotation): [docs/gcp-permissoes.md](docs/gcp-permissoes.md).
- Restrict who can create or edit `TimGcpSmSecret` / `TimGcpSmSecretConfig` CRDs (Kubernetes RBAC).
- Treat `Secret` objects in the cluster as sensitive; use encryption at rest and network policies as appropriate.

## Disclosure policy

We follow coordinated disclosure: we will credit reporters who wish to be named unless anonymity is requested.
