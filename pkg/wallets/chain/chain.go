package chain

import (
	"bytes"
	"context"
	"fmt"
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	"github.com/allinbits/cosmos-cash-agent/pkg/helpers"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"net/http"
	"time"

	"github.com/allinbits/cosmos-cash/v2/app"
	didTypes "github.com/allinbits/cosmos-cash/v2/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type ChainClient struct {
	ctx client.Context
	fct tx.Factory
	qry grpc.ClientConn
}

func Client(cfg config.EdgeConfigSchema, password string) *ChainClient {
	log.Infoln("initializing client")
	chainData, _ := config.GetAppData("chain")

	kr := keyring.NewInMemory()

	armorPath, exists := config.GetAppData("account_armor.json")
	if !exists {
		_, addr, err := kr.NewMnemonic(cfg.ControllerName, keyring.English, sdk.GetConfig().GetFullBIP44Path(), "", hd.Secp256k1)
		if err != nil {
			log.Fatalln("error creating new key", err)
		}
		log.WithFields(log.Fields{
			"mnemonic": addr,
		}).Infoln("created new key for", cfg.ControllerName)
		armorData, err := kr.ExportPrivKeyArmor(cfg.ControllerName, password)
		if err != nil {
			log.Fatalln("error exporting private key", err)
		}

		helpers.WriteJson(armorPath, map[string]interface{}{
			"address": addr,
			"armor":   armorData,
		})
		log.Infoln("exported armored private key to", armorPath)
	} else {
		var accountData map[string]string
		helpers.LoadJson(armorPath, accountData)
		err := kr.ImportPrivKey(cfg.ControllerName, accountData["armor"], password)
		if err != nil {
			log.Fatalln("error loading private key", err)
		}
		log.Infoln("private key loaded from", armorPath)
	}
	// now get the account
	ki, err := kr.Key(cfg.ControllerName)
	if err != nil {
		log.Fatalln("cannot load stored key by uid", err)
	}
	// RPC client for transactions
	netCli, err := client.NewClientFromNode(cfg.NodeURI)
	if err != nil {
		log.Fatalln("error connecting to the node", err)
	}
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithChainID(cfg.ChainID).
		WithKeyring(kr).
		WithHomeDir(chainData).
		WithNodeURI(cfg.NodeURI).
		WithFromName(cfg.ControllerName).
		WithFromAddress(ki.GetAddress()).
		WithSkipConfirmation(true).
		WithClient(netCli)
	//WithLegacyAmino(encodingConfig.Amino).
	//WithInput(os.Stdin).

	log.Infoln("opening existing client for", cfg.ControllerName)
	pf := pflag.NewFlagSet("default", pflag.PanicOnError)
	factory := tx.NewFactoryCLI(initClientCtx, pf).
		WithChainID(initClientCtx.ChainID).
		WithAccountRetriever(initClientCtx.AccountRetriever).
		WithKeybase(initClientCtx.Keyring)

	cc := ChainClient{
		ctx: initClientCtx,
		fct: factory,
	}

	if !exists {
		// get some tokens here
		if !cc.Balance(ki.GetAddress().String()).IsPositive() {
			callFaucet(cfg.FaucetURL, ki.GetAddress().String())

			for i := 0; i < 3; i++ {
				time.Sleep(6 * time.Second)
				if cc.Balance(ki.GetAddress().String()).IsPositive() {
					break
					log.Infoln("got a positive balance")
				}
			}
		}
		msg := initDIDDoc(cfg.ChainID, cfg.ControllerDidID, cfg.CloudAgentPublicURL, ki)
		cc.BroadcastTx(msg)
	}

	return &cc
}

func (cc *ChainClient) BroadcastTx(msgs ...sdk.Msg) {
	log.Infoln("broadcasting messages")
	if err := tx.GenerateOrBroadcastTxWithFactory(cc.ctx, cc.fct, msgs...); err != nil {
		log.Fatalln("failed tx", err)
	}

}

func (cc *ChainClient) Balance(address string) *sdk.Coin {
	bankClient := banktypes.NewQueryClient(cc.ctx)
	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: address, Denom: "cash"},
	)
	if err != nil {
		log.Fatalln("error requesting balance", err)
	}
	log.Infoln("balance for", address, "is", bankRes.GetBalance())
	return bankRes.GetBalance()
}

// Init performs the client initialization
// 1. creates account keypair
// 2. get some coins from the faucet
// 3. creates the did document
func initDIDDoc(chainID, didID, agentURL string, ki keyring.Info) sdk.Msg {
	log.Println("initializing new did document", didID)
	did := didTypes.NewChainDID(chainID, didID)
	// verification method id
	vmID := did.NewVerificationMethodID(ki.GetAddress().String())

	verification := didTypes.NewVerification(
		didTypes.NewVerificationMethod(
			vmID,
			did,
			didTypes.NewPublicKeyMultibase(ki.GetPubKey().Bytes(), didTypes.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
		),
		[]string{didTypes.Authentication},
		nil,
	)

	service := didTypes.NewService(
		didID+"-agent",
		"DIDCommMessaging",
		agentURL,
	)

	return didTypes.NewMsgCreateDidDocument(
		did.String(),
		didTypes.Verifications{verification},
		didTypes.Services{service},
		ki.GetAddress().String(),
	)

	//// add verification
	//return didTypes.NewMsgAddVerification(
	//	did.String(),
	//	verification,
	//	ki.GetAddress().String(),
	//)
}

func callFaucet(faucetURL, address string) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	payload := fmt.Sprintf(`{"address": "%s"}`, address)
	rsp, err := netClient.Post(faucetURL, "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		log.Fatalln("error requesting tokens from the faucet")
	}
	log.Debugf("faucet response: %v", rsp)
	if rsp.StatusCode != http.StatusOK {
		log.Fatalln("error requesting tokens from the faucet", rsp.Status, rsp)
	}

}

func (cc *ChainClient) Run(state *config.State, hub *config.MsgHub) {

}
