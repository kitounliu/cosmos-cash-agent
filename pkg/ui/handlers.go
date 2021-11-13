package ui

import (
	"encoding/json"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/allinbits/cosmos-cash-agent/pkg/helpers"
	"github.com/allinbits/cosmos-cash-agent/pkg/model"
	vcTypes "github.com/allinbits/cosmos-cash/v2/x/verifiable-credential/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	log "github.com/sirupsen/logrus"
)

// this section contains all the databindings to the ui
var (
	appCfg *config.EdgeConfigSchema
	state  *State
	// balances
	balances             = binding.NewStringList()
	balancesChainOfTrust = binding.NewString()
	// Credentials tab
	// TODO need to have separate stuff for public and private credentials
	credentials    = binding.NewStringList()
	credentialData = binding.NewString()
	// Messages tab
	contacts    = binding.NewStringList()
	userCommand = binding.NewString()
	messages    = binding.NewStringList()
	// contactData = binding.NewString()
	// Marketplace tab
	marketplaces    = binding.NewStringList()
	marketplaceData = binding.NewString()
	// logs
	logData = binding.NewString()
)

// dispatcher this reads notifications and updates the
// data binding
func dispatcher(in chan config.AppMsg) {
	state = NewState()
	// first load the statate
	// write the state on file
	statePath, exists := config.GetAppData("state.json")
	if !exists {
		helpers.WriteJson(statePath, state)
	}
	// now load the state
	helpers.LoadJson(statePath, state)



	for {
		m := <-in
		switch m.Typ {
		case config.MsgSaveState:
			log.Debugln("saving state to file")
			helpers.WriteJson(statePath, state)
		case config.MsgBalances:
			// populate the list of balances
			var newBalances []string
			for _, c := range m.Payload.(sdk.Coins) {
				newBalances = append(newBalances, c.String())
			}
			balances.Set(newBalances)
		case config.MsgChainOfTrust:
			// populate the chain of trust for a denom
			vcs := m.Payload.([]vcTypes.VerifiableCredential)
			data, _ := json.MarshalIndent(vcs, "", " ")
			balancesChainOfTrust.Set(string(data))
		case config.MsgPublicVCs:
			// populate public credentials
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
		case config.MsgMarketplaces:
			var mks []string
			for _, c := range m.Payload.([]vcTypes.VerifiableCredential) {
				mks = append(mks, c.GetId())
			}
			marketplaces.Set(mks)
		case config.MsgContactAdded:
			contact := m.Payload.(model.Contact)
			state.Contacts[contact.Name] = contact
			// updae the model
			contacts.Append(contact.Name)
			// request state save
			appCfg.RuntimeMsgs.Notification <- config.NewAppMsg(config.MsgSaveState, nil)
		case config.MsgTextReceived:
			tm := m.Payload.(model.TextMessage)
			contact, _ := state.Contacts[tm.From]
			contact.Texts = append(contact.Texts, tm)
			// save the state
			appCfg.RuntimeMsgs.Notification <- config.NewAppMsg(config.MsgSaveState, nil)

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

// marketplacesSelected gets triggered when an item is selected in the credential list
func marketplacesSelected(iID widget.ListItemID) {
	v, _ := marketplaces.GetValue(iID)
	log.Debugln("marketplace selected", v)
	appCfg.RuntimeMsgs.TokenWalletIn <- config.NewAppMsg(config.MsgMarketplaceData, v)
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
	if val == "spend" {
		// TRANSFER TOKENS TO THE CONTACT
	}

	// get the current selected contact
	name, _ := contacts.GetValue(state.SelectedContact)
	contact := state.Contacts[name]
	// add the message from to the texts
	contact.Texts = append(contact.Texts, model.NewTextMessage(appCfg.ControllerName, val))
	// save the state
	appCfg.RuntimeMsgs.Notification <- config.NewAppMsg(config.MsgSaveState, nil)
	// reset the command
	userCommand.Set("")
}
