package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	de "github.com/hyperledger/aries-framework-go/pkg/controller/command/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
)

// Invitation model for DID Exchange invitation.
type Invitation struct {
	*didexchange.Invitation
}

func request(client *http.Client, method, url string, requestBody io.Reader, val interface{}) {
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	json.Unmarshal(bodyBytes, &val)
	fmt.Printf("---> Request URL:\n %s\n<--- Reply:\n%s\n", url, bodyBytes)
}

func post(client *http.Client, url string, requestBody io.Reader, val interface{}) {
	request(client, "POST", url, requestBody, val)
}
func get(client *http.Client, url string, val interface{}) {
	request(client, "GET", url, nil, val)
}

func main() {

	// DID Exchange
	// https://github.com/hyperledger/aries-framework-go/blob/main/docs/rest/openapi_demo.md#steps-for-didexchange

	var (
		bobAgent = "http://localhost:8090"
		bobDID     = "did:cosmos:net:cash:bob"
		aliceAgent = "http://localhost:7090"
		aliceDID     = "did:cosmos:net:cash:alice"
		url        string
		client     = &http.Client{}
		msg []byte
		err error
	)

	var(
		connection de.QueryConnectionResponse
		connections de.QueryConnectionsResponse
		invite de.CreateInvitationResponse
		receiveInvite de.ReceiveInvitationResponse
		acceptInvite de.AcceptInvitationResponse
		confirmExchange de.ExchangeResponse
	)

	println("DID Exchange")
	//x := de.ImplicitInvitationArgs{
	//	InviterDID:        aliceDID,
	//	InviterLabel:      "AliceAgent",
	//	InviteeDID:        bobDID,
	//	InviteeLabel:      "BobAgent",
	//}
	//
	//msg, err = json.Marshal(x)
	//if err != nil {
	//	panic(err)
	//}

	println("ALICE", aliceDID)
	println("BOB  ", bobDID)
	routerID := fmt.Sprint(rand.Int())

	// Create invitation
	url = fmt.Sprint(aliceAgent, "/connections/create-invitation?public=", bobDID, "&alias=hellothere&router_connection_id=",routerID)
	println("1. ALICE creates an invitation", url)
	post(client, url, nil, &invite)

	msg, err = json.Marshal(invite)
	if err != nil {
		panic(err)
	}

	url = fmt.Sprint(bobAgent, "/connections/receive-invitation")
	println("2. BOB receive the invitation", url)

	post(client, url, bytes.NewBuffer(msg), &receiveInvite)

	// Check connection
	url = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID)
	println("3. BOB inspect the invitation", url)
	get(client, url, &connection)


	// Check connection
	url = fmt.Sprint(bobAgent, "/connections")
	println("4. BOB lists connections ", url)
	get(client, url, &connections)


	url = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID, "/accept-invitation")
	println("5. BOB accepts the connection", url)
	//var accept de.AcceptInvitationResponse
	post(client, url, nil, &acceptInvite)

	// Check connection
	url = fmt.Sprint(aliceAgent, "/connections")
	println("6. ALICE lists connections", url)
	get(client, url, &connections)

	url = fmt.Sprint(aliceAgent, "/connections/", receiveInvite.ConnectionID, "/accept-request")
	println("7. ALICE accepts the connection request (replied from bob)", url)
	post(client, url, nil, &confirmExchange)


	// Check connection
	url = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID)
	println("8.1 BOB inspect the connection", url)
	get(client, url, &connection)

	url = fmt.Sprint(aliceAgent, "/connections/", receiveInvite.ConnectionID)
	println("8.2 ALICE inspect the connection", url)
	get(client, url, &connection)


	print("yey!")
}
