#!/bin/bash
set -e

go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
fieldalignment --fix pkg/domain/models/*.go
