# JSON HTTP Handler

Provides hopefully reasonable defaults for JSON HTTP APIs:

* adds `Content-Type: application/json` to all responses including 404 Not Found
* Recovery from panics with a 500 Internal Server Error and an empty JSON response
* Logging of incoming requests (`Log(string)` and `LogErrorWithTrace(string)` methods are required as well as `Logger() Logger`)
* Passes a struct with pseudo-global variables (db connections, environment variables etc.) to all http handlers (You need to define the `Config() Config` method, `Clients() map[string]interface{}` and `DB(string) *sql.DB`). For more info see the [go-pseudoglobals](https://github.com/renra/go-pseudoglobals) project.

## Usage

See tests for a full example

You can use [manners](https://github.com/braintree/manners) instead of `http.Server`, add signal handling and so on still relying on `jsonHttpHandler` as the low-level tool.


