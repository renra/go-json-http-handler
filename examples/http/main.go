package main

import (
  "os"
  "fmt"
  "net/http"
  "app/jsonHttpHandler"
)

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
