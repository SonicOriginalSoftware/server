package env

import (
	"fmt"
	"os"
	"strings"
)

// Lookup a pre-formatted env variable or return a default value
func Lookup(prefix, attribute, defaultValue string) (value string) {
	isSet := false
	envName := fmt.Sprintf("%v_SERVE_%v", strings.ToUpper(prefix), attribute)
	if value, isSet = os.LookupEnv(envName); !isSet {
		value = defaultValue
	}
	return
}

// Protocol returns the protocol of the prefixed env variable
func Protocol(prefix, defaultProtocol string) (address string) {
	return Lookup(prefix, "PROTOCOL", defaultProtocol)
}

// Address returns the address of the prefixed env variable
func Address(prefix, defaultAddress string) (address string) {
	return Lookup(prefix, "ADDRESS", defaultAddress)
}

// Port returns the port of the prefixed env variable
func Port(prefix, defaultPort string) (address string) {
	return Lookup(prefix, "PORT", defaultPort)
}
