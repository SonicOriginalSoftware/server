//revive:disable:package-comments

package lib

import (
	"crypto/tls"
	"os"
)

// Config - the web application configuration
type Config struct {
	Address      string
	Port         string
	Certificates []tls.Certificate
}

// NewConfig returns a Config instance of an App with sane defaults
func NewConfig(certs []tls.Certificate) (config *Config) {
	var address, port string

	isSet := false

	if address, isSet = os.LookupEnv("ADDRESS"); !isSet {
		address = ""
	}

	if port, isSet = os.LookupEnv("PORT"); !isSet {
		port = "4430"
	}

	config = &Config{
		Address:      address,
		Port:         port,
		Certificates: certs,
	}

	return
}
