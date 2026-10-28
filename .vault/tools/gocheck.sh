#!/usr/bin/env bash
# Usage: gocheck.sh <path-to-whatsapp-bridge> [extra go test args]
# Runs go build / vet / test in a Linux container with CGO (real sqlite3).
set -e
dir=$(cd "$1" && pwd -W)
shift
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=1 -e GOFLAGS=-buildvcs=false \
  -v gomodcache:/go/pkg/mod -v gobuildcache:/root/.cache/go-build \
  -v "$dir:/src" -w /src golang:1.27-bookworm \
  bash -c "go build ./... && go vet ./... && go test -count=1 $* ./... && echo GOCHECK_OK"
