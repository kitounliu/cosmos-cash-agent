package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	didTypes "github.com/allinbits/cosmos-cash/v2/x/did/types"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
)

// TextMessage a text message exchanged between a contact and the user
type TextMessage struct {
	// Channel is used to identify which chat the message belongs to
	Channel     string    `json:"-"`
	From        string    `json:"fr"`
	Content     string    `json:"cn"`
	ProcessedAt time.Time `json:"at"`
}

func (tm TextMessage) String() string {
	return fmt.Sprintf("%s\n%10s | %s", tm.ProcessedAt.Format(time.RFC822), tm.From, tm.Content)
}

func NewTextMessageWithTime(channel, from, content string, processedAt time.Time) TextMessage {
	return TextMessage{
		Channel:     channel,
		From:        from,
		Content:     content,
		ProcessedAt: processedAt,
	}
}

func NewTextMessage(channel, from, content string) TextMessage {
	return NewTextMessageWithTime(channel, from, content, time.Now())
}

// Contact represent SSI contact
type Contact struct {
	DID        string                 `json:"did"`
	Address    string                 `json:"address"`
	Name       string                 `json:"name"`
	Connection didexchange.Connection `json:"connection"`
	Texts      []TextMessage          `json:"texts"`
}

func NewContact(name, didID string, connection didexchange.Connection) Contact {
	return Contact{
		DID:        didID,
		Name:       name,
		Connection: connection,
		Texts:      make([]TextMessage, 0),
	}
}

// deal with public keys

type AriesPubKey interface {
	PubKeyBytes() []byte
	KeyID() string
	VerificationMaterialType() didTypes.VerificationMaterialType
	DIDRelationships() []string
}

type X25519 struct {
	*wallet.KeyPair
}

func (x X25519) VerificationMaterialType() didTypes.VerificationMaterialType {
	return didTypes.DIDVMethodTypeX25519KeyAgreementKey2019
}

func (x X25519) DIDRelationships() []string {
	return []string{didTypes.KeyAgreement}
}

func (x X25519) KeyID() string {
	return x.KeyPair.KeyID
}

func (x X25519) PubKeyBytes() []byte {
	// decode the pub key base64
	b, err := base64.RawURLEncoding.DecodeString(x.PublicKey)
	if err != nil {
		log.Fatalln("cannot decode X25519 pub key base64 string", err)
	}
	// parse the pub key structure
	var x25519 struct {
		Kid   string `json:"kid"`
		X     string `json:"x"`
		Curve string `json:"curve"`
		Type  string `json:"type"`
	}
	if err := json.Unmarshal(b, &x25519); err != nil {
		log.Fatalln("cannot parse X25519 pub key data", err)
	}
	// export the pub key bytes
	pk, err := base64.StdEncoding.DecodeString(x25519.X)
	if err != nil {
		log.Fatalln("cannot decode X25519 X component", err)
	}
	return pk

}

type ED25519 struct {
	*wallet.KeyPair
}

func (x ED25519) VerificationMaterialType() didTypes.VerificationMaterialType {
	return didTypes.DIDVMethodTypeEd25519VerificationKey2018
}

func (x ED25519) DIDRelationships() []string {
	return []string{didTypes.AssertionMethod}
}

func (x ED25519) KeyID() string {
	return x.KeyPair.KeyID
}

func (x ED25519) PubKeyBytes() []byte {
	// decode the pub key base64
	b, err := base64.RawURLEncoding.DecodeString(x.PublicKey)
	if err != nil {
		log.Fatalln("cannot decode ED25519 pub key base64 string", err)
	}
	return b
}

// UTLITY MESSAGES

// ChargedEnvelope this is used to send messages that should be sent back to the sender
type ChargedEnvelope struct {
	DataIn   interface{}
	Callback func(message string)
}

// CREDENTIALS

// AccountCredentialSubject represent a subject for a blockchain account
type AccountCredentialSubject struct {
	ID      string
	Address string
	Name    string
}

// ChainAccountCredential creates a verifiable credential for a blockchain account
func ChainAccountCredential(chainID, address, did, name string) *verifiable.Credential {
	return &verifiable.Credential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
		},
		ID: didTypes.NewBlockchainAccountID(chainID, address).EncodeToString(),
		Types: []string{
			"VerifiableCredential",
			// "CosmosAccountAddressCredential",
		},
		Subject: AccountCredentialSubject{
			ID:      did,
			Name:    name,
			Address: address,
		},
		Issuer: verifiable.Issuer{
			ID: did,
		},
		Issued: util.NewTime(time.Now()),
	}
}

// ChainAccountCredentialRaw creates a verifiable credential vor a blockchain account
// and returns it as a json.RawMessage ready to be used in the aries wallet context
func ChainAccountCredentialRaw(chainID, address, did, name string) json.RawMessage {
	vc := ChainAccountCredential(chainID, address, did, name)
	rawVC, _ := json.Marshal(vc)
	return json.RawMessage(string(rawVC))
}
