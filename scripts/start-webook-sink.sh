#!/bin/bash

echo "Starting webhook sink"
#go run ../cmd/elesto-agent/main.go start \
go run ../cmd/webhook-relayer/main.go -listen :7082

