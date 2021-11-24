package model

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
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

// X25519ECDHKWPub pub key for key agreement format
type X25519ECDHKWPub struct {
	Kid   string `json:"kid"`
	X     string `json:"x"`
	Curve string `json:"curve"`
	Type  string `json:"type"`
}

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

func ChainAccountCredential(address, did, name string) *verifiable.Credential {
	return &verifiable.Credential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
		},
		ID: fmt.Sprint("cosmos:account:", address),
		Types: []string{
			"VerifiableCredential",
			"CosmosAccountAddressCredential",
		},
		Subject: AccountCredentialSubject{
			ID:      did,
			Name:    name,
			Address: address,
		},
		Issuer: verifiable.Issuer{
			ID: did,
		},
		Issued:  util.NewTime(time.Now()),
	}
}

func ChainAccountCredentialRaw(address, did, name string) json.RawMessage {
	vc := ChainAccountCredential(address, did, name)
	rawVC, _ := json.Marshal(vc)
	return json.RawMessage(string(rawVC))
}

