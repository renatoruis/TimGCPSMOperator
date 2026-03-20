#!/usr/bin/env bash
set -euo pipefail

echo "TIMGCPSMOPERATOR — quick start"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

kubectl apply -f config/crd/timgcpsmsecret-crd.yaml
kubectl apply -f config/crd/timgcpsmsecretconfig-crd.yaml
kubectl apply -f config/crd/timgcpsmclusterconfig-crd.yaml
kubectl apply -f config/manager/namespace.yaml
kubectl apply -f config/rbac/service_account.yaml
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/role_binding.yaml

IMAGE_NAME="${IMAGE_NAME:-timgcpsm-operator:latest}"
if command -v docker &>/dev/null; then
  docker build -t "$IMAGE_NAME" .
  kind load docker-image "$IMAGE_NAME" 2>/dev/null || true
fi

kubectl set image deployment/timgcpsm-operator-controller \
  manager="$IMAGE_NAME" \
  -n timgcpsm-operator-system 2>/dev/null || \
  sed "s|image: timgcpsm-operator:latest|image: $IMAGE_NAME|g" config/manager/deployment.yaml | kubectl apply -f -

echo "Done. See examples/timgcpsmsecret-*.yaml"
