//revive:disable:package-comments

package lib

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
