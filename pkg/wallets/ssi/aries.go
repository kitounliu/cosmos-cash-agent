package ssi

import (
	"net/http"
	"time"

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
func (cw *SSIWallet) Run(hub *config.MsgHub) {

	// send updates about verifiable credentials
	t0 := time.NewTicker(30 * time.Second)
	go func() {
		for {
			log.Infoln("ticker! retrieving verifiable credentials")
			vcs := []string{}
			hub.Notification <- config.NewAppMsg(config.MsgVCs, vcs)
			<-t0.C
		}
	}()

	// send updates about contacts
	t1 := time.NewTicker(30 * time.Second)
	go func() {
		for {
			// TODO handle contacts
			<-t1.C
		}
	}()


	// here an example how to listen to internal messages
	for {
		m := <-hub.AgentWalletIn
		log.Debugln("received message", m)
		switch m.Typ {
		case config.MsgVCData:
			vcID := m.Payload.(string)
			// TODO: retrieve the verifiable credential
			// vc := cc.GetPublicVC(m.Payload.(string))
			log.Debugln("AgentWallet received MsgVCData msg for ", vcID)
			vc := struct{}{} // <-- fake credential
			// always send to the notification channel for the UI
			// handle the notification in the ui/handlers.go dispatcher function
			hub.Notification <- config.NewAppMsg(m.Typ, vc)
		}
	}
}

// AcceptContactRequest
// SendContactRequest
// AcceptVC
// RequestVC
