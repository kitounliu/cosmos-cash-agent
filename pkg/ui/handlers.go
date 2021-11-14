package ui

import (
	"encoding/json"
	"fyne.io/fyne/v2/widget"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	vcTypes "github.com/allinbits/cosmos-cash/v2/x/verifiable-credential/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	log "github.com/sirupsen/logrus"
	"strings"
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
		case config.MsgHandleInvitation:
			newContact := m.Payload.(*didexchange.Connection)
			contacts.Append(newContact.TheirLabel + " " + newContact.ConnectionID)
		}
	}

}

// balancesSelected gets triggered when an item is selected in the balance list
func balancesSelected(iID widget.ListItemID) {
	v, _ := balances.GetValue(iID)
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

	s := strings.Split(val, " ")

	switch s[0] {
	case "ssi":
	case "s":
		switch s[1] {
		case "invitation":
		case "i":
			switch s[2] {
			case "handle":
			case "h":
				ns := strings.Join(s, " ")
				log.Infoln("command handler", ns)
				appCfg.RuntimeMsgs.AgentWalletIn <- config.NewAppMsg(config.MsgHandleInvitation, s[3])
			}
		}
	case "chain":
		//hub.Notification <- "chain"
	}

	// reset the command
	userCommand.Set("")
}
