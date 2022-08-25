# server

This module will eventually be hosted on a `git` server outside of GitHub. For the time being there is a replace directive to point the module at the repository hosted on GitHub.

# Subdomain Serving

Run the server by passing in a map of service prefixes (as `string`) and their respective `http.Handler`s to the `Run` method along with any desired certs. Leaving an empty array of certs will serve over `http` rather than `https` (this also means no `HTTP/2`).

It does not serve on an empty top-level domain (e.g. `NOT_IMPLEMENTED` response for `localhost:4430` vs. `app.localhost:4430`)

## Additional Subdomains

Use any custom `http.Handler` of your choice. Just attach the handler to a `Prefix` string when passing through to the `Run` function.

Check out these premade handlers to get you going!

- [web-app](https://github.com/SonicOriginalSoftware/server-routes-app)
- [git server](https://github.com/SonicOriginalSoftware/server-routes-git)
- [grpc server](https://github.com/SonicOriginalSoftware/server-routes-grpc)
- [graphql endpoint](https://github.com/SonicOriginalSoftware/server-routes-graphql)

# Running and Stopping

To stop the server send the program an `SIGINT` or `SIGKILL` - either natively through another process or with a `Ctrl-C`.

Stopping the server is not "graceful" - it does not await any open connections and will likely attempt to close down IO/handles without consideration of consumers waiting their results.

# HTTP/2 and Certs

Deploy production-grade cert files (backed by a trusted CA) to the same machine as the server binary.

# Consuming

This is a _very_ high-level server library. Using it requires the import and use of a single `Run` function:

```go
import (
	lib "git.nathanblair.rocks/server"
	"git.nathanblair.rocks/server/handler"
)

func main() {
	subdomains := handler.Handlers{
    // TODO Import your desired subdomain services here
    // This is a map so you would set a string key
    // associated with the handler
    // and this key identifies the top-level domain
    // where that service can be reached
    // e.g. if importing the 'git' handler, use
    // git.Prefix: git.New(),
  }

  var cert, key []byte

  // TODO Load your cert and key or skip and just use
  // var certs []tls.Certificate

	cert, err := tls.X509KeyPair(cert, key)
	if err != nil {
    // Handle a certificate server failure for your app here
	}

	certs := []tls.Certificate{cert}
	ctx, cancelCtx := context.WithCancel(context.Background())

  if exitCode := lib.Run(ctx, subdomains, certs); exitCode != 0 {
    // Handle the server error according to your needs
	}

  cancelCtx()
}
```
