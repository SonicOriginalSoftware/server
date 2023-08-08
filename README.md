# server

Use any `http.Handler` implementation of your choice. Register the handler for a given path using the `server.RegisterHandler` function.

Check out these premade handlers to get you going!

- [http server](https://git.sonicoriginal.software/server-routes-app)
- [git server](https://git.sonicoriginal.software/server-routes-git)
- [grpc server](https://git.sonicoriginal.software/server-routes-grpc)
- [graphql endpoint](https://git.sonicoriginal.software/server-routes-graphql)

# Running and Stopping

To stop the server send the program an `SIGINT` or `SIGKILL` - either natively through another process or with a `Ctrl-C`.

Stopping the server is not "graceful" - it does not await any open connections and will likely attempt to close down IO/handles without consideration of consumers waiting their results.

# HTTP/2 and Certs

Deploy production-grade cert files (backed by a trusted CA) to the same machine as the server binary.

Load them in and pass them to the `server.Run` function.

**NOTE** Still in process of testing.

# Consuming

This is a _very_ high-level server library. Using it requires the import and use of a single `Run` function:

```go
import (
  "context"
  "crypto/tls"

  "git.sonicoriginal.software/server.git/v2"
  "git.sonicoriginal.software/server.git/v2/handler"
)

func main() {
  const portEnvKey = "APP_PORT"
  os.Setenv(portEnvKey, "4430") // Default

  // TODO Import your desired handlers and register them here
  // e.g. if importing the 'app' handler, use
  // _ = app.New()

  var certs []tls.Certificate

  // TODO Load your cert and key or skip and just use
  // cert, err := tls.X509KeyPair(cert, key)
  // if err != nil {
  //   // Handle a certificate server failure for your app here
  // }
  // certs = []tls.Certificate{cert}

  ctx, cancelContext := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)

  // Do other stuff while your server runs

  // Wait for your server to close (through a signal or internal error)
	serverError := <-serverErrorChannel
	if serverError.Close != nil {
    // Handle closing server error
	}

	contextError := serverError.Context.Error()

	if serverError.Context.Error() != nil {
    // Handle server failing unexpectedly
	}

  cancelCtx()
}
```
