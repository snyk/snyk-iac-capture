#!/usr/bin/env bash

set -eo pipefail

if ! which goreleaser >/dev/null ; then
    go install github.com/goreleaser/goreleaser@v1.14.1
fi

# Check configuration
goreleaser check

FLAGS=""
FLAGS+="--rm-dist --parallelism 2 "

# Only CI system should publish artifacts
if [ "$CI" != true ]; then
    FLAGS+="--skip-announce "
    FLAGS+="--skip-publish "
    FLAGS+="--snapshot "
fi

CMD="goreleaser release ${FLAGS}"

echo "+ Using goreleaser"
echo "+ CMD=${CMD}"

$CMD
