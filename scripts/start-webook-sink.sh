#!/bin/bash

echo "Starting webhook sink"
#go run ../cmd/elesto-agent/main.go start \
go run ../pkg/httpsink/main.go -listen :7082

