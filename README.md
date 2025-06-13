# server

Use any `http.Handler` implementation of your choice. Register the handler for a given path using
the `server.RegisterHandler` function.

Check out these premade handlers to get you going!

- [http server](https://git.sonicoriginal.software/server-routes-app)
- [git server](https://git.sonicoriginal.software/server-routes-git)
- [grpc server](https://git.sonicoriginal.software/server-routes-grpc)
- [graphql endpoint](https://git.sonicoriginal.software/server-routes-graphql)

# Running and Stopping

To stop the server send the program an `SIGINT` or `SIGKILL` - either natively through another
process or with a `Ctrl-C`.

Stopping the server is not "graceful" - it does not await any open connections and will likely
attempt to close down IO/handles without consideration of consumers waiting their results.

# HTTP/2 and Certs

Deploy production-grade cert files (backed by a trusted CA) to the same machine as the server
binary.

Load them in and pass them to the `server.Run` function.

**NOTE** Still in process of testing.

# Consuming

This is a _very_ high-level server library. Using it requires the import and use of a single `Run`
function:

```go
import (
  "context"
  "crypto/tls"
  "os"

  "git.sonicoriginal.software/server/v2"
	"git.sonicoriginal.software/server/v2/logging"
)

const (
	portEnvKey = "APP_PORT"
  enableHeartBeat = true
)

var (
	certs               []tls.Certificate
	mux                 = http.NewServeMux()
	ctx, cancelFunction = context.WithCancel(context.Background())
)

func main() {
	defer cancelFunction()

  // TODO Load your cert and key or skip and just use
  // cert, err := tls.X509KeyPair(cert, key)
  // if err != nil {
  //   // Handle a certificate server failure for your app here
  // }

	if enableHeartBeat {
		server.RegisterHeartBeat(mux, nil, nil)
	}

	// For examples of how to register other Handlers, see server.RegisterHeartBeat

	address, serverErrorChannel := server.Run(ctx, &certs, mux, portEnvKey)
	logging.TextLogger.InfoContext(ctx, "Serving", slog.String("address", address))

	serverError := <-serverErrorChannel

	if serverError.Close != nil {
		logging.TextLogger.ErrorContext(
			ctx,
			"Error closing server",
			slog.String("error", serverError.Close.Error()),
		)
	}

	if contextError := serverError.Context; contextError != nil {
		if errors.Is(contextError, server.ErrContextCancelled) {
			logging.TextLogger.ErrorContext(
				ctx,
				"Server failed unexpectedly",
				slog.String("error", contextError.Error()),
			)
		}
	}
}
```
