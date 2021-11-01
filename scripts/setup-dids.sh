#!/bin/bash



cosmos-cashd tx did create-did alice --from alice --chain-id cash -y
sleep 5
cosmos-cashd tx did link-aries-agent alice http://localhost:7090 http://localhost:7091 --from alice --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq


cosmos-cashd tx did create-did bob --from bob --chain-id cash -y
sleep 5
cosmos-cashd tx did link-aries-agent bob http://localhost:8090 http://localhost:8091 --from bob --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq


# Used for did message routing

cosmos-cashd tx did create-did carl --from carl --chain-id cash -y
sleep 5
cosmos-cashd tx did link-aries-agent carl http://localhost:10090 http://localhost:10091 --from carl --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq

cosmos-cashd tx did create-did dave --from dave --chain-id cash -y
sleep 5
cosmos-cashd tx did link-aries-agent dave http://localhost:11090 http://localhost:11091 --from dave --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq
