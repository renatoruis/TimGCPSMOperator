#!/bin/bash
set -e

echo "🔧 Generating CRDs from Go types..."

# Check if controller-gen is available
if ! command -v controller-gen &> /dev/null; then
    echo "📦 Installing controller-gen..."
    go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
fi

# Generate CRDs
echo "📝 Running controller-gen..."
controller-gen crd:crdVersions=v1 paths="./api/..." output:crd:artifacts:config=config/crd

# Rename generated files to our naming convention
if [ -f "config/crd/secrets.tim.operator_timgcpsmsecrets.yaml" ]; then
    mv config/crd/secrets.tim.operator_timgcpsmsecrets.yaml config/crd/timgcpsmsecret-crd.yaml
    echo "✅ Generated timgcpsmsecret-crd.yaml"
fi

if [ -f "config/crd/secrets.tim.operator_timgcpsmsecretconfigs.yaml" ]; then
    mv config/crd/secrets.tim.operator_timgcpsmsecretconfigs.yaml config/crd/timgcpsmsecretconfig-crd.yaml
    echo "✅ Generated timgcpsmsecretconfig-crd.yaml"
fi

echo "✅ CRDs generated successfully!"

