package main

import (
	"bytes"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	de "github.com/hyperledger/aries-framework-go/pkg/controller/command/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
)

// Invitation model for DID Exchange invitation.
type Invitation struct {
	*didexchange.Invitation
}

func main() {
	fmt.Println("Calling API...")
	client := &http.Client{}

	// Create invitation
	req, err := http.NewRequest("POST", "http://localhost:8090/connections/create-invitation?public=did:cosmos:net:cash:emti&alias=hellothere&router_connection_id=wut", nil)
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
	var responseObject de.CreateInvitationResponse
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject.Invitation.DID)

	msg, err := json.Marshal(responseObject.Invitation)
	if err != nil {
		panic(err)
	}

	// Recieve Invitation
	req, err = http.NewRequest("POST", "http://localhost:7090/connections/receive-invitation", bytes.NewBuffer(msg))
	resp, err = client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var result de.ReceiveInvitationResponse
	json.Unmarshal(bodyBytes, &result)
	fmt.Printf("API Response receive invitation as struct %+v\n", result)

	// Check connection
	req, err = http.NewRequest("GET", "http://localhost:7090/connections/"+result.ConnectionID, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var connection interface{}
	json.Unmarshal(bodyBytes, &connection)
	fmt.Printf("API Response as struct %+v\n", connection)
}
