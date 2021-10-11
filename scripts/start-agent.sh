#!/bin/bash

echo "Starting alice agent on port: 7090"
go run ../cmd/elesto-agent/main.go start --api-host localhost:7090 --inbound-host http@localhost:7091,ws@localhost:7092 --inbound-host-external http@https://example.com:7091,ws@ws://localhost:7092 --webhook-url http://localhost:7082 --agent-default-label AliceAgent --database-type mem --http-resolver-url cosmos@http://localhost:2109/identifier/ &

echo "Starting bob agent on port: 8090"
go run ../cmd/elesto-agent/main.go start --api-host localhost:8090 --inbound-host http@localhost:8091,ws@localhost:8092 --inbound-host-external http@https://example.com:8091,ws@ws://localhost:8092 --webhook-url localhost:8082 --agent-default-label BobAgent --database-type leveldb --http-resolver-url cosmos@http://localhost:2109/identifier/
