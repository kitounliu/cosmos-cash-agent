package config


// MsgHub contains channels used by the components to send messages to each others
type MsgHub struct {
	Notification   chan AppMsg
	AgentWalletIn  chan AppMsg
	AgentWalletOut chan AppMsg
	TokenWalletIn  chan AppMsg
	TokenWalletOut chan AppMsg
}

func NewMsgHub() *MsgHub {
	return &MsgHub{
		Notification:   make(chan AppMsg, 4096),
		AgentWalletIn:  make(chan AppMsg, 4096),
		AgentWalletOut: make(chan AppMsg, 4096),
		TokenWalletIn:  make(chan AppMsg, 4096),
		TokenWalletOut: make(chan AppMsg, 4096),
	}
}

const (
	// MsgBalances from chain to ui with a list of update balances
	MsgBalances = iota
	// MsgDidDoc from chain to ui with a did document (at startup)
	MsgDidDoc
	// MsgChainOfTrust returns the list of verifiable credentials for a DENOM
	MsgChainOfTrust
	// MsgPublicVCs returns the list of verifiable credentials
	MsgPublicVCs
	// MsgPublicVCData returns the details of a verifiable credential
	MsgPublicVCData
	// MsgMarketplaces used for marketplace listing
	MsgMarketplaces
	// MsgMarketplaceData used for details of a marketplace
	MsgMarketplaceData
	// MsgVCs returns the list of verifiable credentials
	MsgVCs
	// MsgVCData returns the details of a verifiable credential
	MsgVCData

)

// AppMsg are messages that are exchanged within the app
type AppMsg struct {
	Typ int
	Payload interface{}
}

func NewAppMsg(typ int, payload interface{}) AppMsg{
	return AppMsg{
		Typ:     typ,
		Payload: payload,
	}
}