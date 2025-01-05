#!/bin/bash
set -e

go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

max_attempts=4

for attempt in $(seq 1 $max_attempts); do
    echo "$attempt of $max_attempts..."
    fieldalignment --json --fix pkg/domain/models/*.go
done
