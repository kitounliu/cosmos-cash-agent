package chain

import (
	"github.com/allinbits/cosmos-cash-agent/pkg/config"
	didTypes "github.com/allinbits/cosmos-cash/v2/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type ChainClient struct {
	ctx client.Context
}

func Client(cfg config.EdgeConfigSchema, password string) ChainClient {
	log.Infoln("initializing client")
	chainData, _ := config.GetAppData("chain")
	keyringData, exists := config.GetAppData("keyring")

	initClientCtx := client.Context{}.
		WithHomeDir(chainData).
		WithChainID(cfg.ChainID).
		WithKeyringDir(keyringData)

	if !exists {
		//initWallet(initClientCtx, cfg, password)
	}
	log.Infoln("opening existing client for", cfg.ControllerName)

	//WithCodec(encodingConfig.Marshaler).
	//	WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
	//	WithTxConfig(encodingConfig.TxConfig).
	//	WithLegacyAmino(encodingConfig.Amino).
	//	WithInput(os.Stdin).
	//	WithAccountRetriever(types.AccountRetriever{}).
	//	WithBroadcastMode(flags.BroadcastBlock).
	//	WithHomeDir(app.DefaultNodeHome(appName)).
	//	WithViper("")
	return ChainClient{
		ctx: initClientCtx,
	}
}

// Init performs the client initialization
// 1. creates account keypair
// 2. get some coins from the faucet
// 3. creates the did document
func initWallet(ctx client.Context, cfg config.EdgeConfigSchema, password string) {
	log.Infoln("create existing chain client for", cfg.ControllerName)
	ki, m, err := ctx.Keyring.NewMnemonic(cfg.ControllerName, keyring.English, sdk.GetConfig().GetFullBIP44Path(), password, hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	log.Infoln("!WARNING! mnemonic for the account is ", m)
	log.Infoln("account address is ", ki.GetAddress())
	// didvim c
	did := didTypes.NewChainDID(ctx.ChainID, cfg.ControllerDidID)
	// verification method id
	vmID := did.NewVerificationMethodID(sdk.MustBech32ifyAddressBytes(
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		ki.GetAddress(),
	))

	verification := didTypes.NewVerification(
		didTypes.NewVerificationMethod(
			vmID,
			did,
			didTypes.NewPublicKeyMultibase(ki.GetPubKey().Bytes(), didTypes.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
		),
		[]string{didTypes.Authentication},
		nil,
	)

	// add verification
	msg := didTypes.NewMsgAddVerification(
		did.String(),
		verification,
		ki.GetAddress().String(),
	)

	fs := pflag.FlagSet{}


	if err := tx.GenerateOrBroadcastTxCLI(ctx, &fs, msg); err != nil {
		log.Error(err)
	}

	log.Infoln("token wallet initialization completed")



}

func (cc *ChainClient) Run(state *config.State, hub *config.MsgHub) {

}
