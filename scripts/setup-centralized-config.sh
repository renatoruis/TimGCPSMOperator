#!/usr/bin/env bash
# Example: create TimGcpSmSecretConfig with a default GCP project id
set -euo pipefail

PROJECT_ID="${1:-}"
if [[ -z "$PROJECT_ID" ]]; then
  echo "Usage: $0 <gcp-project-id>"
  exit 1
fi

kubectl apply -f - <<EOF
apiVersion: secrets.tim.operator/v1alpha1
kind: TimGcpSmSecretConfig
metadata:
  name: gcp-default
  namespace: default
spec:
  projectId: ${PROJECT_ID}
EOF

echo "Created TimGcpSmSecretConfig gcp-default. Reference it from TimGcpSmSecret.spec.gcpSmConfig."
