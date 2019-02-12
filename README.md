# JSON HTTP Handler

Provides hopefully reasonable defaults for JSON HTTP APIs:

* adds `Content-Type: application/json` to all responses including 404 Not Found
* Recovery from panics with a 500 Internal Server Error and an empty JSON response
* Logging of incoming requests (`Log(string)` and `LogErrorWithTrace(string)` methods are required as well as Logger())
* Passes a struct with pseudo-global variables (db connections, environment variables etc.) to all http handlers (You need to define the Config() method)

## Usage

```go
package main

import (
  "os"
  "fmt"
  "net/http"
  "github.com/renra/go-json-http-handler/jsonHttpHandler"
)

// First you need to define some structs to implement the desired behaviour.
//  You can use other libraries that already implement it for you.
type Logger struct {
}

func (l *Logger) LogWithSeverity(data map[string]string, severity int) {
  fmt.Println(fmt.Sprintf("[%d] %v", severity, data))
}

type Config struct {
}

func (c *Config) Get(key string) interface {} {
  return key
}

func (c *Config) GetString(key string) string {
  return key
}

// Globals is a struct which provides access to utilities and a logger, it is eventually passed down to handlers
type Globals struct {
}

func (g *Globals) Config() jsonHttpHandler.Config {
  return &Config{}
}

func (g *Globals) Logger() jsonHttpHandler.Logger {
  return &Logger{}
}

func (g *Globals) Log(msg string) {
  g.Logger().LogWithSeverity(map[string]string{"msg": msg}, 1)
}

func (g *Globals) LogErrorWithTrace(msg string, trace string) {
  g.Logger().LogWithSeverity(map[string]string{"msg": msg, "trace": trace}, 0)
}

// Your own handlers receive globals
func statusHandler (g jsonHttpHandler.Globals) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    g.Log("I'm inside a handler")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{\"status\":\"ok\"}")
  }
}

func main() {
  globals := &Globals{}

  handler := jsonHttpHandler.New(globals, map[string]jsonHttpHandler.GlobalsReceivingHandlerFunc{
    "/status": statusHandler,
  })

  port := os.Getenv("PORT")

  s := &http.Server{
    Addr: fmt.Sprintf("0.0.0.0:%s", port),
    Handler: handler,
  }

  globals.Log(fmt.Sprintf("About to listen on port %s", port))
  s.ListenAndServe()
}
```

You can use [manners](https://github.com/braintree/manners) instead of `http.Server`, add signal handling and so on still relying on `jsonHttpHandler` as the low-level tool.


