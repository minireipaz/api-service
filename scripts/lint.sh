#!/bin/bash
set -e

if [ "$2" == "-install" ]; then
  # binary will be $(go env GOPATH)/bin/golangci-lint
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1

  golangci-lint --version
  go install github.com/roblaszczak/go-cleanarch@latest
fi

# readonly service="$1"
# cd "./pkg"
golangci-lint run -v --tests=false --timeout=2m --config ./.golangci.yaml