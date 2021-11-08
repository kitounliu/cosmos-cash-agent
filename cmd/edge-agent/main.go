package main

import (
	"bufio"
	"fmt"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/allinbits/cosmos-cash-agent/pkg/helpers"
	"github.com/allinbits/cosmos-cash-agent/pkg/ui"
	"github.com/allinbits/cosmos-cash-agent/pkg/wallets/chain"
	"github.com/allinbits/cosmos-cash-agent/pkg/wallets/ssi"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("_private/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		log.SetOutput(os.Stderr)
		log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	// check if the config file exists
	// if not ask for the account name
	// create the config
	// init the web ui

	cfg := setup()

	pwd := "a_password"

	// aries wallet creation
	// https://github.com/hyperledger/aries-framework-go/blob/main/docs/vc_wallet.md
	agent := ssi.Agent(cfg.ControllerName, pwd)
	go agent.Run(cfg.RuntimeState, cfg.RuntimeMsgs)

	// cosmos-sdk keystore
	// https://github.com/cosmos/cosmos-sdk/blob/master/client/keys/add.go
	wallet := chain.Client(cfg, pwd)
	go wallet.Run(cfg.RuntimeState, cfg.RuntimeMsgs)

	// render the app
	ui.Render(cfg.RuntimeState, cfg.RuntimeMsgs)
}

// setup creates the app config folder
func setup() (cfg config.EdgeConfigSchema) {
	cfgDir, exists := config.GetAppConfig()
	if !exists {
		if err := os.MkdirAll(cfgDir, 0700); err != nil {
			panic(fmt.Sprintln("cannot create the config directory", err))
		}
	}
	dataDir, exists := config.GetAppData()
	if !exists {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			panic(fmt.Sprintln("cannot create the data directory", err))
		}
	}
	// load wallet config
	agentCfg, exists := config.GetAppConfig("edget-agent.json")
	if !exists {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Hi there! looks like you are new here, welcome!")
		fmt.Println("let's begin with the formality, what is your name? ")
		fmt.Print("> ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			panic("too bad that you don't have a name :(")
		}
		fmt.Println("Great", name, "! strap-on and let's go!")
		cfg = config.NewEdgeConfigSchema(name)
		helpers.WriteJson(agentCfg, cfg)
	} else {
		helpers.LoadJson(agentCfg, &cfg)
	}
	// load app state
	cfg.RuntimeState = config.NewState()
	appState, exists := config.GetAppData("state.json")
	if !exists {
		helpers.WriteJson(appState, cfg.RuntimeState)
	} else {
		helpers.LoadJson(appState, cfg.RuntimeState)
	}
	cfg.RuntimeMsgs = config.NewMsgHub()
	return
}
