package ui

import (
	"encoding/json"
	"fyne.io/fyne/v2/widget"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	vcTypes "github.com/allinbits/cosmos-cash/v2/x/verifiable-credential/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	log "github.com/sirupsen/logrus"
)

// dispatcher this reads notifications and updates the
// data binding
func dispatcher(in chan config.AppMsg) {

	for {
		m := <-in
		switch m.Typ {
		case config.MsgBalances:
			var newBalances []string
			for _, c := range m.Payload.(sdk.Coins) {
				newBalances = append(newBalances, c.String())
			}
			balances.Set(newBalances)
		case config.MsgChainOfTrust:
			vcs := m.Payload.([]vcTypes.VerifiableCredential)
			data, _ := json.Marshal(vcs)
			balancesChainOfTrust.Set(string(data))
		case config.MsgPublicVCs:
			var credentialIDs []string
			for _, c := range m.Payload.([]vcTypes.VerifiableCredential) {
				credentialIDs = append(credentialIDs, c.Id)
			}
			credentials.Set(credentialIDs)
		case config.MsgPublicVCData:
			vcs := m.Payload.(vcTypes.VerifiableCredential)
			data, _ := json.MarshalIndent(vcs, "", " ")
			credentialData.Set(string(data))
		}
	}

}

func balancesClick(iID widget.ListItemID) {
	v, _ := balances.GetValue(iID)
	appCfg.RuntimeMsgs.TokenWalletIn <- config.NewAppMsg(config.MsgChainOfTrust, v)
}

func credentialsSelected(iID widget.ListItemID) {
	v, _ := credentials.GetValue(iID)
	log.Debugln("credential selected", v)
	appCfg.RuntimeMsgs.TokenWalletIn <- config.NewAppMsg(config.MsgPublicVCData, v)
}

// This get executed every time the text input field get executed
func executeCmd() {
	val, _ := userCommand.Get()
	log.WithFields(log.Fields{"command": val}).Infoln("user command received")
	// parse the command
	if val == "add" {

	}
	// reset the command
	userCommand.Set("")
}
