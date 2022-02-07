package env

import (
	"fmt"
	"os"
	"strings"
)

func lookup(prefix, attribute, def string) (value string) {
	isSet := false
	envName := fmt.Sprintf("%v_SERVE_%v", strings.ToUpper(prefix), attribute)
	if value, isSet = os.LookupEnv(envName); !isSet {
		value = def
	}
	return
}

// Protocol returns the protocol of the prefixed env variable
func Protocol(prefix, def string) (address string) {
	return lookup(prefix, "PROTOCOL", def)
}

// Address returns the address of the prefixed env variable
func Address(prefix, def string) (address string) {
	return lookup(prefix, "ADDRESS", def)
}

// Port returns the port of the prefixed env variable
func Port(prefix, def string) (address string) {
	return lookup(prefix, "PORT", def)
}
