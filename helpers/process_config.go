package helpers

import (
	"os"
)

// Reads environmental vars. if they're empty use default value
func GetEnvVar(e string, def interface{}) interface{} {

	var env interface{}
	env = os.Getenv(e)

	if env == "" {
		env = def
	}

	return env
}
