package helpers

import (
	"os"
	"strings"
)

func Env(name, defaultValue string) (val string) {
	val = os.Getenv(name)
	if strings.TrimSpace(val) == "" {
		val = defaultValue
	}
	return
}