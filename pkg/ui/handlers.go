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
			data, _ := json.MarshalIndent(vcs, "", " ")
			balancesChainOfTrust.Set(string(data))
		case config.MsgPublicVCs:
			var credentialIDs []string
			for _, c := range m.Payload.([]vcTypes.VerifiableCredential) {
				credentialIDs = append(credentialIDs, c.Id)
			}
			credentials.Set(credentialIDs)
		case config.MsgPublicVCData:
			if m.Payload == nil {
				credentialData.Set(string("No data"))
				continue
			}
			vcs := m.Payload.(vcTypes.VerifiableCredential)
			data, _ := json.MarshalIndent(vcs, "", " ")
			credentialData.Set(string(data))
		}
	}

}

// balancesSelected gets triggered when an item is selected in the balance list
func balancesSelected(iID widget.ListItemID) {
	v, _ := balances.GetValue(iID)
	log.Debugln("token selected", v)
	appCfg.RuntimeMsgs.TokenWalletIn <- config.NewAppMsg(config.MsgChainOfTrust, v)
}

// credentialsSelected gets triggered when an item is selected in the credential list
func credentialsSelected(iID widget.ListItemID) {
	v, _ := credentials.GetValue(iID)
	log.Debugln("credential selected", v)
	appCfg.RuntimeMsgs.TokenWalletIn <- config.NewAppMsg(config.MsgPublicVCData, v)
}

// contactSelected gets triggered when an item is selected in the contact list
func contactSelected(iID widget.ListItemID) {
	// TODO what should happen when a contact is selected?
	// the messages should be sent to that contact
	// the payments should be sent to that contact
}

// executeCmd get executed every time the text input field get executed
func executeCmd() {
	val, _ := userCommand.Get()
	log.WithFields(log.Fields{"command": val}).Infoln("user command received")
	// parse the command
	if val == "add" {

	}
	// reset the command
	userCommand.Set("")
}
