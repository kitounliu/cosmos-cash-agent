package ssi

import (
	"net/http"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"

	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"
	arieshttp "github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/http"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/httpbinding"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"

	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"

	log "github.com/sirupsen/logrus"
)

var (
	w *wallet.Wallet
)

type SSIWallet struct {
	w   *wallet.Wallet
	ctx *context.Provider
}

func Agent(name, pass, resolverURL string) *SSIWallet {
	// datastore
	provider := mem.NewProvider()
	stateProvider := mem.NewProvider()

	// ws inbound, outbound
	var transports []transport.OutboundTransport
	outboundHTTP, err := arieshttp.NewOutbound(arieshttp.WithOutboundHTTPClient(&http.Client{}))
	transports = append(transports, outboundHTTP)
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
		aries.WithVDR(CosmosVDR{}),
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
	return &SSIWallet{w: w, ctx: ctx}
}

// Run should be called as a goroutine, the parameters are:
// State: the local state of the app that should be stored on disk
// Hub: is the messages where the 3 components (ui, wallet, agent) can exchange messages
func (cw *SSIWallet) Run(state *config.State, hub *config.MsgHub) {
	// here an example how to listen to internal messages
	for {
		m := <-hub.AgentWalletIn
		log.Infoln("received message", m)
	}
	// here an example how to send the messages to the wallet
	//hub.TokenWalletIn <- config.AppMsg{config.MsgBalances, "add verification method xyz"}

	// here an example how to send a notification to the ui
	// hub.Notification <- "connection with bob agent established, bob is now a contact"

}

// AcceptContactRequest
// SendContactRequest
// AcceptVC
// RequestVC
