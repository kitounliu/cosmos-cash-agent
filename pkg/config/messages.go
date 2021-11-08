package config


// MsgHub contains channels used by the components to send messages to each others
type MsgHub struct {
	Notification   chan string
	AgentWalletIn  chan string
	AgentWalletOut chan string
	TokenWalletIn  chan string
	TokenWalletOut chan string
}

func NewMsgHub() *MsgHub {
	return &MsgHub{
		Notification:   make(chan string, 4096),
		AgentWalletIn:  make(chan string, 4096),
		AgentWalletOut: make(chan string, 4096),
		TokenWalletIn:  make(chan string, 4096),
		TokenWalletOut: make(chan string, 4096),
	}
}
