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
	app_route.Setup()
	api_route.Setup()
	auth_route.Setup()
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
	// FIXME Yank address and certificate paths from env variables
	app := &App{
		address:  ":8080",
		certPath: "./cert.pem",
		keyPath:  "./key.pem",
	}
	app.setup()
	return app
}
