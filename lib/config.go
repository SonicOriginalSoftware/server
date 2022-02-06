package lib

import (
	"fmt"
	"log"
	"os"
	"path"
)

// Config - the web application configuration
type Config struct {
	Address  string
	Port     string
	CertPath string
	KeyPath  string
}

// NewConfig returns an instance of an App with sane defaults
func NewConfig(outlog *log.Logger, errlog *log.Logger) (app *Config, err error) {
	var address, port, executablePath, certPath, keyPath string

	isSet := false

	if address, isSet = os.LookupEnv("ADDRESS"); !isSet {
		address = ""
	}

	if port, _ = os.LookupEnv("PORT"); !isSet {
		port = "4430"
	}

	if executablePath, err = os.Executable(); err != nil {
		errlog.Printf("Could not get working directory of executable!")
		return
	}

	workingDirectory := path.Dir(executablePath)

	if certPath, isSet = os.LookupEnv("CERT_PATH"); !isSet {
		certPath = fmt.Sprintf("%v/cert.pem", workingDirectory)
	}

	if keyPath, isSet = os.LookupEnv("KEY_PATH"); !isSet {
		keyPath = fmt.Sprintf("%v/key.pem", workingDirectory)
	}

	app = &Config{
		Address:  address,
		Port:     port,
		CertPath: certPath,
		KeyPath:  keyPath,
	}

	return
}
