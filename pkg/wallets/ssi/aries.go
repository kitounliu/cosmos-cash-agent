package ssi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/messaging"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/messaging/msghandler"

	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/mediator"
	de "github.com/hyperledger/aries-framework-go/pkg/controller/command/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"

	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/httpbinding"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"

	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"

	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"

	log "github.com/sirupsen/logrus"
)

var (
	w      *wallet.Wallet
	client = &http.Client{}
	reqURL string
)

type SSIWallet struct {
	w                 *wallet.Wallet
	ctx               *context.Provider
	didExchangeClient *didexchange.Client
	routeClient       *mediator.Client
	messagingClient   *messaging.Client
}

func createDIDExchangeClient(ctx *context.Provider) *didexchange.Client {
	// create a new did exchange client
	didExchange, err := didexchange.New(ctx)
	if err != nil {
		panic(err)
	}

	actions := make(chan service.DIDCommAction, 1)

	err = didExchange.RegisterActionEvent(actions)
	if err != nil {
		panic(err)
	}

	go func() {
		service.AutoExecuteActionEvent(actions)
	}()

	return didExchange
}

func createRoutingClient(ctx *context.Provider) *mediator.Client {
	// create the mediator client this client handler routing between edge and cloud agents
	routeClient, err := mediator.New(ctx)
	if err != nil {
		panic(err)
	}
	events := make(chan service.DIDCommAction)

	err = routeClient.RegisterActionEvent(events)
	if err != nil {
		panic(err)
	}
	go func() {
		service.AutoExecuteActionEvent(events)
	}()

	return routeClient
}

func createMessagingClient(ctx *context.Provider) *messaging.Client {
	n := LocalNotifier{}
	registrar := msghandler.NewRegistrar()

	msgClient, err := messaging.New(ctx, registrar, n)
	if err != nil {
		panic(err)
	}

	return msgClient

}

func Agent(name, pass, resolverURL string) *SSIWallet {
	// datastore
	provider := mem.NewProvider()
	stateProvider := mem.NewProvider()

	// ws outbound
	var transports []transport.OutboundTransport
	outboundWs := ws.NewOutbound()
	transports = append(transports, outboundWs)

	// resolver
	httpVDR, err := httpbinding.New(resolverURL,
		httpbinding.WithAccept(func(method string) bool { return method == "cosmos" }))
	if err != nil {
		panic(err.Error())
	}

	// create framework
	framework, err := aries.New(
		aries.WithStoreProvider(provider),
		aries.WithProtocolStateStoreProvider(stateProvider),
		aries.WithOutboundTransports(transports...),
		aries.WithTransportReturnRoute("all"),
		aries.WithVDR(httpVDR),
	//	aries.WithVDR(CosmosVDR{}),
	)
	// get the context
	ctx, err := framework.Context()
	if err != nil {
		panic(err)
	}
	// creating wallet profile using local KMS passphrase
	err = wallet.CreateProfile(name, ctx, wallet.WithPassphrase(pass))
	if err != nil {
		panic(err)
	}

	// creating vcwallet instance for user with local KMS settings.
	w, err = wallet.New(name, ctx)
	if err != nil {
		panic(err)
	}

	didExchangeClient := createDIDExchangeClient(ctx)
	routeClient := createRoutingClient(ctx)
	messagingClient := createMessagingClient(ctx)

	return &SSIWallet{
		w:                 w,
		ctx:               ctx,
		didExchangeClient: didExchangeClient,
		routeClient:       routeClient,
		messagingClient:   messagingClient,
	}
}

func (cw *SSIWallet) HandleInvitation(
	invitation *de.CreateInvitationResponse,
) *didexchange.Connection {
	connectionID, err := cw.didExchangeClient.HandleInvitation(invitation.Invitation)
	if err != nil {
		panic(err)
	}

	connection, err := cw.didExchangeClient.GetConnection(connectionID)
	if err != nil {
		panic(err)
	}
	log.Infoln("Connection created", connection)

	return connection

}

// Run should be called as a goroutine, the parameters are:
// State: the local state of the app that should be stored on disk
// Hub: is the messages where the 3 components (ui, wallet, agent) can exchange messages
func (cw *SSIWallet) Run(state *config.State, hub *config.MsgHub) {
	// here an example how to listen to internal messages
	for {
		m := <-hub.AgentWalletIn
		switch m.Typ {
		case config.MsgHandleInvitation:
			log.Debugln(
				"TokenWallet received MsgHandleInvitation msg for ",
				m.Payload.(string),
			)
			var invite de.CreateInvitationResponse
			reqURL = fmt.Sprint(
				"http://localhost:8090",
				"/connections/create-invitation?&label=BobMediatorEdgeAgent",
			)
			post(client, reqURL, nil, &invite)

			// TODO: validate invitation is correct
			connection := cw.HandleInvitation(&invite)

			hub.Notification <- config.NewAppMsg(config.MsgHandleInvitation, connection)
		}

	}
}

// TODO remove in favor of public did exchange here for test purposes
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
func bitify(in interface{}) io.Reader {
	v, err := json.Marshal(in)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewBuffer(v)
}

// AcceptContactRequest
// SendContactRequest
// AcceptVC
// RequestVC
