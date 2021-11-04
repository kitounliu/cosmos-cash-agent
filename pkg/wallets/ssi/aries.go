package ssi

import (
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"
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

func (cw *CredentialsWallet) Run(state *config.State) {

}

// AcceptContactRequest
// SendContactRequest
// AcceptVC
// RequestVC
