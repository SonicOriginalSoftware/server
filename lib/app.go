package lib

import (
	"fmt"
	"net/http"
	"os"

	api_route "pwa-server/routes/api"
	app_route "pwa-server/routes/app"
	auth_route "pwa-server/routes/auth"
)

// App - the web application
type App struct {
	address  string
	certPath string
	keyPath  string
}

func (app *App) setup() {
	_, isSet := os.LookupEnv("NO_SERVE_APP")
	if !isSet {
		fmt.Fprintln(os.Stdout, "Registering App service...")
		app_route.Register()
	}

	_, isSet = os.LookupEnv("NO_SERVE_API")
	if !isSet {
		fmt.Fprintln(os.Stdout, "Registering API service...")
		api_route.Register()
	}

	_, isSet = os.LookupEnv("NO_SERVE_AUTH")
	if !isSet {
		fmt.Fprintln(os.Stdout, "Registering Auth service...")
		auth_route.Register()
	}
}

// Serve the App
func (app *App) Serve() {
	error := http.ListenAndServeTLS(app.address, app.certPath, app.keyPath, nil)
	if error != nil {
		fmt.Fprint(os.Stderr, error)
	}
}

// NewApp returns an instance of an App with sane defaults
func NewApp() *App {
	port, isSet := os.LookupEnv("PORT")
	if !isSet {
		port = "8080"
	}

	certPath, isSet := os.LookupEnv("CERT_PATH")
	if !isSet {
		certPath = "./cert.pem"
	}

	keyPath, isSet := os.LookupEnv("KEY_PATH")
	if !isSet {
		keyPath = "./key.pem"
	}

	// FIXME Yank address and certificate paths from env variables
	app := &App{
		address:  fmt.Sprintf(":%v", port),
		certPath: fmt.Sprintf("%v", certPath),
		keyPath:  fmt.Sprintf("%v", keyPath),
	}
	app.setup()
	return app
}
