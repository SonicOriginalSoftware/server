package env

import (
	"fmt"
	"os"
	"strings"
)

func lookup(prefix, attribute string) (value string) {
	isSet := false
	envName := fmt.Sprintf("%v_SERVE_%v", strings.ToUpper(prefix), attribute)
	if value, isSet = os.LookupEnv(envName); !isSet {
		value = ""
	}
	return
}

// Address returns the address the Handler will service
func Address(prefix string) (address string) {
	return lookup(prefix, "ADDRESS")
}

// Port returns the port the Handler will service
func Port(prefix string) (address string) {
	return lookup(prefix, "PORT")
}
