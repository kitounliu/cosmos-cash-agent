package ssi

import (
	log "github.com/sirupsen/logrus"
)

// LocalNotifier handles message events
type LocalNotifier struct {
}

// Notify handlers all incoming message events.
func (n LocalNotifier) Notify(topic string, message []byte) error {
	log.Infoln("local notification:", topic, message)
	return nil
}
