//revive:disable:package-comments

package handler

import (
	"net/http"
)

// Handlers is a map of service handlers along with their service prefixes
type Handlers map[string]http.Handler
