package main

import (
  "os"
  "fmt"
  "net/http"
  "database/sql"
  "app/jsonHttpHandler"
  "github.com/renra/go-errtrace/errtrace"
)

type Logger struct {
}

func (l *Logger) LogWithSeverity(data map[string]string, severity int) {
  fmt.Println(fmt.Sprintf("[%d] %v", severity, data))
}

type Config struct {
}

func (ci *Config) Get(key string) (interface{}, *errtrace.Error) {
  return key, nil;
}

func (ci *Config) GetP(key string) interface{} {
  return key;
}

func (ci *Config) GetString(key string) (string, *errtrace.Error) {
  return key, nil;
}

func (ci *Config) GetStringP(key string) string {
  return key;
}

func (ci *Config) GetInt(key string) (int, *errtrace.Error) {
  return 4, nil;
}

func (ci *Config) GetIntP(key string) int {
  return 4;
}

func (ci *Config) GetFloat(key string) (float64, *errtrace.Error) {
  return 3.14, nil;
}

func (ci *Config) GetFloatP(key string) float64 {
  return 3.14;
}

func (ci *Config) GetBool(key string) (bool, *errtrace.Error) {
  return true, nil;
}

func (ci *Config) GetBoolP(key string) bool {
  return true;
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

func (g *Globals) LogErrorWithTrace(e *errtrace.Error) {
  g.Logger().LogWithSeverity(map[string]string{"msg": e.Error(), "trace": e.StringStack()}, 0)
}

func (g *Globals) DB(name string) *sql.DB {
  conn, _ := sql.Open("postgres", "whatever")
  return conn
}

func (g *Globals) Clients() map[string]interface{} {
  return map[string]interface{}{}
}

func statusHandler (g jsonHttpHandler.Globals) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    g.Log("I'm inside a handler")
    g.Log(fmt.Sprintf("Here are the clients: %v", g.Clients()))
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
