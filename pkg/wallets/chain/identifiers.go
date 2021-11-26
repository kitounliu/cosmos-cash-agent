package chain

import (
	"context"
	didTypes "github.com/allinbits/cosmos-cash/v2/x/did/types"
	log "github.com/sirupsen/logrus"
)

// DIDDoc retrieve a did document for a
func (cc *ChainClient) DIDDoc(didID string) didTypes.DidDocument {
	client := didTypes.NewQueryClient(cc.ctx)
	res, err := client.DidDocument(context.Background(), &didTypes.QueryDidDocumentRequest{Id: didID})
	if err != nil {
		log.Fatalln("error requesting balance", err)
	}
	log.Infoln("did document for", didID, "is", res.GetDidDocument())
	return res.GetDidDocument()
}

// DIDAddVerification add verification to a a DID document
func (cc *ChainClient) DIDAddVerification(VMIDFragment string, pubKey []byte, vmType didTypes.VerificationMaterialType, relationships ... string) {
	vmID := cc.did.NewVerificationMethodID(VMIDFragment)
	v := didTypes.NewVerification(
		didTypes.NewVerificationMethod(
			vmID,
			cc.did,
			didTypes.NewPublicKeyMultibase(
				pubKey,
				vmType),
		),
		relationships,
		nil,
	)
	cc.BroadcastTx(didTypes.NewMsgAddVerification(cc.did.String(), v, cc.acc.String()))
}
