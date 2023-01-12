#!/usr/bin/env bash

set -eo pipefail

if ! which goreleaser >/dev/null ; then
    go install github.com/goreleaser/goreleaser@v1.9.2
fi

goreleaser build --single-target --snapshot --rm-dist
