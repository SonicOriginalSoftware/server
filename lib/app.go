package lib

import (
	"fmt"
	"net/http"
	"os"
	"path"

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
		fmt.Fprintln(os.Stdout, "App service registered!")
	}

	if _, isSet := os.LookupEnv("NO_SERVE_API"); !isSet {
		fmt.Fprintln(os.Stdout, "Registering API service...")
		api_route.Register()
		fmt.Fprintln(os.Stdout, "API service registered!")
	}

	if _, isSet := os.LookupEnv("NO_SERVE_AUTH"); !isSet {
		fmt.Fprintln(os.Stdout, "Registering Auth service...")
		auth_route.Register()
		fmt.Fprintln(os.Stdout, "Auth service registered!")
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

	executablePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get working directory of executable!")
	}
	workingDirectory := path.Dir(executablePath)

	certPath, isSet := os.LookupEnv("CERT_PATH")
	if !isSet {
		certPath = fmt.Sprintf("%v/cert.pem", workingDirectory)
	}

	keyPath, isSet := os.LookupEnv("KEY_PATH")
	if !isSet {
		keyPath = fmt.Sprintf("%v/key.pem", workingDirectory)
	}

	app := &App{
		address:  fmt.Sprintf(":%v", port),
		certPath: fmt.Sprintf("%v", certPath),
		keyPath:  fmt.Sprintf("%v", keyPath),
	}
	app.setup()
	return app
}
