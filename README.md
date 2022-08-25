# server

# Subdomain Serving

Run the server by passing in a map of service prefixes (as `string`) and their respective `http.Handler`s to the `Run` method along with any desired certs (leaving an empty array of certs will serve over `http` rather than `https` - note that this also means no `HTTP/2`).

It does not serve on the top-level domain.

## Additional Subdomains

Use any custom `http.Handler` of your choice. Just attach the handler to a `Prefix` string when passing through to the `Run` function.

# Running and Stopping

To stop the server send the program an `SIGINT` or `SIGKILL` - either natively through another process or with a `Ctrl-C`.

Stopping the server is not "graceful" - it does not await any open connections and will likely attempt to close down IO/handles without consideration of consumers waiting their results.

# HTTP/2 and Certs

Deploy production-grade cert files (backed by a trusted CA) to the same machine as the server binary.
