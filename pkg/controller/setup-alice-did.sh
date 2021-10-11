#!/bin/bash



cosmos-cashd tx did create-did alice --from alice --chain-id cash
sleep 5
cosmos-cashd query did dids --output json | jq

cosmos-cashd tx did set-verification-relationship alice $(cosmos-cashd keys show alice -a)  --relationship keyAgreement  --from alice --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq

cosmos-cashd tx did add-service alice agent DIDCommMessaging "http://localhost:7091" --from alice --chain-id cash -y
sleep 5
cosmos-cashd query did dids --output json | jq
