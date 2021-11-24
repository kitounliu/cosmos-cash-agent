package chain

import (
	"context"
	"encoding/base64"
	"github.com/allinbits/cosmos-cash-agent/pkg/model"
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

// DIDAddVerification ad a verification to a did document
func (cc *ChainClient) DIDAddVerification(xPubKey model.X25519ECDHKWPub, relationships ...string) {
	vmID := cc.did.NewVerificationMethodID(xPubKey.Kid)
	// now convert the base64 encoding
	rawPubKey, err := base64.StdEncoding.DecodeString(xPubKey.X)
	if err != nil {
		log.Fatalln(err)
	}
	verificationKeyAgreement := didTypes.NewVerification(
		didTypes.NewVerificationMethod(
			vmID,
			cc.did,
			didTypes.NewPublicKeyMultibase(
				rawPubKey,
				didTypes.DIDVMethodTypeX25519KeyAgreementKey2019),
		),
		[]string{didTypes.KeyAgreement},
		nil,
	)
	cc.BroadcastTx(didTypes.NewMsgAddVerification(cc.did.String(), verificationKeyAgreement, cc.acc.String()))
}
