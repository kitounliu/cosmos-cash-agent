package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	de "github.com/hyperledger/aries-framework-go/pkg/controller/command/didexchange"
	ks "github.com/hyperledger/aries-framework-go/pkg/controller/command/kms"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command/messaging"
	presentproofcmd "github.com/hyperledger/aries-framework-go/pkg/controller/command/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
)

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
	fmt.Printf("---> Request URL:\n %s\nPayload:\n%s\n", url, requestBody)
	fmt.Printf("<--- Reply:\n%s\n", bodyBytes)
}

func post(client *http.Client, url string, requestBody, val interface{}) {
	if requestBody != nil {
		request(client, "POST", url, bitify(requestBody), val)
	} else {
		request(client, "POST", url, nil, val)
	}

}
func get(client *http.Client, url string, val interface{}) {
	request(client, "GET", url, nil, val)
}

func bitify(in interface{}) io.Reader {
	v, err := json.Marshal(in)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewBuffer(v)
}

func main() {

	var (
		bobAgent   = "http://localhost:8090"
		bobDID     = "did:cosmos:net:cash:bob"
		aliceAgent = "http://localhost:7090"
		aliceDID   = "did:cosmos:net:cash:alice"
	)

	bobConnID, aliceConnID, bobPeerDID, alicePeerDID := DIDExchange(bobAgent, bobDID, aliceAgent, aliceDID)
	println("Bob connid", bobConnID)
	println("BOB peerid", bobPeerDID)
	println("ALICE connid", aliceConnID)
	println("ALICE peerid", alicePeerDID)

	//time.Sleep(3 * time.Second)
	//	DIDMessaging(bobAgent, aliceAgent, bobConnID)

	time.Sleep(3 * time.Second)
	CredentialPresentation(bobAgent, alicePeerDID, bobDID, aliceAgent, bobPeerDID, aliceDID)

	print("yey!")
}

type VC struct {
	Context              []string `json:"@context"`
	Holder               string   `json:"holder"`
	Id                   string   `json:"@id"`
	Issuer               string   `json:"issuer"`
	Type                 []string `json:"type"`
	VerifiableCredential aCred    `json:"verifiableCredential"`
}

type aCred struct {
	Id     string `json:"id"`
	Holder string `json:"holder"`
}

func CredentialPresentation(bobAgent, bobPeerDID, bobPublicDID, aliceAgent, alicePeerDID, alicePublicDID string) {
	// Exchange a presentation through Present Proof protocol
	// https://github.com/hyperledger/aries-framework-go/blob/main/docs/rest/openapi_demo.md#how-to-exchange-a-presentation-through-the-present-proof-protocol
	message := `
**********************
CredentialPresentation
**********************
`

	fmt.Printf("%s", message)
	var (
		client = &http.Client{}
		resp   interface{}
		reqURL string
	)

	requestPresentation := presentproof.RequestPresentation{}
	var (
		acceptPresentation VC
		presentation       presentproof.Presentation
	)

	presentproofRequest := presentproofcmd.SendRequestPresentationArgs{
		MyDID:               bobPeerDID,
		TheirDID:            alicePeerDID,
		RequestPresentation: &requestPresentation,
	}

	// Request a presentation of a credential
	var requestResp presentproofcmd.SendRequestPresentationResponse

	reqURL = fmt.Sprint(aliceAgent, "/presentproof/send-request-presentation")
	println("1. Send a request presentation to ALICE", reqURL)
	post(client, reqURL, presentproofRequest, &requestResp)
	pid := requestResp.PIID
	println(requestResp.PIID)

	// Check the credential request worked
	reqURL = fmt.Sprint(bobAgent, "/presentproof/actions")
	println("2. Check the presentation request in BOBs agent", reqURL)
	get(client, reqURL, resp)

	acceptPresentation = VC{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
		},
		Holder: bobPublicDID,
		Id:     "user:cred:1",
		Type: []string{
			"VerifiablePresentation",
			"CredentialManagerPresentation",
		},
		VerifiableCredential: aCred{
			"asd",
			"asd",
		},
	}
	b, err := json.Marshal(acceptPresentation)
	if err != nil {
		panic(err)
	}

	presentation = presentproof.Presentation{
		PresentationsAttach: []decorator.Attachment{
			{Data: decorator.AttachmentData{Base64: base64.StdEncoding.EncodeToString([]byte(b))}},
		},
	}
	payload := presentproofcmd.AcceptRequestPresentationArgs{
		PIID:         pid,
		Presentation: &presentation,
	}

	// Accept the presentation request for BOB
	reqURL = fmt.Sprint(bobAgent, "/presentproof/", pid, "/accept-request-presentation")
	println("3. Confirm presentation request for BOB", reqURL)
	post(client, reqURL, payload, resp)

	// NOTE: latency between agents so wait 3 seconds
	time.Sleep(3 * time.Second)

	// Check the presentation on ALICEs agent
	reqURL = fmt.Sprint(aliceAgent, "/presentproof/actions")
	println("4. Get pid from Alices agent", reqURL)
	get(client, reqURL, resp)

	acceptPresentationPayload := presentproofcmd.AcceptPresentationArgs{
		PIID:  pid,
		Names: []string{"demo"},
	}

	// Accept the presentation on ALICEs agent
	reqURL = fmt.Sprint(aliceAgent, "/presentproof/", pid, "/accept-presentation")
	println("5. Accept the presentation by ALICE", reqURL)
	post(client, reqURL, acceptPresentationPayload, resp)

	// NOTE: latency between agents so wait 3 seconds
	time.Sleep(3 * time.Second)

	// Check the credential request worked
	reqURL = fmt.Sprint(aliceAgent, "/verifiable/presentations")
	println("6. Check the verifiable presentation in ALICEs agent", reqURL)
	get(client, reqURL, resp)
}

type genericInviteMsg struct {
	ID      string   `json:"@id"`
	Type    string   `json:"@type"`
	Purpose []string `json:"~purpose"`
	Message string   `json:"message"`
	From    string   `json:"from"`
}

func DIDMessaging(bobAgent, aliceAgent, connID string) {
	// DID Messaging
	// https://github.com/hyperledger/aries-framework-go/blob/main/docs/rest/openapi_demo.md#steps-for-custom-message-handling

	message := `
**********************
DIDComm Messaging
**********************
`

	fmt.Printf("%s", message)

	var (
		client = &http.Client{}
		reqURL string
	)

	var (
		createService messaging.RegisterMsgSvcArgs
		genericMsg    genericInviteMsg
		request       messaging.SendNewMessageArgs
	)

	// Messaging service
	createService.Type = "https://didcomm.org/generic/1.0/message"
	createService.Purpose = []string{"meeting", "appointment", "event"}
	createService.Name = "generic-invite"

	// Create a service to use for communication
	var resp interface{}
	reqURL = fmt.Sprint(aliceAgent, "/message/register-service")
	println("1. ALICE creates a service for BOB to send messages", reqURL)
	post(client, reqURL, createService, resp)

	// Check the service has been created
	reqURL = fmt.Sprint(aliceAgent, "/message/services")
	println("2. ALICE verifies the service has been created", reqURL)
	get(client, reqURL, resp)

	genericMsg.ID = "12123123213213"
	genericMsg.Type = "https://didcomm.org/generic/1.0/message"
	genericMsg.Purpose = []string{"meeting"}
	genericMsg.Message = "fight me you coward"
	genericMsg.From = "Bob"

	rawBytes, _ := json.Marshal(genericMsg)

	request.ConnectionID = connID
	request.MessageBody = rawBytes

	// send a message to the previously created service
	reqURL = fmt.Sprint(bobAgent, "/message/send")
	println("3. BOB sends a message of type generic invite to ALICE", reqURL)
	post(client, reqURL, request, resp)
}

func DIDExchange(bobAgent, bobDID, aliceAgent, aliceDID string) (string, string, string, string) {
	// DID Exchange
	// https://github.com/hyperledger/aries-framework-go/blob/main/docs/rest/openapi_demo.md#steps-for-didexchange
	message := `
**********************
DID Exchange
**********************
`

	fmt.Printf("%s", message)

	var (
		client = &http.Client{}
		reqURL string
		params url.Values
	)

	var (
		keySetRsp      ks.CreateKeySetResponse
		connection     de.QueryConnectionResponse
		bobconnections de.QueryConnectionsResponse
		connections    de.QueryConnectionsResponse
		invite         de.CreateInvitationResponse
		//implicitInvite de.ImplicitInvitationResponse
		receiveInvite   de.ReceiveInvitationResponse
		acceptInvite    de.AcceptInvitationResponse
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
	println("router id", routerID)

	v, _ := base64.StdEncoding.DecodeString(keySetRsp.PublicKey)
	println("keyID", keySetRsp.KeyID)
	println("keyPub (base64)", keySetRsp.PublicKey)
	println("keyPub (hex)", hex.EncodeToString(v))

	// Create invitation
	params = url.Values{}
	params.Add("public", aliceDID)
	reqURL = fmt.Sprint(aliceAgent, "/connections/create-invitation?public=", aliceDID, "&label=AliceAgent")
	println("1. ALICE creates an invitation", reqURL)
	post(client, reqURL, nil, &invite)

	//params = url.Values{}
	//params.Add("their_did", aliceDID)
	//params.Add("their_label", "AliceAgent")
	//params.Add("their_did", bobDID)
	//params.Add("their_did", "BobAgent")
	//reqURL = fmt.Sprint(aliceAgent, "/connections/create-implicit-invitation?", params.Encode())
	//println("1. ALICE creates an implicit invitation", reqURL)
	//post(client, reqURL, nil, &implicitInvite)

	reqURL = fmt.Sprint(bobAgent, "/connections/receive-invitation")
	println("2. BOB receive the invitation", reqURL)
	post(client, reqURL, invite.Invitation, &receiveInvite)

	// Check connection
	reqURL = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID)
	println("3. BOB inspect the invitation", reqURL)
	get(client, reqURL, &connection)

	// Check connection
	reqURL = fmt.Sprint(bobAgent, "/connections")
	println("4. BOB lists connections ", reqURL)
	get(client, reqURL, &bobconnections)

	reqURL = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID, "/accept-invitation")
	println("5. BOB accepts the connection", reqURL)
	//var accept de.AcceptInvitationResponse
	post(client, reqURL, nil, &acceptInvite)

	// Check connection
	reqURL = fmt.Sprint(aliceAgent, "/connections")
	println("6. ALICE lists connections", reqURL)
	get(client, reqURL, &connections)

	var aliceConnID, alicePeerDID, bobPeerDID string
	for _, c := range connections.Results {
		if c.State == "requested" {
			reqURL = fmt.Sprint(aliceAgent, "/connections/", c.ConnectionID, "/accept-request")
			println("7. ALICE accepts the connection request (replied from bob)", reqURL)
			post(client, reqURL, nil, &confirmExchange)

			reqURL = fmt.Sprint(bobAgent, "/connections/", receiveInvite.ConnectionID)
			println("8.1 BOB get connection", receiveInvite.ConnectionID)
			//var accept de.AcceptInvitationResponse
			get(client, reqURL, &connection)
			println("8.1 Connection state", connection.Result.State)

			reqURL = fmt.Sprint(aliceAgent, "/connections/", c.ConnectionID)
			println("8.2 ALICE get connection", c.ConnectionID)
			//var accept de.AcceptInvitationResponse
			get(client, reqURL, &connection)
			aliceConnID = connection.Result.ConnectionID
			alicePeerDID = connection.Result.MyDID
			bobPeerDID = connection.Result.TheirDID
		}
	}

	return receiveInvite.ConnectionID, aliceConnID, bobPeerDID, alicePeerDID
}
