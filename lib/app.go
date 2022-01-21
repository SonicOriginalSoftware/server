package lib

import (
	"fmt"
	"net/http"
	"os"

	api_route "api-server/routes/api"
	app_route "api-server/routes/app"
	auth_route "api-server/routes/auth"
)

// App - the web application
type App struct {
	address  string
	certPath string
	keyPath  string
}

func (app *App) setup() {
	if _, isSet := os.LookupEnv("NO_SERVE_APP"); !isSet {
		fmt.Fprintln(os.Stdout, "Registering App service...")
		app_route.Register()
	}

	if _, isSet := os.LookupEnv("NO_SERVE_API"); !isSet {
		fmt.Fprintln(os.Stdout, "Registering API service...")
		api_route.Register()
	}

	if _, isSet := os.LookupEnv("NO_SERVE_AUTH"); !isSet {
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

	app := &App{
		address:  fmt.Sprintf(":%v", port),
		certPath: fmt.Sprintf("%v", certPath),
		keyPath:  fmt.Sprintf("%v", keyPath),
	}
	app.setup()
	return app
}
