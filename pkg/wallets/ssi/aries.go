package ssi

import (
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"
	log "github.com/sirupsen/logrus"
)


var (
	w *wallet.Wallet
)

type CredentialsWallet struct {
	w *wallet.Wallet
}

func Agent(name, pass string) *CredentialsWallet {
	framework, err := aries.New(
		aries.WithVDR(CosmosVDR{}),
	)
	if err != nil {
		panic(err)
	}
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
	return &CredentialsWallet{w: w}

}

// Run should be called as a goroutine, the parameters are:
// State: the local state of the app that should be stored on disk
// Hub: is the messages where the 3 components (ui, wallet, agent) can exchange messages
func (cw *CredentialsWallet) Run(state *config.State, hub *config.MsgHub) {
	// here an example how to listen to internal messages
	for {
		m := <- hub.AgentWalletIn
		log.Infoln("received message", m)
	}
	// here an example how to send the messages to the wallet
	hub.TokenWalletIn <- "add verification method xyz"

	// here an example how to send a notification to the ui
	hub.Notification <- "connection with bob agent established, bob is now a contact"

}

// AcceptContactRequest
// SendContactRequest
// AcceptVC
// RequestVC
