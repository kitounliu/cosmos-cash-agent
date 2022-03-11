#!/bin/bash
	

	
echo "Starting bob agent on port: 8090"
go run ../cmd/elesto-agent/main.go start --api-host localhost:8090 --inbound-host ws@localhost:8092 \
	--inbound-host-external ws@ws://localhost:8092 \
	--outbound-transport ws \
	--webhook-url http://localhost:8082/wh/bob \
	--auto-accept true \
	--transport-return-route all \
	--agent-default-label BobAgent --database-type leveldb \
	--http-resolver-url cosmos@http://localhost:2109/identifier/aries/
