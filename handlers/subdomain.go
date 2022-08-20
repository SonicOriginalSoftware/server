//revive:disable:package-comments

package handlers

import (
	"net/http"
)

// SubdomainHandler defines an extension of the http.Handler interface
type SubdomainHandler interface {
	http.Handler
	Prefix() string
	Address() string
}
