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

// Address returns the address the Handler will service
func Address(prefix, def string) (address string) {
	return lookup(prefix, "ADDRESS", def)
}

// Port returns the port the Handler will service
func Port(prefix, def string) (address string) {
	return lookup(prefix, "PORT", def)
}
