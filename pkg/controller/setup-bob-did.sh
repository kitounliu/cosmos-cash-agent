#!/bin/bash



cosmos-cashd tx did create-did bob --from bob --chain-id cash
sleep 5
cosmos-cashd query did dids --output json | jq

cosmos-cashd tx did set-verification-relationship bob $(cosmos-cashd keys show bob -a)  --relationship keyAgreement  --from bob --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq

cosmos-cashd tx did add-service bob agent DIDCommMessaging "http://localhost:8091" --from bob --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq
