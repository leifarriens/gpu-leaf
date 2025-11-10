#!/usr/bin/env bash
set -euo pipefail

mkdir -p bin

# Determine version (can be overridden by env var VERSION)
VERSION=${VERSION:-$(git describe --tags --always 2>/dev/null || echo dev)}
LDFLAGS="-X github.com/leifarriens/gpu-leaf/internal/version.Version=${VERSION}"

GO111MODULE=on go build -ldflags "$LDFLAGS" -o bin/gpu-leaf ./cmd/gpuleaf