package config

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func NewEdgeConfigSchema(controllerName string) EdgeConfigSchema {
	return EdgeConfigSchema{
		ControllerName:  controllerName,
		NodeURL:         "https://cosmos-cash.app.beta.starport.cloud",
		ChainID:         "cosmoscash-testnet",
		CloudAgentWsURL: "https://ws.router.cosmos-cash.app.beta.starport.cloud",
		ControllerDidID: uuid.New().String(),
		// runtime
		RuntimeMsgs: NewMsgHub(),
	}
}

type EdgeConfigSchema struct {
	ControllerName  string `json:"controller_name"`
	ControllerDidID string `json:"controller_did"`
	NodeURL         string `json:"node_url"`
	ChainID         string `json:"chain_id"`
	CloudAgentWsURL string `json:"cloud_agent_ws_url"`

	// Runtime
	RuntimeState *State  `json:"-"`
	RuntimeMsgs  *MsgHub `json:"-"`
}

func GetAppData(subPath ...string) (string, bool) {
	v := []string{"data"}
	v = append(v, subPath...)
	return GetAppConfig(v...)
}

func GetAppConfig(subPath ...string) (cfgPath string, exists bool) {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	cfgPath = path.Join(cfgPath, "cosmos-cash-agent")
	for _, sp := range subPath {
		cfgPath = path.Join(cfgPath, sp)
	}
	exists = true
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		exists = false
	}
	return
}
