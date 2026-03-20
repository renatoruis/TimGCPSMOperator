# Installation — TIMGCPSMOPERATOR

O bundle cria o namespace **`timgcpsm-operator-system`** (exclusivo para o controller). Não mistures outras cargas de trabalho nesse namespace.

## Quick install (release)

```bash
kubectl apply -f https://github.com/renatoruis/TimGCPSMOperator/releases/latest/download/install.yaml
```

The bundle installs:

- CRDs: `timgcpsmsecrets.secrets.tim.operator`, `timgcpsmsecretconfigs.secrets.tim.operator`
- Namespace: `timgcpsm-operator-system`
- ClusterRole `timgcpsm-operator-role`, ServiceAccount `timgcpsm-operator`, ClusterRoleBinding
- Deployment `timgcpsm-operator-controller`

## Manual install (CRDs only)

```bash
kubectl apply -f https://github.com/renatoruis/TimGCPSMOperator/releases/latest/download/timgcpsmsecret-crd.yaml
kubectl apply -f https://github.com/renatoruis/TimGCPSMOperator/releases/latest/download/timgcpsmsecretconfig-crd.yaml
```

## From this repo

```bash
kubectl apply -f config/crd/timgcpsmsecret-crd.yaml
kubectl apply -f config/crd/timgcpsmsecretconfig-crd.yaml
kubectl create namespace timgcpsm-operator-system
kubectl apply -f config/rbac/service_account.yaml
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/role_binding.yaml
kubectl apply -f config/manager/deployment.yaml
```

## Verify

```bash
kubectl get pods -n timgcpsm-operator-system
kubectl logs -n timgcpsm-operator-system deployment/timgcpsm-operator-controller
kubectl get crd | grep secrets.tim.operator
```

## Upgrade

Re-apply the same `install.yaml` for the target version (CRDs are backward-compatible within `v1alpha1` unless release notes say otherwise).

## Uninstall

```bash
kubectl delete -f https://github.com/renatoruis/TimGCPSMOperator/releases/latest/download/install.yaml
# or delete resources in reverse order; then CRDs if desired
```

## First resources

```bash
kubectl apply -f examples/timgcpsmsecretconfig-example.yaml
kubectl apply -f examples/timgcpsmsecret-with-config.yaml
kubectl get tgs
```

## Acesso ao Secret Manager (GCP)

O `kubectl apply` **não** configura permissões na GCP. Tens de configurar **Workload Identity** (GKE) ou outra forma de ADC e IAM — passo a passo em **[docs/gcp-permissoes.md](docs/gcp-permissoes.md)** (inclui exemplo de `gcloud` e anotação no `ServiceAccount` `timgcpsm-operator`).
