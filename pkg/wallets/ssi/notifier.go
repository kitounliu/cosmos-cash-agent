package ssi

import (
	"fmt"
)

// LocalNotifier handles message events
type LocalNotifier struct {
}

// Notify handlers all incoming message events.
func (n LocalNotifier) Notify(topic string, message []byte) error {
	var allErrs error

	fmt.Println(topic, message)
	return allErrs
}
