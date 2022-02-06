## WIP multi-plexed server

### Currently written in `go`

# Subdomain Serving

The server, by default, serves the subdomains `api` and `app`, and `auth`.

It does not serve on the top-level domain.

## `app`

For file/web-app serving, the `app` subdomain routing is appropriate.

This will serve content, by default, from a `public` folder located in whatever the current directory context is. It can be configured by setting/`export`ing a `SERVE_PATH` `env` variable.

Requests to the `/` root path will respond automatically with the `/index.html` file.

## `api`

**TODO**

In the future, this would handle requests for API-type operations:

- CRUD operations on a database
- Applying an operation on the system using some kind of tool (provisioning storage space, looking up file system/configuration information, etc.)

## `auth`

**TODO**

In the future this would be a middleware for `OAuth/2` type systems.

It is doubtful this will ever offer an `OAuth/2` implementation of its own and instead rely on interfacing with existing federated providers (Apple, Google, GitHub, etc.).

This could act as either or both an endpoint for your service to make requests to a federated provider on behalf of your service or as a callback listener for a federated provider that requires a URL for sending `auth` requests back to (Apple...).

## Additional Subdomains

Support is not yet implemented for more subdomains or extending the existing subdomains.

This might get added some time in the future by making the `mux` functionality a library that can be `go get`'ed into an underlying server framework (that also is not yet implemented but may be in the future).

# Building

Run `make`

This will build a `static`ally compiled, PIE executable, along with associated demo certs so the server can serve over `HTTP/2` by default.

# Installing

The binary can live wherever it will be exposed to your `$PATH`

# Containerization

The included `Dockerfile` will build an image that runs the binary from a `scratch` image; i.e. you can run the binary without loading it on top of _any_ Linux-based distro. *Frickin' sweet*.

# Running and Stopping

Run the output of `make`.

To stop the server send the program an `SIGINT` or `SIGKILL` - either natively through another process or with a `Ctrl-C`.

Stopping the server is not "graceful" - it does not await any open connections and will likely attempt to close down IO/handles without consideration of consumers waiting their results.

# HTTP/2 and Certs

Please don't use the generated cert files for production-grade servers. Deploy production-grade cert files (backed by a trusted CA) to the same machine as the server binary and set the path to those files using/exporting `CERT_PATH` and `KEY_PATH` `env` variables.
