#!/bin/bash
	

	
echo "Starting bob agent on port: 8090"
dlv debug ../cmd/elesto-agent/main.go -- start --api-host localhost:8090 --inbound-host http@localhost:8091,ws@localhost:8092 \
	--inbound-host-external http@https://example.com:8091,ws@ws://localhost:8092 \
	--webhook-url http://localhost:8082/wh/bob \
	--agent-default-label BobAgent --database-type leveldb \
	--http-resolver-url cosmos@http://localhost:2109/identifier/aries/
