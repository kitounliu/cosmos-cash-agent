package model

import (
	"fmt"
	"time"
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
	DID     string        `json:"did"`
	Address string        `json:"address"`
	Name    string        `json:"name"`
	Texts   []TextMessage `json:"texts"`
}

func NewContact(name, didID string) Contact {
	return Contact{
		DID:   didID,
		Name:  name,
		Texts: make([]TextMessage, 0),
	}
}
