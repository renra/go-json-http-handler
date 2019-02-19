package jsonHttpHandler

import (
  "fmt"
  "net/http"
  "database/sql"
  "github.com/go-errors/errors"
)

type Config interface {
  Get(string) interface{}
  GetString(string) string
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
      err := errors.Wrap(r, 2)
      h.globals.LogErrorWithTrace(err.Err.Error(), err.ErrorStack())

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

