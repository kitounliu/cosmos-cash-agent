#!/bin/bash

go run ../cmd/elesto-agent/main.go start --api-host localhost:9090 --inbound-host http@localhost:9091,ws@localhost:9092 --inbound-host-external http@https://example.com:9091,ws@ws://localhost:9092 --webhook-url localhost:8082 --agent-default-label MyAgent --database-type leveldb
