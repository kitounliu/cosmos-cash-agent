#!/bin/bash

curl -d '{ "serviceEndpoint":"http://localhost:7090","recipientKeys":[ "e3xUwgT9Qjb8KamK3kmvfRfmFf6LdYZMC6SeeV2oUnV"],"@id":"90c46677-27a0-41e1-a272-68eb15bb1984","label":"alice-agent","@type":"https://didcomm.org/didexchange/1.0/invitation"}' -H "Content-Type: application/json" -X POST http://localhost:9090/connections/create-invitation

