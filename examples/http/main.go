package main

import (
  "os"
  "fmt"
  "net/http"
  "app/jsonHttpHandler"
)

type Globals struct {
}

func (g *Globals) Log(msg string) {
  fmt.Println(fmt.Sprintf("[Logger] %s", msg))
}

func (g *Globals) LogErrorWithTrace(msg string, trace string) {
  fmt.Println(fmt.Sprintf("[Logger] msg=%s trace=%s", msg, trace))
}

func statusHandler (g jsonHttpHandler.Globals) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    g.Log("Status")
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
