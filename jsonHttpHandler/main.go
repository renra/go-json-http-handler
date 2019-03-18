package jsonHttpHandler

import (
  "fmt"
  "net/http"
  "database/sql"
  "github.com/renra/go-errtrace/errtrace"
)

type Config interface {
  Get(string) (interface{}, *errtrace.Error)
  GetP(string) interface{}
  GetString(string) (string, *errtrace.Error)
  GetStringP(string) string
  GetInt(string) (int, *errtrace.Error)
  GetIntP(string) int
  GetFloat(string) (float64, *errtrace.Error)
  GetFloatP(string) float64
  GetBool(string) (bool, *errtrace.Error)
  GetBoolP(string) bool
}

type Logger interface {
  LogWithSeverity(map[string]string, int)
}

type Globals interface {
  Config() Config
  Logger() Logger
  Log(string)
  LogErrorWithTrace(string, string)
  DB(string) *sql.DB
  Clients() map[string]interface{}
}

type GlobalsReceivingHandlerFunc func(Globals) http.HandlerFunc

type JsonHttpHandler struct {
  globals Globals
  handlers map[string]GlobalsReceivingHandlerFunc
}

func (h JsonHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  h.globals.Log(fmt.Sprintf("%s", r.URL))

  defer func() {
    if r := recover(); r != nil {
      err := errtrace.Wrap(r)
      h.globals.LogErrorWithTrace(err.Error(), err.StringStack())

      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{}")
    }
  }()

  handler := h.handlers[r.URL.Path]

  if handler != nil {
    handler(h.globals)(w, r)
  } else {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "{}")
  }
}

func New(g Globals, handlers map[string]GlobalsReceivingHandlerFunc) JsonHttpHandler {
  return JsonHttpHandler{globals: g, handlers: handlers}
}

