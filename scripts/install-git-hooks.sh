#!/usr/bin/env bash
# Install git hooks for TIMGCPSMOPERATOR
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
git config core.hooksPath .githooks
echo "Hooks path set to .githooks"
