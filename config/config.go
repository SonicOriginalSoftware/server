//revive:disable:package-comments

package config

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

// New returns a Config instance of an App with sane defaults
func New(certs []tls.Certificate) (config *Config, err error) {
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
