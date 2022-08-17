//revive:disable:package-comments

package local

import (
	"fmt"
)

// Host is a const representing a proper path for hosting locally
const Host = "localhost"

// Path returns a localhost domain formatted properly with an optional port
func Path(port string) (path string) {
	path = Host
	if port != "" {
		path = fmt.Sprintf("%v:%v", path, port)
	}

	path = fmt.Sprintf("%v/", path)
	return
}
