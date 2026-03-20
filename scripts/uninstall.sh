#!/usr/bin/env bash
set -euo pipefail

echo "TIMGCPSMOPERATOR — uninstall"

kubectl delete timgcpsmsecrets --all --all-namespaces 2>/dev/null || true
kubectl delete timgcpsmsecretconfigs --all --all-namespaces 2>/dev/null || true

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
kubectl delete -f "$ROOT/config/manager/deployment.yaml" 2>/dev/null || true
kubectl delete -f "$ROOT/config/rbac/role_binding.yaml" 2>/dev/null || true
kubectl delete -f "$ROOT/config/rbac/role.yaml" 2>/dev/null || true
kubectl delete -f "$ROOT/config/rbac/service_account.yaml" 2>/dev/null || true
kubectl delete -f "$ROOT/config/crd/timgcpsmsecret-crd.yaml" 2>/dev/null || true
kubectl delete -f "$ROOT/config/crd/timgcpsmsecretconfig-crd.yaml" 2>/dev/null || true
kubectl delete namespace timgcpsm-operator-system 2>/dev/null || true

echo "Uninstall finished."
