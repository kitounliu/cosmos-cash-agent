package ui

import (
	"github.com/sirupsen/logrus"
)

// Hook to send logs to a PostgreSQL database
type Hook struct {
	LogLevels []logrus.Level
}

func NewHook(levels ...logrus.Level) *Hook {
	return &Hook{LogLevels: levels}
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	for _, l := range hook.LogLevels {
		if l == entry.Level {
			logData.Append(entry.Message)
			break
		}
	}
	return nil
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
