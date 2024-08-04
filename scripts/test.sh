#!/bin/bash
set -e

readonly env_file="$1"
# readonly service="$2"

cd "./tests"
env $(cat "../.env" "../$env_file" | grep -Ev '^#' | xargs) go test -count=1 -p=8 -parallel=8 -race ./...