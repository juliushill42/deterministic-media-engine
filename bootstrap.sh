#!/usr/bin/env bash
set -Eeuo pipefail
command -v go >/dev/null
mkdir -p dist
go test ./...
go vet ./...
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/dme ./cmd/dme
./scripts/smoke.sh
